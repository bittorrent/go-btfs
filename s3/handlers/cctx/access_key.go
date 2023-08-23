package cctx

import (
	"github.com/bittorrent/go-btfs/s3/services"
	"net/http"
)

func SetAccessKey(r *http.Request, ack *services.AccessKey) {
	set(r, keyOfAccessKey, ack)
	return
}

func GetAccessKey(r *http.Request) (ack *services.AccessKey) {
	v := get(r, keyOfAccessKey)
	ack, _ = v.(*services.AccessKey)
	return
}
