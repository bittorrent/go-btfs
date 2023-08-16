package handlers

import "errors"

var (
	ErrSginVersionNotSupport = errors.New("sign version is not support")

	// bucket
	ErrBucketNotFound       = errors.New("bucket is not found")
	ErrBucketAccessDenied   = errors.New("bucket access denied. ")
	ErrSetBucketEmptyFailed = errors.New("set bucket empty failed. ")
)
