package store

import (
	"context"
	"github.com/bittorrent/go-btfs/s3d/lock"
)

type Service interface {
	NewNSLock(bucket string) lock.RWLocker
	SetEmptyBucket(emptyBucket func(ctx context.Context, bucket string) (bool, error))
	CreateBucket(ctx context.Context, bucket, region, accessKey, acl string) error
	GetBucketMeta(ctx context.Context, bucket string) (meta BucketMetadata, err error)
	HasBucket(ctx context.Context, bucket string) bool
	DeleteBucket(ctx context.Context, bucket string) error
	GetAllBucketsOfUser(ctx context.Context, username string) ([]BucketMetadata, error)
}
