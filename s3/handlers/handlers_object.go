package handlers

import (
	"fmt"
	s3action "github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/etag"
	"github.com/bittorrent/go-btfs/s3/requests"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"net/http"
)

// PutObjectHandler http://docs.aws.amazon.com/AmazonS3/latest/dev/UploadingObjects.html
func (h *Handlers) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	// X-Amz-Copy-Source shouldn't be set for this call.
	if _, ok := r.Header[consts.AmzCopySource]; ok {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidCopySource)
		return
	}

	buc, obj, err := requests.ParseBucketAndObject(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
		return
	}

	clientETag, err := etag.FromContentMD5(r.Header)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidDigest)
		return
	}
	_ = clientETag

	size := r.ContentLength
	// todo: streaming signed

	if size == -1 {
		responses.WriteErrorResponse(w, r, responses.ErrMissingContentLength)
		return
	}
	if size == 0 {
		responses.WriteErrorResponse(w, r, responses.ErrEntityTooSmall)
		return
	}

	if size > consts.MaxObjectSize {
		responses.WriteErrorResponse(w, r, responses.ErrEntityTooLarge)
		return
	}

	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	err = h.bucsvc.CheckACL(ack, buc, s3action.PutObjectAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	// todo: convert error
	err = s3utils.CheckPutObjectArgs(ctx, buc, obj)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	// todo
	fmt.Println("need put object...", buc, obj)
}
