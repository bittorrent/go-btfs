package cctx

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/services"
	"net/http"
)

func SetAccessKey(r *http.Request, ack *services.AccessKey) {
	ctx := context.WithValue(r.Context(), keyOfAccessKey, ack)
	r.WithContext(ctx)
}

func GetAccessKey(r *http.Request) (ack *services.AccessKey) {
	v := r.Context().Value(keyOfAccessKey)
	if v == nil {
		return
	}
	ack, _ = v.(*services.AccessKey)
	return
}
