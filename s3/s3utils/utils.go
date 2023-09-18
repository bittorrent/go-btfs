package s3utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"strings"
	"unicode/utf8"
)

// GenericError - generic object layer error.
type GenericError struct {
	Bucket string
	Object string
	Err    error
}

// Bucket related errors.

// BucketNameInvalid - bucket name provided is invalid.
type BucketNameInvalid GenericError

// Error returns string an error formatted as the given text.
func (e BucketNameInvalid) Error() string {
	return "bucket name invalid: " + e.Bucket
}

// Object related errors.

// ObjectNameInvalid - object name provided is invalid.
type ObjectNameInvalid GenericError

// ObjectNameTooLong - object name too long.
type ObjectNameTooLong GenericError

// ObjectNamePrefixAsSlash - object name has a slash as prefix.
type ObjectNamePrefixAsSlash GenericError

// Error returns string an error formatted as the given text.
func (e ObjectNameInvalid) Error() string {
	return "Object name invalid: " + e.Bucket + "/" + e.Object
}

// Error returns string an error formatted as the given text.
func (e ObjectNameTooLong) Error() string {
	return "Object name too long: " + e.Bucket + "/" + e.Object
}

// Error returns string an error formatted as the given text.
func (e ObjectNamePrefixAsSlash) Error() string {
	return "Object name contains forward slash as prefix: " + e.Bucket + "/" + e.Object
}

// InvalidUploadIDKeyCombination - invalid upload id and key marker combination.
type InvalidUploadIDKeyCombination struct {
	UploadIDMarker, KeyMarker string
}

func (e InvalidUploadIDKeyCombination) Error() string {
	return fmt.Sprintf("Invalid combination of uploadID marker '%s' and marker '%s'", e.UploadIDMarker, e.KeyMarker)
}

// InvalidMarkerPrefixCombination - invalid marker and prefix combination.
type InvalidMarkerPrefixCombination struct {
	Marker, Prefix string
}

func (e InvalidMarkerPrefixCombination) Error() string {
	return fmt.Sprintf("Invalid combination of marker '%s' and prefix '%s'", e.Marker, e.Prefix)
}

// Multipart related errors.

// MalformedUploadID malformed upload id.
type MalformedUploadID struct {
	UploadID string
}

func (e MalformedUploadID) Error() string {
	return "Malformed upload id " + e.UploadID
}

// InvalidUploadID invalid upload id.
type InvalidUploadID struct {
	Bucket   string
	Object   string
	UploadID string
}

func (e InvalidUploadID) Error() string {
	return "Invalid upload id " + e.UploadID
}

// InvalidPart One or more of the specified parts could not be found
type InvalidPart struct {
	PartNumber int
	ExpETag    string
	GotETag    string
}

func (e InvalidPart) Error() string {
	return fmt.Sprintf("Specified part could not be found. PartNumber %d, Expected %s, got %s",
		e.PartNumber, e.ExpETag, e.GotETag)
}

// PartTooSmall - error if part size is less than 5MB.
type PartTooSmall struct {
	PartSize   int64
	PartNumber int
	PartETag   string
}

func (e PartTooSmall) Error() string {
	return fmt.Sprintf("Part size for %d should be at least 5MB", e.PartNumber)
}

// PartTooBig returned if size of part is bigger than the allowed limit.
type PartTooBig struct{}

func (e PartTooBig) Error() string {
	return "Part size bigger than the allowed limit"
}

// We support '.' with bucket names but we fallback to using path
// style requests instead for such buckets.
var (
	validBucketName       = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9\.\-\_\:]{1,61}[A-Za-z0-9]$`)
	validBucketNameStrict = regexp.MustCompile(`^[a-z0-9][a-z0-9\.\-]{1,61}[a-z0-9]$`)
	ipAddress             = regexp.MustCompile(`^(\d+\.){3}\d+$`)
)

// Common checker for both stricter and basic validation.
func checkBucketName(bucketName string, strict bool) (err error) {
	if strings.TrimSpace(bucketName) == "" {
		return errors.New("bucket name cannot be empty")
	}
	if len(bucketName) < 3 {
		return errors.New("bucket name cannot be shorter than 3 characters")
	}
	if len(bucketName) > 63 {
		return errors.New("bucket name cannot be longer than 63 characters")
	}
	if ipAddress.MatchString(bucketName) {
		return errors.New("bucket name cannot be an ip address")
	}
	if strings.Contains(bucketName, "..") || strings.Contains(bucketName, ".-") || strings.Contains(bucketName, "-.") {
		return errors.New("bucket name contains invalid characters")
	}
	if strict {
		if !validBucketNameStrict.MatchString(bucketName) {
			err = errors.New("bucket name contains invalid characters")
		}
		return err
	}
	if !validBucketName.MatchString(bucketName) {
		err = errors.New("bucket name contains invalid characters")
	}
	return err
}

// CheckValidBucketName - checks if we have a valid input bucket name.
func CheckValidBucketName(bucketName string) (err error) {
	return checkBucketName(bucketName, false)
}

// CheckValidBucketNameStrict - checks if we have a valid input bucket name.
// This is a stricter version.
// - http://docs.aws.amazon.com/AmazonS3/latest/dev/UsingBucket.html
func CheckValidBucketNameStrict(bucketName string) (err error) {
	return checkBucketName(bucketName, true)
}

// Checks on GetObject arguments, bucket and object.
func CheckGetObjArgs(ctx context.Context, bucket, object string) error {
	return checkBucketAndObjectNames(ctx, bucket, object)
}

// Checks on DeleteObject arguments, bucket and object.
func CheckDelObjArgs(ctx context.Context, bucket, object string) error {
	return checkBucketAndObjectNames(ctx, bucket, object)
}

// Checks bucket and object name validity, returns nil if both are valid.
func checkBucketAndObjectNames(ctx context.Context, bucket, object string) error {
	// Verify if bucket is valid.
	if CheckValidBucketName(bucket) != nil {
		return BucketNameInvalid{Bucket: bucket}
	}
	// Verify if object is valid.
	if len(object) == 0 {
		return ObjectNameInvalid{Bucket: bucket, Object: object}
	}
	if !IsValidObjectPrefix(object) {
		return ObjectNameInvalid{Bucket: bucket, Object: object}
	}
	return nil
}

// Checks for all ListObjects arguments validity.
func CheckListObjsArgs(ctx context.Context, bucket, prefix, marker string) error {
	// Validates object prefix validity after bucket exists.
	if !IsValidObjectPrefix(prefix) {
		return ObjectNameInvalid{
			Bucket: bucket,
			Object: prefix,
		}
	}
	// Verify if marker has prefix.
	if marker != "" && !strings.HasPrefix(marker, prefix) {
		return InvalidMarkerPrefixCombination{
			Marker: marker,
			Prefix: prefix,
		}
	}
	return nil
}

// Checks for all ListMultipartUploads arguments validity.
func CheckListMultipartArgs(ctx context.Context, bucket, prefix, keyMarker, uploadIDMarker, delimiter string) error {
	if err := CheckListObjsArgs(ctx, bucket, prefix, keyMarker); err != nil {
		return err
	}
	if uploadIDMarker != "" {
		if strings.HasSuffix(keyMarker, SlashSeparator) {
			return InvalidUploadIDKeyCombination{
				UploadIDMarker: uploadIDMarker,
				KeyMarker:      keyMarker,
			}
		}
		if _, err := uuid.Parse(uploadIDMarker); err != nil {
			return MalformedUploadID{
				UploadID: uploadIDMarker,
			}
		}
	}
	return nil
}

// Checks for NewMultipartUpload arguments validity, also validates if bucket exists.
func CheckNewMultipartArgs(ctx context.Context, bucket, object string) error {
	return checkObjectArgs(ctx, bucket, object)
}

// Checks for PutObjectPart arguments validity, also validates if bucket exists.
func CheckPutObjectPartArgs(ctx context.Context, bucket, object string) error {
	return checkObjectArgs(ctx, bucket, object)
}

// Checks for ListParts arguments validity, also validates if bucket exists.
func CheckListPartsArgs(ctx context.Context, bucket, object string) error {
	return checkObjectArgs(ctx, bucket, object)
}

// Checks for CompleteMultipartUpload arguments validity, also validates if bucket exists.
func CheckCompleteMultipartArgs(ctx context.Context, bucket, object string) error {
	return checkObjectArgs(ctx, bucket, object)
}

// Checks for AbortMultipartUpload arguments validity, also validates if bucket exists.
func CheckAbortMultipartArgs(ctx context.Context, bucket, object string) error {
	return checkObjectArgs(ctx, bucket, object)
}

// Checks Object arguments validity, also validates if bucket exists.
func checkObjectArgs(ctx context.Context, bucket, object string) error {
	if err := checkObjectNameForLengthAndSlash(bucket, object); err != nil {
		return err
	}

	// Validates object name validity after bucket exists.
	if !IsValidObjectName(object) {
		return ObjectNameInvalid{
			Bucket: bucket,
			Object: object,
		}
	}

	return nil
}

// Checks for PutObject arguments validity, also validates if bucket exists.
func CheckPutObjectArgs(ctx context.Context, bucket, object string) error {
	if err := checkObjectNameForLengthAndSlash(bucket, object); err != nil {
		return err
	}
	if len(object) == 0 ||
		!IsValidObjectPrefix(object) {
		return ObjectNameInvalid{
			Bucket: bucket,
			Object: object,
		}
	}
	return nil
}

// SlashSeparator - slash separator.
const SlashSeparator = "/"

// IsValidObjectName verifies an object name in accordance with Amazon's
// requirements. It cannot exceed 1024 characters and must be a valid UTF8
// string.
//
// See:
// http://docs.aws.amazon.com/AmazonS3/latest/dev/UsingMetadata.html
//
// You should avoid the following characters in a key name because of
// significant special handling for consistency across all
// applications.
//
// Rejects strings with following characters.
//
// - Backslash ("\")
//
// additionally minio does not support object names with trailing SlashSeparator.
func IsValidObjectName(object string) bool {
	if len(object) == 0 {
		return false
	}
	if strings.HasSuffix(object, SlashSeparator) {
		return false
	}
	return IsValidObjectPrefix(object)
}

// IsValidObjectPrefix verifies whether the prefix is a valid object name.
// Its valid to have a empty prefix.
func IsValidObjectPrefix(object string) bool {
	if hasBadPathComponent(object) {
		return false
	}
	if !utf8.ValidString(object) {
		return false
	}
	if strings.Contains(object, `//`) {
		return false
	}
	return true
}

// checkObjectNameForLengthAndSlash -check for the validity of object name length and prefis as slash
func checkObjectNameForLengthAndSlash(bucket, object string) error {
	// Check for the length of object name
	if len(object) > 1024 {
		return ObjectNameTooLong{
			Bucket: bucket,
			Object: object,
		}
	}
	// Check for slash as prefix in object name
	if strings.HasPrefix(object, SlashSeparator) {
		return ObjectNamePrefixAsSlash{
			Bucket: bucket,
			Object: object,
		}
	}
	return nil
}

// Bad path components to be rejected by the path validity handler.
const (
	dotdotComponent = ".."
	dotComponent    = "."
)

// Check if the incoming path has bad path components,
// such as ".." and "."
func hasBadPathComponent(path string) bool {
	path = strings.TrimSpace(path)
	for _, p := range strings.Split(path, SlashSeparator) {
		switch strings.TrimSpace(p) {
		case dotdotComponent:
			return true
		case dotComponent:
			return true
		}
	}
	return false
}
