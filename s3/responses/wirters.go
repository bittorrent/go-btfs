package responses

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/services/bucket"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"net/http"
)

func WritePutBucketResponse(w http.ResponseWriter, r *http.Request) {
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

func WriteListBucketsResponse(w http.ResponseWriter, r *http.Request, bucketMetas []*bucket.Bucket) {
	var buckets []*s3.Bucket
	for _, b := range bucketMetas {
		buckets = append(buckets, &s3.Bucket{
			Name:         aws.String(b.Name),
			CreationDate: aws.Time(b.Created),
		})
	}

	resp := ListAllMyBucketsResult{
		Owner: &s3.Owner{
			ID:          aws.String(consts.DefaultOwnerID),
			DisplayName: aws.String(consts.DisplayName),
		},
		Buckets: buckets,
	}

	WriteSuccessResponseXML(w, r, resp)
	return
}

func WriteGetBucketAclResponse(w http.ResponseWriter, r *http.Request, key string, acl string) {
	resp := GetBucketAclResponse{}
	fmt.Printf(" -1- get acl resp: %+v \n", resp)

	id := key
	if resp.Owner.DisplayName == "" {
		resp.Owner.DisplayName = key
		resp.Owner.ID = id
	}
	fmt.Printf(" -2- get acl resp: %+v \n", resp)

	resp.AccessControlList.Grant = make([]Grant, 0)
	resp.AccessControlList.Grant = append(resp.AccessControlList.Grant, Grant{
		Grantee: Grantee{
			ID:          id,
			DisplayName: key,
			Type:        "CanonicalUser",
			XMLXSI:      "CanonicalUser",
			XMLNS:       "http://www.w3.org/2001/XMLSchema-instance"},
		Permission: Permission(acl), //todo change
	})
	fmt.Printf(" -3- get acl resp: %+v \n", resp)

	fmt.Printf("get acl resp: %+v \n", resp)

	WriteSuccessResponseXML(w, r, resp)
	return
}

func WritePutBucketAclResponse(w http.ResponseWriter, r *http.Request) {
	WriteSuccessResponse(w, r)
	return
}

func WritePutObjectResponse(w http.ResponseWriter, r *http.Request, obj object.Object) {
	setPutObjHeaders(w, obj.ETag, obj.Cid, false)
	WriteSuccessResponseHeadersOnly(w, r)
}

func WriteCreateMultipartUploadResponse(w http.ResponseWriter, r *http.Request, bucname, objname, uploadID string) {
	resp := GenerateInitiateMultipartUploadResponse(bucname, objname, uploadID)
	WriteSuccessResponseXML(w, r, resp)
}

func WriteAbortMultipartUploadResponse(w http.ResponseWriter, r *http.Request) {
	WriteSuccessNoContent(w)
}

func WriteUploadPartResponse(w http.ResponseWriter, r *http.Request, part object.ObjectPart) {
	setPutObjHeaders(w, part.ETag, part.Cid, false)
	WriteSuccessResponseHeadersOnly(w, r)
}

func WriteCompleteMultipartUploadResponse(w http.ResponseWriter, r *http.Request, bucname, objname, region string, obj object.Object) {
	resp := GenerateCompleteMultipartUploadResponse(bucname, objname, region, obj)
	setPutObjHeaders(w, obj.ETag, obj.Cid, false)
	WriteSuccessResponseXML(w, r, resp)
}