package requests

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/responses"
	"net/http"
)

// CreateBucketRequest .
type CreateBucketRequest struct {
	AccessKey string
	Bucket    string
	ACL       string
	Region    string
}

func ParseCreateBucketRequest(r *http.Request) (req *CreateBucketRequest, rerr *responses.Error) {
	req = &CreateBucketRequest{}
	req.AccessKey = cctx.GetAccessKey(r)
	req.Bucket, rerr = ParseBucket(r)
	if rerr != nil {
		return
	}
	req.ACL, rerr = ParseBucketACL(r)
	if rerr != nil {
		return
	}
	req.Region, rerr = ParseLocation(r)
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
	req.Bucket, rerr = ParseBucket(r)
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
	req.Bucket, rerr = ParseBucket(r)
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
	req.Bucket, rerr = ParseBucket(r)
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
	req.Bucket, rerr = ParseBucket(r)
	if rerr != nil {
		return
	}
	req.ACL, rerr = ParseBucketACL(r)
	return
}

func ParsePutObjectRequest(r *http.Request) (req *s3.PutObjectInput, rerr *responses.Error) {
	err := responses.ParseRequest(r, &req)
	if err != nil {
		rerr = responses.ErrInvalidRequestParameter
		return
	}

	fmt.Printf("%+v", *req)
	return
}
