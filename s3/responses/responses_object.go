package responses

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"io"
	"net/http"
)

func WritePutObjectResponse(w http.ResponseWriter, r *http.Request, obj *object.Object) {
	output := new(s3.PutObjectOutput)
	output.SetETag(`"` + obj.ETag + `"`)
	w.Header().Set(consts.CID, obj.CID)
	WriteSuccessResponse(w, output, "")
}

func WriteHeadObjectResponse(w http.ResponseWriter, r *http.Request, obj *object.Object) {
	output := new(s3.HeadObjectOutput)
	w.Header().Set(consts.CID, obj.CID)
	SetObjectHeaders(w, r, obj)
	SetHeadGetRespHeaders(w, r.Form)
	WriteSuccessResponse(w, output, "")
}

func WriteCopyObjectResponse(w http.ResponseWriter, r *http.Request, obj *object.Object) {
	output := new(s3.CopyObjectResult)
	output.SetETag(`"` + obj.ETag + `"`)
	output.SetLastModified(obj.ModTime)
	w.Header().Set(consts.CID, obj.CID)
	WriteSuccessResponse(w, output, "")
}

func WriteDeleteObjectResponse(w http.ResponseWriter, r *http.Request, obj *object.Object) {
	output := new(s3.DeleteObjectOutput)
	output.SetDeleteMarker(true)
	WriteSuccessResponse(w, output, "")
}

func WriteGetObjectResponse(w http.ResponseWriter, r *http.Request, obj *object.Object, body io.ReadCloser) {
	output := new(s3.GetObjectOutput)
	output.SetContentLength(obj.Size)
	output.SetBody(body)
	output.SetMetadata(map[string]*string{
		consts.CID: &obj.CID,
	})
	w.Header().Set(consts.CID, obj.CID)
	SetObjectHeaders(w, r, obj)
	SetHeadGetRespHeaders(w, r.Form)
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

func WriteListObjectsResponse(w http.ResponseWriter, r *http.Request, accessKey, bucname, prefix, marker, delimiter, encodingType string, maxKeys int64, list *object.ObjectsList) {
	out := new(s3.ListObjectsOutput)
	out.SetName(bucname)
	out.SetPrefix(prefix)
	out.SetMarker(marker)
	out.SetDelimiter(delimiter)
	out.SetEncodingType(encodingType)
	out.SetMaxKeys(maxKeys)
	s3Objs := make([]*s3.Object, len(list.Objects))
	for i, obj := range list.Objects {
		s3Obj := new(s3.Object)
		s3Obj.SetETag(`"` + obj.ETag + `"`)
		s3Obj.SetOwner(owner(accessKey))
		s3Obj.SetLastModified(obj.ModTime)
		s3Obj.SetKey(obj.Name)
		s3Obj.SetSize(obj.Size)
		s3Objs[i] = s3Obj
	}
	out.SetContents(s3Objs)
	s3CommPrefixes := make([]*s3.CommonPrefix, len(list.Prefixes))
	for i, cpf := range list.Prefixes {
		pfx := new(s3.CommonPrefix)
		pfx.SetPrefix(cpf)
		s3CommPrefixes[i] = pfx
	}
	out.SetCommonPrefixes(s3CommPrefixes)
	out.SetIsTruncated(list.IsTruncated)
	out.SetNextMarker(list.NextMarker)
	WriteSuccessResponse(w, out, "")
}
