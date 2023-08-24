package bucket

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/services"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"time"

	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/ctxmu"
	"github.com/bittorrent/go-btfs/s3/policy"
)

const (
	bucketPrefix           = "bkt/"
	defaultUpdateTimeoutMS = 200
)

var _ Service = (*service)(nil)

// service captures all bucket metadata for a given cluster.
type service struct {
	providers     providers.Providerser
	emptyBucket   func(ctx context.Context, bucket string) (bool, error)
	locks         *ctxmu.MultiCtxRWMutex
	updateTimeout time.Duration
}

// NewService - creates new policy system.
func NewService(providers providers.Providerser, options ...Option) Service {
	s := &service{
		providers:     providers,
		locks:         ctxmu.NewDefaultMultiCtxRWMutex(),
		updateTimeout: time.Duration(defaultUpdateTimeoutMS) * time.Millisecond,
	}
	for _, option := range options {
		option(s)
	}
	return s
}

func (s *service) CheckACL(ack *accesskey.AccessKey, bucketName string, act action.Action) (err error) {
	if act == action.ListBucketAction {
		if ack.Key == "" {
			err = services.RespErrAccessDenied
		}
		return
	}

	//需要判断bucketName是否为空字符串
	if bucketName == "" {
		return services.RespErrNoSuchBucket
	}

	bucketMeta, err := s.GetBucketMeta(context.Background(), bucketName)
	if err != nil {
		return err
	}

	if policy.IsAllowed(bucketMeta.Owner == ack.Key, bucketMeta.Acl, act) == false {
		return services.RespErrAccessDenied
	}
	return
}

// NewBucketMetadata creates handlers.Bucket with the supplied name and Created to Now.
func (s *service) NewBucketMetadata(name, region, accessKey, acl string) *Bucket {
	return &Bucket{
		Name:    name,
		Region:  region,
		Owner:   accessKey,
		Acl:     acl,
		Created: time.Now().UTC(),
	}
}

// lockSetBucketMeta - sets a new metadata in-db
func (s *service) lockSetBucketMeta(bucket string, meta *Bucket) error {
	return s.providers.GetStateStore().Put(bucketPrefix+bucket, meta)
}

// CreateBucket - create a new Bucket
func (s *service) CreateBucket(ctx context.Context, bucket, region, accessKey, acl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.updateTimeout)
	defer cancel()

	err := s.locks.Lock(ctx, bucket)
	if err != nil {
		return err
	}
	defer s.locks.Unlock(bucket)

	return s.lockSetBucketMeta(bucket, s.NewBucketMetadata(bucket, region, accessKey, acl))
}

func (s *service) lockGetBucketMeta(bucket string) (meta Bucket, err error) {
	err = s.providers.GetStateStore().Get(bucketPrefix+bucket, &meta)
	if errors.Is(err, providers.ErrStateStoreNotFound) {
		err = services.RespErrNoSuchBucket
	}
	return
}

// GetBucketMeta metadata for a bucket.
func (s *service) GetBucketMeta(ctx context.Context, bucket string) (meta Bucket, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.updateTimeout)
	defer cancel()

	err = s.locks.RLock(ctx, bucket)
	if err != nil {
		return Bucket{Name: bucket}, err
	}
	defer s.locks.RUnlock(bucket)

	return s.lockGetBucketMeta(bucket)
}

// HasBucket  metadata for a bucket.
func (s *service) HasBucket(ctx context.Context, bucket string) bool {
	_, err := s.GetBucketMeta(ctx, bucket)
	return err == nil
}

// DeleteBucket bucket.
func (s *service) DeleteBucket(ctx context.Context, bucket string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.updateTimeout)
	defer cancel()

	err := s.locks.Lock(ctx, bucket)
	if err != nil {
		return err
	}
	defer s.locks.Unlock(bucket)

	if _, err = s.lockGetBucketMeta(bucket); err != nil {
		return err
	}

	empty, err := s.emptyBucket(ctx, bucket)
	if err != nil {
		return err
	}

	if !empty {
		return errors.New("bucket not empty")
	}

	return s.providers.GetStateStore().Delete(bucketPrefix + bucket)
}

func (s *service) SetEmptyBucket(emptyBucket func(ctx context.Context, bucket string) (bool, error)) {
	s.emptyBucket = emptyBucket
}

// GetAllBucketsOfUser metadata for all bucket.
func (s *service) GetAllBucketsOfUser(username string) (list []*Bucket, err error) {
	err = s.providers.GetStateStore().Iterate(bucketPrefix, func(key, _ []byte) (stop bool, er error) {
		record := &Bucket{}
		er = s.providers.GetStateStore().Get(string(key), record)
		if er != nil {
			return
		}
		if record.Owner == username {
			list = append(list, record)
		}

		return
	})

	return
}

// UpdateBucketAcl .
func (s *service) UpdateBucketAcl(ctx context.Context, bucket, acl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.updateTimeout)
	defer cancel()

	err := s.locks.Lock(ctx, bucket)
	if err != nil {
		return err
	}
	defer s.locks.Unlock(bucket)

	meta, err := s.lockGetBucketMeta(bucket)
	if err != nil {
		return err
	}

	meta.Acl = acl
	return s.lockSetBucketMeta(bucket, &meta)
}

// GetBucketAcl .
func (s *service) GetBucketAcl(ctx context.Context, bucket string) (string, error) {
	meta, err := s.GetBucketMeta(ctx, bucket)
	if err != nil {
		return "", err
	}
	return meta.Acl, nil
}

// EmptyBucket object中后续添加
func (s *service) EmptyBucket(ctx context.Context, bucket string) (bool, error) {
	//loi, err := s.ListObjects(ctx, bucket, "", "", "", 1)
	//if err != nil {
	//	return false, err
	//}
	//return len(loi.Objects) == 0, nil
	return true, nil
}
