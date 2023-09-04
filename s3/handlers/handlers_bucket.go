package handlers

import (
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/requests"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"net/http"
)

var errToRespErr = map[error]*responses.Error{
	object.ErrBucketNotFound:      responses.ErrNoSuchBucket,
	object.ErrObjectNotFound:      responses.ErrNoSuchKey,
	object.ErrUploadNotFound:      responses.ErrNoSuchUpload,
	object.ErrBucketAlreadyExists: responses.ErrBucketAlreadyExists,
	object.ErrNotAllowed:          responses.ErrAccessDenied,
}

func (h *Handlers) respErr(err error) (rerr *responses.Error) {
	rerr, ok := errToRespErr[err]
	if ok {
		return
	}
	switch err.(type) {
	case s3utils.BucketNameInvalid:
		rerr = responses.ErrInvalidBucketName
	case s3utils.ObjectNameInvalid:
		rerr = responses.ErrInvalidObjectName
	case s3utils.InvalidPart:
		rerr = responses.ErrInvalidPart
	case s3utils.InvalidUploadID:
		rerr = responses.ErrNoSuchUpload
	case s3utils.InvalidMarkerPrefixCombination:
		rerr = responses.ErrInvalidRequestParameter
	default:
		rerr = responses.ErrInternalError
	}
	return
}

func (h *Handlers) CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	req, rerr := requests.ParseCreateBucketRequest(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	_, err = h.objsvc.CreateBucket(r.Context(), req.AccessKey, req.Bucket, req.Region, req.ACL)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteCreateBucketResponse(w, r)

	return
}

func (h *Handlers) HeadBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	req, rerr := requests.ParseHeadBucketRequest(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	_, err = h.objsvc.GetBucket(r.Context(), req.AccessKey, req.Bucket)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteHeadBucketResponse(w, r)

	return
}

func (h *Handlers) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	req, rerr := requests.ParseDeleteBucketRequest(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	err = h.objsvc.DeleteBucket(r.Context(), req.AccessKey, req.Bucket)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteDeleteBucketResponse(w)

	return
}

func (h *Handlers) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	req, rerr := requests.ParseListBucketsRequest(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	list, err := h.objsvc.GetAllBuckets(r.Context(), req.AccessKey)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteListBucketsResponse(w, r, req.AccessKey, list)

	return
}

func (h *Handlers) GetBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	req, rerr := requests.ParseGetBucketACLRequest(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	acl, err := h.objsvc.GetBucketACL(r.Context(), req.AccessKey, req.Bucket)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WriteGetBucketACLResponse(w, r, req.AccessKey, acl)
}

func (h *Handlers) PutBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	req, rerr := requests.ParsePutBucketAclRequest(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	err = h.objsvc.PutBucketACL(r.Context(), req.AccessKey, req.Bucket, req.ACL)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WritePutBucketAclResponse(w, r)
}
