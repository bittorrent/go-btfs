package cors

import (
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3d/consts"
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

var _ handlers.CorsService = (*Service)(nil)

type Service struct {
	allowOrigins []string
	allowMethods []string
	allowHeaders []string
}

func NewService(options ...Option) (svc *Service) {
	svc = &Service{
		allowOrigins: defaultAllowOrigins,
		allowMethods: defaultAllowMethods,
		allowHeaders: defaultAllowHeaders,
	}
	for _, option := range options {
		option(svc)
	}
	return
}

func (svc *Service) GetAllowOrigins() []string {
	return svc.allowOrigins
}

func (svc *Service) GetAllowMethods() []string {
	return svc.allowMethods
}

func (svc *Service) GetAllowHeaders() []string {
	return svc.allowHeaders
}
