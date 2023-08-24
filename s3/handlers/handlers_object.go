package handlers

import (
	"context"
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

// PutObjectHandler http://docs.aws.amazon.com/AmazonS3/latest/dev/UploadingObjects.html
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

	// lock
	runlock, err := h.rlock(ctx, bucname, w, r)
	if err != nil {
		return
	}
	defer runlock()

	err = h.bucsvc.CheckACL(ack, bucname, action.PutObjectAction)
	if errors.Is(err, bucket.ErrNotFound) {
		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
		return
	}

	hrdr, ok := r.Body.(*hash.Reader)
	if !ok {
		responses.WriteErrorResponse(w, r, responses.ErrInternalError)
		return
	}

	metadata, err := extractMetadata(ctx, r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequest)
		return
	}

	obj, err := h.objsvc.StoreObject(ctx, bucname, objname, hrdr, r.ContentLength, metadata)

	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	responses.WritePutObjectResponse(w, r, obj, false)
}

func (h *Handlers) rlock(ctx context.Context, key string, w http.ResponseWriter, r *http.Request) (runlock func(), err error) {
	ctx, cancel := context.WithTimeout(ctx, lockWaitTimeout)
	err = h.nslock.RLock(ctx, key)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		cancel()
		return
	}
	runlock = func() {
		h.nslock.RUnlock(key)
		cancel()
	}
	return
}
