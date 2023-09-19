package handlers

import (
	"github.com/bittorrent/go-btfs/s3/api/contexts"
	"github.com/bittorrent/go-btfs/s3/api/requests"
	"github.com/bittorrent/go-btfs/s3/api/responses"
	"github.com/bittorrent/go-btfs/s3/api/services/object"
	"net/http"
)

// PutObjectHandler .
func (h *Handlers) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.PutObjectArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParsePutObjectRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	obj, err := h.objsvc.PutObject(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WritePutObjectResponse(w, r, obj)
	return
}

// CopyObjectHandler .
func (h *Handlers) CopyObjectHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.CopyObjectArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseCopyObjectRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	obj, err := h.objsvc.CopyObject(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteCopyObjectResponse(w, r, obj)
	return
}

// HeadObjectHandler .
func (h *Handlers) HeadObjectHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.GetObjectArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseHeadObjectRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	obj, _, err := h.objsvc.GetObject(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteHeadObjectResponse(w, r, obj)
	return
}

// GetObjectHandler .
func (h *Handlers) GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.GetObjectArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseGetObjectRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	obj, body, err := h.objsvc.GetObject(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteGetObjectResponse(w, r, obj, body)
	return
}

// DeleteObjectHandler .
func (h *Handlers) DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.DeleteObjectArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseDeleteObjectRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}
	err = h.objsvc.DeleteObject(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteDeleteObjectResponse(w, r, nil)
	return
}

// DeleteObjectsHandler .
func (h *Handlers) DeleteObjectsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.DeleteObjectsArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseDeleteObjectsRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	deletes, err := h.objsvc.DeleteObjects(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteDeleteObjectsResponse(w, r, h.toResponseErr, deletes)
	return
}

// ListObjectsHandler .
func (h *Handlers) ListObjectsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.ListObjectsArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseListObjectsRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	list, err := h.objsvc.ListObjects(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteListObjectsResponse(w, r, list)
	return
}

// ListObjectsV2Handler .
func (h *Handlers) ListObjectsV2Handler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.ListObjectsV2Args
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseListObjectsV2Request(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	list, err := h.objsvc.ListObjectsV2(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteListObjectsV2Response(w, r, list)
	return
}

// GetObjectACLHandler - GET Object ACL
func (h *Handlers) GetObjectACLHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var args *object.GetObjectACLArgs
	defer func() {
		contexts.SetHandleInf(r, h.name(), err, args)
	}()

	args, err = requests.ParseGetObjectACLRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	acl, err := h.objsvc.GetObjectACL(r.Context(), args)
	if err != nil {
		responses.WriteErrorResponse(w, r, h.toResponseErr(err))
		return
	}

	responses.WriteGetObjectACLResponse(w, r, acl)
	return
}
