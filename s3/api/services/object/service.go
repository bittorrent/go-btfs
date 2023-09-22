package object

import (
	"context"
	"fmt"
	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/api/providers"
	"github.com/bittorrent/go-btfs/s3/ctxmu"
	"github.com/bittorrent/go-btfs/s3/policy"
	"io"
	"strings"
	"time"
)

var _ Service = (*service)(nil)

type service struct {
	providers        providers.Providerser
	lock             ctxmu.MultiCtxRWLocker
	keySeparator     string
	bucketSpace      string
	objectSpace      string
	uploadSpace      string
	cidrefSpace      string
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
		cidrefSpace:      defaultCidrefSpace,
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
	key = s.getAllUploadsKeyPrefix(bucname) + strings.Join([]string{objname, uploadid}, s.keySeparator)
	return
}

func (s *service) getUploadPartKey(uplkey string, idx int) (key string) {
	key = fmt.Sprintf("%s_%d", uplkey, idx)
	return
}

func (s *service) getAllCidrefsKeyPrefix(cid string) (prefix string) {
	prefix = strings.Join([]string{s.cidrefSpace, cid, ""}, s.keySeparator)
	return
}

func (s *service) getCidrefKey(cid, to string) (key string) {
	key = s.getAllCidrefsKeyPrefix(cid) + to
	return
}

func (s *service) opctx(parent context.Context) (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = context.WithTimeout(parent, s.operationTimeout)
	return
}

func (s *service) checkACL(owner, acl, user string, act action.Action) (allow bool) {
	own := user != "" && user == owner
	allow = policy.IsAllowed(own, acl, act)
	return
}

func (s *service) addBodyRef(ctx context.Context, cid, tokey string) (err error) {
	// Cid reference key
	crfkey := s.getCidrefKey(cid, tokey)

	// Add cid reference
	err = s.providers.StateStore().Put(crfkey, nil)

	return
}

func (s *service) removeBodyRef(ctx context.Context, cid, tokey string) (err error) {
	// This object cid reference key
	crfkey := s.getCidrefKey(cid, tokey)

	// Delete cid ref of this object
	err = s.providers.StateStore().Delete(crfkey)

	return
}

func (s *service) storeBody(ctx context.Context, body io.Reader, tokey string) (cid string, err error) {
	// RLock all cid refs to enable no cid will be deleted
	err = s.lock.RLock(ctx, s.cidrefSpace)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(s.cidrefSpace)

	// Store body and get the cid
	cid, err = s.providers.FileStore().Store(body)
	if err != nil {
		return
	}

	// Add cid reference
	err = s.addBodyRef(ctx, cid, tokey)

	return
}

func (s *service) removeBody(ctx context.Context, cid, tokey string) (err error) {
	// Flag to mark cid be referenced by other object
	otherRef := false

	// Log removing
	defer func() {
		fmt.Printf("s3-api: remove <%s>, ref <%s>, other-ref - %v, err: %v\n", cid, tokey, otherRef, err)
	}()

	// Lock all cid refs to enable new cid reference can not be added when
	// remove is executing
	err = s.lock.Lock(ctx, s.cidrefSpace)
	if err != nil {
		return
	}
	defer s.lock.Unlock(s.cidrefSpace)

	// Remove cid ref of this object
	err = s.removeBodyRef(ctx, cid, tokey)
	if err != nil {
		return
	}

	// All this cid references prefix
	allRefsPrefix := s.getAllCidrefsKeyPrefix(cid)

	// Iterate all this cid refs, if exists other object's ref, set
	// the otherRef mark to true
	err = s.providers.StateStore().Iterate(allRefsPrefix, func(key, _ []byte) (stop bool, err error) {
		otherRef = true
		stop = true
		return
	})
	if err != nil {
		return
	}

	// Exists other refs, cid body can not be removed
	if otherRef {
		return
	}

	// No other refs to this cid, remove it
	err = s.providers.FileStore().Remove(cid)

	return
}
