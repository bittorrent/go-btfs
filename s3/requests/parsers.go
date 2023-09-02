package requests

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/responses"
	"net/http"
	"reflect"
)

// CreateBucketRequest .
type CreateBucketRequest struct {
	AccessKey string
	Bucket    string
	ACL       string
	Region    string
}

// todo: parse aws request use aws struct
func ParseS3Request(r *http.Request, v interface{}) (err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		err = errors.New("invalid value must be non nil pointer")
		return
	}

	rt := reflect.TypeOf(v).Elem()
	n := rt.NumField()
	for i := 0; i < n; i++ {
		f := rt.Field(i)
		fmt.Println(f)
	}
	return
}

func ParseCreateBucketRequest(r *http.Request) (req *CreateBucketRequest, rerr *responses.Error) {
	req = &CreateBucketRequest{}
	req.AccessKey = cctx.GetAccessKey(r)
	req.Bucket, rerr = parseBucket(r)
	if rerr != nil {
		return
	}
	req.ACL, rerr = parseBucketACL(r)
	if rerr != nil {
		return
	}
	req.Region, rerr = parseLocation(r)
	return
}

// DeleteBucketRequest .
type DeleteBucketRequest struct {
	AccessKey string
	Bucket    string
}

func ParseDeleteBucketRequest(r *http.Request) (req *DeleteBucketRequest, rerr *responses.Error) {
	req = &DeleteBucketRequest{}
	req.AccessKey = cctx.GetAccessKey(r)
	req.Bucket, rerr = parseBucket(r)
	return
}

// HeadBucketRequest .
type HeadBucketRequest struct {
	AccessKey string
	Bucket    string
}

func ParseHeadBucketRequest(r *http.Request) (req *HeadBucketRequest, rerr *responses.Error) {
	req = &HeadBucketRequest{}
	req.AccessKey = cctx.GetAccessKey(r)
	req.Bucket, rerr = parseBucket(r)
	return
}

// ListBucketsRequest .
type ListBucketsRequest struct {
	AccessKey string
}

func ParseListBucketsRequest(r *http.Request) (req *ListBucketsRequest, rerr *responses.Error) {
	req = &ListBucketsRequest{}
	req.AccessKey = cctx.GetAccessKey(r)
	return
}

// GetBucketACLRequest .
type GetBucketACLRequest struct {
	AccessKey string
	Bucket    string
}

func ParseGetBucketACLRequest(r *http.Request) (req *GetBucketACLRequest, rerr *responses.Error) {
	req = &GetBucketACLRequest{}
	req.AccessKey = cctx.GetAccessKey(r)
	req.Bucket, rerr = parseBucket(r)
	return
}

// PutBucketACLRequest .
type PutBucketACLRequest struct {
	AccessKey string
	Bucket    string
	ACL       string
}

func ParsePutBucketAclRequest(r *http.Request) (req *PutBucketACLRequest, rerr *responses.Error) {
	req = &PutBucketACLRequest{}
	req.AccessKey = cctx.GetAccessKey(r)
	req.Bucket, rerr = parseBucket(r)
	if rerr != nil {
		return
	}
	req.ACL, rerr = parseBucketACL(r)
	return
}
