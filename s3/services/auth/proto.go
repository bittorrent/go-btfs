package auth

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"net/http"
)

type Service interface {
	VerifySignature(ctx context.Context, r *http.Request) (ack *accesskey.AccessKey, err error)
}
