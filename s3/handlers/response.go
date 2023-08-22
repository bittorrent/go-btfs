package handlers

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/consts"
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

func WriteListBucketsResponse(w http.ResponseWriter, r *http.Request, bucketMetas []*BucketMetadata) {
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

func WriteGetBucketAclResponse(w http.ResponseWriter, r *http.Request, accessKeyRecord *AccessKeyRecord, acl string) {
	resp := AccessControlPolicy{}
	fmt.Printf(" -1- get acl resp: %+v \n", resp)

	id := accessKeyRecord.Key
	if resp.Owner.DisplayName == "" {
		resp.Owner.DisplayName = accessKeyRecord.Key
		resp.Owner.ID = id
	}
	fmt.Printf(" -2- get acl resp: %+v \n", resp)

	resp.AccessControlList.Grant = make([]Grant, 0)
	resp.AccessControlList.Grant = append(resp.AccessControlList.Grant, Grant{
		Grantee: Grantee{
			ID:          id,
			DisplayName: accessKeyRecord.Key,
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
