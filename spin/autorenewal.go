package spin

import (
	"fmt"

	"github.com/bittorrent/go-btfs/core"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/upload"

	cmds "github.com/bittorrent/go-btfs-cmds"
	logging "github.com/ipfs/go-log/v2"
)

var autoRenewalLog = logging.Logger("spin-autorenewal")

// AutoRenewalService starts the auto-renewal service for storage files
func AutoRenewalService(node *core.IpfsNode, req *cmds.Request, env cmds.Environment) {
	fmt.Println("Initializing auto-renewal service...")

	// Check if storage client is enabled
	cfg, err := node.Repo.Config()
	if err != nil {
		autoRenewalLog.Errorf("Failed to get node config: %v", err)
		return
	}

	if !cfg.Experimental.StorageClientEnabled {
		autoRenewalLog.Debug("Storage client is disabled, skipping auto-renewal service")
		return
	}

	// Extract context parameters
	ctxParams, err := uh.ExtractContextParams(req, env)
	if err != nil {
		autoRenewalLog.Errorf("Failed to extract context parameters: %v", err)
		return
	}

	// Start the auto-renewal service in a goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				autoRenewalLog.Errorf("Auto-renewal service panic: %v", r)
			}
		}()

		err := upload.InitializeRenewalService(ctxParams)
		if err != nil {
			autoRenewalLog.Errorf("Failed to initialize auto-renewal service: %v", err)
			return
		}

		fmt.Println("Auto-renewal service started successfully")
	}()
}

// StopAutoRenewalService stops the auto-renewal service
func StopAutoRenewalService() {
	autoRenewalLog.Info("Stopping auto-renewal service...")

	err := upload.ShutdownRenewalService()
	if err != nil {
		autoRenewalLog.Errorf("Failed to stop auto-renewal service: %v", err)
		return
	}

	autoRenewalLog.Info("Auto-renewal service stopped successfully")
}

// GetAutoRenewalServiceStatus returns the status of the auto-renewal service
func GetAutoRenewalServiceStatus() *upload.RenewalServiceStatus {
	return upload.GetRenewalServiceStatus()
}

// RestartAutoRenewalService restarts the auto-renewal service
func RestartAutoRenewalService(node *core.IpfsNode, req *cmds.Request, env cmds.Environment) error {
	autoRenewalLog.Info("Restarting auto-renewal service...")

	// Stop the current service
	StopAutoRenewalService()

	// Extract context parameters
	ctxParams, err := uh.ExtractContextParams(req, env)
	if err != nil {
		return fmt.Errorf("failed to extract context parameters: %v", err)
	}

	// Restart the service
	manager := upload.GetGlobalRenewalServiceManager()
	err = manager.RestartService(ctxParams)
	if err != nil {
		return fmt.Errorf("failed to restart auto-renewal service: %v", err)
	}

	autoRenewalLog.Info("Auto-renewal service restarted successfully")
	return nil
}
