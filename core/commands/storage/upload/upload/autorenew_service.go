package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	nodepb "github.com/bittorrent/go-btfs-common/protos/node"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/protos/metadata"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	logging "github.com/ipfs/go-log/v2"
)

const (
	userFileShard = "/btfs/%s/shards/file/%s"
)

var (
	autoRenewLog = logging.Logger("autorenew")
)

// AutoRenewalService manages automatic renewal of storage contracts
type AutoRenewalService struct {
	ctxParams     *uh.ContextParams
	ctx           context.Context
	cancel        context.CancelFunc
	ticker        *time.Ticker
	mu            sync.RWMutex
	running       bool
	checkInterval time.Duration
}

// NewAutoRenewalService creates a new auto-renewal service
func NewAutoRenewalService(ctxParams *uh.ContextParams) *AutoRenewalService {
	ctx, cancel := context.WithCancel(ctxParams.Ctx)

	return &AutoRenewalService{
		ctxParams:     ctxParams,
		ctx:           ctx,
		cancel:        cancel,
		checkInterval: 1 * time.Minute, // Check every hour
		running:       false,
	}
}

// Start begins the auto-renewal service
func (ars *AutoRenewalService) Start() error {
	ars.mu.Lock()
	defer ars.mu.Unlock()

	if ars.running {
		return fmt.Errorf("auto-renewal service is already running")
	}

	ars.ticker = time.NewTicker(ars.checkInterval)
	ars.running = true

	go ars.run()

	autoRenewLog.Info("Auto-renewal service started")
	return nil
}

// Stop stops the auto-renewal service
func (ars *AutoRenewalService) Stop() error {
	ars.mu.Lock()
	defer ars.mu.Unlock()

	if !ars.running {
		return fmt.Errorf("auto-renewal service is not running")
	}

	ars.cancel()
	if ars.ticker != nil {
		ars.ticker.Stop()
	}
	ars.running = false

	autoRenewLog.Info("Auto-renewal service stopped")
	return nil
}

// IsRunning returns whether the service is currently running
func (ars *AutoRenewalService) IsRunning() bool {
	ars.mu.RLock()
	defer ars.mu.RUnlock()
	return ars.running
}

// run is the main loop for the auto-renewal service
func (ars *AutoRenewalService) run() {
	defer func() {
		if r := recover(); r != nil {
			autoRenewLog.Errorf("Auto-renewal service panic: %v", r)
		}
	}()

	for {
		select {
		case <-ars.ctx.Done():
			autoRenewLog.Info("Auto-renewal service context cancelled")
			return
		case <-ars.ticker.C:
			ars.checkAndRenewFiles()
		}
	}
}

// checkAndRenewFiles checks for files that need renewal and processes them
func (ars *AutoRenewalService) checkAndRenewFiles() {
	autoRenewLog.Debug("Checking for files that need renewal")

	// configs, err := ars.getAutoRenewalConfigs()
	contracts, err := sessions.ListShardsContracts(ars.ctxParams.N.Repo.Datastore(), ars.ctxParams.N.Identity.String(), nodepb.ContractStat_RENTER.String())
	if err != nil {
		autoRenewLog.Errorf("Failed to get auto-renewal configs: %v", err)
		return
	}

	now := time.Now()
	renewalThreshold := 24 * time.Hour // Renew 24 hours before expiration

	for _, contract := range contracts {
		if !contract.Meta.AutoRenewal {
			continue
		}

		// Check if renewal is needed (within threshold of expiration)
		if now.Add(renewalThreshold).After(time.Unix(int64(contract.Meta.StorageEnd), 0)) {
			cid, err := ars.ctxParams.N.Repo.Datastore().Get(ars.ctx, datastore.NewKey(fmt.Sprintf(userFileShard, ars.ctxParams.N.Identity.String(), contract.Meta.ContractId)))
			if err != nil {
				autoRenewLog.Errorf("Failed to get file CID: %v", err)
				continue
			}
			autoRenewLog.Infof("Processing auto-renewal for file: %s", cid)

			err = ars.processAutoRenewal(contract)
			if err != nil {
				autoRenewLog.Errorf("Failed to auto-renew file %s: %v", cid, err)
			} else {
				autoRenewLog.Infof("Successfully auto-renewed file: %s", cid)
			}
		}
	}
}

// getAutoRenewalConfigs retrieves all auto-renewal configurations
func (ars *AutoRenewalService) getAutoRenewalConfigs() ([]AutoRenewalConfig, error) {
	prefix := fmt.Sprintf("/btfs/%s/autorenew/", ars.ctxParams.N.Identity.String())
	q := query.Query{
		Prefix: prefix,
	}

	results, err := ars.ctxParams.N.Repo.Datastore().Query(ars.ctx, q)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var configs []AutoRenewalConfig

	for result := range results.Next() {
		if result.Error != nil {
			autoRenewLog.Errorf("Error reading auto-renewal config: %v", result.Error)
			continue
		}

		var config AutoRenewalConfig
		if err := json.Unmarshal(result.Value, &config); err != nil {
			autoRenewLog.Errorf("Failed to unmarshal auto-renewal config: %v", err)
			continue
		}

		configs = append(configs, config)
	}

	return configs, nil
}

// processAutoRenewal processes the automatic renewal for a specific file
func (ars *AutoRenewalService) processAutoRenewal(contract *metadata.Contract) error {
	renewReq := &RenewRequest{
		FileHash:    contract.Meta.ContractId,
		Token:       contract.Meta.Token,
		Price:       contract.Meta.Price,
		SpId:        contract.Meta.SpId,
		RenterID:    ars.ctxParams.N.Identity,
		ShardSize:   int64(contract.Meta.ShardSize),
		OriginalEnd: time.Unix(int64(contract.Meta.StorageEnd), 0),
		NewEnd:      time.Unix(int64(contract.Meta.StorageEnd), 0).Add(time.Duration(contract.Meta.StorageEnd - contract.Meta.StorageStart)),
		Duration:    int(contract.Meta.StorageEnd-contract.Meta.StorageStart) / 86400,
	}

	// Execute renewal
	_, err := executeRenewal(ars.ctxParams, renewReq)
	if err != nil {
		return fmt.Errorf("renewal execution failed: %v", err)
	}

	return nil
}

// updateAutoRenewalConfig updates an auto-renewal configuration
func (ars *AutoRenewalService) updateAutoRenewalConfig(config *AutoRenewalConfig) error {
	configKey := fmt.Sprintf("/btfs/%s/autorenew/%s", ars.ctxParams.N.Identity.String(), config.FileHash)

	configData, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return ars.ctxParams.N.Repo.Datastore().Put(ars.ctx, datastore.NewKey(configKey), configData)
}

// GetAutoRenewalStatus returns the status of auto-renewal for a specific file
func (ars *AutoRenewalService) GetAutoRenewalStatus(fileHash string) (*AutoRenewalConfig, error) {
	configKey := fmt.Sprintf("/btfs/%s/autorenew/%s", ars.ctxParams.N.Identity.String(), fileHash)

	data, err := ars.ctxParams.N.Repo.Datastore().Get(ars.ctx, datastore.NewKey(configKey))
	if err != nil {
		if err == datastore.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	var config AutoRenewalConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// DisableAutoRenewal disables auto-renewal for a specific file
func (ars *AutoRenewalService) DisableAutoRenewal(fileHash string) error {
	config, err := ars.GetAutoRenewalStatus(fileHash)
	if err != nil {
		return err
	}

	if config == nil {
		return fmt.Errorf("no auto-renewal config found for file: %s", fileHash)
	}

	config.Enabled = false
	return ars.updateAutoRenewalConfig(config)
}

// EnableAutoRenewal enables auto-renewal for a specific file
func (ars *AutoRenewalService) EnableAutoRenewal(fileHash string) error {
	config, err := ars.GetAutoRenewalStatus(fileHash)
	if err != nil {
		return err
	}

	if config == nil {
		return fmt.Errorf("no auto-renewal config found for file: %s", fileHash)
	}

	config.Enabled = true
	return ars.updateAutoRenewalConfig(config)
}

// Global auto-renewal service instance
var globalAutoRenewalService *AutoRenewalService
var autoRenewalServiceMu sync.Mutex

// GetGlobalAutoRenewalService returns the global auto-renewal service instance
func GetGlobalAutoRenewalService() *AutoRenewalService {
	autoRenewalServiceMu.Lock()
	defer autoRenewalServiceMu.Unlock()
	return globalAutoRenewalService
}

// SetGlobalAutoRenewalService sets the global auto-renewal service instance
func SetGlobalAutoRenewalService(service *AutoRenewalService) {
	autoRenewalServiceMu.Lock()
	defer autoRenewalServiceMu.Unlock()
	globalAutoRenewalService = service
}
