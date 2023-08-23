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

	root.Use(
		routers.handlers.Cors,
		routers.handlers.Auth,
	)

	bucket := root.PathPrefix("/{bucket}").Subrouter()
	bucket.Methods(http.MethodGet).HandlerFunc(routers.handlers.GetBucketAclHandler).Queries("acl", "")
	bucket.Methods(http.MethodPut).HandlerFunc(routers.handlers.PutBucketAclHandler).Queries("acl", "")

	bucket.Methods(http.MethodPut).HandlerFunc(routers.handlers.PutBucketHandler)
	bucket.Methods(http.MethodHead).HandlerFunc(routers.handlers.HeadBucketHandler)
	bucket.Methods(http.MethodDelete).HandlerFunc(routers.handlers.DeleteBucketHandler)

	root.Methods(http.MethodGet).Path("/").HandlerFunc(routers.handlers.ListBucketsHandler)

	//object
	//bucket.Methods(http.MethodPut).Path("/{object:.+}").HandlerFunc(routers.handlers.PutObjectHandler)

	return root
}
