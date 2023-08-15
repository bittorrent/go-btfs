package handlers

import (
	"context"
	"github.com/yann-y/fds/internal/lock"
	"github.com/yann-y/fds/internal/store"
	"github.com/yann-y/fds/internal/utils/hash"
	"github.com/yann-y/fds/pkg/s3utils"
	"golang.org/x/xerrors"
	"net/url"
)

// NotImplemented If a feature is not implemented
type NotImplemented struct {
	Message string
}

// ContextCanceled returns whether a context is canceled.
func ContextCanceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func ToApiError(ctx context.Context, err error) ErrorCode {
	if ContextCanceled(ctx) {
		if ctx.Err() == context.Canceled {
			return ErrClientDisconnected
		}
	}
	errCode := ErrInternalError
	switch err.(type) {
	case lock.OperationTimedOut:
		errCode = ErrOperationTimedOut
	case hash.SHA256Mismatch:
		errCode = ErrContentSHA256Mismatch
	case hash.BadDigest:
		errCode = ErrBadDigest
	case store.BucketNotFound:
		errCode = ErrNoSuchBucket
	case store.BucketPolicyNotFound:
		errCode = ErrNoSuchBucketPolicy
	case store.BucketTaggingNotFound:
		errCode = ErrBucketTaggingNotFound
	case s3utils.BucketNameInvalid:
		errCode = ErrInvalidBucketName
	case s3utils.ObjectNameInvalid:
		errCode = ErrInvalidObjectName
	case s3utils.ObjectNameTooLong:
		errCode = ErrKeyTooLongError
	case s3utils.ObjectNamePrefixAsSlash:
		errCode = ErrInvalidObjectNamePrefixSlash
	case s3utils.InvalidUploadIDKeyCombination:
		errCode = ErrNotImplemented
	case s3utils.InvalidMarkerPrefixCombination:
		errCode = ErrNotImplemented
	case s3utils.MalformedUploadID:
		errCode = ErrNoSuchUpload
	case s3utils.InvalidUploadID:
		errCode = ErrNoSuchUpload
	case s3utils.InvalidPart:
		errCode = ErrInvalidPart
	case s3utils.PartTooSmall:
		errCode = ErrEntityTooSmall
	case s3utils.PartTooBig:
		errCode = ErrEntityTooLarge
	case url.EscapeError:
		errCode = ErrInvalidObjectName
	default:
		if xerrors.Is(err, store.ErrObjectNotFound) {
			errCode = ErrNoSuchKey
		} else if xerrors.Is(err, store.ErrBucketNotEmpty) {
			errCode = ErrBucketNotEmpty
		}
	}
	return errCode
}
