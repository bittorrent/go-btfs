package handlers

import (
	"errors"
	s3action "github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/requests"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"net/http"
)

func (h *Handlers) PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	req, err := requests.ParsePutBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestBody)
		return
	}

	// issue: lock for check
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	if err = s3utils.CheckValidBucketNameStrict(req.Bucket); err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidBucketName)
		return
	}

	if !requests.CheckAclPermissionType(&req.ACL) {
		err = responses.ErrNotImplemented
		responses.WriteErrorResponse(w, r, responses.ErrNotImplemented)
		return
	}

	if ok := h.bucsvc.HasBucket(ctx, req.Bucket); ok {
		err = responses.ErrBucketAlreadyExists
		responses.WriteErrorResponseHeadersOnly(w, r, responses.ErrBucketAlreadyExists)
		return
	}

	err = h.bucsvc.CreateBucket(ctx, req.Bucket, req.Region, ack, req.ACL)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInternalError)
		return
	}

	// Make sure to add Location information here only for bucket
	if cp := requests.PathClean(r.URL.Path); cp != "" {
		w.Header().Set(consts.Location, cp) // Clean any trailing slashes.
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
