// Package handlers is an implementation of Handlerser
package handlers

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"github.com/bittorrent/go-btfs/s3/services/sign"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"net/http"
	"net/url"
	"runtime"
)

var _ Handlerser = (*Handlers)(nil)

type Handlers struct {
	headers http.Header
	acksvc  accesskey.Service
	sigsvc  sign.Service
	objsvc  object.Service
}

func NewHandlers(acksvc accesskey.Service, sigsvc sign.Service, objsvc object.Service, options ...Option) (handlers *Handlers) {
	handlers = &Handlers{
		headers: defaultHeaders,
		acksvc:  acksvc,
		sigsvc:  sigsvc,
		objsvc:  objsvc,
	}
	for _, option := range options {
		option(handlers)
	}
	return
}

func (h *Handlers) name() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func (h *Handlers) respErr(err error) (rerr *responses.Error) {
	switch err {
	case object.ErrBucketNotFound:
		rerr = responses.ErrNoSuchBucket
	case object.ErrBucketeNotEmpty:
		rerr = responses.ErrBucketNotEmpty
	case object.ErrObjectNotFound:
		rerr = responses.ErrNoSuchKey
	case object.ErrUploadNotFound:
		rerr = responses.ErrNoSuchUpload
	case object.ErrBucketAlreadyExists:
		rerr = responses.ErrBucketAlreadyExists
	case object.ErrNotAllowed:
		rerr = responses.ErrAccessDenied
	case context.Canceled:
		rerr = responses.ErrClientDisconnected
	case context.DeadlineExceeded:
		rerr = responses.ErrOperationTimedOut
	default:
		switch err.(type) {
		case hash.SHA256Mismatch:
			rerr = responses.ErrContentSHA256Mismatch
		case hash.BadDigest:
			rerr = responses.ErrBadDigest
		case s3utils.BucketNameInvalid:
			rerr = responses.ErrInvalidBucketName
		case s3utils.ObjectNameInvalid:
			rerr = responses.ErrInvalidObjectName
		case s3utils.ObjectNameTooLong:
			rerr = responses.ErrKeyTooLongError
		case s3utils.ObjectNamePrefixAsSlash:
			rerr = responses.ErrInvalidObjectNamePrefixSlash
		case s3utils.InvalidUploadIDKeyCombination:
			rerr = responses.ErrNotImplemented
		case s3utils.InvalidMarkerPrefixCombination:
			rerr = responses.ErrNotImplemented
		case s3utils.MalformedUploadID:
			rerr = responses.ErrNoSuchUpload
		case s3utils.InvalidUploadID:
			rerr = responses.ErrNoSuchUpload
		case s3utils.InvalidPart:
			rerr = responses.ErrInvalidPart
		case s3utils.PartTooSmall:
			rerr = responses.ErrEntityTooSmall
		case s3utils.PartTooBig:
			rerr = responses.ErrEntityTooLarge
		case url.EscapeError:
			rerr = responses.ErrInvalidObjectName
		default:
			rerr = responses.ErrInternalError
		}
	}
	return
}
