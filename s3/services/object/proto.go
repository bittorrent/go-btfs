package object

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"io"
	"time"
)

var (
	ErrBucketNotFound      = errors.New("bucket not found")
	ErrObjectNotFound      = errors.New("object not found")
	ErrUploadNotFound      = errors.New("upload not found")
	ErrNotAllowed          = errors.New("not allowed")
	ErrBucketAlreadyExists = errors.New("bucket already exists")
)

type Service interface {
	// bucket
	CreateBucket(ctx context.Context, bucket, region, accessKey, acl string) error
	GetBucketMeta(ctx context.Context, bucket string) (meta Bucket, err error)
	HasBucket(ctx context.Context, bucket string) bool
	SetEmptyBucket(emptyBucket func(ctx context.Context, bucket string) (bool, error))
	DeleteBucket(ctx context.Context, bucket string) error
	GetAllBucketsOfUser(username string) (list []*Bucket, err error)
	UpdateBucketAcl(ctx context.Context, bucket, acl string) error
	GetBucketAcl(ctx context.Context, bucket string) (string, error)
	EmptyBucket(emptyBucket func(ctx context.Context, bucket string) (bool, error))

	// object
	PutObject(ctx context.Context, bucname, objname string, reader *hash.Reader, size int64, meta map[string]string) (obj Object, err error)
	CopyObject(ctx context.Context, bucket, object string, info Object, size int64, meta map[string]string) (Object, error)
	GetObject(ctx context.Context, bucket, object string) (Object, io.ReadCloser, error)
	GetObjectInfo(ctx context.Context, bucket, object string) (Object, error)
	DeleteObject(ctx context.Context, bucket, object string) error
	ListObjects(ctx context.Context, bucket string, prefix string, marker string, delimiter string, maxKeys int) (loi ListObjectsInfo, err error)
	ListObjectsV2(ctx context.Context, bucket string, prefix string, continuationToken string, delimiter string, maxKeys int, owner bool, startAfter string) (ListObjectsV2Info, error)

	// martipart
	CreateMultipartUpload(ctx context.Context, bucname string, objname string, meta map[string]string) (mtp Multipart, err error)
	AbortMultipartUpload(ctx context.Context, bucname string, objname string, uploadID string) (err error)
	UploadPart(ctx context.Context, bucname string, objname string, uploadID string, partID int, reader *hash.Reader, size int64, meta map[string]string) (part ObjectPart, err error)
	CompleteMultiPartUpload(ctx context.Context, bucname string, objname string, uploadID string, parts []CompletePart) (obj Object, err error)
	GetMultipart(ctx context.Context, bucname string, objname string, uploadID string) (mtp Multipart, err error)
}

// Bucket contains bucket metadata.
type Bucket struct {
	Name    string
	Region  string
	Owner   string
	Acl     string
	Created time.Time
}

type Object struct {
	Bucket           string
	Name             string
	ModTime          time.Time
	Size             int64
	IsDir            bool
	ETag             string
	Cid              string
	Acl              string
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
	Bucket    string
	Object    string
	UploadID  string
	Initiated time.Time
	MetaData  map[string]string
	Parts     []ObjectPart
}

type ObjectPart struct {
	ETag    string    `json:"etag,omitempty"`
	Cid     string    `json:"cid,omitempty"`
	Number  int       `json:"number"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

type CompletePart struct {
	PartNumber     int
	ETag           string
	ChecksumCRC32  string
	ChecksumCRC32C string
	ChecksumSHA1   string
	ChecksumSHA256 string
}

type CompletedParts []CompletePart

func (a CompletedParts) Len() int           { return len(a) }
func (a CompletedParts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CompletedParts) Less(i, j int) bool { return a[i].PartNumber < a[j].PartNumber }

type CompleteMultipartUpload struct {
	Parts []CompletePart `xml:"Part"`
}
