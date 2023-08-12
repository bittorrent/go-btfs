// Package handlers is an implementation of s3.Handlerser
package handlers

import (
	"github.com/bittorrent/go-btfs/s3/routers"
	"github.com/rs/cors"
	"net/http"
)

var _ routers.Handlerser = (*Handlers)(nil)

type Handlers struct {
	corsSvc      CorsService
	authSvc      AuthService
	bucketSvc    BucketService
	objectSvc    ObjectService
	multipartSvc MultipartService
}

func NewHandlers(
	corsSvc CorsService,
	authSvc AuthService,
	bucketSvc BucketService,
	objectSvc ObjectService,
	multipartSvc MultipartService,
	options ...Option,
) (handlers *Handlers) {
	handlers = &Handlers{
		corsSvc:      corsSvc,
		authSvc:      authSvc,
		bucketSvc:    bucketSvc,
		objectSvc:    objectSvc,
		multipartSvc: multipartSvc,
	}
	for _, option := range options {
		option(handlers)
	}
	return
}

func (handlers *Handlers) Cors(handler http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   handlers.corsSvc.GetAllowOrigins(),
		AllowedMethods:   handlers.corsSvc.GetAllowMethods(),
		AllowedHeaders:   handlers.corsSvc.GetAllowHeaders(),
		ExposedHeaders:   handlers.corsSvc.GetAllowHeaders(),
		AllowCredentials: true,
	}).Handler(handler)
}

func (handlers *Handlers) Sign(handler http.Handler) http.Handler {
	return nil
}

func (handlers *Handlers) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	return
}
