package auth

import (
	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/services"
	"net/http"
)

var _ handlers.AuthService = (*Service)(nil)

type Service struct {
	providers    services.Providerser
	accesskeySvc handlers.AccessKeyService
}

func NewService(providers services.Providerser, accesskeySvc handlers.AccessKeyService, options ...Option) (s *Service) {
	s = &Service{
		providers:    providers,
		accesskeySvc: accesskeySvc,
	}
	for _, option := range options {
		option(s)
	}
	return
}

func (s *Service) VerifySignature(r *http.Request) (accessKeyRecord *handlers.AccessKeyRecord, err error) {
	return
}

func (s *Service) CheckACL(accessKeyRecord *handlers.AccessKeyRecord, bucketMeta *handlers.BucketMeta, action action.Action) (err error) {
	return
}
