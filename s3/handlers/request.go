package handlers

import (
	"io"
	"net/http"
)

type RequestBinder interface {
	Bind(r *http.Request) (err error)
}

type PutObjectRequest struct {
	Bucket string
	Object string
	Body   io.Reader
}

func (req *PutObjectRequest) Bind(r *http.Request) (err error) {
	return
}
