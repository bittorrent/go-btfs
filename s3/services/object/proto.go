package object

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"time"
)

var ErrNotFound = errors.New("object not found")

type Service interface {
	StoreObject(ctx context.Context, bucname, objname string, reader *hash.Reader, size int64, meta map[string]string) (obj Object, err error)
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
