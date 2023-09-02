package requests

import (
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/responses"
	"net/http"
	"path"
)

// PutBucketRequest .
type PutBucketRequest struct {
	AccessKey string
	Bucket    string
	ACL       string
	Region    string
}

func ParsePutBucketRequest(r *http.Request) (req *PutBucketRequest, rerr *responses.Error) {
	req = &PutBucketRequest{}
	req.AccessKey = cctx.GetAccessKey(r)
	req.Bucket, rerr = parseBucket(r)
	if rerr != nil {
		return
	}
	req.ACL, rerr = parseAcl(r)
	if rerr != nil {
		return
	}
	req.Region, rerr = parseLocationConstraint(r)
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

// GetBucketAclRequest .
type GetBucketAclRequest struct {
	AccessKey string
	Bucket    string
}

func ParseGetBucketAclRequest(r *http.Request) (req *GetBucketAclRequest, rerr *responses.Error) {
	req = &GetBucketAclRequest{}
	req.AccessKey = cctx.GetAccessKey(r)
	req.Bucket, rerr = parseBucket(r)
	return
}

// PutBucketAclRequest .
type PutBucketAclRequest struct {
	AccessKey string
	Bucket    string
	ACL       string
}

func ParsePutBucketAclRequest(r *http.Request) (req *PutBucketAclRequest, rerr *responses.Error) {
	req = &PutBucketAclRequest{}
	req.AccessKey = cctx.GetAccessKey(r)
	req.Bucket, rerr = parseBucket(r)
	if rerr != nil {
		return
	}
	req.ACL, rerr = parseAcl(r)
	return
}

// pathClean is like path.Clean but does not return "." for
// empty inputs, instead returns "empty" as is.
func PathClean(p string) string {
	cp := path.Clean(p)
	if cp == "." {
		return ""
	}
	return cp
}

//func unmarshalXML(reader io.Reader, isObject bool) (*store.Tags, error) {
//	tagging := &store.Tags{
//		TagSet: &store.TagSet{
//			TagMap:   make(map[string]string),
//			IsObject: isObject,
//		},
//	}
//
//	if err := xml.NewDecoder(reader).Decode(tagging); err != nil {
//		return nil, err
//	}
//
//	return tagging, nil
//}

func checkAcl(acl string) (ok bool) {
	_, ok = supportAcls[acl]
	return
}
