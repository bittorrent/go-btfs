// Package handlers is an implementation of s3.Handlerser
package handlers

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/ctxmu"
	"github.com/bittorrent/go-btfs/s3/etag"
	"github.com/bittorrent/go-btfs/s3/requests"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"github.com/bittorrent/go-btfs/s3/services/bucket"
	"github.com/bittorrent/go-btfs/s3/services/cors"
	"github.com/bittorrent/go-btfs/s3/services/sign"
	"net/http"
	"runtime"

	s3action "github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	rscors "github.com/rs/cors"
)

var _ Handlerser = (*Handlers)(nil)

type Handlers struct {
	corsvc cors.Service
	acksvc accesskey.Service
	sigsvc sign.Service
	bucsvc bucket.Service
	nslock ctxmu.MultiCtxRWLocker
}

func NewHandlers(corsvc cors.Service, acksvc accesskey.Service, sigsvc sign.Service, bucsvc bucket.Service, options ...Option) (handlers *Handlers) {
	handlers = &Handlers{
		corsvc: corsvc,
		acksvc: acksvc,
		sigsvc: sigsvc,
		bucsvc: bucsvc,
		nslock: ctxmu.NewDefaultMultiCtxRWMutex(),
	}
	for _, option := range options {
		option(handlers)
	}
	return
}

func (h *Handlers) Cors(handler http.Handler) http.Handler {
	return rscors.New(rscors.Options{
		AllowedOrigins:   h.corsvc.GetAllowOrigins(),
		AllowedMethods:   h.corsvc.GetAllowMethods(),
		AllowedHeaders:   h.corsvc.GetAllowHeaders(),
		ExposedHeaders:   h.corsvc.GetAllowHeaders(),
		AllowCredentials: true,
	}).Handler(handler)
}

func (h *Handlers) Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[REQ] <%4s> | %s\n", r.Method, r.URL)
		handler.ServeHTTP(w, r)
		hname, herr := cctx.GetHandleInf(r)
		fmt.Printf("[RSP] <%4s> | %s | %s | %v\n", r.Method, r.URL, hname, herr)
	})
}

func (h *Handlers) Sign(handler http.Handler) http.Handler {
	h.sigsvc.SetSecretGetter(func(key string) (secret string, exists, enable bool, err error) {
		ack, err := h.acksvc.Get(key)
		if errors.Is(err, accesskey.ErrNotFound) {
			exists = false
			enable = true
			err = nil
			return
		}
		if err != nil {
			return
		}
		exists = true
		secret = ack.Secret
		enable = ack.Enable
		return
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err *responses.Error
		defer func() {
			if err != nil {
				cctx.SetHandleInf(r, fnName(), err)
			}
		}()

		ack, err := h.sigsvc.VerifyRequestSignature(r)
		if err != nil {
			responses.WriteErrorResponse(w, r, err)
			return
		}

		cctx.SetAccessKey(r, ack)

		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, fnName(), err)
	}()

	req, err := requests.ParsePutBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestBody)
		return
	}

	// issue: lock for check
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	if err = s3utils.CheckValidBucketNameStrict(req.Bucket); err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidBucketName)
		return
	}

	if !requests.CheckAclPermissionType(&req.ACL) {
		err = responses.ErrNotImplemented
		responses.WriteErrorResponse(w, r, responses.ErrNotImplemented)
		return
	}

	if ok := h.bucsvc.HasBucket(ctx, req.Bucket); ok {
		err = responses.ErrBucketAlreadyExists
		responses.WriteErrorResponseHeadersOnly(w, r, responses.ErrBucketAlreadyExists)
		return
	}

	err = h.bucsvc.CreateBucket(ctx, req.Bucket, req.Region, ack, req.ACL)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInternalError)
		return
	}

	// Make sure to add Location information here only for bucket
	if cp := requests.PathClean(r.URL.Path); cp != "" {
		w.Header().Set(consts.Location, cp) // Clean any trailing slashes.
	}

	responses.WritePutBucketResponse(w, r)

	return
}

func (h *Handlers) HeadBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, fnName(), err)
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

func (h *Handlers) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, fnName(), err)
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
		cctx.SetHandleInf(r, fnName(), err)
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
		cctx.SetHandleInf(r, fnName(), err)
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
		cctx.SetHandleInf(r, fnName(), err)
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

// PutObjectHandler http://docs.aws.amazon.com/AmazonS3/latest/dev/UploadingObjects.html
func (h *Handlers) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, fnName(), err)
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

func fnName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
