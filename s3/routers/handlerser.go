package routers

import (
	"net/http"
)

type Handlerser interface {
	Cors(handler http.Handler) http.Handler
	Sign(handler http.Handler) http.Handler

	PutBucketHandler(w http.ResponseWriter, r *http.Request)
	HeadBucketHandler(w http.ResponseWriter, r *http.Request)
	DeleteBucketHandler(w http.ResponseWriter, r *http.Request)
	ListBucketsHandler(w http.ResponseWriter, r *http.Request)
	GetBucketAclHandler(w http.ResponseWriter, r *http.Request)
	PutBucketAclHandler(w http.ResponseWriter, r *http.Request)

	//PutObjectHandler(w http.ResponseWriter, r *http.Request)
}
