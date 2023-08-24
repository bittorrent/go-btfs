package handlers

import (
	"errors"
	s3action "github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/requests"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/services/bucket"
	"net/http"
)

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
	if errors.Is(err, bucket.ErrNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrAccessDenied)
		return
	}

	responses.WriteHeadBucketResponse(w, r)
}
