package handlers

import (
	"net/http"
)

type Handlerser interface {
	// middlewares
	Cors(handler http.Handler) http.Handler
	Sign(handler http.Handler) http.Handler
	Log(handler http.Handler) http.Handler

	// bucket
	PutBucketHandler(w http.ResponseWriter, r *http.Request)
	HeadBucketHandler(w http.ResponseWriter, r *http.Request)
	DeleteBucketHandler(w http.ResponseWriter, r *http.Request)
	ListBucketsHandler(w http.ResponseWriter, r *http.Request)
	GetBucketAclHandler(w http.ResponseWriter, r *http.Request)
	PutBucketAclHandler(w http.ResponseWriter, r *http.Request)

	// object
	PutObjectHandler(w http.ResponseWriter, r *http.Request)
	HeadObjectHandler(w http.ResponseWriter, r *http.Request)
	CopyObjectHandler(w http.ResponseWriter, r *http.Request)
	DeleteObjectHandler(w http.ResponseWriter, r *http.Request)
	GetObjectHandler(w http.ResponseWriter, r *http.Request)
	GetObjectACLHandler(w http.ResponseWriter, r *http.Request)
	ListObjectsV1Handler(w http.ResponseWriter, r *http.Request)
	ListObjectsV2Handler(w http.ResponseWriter, r *http.Request)

	// multipart
	CreateMultipartUploadHandler(w http.ResponseWriter, r *http.Request)
	UploadPartHandler(w http.ResponseWriter, r *http.Request)
	AbortMultipartUploadHandler(w http.ResponseWriter, r *http.Request)
	CompleteMultipartUploadHandler(w http.ResponseWriter, r *http.Request)
}
