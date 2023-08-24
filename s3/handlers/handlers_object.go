package handlers

import (
	"errors"
	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/requests"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/bittorrent/go-btfs/s3/services/bucket"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"net/http"
	"time"
)

const lockWaitTimeout = 5 * time.Minute

func (h *Handlers) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	// X-Amz-Copy-Source shouldn't be set for this call.
	if _, ok := r.Header[consts.AmzCopySource]; ok {
		err = errors.New("shouldn't be copy")
		responses.WriteErrorResponse(w, r, responses.ErrInvalidCopySource)
		return
	}

	aclHeader := r.Header.Get(consts.AmzACL)
	if aclHeader != "" {
		err = errors.New("object acl can only set to default")
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}

	bucname, objname, err := requests.ParseBucketAndObject(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}

	err = s3utils.CheckPutObjectArgs(ctx, bucname, objname)
	if err != nil { // todo: convert error
		responses.WriteErrorResponse(w, r, err)
		return
	}

	meta, err := extractMetadata(ctx, r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequest)
		return
	}

	if r.ContentLength == 0 {
		responses.WriteErrorResponse(w, r, responses.ErrEntityTooSmall)
		return
	}

	hrdr, ok := r.Body.(*hash.Reader)
	if !ok {
		responses.WriteErrorResponse(w, r, responses.ErrInternalError)
		return
	}

	// rlock bucket
	runlock, err := h.rlock(ctx, bucname, w, r)
	if err != nil {
		return
	}
	defer runlock()

	// lock object
	unlock, err := h.lock(ctx, bucname+"/"+objname, w, r)
	if err != nil {
		return
	}
	defer unlock()

	err = h.bucsvc.CheckACL(ack, bucname, action.PutObjectAction)
	if errors.Is(err, bucket.ErrNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	obj, err := h.objsvc.PutObject(ctx, bucname, objname, hrdr, r.ContentLength, meta)

	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	responses.WritePutObjectResponse(w, r, obj)

	return
}
