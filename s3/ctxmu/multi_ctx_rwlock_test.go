package ctxmu

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"testing"
	"time"
)

func TestMultiCtxRWMutex_Lock(t *testing.T) {
	locks := NewMultiCtxRWMutex(func() CtxRWLocker {
		return &CtxRWMutex{}
	})
	eg := errgroup.Group{}
	key := "test_key"
	err := locks.Lock(context.Background(), key)
	if err != nil {
		t.Fatalf("can not lock")
	} else {
		t.Logf("can lock")
	}
	eg.Go(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		lerr := locks.Lock(ctx, key)
		if lerr == nil {
			t.Fatalf("can lock after locked")
		} else if !errors.Is(lerr, context.DeadlineExceeded) {
			t.Fatalf("timout lock return non DeadlineExceeded error: %v", lerr)
		} else {
			t.Logf("can not lock after locked")
		}
		lerr = locks.RLock(ctx, key)
		if lerr == nil {
			t.Fatalf("can rlock after locked")
		} else if !errors.Is(lerr, context.DeadlineExceeded) {
			t.Fatalf("timout rlock return non DeadlineExceeded error: %v", lerr)
		} else {
			t.Logf("can not rlock after locked")
		}
		locks.Unlock(key)
		lerr = locks.Lock(context.Background(), key)
		if lerr != nil {
			t.Fatalf("can not lock after unlocked")
		} else {
			t.Logf("can lock after unlocked")
		}
		locks.Unlock(key)
		lerr = locks.RLock(context.Background(), key)
		if lerr != nil {
			t.Fatalf("can not rlock after unlocked")
		} else {
			t.Logf("can rlock after unlocked")
		}
		locks.RUnlock(key)
		return nil
	})

	_ = eg.Wait()
}

func TestMultiCtxRWMutex_LockWithTimout(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	eg := errgroup.Group{}
	timeout := 50 * time.Millisecond
	locks := NewMultiCtxRWMutex(func() CtxRWLocker {
		return &CtxRWMutex{}
	})
	for i := 0; i < 1000; i++ {
		okey := fmt.Sprintf("key_%d", i)
		for j := 0; j < 100; j++ {
			key := okey
			n := j
			wt := rand.Intn(200)
			if j == 0 || j == 30 {
				eg.Go(func() error {
					lerr := locks.LockWithTimout(timeout, key)
					if lerr == nil {
						defer func() {
							t.Logf("%s %d Unlock: %v, %d", key, n, lerr, wt)
							locks.Unlock(key)
						}()
					}
					t.Logf("%s %d Lock: %v, %d", key, n, lerr, wt)
					time.Sleep(time.Duration(wt) * time.Millisecond)
					return nil
				})
			} else {
				eg.Go(func() error {
					lerr := locks.RLockWithTimout(timeout, key)
					if lerr == nil {
						defer func() {
							t.Logf("%s %d RLock: %v, %d", key, n, lerr, wt)
							locks.RUnlock(key)
						}()
					}
					t.Logf("%s %d RLock: %v, %d", key, n, lerr, wt)
					time.Sleep(time.Duration(wt) * time.Millisecond)
					return nil
				})

			}
		}
	}
	_ = eg.Wait()
}
