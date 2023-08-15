package handlers

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/apierrors"
	"github.com/bittorrent/go-btfs/s3/lock"
	"net/http"
)

type CorsService interface {
	GetAllowOrigins() []string
	GetAllowMethods() []string
	GetAllowHeaders() []string
}

type AccessKeyService interface {
	Generate() (record *AccessKeyRecord, err error)
	Enable(key string) (err error)
	Disable(key string) (err error)
	Reset(key string) (err error)
	Delete(key string) (err error)
	Get(key string) (record *AccessKeyRecord, err error)
	List() (list []*AccessKeyRecord, err error)
}

type AuthService interface {
	VerifySignature(ctx context.Context, r *http.Request) (accessKeyRecord *AccessKeyRecord, err apierrors.ErrorCode)
}

type BucketService interface {
	CheckACL(accessKeyRecord *AccessKeyRecord, bucketName string, action action.Action) (err error)
	NewNSLock(bucket string) lock.RWLocker
	SetEmptyBucket(emptyBucket func(ctx context.Context, bucket string) (bool, error))
	CreateBucket(ctx context.Context, bucket, region, accessKey, acl string) error
	GetBucketMeta(ctx context.Context, bucket string) (meta BucketMetadata, err error)
	HasBucket(ctx context.Context, bucket string) bool
	DeleteBucket(ctx context.Context, bucket string) error
	GetAllBucketsOfUser(ctx context.Context, username string) ([]BucketMetadata, error)
}

type ObjectService interface {
}

type MultipartService interface {
}
