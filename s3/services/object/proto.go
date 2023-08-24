package object

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"time"
)

var (
	ErrObjectNotFound = errors.New("object not found")
	ErrUploadNotFound = errors.New("upload not found")
)

type Service interface {
	PutObject(ctx context.Context, bucname, objname string, reader *hash.Reader, size int64, meta map[string]string) (obj Object, err error)

	// martipart
	CreateMultipartUpload(ctx context.Context, bucname string, objname string, meta map[string]string) (mtp Multipart, err error)
	AbortMultipartUpload(ctx context.Context, bucname string, objname string, uploadID string) (err error)
	UploadPart(ctx context.Context, bucname string, objname string, uploadID string, partID int, reader *hash.Reader, size int64, meta map[string]string) (part ObjectPart, err error)
	CompleteMultiPartUpload(ctx context.Context, bucname string, objname string, uploadID string, parts []CompletePart) (obj Object, err error)
	GetMultipart(ctx context.Context, bucname string, objname string, uploadID string) (mtp Multipart, err error)
}

type Object struct {
	// Name of the bucket.
	Bucket string

	// Name of the object.
	Name string

	// Date and time when the object was last modified.
	ModTime time.Time

	// Total object size.
	Size int64

	// IsDir indicates if the object is prefix.
	IsDir bool

	// Hex encoded unique entity tag of the object.
	ETag string

	// ipfs key
	Cid string
	Acl string
	// Version ID of this object.
	VersionID string

	// IsLatest indicates if this is the latest current version
	// latest can be true for delete marker or a version.
	IsLatest bool

	// DeleteMarker indicates if the versionId corresponds
	// to a delete marker on an object.
	DeleteMarker bool

	// A standard MIME type describing the format of the object.
	ContentType string

	// Specifies what content encodings have been applied to the object and thus
	// what decoding mechanisms must be applied to obtain the object referenced
	// by the Content-Type header field.
	ContentEncoding string

	// Date and time at which the object is no longer able to be cached
	Expires time.Time

	// Date and time when the object was last accessed.
	AccTime time.Time

	//  The mod time of the successor object version if any
	SuccessorModTime time.Time
}

type Multipart struct {
	Bucket    string
	Object    string
	UploadID  string
	Initiated time.Time
	MetaData  map[string]string
	// List of individual parts, maximum size of upto 10,000
	Parts []ObjectPart
}

// objectPartInfo Info of each part kept in the multipart metadata
// file after CompleteMultipartUpload() is called.
type ObjectPart struct {
	ETag    string    `json:"etag,omitempty"`
	Cid     string    `json:"cid,omitempty"`
	Number  int       `json:"number"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

// CompletePart - represents the part that was completed, this is sent by the client
// during CompleteMultipartUpload request.
type CompletePart struct {
	// Part number identifying the part. This is a positive integer between 1 and
	// 10,000
	PartNumber int

	// Entity tag returned when the part was uploaded.
	ETag string

	// Checksum values. Optional.
	ChecksumCRC32  string
	ChecksumCRC32C string
	ChecksumSHA1   string
	ChecksumSHA256 string
}

// CompletedParts - is a collection satisfying sort.Interface.
type CompletedParts []CompletePart

func (a CompletedParts) Len() int           { return len(a) }
func (a CompletedParts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CompletedParts) Less(i, j int) bool { return a[i].PartNumber < a[j].PartNumber }

// CompleteMultipartUpload - represents list of parts which are completed, this is sent by the
// client during CompleteMultipartUpload request.
type CompleteMultipartUpload struct {
	Parts []CompletePart `xml:"Part"`
}
