package handlers

import (
	"github.com/bittorrent/go-btfs/s3/handlers/responses"
	"github.com/bittorrent/go-btfs/s3/services"
)

var toResponseErr = map[error]*responses.Error{
	services.ErrBucketNotFound: responses.ErrNoSuchBucket,
}

func ToResponseErr(err error) (rerr *responses.Error) {
	rerr, ok := toResponseErr[err]
	if !ok {
		rerr = responses.ErrInternalError
	}
	return
}
