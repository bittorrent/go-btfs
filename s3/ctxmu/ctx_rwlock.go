package ctxmu

import (
	"golang.org/x/net/context"
	"math/rand"
	"sync"
	"time"
)

const lockRetryInterval = 50 * time.Millisecond

type CtxLocker interface {
	Lock(ctx context.Context) (err error)
	Unlock()
}

type CtxRWLocker interface {
	CtxLocker
	RLock(ctx context.Context) (err error)
	RUnlock()
}

type CtxRWMutex struct {
	lock sync.RWMutex
}

func (c *CtxRWMutex) Lock(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if c.lock.TryLock() {
				return nil
			}
		}
		time.Sleep(c.getRandInterval())
	}
}

func (c *CtxRWMutex) Unlock() {
	c.lock.Unlock()
}

func (c *CtxRWMutex) RLock(ctx context.Context) (err error) {
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		default:
			if c.lock.TryRLock() {
				return
			}
		}
		time.Sleep(c.getRandInterval())
	}
}

func (c *CtxRWMutex) RUnlock() {
	c.lock.RUnlock()
}

func (c *CtxRWMutex) getRandInterval() time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Duration(r.Float64() * float64(lockRetryInterval))
}
