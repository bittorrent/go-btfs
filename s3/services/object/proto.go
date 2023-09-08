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
	ErrBucketeNotEmpty     = errors.New("bucket not empty")
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
	PutBucketACL(ctx context.Context, user, bucname, acl string) (err error)
	GetBucketACL(ctx context.Context, user, bucname string) (acl string, err error)

	PutObject(ctx context.Context, user, bucname, objname string, body *hash.Reader, size int64, meta map[string]string) (object *Object, err error)
	CopyObject(ctx context.Context, user, srcBucname, srcObjname, dstBucname, dstObjname string, meta map[string]string) (dstObject *Object, err error)
	GetObject(ctx context.Context, user, bucname, objname string, withBody bool) (object *Object, body io.ReadCloser, err error)
	DeleteObject(ctx context.Context, user, bucname, objname string) (err error)
	ListObjects(ctx context.Context, user, bucname, prefix, delimiter, marker string, max int64) (list *ObjectsList, err error)
	ListObjectsV2(ctx context.Context, user string, bucket string, prefix string, token, delimiter string, max int64, owner bool, after string) (list *ObjectsListV2, err error)

	CreateMultipartUpload(ctx context.Context, user, bucname, objname string, meta map[string]*string) (multipart *Multipart, err error)
	UploadPart(ctx context.Context, user, bucname, objname, uplid string, partId int, reader *hash.Reader, size int64) (part *Part, err error)
	AbortMultipartUpload(ctx context.Context, user, bucname, objname, uplid string) (err error)
	CompleteMultiPartUpload(ctx context.Context, user, bucname, objname, uplid string, parts []*CompletePart) (object *Object, err error)
}

// Bucket contains bucket metadata.
type Bucket struct {
	Name    string
	Region  string
	Owner   string
	ACL     string
	Created time.Time
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
	Bucket    string
	Object    string
	UploadID  string
	Initiated time.Time
	MetaData  map[string]*string
	Parts     []*Part
}

type Part struct {
	ETag    string    `json:"etag,omitempty"`
	CID     string    `json:"cid,omitempty"`
	Number  int       `json:"number"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

type ObjectsList struct {
	IsTruncated bool
	NextMarker  string
	Objects     []*Object
	Prefixes    []string
}

type ObjectsListV2 struct {
	IsTruncated           bool
	ContinuationToken     string
	NextContinuationToken string
	Objects               []*Object
	Prefixes              []string
}

type CompletePart struct {
	PartNumber     int
	ETag           string
	ChecksumCRC32  string
	ChecksumCRC32C string
	ChecksumSHA1   string
	ChecksumSHA256 string
}

type CompletedParts []*CompletePart

func (a CompletedParts) Len() int           { return len(a) }
func (a CompletedParts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CompletedParts) Less(i, j int) bool { return a[i].PartNumber < a[j].PartNumber }

type CompleteMultipartUpload struct {
	Parts []*CompletePart `xml:"Part"`
}
