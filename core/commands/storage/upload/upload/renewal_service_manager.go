package upload

import (
	"fmt"
	"sync"

	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	logging "github.com/ipfs/go-log/v2"
)

var (
	renewalServiceLog = logging.Logger("renewal-service-manager")
)

// RenewalServiceManager manages the lifecycle of the auto-renewal service
type RenewalServiceManager struct {
	service *AutoRenewalService
	mu      sync.RWMutex
	started bool
}

// NewRenewalServiceManager creates a new renewal service manager
func NewRenewalServiceManager() *RenewalServiceManager {
	return &RenewalServiceManager{
		started: false,
	}
}

// StartService starts the auto-renewal service if it's not already running
func (rsm *RenewalServiceManager) StartService(ctxParams *uh.ContextParams) error {
	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	if rsm.started && rsm.service != nil {
		renewalServiceLog.Debug("Auto-renewal service is already running")
		return nil
	}

	renewalServiceLog.Info("Starting auto-renewal service...")

	// Create new service instance
	rsm.service = NewAutoRenewalService(ctxParams)

	// Start the service
	err := rsm.service.Start()
	if err != nil {
		renewalServiceLog.Errorf("Failed to start auto-renewal service: %v", err)
		return err
	}

	rsm.started = true
	renewalServiceLog.Info("Auto-renewal service started successfully")

	return nil
}

// StopService stops the auto-renewal service if it's running
func (rsm *RenewalServiceManager) StopService() error {
	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	if !rsm.started || rsm.service == nil {
		renewalServiceLog.Debug("Auto-renewal service is not running")
		return nil
	}

	renewalServiceLog.Info("Stopping auto-renewal service...")

	err := rsm.service.Stop()
	if err != nil {
		renewalServiceLog.Errorf("Failed to stop auto-renewal service: %v", err)
		return err
	}

	rsm.service = nil
	rsm.started = false
	renewalServiceLog.Info("Auto-renewal service stopped successfully")

	return nil
}

// IsRunning returns whether the auto-renewal service is currently running
func (rsm *RenewalServiceManager) IsRunning() bool {
	rsm.mu.RLock()
	defer rsm.mu.RUnlock()
	return rsm.started && rsm.service != nil && rsm.service.IsRunning()
}

// GetService returns the current auto-renewal service instance
func (rsm *RenewalServiceManager) GetService() *AutoRenewalService {
	rsm.mu.RLock()
	defer rsm.mu.RUnlock()
	return rsm.service
}

// RestartService restarts the auto-renewal service
func (rsm *RenewalServiceManager) RestartService(ctxParams *uh.ContextParams) error {
	renewalServiceLog.Info("Restarting auto-renewal service...")

	// Stop the current service
	err := rsm.StopService()
	if err != nil {
		return err
	}

	// Start a new service
	return rsm.StartService(ctxParams)
}

// Global renewal service manager instance
var globalRenewalServiceManager *RenewalServiceManager
var renewalServiceManagerMu sync.Mutex

// GetGlobalRenewalServiceManager returns the global renewal service manager instance
func GetGlobalRenewalServiceManager() *RenewalServiceManager {
	renewalServiceManagerMu.Lock()
	defer renewalServiceManagerMu.Unlock()

	if globalRenewalServiceManager == nil {
		globalRenewalServiceManager = NewRenewalServiceManager()
	}

	return globalRenewalServiceManager
}

// InitializeRenewalService initializes and starts the global renewal service
func InitializeRenewalService(ctxParams *uh.ContextParams) error {
	manager := GetGlobalRenewalServiceManager()
	return manager.StartService(ctxParams)
}

// ShutdownRenewalService stops the global renewal service
func ShutdownRenewalService() error {
	manager := GetGlobalRenewalServiceManager()
	return manager.StopService()
}

// IsRenewalServiceRunning checks if the global renewal service is running
func IsRenewalServiceRunning() bool {
	manager := GetGlobalRenewalServiceManager()
	return manager.IsRunning()
}

// GetRenewalService returns the current renewal service instance
func GetRenewalService() *AutoRenewalService {
	manager := GetGlobalRenewalServiceManager()
	return manager.GetService()
}

// RenewalServiceStatus represents the status of the renewal service
type RenewalServiceStatus struct {
	Running       bool   `json:"running"`
	CheckInterval string `json:"check_interval"`
	Message       string `json:"message"`
}

// GetRenewalServiceStatus returns the current status of the renewal service
func GetRenewalServiceStatus() *RenewalServiceStatus {
	manager := GetGlobalRenewalServiceManager()
	service := manager.GetService()

	status := &RenewalServiceStatus{
		Running: manager.IsRunning(),
	}

	if service != nil {
		status.CheckInterval = service.checkInterval.String()
		if status.Running {
			status.Message = "Auto-renewal service is running normally"
		} else {
			status.Message = "Auto-renewal service is stopped"
		}
	} else {
		status.Message = "Auto-renewal service is not initialized"
	}

	return status
}

// EnableAutoRenewalForFile enables auto-renewal for a specific file
func EnableAutoRenewalForFile(ctxParams *uh.ContextParams, fileHash string) error {
	service := GetRenewalService()
	if service == nil {
		return fmt.Errorf("auto-renewal service is not running")
	}

	return service.EnableAutoRenewal(fileHash)
}

// DisableAutoRenewalForFile disables auto-renewal for a specific file
func DisableAutoRenewalForFile(ctxParams *uh.ContextParams, fileHash string) error {
	service := GetRenewalService()
	if service == nil {
		return fmt.Errorf("auto-renewal service is not running")
	}

	return service.DisableAutoRenewal(fileHash)
}

// GetAutoRenewalStatusForFile gets the auto-renewal status for a specific file
func GetAutoRenewalStatusForFile(ctxParams *uh.ContextParams, fileHash string) (*AutoRenewalConfig, error) {
	service := GetRenewalService()
	if service == nil {
		return nil, fmt.Errorf("auto-renewal service is not running")
	}

	return service.GetAutoRenewalStatus(fileHash)
}
