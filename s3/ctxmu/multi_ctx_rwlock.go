package ctxmu

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type MultiCtxLocker interface {
	Lock(ctx context.Context, key interface{}) (err error)
	Unlock(key interface{})
}

type MultiCtxRWLocker interface {
	MultiCtxLocker
	RLock(ctx context.Context, key interface{}) (err error)
	RUnlock(key interface{})
}

type MultiCtxRWMutex struct {
	locks sync.Map
	pool  sync.Pool
}

func NewDefaultMultiCtxRWMutex() *MultiCtxRWMutex {
	return NewMultiCtxRWMutex(func() CtxRWLocker {
		return &CtxRWMutex{}
	})
}

func NewMultiCtxRWMutex(newCtxRWLock func() CtxRWLocker) *MultiCtxRWMutex {
	return &MultiCtxRWMutex{
		locks: sync.Map{},
		pool: sync.Pool{
			New: func() interface{} {
				return newCtxRWLock()
			},
		},
	}
}

type ctxRWLockRefCounter struct {
	count int64
	lock  CtxRWLocker
}

func (m *MultiCtxRWMutex) Lock(ctx context.Context, key interface{}) (err error) {
	counter, err := m.incrGetRWLockRefCounter(ctx, key)
	if err != nil {
		return
	}
	err = (counter.lock).Lock(ctx)
	if err != nil {
		m.decrPutRWLockRefCounter(key, counter)
	}
	return
}

func (m *MultiCtxRWMutex) LockWithTimout(timeout time.Duration, key interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = m.Lock(ctx, key)
	return
}

func (m *MultiCtxRWMutex) Unlock(key interface{}) {
	counter := m.mustGetCounter(key)
	counter.lock.Unlock()
	m.decrPutRWLockRefCounter(key, counter)
	return
}

func (m *MultiCtxRWMutex) RLock(ctx context.Context, key interface{}) (err error) {
	counter, err := m.incrGetRWLockRefCounter(ctx, key)
	if err != nil {
		return
	}
	err = (counter.lock).RLock(ctx)
	if err != nil {
		m.decrPutRWLockRefCounter(key, counter)
	}
	return
}

func (m *MultiCtxRWMutex) RLockWithTimout(timeout time.Duration, key interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = m.RLock(ctx, key)
	return
}

func (m *MultiCtxRWMutex) RUnlock(key interface{}) {
	counter := m.mustGetCounter(key)
	counter.lock.RUnlock()
	m.decrPutRWLockRefCounter(key, counter)
	return
}

func (m *MultiCtxRWMutex) mustGetCounter(key interface{}) (counter *ctxRWLockRefCounter) {
	actual, ok := m.locks.Load(key)
	if !ok {
		panic("key's lock has been invalidly freed")
	}
	counter = actual.(*ctxRWLockRefCounter)
	return
}

func (m *MultiCtxRWMutex) incrGetRWLockRefCounter(ctx context.Context, key interface{}) (counter *ctxRWLockRefCounter, err error) {
	for {
		err = ctx.Err()
		if err != nil {
			return
		}
		actual, _ := m.locks.LoadOrStore(key, &ctxRWLockRefCounter{
			count: 0,
			lock:  m.pool.Get().(*CtxRWMutex),
		})
		counter = actual.(*ctxRWLockRefCounter)
		old := counter.count
		if old < 0 {
			continue
		}
		if atomic.CompareAndSwapInt64(&counter.count, old, old+1) {
			break
		}
	}
	return
}

func (m *MultiCtxRWMutex) decrPutRWLockRefCounter(key interface{}, counter *ctxRWLockRefCounter) {
	atomic.AddInt64(&counter.count, -1)
	if atomic.CompareAndSwapInt64(&counter.count, 0, -1) {
		m.pool.Put(counter.lock)
		m.locks.Delete(key)
	}
}
