package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/requests"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/bittorrent/go-btfs/s3/services/bucket"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
)

const lockWaitTimeout = 5 * time.Minute

func (h *Handlers) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	// X-Amz-Copy-Source shouldn't be set for this call.
	if _, ok := r.Header[consts.AmzCopySource]; ok {
		err = errors.New("shouldn't be copy")
		responses.WriteErrorResponse(w, r, responses.ErrInvalidCopySource)
		return
	}

	aclHeader := r.Header.Get(consts.AmzACL)
	if aclHeader != "" {
		err = errors.New("object acl can only set to default")
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}

	bucname, objname, err := requests.ParseBucketAndObject(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}

	err = s3utils.CheckPutObjectArgs(ctx, bucname, objname)
	if err != nil { // todo: convert error
		responses.WriteErrorResponse(w, r, err)
		return
	}

	meta, err := extractMetadata(ctx, r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequest)
		return
	}

	if r.ContentLength == 0 {
		responses.WriteErrorResponse(w, r, responses.ErrEntityTooSmall)
		return
	}

	hrdr, ok := r.Body.(*hash.Reader)
	if !ok {
		responses.WriteErrorResponse(w, r, responses.ErrInternalError)
		return
	}

	// rlock bucket
	runlock, err := h.rlock(ctx, bucname, w, r)
	if err != nil {
		return
	}
	defer runlock()

	// lock object
	unlock, err := h.lock(ctx, bucname+"/"+objname, w, r)
	if err != nil {
		return
	}
	defer unlock()

	err = h.bucsvc.CheckACL(ack, bucname, action.PutObjectAction)
	if errors.Is(err, bucket.ErrNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	//objsvc
	obj, err := h.objsvc.PutObject(ctx, bucname, objname, hrdr, r.ContentLength, meta)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	responses.WritePutObjectResponse(w, r, obj)

	return
}

// HeadObjectHandler - HEAD Object
func (h *Handlers) HeadObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	bucname, objname, err := requests.ParseBucketAndObject(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}

	if err := s3utils.CheckGetObjArgs(ctx, bucname, objname); err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}

	err = h.bucsvc.CheckACL(ack, bucname, action.HeadObjectAction)
	if errors.Is(err, bucket.ErrNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	// rlock bucket
	runlock, err := h.rlock(ctx, bucname, w, r)
	if err != nil {
		return
	}
	defer runlock()

	// rlock object
	runlockObj, err := h.rlock(ctx, bucname+"/"+objname, w, r)
	if err != nil {
		return
	}
	defer runlockObj()

	//objsvc
	obj, err := h.objsvc.GetObjectInfo(ctx, bucname, objname)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}
	w.Header().Set(consts.AmzServerSideEncryption, consts.AmzEncryptionAES)

	// Set standard object headers.
	responses.SetObjectHeaders(w, r, obj)
	// Set any additional requested response headers.
	responses.SetHeadGetRespHeaders(w, r.Form)

	// Successful response.
	w.WriteHeader(http.StatusOK)
}

// CopyObjectHandler - Copy Object
func (h *Handlers) CopyObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	dstBucket, dstObject, err := requests.ParseBucketAndObject(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}
	if err := s3utils.CheckPutObjectArgs(ctx, dstBucket, dstObject); err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}
	err = h.bucsvc.CheckACL(ack, dstBucket, action.CopyObjectAction)
	if errors.Is(err, bucket.ErrNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	// Copy source path.
	cpSrcPath, err := url.QueryUnescape(r.Header.Get(consts.AmzCopySource))
	if err != nil {
		// Save unescaped string as is.
		cpSrcPath = r.Header.Get(consts.AmzCopySource)
	}
	srcBucket, srcObject := pathToBucketAndObject(cpSrcPath)
	// If source object is empty or bucket is empty, reply back invalid copy source.
	if srcObject == "" || srcBucket == "" {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidCopySource)
		return
	}
	if err = s3utils.CheckGetObjArgs(ctx, srcBucket, srcObject); err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}
	if srcBucket == dstBucket && srcObject == dstObject {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidCopyDest)
		return
	}
	err = h.bucsvc.CheckACL(ack, srcBucket, action.CopyObjectAction)
	if errors.Is(err, bucket.ErrNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	fmt.Printf("CopyObjectHandler %s %s => %s %s \n", srcBucket, srcObject, dstBucket, dstObject)

	// rlock bucket 1
	runlock1, err := h.rlock(ctx, srcBucket, w, r)
	if err != nil {
		return
	}
	defer runlock1()

	// rlock object 1
	runlockObj1, err := h.rlock(ctx, srcBucket+"/"+srcObject, w, r)
	if err != nil {
		return
	}
	defer runlockObj1()

	// rlock bucket 2
	runlock2, err := h.rlock(ctx, dstBucket, w, r)
	if err != nil {
		return
	}
	defer runlock2()

	// lock object 2
	unlockObj2, err := h.lock(ctx, dstBucket+"/"+dstObject, w, r)
	if err != nil {
		return
	}
	defer unlockObj2()

	//objsvc
	srcObjInfo, err := h.objsvc.GetObjectInfo(ctx, srcBucket, srcObject)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	metadata := make(map[string]string)
	metadata[strings.ToLower(consts.ContentType)] = srcObjInfo.ContentType
	metadata[strings.ToLower(consts.ContentEncoding)] = srcObjInfo.ContentEncoding
	if isReplace(r) {
		inputMeta, err := extractMetadata(ctx, r)
		if err != nil {
			responses.WriteErrorResponse(w, r, err)
			return
		}
		for key, val := range inputMeta {
			metadata[key] = val
		}
	}

	//objsvc
	obj, err := h.objsvc.CopyObject(ctx, dstBucket, dstObject, srcObjInfo, srcObjInfo.Size, metadata)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	resp := responses.CopyObjectResult{
		ETag:         "\"" + obj.ETag + "\"",
		LastModified: obj.ModTime.UTC().Format(consts.Iso8601TimeFormat),
	}

	setPutObjHeaders(w, obj, false)

	responses.WriteSuccessResponseXML(w, r, resp)
}

// DeleteObjectHandler - delete an object
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html
func (h *Handlers) DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	bucname, objname, err := requests.ParseBucketAndObject(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}
	if err := s3utils.CheckDelObjArgs(ctx, bucname, objname); err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	err = h.bucsvc.CheckACL(ack, bucname, action.DeleteObjectAction)
	if errors.Is(err, bucket.ErrNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	// rlock bucket
	runlock, err := h.rlock(ctx, bucname, w, r)
	if err != nil {
		return
	}
	defer runlock()

	// lock object
	unlock, err := h.lock(ctx, bucname+"/"+objname, w, r)
	if err != nil {
		return
	}
	defer unlock()

	//objsvc
	obj, err := h.objsvc.GetObjectInfo(ctx, bucname, objname)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}
	//objsvc
	err = h.objsvc.DeleteObject(ctx, bucname, objname)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}
	setPutObjHeaders(w, obj, true)
	responses.WriteSuccessNoContent(w)
}

// DeleteObjectsHandler - delete objects
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObjects.html
func (h *Handlers) DeleteObjectsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	bucname, objname, err := requests.ParseBucketAndObject(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}
	if err := s3utils.CheckDelObjArgs(ctx, bucname, objname); err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	err = h.bucsvc.CheckACL(ack, bucname, action.DeleteObjectAction)
	if errors.Is(err, bucket.ErrNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	// rlock bucket
	runlock, err := h.rlock(ctx, bucname, w, r)
	if err != nil {
		return
	}
	defer runlock()

	// lock object
	unlock, err := h.lock(ctx, bucname+"/"+objname, w, r)
	if err != nil {
		return
	}
	defer unlock()

	//objsvc
	obj, err := h.objsvc.GetObjectInfo(ctx, bucname, objname)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}
	//objsvc
	err = h.objsvc.DeleteObject(ctx, bucname, objname)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}
	setPutObjHeaders(w, obj, true)
	responses.WriteSuccessNoContent(w)
}

// GetObjectHandler - GET Object
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html
func (h *Handlers) GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	bucname, objname, err := requests.ParseBucketAndObject(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}
	if err = s3utils.CheckGetObjArgs(ctx, bucname, objname); err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	err = h.bucsvc.CheckACL(ack, bucname, action.GetObjectAction)
	if errors.Is(err, bucket.ErrNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	// rlock bucket
	runlock, err := h.rlock(ctx, bucname, w, r)
	if err != nil {
		return
	}
	defer runlock()

	// rlock object
	runlockObj, err := h.rlock(ctx, bucname+"/"+objname, w, r)
	if err != nil {
		return
	}
	defer runlockObj()

	//objsvc
	obj, reader, err := h.objsvc.GetObject(ctx, bucname, objname)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}
	//w.Header().Set(consts.AmzServerSideEncryption, consts.AmzEncryptionAES)

	responses.SetObjectHeaders(w, r, obj)
	w.Header().Set(consts.ContentLength, strconv.FormatInt(obj.Size, 10))
	responses.SetHeadGetRespHeaders(w, r.Form)
	_, err = io.Copy(w, reader)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInternalError)
		return
	}
}

// GetObjectACLHandler - GET Object ACL
func (h *Handlers) GetObjectACLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	bucname, _, err := requests.ParseBucketAndObject(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}

	err = h.bucsvc.CheckACL(ack, bucname, action.GetBucketAclAction)
	if errors.Is(err, bucket.ErrNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	// rlock bucket
	runlock, err := h.rlock(ctx, bucname, w, r)
	if err != nil {
		return
	}
	defer runlock()

	acl, err := h.bucsvc.GetBucketAcl(ctx, bucname)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	responses.WriteGetBucketAclResponse(w, r, ack, acl)
}

func (h *Handlers) ListObjectsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	bucname, _, err := requests.ParseBucketAndObject(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}

	// Extract all the litsObjectsV1 query params to their native values.
	prefix, marker, delimiter, maxKeys, encodingType, s3Error := getListObjectsV1Args(r.Form)
	if s3Error != nil {
		responses.WriteErrorResponse(w, r, s3Error)
		return
	}

	if err := s3utils.CheckListObjsArgs(ctx, bucname, prefix, marker); err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	// rlock bucket
	runlock, err := h.rlock(ctx, bucname, w, r)
	if err != nil {
		return
	}
	defer runlock()

	err = h.bucsvc.CheckACL(ack, bucname, action.ListObjectsAction)
	if errors.Is(err, bucket.ErrNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	//objsvc
	objs, err := h.objsvc.ListObjects(ctx, bucname, prefix, marker, delimiter, maxKeys)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}
	resp := responses.GenerateListObjectsV1Response(bucname, prefix, marker, delimiter, encodingType, maxKeys, objs)
	// Write success response.
	responses.WriteSuccessResponseXML(w, r, resp)
}

func (h *Handlers) ListObjectsV2Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	bucname, _, err := requests.ParseBucketAndObject(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}

	err = h.bucsvc.CheckACL(ack, bucname, action.ListObjectsAction)
	if errors.Is(err, bucket.ErrNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	urlValues := r.Form
	// Extract all the listObjectsV2 query params to their native values.
	prefix, token, startAfter, delimiter, fetchOwner, maxKeys, encodingType, errCode := getListObjectsV2Args(urlValues)
	if errCode != nil {
		responses.WriteErrorResponse(w, r, errCode)
		return
	}

	marker := token
	if marker == "" {
		marker = startAfter
	}
	if err := s3utils.CheckListObjsArgs(ctx, bucname, prefix, marker); err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	// Validate the query params before beginning to serve the request.
	// fetch-owner is not validated since it is a boolean
	s3Error := validateListObjectsArgs(token, delimiter, encodingType, maxKeys)
	if s3Error != nil {
		responses.WriteErrorResponse(w, r, s3Error)
		return
	}

	// rlock bucket
	runlock, err := h.rlock(ctx, bucname, w, r)
	if err != nil {
		return
	}
	defer runlock()

	// Initiate a list objects operation based on the input params.
	// On success would return back ListObjectsInfo object to be
	// marshaled into S3 compatible XML header.
	//objsvc
	listObjectsV2Info, err := h.objsvc.ListObjectsV2(ctx, bucname, prefix, token, delimiter,
		maxKeys, fetchOwner, startAfter)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	resp := responses.GenerateListObjectsV2Response(
		bucname, prefix, token, listObjectsV2Info.NextContinuationToken, startAfter,
		delimiter, encodingType, listObjectsV2Info.IsTruncated,
		maxKeys, listObjectsV2Info.Objects, listObjectsV2Info.Prefixes)

	// Write success response.
	responses.WriteSuccessResponseXML(w, r, resp)
}

// setPutObjHeaders sets all the necessary headers returned back
// upon a success Put/Copy/CompleteMultipart/Delete requests
// to activate delete only headers set delete as true
func setPutObjHeaders(w http.ResponseWriter, obj object.Object, delete bool) {
	// We must not use the http.Header().Set method here because some (broken)
	// clients expect the ETag header key to be literally "ETag" - not "Etag" (case-sensitive).
	// Therefore, we have to set the ETag directly as map entry.
	if obj.ETag != "" && !delete {
		w.Header()[consts.ETag] = []string{`"` + obj.ETag + `"`}
	}

	// Set the relevant version ID as part of the response header.
	if obj.VersionID != "" {
		w.Header()[consts.AmzVersionID] = []string{obj.VersionID}
		// If version is a deleted marker, set this header as well
		if obj.DeleteMarker && delete { // only returned during delete object
			w.Header()[consts.AmzDeleteMarker] = []string{strconv.FormatBool(obj.DeleteMarker)}
		}
	}

	if obj.Bucket != "" && obj.Name != "" {
		// do something
	}
}

func pathToBucketAndObject(path string) (bucket, object string) {
	path = strings.TrimPrefix(path, consts.SlashSeparator)
	idx := strings.Index(path, consts.SlashSeparator)
	if idx < 0 {
		return path, ""
	}
	return path[:idx], path[idx+len(consts.SlashSeparator):]
}

func isReplace(r *http.Request) bool {
	return r.Header.Get("X-Amz-Metadata-Directive") == "REPLACE"
}

// Parse bucket url queries
func getListObjectsV1Args(values url.Values) (
	prefix, marker, delimiter string, maxkeys int, encodingType string, errCode error) {

	if values.Get("max-keys") != "" {
		var err error
		if maxkeys, err = strconv.Atoi(values.Get("max-keys")); err != nil {
			errCode = responses.ErrInvalidMaxKeys
			return
		}
	} else {
		maxkeys = consts.MaxObjectList
	}

	prefix = trimLeadingSlash(values.Get("prefix"))
	marker = trimLeadingSlash(values.Get("marker"))
	delimiter = values.Get("delimiter")
	encodingType = values.Get("encoding-type")
	return
}

// Parse bucket url queries for ListObjects V2.
func getListObjectsV2Args(values url.Values) (
	prefix, token, startAfter, delimiter string,
	fetchOwner bool, maxkeys int, encodingType string, errCode error) {

	// The continuation-token cannot be empty.
	if val, ok := values["continuation-token"]; ok {
		if len(val[0]) == 0 {
			errCode = responses.ErrInvalidToken
			return
		}
	}

	if values.Get("max-keys") != "" {
		var err error
		if maxkeys, err = strconv.Atoi(values.Get("max-keys")); err != nil {
			errCode = responses.ErrInvalidMaxKeys
			return
		}
		// Over flowing count - reset to maxObjectList.
		if maxkeys > consts.MaxObjectList {
			maxkeys = consts.MaxObjectList
		}
	} else {
		maxkeys = consts.MaxObjectList
	}

	prefix = trimLeadingSlash(values.Get("prefix"))
	startAfter = trimLeadingSlash(values.Get("start-after"))
	delimiter = values.Get("delimiter")
	fetchOwner = values.Get("fetch-owner") == "true"
	encodingType = values.Get("encoding-type")

	if token = values.Get("continuation-token"); token != "" {
		decodedToken, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			errCode = responses.ErrIncorrectContinuationToken
			return
		}
		token = string(decodedToken)
	}
	return
}

func trimLeadingSlash(ep string) string {
	if len(ep) > 0 && ep[0] == '/' {
		// Path ends with '/' preserve it
		if ep[len(ep)-1] == '/' && len(ep) > 1 {
			ep = path.Clean(ep)
			ep += "/"
		} else {
			ep = path.Clean(ep)
		}
		ep = ep[1:]
	}
	return ep
}

// Validate all the ListObjects query arguments, returns an APIErrorCode
// if one of the args do not meet the required conditions.
//   - delimiter if set should be equal to '/', otherwise the request is rejected.
//   - marker if set should have a common prefix with 'prefix' param, otherwise
//     the request is rejected.
func validateListObjectsArgs(marker, delimiter, encodingType string, maxKeys int) error {
	// Max keys cannot be negative.
	if maxKeys < 0 {
		return responses.ErrInvalidMaxKeys
	}

	if encodingType != "" {
		// AWS S3 spec only supports 'url' encoding type
		if !strings.EqualFold(encodingType, "url") {
			return responses.ErrInvalidEncodingMethod
		}
	}

	return nil
}
