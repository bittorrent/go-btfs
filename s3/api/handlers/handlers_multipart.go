package handlers

import (
	"github.com/bittorrent/go-btfs/s3/api/contexts"
	"github.com/bittorrent/go-btfs/s3/api/requests"
	"github.com/bittorrent/go-btfs/s3/api/responses"
	"net/http"
)

func (h *Handlers) CreateMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		contexts.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseCreateMultipartUploadRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	multipart, err := h.objsvc.CreateMultipartUpload(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteCreateMultipartUploadResponse(w, r, multipart)
	return
}

func (h *Handlers) UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		contexts.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseUploadPartRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	part, err := h.objsvc.UploadPart(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteUploadPartResponse(w, r, part)
	return
}

func (h *Handlers) AbortMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		contexts.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseAbortMultipartUploadRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	err = h.objsvc.AbortMultipartUpload(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteAbortMultipartUploadResponse(w, r)
	return
}

func (h *Handlers) CompleteMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		contexts.SetHandleInf(r, h.name(), err)
	}()

	args, err := requests.ParseCompleteMultipartUploadRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	obj, err := h.objsvc.CompleteMultiPartUpload(ctx, args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteCompleteMultipartUploadResponse(w, r, obj)
	return
}
