package handlers

import (
	"encoding/base64"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/requests"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

// PutObjectHandler .
func (h *Handlers) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParsePutObjectRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	obj, err := h.objsvc.PutObject(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	responses.WritePutObjectResponse(w, r, obj)
	return
}

// CopyObjectHandler .
func (h *Handlers) CopyObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseCopyObjectRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	obj, err := h.objsvc.CopyObject(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	responses.WriteCopyObjectResponse(w, r, obj)
	return
}

// HeadObjectHandler .
func (h *Handlers) HeadObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseHeadObjectRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	obj, _, err := h.objsvc.GetObject(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	responses.WriteHeadObjectResponse(w, r, obj)
	return
}

// GetObjectHandler .
func (h *Handlers) GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseGetObjectRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	obj, body, err := h.objsvc.GetObject(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	responses.WriteGetObjectResponse(w, r, obj, body)
	return
}

// DeleteObjectHandler .
func (h *Handlers) DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseDeleteObjectRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}
	err = h.objsvc.DeleteObject(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	responses.WriteDeleteObjectResponse(w, r, nil)
	return
}

// DeleteObjectsHandler .
func (h *Handlers) DeleteObjectsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseDeleteObjectsRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	deletedObjects, err := h.objsvc.DeleteObjects(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	responses.WriteDeleteObjectsResponse(w, r, h.toRespErr, deletedObjects)
	return
}

// ListObjectsHandler .
func (h *Handlers) ListObjectsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseListObjectsRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	list, err := h.objsvc.ListObjects(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	responses.WriteListObjectsResponse(w, r, ack, list)
	return
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
		rerr = h.toRespErr(err)
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
		rerr = h.toRespErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteListObjectsV2Response(w, r, ack, bucname, prefix, token, startAfter,
		delimiter, encodingType, maxKeys, list)
}

// GetObjectACLHandler - GET Object ACL
func (h *Handlers) GetObjectACLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseGetBucketACLRequest()

	bucname, objname, rerr := requests.ParseBucketAndObject(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	acl, err := h.objsvc.GetObjectACL(ctx, ack, bucname, objname)
	if err != nil {
		rerr = h.toRespErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteGetObjectACLResponse(w, r, ack, acl)
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
