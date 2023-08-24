package bucket

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"time"
)

type Service interface {
	CheckACL(accessKeyRecord *accesskey.AccessKey, bucketName string, action action.Action) (err error)
	CreateBucket(ctx context.Context, bucket, region, accessKey, acl string) error
	GetBucketMeta(ctx context.Context, bucket string) (meta Bucket, err error)
	HasBucket(ctx context.Context, bucket string) bool
	SetEmptyBucket(emptyBucket func(ctx context.Context, bucket string) (bool, error))
	DeleteBucket(ctx context.Context, bucket string) error
	GetAllBucketsOfUser(username string) (list []*Bucket, err error)
	UpdateBucketAcl(ctx context.Context, bucket, acl string) error
	GetBucketAcl(ctx context.Context, bucket string) (string, error)
	EmptyBucket(ctx context.Context, bucket string) (bool, error)
}

// Bucket contains bucket metadata.
type Bucket struct {
	Name    string
	Region  string
	Owner   string
	Acl     string
	Created time.Time
}
