// Package handlers is an implementation of s3.Handlerser
package handlers

import (
	"net/http"

	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/policy"
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

func (handlers *Handlers) Cors(handler http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   handlers.corsSvc.GetAllowOrigins(),
		AllowedMethods:   handlers.corsSvc.GetAllowMethods(),
		AllowedHeaders:   handlers.corsSvc.GetAllowHeaders(),
		ExposedHeaders:   handlers.corsSvc.GetAllowHeaders(),
		AllowCredentials: true,
	}).Handler(handler)
}

func (handlers *Handlers) Sign(handler http.Handler) http.Handler {
	return nil
}

//func (handlers *Handlers) parsePutObjectReq(r *http.Request) (arg *PutObjectReq, err error) {
//	return
//}
//
//func (handlers *Handlers) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
//	req := &PutObjectRequest{}
//	err := req.Bind(r)
//	if err != nil {
//		return
//	}
//	//....
//
//	WritePutObjectResponse(w, object)
//
//	return
//}

func (handlers *Handlers) PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &PutBucketRequest{}
	err := req.Bind(r)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrInvalidArgument))
		return
	}

	accessKeyRecord, errCode := handlers.authSvc.VerifySignature(ctx, r)
	if errCode != ErrCodeNone {
		WriteErrorResponse(w, r, errCode)
		return
	}

	if err := s3utils.CheckValidBucketNameStrict(req.Bucket); err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrInvalidBucketName))
		return
	}

	if !checkPermissionType(req.ACL) {
		req.ACL = policy.Private
	}

	err = handlers.bucketSvc.CreateBucket(ctx, req.Bucket, req.Region, accessKeyRecord.Key, req.ACL)
	if err != nil {
		log.Errorf("PutBucketHandler create bucket error:%v", err)
		WriteErrorResponse(w, r, ToApiError(ctx, ErrCreateBucket))
		return
	}

	// Make sure to add Location information here only for bucket
	if cp := pathClean(r.URL.Path); cp != "" {
		w.Header().Set(consts.Location, cp) // Clean any trailing slashes.
	}

	WriteSuccessResponseHeadersOnly(w, r)

	return
}
