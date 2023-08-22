package handlers

import (
	"context"
	"net/http"

	"github.com/bittorrent/go-btfs/s3/action"
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
	VerifySignature(ctx context.Context, r *http.Request) (accessKeyRecord *AccessKeyRecord, err ErrorCode)
}

type BucketService interface {
	CheckACL(accessKeyRecord *AccessKeyRecord, bucketName string, action action.Action) (err error)
	CreateBucket(ctx context.Context, bucket, region, accessKey, acl string) error
	GetBucketMeta(ctx context.Context, bucket string) (meta BucketMetadata, err error)
	HasBucket(ctx context.Context, bucket string) bool
	SetEmptyBucket(emptyBucket func(ctx context.Context, bucket string) (bool, error))
	DeleteBucket(ctx context.Context, bucket string) error
	GetAllBucketsOfUser(username string) (list []*BucketMetadata, err error)
	UpdateBucketAcl(ctx context.Context, bucket, acl string) error
	GetBucketAcl(ctx context.Context, bucket string) (string, error)
}

type ObjectService interface {
}

type MultipartService interface {
}
