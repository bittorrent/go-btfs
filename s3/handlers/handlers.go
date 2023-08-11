// Package handlers is an implementation of s3.Handlerser
package handlers

import (
	"github.com/bittorrent/go-btfs/s3"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/services"
	"github.com/rs/cors"
	"net/http"
)

var (
	defaultCorsAllowOrigins = []string{"*"}
	defaultCorsAllowHeaders = []string{
		consts.Date,
		consts.ETag,
		consts.ServerInfo,
		consts.Connection,
		consts.AcceptRanges,
		consts.ContentRange,
		consts.ContentEncoding,
		consts.ContentLength,
		consts.ContentType,
		consts.ContentDisposition,
		consts.LastModified,
		consts.ContentLanguage,
		consts.CacheControl,
		consts.RetryAfter,
		consts.AmzBucketRegion,
		consts.Expires,
		consts.Authorization,
		consts.Action,
		consts.Range,
		"X-Amz*",
		"x-amz*",
		"*",
	}
	defaultCorsAllowMethods = []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodHead,
		http.MethodPost,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodPatch,
	}
)

var _ s3.Handlerser = (*Handlers)(nil)

type Handlers struct {
	corsAllowOrigins []string
	corsAllowHeaders []string
	corsAllowMethods []string
	authSvc          services.SignService
	bucketSvc        services.BucketService
	objectSvc        services.ObjectService
	multipartSvc     services.MultipartService
}

func NewHandlers(
	authSvc services.SignService, bucketSvc services.BucketService,
	objectSvc services.ObjectService, multipartSvc services.MultipartService,
	options ...Option,
) (handlers *Handlers) {
	handlers = &Handlers{
		corsAllowOrigins: defaultCorsAllowOrigins,
		corsAllowHeaders: defaultCorsAllowHeaders,
		corsAllowMethods: defaultCorsAllowMethods,
		authSvc:          authSvc,
		bucketSvc:        bucketSvc,
		objectSvc:        objectSvc,
		multipartSvc:     multipartSvc,
	}
	for _, option := range options {
		option(handlers)
	}
	return
}

func (s *Handlers) Cors(handler http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   s.corsAllowOrigins,
		AllowedMethods:   s.corsAllowMethods,
		AllowedHeaders:   s.corsAllowHeaders,
		ExposedHeaders:   s.corsAllowHeaders,
		AllowCredentials: true,
	}).Handler(handler)
}

func (s *Handlers) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	return
}
