package requests

import "errors"

var (
	ErrBucketNameInvalid              = errors.New("the bucket name is invalid")
	ErrObjectNameInvalid              = errors.New("the object name is invalid")
	ErrObjectNameTooLong              = errors.New("the object name cannot be longer than 1024 characters")
	ErrObjectNamePrefixSlash          = errors.New("the object name cannot start with slash")
	ErrRegionUnsupported              = errors.New("the location is not supported by this server")
	ErrACLUnsupported                 = errors.New("the ACL is not supported by this server")
	ErrInvalidContentMd5              = errors.New("the content md5 is invalid")
	ErrInvalidChecksumSha256          = errors.New("the checksum-sha256 is invalid")
	ErrContentLengthMissing           = errors.New("the content-length is missing")
	ErrContentLengthTooSmall          = errors.New("the content-length is too small")
	ErrContentLengthTooLarge          = errors.New("the content-length is too large")
	ErrCopySrcInvalid                 = errors.New("the copy-source is invalid")
	ErrCopyDestInvalid                = errors.New("the copy-destination is invalid")
	ErrMaxKeysInvalid                 = errors.New("the max-keys is invalid")
	ErrEncodingTypeInvalid            = errors.New("the encoding-type is invalid")
	ErrMarkerPrefixCombinationInvalid = errors.New("the marker-prefix combination is invalid")
)
