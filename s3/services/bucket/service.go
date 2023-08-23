package bucket

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/services"
	"time"

	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/ctxmu"
	"github.com/bittorrent/go-btfs/s3/policy"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	bucketPrefix           = "bkt/"
	defaultUpdateTimeoutMS = 200
)

var _ services.BucketService = (*Service)(nil)

// Service captures all bucket metadata for a given cluster.
type Service struct {
	providers     providers.Providerser
	emptyBucket   func(ctx context.Context, bucket string) (bool, error)
	locks         *ctxmu.MultiCtxRWMutex
	updateTimeout time.Duration
}

// NewService - creates new policy system.
func NewService(providers providers.Providerser, options ...Option) (s *Service) {
	s = &Service{
		providers:     providers,
		locks:         ctxmu.NewDefaultMultiCtxRWMutex(),
		updateTimeout: time.Duration(defaultUpdateTimeoutMS) * time.Millisecond,
	}
	for _, option := range options {
		option(s)
	}
	return s
}

func (s *Service) CheckACL(accessKeyRecord *services.AccessKey, bucketName string, action action.Action) (err error) {
	//需要判断bucketName是否为空字符串
	if bucketName == "" {
		return services.ErrBucketNotFound
	}

	bucketMeta, err := s.GetBucketMeta(context.Background(), bucketName)
	if err != nil {
		return err
	}

	if policy.IsAllowed(bucketMeta.Owner == accessKeyRecord.Key, bucketMeta.Acl, action) == false {
		return services.ErrBucketAccessDenied
	}
	return
}

// NewBucketMetadata creates handlers.BucketMetadata with the supplied name and Created to Now.
func (s *Service) NewBucketMetadata(name, region, accessKey, acl string) *services.BucketMetadata {
	return &services.BucketMetadata{
		Name:    name,
		Region:  region,
		Owner:   accessKey,
		Acl:     acl,
		Created: time.Now().UTC(),
	}
}

// lockSetBucketMeta - sets a new metadata in-db
func (s *Service) lockSetBucketMeta(bucket string, meta *services.BucketMetadata) error {
	return s.providers.GetStateStore().Put(bucketPrefix+bucket, meta)
}

// CreateBucket - create a new Bucket
func (s *Service) CreateBucket(ctx context.Context, bucket, region, accessKey, acl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.updateTimeout)
	defer cancel()

	err := s.locks.Lock(ctx, bucket)
	if err != nil {
		return err
	}
	defer s.locks.Unlock(bucket)

	return s.lockSetBucketMeta(bucket, s.NewBucketMetadata(bucket, region, accessKey, acl))
}

func (s *Service) lockGetBucketMeta(bucket string) (meta services.BucketMetadata, err error) {
	err = s.providers.GetStateStore().Get(bucketPrefix+bucket, &meta)
	if err == leveldb.ErrNotFound {
		err = services.ErrBucketNotFound
	}
	return meta, err
}

// GetBucketMeta metadata for a bucket.
func (s *Service) GetBucketMeta(ctx context.Context, bucket string) (meta services.BucketMetadata, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.updateTimeout)
	defer cancel()

	err = s.locks.RLock(ctx, bucket)
	if err != nil {
		return services.BucketMetadata{Name: bucket}, err
	}
	defer s.locks.RUnlock(bucket)

	return s.lockGetBucketMeta(bucket)
}

// HasBucket  metadata for a bucket.
func (s *Service) HasBucket(ctx context.Context, bucket string) bool {
	_, err := s.GetBucketMeta(ctx, bucket)
	return err == nil
}

// DeleteBucket bucket.
func (s *Service) DeleteBucket(ctx context.Context, bucket string) error {
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

	if empty, err := s.emptyBucket(ctx, bucket); err != nil {
		return err
	} else if !empty {
		return services.ErrSetBucketEmptyFailed
	}

	return s.providers.GetStateStore().Delete(bucketPrefix + bucket)
}

func (s *Service) SetEmptyBucket(emptyBucket func(ctx context.Context, bucket string) (bool, error)) {
	s.emptyBucket = emptyBucket
}

// GetAllBucketsOfUser metadata for all bucket.
func (s *Service) GetAllBucketsOfUser(username string) (list []*services.BucketMetadata, err error) {
	err = s.providers.GetStateStore().Iterate(bucketPrefix, func(key, _ []byte) (stop bool, er error) {
		record := &services.BucketMetadata{}
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
func (s *Service) UpdateBucketAcl(ctx context.Context, bucket, acl string) error {
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
func (s *Service) GetBucketAcl(ctx context.Context, bucket string) (string, error) {
	meta, err := s.GetBucketMeta(ctx, bucket)
	if err != nil {
		return "", err
	}
	return meta.Acl, nil
}
