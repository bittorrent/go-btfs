// Package handlers is an implementation of s3.Handlerser
package handlers

import (
	"github.com/bittorrent/go-btfs/s3/common/consts"
	"github.com/bittorrent/go-btfs/s3/server"
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

var _ server.Handlerser = (*Handlers)(nil)

type Handlers struct {
	corsAllowOrigins []string
	corsAllowHeaders []string
	corsAllowMethods []string
	signSvc          SignService
	bucketSvc        BucketService
	objectSvc        ObjectService
	multipartSvc     MultipartService
}

func NewHandlers(
	signSvc SignService, bucketSvc BucketService,
	objectSvc ObjectService, multipartSvc MultipartService,
	options ...Option,
) (handlers *Handlers) {
	handlers = &Handlers{
		corsAllowOrigins: defaultCorsAllowOrigins,
		corsAllowHeaders: defaultCorsAllowHeaders,
		corsAllowMethods: defaultCorsAllowMethods,
		signSvc:          signSvc,
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

func (s *Handlers) Sign(handler http.Handler) http.Handler {
	return nil
}

func (s *Handlers) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	return
}
