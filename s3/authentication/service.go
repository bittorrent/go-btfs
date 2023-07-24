package authentication

import "github.com/bittorrent/go-btfs/s3/adaptor"

type Service interface {
	VerifyRequest(req *adaptor.Request) *AuthErr
}
