package responses

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/services/object"
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

func WriteDeleteObjectsResponse(w http.ResponseWriter, r *http.Request, toErr func(error) *Error, deletedObjects []*object.DeletedObject) {
	output := new(s3.DeleteObjectsOutput)
	objs := make([]*s3.DeletedObject, 0)
	errs := make([]*s3.Error, 0)
	for _, obj := range deletedObjects {
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

func WriteGetObjectACLResponse(w http.ResponseWriter, r *http.Request, accessKey, acl string) {
	output := new(s3.GetObjectAclOutput)
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
	WriteSuccessResponse(w, output, "AccessControlPolicy")
	return
}

func WriteListObjectsResponse(w http.ResponseWriter, r *http.Request, accessKey string, list *object.ObjectsList) {
	out := new(s3.ListObjectsOutput)
	out.SetName(list.Bucket)
	out.SetEncodingType(list.EncodingType)
	out.SetPrefix(utils.S3Encode(list.Prefix, list.EncodingType))
	out.SetMarker(utils.S3Encode(list.Marker, list.EncodingType))
	out.SetDelimiter(utils.S3Encode(list.Delimiter, list.EncodingType))
	out.SetMaxKeys(list.MaxKeys)
	out.SetNextMarker(list.NextMarker)
	out.SetIsTruncated(list.IsTruncated)
	s3Objs := make([]*s3.Object, len(list.Objects))
	for i, obj := range list.Objects {
		s3Obj := new(s3.Object)
		s3Obj.SetETag(`"` + obj.ETag + `"`)
		s3Obj.SetOwner(owner(accessKey))
		s3Obj.SetLastModified(obj.ModTime)
		s3Obj.SetKey(utils.S3Encode(obj.Name, list.EncodingType))
		s3Obj.SetSize(obj.Size)
		s3Obj.SetStorageClass("")
		s3Objs[i] = s3Obj
		w.Header().Add(consts.Cid, obj.CID)
	}
	out.SetContents(s3Objs)
	s3CommPrefixes := make([]*s3.CommonPrefix, len(list.Prefixes))
	for i, cpf := range list.Prefixes {
		pfx := new(s3.CommonPrefix)
		pfx.SetPrefix(utils.S3Encode(cpf, list.EncodingType))
		s3CommPrefixes[i] = pfx
	}
	out.SetCommonPrefixes(s3CommPrefixes)
	WriteSuccessResponse(w, out, "ListBucketResult")
}

func WriteListObjectsV2Response(w http.ResponseWriter, r *http.Request, accessKey, bucname, prefix, token, startAfter, delimiter, encodingType string, maxKeys int64, list *object.ObjectsListV2) {
	out := new(s3.ListObjectsV2Output)
	out.SetName(bucname)
	out.SetEncodingType(encodingType)
	out.SetStartAfter(utils.S3Encode(startAfter, encodingType))
	out.SetDelimiter(utils.S3Encode(delimiter, encodingType))
	out.SetPrefix(utils.S3Encode(prefix, encodingType))
	out.SetMaxKeys(maxKeys)
	out.SetContinuationToken(base64.StdEncoding.EncodeToString([]byte(token)))
	out.SetNextContinuationToken(base64.StdEncoding.EncodeToString([]byte(list.NextContinuationToken)))
	out.SetIsTruncated(list.IsTruncated)
	s3Objs := make([]*s3.Object, len(list.Objects))
	for i, obj := range list.Objects {
		s3Obj := new(s3.Object)
		s3Obj.SetETag(`"` + obj.ETag + `"`)
		s3Obj.SetOwner(owner(accessKey))
		s3Obj.SetLastModified(obj.ModTime)
		s3Obj.SetKey(utils.S3Encode(obj.Name, encodingType))
		s3Obj.SetSize(obj.Size)
		s3Obj.SetStorageClass("")
		s3Objs[i] = s3Obj
		w.Header().Add(consts.Cid, obj.CID)
	}
	out.SetContents(s3Objs)
	s3CommPrefixes := make([]*s3.CommonPrefix, len(list.Prefixes))
	for i, cpf := range list.Prefixes {
		pfx := new(s3.CommonPrefix)
		pfx.SetPrefix(utils.S3Encode(cpf, encodingType))
		s3CommPrefixes[i] = pfx
	}
	out.SetCommonPrefixes(s3CommPrefixes)
	out.SetKeyCount(int64(len(list.Objects) + len(list.Prefixes)))
	WriteSuccessResponse(w, out, "ListBucketResult")
}
