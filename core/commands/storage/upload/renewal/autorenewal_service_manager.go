package renewal

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log/v2"
)

var (
	renewalServiceLog = logging.Logger("renewal-service-manager")
)

const (
	serviceStateKey = "/btfs/renewal-service/state"
)

// PersistentServiceState represents the persistent state of the renewal service
type PersistentServiceState struct {
	Running     bool      `json:"running"`
	PID         int       `json:"pid"`
	StartTime   time.Time `json:"start_time"`
	NodeID      string    `json:"node_id"`
	LastUpdated time.Time `json:"last_updated"`
}

// SaveServiceState saves the service state to datastore
func SaveServiceState(ctxParams *uh.ContextParams, state *PersistentServiceState) error {
	state.LastUpdated = time.Now()
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return ctxParams.N.Repo.Datastore().Put(ctxParams.Ctx, datastore.NewKey(serviceStateKey), data)
}

// LoadServiceState loads the service state from datastore
func LoadServiceState(ctxParams *uh.ContextParams) (*PersistentServiceState, error) {
	data, err := ctxParams.N.Repo.Datastore().Get(ctxParams.Ctx, datastore.NewKey(serviceStateKey))
	if err != nil {
		if err == datastore.ErrNotFound {
			return &PersistentServiceState{Running: false}, nil
		}
		return nil, err
	}

	var state PersistentServiceState
	err = json.Unmarshal(data, &state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

// ClearServiceState removes the service state from datastore
func ClearServiceState(ctxParams *uh.ContextParams) error {
	return ctxParams.N.Repo.Datastore().Delete(ctxParams.Ctx, datastore.NewKey(serviceStateKey))
}

// IsProcessRunning checks if a process with the given PID is running
func IsProcessRunning(pid int) bool {
	if pid <= 0 {
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// On Unix systems, sending signal 0 checks if process exists
	err = process.Signal(os.Signal(nil))
	return err == nil
}

// ServiceManager manages the lifecycle of the auto-renewal service
type ServiceManager struct {
	service *AutoRenewalService
	mu      sync.RWMutex
	started bool
	stop    chan struct{}
}

// NewRenewalServiceManager creates a new renewal service manager
func NewRenewalServiceManager() *ServiceManager {
	return &ServiceManager{
		started: false,
	}
}

// StartService starts the auto-renewal service if it's not already running
func (rsm *ServiceManager) StartService(ctxParams *uh.ContextParams) error {
	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	// Check persistent state first
	persistentState, err := LoadServiceState(ctxParams)
	if err != nil {
		renewalServiceLog.Errorf("Failed to load service state: %v", err)
		return err
	}

	// Check if service is already running based on persistent state
	if persistentState.Running && IsProcessRunning(persistentState.PID) {
		renewalServiceLog.Debug("Auto-renewal service is already running")
		return nil
	}

	renewalServiceLog.Info("Starting auto-renewal service...")

	// Create new service instance
	rsm.service = NewAutoRenewalService(ctxParams)

	// Start the service
	err = rsm.service.Start()
	if err != nil {
		renewalServiceLog.Errorf("Failed to start auto-renewal service: %v", err)
		return err
	}

	rsm.started = true

	// Save persistent state
	newState := &PersistentServiceState{
		Running:   true,
		PID:       os.Getpid(),
		StartTime: time.Now(),
		NodeID:    ctxParams.N.Identity.String(),
	}

	err = SaveServiceState(ctxParams, newState)
	if err != nil {
		renewalServiceLog.Errorf("Failed to save service state: %v", err)
		// Don't return error here, service is still running
	}

	renewalServiceLog.Info("Auto-renewal service started successfully")

	return nil
}

// StopService stops the auto-renewal service if it's running
func (rsm *ServiceManager) StopService(ctxParams *uh.ContextParams) error {
	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	// Load persistent state
	persistentState, err := LoadServiceState(ctxParams)
	if err != nil {
		renewalServiceLog.Errorf("Failed to load service state: %v", err)
		return err
	}

	if !persistentState.Running {
		renewalServiceLog.Debug("Auto-renewal service is not running")
		return nil
	}

	renewalServiceLog.Info("Stopping auto-renewal service...")

	// Stop the service if it's running in current process
	if rsm.started && rsm.service != nil {
		err := rsm.service.Stop()
		if err != nil {
			renewalServiceLog.Errorf("Failed to stop auto-renewal service: %v", err)
		}
		rsm.service = nil
		rsm.started = false
	}

	// Clear persistent state
	err = ClearServiceState(ctxParams)
	if err != nil {
		renewalServiceLog.Errorf("Failed to clear service state: %v", err)
		return err
	}

	renewalServiceLog.Info("Auto-renewal service stopped successfully")

	rsm.stop <- struct{}{}

	return nil
}

// IsRunning returns whether the auto-renewal service is currently running
func (rsm *ServiceManager) IsRunning(ctxParams *uh.ContextParams) bool {
	rsm.mu.RLock()
	defer rsm.mu.RUnlock()

	// Check persistent state
	persistentState, err := LoadServiceState(ctxParams)
	if err != nil {
		renewalServiceLog.Errorf("Failed to load service state: %v", err)
		return false
	}

	// Verify the process is actually running
	if persistentState.Running && IsProcessRunning(persistentState.PID) {
		return true
	}

	// If persistent state says running but process is not found, clean up
	if persistentState.Running && !IsProcessRunning(persistentState.PID) {
		renewalServiceLog.Warn("Service marked as running but process not found, cleaning up state")
		ClearServiceState(ctxParams)
	}

	return false
}

// GetService returns the current auto-renewal service instance
func (rsm *ServiceManager) GetService() *AutoRenewalService {
	rsm.mu.RLock()
	defer rsm.mu.RUnlock()
	return rsm.service
}

// RestartService restarts the auto-renewal service
func (rsm *ServiceManager) RestartService(ctxParams *uh.ContextParams) error {
	renewalServiceLog.Info("Restarting auto-renewal service...")

	// Stop the current service
	err := rsm.StopService(ctxParams)
	if err != nil {
		return err
	}

	// Start a new service
	return rsm.StartService(ctxParams)
}

// Global renewal service manager instance
var globalRenewalServiceManager *ServiceManager
var renewalServiceManagerMu sync.Mutex

// GetGlobalRenewalServiceManager returns the global renewal service manager instance
func GetGlobalRenewalServiceManager() *ServiceManager {
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
func ShutdownRenewalService(ctxParams *uh.ContextParams) error {
	manager := GetGlobalRenewalServiceManager()
	return manager.StopService(ctxParams)
}

// IsRenewalServiceRunning checks if the global renewal service is running
func IsRenewalServiceRunning(ctxParams *uh.ContextParams) bool {
	manager := GetGlobalRenewalServiceManager()
	return manager.IsRunning(ctxParams)
}

// GetAutoRenewalService returns the current renewal service instance
func GetAutoRenewalService() *AutoRenewalService {
	manager := GetGlobalRenewalServiceManager()
	return manager.GetService()
}

// ServiceStatus represents the status of the renewal service
type ServiceStatus struct {
	Running       bool   `json:"running"`
	CheckInterval string `json:"check_interval"`
	Message       string `json:"message"`
}

// GetRenewalServiceStatus returns the current status of the renewal service
func GetRenewalServiceStatus(ctxParams *uh.ContextParams) *ServiceStatus {
	manager := GetGlobalRenewalServiceManager()
	service := manager.GetService()

	status := &ServiceStatus{
		Running: manager.IsRunning(ctxParams),
	}

	// Load persistent state for additional info
	persistentState, err := LoadServiceState(ctxParams)
	if err == nil && persistentState.Running {
		status.CheckInterval = "1h0m0s" // Default interval
		if status.Running {
			status.Message = fmt.Sprintf("Auto-renewal service is running (PID: %d, started: %s)",
				persistentState.PID, persistentState.StartTime.Format("2006-01-02 15:04:05"))
		} else {
			status.Message = "Auto-renewal service process not found"
		}
	} else {
		if service != nil {
			status.CheckInterval = service.checkInterval.String()
		}
		if status.Running {
			status.Message = "Auto-renewal service is running normally"
		} else {
			status.Message = "Auto-renewal service is not running"
		}
	}

	return status
}

// EnableAutoRenewalForFile enables auto-renewal for a specific file
func EnableAutoRenewalForFile(ctxParams *uh.ContextParams, fileHash string) error {
	service := GetAutoRenewalService()
	if service == nil {
		return fmt.Errorf("auto-renewal service is not running")
	}

	return service.EnableAutoRenewal(fileHash)
}

// DisableAutoRenewalForFile disables auto-renewal for a specific file
func DisableAutoRenewalForFile(ctxParams *uh.ContextParams, fileHash string) error {
	service := GetAutoRenewalService()
	if service == nil {
		return fmt.Errorf("auto-renewal service is not running")
	}

	return service.DisableAutoRenewal(fileHash)
}

// GetAutoRenewalStatusForFile gets the auto-renewal status for a specific file
func GetAutoRenewalStatusForFile(ctxParams *uh.ContextParams, fileHash string) (*RenewalInfo, error) {
	service := GetAutoRenewalService()
	if service == nil {
		return nil, fmt.Errorf("auto-renewal service is not running")
	}

	return service.GetAutoRenewalStatus(fileHash)
}
