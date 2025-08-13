package renewal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	"github.com/bittorrent/go-btfs/utils"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p/core/peer"
)

var log = logging.Logger("renew")

const (
	renewDurationOptionName = "duration"
	renewTokenOptionName    = "renew-token"
	renewPriceOptionName    = "renew-price"
)

// RenewRequest represents a file renewal request
type RenewRequest struct {
	CID         string    `json:"cid"`
	Token       string    `json:"token"`
	Price       uint64    `json:"price"`
	Duration    int       `json:"duration"`
	ShardSize   int64     `json:"shard_size"`
	ShardId     string    `json:"shard_id"`
	SpId        string    `json:"sp_id"`
	RenterID    peer.ID   `json:"renter_id"`
	OriginalEnd time.Time `json:"original_end"`
	NewEnd      time.Time `json:"new_end"`
	TotalCost   int64     `json:"total_cost"`
}

// RenewResponse represents the response of a renewal operation
type RenewResponse struct {
	Success       bool      `json:"success"`
	SessionID     string    `json:"session_id"`
	FileHash      string    `json:"file_hash"`
	NewExpiration time.Time `json:"new_expiration"`
	TotalCost     int64     `json:"total_cost"`
	Message       string    `json:"message"`
}

// StorageRenewCmd implements the storage renew command
var StorageRenewCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Renew storage duration for uploaded files.",
		ShortDescription: `
This command allows users to extend the storage duration of previously uploaded files
without re-uploading the content. The renewal extends the storage contract with
existing storage providers.

Examples:
    # Renew a file for 30 days
    $ btfs storage upload renew <file-hash> --duration 30

    # Renew with specific token and price
    $ btfs storage upload renew <file-hash> --duration 60 --token WBTT --price 1000
`,
	},
	Subcommands: map[string]*cmds.Command{
		"enable":  StorageRenewEnableCmd,
		"disable": StorageRenewDiableCmd,
		"status":  StorageRenewStatusCmd,
		"list":    StorageRenewListCmd,
		"service": StorageRenewServiceCmd,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("file-hash", true, false, "Hash of the file to renew."),
	},
	Options: []cmds.Option{
		cmds.IntOption(renewDurationOptionName, "d", "Renewal duration in days.").WithDefault(30),
		cmds.StringOption(renewTokenOptionName, "rt", "Token type for payment (WBTT/TRX/USDD/USDT).").WithDefault("WBTT"),
		cmds.Int64Option(renewPriceOptionName, "rp", "Max price per GiB per day in µBTT."),
	},
	RunTimeout: 10 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		if !nd.IsOnline {
			return errors.New("node must be online to renew storage")
		}

		err = utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		// Extract parameters
		cid := req.Arguments[0]
		duration, _ := req.Options[renewDurationOptionName].(int)
		tokenStr, _ := req.Options[renewTokenOptionName].(string)
		priceOpt, hasPriceOpt := req.Options[renewPriceOptionName].(int64)

		// Validate parameters
		if duration <= 0 {
			return errors.New("duration must be positive")
		}

		// Get context parameters
		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}

		// check if the cid enabled autorenew
		info, err := getRenewalInfo(ctxParams, cid)
		if err != nil {
			return err
		}
		if info != nil {
			return fmt.Errorf("file %s is already auto-renewed and cannot be renewed manually", cid)
		}
		// Get token address
		_, exists := tokencfg.MpTokenAddr[tokenStr]
		if !exists {
			return fmt.Errorf("invalid token type: %s", tokenStr)
		}

		// Get current price if not specified
		var price int64
		if hasPriceOpt {
			price = priceOpt
		} else {
			price, _, err = uh.GetPriceAndMinStorageLength(ctxParams)
			if err != nil {
				return fmt.Errorf("failed to get current price: %v", err)
			}
		}

		// Create renewal request
		renewReq := &RenewRequest{
			CID:         cid,
			Duration:    duration,
			Token:       tokenStr,
			Price:       uint64(price),
			RenterID:    ctxParams.N.Identity,
			OriginalEnd: time.Now().Add(time.Duration(duration) * 24 * time.Hour), // This should be calculated from existing contract
			NewEnd:      time.Now().Add(time.Duration(duration) * 24 * time.Hour),
		}

		// Execute renewal
		renewResp, err := executeRenewal(ctxParams, renewReq)
		if err != nil {
			return fmt.Errorf("renewal failed: %v", err)
		}

		return res.Emit(renewResp)
	},
	Type: RenewResponse{},
}

// executeRenewal performs the actual renewal operation
func executeRenewal(ctxParams *uh.ContextParams, renewReq *RenewRequest) (*RenewResponse, error) {
	// Calculate total cost
	rate, err := chain.SettleObject.OracleService.CurrentRate(common.HexToAddress(renewReq.Token))
	if err != nil {
		return nil, fmt.Errorf("failed to get token rate: %v", err)
	}

	totalCost, err := uh.TotalPay(renewReq.ShardSize, int64(renewReq.Price), renewReq.Duration, rate)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate total cost: %v", err)
	}

	renewReq.TotalCost = totalCost

	// Check available balance
	if err := checkAvailableBalance(ctxParams.Ctx, totalCost, common.HexToAddress(renewReq.Token)); err != nil {
		return nil, fmt.Errorf("insufficient balance: %v", err)
	}

	// Execute renewal with storage providers
	err = payRenewalCheque(ctxParams, renewReq, renewReq.ShardId, totalCost)
	if err != nil {
		return nil, fmt.Errorf("failed to execute renewal with providers: %v", err)
	}

	return &RenewResponse{
		Success:       true,
		FileHash:      renewReq.CID,
		NewExpiration: renewReq.NewEnd,
		TotalCost:     totalCost,
		Message:       fmt.Sprintf("File %s renewed successfully for %d days", renewReq.CID, renewReq.Duration),
	}, nil
}

// payRenewalCheque pays renewal fee directly via cheque to storage provider
func payRenewalCheque(ctxParams *uh.ContextParams, renewReq *RenewRequest, shardHash string, paymentAmount int64) error {
	// Get the original shard contract to find the storage provider
	// Get storage provider ID
	spId := renewReq.SpId
	if spId == "" {
		return fmt.Errorf("no storage provider ID found in contract")
	}

	log.Infof("Paying renewal cheque for shard %s to sp %s, amount: %d", shardHash, spId, paymentAmount)

	// Check available balance before issuing cheque
	err := checkAvailableBalance(ctxParams.Ctx, paymentAmount, common.HexToAddress(renewReq.Token))
	if err != nil {
		return fmt.Errorf("insufficient balance for renewal payment: %v", err)
	}

	// Issue cheque directly to storage provider for renewal
	err = issueRenewalCheque(ctxParams, spId, paymentAmount, common.HexToAddress(renewReq.Token), shardHash, renewReq.Duration)
	if err != nil {
		return fmt.Errorf("failed to issue renewal cheque to provider %s: %v", spId, err)
	}

	// Update shard renewal information
	// err = updateShardRenewalInfo(ctxParams, session.SsId, shardHash, shardIndex, renewReq.Duration, paymentAmount)
	// if err != nil {
	// 	log.Errorf("Failed to update shard renewal info: %v", err)
	// 	// Don't fail the payment for info update issues
	// }

	log.Infof("Successfully issued renewal cheque for shard %s", shardHash)
	return nil
}

// issueRenewalCheque issues a cheque directly to storage provider for renewal payment
func issueRenewalCheque(ctxParams *uh.ContextParams, providerID string, amount int64, token common.Address, shardHash string, duration int) error {
	log.Infof("Issuing renewal cheque to provider %s for shard %s, amount: %d, duration: %d days", providerID, shardHash, amount, duration)

	// Get settlement service
	if chain.SettleObject.SwapService == nil {
		return fmt.Errorf("settlement service not available")
	}

	// Convert amount to big.Int
	paymentAmount := big.NewInt(amount)

	// Generate a renewal contract ID for tracking
	renewalContractID := fmt.Sprintf("renewal_%s_%s_%d", shardHash, providerID, time.Now().Unix())

	// Issue cheque through settlement service
	// This directly pays the provider without creating a new contract
	// Pay method signature: Pay(ctx context.Context, peer string, amount *big.Int, contractId string, token common.Address)
	chain.SettleObject.SwapService.Settle(providerID, paymentAmount, renewalContractID, token)

	log.Infof("Successfully issued renewal cheque to provider %s for contract %s", providerID, renewalContractID)

	// Store cheque information for tracking
	err := storeRenewalChequeInfo(ctxParams, providerID, shardHash, renewalContractID, amount, duration)
	if err != nil {
		log.Errorf("Failed to store renewal cheque info: %v", err)
		// Don't fail the payment for storage issues
	}

	return nil
}

// updateShardRenewalInfo updates the shard information with renewal details
// func updateShardRenewalInfo(ctxParams *uh.ContextParams, sessionID, shardHash string, shardIndex int, duration int, amount int64) error {
// 	// Verify shard exists (we don't need to use the shard object, just verify it exists)
// 	_, err := sessions.GetUserShard(ctxParams, sessionID, shardHash, shardIndex)
// 	if err != nil {
// 		return fmt.Errorf("failed to get shard info: %v", err)
// 	}
//
// 	// Update shard with renewal information
// 	// This extends the storage period without creating new contracts
// 	renewalInfo := map[string]interface{}{
// 		"renewed_at":       time.Now().Unix(),
// 		"renewal_duration": duration,
// 		"renewal_amount":   amount,
// 		"new_expiry":       time.Now().Add(time.Duration(duration) * 24 * time.Hour).Unix(),
// 	}
//
// 	// Store renewal info in shard metadata
// 	shardKey := fmt.Sprintf("/btfs/%s/shards/%s/%s/renewal", ctxParams.N.Identity.String(), sessionID, shardHash)
// 	renewalData, err := json.Marshal(renewalInfo)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal renewal info: %v", err)
// 	}
//
// 	err = ctxParams.N.Repo.Datastore().Put(ctxParams.Ctx, datastore.NewKey(shardKey), renewalData)
// 	if err != nil {
// 		return fmt.Errorf("failed to store shard renewal info: %v", err)
// 	}
//
// 	log.Infof("Updated renewal info for shard %s in session %s", shardHash, sessionID)
// 	return nil
// }

// storeRenewalChequeInfo stores cheque information for renewal tracking
func storeRenewalChequeInfo(ctxParams *uh.ContextParams, providerID, shardHash string, contractID string, amount int64, duration int) error {
	chequeInfo := map[string]interface{}{
		"provider_id": providerID,
		"shard_hash":  shardHash,
		"amount":      amount,
		"duration":    duration,
		"issued_at":   time.Now().Unix(),
		"contract_id": contractID, // Store the renewal contract ID
	}

	chequeKey := fmt.Sprintf("/btfs/%s/renewal_cheques/%s_%s_%d",
		ctxParams.N.Identity.String(), providerID, shardHash, time.Now().Unix())

	chequeData, err := json.Marshal(chequeInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal cheque info: %v", err)
	}

	return ctxParams.N.Repo.Datastore().Put(ctxParams.Ctx, datastore.NewKey(chequeKey), chequeData)
}

// Note: We no longer use contract-based renewal. Instead, we directly issue cheques to providers.
// This simplifies the renewal process and avoids the complexity of contract renegotiation.

// AutoRenewalConfig represents auto-renewal configuration for a file
type AutoRenewalConfig struct {
	FileHash        string         `json:"file_hash"`
	SessionID       string         `json:"session_id"`
	SpId            string         `json:"sp_id"`
	ShardId         string         `json:"shard_id"`
	ShardSize       int            `json:"shard_size"`
	RenewalDuration int            `json:"renewal_duration"`
	Token           common.Address `json:"token"`
	Price           int64          `json:"price"`
	Enabled         bool           `json:"enabled"`
	CreatedAt       time.Time      `json:"created_at"`
	LastRenewalAt   *time.Time     `json:"last_renewal_at,omitempty"`
	NextRenewalAt   time.Time      `json:"next_renewal_at"`
}

// StoreAutoRenewalConfig stores auto-renewal configuration for a file
func StoreAutoRenewalConfig(ctxParams *uh.ContextParams, fileHash string, duration int, token common.Address, price int64) error {
	config := &AutoRenewalConfig{
		FileHash:        fileHash,
		RenewalDuration: duration,
		Token:           token,
		Price:           price,
		Enabled:         true,
		CreatedAt:       time.Now(),
		NextRenewalAt:   time.Now().Add(time.Duration(duration) * 24 * time.Hour),
	}

	configKey := fmt.Sprintf("/btfs/%s/autorenew/%s", ctxParams.N.Identity.String(), fileHash)

	configData, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return ctxParams.N.Repo.Datastore().Put(ctxParams.Ctx, datastore.NewKey(configKey), configData)
}

func checkAvailableBalance(ctx context.Context, amount int64, token common.Address) error {
	realAmount, err := getRealAmount(amount, token)
	if err != nil {
		return err
	}

	// token: get available balance of token.
	// AvailableBalance, err := chain.SettleObject.VaultService.AvailableBalance(ctx, token)
	AvailableBalance, err := chain.SettleObject.VaultService.AvailableBalance(ctx, token)
	if err != nil {
		return err
	}

	fmt.Printf("check,  balance=%v, realAmount=%v \n", AvailableBalance, realAmount)
	if AvailableBalance.Cmp(realAmount) < 0 {
		fmt.Println("check, err: ", vault.ErrInsufficientFunds)
		return vault.ErrInsufficientFunds
	}
	return nil
}

func getRealAmount(amount int64, token common.Address) (*big.Int, error) {
	// this is price's rate [Compatible with older versions]
	rateObj, err := chain.SettleObject.OracleService.CurrentRate(token)
	if err != nil {
		return nil, err
	}

	realAmount := big.NewInt(0).Mul(big.NewInt(amount), rateObj)
	return realAmount, nil
}
