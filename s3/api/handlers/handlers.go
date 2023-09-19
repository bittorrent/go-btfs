// Package handlers is an implementation of Handlerser
package handlers

import (
	"github.com/bittorrent/go-btfs/s3/api/requests"
	"github.com/bittorrent/go-btfs/s3/api/responses"
	"github.com/bittorrent/go-btfs/s3/api/services/accesskey"
	"github.com/bittorrent/go-btfs/s3/api/services/object"
	"github.com/bittorrent/go-btfs/s3/api/services/sign"
	"github.com/bittorrent/go-btfs/s3/hash"
	"net/http"
	"runtime"
	"strings"
)

var _ Handlerser = (*Handlers)(nil)

type Handlers struct {
	headers http.Header
	acksvc  accesskey.Service
	sigsvc  sign.Service
	objsvc  object.Service
}

func NewHandlers(
	acksvc accesskey.Service, sigsvc sign.Service, objsvc object.Service,
	options ...Option) (handlers *Handlers) {
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

// name returns name of the handler function
func (h *Handlers) name() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	ps := strings.Split(f.Name(), ".")
	if len(ps) > 0 {
		return ps[len(ps)-1]
	}
	return "UnknownHandler"
}

// toResponseErr convert internal error to response error
func (h *Handlers) toResponseErr(err error) (rerr *responses.Error) {
	switch err {
	// Errors from requests
	case requests.ErrBucketNameInvalid:
		rerr = responses.ErrInvalidBucketName
	case requests.ErrObjectNameInvalid:
		rerr = responses.ErrInvalidObjectName
	case requests.ErrObjectNameTooLong:
		rerr = responses.ErrKeyTooLongError
	case requests.ErrObjectNamePrefixSlash:
		rerr = responses.ErrInvalidObjectNamePrefixSlash
	case requests.ErrRegionUnsupported:
		rerr = responses.ErrInvalidRegion
	case requests.ErrACLUnsupported:
		rerr = responses.ErrMalformedACLError
	case requests.ErrContentMd5Invalid:
		rerr = responses.ErrInvalidDigest
	case requests.ErrChecksumSha256Invalid:
		rerr = responses.ErrContentSHA256Mismatch
	case requests.ErrContentLengthMissing:
		rerr = responses.ErrMissingContentLength
	case requests.ErrContentLengthTooSmall:
		rerr = responses.ErrEntityTooSmall
	case requests.ErrContentLengthTooLarge:
		rerr = responses.ErrEntityTooLarge
	case requests.ErrCopySrcInvalid:
		rerr = responses.ErrInvalidCopySource
	case requests.ErrCopyDestInvalid:
		rerr = responses.ErrInvalidCopyDest
	case requests.ErrDeletesCountInvalid:
		rerr = responses.ErrInvalidRequest
	case requests.ErrMaxKeysInvalid:
		rerr = responses.ErrInvalidMaxKeys
	case requests.ErrPrefixInvalid:
		rerr = responses.ErrInvalidRequest
	case requests.ErrMarkerInvalid:
		rerr = responses.ErrInvalidRequest
	case requests.ErrMarkerPrefixCombinationInvalid:
		rerr = responses.ErrInvalidRequest
	case requests.ErrContinuationTokenInvalid:
		rerr = responses.ErrIncorrectContinuationToken
	case requests.ErrStartAfterInvalid:
		rerr = responses.ErrInvalidRequest
	case requests.ErrPartNumberInvalid:
		rerr = responses.ErrInvalidPartNumber
	case requests.ErrPartsCountInvalid:
		rerr = responses.ErrInvalidRequest
	case requests.ErrPartInvalid:
		rerr = responses.ErrInvalidPart
	case requests.ErrPartOrderInvalid:
		rerr = responses.ErrInvalidPartOrder
	// Errors from Object service
	case object.ErrBucketNotFound:
		rerr = responses.ErrNoSuchBucket
	case object.ErrBucketNotEmpty:
		rerr = responses.ErrBucketNotEmpty
	case object.ErrObjectNotFound:
		rerr = responses.ErrNoSuchKey
	case object.ErrUploadNotFound:
		rerr = responses.ErrNoSuchUpload
	case object.ErrBucketAlreadyExists:
		rerr = responses.ErrBucketAlreadyExists
	case object.ErrNotAllowed:
		rerr = responses.ErrAccessDenied
	case object.ErrPartNotExists:
		rerr = responses.ErrInvalidPart
	case object.ErrPartETagNotMatch:
		rerr = responses.ErrInvalidPart
	case object.ErrPartTooSmall:
		rerr = responses.ErrEntityTooSmall
	case object.ErrCanceled:
		rerr = responses.ErrClientDisconnected
	case object.ErrTimout:
		rerr = responses.ErrOperationTimedOut
	// Others
	default:
		switch nerr := err.(type) {
		case requests.ErrFailedParseValue:
			rerr = responses.ErrInvalidRequest
		case requests.ErrFailedDecodeXML:
			rerr = responses.ErrMalformedXML
		case requests.ErrMissingRequiredParam:
			rerr = responses.ErrInvalidRequest
		case requests.ErrWithUnsupportedParam:
			rerr = responses.ErrNotImplemented
		case hash.SHA256Mismatch:
			rerr = responses.ErrContentSHA256Mismatch
		case hash.BadDigest:
			rerr = responses.ErrBadDigest
		case hash.ErrSizeMismatch:
			if nerr.Got < nerr.Want {
				rerr = responses.ErrIncompleteBody
			} else {
				rerr = responses.ErrMissingContentLength
			}
		default:
			rerr = responses.ErrInternalError
		}
	}
	return
}
