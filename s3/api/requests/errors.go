package requests

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrBucketNameInvalid              = errors.New("the bucket name is invalid")
	ErrObjectNameInvalid              = errors.New("the object name is invalid")
	ErrObjectNameTooLong              = errors.New("the object name cannot be longer than 1024 characters")
	ErrObjectNamePrefixSlash          = errors.New("the object name cannot start with slash")
	ErrRegionUnsupported              = errors.New("the location is not supported by this server")
	ErrACLUnsupported                 = errors.New("the ACL is not supported by this server")
	ErrContentMd5Invalid              = errors.New("the content md5 is invalid")
	ErrChecksumSha256Invalid          = errors.New("the checksum-sha256 is invalid")
	ErrContentLengthMissing           = errors.New("the content-length is missing")
	ErrContentLengthTooSmall          = errors.New("the content-length is too small")
	ErrContentLengthTooLarge          = errors.New("the content-length is too large")
	ErrCopySrcInvalid                 = errors.New("the copy-source is invalid")
	ErrCopyDestInvalid                = errors.New("the copy-destination is invalid")
	ErrDeletesCountInvalid            = errors.New("the deletes-count is invalid")
	ErrMaxKeysInvalid                 = errors.New("the max-keys is invalid")
	ErrEncodingTypeInvalid            = errors.New("the encoding-type is invalid")
	ErrPrefixInvalid                  = errors.New("the prefix is invalid")
	ErrMarkerInvalid                  = errors.New("the marker is invalid")
	ErrMarkerPrefixCombinationInvalid = errors.New("the marker-prefix combination is invalid")
	ErrContinuationTokenInvalid       = errors.New("the continuation-token is invalid")
	ErrStartAfterInvalid              = errors.New("the start-after is invalid")
	ErrPartNumberInvalid              = errors.New("the part-number is invalid")
	ErrPartsCountInvalid              = errors.New("the parts-count is invalid")
	ErrPartInvalid                    = errors.New("the part is invalid")
	ErrPartOrderInvalid               = errors.New("the part-order is invalid")
)

// ErrInvalidInputValue .
type ErrInvalidInputValue struct {
	msg string
}

func (err ErrInvalidInputValue) Error() string {
	return fmt.Sprintf("invalid input value: %s", err.msg)
}

// ErrTypeNotSet .
type ErrTypeNotSet struct {
	typ reflect.Type
}

func (err ErrTypeNotSet) Error() string {
	return fmt.Sprintf("type <%s> not set", err.typ.String())
}

// ErrPayloadNotSet .
type ErrPayloadNotSet struct {
	el string
}

func (err ErrPayloadNotSet) Error() string {
	return fmt.Sprintf("payload <%s> not set", err.el)
}

// ErrFailedDecodeXML .
type ErrFailedDecodeXML struct {
	err error
}

func (err ErrFailedDecodeXML) Error() string {
	return fmt.Sprintf("decode xml: %v", err.err)
}

// ErrWithUnsupportedParam .
type ErrWithUnsupportedParam struct {
	param string
}

func (err ErrWithUnsupportedParam) Error() string {
	return fmt.Sprintf("param %s is unsported", err.param)
}

// ErrFailedParseValue .
type ErrFailedParseValue struct {
	name string
	err  error
}

func (err ErrFailedParseValue) Name() string {
	return err.name
}

func (err ErrFailedParseValue) Error() string {
	return fmt.Sprintf("parse <%s> value: %v", err.name, err.err)
}

// ErrMissingRequiredParam .
type ErrMissingRequiredParam struct {
	param string
}

func (err ErrMissingRequiredParam) Error() string {
	return fmt.Sprintf("missing required param <%s>", err.param)
}
