package sign

import (
	"github.com/bittorrent/go-btfs/s3/api/responses"
	"net/http"
)

type Service interface {
	SetSecretGetter(f func(key string) (secret string, exists, enable bool, err error))
	VerifyRequestSignature(r *http.Request) (ack string, rerr *responses.Error)
}
