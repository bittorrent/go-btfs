package routers

import (
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

type Routers struct {
	handlers handlers.Handlerser
}

func NewRouters(handlers handlers.Handlerser, options ...Option) (routers *Routers) {
	routers = &Routers{
		handlers: handlers,
	}
	for _, option := range options {
		option(routers)
	}
	return
}

func (routers *Routers) Register() http.Handler {
	hs := routers.handlers

	root := mux.NewRouter()
	root.Use(
		hs.Cors,
		hs.Log,
		hs.Sign,
	)

	bucket := root.PathPrefix("/{bucket}").Subrouter()

	// multipart object...
	// CreateMultipart
	//bucket.Methods(http.MethodPost).Path("/{object:.+}").HandlerFunc(hs.CreateMultipartUploadHandler).Queries("uploads", "")
	//// UploadPart
	//bucket.Methods(http.MethodPut).Path("/{object:.+}").HandlerFunc(hs.UploadPartHandler).Queries("partNumber", "{partNumber:[0-9]+}", "uploadId", "{uploadId:.*}")
	//// CompleteMultipartUpload
	//bucket.Methods(http.MethodPost).Path("/{object:.+}").HandlerFunc(hs.CompleteMultipartUploadHandler).Queries("uploadId", "{uploadId:.*}")
	//// AbortMultipart
	//bucket.Methods(http.MethodDelete).Path("/{object:.+}").HandlerFunc(hs.AbortMultipartUploadHandler).Queries("uploadId", "{uploadId:.*}")

	////object...
	//// ListObjectsV2
	//bucket.Methods(http.MethodGet).HandlerFunc(hs.ListObjectsV2Handler).Queries("list-type", "2")
	// HeadObject
	bucket.Methods(http.MethodHead).Path("/{object:.+}").HandlerFunc(hs.HeadObjectHandler)
	// PutObject
	bucket.Methods(http.MethodPut).Path("/{object:.+}").HandlerFunc(hs.PutObjectHandler)
	// CopyObject
	bucket.Methods(http.MethodPut).Path("/{object:.+}").HeadersRegexp("X-Amz-Copy-Source", ".*?(\\/|%2F).*?").HandlerFunc(hs.CopyObjectHandler)
	// DeleteObject
	bucket.Methods(http.MethodDelete).Path("/{object:.+}").HandlerFunc(hs.DeleteObjectHandler)
	////todo DeleteObjects new ?
	//bucket.Methods(http.MethodDelete).Path("/{object:.+}").HandlerFunc(hs.DeleteObjectHandler)
	// GetObject
	bucket.Methods(http.MethodGet).Path("/{object:.+}").HandlerFunc(hs.GetObjectHandler)
	// GetObjectACL
	bucket.Methods(http.MethodGet).Path("/{object:.+}").HandlerFunc(hs.GetObjectACLHandler).Queries("acl", "")
	// ListObjects
	bucket.Methods(http.MethodGet).HandlerFunc(hs.ListObjectsHandler)

	//bucket...
	// GetBucketACL
	bucket.Methods(http.MethodGet).HandlerFunc(hs.GetBucketAclHandler).Queries("acl", "")
	// PutBucketACL
	bucket.Methods(http.MethodPut).HandlerFunc(hs.PutBucketAclHandler).Queries("acl", "")
	// CreateBucket
	bucket.Methods(http.MethodPut).HandlerFunc(hs.CreateBucketHandler)
	// HeadBucket
	bucket.Methods(http.MethodHead).HandlerFunc(hs.HeadBucketHandler)
	// DeleteBucket
	bucket.Methods(http.MethodDelete).HandlerFunc(hs.DeleteBucketHandler)
	// ListBuckets
	root.Methods(http.MethodGet).Path("/").HandlerFunc(hs.ListBucketsHandler)

	return root
}
