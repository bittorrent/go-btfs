package handlers

import (
	"net/http"
)

type Handlerser interface {
	// middlewares
	Cors(handler http.Handler) http.Handler
	Auth(handler http.Handler) http.Handler
	Log(handler http.Handler) http.Handler

	// handlers
	PutBucketHandler(w http.ResponseWriter, r *http.Request)
	HeadBucketHandler(w http.ResponseWriter, r *http.Request)
	DeleteBucketHandler(w http.ResponseWriter, r *http.Request)
	ListBucketsHandler(w http.ResponseWriter, r *http.Request)
	GetBucketAclHandler(w http.ResponseWriter, r *http.Request)
	PutBucketAclHandler(w http.ResponseWriter, r *http.Request)

	//PutObjectHandler(w http.ResponseWriter, r *http.Request)
}
