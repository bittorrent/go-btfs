package responses

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/api/services/object"
	"github.com/bittorrent/go-btfs/s3/consts"
	"net/http"
)

func newS3Owner(userId string) *s3.Owner {
	return new(s3.Owner).SetID(userId).SetDisplayName(userId)
}

func newS3FullControlGrant(userId string) *s3.Grant {
	return new(s3.Grant).SetGrantee(new(s3.Grantee).SetType(s3.TypeCanonicalUser).SetID(userId).SetDisplayName(userId)).SetPermission(s3.PermissionFullControl)
}

var (
	s3AllUsersReadGrant  = new(s3.Grant).SetGrantee(new(s3.Grantee).SetType(s3.TypeGroup).SetURI(consts.AllUsersURI)).SetPermission(s3.PermissionRead)
	s3AllUsersWriteGrant = new(s3.Grant).SetGrantee(new(s3.Grantee).SetType(s3.TypeGroup).SetURI(consts.AllUsersURI)).SetPermission(s3.PermissionWrite)
)

func WriteCreateBucketResponse(w http.ResponseWriter, r *http.Request, buc *object.Bucket) {
	output := new(s3.CreateBucketOutput).SetLocation(r.URL.Path)
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
	WriteSuccessResponse(w, output, "")
	return
}

func WriteListBucketsResponse(w http.ResponseWriter, r *http.Request, list *object.BucketList) {
	output := new(s3.ListBucketsOutput)
	output.SetOwner(newS3Owner(list.Owner))
	s3Buckets := make([]*s3.Bucket, 0)
	for _, buc := range list.Buckets {
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

func WriteGetBucketACLResponse(w http.ResponseWriter, r *http.Request, acl *object.ACL) {
	output := new(s3.GetBucketAclOutput)
	output.SetOwner(newS3Owner(acl.Owner))
	grants := make([]*s3.Grant, 0)
	grants = append(grants, newS3FullControlGrant(acl.Owner))
	switch acl.ACL {
	case s3.BucketCannedACLPrivate:
	case s3.BucketCannedACLPublicRead:
		grants = append(grants, s3AllUsersReadGrant)
	case s3.BucketCannedACLPublicReadWrite:
		grants = append(grants, s3AllUsersReadGrant, s3AllUsersWriteGrant)
	default:
		panic(fmt.Sprintf("unknwon acl <%s>", acl.ACL))
	}
	output.SetGrants(grants)
	w.Header().Add(consts.AmzACL, acl.ACL)
	WriteSuccessResponse(w, output, "AccessControlPolicy")
	return
}
