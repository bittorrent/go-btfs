package responses

import (
	"github.com/aws/aws-sdk-go/aws"
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

func WriteListBucketsResponse(w http.ResponseWriter, r *http.Request, userId, username string, buckets []*object.Bucket) {
	resp := s3.ListBucketsOutput{
		Owner: &s3.Owner{
			ID:          aws.String(userId),
			DisplayName: aws.String(username),
		},
		Buckets: []*s3.Bucket{},
	}

	for _, buc := range buckets {
		resp.Buckets = append(resp.Buckets, &s3.Bucket{
			Name:         aws.String(buc.Name),
			CreationDate: aws.Time(buc.Created),
		})
	}

	WriteSuccessResponseXML(w, r, resp)

	return
}

func WriteGetBucketAclResponse(w http.ResponseWriter, r *http.Request, userId, username, acl string) {
	resp := s3.GetBucketAclOutput{
		Owner: &s3.Owner{
			ID:          aws.String(userId),
			DisplayName: aws.String(username),
		},
		Grants: []*s3.Grant{
			{
				Grantee: &s3.Grantee{
					ID:          aws.String(userId),
					DisplayName: aws.String(userId),
					Type:        aws.String("CanonicalUser"),
				},
				Permission: aws.String("public-read"),
			},
		},
	}
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
