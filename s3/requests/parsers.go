package requests

import (
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/responses"
	"net/http"
	"path"

	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/gorilla/mux"
)

//type PutObjectRequest struct {
//	Bucket string
//	Object string
//	Body   io.Reader
//}
//
//func (req *PutObjectRequest) Bind(r *http.Request) (err error) {
//	return
//}

func ParsePutBucketRequest(r *http.Request) (req *PutBucketRequest, rerr *responses.Error) {
	req = &PutBucketRequest{}
	req.User = cctx.GetAccessKey(r)
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

func ParseHeadBucketRequest(r *http.Request) (req *HeadBucketRequest, err error) {
	req = &HeadBucketRequest{}
	vars := mux.Vars(r)
	bucket := vars["bucket"]

	//set request
	req.Bucket = bucket
	return
}

// DeleteBucketRequest .
type DeleteBucketRequest struct {
	Bucket string
}

func ParseDeleteBucketRequest(r *http.Request) (req *DeleteBucketRequest, err error) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]

	//set request
	req = &DeleteBucketRequest{}
	req.Bucket = bucket
	return
}

// ListBucketsRequest .
type ListBucketsRequest struct {
	Bucket string
}

func ParseListBucketsRequest(r *http.Request) (req *ListBucketsRequest, err error) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]

	//set request
	req = &ListBucketsRequest{}
	req.Bucket = bucket
	return
}

// GetBucketAclRequest .
type GetBucketAclRequest struct {
	Bucket string
}

func ParseGetBucketAclRequest(r *http.Request) (req *GetBucketAclRequest, err error) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]

	//set request
	req = &GetBucketAclRequest{}
	req.Bucket = bucket
	return
}

// PutBucketAclRequest .
type PutBucketAclRequest struct {
	Bucket string
	ACL    string
}

func ParsePutBucketAclRequest(r *http.Request) (req *PutBucketAclRequest, err error) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]

	acl := r.Header.Get(consts.AmzACL)

	//set request
	req = &PutBucketAclRequest{}
	req.Bucket = bucket
	req.ACL = acl
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
