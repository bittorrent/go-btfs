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
	CreateBucket(ctx context.Context, user, bucname, region, acl string) (bucket *Bucket, err error)
	GetBucket(ctx context.Context, user, bucname string) (bucket *Bucket, err error)
	DeleteBucket(ctx context.Context, user, bucname string) (err error)
	GetAllBuckets(ctx context.Context, user string) (list []*Bucket, err error)
	PutBucketAcl(ctx context.Context, user, bucname, acl string) (err error)
	GetBucketAcl(ctx context.Context, user, bucname string) (acl string, err error)
	EmptyBucket(ctx context.Context, user, bucname string) (empty bool, err error)

	PutObject(ctx context.Context, user, bucname, objname string, body *hash.Reader, size int64, meta map[string]string) (object *Object, err error)
	CopyObject(ctx context.Context, user, srcBucname, srcObjname, dstBucname, dstObjname string, meta map[string]string) (dstObject *Object, err error)
	GetObject(ctx context.Context, user, bucname, objname string) (object *Object, body io.ReadCloser, err error)
	DeleteObject(ctx context.Context, user, bucname, objname string) (err error)
	// todo: DeleteObjects
	ListObjects(ctx context.Context, user, bucname, prefix, delimiter, marker string, max int) (list *ObjectsList, err error)

	CreateMultipartUpload(ctx context.Context, user, bucname, objname string, meta map[string]string) (multipart *Multipart, err error)
	UploadPart(ctx context.Context, user, bucname, objname, uplid string, partId int, reader *hash.Reader, size int64, meta map[string]string) (part *ObjectPart, err error)
	AbortMultipartUpload(ctx context.Context, user, bucname, objname, uplid string) (err error)
	CompleteMultiPartUpload(ctx context.Context, user, bucname, objname, uplid string, parts []*CompletePart) (object *Object, err error)
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
	Parts     []*ObjectPart
}

type ObjectPart struct {
	ETag    string    `json:"etag,omitempty"`
	Cid     string    `json:"cid,omitempty"`
	Number  int       `json:"number"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

// ListObjectsInfo - container for list objects.
type ObjectsList struct {
	// Indicates whether the returned list objects response is truncated. A
	// value of true indicates that the list was truncated. The list can be truncated
	// if the number of objects exceeds the limit allowed or specified
	// by max keys.
	IsTruncated bool

	// When response is truncated (the IsTruncated element value in the response is true),
	// you can use the key name in this field as marker in the subsequent
	// request to get next set of objects.
	//
	// NOTE: AWS S3 returns NextMarker only if you have delimiter request parameter specified,
	NextMarker string

	// List of objects info for this request.
	Objects []*Object

	// List of prefixes for this request.
	Prefixes []string
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
