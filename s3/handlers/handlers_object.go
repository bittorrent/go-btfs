package handlers

import (
	"encoding/base64"
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/requests"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

func (h *Handlers) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	if _, ok := r.Header[consts.AmzCopySource]; ok {
		err = errors.New("shouldn't be copy")
		responses.WriteErrorResponse(w, r, responses.ErrInvalidCopySource)
		return
	}

	bucname, objname, rerr := requests.ParseBucketAndObject(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	_, rerr = requests.ParseObjectACL(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	err = s3utils.CheckPutObjectArgs(ctx, bucname, objname)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
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

	body, ok := r.Body.(*hash.Reader)
	if !ok {
		responses.WriteErrorResponse(w, r, responses.ErrInternalError)
		return
	}

	obj, err := h.objsvc.PutObject(ctx, ack, bucname, objname, body, r.ContentLength, meta)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WritePutObjectResponse(w, r, obj)

	return
}

func (h *Handlers) HeadObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	bucname, objname, rerr := requests.ParseBucketAndObject(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	err = s3utils.CheckGetObjArgs(ctx, bucname, objname)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	//objsvc
	obj, _, err := h.objsvc.GetObject(ctx, ack, bucname, objname, false)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteHeadObjectResponse(w, r, obj)
}

func (h *Handlers) CopyObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	dstBucket, dstObject, rerr := requests.ParseBucketAndObject(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	err = s3utils.CheckPutObjectArgs(ctx, dstBucket, dstObject)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	// Copy source path.
	cpSrcPath, err := url.QueryUnescape(r.Header.Get(consts.AmzCopySource))
	if err != nil {
		// Save unescaped string as is.
		cpSrcPath = r.Header.Get(consts.AmzCopySource)
		err = nil
	}

	srcBucket, srcObject := pathToBucketAndObject(cpSrcPath)
	// If source object is empty or bucket is empty, reply back invalid copy source.
	if srcObject == "" || srcBucket == "" {
		err = responses.ErrInvalidCopySource
		responses.WriteErrorResponse(w, r, responses.ErrInvalidCopySource)
		return
	}
	if err = s3utils.CheckGetObjArgs(ctx, srcBucket, srcObject); err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}
	if srcBucket == dstBucket && srcObject == dstObject {
		err = responses.ErrInvalidCopyDest
		responses.WriteErrorResponse(w, r, responses.ErrInvalidCopyDest)
		return
	}

	metadata := make(map[string]string)
	if isReplace(r) {
		var inputMeta map[string]string
		inputMeta, err = extractMetadata(ctx, r)
		if err != nil {
			rerr = h.respErr(err)
			responses.WriteErrorResponse(w, r, rerr)
			return
		}
		for key, val := range inputMeta {
			metadata[key] = val
		}
	}

	//objsvc
	obj, err := h.objsvc.CopyObject(ctx, ack, srcBucket, srcObject, dstBucket, dstObject, metadata)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteCopyObjectResponse(w, r, obj)
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

	bucname, objname, rerr := requests.ParseBucketAndObject(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	err = s3utils.CheckDelObjArgs(ctx, bucname, objname)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	err = h.objsvc.DeleteObject(ctx, ack, bucname, objname)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteDeleteObjectResponse(w, r, nil)
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

	var input s3.DeleteObjectsInput

	err = responses.ParseRequest(r, &input)
	if err != nil {
		rerr := h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	if input.Delete == nil ||
		len(input.Delete.Objects) == 0 ||
		len(input.Delete.Objects) > consts.MaxObjectList {
		rerr := responses.ErrMalformedXML
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	bucname := *input.Bucket

	_, err = h.objsvc.GetBucket(ctx, ack, bucname)
	if err != nil {
		rerr := h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	output := new(s3.DeleteObjectsOutput)
	delObjs := make([]*s3.DeletedObject, 0)
	delErrs := make([]*s3.Error, 0)
	for _, obj := range input.Delete.Objects {
		objname := *obj.Key
		er := s3utils.CheckDelObjArgs(ctx, bucname, objname)
		if er != nil {
			rerr := h.respErr(er)
			derr := new(s3.Error)
			derr.SetCode(rerr.Code())
			derr.SetMessage(rerr.Description())
			derr.SetKey(objname)
			delErrs = append(delErrs, derr)
			continue
		}
		er = h.objsvc.DeleteObject(ctx, ack, bucname, objname)
		if er != nil {
			rerr := h.respErr(er)
			derr := new(s3.Error)
			derr.SetCode(rerr.Code())
			derr.SetMessage(rerr.Description())
			derr.SetKey(objname)
			delErrs = append(delErrs, derr)
		} else {
			dobj := new(s3.DeletedObject)
			dobj.SetKey(objname)
			delObjs = append(delObjs, dobj)
		}
	}

	output.SetDeleted(delObjs)
	output.SetErrors(delErrs)

	responses.WriteSuccessResponse(w, output, "DeleteResult")
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

	bucname, objname, rerr := requests.ParseBucketAndObject(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	if err = s3utils.CheckGetObjArgs(ctx, bucname, objname); err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	obj, body, err := h.objsvc.GetObject(ctx, ack, bucname, objname, true)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteGetObjectResponse(w, r, obj, body)
}

// GetObjectACLHandler - GET Object ACL
func (h *Handlers) GetObjectACLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	bucname, _, rerr := requests.ParseBucketAndObject(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	acl, err := h.objsvc.GetBucketACL(ctx, ack, bucname)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteGetObjectACLResponse(w, r, ack, acl)
}

func (h *Handlers) ListObjectsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	bucname, rerr := requests.ParseBucket(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	// Extract all the listsObjectsV1 query params to their native values.
	prefix, marker, delimiter, maxKeys, encodingType, rerr := getListObjectsV1Args(r.Form)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	err = s3utils.CheckListObjsArgs(ctx, bucname, prefix, marker)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}
	list, err := h.objsvc.ListObjects(ctx, ack, bucname, prefix, delimiter, marker, maxKeys)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteListObjectsResponse(w, r, ack, bucname, prefix, marker, delimiter, encodingType, maxKeys, list)
}

func (h *Handlers) ListObjectsV2Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	bucname, rerr := requests.ParseBucket(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	urlValues := r.Form
	// Extract all the listObjectsV2 query params to their native values.
	prefix, token, startAfter, delimiter, fetchOwner, maxKeys, encodingType, rerr := getListObjectsV2Args(urlValues)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	marker := token
	if marker == "" {
		marker = startAfter
	}
	err = s3utils.CheckListObjsArgs(ctx, bucname, prefix, marker)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	// Validate the query params before beginning to serve the request.
	// fetch-owner is not validated since it is a boolean
	rerr = validateListObjectsArgs(token, delimiter, encodingType, maxKeys)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	list, err := h.objsvc.ListObjectsV2(ctx, ack, bucname, prefix, token, delimiter,
		maxKeys, fetchOwner, startAfter)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteListObjectsV2Response(w, r, ack, bucname, prefix, token, startAfter,
		delimiter, encodingType, maxKeys, list)
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
	prefix, marker, delimiter string, maxkeys int64, encodingType string, rerr *responses.Error) {

	if values.Get("max-keys") != "" {
		var err error
		if maxkeys, err = strconv.ParseInt(values.Get("max-keys"), 10, 64); err != nil {
			rerr = responses.ErrInvalidMaxKeys
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
	fetchOwner bool, maxkeys int64, encodingType string, rerr *responses.Error) {

	// The continuation-token cannot be empty.
	if val, ok := values["continuation-token"]; ok {
		if len(val[0]) == 0 {
			rerr = responses.ErrInvalidToken
			return
		}
	}

	if values.Get("max-keys") != "" {
		var err error
		if maxkeys, err = strconv.ParseInt(values.Get("max-keys"), 10, 64); err != nil {
			rerr = responses.ErrInvalidMaxKeys
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
			rerr = responses.ErrIncorrectContinuationToken
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
func validateListObjectsArgs(marker, delimiter, encodingType string, maxKeys int64) (rerr *responses.Error) {
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
