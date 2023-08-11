package store

import (
	"context"
	"time"
	
	"github.com/bittorrent/go-btfs/s3d/lock"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	bucketPrefix = "bkt/"
)

const (
	globalOperationTimeout = 5 * time.Minute
	deleteOperationTimeout = 1 * time.Minute
)

// BucketMetadata contains bucket metadata.
type BucketMetadata struct {
	Name    string
	Region  string
	Owner   string
	Acl     string
	Created time.Time
}

// NewBucketMetadata creates BucketMetadata with the supplied name and Created to Now.
func NewBucketMetadata(name, region, accessKey, acl string) *BucketMetadata {
	return &BucketMetadata{
		Name:    name,
		Region:  region,
		Owner:   accessKey,
		Acl:     acl,
		Created: time.Now().UTC(),
	}
}

// BucketMetadataSys captures all bucket metadata for a given cluster.
type BucketMetadataSys struct {
	db          storage.StateStorer
	nsLock      *lock.NsLockMap
	emptyBucket func(ctx context.Context, bucket string) (bool, error)
}

// NewBucketMetadataSys - creates new policy system.
func NewBucketMetadataSys(db storage.StateStorer) *BucketMetadataSys {
	return &BucketMetadataSys{
		db:     db,
		nsLock: lock.NewNSLock(),
	}
}

// NewNSLock - initialize a new namespace RWLocker instance.
func (sys *BucketMetadataSys) NewNSLock(bucket string) lock.RWLocker {
	return sys.nsLock.NewNSLock("meta", bucket)
}

func (sys *BucketMetadataSys) SetEmptyBucket(emptyBucket func(ctx context.Context, bucket string) (bool, error)) {
	sys.emptyBucket = emptyBucket
}

// setBucketMeta - sets a new metadata in-db
func (sys *BucketMetadataSys) setBucketMeta(bucket string, meta *BucketMetadata) error {
	return sys.db.Put(bucketPrefix+bucket, meta)
}

// CreateBucket - create a new Bucket
func (sys *BucketMetadataSys) CreateBucket(ctx context.Context, bucket, region, accessKey, acl string) error {
	lk := sys.NewNSLock(bucket)
	lkctx, err := lk.GetLock(ctx, globalOperationTimeout)
	if err != nil {
		return err
	}
	ctx = lkctx.Context()
	defer lk.Unlock(lkctx.Cancel)

	return sys.setBucketMeta(bucket, NewBucketMetadata(bucket, region, accessKey, acl))
}

func (sys *BucketMetadataSys) getBucketMeta(bucket string) (meta BucketMetadata, err error) {
	err = sys.db.Get(bucketPrefix+bucket, &meta)
	if err == leveldb.ErrNotFound {
		err = BucketNotFound{Bucket: bucket, Err: err}
	}
	return meta, err
}

// GetBucketMeta metadata for a bucket.
func (sys *BucketMetadataSys) GetBucketMeta(ctx context.Context, bucket string) (meta BucketMetadata, err error) {
	lk := sys.NewNSLock(bucket)
	lkctx, err := lk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return BucketMetadata{}, err
	}
	ctx = lkctx.Context()
	defer lk.RUnlock(lkctx.Cancel)

	return sys.getBucketMeta(bucket)
}

// HasBucket  metadata for a bucket.
func (sys *BucketMetadataSys) HasBucket(ctx context.Context, bucket string) bool {
	_, err := sys.GetBucketMeta(ctx, bucket)
	return err == nil
}

// DeleteBucket bucket.
func (sys *BucketMetadataSys) DeleteBucket(ctx context.Context, bucket string) error {
	lk := sys.NewNSLock(bucket)
	lkctx, err := lk.GetLock(ctx, deleteOperationTimeout)
	if err != nil {
		return err
	}
	ctx = lkctx.Context()
	defer lk.Unlock(lkctx.Cancel)

	if _, err = sys.getBucketMeta(bucket); err != nil {
		return err
	}

	if empty, err := sys.emptyBucket(ctx, bucket); err != nil {
		return err
	} else if !empty {
		return ErrBucketNotEmpty
	}

	return sys.db.Delete(bucketPrefix + bucket)
}

// GetAllBucketsOfUser metadata for all bucket.
func (sys *BucketMetadataSys) GetAllBucketsOfUser(ctx context.Context, username string) ([]BucketMetadata, error) {
	var m []BucketMetadata
	all, err := sys.db.ReadAllChan(ctx, bucketPrefix, "")
	if err != nil {
		return nil, err
	}
	for entry := range all {
		data := BucketMetadata{}
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
