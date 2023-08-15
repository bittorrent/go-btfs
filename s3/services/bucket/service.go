package bucket

import (
	"context"
	"time"

	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/lock"
	"github.com/bittorrent/go-btfs/s3/policy"
	"github.com/bittorrent/go-btfs/s3/services"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	bucketPrefix           = "bkt/"
	globalOperationTimeout = 5 * time.Minute
	deleteOperationTimeout = 1 * time.Minute
)

var _ handlers.BucketService = (*Service)(nil)

// Service captures all bucket metadata for a given cluster.
type Service struct {
	providers   services.Providerser
	nsLock      *lock.NsLockMap
	emptyBucket func(ctx context.Context, bucket string) (bool, error)
}

// NewService - creates new policy system.
func NewService(providers services.Providerser, options ...Option) (s *Service) {
	s = &Service{
		providers: providers,
		nsLock:    lock.NewNSLock(),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

func (s *Service) CheckACL(accessKeyRecord *handlers.AccessKeyRecord, bucketName string, action action.Action) (err error) {
	//todo 是否需要判断原始的
	if bucketName == "" {
		return handlers.ErrBucketNotFound
	}

	bucketMeta, err := s.GetBucketMeta(context.Background(), bucketName)
	if err != nil {
		return err
	}

	//todo 注意：如果action是CreateBucketAction，HasBucket(ctx, bucketName)进行判断

	if policy.IsAllowed(bucketMeta.Owner == accessKeyRecord.Key, bucketMeta.Acl, action) == false {
		return handlers.ErrBucketAccessDenied
	}
	return
}

// NewBucketMetadata creates handlers.BucketMetadata with the supplied name and Created to Now.
func (s *Service) NewBucketMetadata(name, region, accessKey, acl string) *handlers.BucketMetadata {
	return &handlers.BucketMetadata{
		Name:    name,
		Region:  region,
		Owner:   accessKey,
		Acl:     acl,
		Created: time.Now().UTC(),
	}
}

// NewNSLock - initialize a new namespace RWLocker instance.
func (s *Service) NewNSLock(bucket string) lock.RWLocker {
	return s.nsLock.NewNSLock("meta", bucket)
}

func (s *Service) SetEmptyBucket(emptyBucket func(ctx context.Context, bucket string) (bool, error)) {
	s.emptyBucket = emptyBucket
}

// setBucketMeta - sets a new metadata in-db
func (s *Service) setBucketMeta(bucket string, meta *handlers.BucketMetadata) error {
	return s.providers.GetStateStore().Put(bucketPrefix+bucket, meta)
}

// CreateBucket - create a new Bucket
func (s *Service) CreateBucket(ctx context.Context, bucket, region, accessKey, acl string) error {
	lk := s.NewNSLock(bucket)
	lkctx, err := lk.GetLock(ctx, globalOperationTimeout)
	if err != nil {
		return err
	}
	ctx = lkctx.Context()
	defer lk.Unlock(lkctx.Cancel)

	return s.setBucketMeta(bucket, s.NewBucketMetadata(bucket, region, accessKey, acl))
}

func (s *Service) getBucketMeta(bucket string) (meta handlers.BucketMetadata, err error) {
	err = s.providers.GetStateStore().Get(bucketPrefix+bucket, &meta)
	if err == leveldb.ErrNotFound {
		err = handlers.ErrBucketNotFound
	}
	return meta, err
}

// GetBucketMeta metadata for a bucket.
func (s *Service) GetBucketMeta(ctx context.Context, bucket string) (meta handlers.BucketMetadata, err error) {
	lk := s.NewNSLock(bucket)
	lkctx, err := lk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return handlers.BucketMetadata{}, err
	}
	ctx = lkctx.Context()
	defer lk.RUnlock(lkctx.Cancel)

	return s.getBucketMeta(bucket)
}

// HasBucket  metadata for a bucket.
func (s *Service) HasBucket(ctx context.Context, bucket string) bool {
	_, err := s.GetBucketMeta(ctx, bucket)
	return err == nil
}

// DeleteBucket bucket.
func (s *Service) DeleteBucket(ctx context.Context, bucket string) error {
	lk := s.NewNSLock(bucket)
	lkctx, err := lk.GetLock(ctx, deleteOperationTimeout)
	if err != nil {
		return err
	}
	ctx = lkctx.Context()
	defer lk.Unlock(lkctx.Cancel)

	if _, err = s.getBucketMeta(bucket); err != nil {
		return err
	}

	if empty, err := s.emptyBucket(ctx, bucket); err != nil {
		return err
	} else if !empty {
		return handlers.ErrBucketNotEmpty
	}

	return s.providers.GetStateStore().Delete(bucketPrefix + bucket)
}

// GetAllBucketsOfUser metadata for all bucket.
func (s *Service) GetAllBucketsOfUser(ctx context.Context, username string) ([]handlers.BucketMetadata, error) {
	var m []handlers.BucketMetadata
	all, err := s.providers.GetStateStore().ReadAllChan(ctx, bucketPrefix, "")
	if err != nil {
		return nil, err
	}
	for entry := range all {
		data := handlers.BucketMetadata{}
		if err = entry.UnmarshalValue(&data); err != nil {
			continue
		}
		if data.Owner != username {
			continue
		}
		m = append(m, data)
	}
	return m, nil
}

// UpdateBucketAcl .
func (s *Service) UpdateBucketAcl(ctx context.Context, bucket, acl, accessKey string) error {
	lk := s.NewNSLock(bucket)
	lkctx, err := lk.GetLock(ctx, globalOperationTimeout)
	if err != nil {
		return err
	}
	ctx = lkctx.Context()
	defer lk.Unlock(lkctx.Cancel)

	meta, err := s.getBucketMeta(bucket)
	if err != nil {
		return err
	}

	meta.Acl = acl
	return s.setBucketMeta(bucket, &meta)
}

// GetBucketAcl .
func (s *Service) GetBucketAcl(ctx context.Context, bucket string) (string, error) {
	meta, err := s.GetBucketMeta(ctx, bucket)
	if err != nil {
		return "", err
	}
	return meta.Acl, nil
}
