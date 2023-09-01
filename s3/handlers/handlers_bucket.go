package handlers

import (
	"errors"
	s3action "github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/requests"
	"github.com/bittorrent/go-btfs/s3/responses"
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
	if !ok {
		err = responses.ErrInternalError
	}
	return
}

func (h *Handlers) PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	ctx := r.Context()

	req, rerr := requests.ParsePutBucketRequest(r)
	if rerr != nil {
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	_, err = h.objsvc.PutBucket(ctx, req.User, req.Bucket, req.Region, req.ACL)
	if err != nil {
		rerr = h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	responses.WritePutBucketResponse(w, r)

	return
}

func (h *Handlers) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	req, err := requests.ParseDeleteBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestBody)
		return
	}

	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	err = h.bucsvc.CheckACL(ack, req.Bucket, s3action.HeadBucketAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	//todo check all errors.
	err = h.bucsvc.DeleteBucket(ctx, req.Bucket)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	responses.WriteDeleteBucketResponse(w)
}

func (h *Handlers) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	ack := cctx.GetAccessKey(r)
	if ack == "" {
		responses.WriteErrorResponse(w, r, responses.ErrNoAccessKey)
		return
	}

	//todo check all errors
	bucketMetas, err := h.bucsvc.GetAllBucketsOfUser(ack)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	responses.WriteListBucketsResponse(w, r, bucketMetas)
}

func (h *Handlers) GetBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	req, err := requests.ParseGetBucketAclRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestBody)
		return
	}

	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	if !h.bucsvc.HasBucket(ctx, req.Bucket) {
		responses.WriteErrorResponseHeadersOnly(w, r, responses.ErrNoSuchBucket)
		return
	}

	err = h.bucsvc.CheckACL(ack, req.Bucket, s3action.GetBucketAclAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	//todo check all errors
	acl, err := h.bucsvc.GetBucketAcl(ctx, req.Bucket)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	responses.WriteGetBucketAclResponse(w, r, ack, acl)
}

func (h *Handlers) PutBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	req, err := requests.ParsePutBucketAclRequest(r)
	if err != nil || len(req.ACL) == 0 || len(req.Bucket) == 0 {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestBody)
		return
	}

	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	err = h.bucsvc.CheckACL(ack, req.Bucket, s3action.PutBucketAclAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	if !requests.CheckAclPermissionType(&req.ACL) {
		responses.WriteErrorResponse(w, r, responses.ErrNotImplemented)
		return
	}

	//todo check all errors
	err = h.bucsvc.UpdateBucketAcl(ctx, req.Bucket, req.ACL)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	//todo check no return?
	responses.WritePutBucketAclResponse(w, r)
}

func (h *Handlers) HeadBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	req, err := requests.ParseHeadBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestBody)
		return
	}

	ack := cctx.GetAccessKey(r)

	err = h.bucsvc.CheckACL(ack, req.Bucket, s3action.HeadBucketAction)
	if errors.Is(err, object.ErrBucketNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrAccessDenied)
		return
	}

	responses.WriteHeadBucketResponse(w, r)
}
