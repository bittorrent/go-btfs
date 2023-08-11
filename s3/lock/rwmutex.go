package lock

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"time"
)

// A TRWMutex is a mutual exclusion lock with timeouts.
type TRWMutex struct {
	isWriteLock bool
	ref         int
	mu          sync.Mutex // Mutex to prevent multiple simultaneous locks
}

// NewTRWMutex - initializes a new lsync RW mutex.
func NewTRWMutex() *TRWMutex {
	return &TRWMutex{}
}

// Lock holds a write lock on lm.
//
// If the lock is already in use, the calling go routine
// blocks until the mutex is available.
func (m *TRWMutex) Lock() {
	const isWriteLock = true
	m.lockLoop(context.Background(), math.MaxInt64, isWriteLock)
}

// GetLock tries to get a write lock on lm before the timeout occurs.
func (m *TRWMutex) GetLock(ctx context.Context, timeout time.Duration) (locked bool) {
	const isWriteLock = true
	return m.lockLoop(ctx, timeout, isWriteLock)
}

// RLock holds a read lock on lm.
//
// If one or more read lock are already in use, it will grant another lock.
// Otherwise the calling go routine blocks until the mutex is available.
func (m *TRWMutex) RLock() {
	const isWriteLock = false
	m.lockLoop(context.Background(), 1<<63-1, isWriteLock)
}

// GetRLock tries to get a read lock on lm before the timeout occurs.
func (m *TRWMutex) GetRLock(ctx context.Context, timeout time.Duration) (locked bool) {
	const isWriteLock = false
	return m.lockLoop(ctx, timeout, isWriteLock)
}

func (m *TRWMutex) lock(isWriteLock bool) (locked bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if isWriteLock {
		if m.ref == 0 && !m.isWriteLock {
			m.ref = 1
			m.isWriteLock = true
			locked = true
		}
	} else {
		if !m.isWriteLock {
			m.ref++
			locked = true
		}
	}

	return locked
}

const (
	lockRetryInterval = 50 * time.Millisecond
)

// lockLoop will acquire either a read or a write lock
//
// The call will block until the lock is granted using a built-in
// timing randomized back-off algorithm to try again until successful
func (m *TRWMutex) lockLoop(ctx context.Context, timeout time.Duration, isWriteLock bool) (locked bool) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	retryCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		select {
		case <-retryCtx.Done():
			// Caller context canceled or we timedout,
			// return false anyways for both situations.
			return false
		default:
			if m.lock(isWriteLock) {
				return true
			}
			time.Sleep(time.Duration(r.Float64() * float64(lockRetryInterval)))
		}
	}
}

// Unlock unlocks the write lock.
//
// It is a run-time error if lm is not locked on entry to Unlock.
func (m *TRWMutex) Unlock() {
	isWriteLock := true
	success := m.unlock(isWriteLock)
	if !success {
		panic("Trying to Unlock() while no Lock() is active")
	}
}

// RUnlock releases a read lock held on lm.
//
// It is a run-time error if lm is not locked on entry to RUnlock.
func (m *TRWMutex) RUnlock() {
	isWriteLock := false
	success := m.unlock(isWriteLock)
	if !success {
		panic("Trying to RUnlock() while no RLock() is active")
	}
}

func (m *TRWMutex) unlock(isWriteLock bool) (unlocked bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Try to release lock.
	if isWriteLock {
		if m.isWriteLock && m.ref == 1 {
			m.ref = 0
			m.isWriteLock = false
			unlocked = true
		}
	} else {
		if !m.isWriteLock {
			if m.ref > 0 {
				m.ref--
				unlocked = true
			}
		}
	}

	return unlocked
}

// ForceUnlock will forcefully clear a write or read lock.
func (m *TRWMutex) ForceUnlock() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ref = 0
	m.isWriteLock = false
}
