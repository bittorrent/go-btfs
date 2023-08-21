package handlers

import (
	"context"
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
			return ErrCodeClientDisconnected
		}
	}
	errCode := ErrCodeInternalError
	switch err {
	case ErrInvalidArgument:
		errCode = ErrCodeInvalidRequestBody //实际是request请求信息， header or query uri 信息。
	case ErrInvalidBucketName:
		errCode = ErrCodeInvalidBucketName
	case ErrBucketNotFound:
		errCode = ErrCodeNoSuchBucket
	case ErrBucketAccessDenied:
		errCode = ErrCodeAccessDenied
	case ErrSetBucketEmptyFailed:
	case ErrCreateBucket:
		errCode = ErrCodeInternalError
	case ErrNotImplemented:
		errCode = ErrCodeNotImplemented
	case ErrBucketAlreadyExists:
		errCode = ErrCodeBucketAlreadyExists
		//case lock.OperationTimedOut:
		//	errCode = ErrCodeOperationTimedOut
		//case hash.SHA256Mismatch:
		//	errCode = ErrCodeContentSHA256Mismatch
		//case hash.BadDigest:
		//	errCode = ErrCodeBadDigest
		//case store.BucketPolicyNotFound:
		//	errCode = ErrCodeNoSuchBucketPolicy
		//case store.BucketTaggingNotFound:
		//	errCode = ErrBucketTaggingNotFound
		//case s3utils.BucketNameInvalid:
		//	errCode = ErrCodeInvalidBucketName
		//case s3utils.ObjectNameInvalid:
		//	errCode = ErrCodeInvalidObjectName
		//case s3utils.ObjectNameTooLong:
		//	errCode = ErrCodeKeyTooLongError
		//case s3utils.ObjectNamePrefixAsSlash:
		//	errCode = ErrCodeInvalidObjectNamePrefixSlash
		//case s3utils.InvalidUploadIDKeyCombination:
		//	errCode = ErrCodeNotImplemented
		//case s3utils.InvalidMarkerPrefixCombination:
		//	errCode = ErrCodeNotImplemented
		//case s3utils.MalformedUploadID:
		//	errCode = ErrCodeNoSuchUpload
		//case s3utils.InvalidUploadID:
		//	errCode = ErrCodeNoSuchUpload
		//case s3utils.InvalidPart:
		//	errCode = ErrCodeInvalidPart
		//case s3utils.PartTooSmall:
		//	errCode = ErrCodeEntityTooSmall
		//case s3utils.PartTooBig:
		//	errCode = ErrCodeEntityTooLarge
		//case url.EscapeError:
		//	errCode = ErrCodeInvalidObjectName
		//default:
		//	if xerrors.Is(err, store.ErrObjectNotFound) {
		//		errCode = ErrCodeNoSuchKey
		//	} else if xerrors.Is(err, store.ErrBucketNotEmpty) {
		//		errCode = ErrCodeBucketNotEmpty
		//	}
	}
	return errCode
}
