package bucket

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/action"
	"time"
)

var ErrNotFound = errors.New("bucket not found")

type Service interface {
	CheckACL(accessKey string, bucketName string, action action.Action) (err error)
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
