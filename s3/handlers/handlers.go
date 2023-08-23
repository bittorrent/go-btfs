// Package handlers is an implementation of s3.Handlerser
package handlers

import (
	"fmt"
	"github.com/bittorrent/go-btfs/s3/handlers/cctx"
	"github.com/bittorrent/go-btfs/s3/handlers/requests"
	"github.com/bittorrent/go-btfs/s3/handlers/responses"
	"github.com/bittorrent/go-btfs/s3/services"
	"net/http"

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
		ack, rErr := h.authSvc.VerifySignature(r.Context(), r)
		if rErr != nil {
			responses.WriteErrorResponse(w, r, rErr)
			return
		}
		cctx.SetAccessKey(r, ack)
		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("... PutBucketHandler: begin")

	ctx := r.Context()

	req, err := requests.ParsePubBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestBody)
		return
	}

	//err = h.bucketSvc.CheckACL(accessKeyRecord, req.Bucket, s3action.CreateBucketAction)
	//if err != nil {
	//	WriteErrorResponse(w, r, ToApiError(ctx, ErrNoSuchUserPolicy))
	//	return
	//}

	if err = s3utils.CheckValidBucketNameStrict(req.Bucket); err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidBucketName)
		return
	}

	fmt.Println("4")
	if !requests.CheckAclPermissionType(&req.ACL) {
		responses.WriteErrorResponse(w, r, responses.ErrNotImplemented)
		return
	}

	fmt.Println("3")
	if ok := h.bucketSvc.HasBucket(ctx, req.Bucket); ok {
		responses.WriteErrorResponseHeadersOnly(w, r, responses.ErrBucketAlreadyExists)
		return
	}

	fmt.Println("2")
	err = h.bucketSvc.CreateBucket(ctx, req.Bucket, req.Region, cctx.GetAccessKey(r).Key, req.ACL)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInternalError)
		return
	}

	fmt.Println("1")
	// Make sure to add Location information here only for bucket
	if cp := requests.PathClean(r.URL.Path); cp != "" {
		w.Header().Set(consts.Location, cp) // Clean any trailing slashes.
	}

	fmt.Println("0")

	responses.WritePutBucketResponse(w, r)

	return
}

func (h *Handlers) HeadBucketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("... HeadBucketHandler: begin")

	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	req := &requests.HeadBucketRequest{}
	err := req.Bind(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestBody)
		return
	}

	fmt.Println("... head bucket ", req)

	err = h.bucketSvc.CheckACL(ack, req.Bucket, s3action.HeadBucketAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, ToResponseErr(err))
		return
	}

	if ok := h.bucketSvc.HasBucket(ctx, req.Bucket); !ok {
		responses.WriteErrorResponseHeadersOnly(w, r, responses.ErrNoSuchBucket)
		return
	}

	responses.WriteHeadBucketResponse(w, r)
}

func (h *Handlers) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("... DeleteBucketHandler: begin")

	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	req := &requests.DeleteBucketRequest{}
	err := req.Bind(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestBody)
		return
	}

	err = h.bucketSvc.CheckACL(ack, req.Bucket, s3action.HeadBucketAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, ToResponseErr(err))
		return
	}

	//todo check all errors.
	err = h.bucketSvc.DeleteBucket(ctx, req.Bucket)
	if err != nil {
		responses.WriteErrorResponse(w, r, ToResponseErr(err))
		return
	}

	responses.WriteDeleteBucketResponse(w)
}

func (h *Handlers) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("... ListBucketsHandler: begin")

	ack := cctx.GetAccessKey(r)

	req := &requests.ListBucketsRequest{}
	err := req.Bind(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestBody)
		return
	}

	err = h.bucketSvc.CheckACL(ack, req.Bucket, s3action.ListBucketAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, ToResponseErr(err))
		return
	}

	//todo check all errors
	bucketMetas, err := h.bucketSvc.GetAllBucketsOfUser(ack.Key)
	if err != nil {
		responses.WriteErrorResponse(w, r, ToResponseErr(err))
		return
	}

	responses.WriteListBucketsResponse(w, r, bucketMetas)
}

func (h *Handlers) GetBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("... get acl req: begin")

	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	req := &requests.GetBucketAclRequest{}
	err := req.Bind(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestBody)
		return
	}

	fmt.Println("... get acl req: ", req)

	err = h.bucketSvc.CheckACL(ack, req.Bucket, s3action.GetBucketAclAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, ToResponseErr(err))
		return
	}

	if !h.bucketSvc.HasBucket(ctx, req.Bucket) {
		responses.WriteErrorResponseHeadersOnly(w, r, responses.ErrNoSuchBucket)
		return
	}
	//todo check all errors
	acl, err := h.bucketSvc.GetBucketAcl(ctx, req.Bucket)
	if err != nil {
		responses.WriteErrorResponse(w, r, ToResponseErr(err))
		return
	}

	fmt.Println("... get acl = ", req)

	responses.WriteGetBucketAclResponse(w, r, ack, acl)
}

func (h *Handlers) PutBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("... PutBucketAclHandler: begin")

	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	req := &requests.PutBucketAclRequest{}
	err := req.Bind(r)
	if err != nil || len(req.ACL) == 0 || len(req.Bucket) == 0 {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestBody)
		return
	}

	err = h.bucketSvc.CheckACL(ack, req.Bucket, s3action.PutBucketAclAction)
	if err != nil {
		responses.WriteErrorResponse(w, r, ToResponseErr(err))
		return
	}

	if !requests.CheckAclPermissionType(&req.ACL) {
		responses.WriteErrorResponse(w, r, responses.ErrNotImplemented)
		return
	}

	//todo check all errors
	err = h.bucketSvc.UpdateBucketAcl(ctx, req.Bucket, req.ACL)
	if err != nil {
		responses.WriteErrorResponse(w, r, ToResponseErr(err))
		return
	}

	//todo check no return?
	responses.WritePutBucketAclResponse(w, r)
}
