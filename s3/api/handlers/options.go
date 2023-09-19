package handlers

import (
	"github.com/bittorrent/go-btfs/s3/consts"
	"net/http"
)

var defaultCorsMethods = []string{
	http.MethodGet,
	http.MethodPut,
	http.MethodHead,
	http.MethodPost,
	http.MethodDelete,
	http.MethodOptions,
	http.MethodPatch,
}

var defaultCorsHeaders = []string{
	consts.Date,
	consts.ETag,
	consts.ServerInfo,
	consts.Connection,
	consts.AcceptRanges,
	consts.ContentRange,
	consts.ContentEncoding,
	consts.ContentLength,
	consts.ContentType,
	consts.ContentMD5,
	consts.ContentDisposition,
	consts.LastModified,
	consts.ContentLanguage,
	consts.CacheControl,
	consts.Location,
	consts.RetryAfter,
	consts.AmzBucketRegion,
	consts.Expires,
	consts.Authorization,
	consts.Action,
	consts.XRequestWith,
	consts.Range,
	consts.UserAgent,
	consts.Cid,
	"Amz-*",
	"amz-*",
	"X-Amz*",
	"x-amz*",
	"*",
}

const defaultCorsMaxAge = "36000"

var defaultHeaders = map[string][]string{
	consts.AccessControlAllowOrigin:      {"*"},
	consts.AccessControlAllowMethods:     defaultCorsMethods,
	consts.AccessControlAllowHeaders:     defaultCorsHeaders,
	consts.AccessControlExposeHeaders:    defaultCorsHeaders,
	consts.AccessControlAllowCredentials: {"true"},
	consts.AccessControlMaxAge:           {defaultCorsMaxAge},
}

type Option func(handlers *Handlers)

func WithHeaders(headers map[string][]string) Option {
	return func(handlers *Handlers) {
		if headers != nil {
			handlers.headers = headers
		}
	}
}
