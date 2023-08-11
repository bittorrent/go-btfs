package lock

import (
	"context"
	"errors"
	logging "github.com/ipfs/go-log/v2"
	"path"
	"sort"
	"strings"
	"sync"
	"time"
)

var log = logging.Logger("nslocker")

// OperationTimedOut - a timeout occurred.
type OperationTimedOut struct{}

func (e OperationTimedOut) Error() string {
	return "Operation timed out"
}

// RWLocker - locker interface to introduce GetRLock, RUnlock.
type RWLocker interface {
	GetLock(ctx context.Context, timeout time.Duration) (lkCtx LockContext, timedOutErr error)
	Unlock(cancel context.CancelFunc)
	GetRLock(ctx context.Context, timeout time.Duration) (lkCtx LockContext, timedOutErr error)
	RUnlock(cancel context.CancelFunc)
}

// LockContext lock context holds the lock backed context and canceler for the context.
type LockContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// Context returns lock context
func (l LockContext) Context() context.Context {
	return l.ctx
}

// Cancel function calls cancel() function
func (l LockContext) Cancel() {
	if l.cancel != nil {
		l.cancel()
	}
}

// NewNSLock - return a new name space lock map.
func NewNSLock() *NsLockMap {
	return &NsLockMap{
		lockMap: make(map[string]*nsLock),
	}
}

// nsLock - provides primitives for locking critical namespace regions.
type nsLock struct {
	ref int32
	*TRWMutex
}

// NsLockMap - namespace lock map, provides primitives to Lock,
// Unlock, RLock and RUnlock.
type NsLockMap struct {
	lockMap      map[string]*nsLock
	lockMapMutex sync.Mutex
}

// Lock the namespace resource.
func (n *NsLockMap) lock(ctx context.Context, volume string, path string, readLock bool, timeout time.Duration) (locked bool) {
	resource := PathJoin(volume, path)

	n.lockMapMutex.Lock()
	nsLk, found := n.lockMap[resource]
	if !found {
		nsLk = &nsLock{
			TRWMutex: NewTRWMutex(),
		}
	}
	nsLk.ref++
	n.lockMap[resource] = nsLk
	n.lockMapMutex.Unlock()

	// Locking here will block (until timeout).
	if readLock {
		locked = nsLk.GetRLock(ctx, timeout)
	} else {
		locked = nsLk.GetLock(ctx, timeout)
	}

	if !locked { // We failed to get the lock
		// Decrement ref count since we failed to get the lock
		n.lockMapMutex.Lock()
		n.lockMap[resource].ref--
		if n.lockMap[resource].ref < 0 {
			log.Error(errors.New("resource reference count was lower than 0"))
		}
		if n.lockMap[resource].ref == 0 {
			// Remove from the map if there are no more references.
			delete(n.lockMap, resource)
		}
		n.lockMapMutex.Unlock()
	}

	return
}

// Unlock the namespace resource.
func (n *NsLockMap) unlock(volume string, path string, readLock bool) {
	resource := PathJoin(volume, path)

	n.lockMapMutex.Lock()
	defer n.lockMapMutex.Unlock()
	if _, found := n.lockMap[resource]; !found {
		return
	}
	if readLock {
		n.lockMap[resource].RUnlock()
	} else {
		n.lockMap[resource].Unlock()
	}
	n.lockMap[resource].ref--
	if n.lockMap[resource].ref < 0 {
		log.Error(errors.New("resource reference count was lower than 0"))
	}
	if n.lockMap[resource].ref == 0 {
		// Remove from the map if there are no more references.
		delete(n.lockMap, resource)
	}
}

// localLockInstance - frontend/top-level interface for namespace locks.
type localLockInstance struct {
	ns     *NsLockMap
	volume string
	paths  []string
}

// NewNSLock - returns a lock instance for a given volume and
// path. The returned lockInstance object encapsulates the nsLockMap,
// volume, path and operation ID.
func (n *NsLockMap) NewNSLock(volume string, paths ...string) RWLocker {
	sort.Strings(paths)
	return &localLockInstance{n, volume, paths}
}

// GetLock - block until write lock is taken or timeout has occurred.
func (li *localLockInstance) GetLock(ctx context.Context, timeout time.Duration) (_ LockContext, timedOutErr error) {
	const readLock = false
	success := make([]int, len(li.paths))
	for i, path := range li.paths {
		if !li.ns.lock(ctx, li.volume, path, readLock, timeout) {
			for si, sint := range success {
				if sint == 1 {
					li.ns.unlock(li.volume, li.paths[si], readLock)
				}
			}
			return LockContext{}, OperationTimedOut{}
		}
		success[i] = 1
	}
	return LockContext{ctx: ctx, cancel: func() {}}, nil
}

// Unlock - block until write lock is released.
func (li *localLockInstance) Unlock(cancel context.CancelFunc) {
	if cancel != nil {
		cancel()
	}
	const readLock = false
	for _, path := range li.paths {
		li.ns.unlock(li.volume, path, readLock)
	}
}

// GetRLock - block until read lock is taken or timeout has occurred.
func (li *localLockInstance) GetRLock(ctx context.Context, timeout time.Duration) (_ LockContext, timedOutErr error) {
	const readLock = true
	success := make([]int, len(li.paths))
	for i, path := range li.paths {
		if !li.ns.lock(ctx, li.volume, path, readLock, timeout) {
			for si, sint := range success {
				if sint == 1 {
					li.ns.unlock(li.volume, li.paths[si], readLock)
				}
			}
			return LockContext{}, OperationTimedOut{}
		}
		success[i] = 1
	}
	return LockContext{ctx: ctx, cancel: func() {}}, nil
}

// RUnlock - block until read lock is released.
func (li *localLockInstance) RUnlock(cancel context.CancelFunc) {
	if cancel != nil {
		cancel()
	}
	const readLock = true
	for _, path := range li.paths {
		li.ns.unlock(li.volume, path, readLock)
	}
}

// SlashSeparator - slash separator.
const SlashSeparator = "/"

// PathJoin - like path.Join() but retains trailing SlashSeparator of the last element
func PathJoin(elem ...string) string {
	trailingSlash := ""
	if len(elem) > 0 {
		if strings.HasSuffix(elem[len(elem)-1], SlashSeparator) {
			trailingSlash = SlashSeparator
		}
	}
	return path.Join(elem...) + trailingSlash
}
