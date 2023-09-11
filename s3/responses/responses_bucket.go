package responses

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/protocol"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"net/http"
)

func WriteCreateBucketResponse(w http.ResponseWriter, r *http.Request, buc *object.Bucket) {
	output := new(s3.CreateBucketOutput).SetLocation(pathClean(r.URL.Path))
	w.Header().Add(consts.AmzACL, buc.ACL)
	WriteSuccessResponse(w, output, "")
	return
}

func WriteHeadBucketResponse(w http.ResponseWriter, r *http.Request, buc *object.Bucket) {
	output := new(s3.HeadBucketOutput)
	w.Header().Add(consts.AmzACL, buc.ACL)
	WriteSuccessResponse(w, output, "")
	return
}

func WriteDeleteBucketResponse(w http.ResponseWriter) {
	output := new(s3.DeleteBucketOutput)
	_ = protocol.WriteResponse(w, http.StatusOK, output, "")
	return
}

func WriteListBucketsResponse(w http.ResponseWriter, r *http.Request, accessKey string, buckets []*object.Bucket) {
	output := new(s3.ListBucketsOutput)
	output.SetOwner(owner(accessKey))
	s3Buckets := make([]*s3.Bucket, 0)
	for _, buc := range buckets {
		s3Bucket := new(s3.Bucket).SetName(buc.Name).SetCreationDate(buc.Created)
		s3Buckets = append(s3Buckets, s3Bucket)
		w.Header().Add(consts.AmzACL, buc.ACL)
	}
	output.SetBuckets(s3Buckets)
	WriteSuccessResponse(w, output, "ListAllMyBucketsResult")
	return
}

func WritePutBucketAclResponse(w http.ResponseWriter, r *http.Request) {
	output := new(s3.PutBucketAclOutput)
	WriteSuccessResponse(w, output, "")
	return
}

func WriteGetBucketACLResponse(w http.ResponseWriter, r *http.Request, accessKey string, acl string) {
	output := new(s3.GetBucketAclOutput)
	output.SetOwner(owner(accessKey))
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
	output.SetGrants(grants)
	w.Header().Add(consts.AmzACL, acl)
	WriteSuccessResponse(w, output, "AccessControlPolicy")
	return
}
