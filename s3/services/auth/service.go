package auth

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/apierrors"
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/policy"
	"github.com/bittorrent/go-btfs/s3/services"
	"net/http"
)

var _ handlers.AuthService = (*Service)(nil)

type Service struct {
	providers    services.Providerser
	accessKeySvc handlers.AccessKeyService
}

func NewService(providers services.Providerser, accessKeySvc handlers.AccessKeyService, options ...Option) (svc *Service) {
	svc = &Service{
		providers:    providers,
		accessKeySvc: accessKeySvc,
	}
	for _, option := range options {
		option(svc)
	}
	return
}

func (s *Service) VerifySignature(ctx context.Context, r *http.Request) (accessKeyRecord *handlers.AccessKeyRecord, err apierrors.ErrorCode) {
	s.CheckRequestAuthTypeCredential(ctx, r)
	return
}

func (svc *Service) CheckACL(accessKeyRecord *handlers.AccessKeyRecord, bucketMeta *handlers.BucketMeta, action action.Action) (err error) {
	////todo 是否需要判断原始的
	//if bucketName == "" {
	//	return cred, handlers.ErrBucketNotFound
	//}

	//todo 注意：如果action是CreateBucketAction，HasBucket(ctx, bucketName)进行判断

	if policy.IsAllowed(bucketMeta.Owner == accessKeyRecord.Key, bucketMeta.Acl, action) == false {
		return cred, apierrors.ErrAccessDenied
	}
	return
}
