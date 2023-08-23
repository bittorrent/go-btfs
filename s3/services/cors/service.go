package cors

import (
	"github.com/bittorrent/go-btfs/s3/consts"
	"net/http"
)

var (
	defaultAllowOrigins = []string{"*"}
	defaultAllowMethods = []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodHead,
		http.MethodPost,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodPatch,
	}
	defaultAllowHeaders = []string{
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
)

var _ Service = (*service)(nil)

type service struct {
	allowOrigins []string
	allowMethods []string
	allowHeaders []string
}

func NewService(options ...Option) Service {
	svc := &service{
		allowOrigins: defaultAllowOrigins,
		allowMethods: defaultAllowMethods,
		allowHeaders: defaultAllowHeaders,
	}
	for _, option := range options {
		option(svc)
	}
	return svc
}

func (svc *service) GetAllowOrigins() []string {
	return svc.allowOrigins
}

func (svc *service) GetAllowMethods() []string {
	return svc.allowMethods
}

func (svc *service) GetAllowHeaders() []string {
	return svc.allowHeaders
}
