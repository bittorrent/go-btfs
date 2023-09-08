package handlers

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/protocol"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"net/http"
	"sort"
)

func (h *Handlers) CreateMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	var input s3.CreateMultipartUploadInput

	err = protocol.ParseRequest(r, &input)
	if err != nil {
		rerr := responses.ErrBadRequest
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	bucname, objname := *input.Bucket, *input.Key

	err = s3utils.CheckNewMultipartArgs(ctx, bucname, objname)
	if err != nil {
		rerr := h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	meta := input.Metadata

	mtp, err := h.objsvc.CreateMultipartUpload(ctx, ack, bucname, objname, meta)
	if err != nil {
		rerr := h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	output := new(s3.CreateMultipartUploadOutput)
	output.SetBucket(bucname)
	output.SetKey(objname)
	output.SetUploadId(mtp.UploadID)

	responses.WriteSuccessResponse(w, output, "InitiateMultipartUploadResult")

	return
}

func (h *Handlers) UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	var input s3.UploadPartInput

	err = protocol.ParseRequest(r, &input)
	if err != nil {
		rerr := responses.ErrBadRequest
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	bucname, objname := *input.Bucket, *input.Key

	err = s3utils.CheckPutObjectPartArgs(ctx, bucname, objname)
	if err != nil {
		rerr := h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	uploadId, partId := *input.UploadId, int(*input.PartNumber)
	if partId > consts.MaxPartID {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidMaxParts)
		return
	}

	size := r.ContentLength

	if size == 0 {
		responses.WriteErrorResponse(w, r, responses.ErrEntityTooSmall)
		return
	}

	if size > consts.MaxPartSize {
		responses.WriteErrorResponse(w, r, responses.ErrEntityTooLarge)
		return
	}

	hrdr, ok := r.Body.(*hash.Reader)
	if !ok {
		responses.WriteErrorResponse(w, r, responses.ErrInternalError)
		return
	}

	part, err := h.objsvc.UploadPart(ctx, ack, bucname, objname, uploadId, partId, hrdr, size)
	if err != nil {
		rerr := h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	output := new(s3.UploadPartOutput)
	output.SetETag(`"` + part.ETag + `"`)
	w.Header().Set(consts.Cid, part.CID)

	responses.WriteSuccessResponse(w, output, "")

	return
}

func (h *Handlers) AbortMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	var input s3.AbortMultipartUploadInput

	err = protocol.ParseRequest(r, &input)
	if err != nil {
		rerr := responses.ErrBadRequest
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	bucname, objname := *input.Bucket, *input.Key

	err = s3utils.CheckAbortMultipartArgs(ctx, bucname, objname)
	if err != nil {
		rerr := h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	uploadId := *input.UploadId

	err = h.objsvc.AbortMultipartUpload(ctx, ack, bucname, objname, uploadId)
	if err != nil {
		rerr := h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	output := new(s3.AbortMultipartUploadOutput)

	responses.WriteSuccessResponse(w, output, "")

	return
}

func (h *Handlers) CompleteMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	var input s3.CompleteMultipartUploadInput

	err = protocol.ParseRequest(r, &input)
	if err != nil {
		rerr := responses.ErrBadRequest
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	bucname, objname := *input.Bucket, *input.Key

	err = s3utils.CheckCompleteMultipartArgs(ctx, bucname, objname)
	if err != nil {
		rerr := h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	// Content-Length is required and should be non-zero
	if r.ContentLength <= 0 {
		responses.WriteErrorResponse(w, r, responses.ErrMissingContentLength)
		return
	}

	if len(input.MultipartUpload.Parts) == 0 {
		rerr := responses.ErrMalformedXML
		err = rerr
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	complUpload := new(object.CompleteMultipartUpload)

	for _, part := range input.MultipartUpload.Parts {
		complUpload.Parts = append(complUpload.Parts, &object.CompletePart{
			PartNumber: int(*part.PartNumber),
			ETag:       *part.ETag,
		})

	}

	if !sort.IsSorted(object.CompletedParts(complUpload.Parts)) {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidPartOrder)
		return
	}

	uploadId := *input.UploadId

	obj, err := h.objsvc.CompleteMultiPartUpload(ctx, ack, bucname, objname, uploadId, complUpload.Parts)
	if err != nil {
		rerr := h.respErr(err)
		responses.WriteErrorResponse(w, r, rerr)
		return
	}

	output := new(s3.CompleteMultipartUploadOutput)
	output.SetBucket(bucname)
	output.SetKey(objname)
	output.SetETag(`"` + obj.ETag + `"`)
	w.Header().Set(consts.Cid, obj.CID)

	responses.WriteSuccessResponse(w, output, "")
}

