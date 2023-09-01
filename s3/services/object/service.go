package object

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/ctxmu"
	"github.com/bittorrent/go-btfs/s3/policy"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/s3/providers"
)

var _ Service = (*service)(nil)

// service captures all bucket metadata for a given cluster.
type service struct {
	providers        providers.Providerser
	lock             ctxmu.MultiCtxRWLocker
	keySeparator     string
	bucketSpace      string
	objectSpace      string
	uploadSpace      string
	operationTimeout time.Duration
	closeBodyTimeout time.Duration
}

func NewService(providers providers.Providerser, options ...Option) Service {
	s := &service{
		providers:        providers,
		lock:             defaultLock,
		keySeparator:     defaultKeySeparator,
		bucketSpace:      defaultBucketSpace,
		objectSpace:      defaultObjectSpace,
		uploadSpace:      defaultUploadSpace,
		operationTimeout: defaultOperationTimeout,
		closeBodyTimeout: defaultCloseBodyTimeout,
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// common helper methods

func (s *service) getAllBucketsKeyPrefix() (prefix string) {
	prefix = strings.Join([]string{s.bucketSpace, ""}, s.keySeparator)
	return
}

func (s *service) getBucketKey(bucname string) (key string) {
	key = s.getAllBucketsKeyPrefix() + bucname
	return
}

func (s *service) getAllObjectsKeyPrefix(bucname string) (prefix string) {
	prefix = strings.Join([]string{s.objectSpace, bucname, ""}, s.keySeparator)
	return
}

func (s *service) getObjectKey(bucname, objname string) (key string) {
	key = s.getAllObjectsKeyPrefix(bucname) + objname
	return
}

func (s *service) getAllUploadsKeyPrefix(bucname string) (prefix string) {
	prefix = strings.Join([]string{s.uploadSpace, bucname, ""}, s.keySeparator)
	return
}

func (s *service) getUploadKey(bucname, objname, uploadid string) (key string) {
	key = strings.Join([]string{s.getAllUploadsKeyPrefix(bucname), objname, uploadid}, s.keySeparator)
	return
}

func (s *service) opctx(parent context.Context) (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = context.WithTimeout(parent, s.operationTimeout)
	return
}

func (s *service) checkAcl(owner, acl, user string, act action.Action) (allow bool) {
	own := user != "" && user == owner
	allow = policy.IsAllowed(own, acl, act)
	return
}
