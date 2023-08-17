package handlers

import "errors"

var (
	ErrSginVersionNotSupport = errors.New("sign version is not support")

	ErrInvalidArgument = errors.New("invalid argument")

	ErrInvalidBucketName    = errors.New("bucket name is invalid")
	ErrBucketNotFound       = errors.New("bucket is not found")
	ErrBucketAccessDenied   = errors.New("bucket access denied. ")
	ErrSetBucketEmptyFailed = errors.New("set bucket empty failed. ")
	ErrCreateBucket         = errors.New("create bucket failed")
)
