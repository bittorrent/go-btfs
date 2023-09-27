package object

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/hash"
	"io"
	"time"
)

var (
	ErrBucketNotFound      = errors.New("bucket not found")
	ErrBucketNotEmpty      = errors.New("bucket not empty")
	ErrObjectNotFound      = errors.New("object not found")
	ErrUploadNotFound      = errors.New("upload not found")
	ErrNotAllowed          = errors.New("authentication not allowed")
	ErrBucketAlreadyExists = errors.New("bucket already exists")
	ErrPartNotExists       = errors.New("part not exists")
	ErrPartETagNotMatch    = errors.New("part etag not match")
	ErrPartTooSmall        = errors.New("part size too small")
	ErrCanceled            = context.Canceled
	ErrTimout              = context.DeadlineExceeded
)

type Service interface {
	CreateBucket(ctx context.Context, args *CreateBucketArgs) (bucket *Bucket, err error)
	GetBucket(ctx context.Context, args *GetBucketArgs) (bucket *Bucket, err error)
	DeleteBucket(ctx context.Context, args *DeleteBucketArgs) (err error)
	ListBuckets(ctx context.Context, args *ListBucketsArgs) (list *BucketList, err error)
	PutBucketACL(ctx context.Context, args *PutBucketACLArgs) (err error)
	GetBucketACL(ctx context.Context, args *GetBucketACLArgs) (acl *ACL, err error)

	PutObject(ctx context.Context, args *PutObjectArgs) (object *Object, err error)
	CopyObject(ctx context.Context, args *CopyObjectArgs) (object *Object, err error)
	GetObject(ctx context.Context, args *GetObjectArgs) (object *Object, body io.ReadCloser, err error)
	DeleteObject(ctx context.Context, args *DeleteObjectArgs) (err error)
	DeleteObjects(ctx context.Context, args *DeleteObjectsArgs) (deletes []*DeletedObject, err error)
	ListObjects(ctx context.Context, args *ListObjectsArgs) (list *ObjectsList, err error)
	ListObjectsV2(ctx context.Context, args *ListObjectsV2Args) (list *ObjectsListV2, err error)
	GetObjectACL(ctx context.Context, args *GetObjectACLArgs) (acl *ACL, err error)

	CreateMultipartUpload(ctx context.Context, args *CreateMultipartUploadArgs) (multipart *Multipart, err error)
	UploadPart(ctx context.Context, args *UploadPartArgs) (part *Part, err error)
	AbortMultipartUpload(ctx context.Context, args *AbortMultipartUploadArgs) (err error)
	CompleteMultiPartUpload(ctx context.Context, args *CompleteMultipartUploadArgs) (object *Object, err error)
}

type CreateBucketArgs struct {
	UserId string
	ACL    string
	Bucket string
	Region string
}

type GetBucketArgs struct {
	UserId string
	Bucket string
}

type DeleteBucketArgs struct {
	UserId string
	Bucket string
}

type ListBucketsArgs struct {
	UserId string
}

type GetBucketACLArgs struct {
	UserId string
	Bucket string
}

type PutBucketACLArgs struct {
	UserId string
	ACL    string
	Bucket string
}

type PutObjectArgs struct {
	UserId          string
	Body            *hash.Reader
	Bucket          string
	Object          string
	ContentEncoding string
	ContentLength   int64
	ContentType     string
	Expires         time.Time
}

type CopyObjectArgs struct {
	UserId          string
	Bucket          string
	Object          string
	SrcBucket       string
	SrcObject       string
	ContentEncoding string
	ContentType     string
	Expires         time.Time
	ReplaceMeta     bool
}

type GetObjectArgs struct {
	UserId   string
	Bucket   string
	Object   string
	WithBody bool
}

type DeleteObjectArgs struct {
	UserId string
	Bucket string
	Object string
}

type DeleteObjectsArgs struct {
	UserId          string
	Bucket          string
	ToDeleteObjects []*ToDeleteObject
	Quite           bool
}

type ToDeleteObject struct {
	Object      string
	ValidateErr error
}

type ListObjectsArgs struct {
	UserId       string
	Bucket       string
	MaxKeys      int64
	Marker       string
	Prefix       string
	Delimiter    string
	EncodingType string
}

type ListObjectsV2Args struct {
	UserId       string
	Bucket       string
	MaxKeys      int64
	Prefix       string
	Delimiter    string
	EncodingType string
	Token        string
	After        string
	FetchOwner   bool
}

type GetObjectACLArgs struct {
	UserId string
	Bucket string
	Object string
}

type CreateMultipartUploadArgs struct {
	UserId          string
	Bucket          string
	Object          string
	ContentEncoding string
	ContentType     string
	Expires         time.Time
}

type UploadPartArgs struct {
	UserId        string
	Body          *hash.Reader
	Bucket        string
	Object        string
	UploadId      string
	PartNumber    int64
	ContentLength int64
}

type AbortMultipartUploadArgs struct {
	UserId   string
	Body     *hash.Reader
	Bucket   string
	Object   string
	UploadId string
}

type CompleteMultipartUploadArgs struct {
	UserId         string
	Bucket         string
	Object         string
	UploadId       string
	CompletedParts CompletedParts
}

type ACL struct {
	Owner string
	ACL   string
}

type DeletedObject struct {
	Object    string
	DeleteErr error
}

type Bucket struct {
	Name    string
	Region  string
	Owner   string
	ACL     string
	Created time.Time
}

type BucketList struct {
	Owner   string
	Buckets []*Bucket
}

type Object struct {
	Bucket           string
	Name             string
	ModTime          time.Time
	Size             int64
	IsDir            bool
	ETag             string
	CID              string
	ACL              string
	VersionID        string
	IsLatest         bool
	DeleteMarker     bool
	ContentType      string
	ContentEncoding  string
	Expires          time.Time
	AccTime          time.Time
	SuccessorModTime time.Time
}

type Multipart struct {
	Bucket          string
	Object          string
	UploadID        string
	Initiated       time.Time
	ContentType     string
	ContentEncoding string
	Expires         time.Time
	Parts           []*Part
}

type Part struct {
	ETag    string    `json:"etag,omitempty"`
	CID     string    `json:"cid,omitempty"`
	Number  int64     `json:"number"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

type ObjectsList struct {
	Args        *ListObjectsArgs
	Owner       string
	IsTruncated bool
	NextMarker  string
	Objects     []*Object
	Prefixes    []string
}

type ObjectsListV2 struct {
	Args                  *ListObjectsV2Args
	Owner                 string
	IsTruncated           bool
	NextContinuationToken string
	Objects               []*Object
	Prefixes              []string
}

type CompletePart struct {
	PartNumber int64
	ETag       string
}

type CompletedParts []*CompletePart

func (a CompletedParts) Len() int           { return len(a) }
func (a CompletedParts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CompletedParts) Less(i, j int) bool { return a[i].PartNumber < a[j].PartNumber }
