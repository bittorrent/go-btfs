package accesskey

import (
	"github.com/bittorrent/go-btfs/s3/ctxmu"
	"time"
)

const (
	defaultSecretLength   = 32
	defaultStoreKeyPrefix = "access-keys:"
	defaultWaitLockTimout = 2 * time.Minute
)

var defaultLock = ctxmu.NewDefaultMultiCtxRWMutex()

type Option func(svc *service)

func WithSecretLength(length int) Option {
	return func(svc *service) {
		svc.secretLength = length
	}
}

func WithStoreKeyPrefix(prefix string) Option {
	return func(svc *service) {
		svc.storeKeyPrefix = prefix
	}
}

func WithWaitLockTimout(timout time.Duration) Option {
	return func(svc *service) {
		svc.waitLockTimeout = timout
	}
}

func WithLock(lock ctxmu.MultiCtxRWLocker) Option {
	return func(svc *service) {
		svc.lock = lock
	}
}
