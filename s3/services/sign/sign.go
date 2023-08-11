package sign

import (
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/services"
	"net/http"
)

var _ handlers.SignService = (*Sign)(nil)

type Sign struct {
	providers    services.Providerser
	accesskeySvc handlers.AccessKeyService
}

func NewSign(providers services.Providerser, accesskeySvc handlers.AccessKeyService, options ...Option) (sign *Sign) {
	sign = &Sign{
		providers:    providers,
		accesskeySvc: accesskeySvc,
	}
	for _, option := range options {
		option(sign)
	}
	return
}

func (s *Sign) Verify(r *http.Request) (err error) {
	return
}
