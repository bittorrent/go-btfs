package auth

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/handlers/responses"
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/services"
	"net/http"
)

var _ services.AuthService = (*Service)(nil)

type Service struct {
	providers    providers.Providerser
	accessKeySvc services.AccessKeyService
}

func NewService(providers providers.Providerser, accessKeySvc services.AccessKeyService, options ...Option) (svc *Service) {
	svc = &Service{
		providers:    providers,
		accessKeySvc: accessKeySvc,
	}
	for _, option := range options {
		option(svc)
	}
	return
}

func (s *Service) VerifySignature(ctx context.Context, r *http.Request) (accessKeyRecord *services.AccessKey, err *responses.Error) {
	return s.CheckRequestAuthTypeCredential(ctx, r)
}
