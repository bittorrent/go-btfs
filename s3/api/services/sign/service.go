package sign

import (
	"github.com/bittorrent/go-btfs/s3/api/responses"
	"net/http"
	"sync"
)

var _ Service = (*service)(nil)

type service struct {
	getSecret func(key string) (secret string, exists, enable bool, err error)
	once      sync.Once
}

func NewService(options ...Option) Service {
	svc := &service{
		getSecret: func(key string) (secret string, exists, enable bool, err error) {
			return
		},
		once: sync.Once{},
	}
	for _, option := range options {
		option(svc)
	}
	return svc
}

func (s *service) SetSecretGetter(f func(key string) (secret string, exists, enable bool, err error)) {
	s.once.Do(func() {
		s.getSecret = f
	})
}

func (s *service) VerifyRequestSignature(r *http.Request) (ack string, rerr *responses.Error) {
	switch GetRequestAuthType(r) {
	case AuthTypeUnknown:
		return
	case AuthTypeSigned:
		ack, rerr = s.reqSignatureV4Verify(r, "")
		return
	case AuthTypeStreamingSigned:
		ack, rerr = s.setReqBodySignV4ChunkedReader(r, "")
		return
	default:
		rerr = responses.ErrSignatureVersionNotSupported
		return
	}
}
