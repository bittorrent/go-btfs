package handlers

import (
	"net/http"
)

type Handlerser interface {
	// Middlewares

	Cors(handler http.Handler) http.Handler
	Sign(handler http.Handler) http.Handler
	Log(handler http.Handler) http.Handler

	// Bucket

	CreateBucketHandler(w http.ResponseWriter, r *http.Request)
	HeadBucketHandler(w http.ResponseWriter, r *http.Request)
	DeleteBucketHandler(w http.ResponseWriter, r *http.Request)
	ListBucketsHandler(w http.ResponseWriter, r *http.Request)
	PutBucketACLHandler(w http.ResponseWriter, r *http.Request)
	GetBucketACLHandler(w http.ResponseWriter, r *http.Request)

	// Object

	PutObjectHandler(w http.ResponseWriter, r *http.Request)
	CopyObjectHandler(w http.ResponseWriter, r *http.Request)
	HeadObjectHandler(w http.ResponseWriter, r *http.Request)
	GetObjectHandler(w http.ResponseWriter, r *http.Request)
	DeleteObjectHandler(w http.ResponseWriter, r *http.Request)
	DeleteObjectsHandler(w http.ResponseWriter, r *http.Request)
	ListObjectsHandler(w http.ResponseWriter, r *http.Request)
	ListObjectsV2Handler(w http.ResponseWriter, r *http.Request)
	GetObjectACLHandler(w http.ResponseWriter, r *http.Request)

	// Multipart

	CreateMultipartUploadHandler(w http.ResponseWriter, r *http.Request)
	UploadPartHandler(w http.ResponseWriter, r *http.Request)
	AbortMultipartUploadHandler(w http.ResponseWriter, r *http.Request)
	CompleteMultipartUploadHandler(w http.ResponseWriter, r *http.Request)
}
