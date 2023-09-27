package responses

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/api/services/object"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/utils"
	"io"
	"net/http"
)

func WritePutObjectResponse(w http.ResponseWriter, r *http.Request, obj *object.Object) {
	output := new(s3.PutObjectOutput)
	output.SetETag(`"` + obj.ETag + `"`)
	w.Header().Set(consts.Cid, obj.CID)
	WriteSuccessResponse(w, output, "")
}

func WriteCopyObjectResponse(w http.ResponseWriter, r *http.Request, obj *object.Object) {
	output := new(s3.CopyObjectResult)
	output.SetETag(`"` + obj.ETag + `"`)
	output.SetLastModified(obj.ModTime)
	w.Header().Set(consts.Cid, obj.CID)
	WriteSuccessResponse(w, output, "CopyObjectResult")
}

func WriteHeadObjectResponse(w http.ResponseWriter, r *http.Request, obj *object.Object) {
	output := new(s3.HeadObjectOutput)
	output.SetETag(`"` + obj.ETag + `"`)
	output.SetLastModified(obj.ModTime)
	output.SetContentLength(obj.Size)
	output.SetContentType(obj.ContentType)
	output.SetContentEncoding(obj.ContentEncoding)
	if !obj.Expires.IsZero() {
		output.SetExpiration(obj.Expires.UTC().Format(http.TimeFormat))
	}
	w.Header().Set(consts.Cid, obj.CID)
	output.SetMetadata(map[string]*string{
		consts.Cid: &obj.CID,
	})
	WriteSuccessResponse(w, output, "")
}

func WriteDeleteObjectResponse(w http.ResponseWriter, r *http.Request, obj *object.Object) {
	output := new(s3.DeleteObjectOutput)
	WriteSuccessResponse(w, output, "")
}

func WriteDeleteObjectsResponse(w http.ResponseWriter, r *http.Request, toErr func(error) *Error, deletes []*object.DeletedObject) {
	output := new(s3.DeleteObjectsOutput)
	objs := make([]*s3.DeletedObject, 0)
	errs := make([]*s3.Error, 0)
	for _, obj := range deletes {
		if obj.DeleteErr != nil {
			rerr := toErr(obj.DeleteErr)
			s3Err := new(s3.Error)
			s3Err.SetCode(rerr.Code())
			s3Err.SetMessage(rerr.Description())
			s3Err.SetKey(obj.Object)
			errs = append(errs, s3Err)
			continue
		}
		s3Obj := new(s3.DeletedObject)
		s3Obj.SetKey(obj.Object)
		objs = append(objs, s3Obj)
	}
	if len(errs) > 0 {
		output.SetErrors(errs)
	}
	if len(objs) > 0 {
		output.SetDeleted(objs)
	}
	WriteSuccessResponse(w, output, "DeleteResult")
}

func WriteGetObjectResponse(w http.ResponseWriter, r *http.Request, obj *object.Object, body io.ReadCloser) {
	output := new(s3.GetObjectOutput)
	output.SetLastModified(obj.ModTime)
	output.SetContentLength(obj.Size)
	output.SetContentType(obj.ContentType)
	output.SetContentEncoding(obj.ContentEncoding)
	output.SetBody(body)
	if !obj.Expires.IsZero() {
		output.SetExpiration(obj.Expires.UTC().Format(http.TimeFormat))
	}
	w.Header().Set(consts.Cid, obj.CID)
	output.SetMetadata(map[string]*string{
		consts.Cid: &obj.CID,
	})
	WriteSuccessResponse(w, output, "")
}

func WriteListObjectsResponse(w http.ResponseWriter, r *http.Request, list *object.ObjectsList) {
	out := new(s3.ListObjectsOutput)
	out.SetName(list.Args.Bucket)
	out.SetEncodingType(list.Args.EncodingType)
	out.SetPrefix(utils.S3Encode(list.Args.Prefix, list.Args.EncodingType))
	out.SetMarker(utils.S3Encode(list.Args.Marker, list.Args.EncodingType))
	out.SetDelimiter(utils.S3Encode(list.Args.Delimiter, list.Args.EncodingType))
	out.SetMaxKeys(list.Args.MaxKeys)
	out.SetNextMarker(list.NextMarker)
	out.SetIsTruncated(list.IsTruncated)
	s3Objs := make([]*s3.Object, len(list.Objects))
	for i, obj := range list.Objects {
		s3Obj := new(s3.Object)
		s3Obj.SetETag(`"` + obj.ETag + `"`)
		s3Obj.SetOwner(newS3Owner(list.Owner))
		s3Obj.SetLastModified(obj.ModTime)
		s3Obj.SetKey(utils.S3Encode(obj.Name, list.Args.EncodingType))
		s3Obj.SetSize(obj.Size)
		s3Obj.SetStorageClass("")
		s3Objs[i] = s3Obj
		w.Header().Add(consts.Cid, obj.CID)
	}
	out.SetContents(s3Objs)
	s3CommPrefixes := make([]*s3.CommonPrefix, len(list.Prefixes))
	for i, cpf := range list.Prefixes {
		pfx := new(s3.CommonPrefix)
		pfx.SetPrefix(utils.S3Encode(cpf, list.Args.EncodingType))
		s3CommPrefixes[i] = pfx
	}
	out.SetCommonPrefixes(s3CommPrefixes)
	WriteSuccessResponse(w, out, "ListBucketResult")
}

func WriteListObjectsV2Response(w http.ResponseWriter, r *http.Request, list *object.ObjectsListV2) {
	out := new(s3.ListObjectsV2Output)
	out.SetName(list.Args.Bucket)
	out.SetEncodingType(list.Args.EncodingType)
	out.SetStartAfter(utils.S3Encode(list.Args.After, list.Args.EncodingType))
	out.SetDelimiter(utils.S3Encode(list.Args.Delimiter, list.Args.EncodingType))
	out.SetPrefix(utils.S3Encode(list.Args.Prefix, list.Args.EncodingType))
	out.SetMaxKeys(list.Args.MaxKeys)
	out.SetContinuationToken(base64.StdEncoding.EncodeToString([]byte(list.Args.Token)))
	out.SetNextContinuationToken(base64.StdEncoding.EncodeToString([]byte(list.NextContinuationToken)))
	out.SetIsTruncated(list.IsTruncated)
	s3Objs := make([]*s3.Object, len(list.Objects))
	for i, obj := range list.Objects {
		s3Obj := new(s3.Object)
		s3Obj.SetETag(`"` + obj.ETag + `"`)
		if list.Args.FetchOwner {
			s3Obj.SetOwner(newS3Owner(list.Owner))
		}
		s3Obj.SetLastModified(obj.ModTime)
		s3Obj.SetKey(utils.S3Encode(obj.Name, list.Args.EncodingType))
		s3Obj.SetSize(obj.Size)
		s3Obj.SetStorageClass("")
		s3Objs[i] = s3Obj
		w.Header().Add(consts.Cid, obj.CID)
	}
	out.SetContents(s3Objs)
	s3CommPrefixes := make([]*s3.CommonPrefix, len(list.Prefixes))
	for i, cpf := range list.Prefixes {
		pfx := new(s3.CommonPrefix)
		pfx.SetPrefix(utils.S3Encode(cpf, list.Args.EncodingType))
		s3CommPrefixes[i] = pfx
	}
	out.SetCommonPrefixes(s3CommPrefixes)
	out.SetKeyCount(int64(len(list.Objects) + len(list.Prefixes)))
	WriteSuccessResponse(w, out, "ListBucketResult")
}

func WriteGetObjectACLResponse(w http.ResponseWriter, r *http.Request, acl *object.ACL) {
	output := new(s3.GetObjectAclOutput)
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
		panic(fmt.Sprintf("unknwo acl <%s>", acl.ACL))
	}
	output.SetGrants(grants)
	WriteSuccessResponse(w, output, "AccessControlPolicy")
	return
}
