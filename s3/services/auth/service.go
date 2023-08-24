package auth

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"net/http"
)

var _ Service = (*service)(nil)

type service struct {
	getSecret func(key string) (secret string, disabled bool, err error)
}

func NewService(providers providers.Providerser, accessKeySvc accesskey.Service, options ...Option) Service {
	svc := &service{}
	for _, option := range options {
		option(svc)
	}
	return svc
}

func (s *service) VerifySignature(ctx context.Context, r *http.Request) (accessKeyRecord *accesskey.AccessKey, err error) {
	return s.CheckRequestAuthTypeCredential(ctx, r)
}
