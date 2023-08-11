package sign

import (
	"github.com/bittorrent/go-btfs/s3/providers"
	"net/http"
)

type Sign struct {
	providers providers.Providerser
}

func NewSign(providers providers.Providerser, options ...Option) (sign *Sign) {
	sign = &Sign{
		providers: providers,
	}

	for _, option := range options {
		option(sign)
	}
	return
}

func (s *Sign) Verify(r *http.Request) (err error) {
	return
}
