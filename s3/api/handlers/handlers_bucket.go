package handlers

import (
	"github.com/bittorrent/go-btfs/s3/api/contexts"
	"github.com/bittorrent/go-btfs/s3/api/requests"
	"github.com/bittorrent/go-btfs/s3/api/responses"
	"github.com/bittorrent/go-btfs/s3/api/services/object"
	"net/http"
)

func (h *Handlers) CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.CreateBucketArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseCreateBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	buc, err := h.objsvc.CreateBucket(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteCreateBucketResponse(w, r, buc)
	return
}

func (h *Handlers) HeadBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.GetBucketArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseHeadBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	buc, err := h.objsvc.GetBucket(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteHeadBucketResponse(w, r, buc)
	return
}

func (h *Handlers) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.DeleteBucketArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseDeleteBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	err = h.objsvc.DeleteBucket(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteDeleteBucketResponse(w)
	return
}

func (h *Handlers) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.ListBucketsArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseListBucketsRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	list, err := h.objsvc.ListBuckets(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteListBucketsResponse(w, r, list)
	return
}

func (h *Handlers) PutBucketACLHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.PutBucketACLArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParsePutBucketAclRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	err = h.objsvc.PutBucketACL(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WritePutBucketAclResponse(w, r)
	return
}

func (h *Handlers) GetBucketACLHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.GetBucketACLArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseGetBucketACLRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	acl, err := h.objsvc.GetBucketACL(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteGetBucketACLResponse(w, r, acl)
	return
}
