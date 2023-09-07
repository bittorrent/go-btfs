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
	root := mux.NewRouter()

	hs := routers.handlers


	root.Use(
		hs.Cors,
		hs.Log,
		hs.Sign,
	)

	bucket := root.PathPrefix("/{Bucket}").Subrouter()

	// multipart object...
	// CreateMultipart
	//bucket.Methods(http.MethodPost).Path("/{object:.+}").HandlerFunc(hs.CreateMultipartUploadHandler).Queries("uploads", "")
	//// UploadPart
	//bucket.Methods(http.MethodPut).Path("/{object:.+}").HandlerFunc(hs.UploadPartHandler).Queries("partNumber", "{partNumber:[0-9]+}", "uploadId", "{uploadId:.*}")
	//// CompleteMultipartUpload
	//bucket.Methods(http.MethodPost).Path("/{object:.+}").HandlerFunc(hs.CompleteMultipartUploadHandler).Queries("uploadId", "{uploadId:.*}")
	//// AbortMultipart
	//bucket.Methods(http.MethodDelete).Path("/{object:.+}").HandlerFunc(hs.AbortMultipartUploadHandler).Queries("uploadId", "{uploadId:.*}")

	// HeadObject
	bucket.Methods(http.MethodHead).Path("/{Object:.+}").HandlerFunc(hs.HeadObjectHandler)

	// CopyObject
	bucket.Methods(http.MethodPut).Path("/{Object:.+}").HeadersRegexp("X-Amz-Copy-Source", ".*?(\\/|%2F).*?").HandlerFunc(hs.CopyObjectHandler)

	// PutObject
	bucket.Methods(http.MethodPut).Path("/{Object:.+}").HandlerFunc(hs.PutObjectHandler)

	// DeleteObject
	bucket.Methods(http.MethodDelete).Path("/{Object:.+}").HandlerFunc(hs.DeleteObjectHandler)

	// GetObjectACL
	bucket.Methods(http.MethodGet).Path("/{Object:.+}").HandlerFunc(hs.GetObjectACLHandler).Queries("acl", "")

	// GetObject
	bucket.Methods(http.MethodGet).Path("/{Object:.+}").HandlerFunc(hs.GetObjectHandler)

	// GetBucketACL
	bucket.Methods(http.MethodGet).HandlerFunc(hs.GetBucketAclHandler).Queries("acl", "")

	// ListObjectsV2
	bucket.Methods(http.MethodGet).HandlerFunc(hs.ListObjectsV2Handler).Queries("list-type", "2")

	// ListObjects
	bucket.Methods(http.MethodGet).HandlerFunc(hs.ListObjectsHandler)

	// PutBucketACL
	bucket.Methods(http.MethodPut).HandlerFunc(hs.PutBucketAclHandler).Queries("acl", "")

	// CreateBucket
	bucket.Methods(http.MethodPut).HandlerFunc(hs.CreateBucketHandler)

	// HeadBucket
	bucket.Methods(http.MethodHead).HandlerFunc(hs.HeadBucketHandler)

	// DeleteObjects
	bucket.Methods(http.MethodPost).HandlerFunc(hs.DeleteObjectsHandler).Queries("delete", "")

	// DeleteBucket
	bucket.Methods(http.MethodDelete).HandlerFunc(hs.DeleteBucketHandler)

	// ListBuckets
	root.Methods(http.MethodGet).Path("/").HandlerFunc(hs.ListBucketsHandler)

	// Options
	root.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return root
}
