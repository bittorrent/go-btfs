package object

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/ctxmu"
	"github.com/bittorrent/go-btfs/s3/policy"
	"regexp"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/s3/providers"
)

const (
	defaultKeySeparator     = "/"
	defaultBucketSpace      = "bkt"
	defaultObjectSpace      = "obj"
	defaultUploadSpace      = "upl"
	defaultOperationTimeout = 5 * time.Minute

	bucketPrefix           = "bkt/"
	objectKeyFormat        = "obj/%s/%s"
	objectPrefix           = "obj/%s/"
	allObjectPrefixFormat  = "obj/%s/%s"
	allObjectSeekKeyFormat = "obj/%s/%s"

	uploadKeyFormat        = "uploadObj/%s/%s/%s"
	allUploadPrefixFormat  = "uploadObj/%s/%s"
	allUploadSeekKeyFormat = "uploadObj/%s/%s/%s"

	deleteKeyFormat       = "delObj/%s"
	allDeletePrefixFormat = "delObj/"

	globalOperationTimeout = 5 * time.Minute
	deleteOperationTimeout = 1 * time.Minute

	maxCpuPercent        = 60
	maxUsedMemoryPercent = 80
)

var etagRegex = regexp.MustCompile("\"*?([^\"]*?)\"*?$")

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
}

func NewService(providers providers.Providerser, options ...Option) Service {
	s := &service{
		providers:        providers,
		lock:             ctxmu.NewDefaultMultiCtxRWMutex(),
		keySeparator:     defaultKeySeparator,
		bucketSpace:      defaultBucketSpace,
		objectSpace:      defaultObjectSpace,
		uploadSpace:      defaultUploadSpace,
		operationTimeout: defaultOperationTimeout,
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// common helper methods

func (s *service) getBucketKeyPrefix() (prefix string) {
	prefix = strings.Join([]string{s.bucketSpace}, s.keySeparator)
	return
}

func (s *service) getObjectKeyPrefix(bucname string) (prefix string) {
	prefix = strings.Join([]string{s.objectSpace, bucname}, s.keySeparator)
	return
}

func (s *service) getUploadKeyPrefix(bucname, objname string) (prefix string) {
	prefix = strings.Join([]string{s.uploadSpace, bucname, objname}, s.keySeparator)
	return
}

func (s *service) getBucketKey(bucname string) (key string) {
	key = s.getBucketKeyPrefix() + bucname
	return
}

func (s *service) getObjectKey(bucname, objname string) (key string) {
	key = s.getObjectKeyPrefix(bucname) + objname
	return
}

func (s *service) getUploadKey(bucname, objname, uploadid string) (key string) {
	key = s.getUploadKeyPrefix(bucname, objname) + uploadid
	return
}

func (s *service) checkAcl(owner, acl, user string, act action.Action) (allow bool) {
	own := user != "" && user == owner
	allow = policy.IsAllowed(own, acl, act)
	return
}

func (s *service) opctx(parent context.Context) (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = context.WithTimeout(parent, s.operationTimeout)
	return
}
