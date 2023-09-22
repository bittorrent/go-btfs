package object

import (
	"github.com/bittorrent/go-btfs/s3/ctxmu"
	"time"
)

const (
	defaultKeySeparator     = "/"
	defaultBucketSpace      = "s3:bkt"
	defaultObjectSpace      = "s3:obj"
	defaultUploadSpace      = "s3:upl"
	defaultCidrefSpace      = "s3:cid"
	defaultOperationTimeout = 5 * time.Minute
	defaultCloseBodyTimeout = 10 * time.Minute
)

var defaultLock = ctxmu.NewDefaultMultiCtxRWMutex()

type Option func(svc *service)

func WithKeySeparator(separator string) Option {
	return func(svc *service) {
		svc.keySeparator = separator
	}
}

func WithBucketSpace(space string) Option {
	return func(svc *service) {
		svc.bucketSpace = space
	}
}

func WithObjectSpace(space string) Option {
	return func(svc *service) {
		svc.objectSpace = space
	}
}

func WithUploadSpace(space string) Option {
	return func(svc *service) {
		svc.uploadSpace = space
	}
}

func WithCidrefSpace(space string) Option {
	return func(svc *service) {
		svc.cidrefSpace = space
	}
}

func WithOperationTimeout(timeout time.Duration) Option {
	return func(svc *service) {
		svc.operationTimeout = timeout
	}
}

func WithCloseBodyTimeout(timeout time.Duration) Option {
	return func(svc *service) {
		svc.closeBodyTimeout = timeout
	}
}

func WithLock(lock ctxmu.MultiCtxRWLocker) Option {
	return func(svc *service) {
		svc.lock = lock
	}
}
