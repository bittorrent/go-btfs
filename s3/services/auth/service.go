package auth

import (
	"context"
	"net/http"

	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/services"
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

func (s *Service) VerifySignature(ctx context.Context, r *http.Request) (accessKeyRecord *handlers.AccessKeyRecord, err handlers.ErrorCode) {
	return s.CheckRequestAuthTypeCredential(ctx, r)
}
