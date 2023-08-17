// Package handlers is an implementation of s3.Handlerser
package handlers

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
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

func (h *Handlers) PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &PutBucketRequest{}
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

	err = h.bucketSvc.CheckACL(accessKeyRecord, req.Bucket, s3action.CreateBucketAction)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrNoSuchUserPolicy))
		return
	}

	if err := s3utils.CheckValidBucketNameStrict(req.Bucket); err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrInvalidBucketName))
		return
	}

	if !checkPermissionType(req.ACL) {
		WriteErrorResponse(w, r, ToApiError(ctx, ErrNotImplemented))
	}

	if ok := h.bucketSvc.HasBucket(r.Context(), req.Bucket); !ok {
		WriteErrorResponseHeadersOnly(w, r, ToApiError(ctx, ErrBucketNotFound))
		return
	}

	err = h.bucketSvc.CreateBucket(ctx, req.Bucket, req.Region, accessKeyRecord.Key, req.ACL)
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

func (h *Handlers) HeadBucketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &PutBucketRequest{}
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

	WriteSuccessResponseHeadersOnly(w, r)
}

func (h *Handlers) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &PutBucketRequest{}
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
	WriteSuccessNoContent(w)
}

func (h *Handlers) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &PutBucketRequest{}
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
	var buckets []*s3.Bucket
	for _, b := range bucketMetas {
		buckets = append(buckets, &s3.Bucket{
			Name:         aws.String(b.Name),
			CreationDate: aws.Time(b.Created),
		})
	}

	resp := ListAllMyBucketsResult{
		Owner: &s3.Owner{
			ID:          aws.String(consts.DefaultOwnerID),
			DisplayName: aws.String(consts.DisplayName),
		},
		Buckets: buckets,
	}

	WriteSuccessResponseXML(w, r, resp)
}

func (h *Handlers) GetBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &PutBucketRequest{}
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
	bucketMeta, err := h.bucketSvc.GetBucketMeta(ctx, req.Bucket)
	if err != nil {
		WriteErrorResponse(w, r, ToApiError(ctx, err))
		return
	}
	// 校验桶ACL类型，公共读(PublicRead)，公共读写(PublicReadWrite)，私有(Private)
	acl := bucketMeta.Acl
	if acl == "" {
		acl = "private"
	}

	resp := AccessControlPolicy{}
	id := accessKeyRecord.Key
	if resp.Owner.DisplayName == "" {
		resp.Owner.DisplayName = accessKeyRecord.Key
		resp.Owner.ID = id
	}
	resp.AccessControlList.Grant = append(resp.AccessControlList.Grant, Grant{
		Grantee: Grantee{
			ID:          id,
			DisplayName: accessKeyRecord.Key,
			Type:        "CanonicalUser",
			XMLXSI:      "CanonicalUser",
			XMLNS:       "http://www.w3.org/2001/XMLSchema-instance"},
		Permission: Permission(acl), //todo change
	})
	WriteSuccessResponseXML(w, r, resp)
}

func (h *Handlers) PutBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &PutBucketRequest{}
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
	WriteSuccessNoContent(w)
}
