package handlers

import (
	"github.com/bittorrent/go-btfs/s3/api/contexts"
	"github.com/bittorrent/go-btfs/s3/api/requests"
	"github.com/bittorrent/go-btfs/s3/api/responses"
	"github.com/bittorrent/go-btfs/s3/api/services/object"
	"net/http"
)

func (h *Handlers) CreateMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.CreateMultipartUploadArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseCreateMultipartUploadRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	multipart, err := h.objsvc.CreateMultipartUpload(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteCreateMultipartUploadResponse(w, r, multipart)
	return
}

func (h *Handlers) UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.UploadPartArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseUploadPartRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	part, err := h.objsvc.UploadPart(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteUploadPartResponse(w, r, part)
	return
}

func (h *Handlers) AbortMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.AbortMultipartUploadArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseAbortMultipartUploadRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	err = h.objsvc.AbortMultipartUpload(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteAbortMultipartUploadResponse(w, r)
	return
}

func (h *Handlers) CompleteMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.CompleteMultipartUploadArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseCompleteMultipartUploadRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	obj, err := h.objsvc.CompleteMultiPartUpload(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteCompleteMultipartUploadResponse(w, r, obj)
	return
}
