package handlers

import (
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/requests"
	"github.com/bittorrent/go-btfs/s3/responses"
	"net/http"
)

func (h *Handlers) CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseCreateBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	buc, err := h.objsvc.CreateBucket(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	responses.WriteCreateBucketResponse(w, r, buc)
	return
}

func (h *Handlers) HeadBucketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseHeadBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	buc, err := h.objsvc.GetBucket(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	responses.WriteHeadBucketResponse(w, r, buc)
	return
}

func (h *Handlers) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseDeleteBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	err = h.objsvc.DeleteBucket(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
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

	args, err := requests.ParseListBucketsRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	list, err := h.objsvc.ListBuckets(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	responses.WriteListBucketsResponse(w, r, args.AccessKey, list)
	return
}

func (h *Handlers) GetBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseGetBucketACLRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	acl, err := h.objsvc.GetBucketACL(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	responses.WriteGetBucketACLResponse(w, r, args.AccessKey, acl)
	return
}

func (h *Handlers) PutBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParsePutBucketAclRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	err = h.objsvc.PutBucketACL(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toRespErr(err))
		return
	}

	responses.WritePutBucketAclResponse(w, r)
	return
}
