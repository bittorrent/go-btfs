// Package handlers is an implementation of s3.Handlerser
package handlers

import (
	"fmt"
	"github.com/bittorrent/go-btfs/s3/handlers/cctx"
	"github.com/bittorrent/go-btfs/s3/handlers/requests"
	"github.com/bittorrent/go-btfs/s3/handlers/responses"
	"github.com/bittorrent/go-btfs/s3/services"
	"net/http"
	"runtime"

	s3action "github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/rs/cors"
)

var _ Handlerser = (*Handlers)(nil)

type Handlers struct {
	corsSvc      services.CorsService
	authSvc      services.AuthService
	bucketSvc    services.BucketService
	objectSvc    services.ObjectService
	multipartSvc services.MultipartService
}

func NewHandlers(
	corsSvc services.CorsService,
	authSvc services.AuthService,
	bucketSvc services.BucketService,
	objectSvc services.ObjectService,
	multipartSvc services.MultipartService,
	options ...Option,
) (handlers *Handlers) {
	handlers = &Handlers{
		corsSvc:      corsSvc,
		authSvc:      authSvc,
		bucketSvc:    bucketSvc,
		objectSvc:    objectSvc,
		multipartSvc: multipartSvc,
	}
	for _, option := range options {
		option(handlers)
	}
	return
}

func (h *Handlers) Cors(handler http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   h.corsSvc.GetAllowOrigins(),
		AllowedMethods:   h.corsSvc.GetAllowMethods(),
		AllowedHeaders:   h.corsSvc.GetAllowHeaders(),
		ExposedHeaders:   h.corsSvc.GetAllowHeaders(),
		AllowCredentials: true,
	}).Handler(handler)
}

func (h *Handlers) Auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ack, err := h.authSvc.VerifySignature(r.Context(), r)
		if err != nil {
			responses.WriteErrorResponse(w, r, err)
			return
		}
		cctx.SetAccessKey(r, ack)
		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		hname, herr := cctx.GetHandleInf(r)
		fmt.Printf("[%-4s] %s | %s | %v\n", r.Method, r.URL, hname, herr)
	})
}

func (h *Handlers) PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, fnName(), err)
	}()

	req, err := requests.ParsePubBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, services.ErrInvalidRequestBody)
		return
	}

	ctx := r.Context()

	if err = s3utils.CheckValidBucketNameStrict(req.Bucket); err != nil {
		responses.WriteErrorResponse(w, r, services.ErrInvalidBucketName)
		return
	}

	if !requests.CheckAclPermissionType(&req.ACL) {
		err = services.ErrNotImplemented
		responses.WriteErrorResponse(w, r, services.ErrNotImplemented)
		return
	}

	if ok := h.bucketSvc.HasBucket(ctx, req.Bucket); ok {
		err = services.ErrBucketAlreadyExists
		responses.WriteErrorResponseHeadersOnly(w, r, services.ErrBucketAlreadyExists)
		return
	}

	err = h.bucketSvc.CreateBucket(ctx, req.Bucket, req.Region, cctx.GetAccessKey(r).Key, req.ACL)
	if err != nil {
		responses.WriteErrorResponse(w, r, services.ErrInternalError)
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
		responses.WriteErrorResponse(w, r, services.ErrInvalidRequestBody)
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
		responses.WriteErrorResponseHeadersOnly(w, r, services.ErrNoSuchBucket)
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
		responses.WriteErrorResponse(w, r, services.ErrInvalidRequestBody)
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
		responses.WriteErrorResponse(w, r, services.ErrInvalidRequestBody)
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
		responses.WriteErrorResponse(w, r, services.ErrInvalidRequestBody)
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
		responses.WriteErrorResponseHeadersOnly(w, r, services.ErrNoSuchBucket)
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
		responses.WriteErrorResponse(w, r, services.ErrInvalidRequestBody)
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
		responses.WriteErrorResponse(w, r, services.ErrNotImplemented)
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

func fnName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
