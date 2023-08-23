package cctx

import (
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"net/http"
)

func SetAccessKey(r *http.Request, ack *accesskey.AccessKey) {
	set(r, keyOfAccessKey, ack)
	return
}

func GetAccessKey(r *http.Request) (ack *accesskey.AccessKey) {
	v := get(r, keyOfAccessKey)
	ack, _ = v.(*accesskey.AccessKey)
	return
}
