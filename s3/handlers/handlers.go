// Package handlers is an implementation of s3.Handlerser
package handlers

import (
	"fmt"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/etag"
	"github.com/bittorrent/go-btfs/s3/requests"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/services"
	"github.com/bittorrent/go-btfs/s3/services/auth"
	"github.com/bittorrent/go-btfs/s3/services/bucket"
	"github.com/bittorrent/go-btfs/s3/services/cors"
	"net/http"
	"runtime"

	s3action "github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	rscors "github.com/rs/cors"
)

var _ Handlerser = (*Handlers)(nil)

type Handlers struct {
	corsSvc   cors.Service
	authSvc   auth.Service
	bucketSvc bucket.Service
}

func NewHandlers(
	corsSvc cors.Service,
	authSvc auth.Service,
	bucketSvc bucket.Service,
	options ...Option,
) (handlers *Handlers) {
	handlers = &Handlers{
		corsSvc:   corsSvc,
		authSvc:   authSvc,
		bucketSvc: bucketSvc,
	}
	for _, option := range options {
		option(handlers)
	}
	return
}

func (h *Handlers) Cors(handler http.Handler) http.Handler {
	return rscors.New(rscors.Options{
		AllowedOrigins:   h.corsSvc.GetAllowOrigins(),
		AllowedMethods:   h.corsSvc.GetAllowMethods(),
		AllowedHeaders:   h.corsSvc.GetAllowHeaders(),
		ExposedHeaders:   h.corsSvc.GetAllowHeaders(),
		AllowCredentials: true,
	}).Handler(handler)
}

func (h *Handlers) Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[REQ] %4s | %s\n", r.Method, r.URL)
		handler.ServeHTTP(w, r)
		hname, herr := cctx.GetHandleInf(r)
		fmt.Printf("[RSP] %4s | %s | %s | %v\n", r.Method, r.URL, hname, herr)
	})
}

func (h *Handlers) Auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				cctx.SetHandleInf(r, fnName(), err)
			}
		}()

		ack, err := h.authSvc.VerifySignature(r.Context(), r)
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

	req, err := requests.ParsePubBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, services.RespErrInvalidRequestBody)
		return
	}

	ctx := r.Context()

	if err = s3utils.CheckValidBucketNameStrict(req.Bucket); err != nil {
		responses.WriteErrorResponse(w, r, services.RespErrInvalidBucketName)
		return
	}

	if !requests.CheckAclPermissionType(&req.ACL) {
		err = services.RespErrNotImplemented
		responses.WriteErrorResponse(w, r, services.RespErrNotImplemented)
		return
	}

	if ok := h.bucketSvc.HasBucket(ctx, req.Bucket); ok {
		err = services.RespErrBucketAlreadyExists
		responses.WriteErrorResponseHeadersOnly(w, r, services.RespErrBucketAlreadyExists)
		return
	}

	err = h.bucketSvc.CreateBucket(ctx, req.Bucket, req.Region, cctx.GetAccessKey(r).Key, req.ACL)
	if err != nil {
		responses.WriteErrorResponse(w, r, services.RespErrInternalError)
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
		responses.WriteErrorResponse(w, r, services.RespErrInvalidRequestBody)
		return
	}

	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	err = h.bucketSvc.CheckACL(ack, req.Bucket, s3action.HeadBucketAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	if ok := h.bucketSvc.HasBucket(ctx, req.Bucket); !ok {
		responses.WriteErrorResponseHeadersOnly(w, r, services.RespErrNoSuchBucket)
		return
	}

	responses.WriteHeadBucketResponse(w, r)
}

func (h *Handlers) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, fnName(), err)
	}()

	req := &requests.DeleteBucketRequest{}
	err = req.Bind(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, services.RespErrInvalidRequestBody)
		return
	}

	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	err = h.bucketSvc.CheckACL(ack, req.Bucket, s3action.HeadBucketAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	//todo check all errors.
	err = h.bucketSvc.DeleteBucket(ctx, req.Bucket)
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

	req := &requests.ListBucketsRequest{}
	err = req.Bind(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, services.RespErrInvalidRequestBody)
		return
	}

	ack := cctx.GetAccessKey(r)

	err = h.bucketSvc.CheckACL(ack, req.Bucket, s3action.ListBucketAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	//todo check all errors
	bucketMetas, err := h.bucketSvc.GetAllBucketsOfUser(ack.Key)
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

	req := &requests.GetBucketAclRequest{}
	err = req.Bind(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, services.RespErrInvalidRequestBody)
		return
	}

	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	err = h.bucketSvc.CheckACL(ack, req.Bucket, s3action.GetBucketAclAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	if !h.bucketSvc.HasBucket(ctx, req.Bucket) {
		responses.WriteErrorResponseHeadersOnly(w, r, services.RespErrNoSuchBucket)
		return
	}
	//todo check all errors
	acl, err := h.bucketSvc.GetBucketAcl(ctx, req.Bucket)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	responses.WriteGetBucketAclResponse(w, r, ack.Key, acl)
}

func (h *Handlers) PutBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, fnName(), err)
	}()

	req := &requests.PutBucketAclRequest{}
	err = req.Bind(r)
	if err != nil || len(req.ACL) == 0 || len(req.Bucket) == 0 {
		responses.WriteErrorResponse(w, r, services.RespErrInvalidRequestBody)
		return
	}

	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	err = h.bucketSvc.CheckACL(ack, req.Bucket, s3action.PutBucketAclAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		return
	}

	if !requests.CheckAclPermissionType(&req.ACL) {
		responses.WriteErrorResponse(w, r, services.RespErrNotImplemented)
		return
	}

	//todo check all errors
	err = h.bucketSvc.UpdateBucketAcl(ctx, req.Bucket, req.ACL)
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
		responses.WriteErrorResponse(w, r, services.RespErrInvalidCopySource)
		return
	}

	buc, obj, err := requests.ParseBucketAndObject(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, services.RespErrInvalidRequestParameter)
		return
	}

	clientETag, err := etag.FromContentMD5(r.Header)
	if err != nil {
		responses.WriteErrorResponse(w, r, services.RespErrInvalidDigest)
		return
	}
	_ = clientETag

	size := r.ContentLength
	// todo: streaming signed

	if size == -1 {
		responses.WriteErrorResponse(w, r, services.RespErrMissingContentLength)
		return
	}
	if size == 0 {
		responses.WriteErrorResponse(w, r, services.RespErrEntityTooSmall)
		return
	}

	if size > consts.MaxObjectSize {
		responses.WriteErrorResponse(w, r, services.RespErrEntityTooLarge)
		return
	}

	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	err = h.bucketSvc.CheckACL(ack, buc, s3action.PutObjectAction)
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
