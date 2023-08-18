// Package handlers is an implementation of s3.Handlerser
package handlers

import (
	"fmt"
	"net/http"

	s3action "github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/routers"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/rs/cors"
)

var _ routers.Handlerser = (*Handlers)(nil)

type Handlers struct {
	corsSvc      CorsService
	authSvc      AuthService
	bucketSvc    BucketService
	objectSvc    ObjectService
	multipartSvc MultipartService
}

func NewHandlers(
	corsSvc CorsService,
	authSvc AuthService,
	bucketSvc BucketService,
	objectSvc ObjectService,
	multipartSvc MultipartService,
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

func (h *Handlers) Sign(handler http.Handler) http.Handler {
	return nil
}

func (h *Handlers) PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("receive request")
	ctx := r.Context()
	req := &PutBucketRequest{}
	err := req.Bind(r)

	defer func() {
		fmt.Println("handle err: ", err)
	}()

	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrInvalidArgument))
		return
	}

	accessKeyRecord, errCode := h.authSvc.VerifySignature(ctx, r)
	if errCode != ErrCodeNone {
		WriteErrorResponse(w, r, errCode)
		return
	}

	//err = h.bucketSvc.CheckACL(accessKeyRecord, req.Bucket, s3action.CreateBucketAction)
	//if err != nil {
	//	WriteErrorResponse(w, r, ToApiError(ctx, ErrNoSuchUserPolicy))
	//	return
	//}

	if err = s3utils.CheckValidBucketNameStrict(req.Bucket); err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrInvalidBucketName))
		return
	}

	fmt.Println("4")
	if !checkPermissionType(req.ACL) {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrNotImplemented))
	}

	fmt.Println("3")
	if ok := h.bucketSvc.HasBucket(r.Context(), req.Bucket); ok {
		WriteErrorResponseHeadersOnly(w, r, ToApiError(ctx, ErrBucketNotFound))
		return
	}

	fmt.Println("2")
	fmt.Println(h.bucketSvc, accessKeyRecord)
	err = h.bucketSvc.CreateBucket(ctx, req.Bucket, req.Region, accessKeyRecord.Key, req.ACL)
	if err != nil {
		log.Errorf("PutBucketHandler create bucket error:%v", err)
		WriteErrorResponse(w, r, ToApiError(ctx, ErrCreateBucket))
		return
	}

	fmt.Println("1")
	// Make sure to add Location information here only for bucket
	if cp := pathClean(r.URL.Path); cp != "" {
		w.Header().Set(consts.Location, cp) // Clean any trailing slashes.
	}

	fmt.Println("0")

	WritePutBucketResponse(w, r)

	return
}

func (h *Handlers) HeadBucketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &HeadBucketRequest{}
	err := req.Bind(r)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrInvalidArgument))
		return
	}

	accessKeyRecord, errCode := h.authSvc.VerifySignature(ctx, r)
	if errCode != ErrCodeNone {
		WriteErrorResponse(w, r, errCode)
		return
	}

	err = h.bucketSvc.CheckACL(accessKeyRecord, req.Bucket, s3action.HeadBucketAction)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrNoSuchUserPolicy))
		return
	}

	if ok := h.bucketSvc.HasBucket(r.Context(), req.Bucket); !ok {
		WriteErrorResponseHeadersOnly(w, r, ToApiError(ctx, ErrBucketNotFound))
		return
	}

	WriteHeadBucketResponse(w, r)
}

func (h *Handlers) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &DeleteBucketRequest{}
	err := req.Bind(r)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrInvalidArgument))
		return
	}

	accessKeyRecord, errCode := h.authSvc.VerifySignature(ctx, r)
	if errCode != ErrCodeNone {
		WriteErrorResponse(w, r, errCode)
		return
	}

	err = h.bucketSvc.CheckACL(accessKeyRecord, req.Bucket, s3action.HeadBucketAction)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrNoSuchUserPolicy))
		return
	}

	//todo check all errors.
	err = h.bucketSvc.DeleteBucket(ctx, req.Bucket)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, err))
		return
	}
	WriteDeleteBucketResponse(w)
}

func (h *Handlers) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &ListBucketsRequest{}
	err := req.Bind(r)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrInvalidArgument))
		return
	}

	accessKeyRecord, errCode := h.authSvc.VerifySignature(ctx, r)
	if errCode != ErrCodeNone {
		WriteErrorResponse(w, r, errCode)
		return
	}

	err = h.bucketSvc.CheckACL(accessKeyRecord, req.Bucket, s3action.ListBucketAction)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrNoSuchUserPolicy))
		return
	}

	//todo check all errors
	bucketMetas, err := h.bucketSvc.GetAllBucketsOfUser(ctx, accessKeyRecord.Key)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, err))
		return
	}

	WriteListBucketsResponse(w, r, bucketMetas)
}

func (h *Handlers) GetBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &GetBucketAclRequest{}
	err := req.Bind(r)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrInvalidArgument))
		return
	}

	accessKeyRecord, errCode := h.authSvc.VerifySignature(ctx, r)
	if errCode != ErrCodeNone {
		WriteErrorResponse(w, r, errCode)
		return
	}

	err = h.bucketSvc.CheckACL(accessKeyRecord, req.Bucket, s3action.GetBucketAclAction)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrNoSuchUserPolicy))
		return
	}

	if !h.bucketSvc.HasBucket(ctx, req.Bucket) {
		WriteErrorResponseHeadersOnly(w, r, ToApiError(ctx, ErrBucketNotFound))
		return
	}
	//todo check all errors
	acl, err := h.bucketSvc.GetBucketAcl(ctx, req.Bucket)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, err))
		return
	}

	WriteGetBucketAclResponse(w, r, accessKeyRecord, acl)
}

func (h *Handlers) PutBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &PutBucketAclRequest{}
	err := req.Bind(r)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrInvalidArgument))
		return
	}

	accessKeyRecord, errCode := h.authSvc.VerifySignature(ctx, r)
	if errCode != ErrCodeNone {
		WriteErrorResponse(w, r, errCode)
		return
	}

	err = h.bucketSvc.CheckACL(accessKeyRecord, req.Bucket, s3action.PutBucketAclAction)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrNoSuchUserPolicy))
		return
	}

	if !checkPermissionType(req.ACL) || req.ACL == "" {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrNotImplemented))
		return
	}

	//todo check all errors
	err = h.bucketSvc.UpdateBucketAcl(ctx, req.Bucket, req.ACL)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, err))
		return
	}

	//todo check no return?
	WritePutBucketAclResponse(w, r)
}
