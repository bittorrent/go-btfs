package services

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/handlers/responses"
	"net/http"
	"time"

	"errors"
	"github.com/bittorrent/go-btfs/s3/action"
)

type CorsService interface {
	GetAllowOrigins() []string
	GetAllowMethods() []string
	GetAllowHeaders() []string
}

type AccessKeyService interface {
	Generate() (record *AccessKey, err error)
	Enable(key string) (err error)
	Disable(key string) (err error)
	Reset(key string) (err error)
	Delete(key string) (err error)
	Get(key string) (ack *AccessKey, err error)
	List() (list []*AccessKey, err error)
}

type AuthService interface {
	VerifySignature(ctx context.Context, r *http.Request) (ack *AccessKey, rErr *responses.Error)
}

type BucketService interface {
	CheckACL(accessKeyRecord *AccessKey, bucketName string, action action.Action) (err error)
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

var (
	ErrSignVersionNotSupport = errors.New("sign version is not support")
	ErrInvalidArgument       = errors.New("invalid argument")
	ErrInvalidBucketName     = errors.New("bucket name is invalid")
	ErrBucketNotFound        = errors.New("bucket is not found")
	ErrBucketAlreadyExists   = errors.New("bucket is already exists")
	ErrBucketAccessDenied    = errors.New("bucket access denied")
	ErrSetBucketEmptyFailed  = errors.New("set bucket empty failed")
	ErrCreateBucket          = errors.New("create bucket failed")
	ErrNoSuchUserPolicy      = errors.New("no such user policy")
	ErrNotImplemented        = errors.New("not implemented")
	ErrAccessKeyIsNotFound   = errors.New("access-key is not found")
)

type AccessKey struct {
	Key       string    `json:"key"`
	Secret    string    `json:"secret"`
	Enable    bool      `json:"enable"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BucketMetadata contains bucket metadata.
type BucketMetadata struct {
	Name    string
	Region  string
	Owner   string
	Acl     string
	Created time.Time
}

type ObjectMetadata struct {
}