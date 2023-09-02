package responses

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"net/http"
)

func WritePutBucketResponse(w http.ResponseWriter, r *http.Request) {
	if cp := pathClean(r.URL.Path); cp != "" {
		w.Header().Set(consts.Location, cp)
	}
	WriteSuccessResponse(w, r)
	return
}

func WriteHeadBucketResponse(w http.ResponseWriter, r *http.Request) {
	WriteSuccessResponse(w, r)
	return
}

func WriteDeleteBucketResponse(w http.ResponseWriter) {
	WriteSuccessNoContent(w)
	return
}

type ListBucketResponse struct {
	ListAllMyBucketsResult s3.ListBucketsOutput `xml:"ListAllMyBucketsResult"`
}

func WriteListBucketsResponse(w http.ResponseWriter, r *http.Request, accessKey string, buckets []*object.Bucket) {
	resp := &ListBucketResponse{}
	resp.ListAllMyBucketsResult.SetOwner(owner(accessKey))
	s3Buckets := make([]*s3.Bucket, 0)
	for _, buc := range buckets {
		s3Bucket := new(s3.Bucket).SetName(buc.Name).SetCreationDate(buc.Created)
		s3Buckets = append(s3Buckets, s3Bucket)
	}
	resp.ListAllMyBucketsResult.SetBuckets(s3Buckets)
	WriteSuccessResponseXML(w, r, resp)
	return
}

func WritePutBucketAclResponse(w http.ResponseWriter, r *http.Request) {
	WriteSuccessResponse(w, r)
	return
}

type GetBucketACLResponse struct {
	AccessControlPolicy s3.GetBucketAclOutput `xml:"AccessControlPolicy"`
}

func WriteGetBucketACLResponse(w http.ResponseWriter, r *http.Request, accessKey string, acl string) {
	resp := GetBucketACLResponse{}
	resp.AccessControlPolicy.SetOwner(owner(accessKey))
	grants := make([]*s3.Grant, 0)
	grants = append(grants, ownerFullControlGrant(accessKey))
	switch acl {
	case s3.BucketCannedACLPrivate:
	case s3.BucketCannedACLPublicRead:
		grants = append(grants, allUsersReadGrant)
	case s3.BucketCannedACLPublicReadWrite:
		grants = append(grants, allUsersReadGrant, allUsersWriteGrant)
	default:
		panic("unknown acl")
	}
	resp.AccessControlPolicy.SetGrants(grants)
	WriteSuccessResponseXML(w, r, resp)
	return
}
