package requests

import (
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/api/contexts"
	"github.com/bittorrent/go-btfs/s3/api/services/object"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/hash"
)

var putObjectSupports = fields{
	"Body":   true,
	"Bucket": true,
	"Key":    true,
	// The browser some time automatically add this CacheControl header
	// just allow, do not handle
	"CacheControl":    true,
	"ContentLength":   true,
	"ContentEncoding": true,
	"ContentType":     true,
	"Expires":         true,
	"ContentMD5":      true,
	"ChecksumSHA256":  true,
}

func ParsePutObjectRequest(r *http.Request) (args *object.PutObjectArgs, err error) {
	var input s3.PutObjectInput
	err = ParseLocation(r, &input, putObjectSupports)
	if err != nil {
		return
	}
	args = &object.PutObjectArgs{
		UserId: contexts.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	args.Object, err = ValidateObjectName(input.Key)
	if err != nil {
		return
	}
	args.ContentLength, err = ValidateContentLength(input.ContentLength, consts.MaxObjectSize)
	if err != nil {
		return
	}
	args.ContentType, err = ValidateContentType(input.ContentType)
	if err != nil {
		return
	}
	args.ContentEncoding, err = ValidateContentEncoding(input.ContentEncoding)
	if err != nil {
		return
	}
	args.Expires, err = ValidateExpires(input.Expires)
	if err != nil {
		return
	}
	contentMD5, err := ValidateContentMD5(input.ContentMD5)
	if err != nil {
		return
	}
	checksumSHA256, err := ValidateChecksumSHA256(input.ChecksumSHA256)
	if err != nil {
		return
	}
	args.Body, err = hash.NewReader(
		r.Body, args.ContentLength, contentMD5,
		checksumSHA256, args.ContentLength,
	)
	return
}

var copyObjectSupports = fields{
	"Bucket":     true,
	"Key":        true,
	"CopySource": true,
	// The browser some time automatically add this CacheControl header
	// just allow, do not handle
	"CacheControl":      true,
	"ContentEncoding":   true,
	"ContentType":       true,
	"Expires":           true,
	"MetadataDirective": true,
}

func ParseCopyObjectRequest(r *http.Request) (args *object.CopyObjectArgs, err error) {
	var input s3.CopyObjectInput
	err = ParseLocation(r, &input, copyObjectSupports)
	if err != nil {
		return
	}
	args = &object.CopyObjectArgs{
		UserId: contexts.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	args.Object, err = ValidateObjectName(input.Key)
	if err != nil {
		return
	}
	args.SrcBucket, args.SrcObject, err = ValidateCopySource(input.CopySource)
	if err != nil {
		return
	}
	args.ReplaceMeta, err = ValidateMetadataDirective(input.MetadataDirective)
	if err != nil {
		return
	}
	if args.Bucket == args.SrcBucket && args.Object == args.SrcObject {
		err = ErrCopyDestInvalid
		return
	}
	args.ContentType, err = ValidateContentType(input.ContentType)
	if err != nil {
		return
	}
	args.ContentEncoding, err = ValidateContentEncoding(input.ContentEncoding)
	if err != nil {
		return
	}
	args.Expires, err = ValidateExpires(input.Expires)
	return
}

var headObjectSupports = fields{
	"Bucket": true,
	"Key":    true,
}

func ParseHeadObjectRequest(r *http.Request) (args *object.GetObjectArgs, err error) {
	var input s3.HeadObjectInput
	err = ParseLocation(r, &input, headObjectSupports)
	if err != nil {
		return
	}
	args = &object.GetObjectArgs{
		UserId: contexts.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	args.Object, err = ValidateObjectName(input.Key)
	if err != nil {
		return
	}
	args.WithBody = false
	return
}

var getObjectSupports = fields{
	"Bucket": true,
	"Key":    true,
}

func ParseGetObjectRequest(r *http.Request) (args *object.GetObjectArgs, err error) {
	var input s3.GetObjectInput
	err = ParseLocation(r, &input, getObjectSupports)
	if err != nil {
		return
	}
	args = &object.GetObjectArgs{
		UserId: contexts.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	args.Object, err = ValidateObjectName(input.Key)
	if err != nil {
		return
	}
	args.WithBody = true
	return
}

var deleteObjectSupports = fields{
	"Bucket": true,
	"Key":    true,
}

func ParseDeleteObjectRequest(r *http.Request) (args *object.DeleteObjectArgs, err error) {
	var input s3.DeleteObjectInput
	err = ParseLocation(r, &input, deleteObjectSupports)
	if err != nil {
		return
	}
	args = &object.DeleteObjectArgs{
		UserId: contexts.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	args.Object, err = ValidateObjectName(input.Key)
	return
}

var deleteObjectsSupports = fields{
	"Bucket": true,
	"Delete": true,
}

func ParseDeleteObjectsRequest(r *http.Request) (args *object.DeleteObjectsArgs, err error) {
	var input s3.DeleteObjectsInput
	err = ParseLocation(r, &input, deleteObjectsSupports)
	if err != nil {
		return
	}
	args = &object.DeleteObjectsArgs{
		UserId: contexts.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	size, err := ValidateContentLength(&r.ContentLength, consts.MaxXMLBodySize)
	if err != nil {
		return
	}
	r.Body, err = hash.NewReader(r.Body, size, "", "", size)
	if err != nil {
		return
	}
	err = ParseXMLBody(r, &input)
	if err != nil {
		return
	}
	args.ToDeleteObjects, args.Quite, err = ValidateObjectsDelete(input.Delete)
	return
}

var listObjectsSupports = fields{
	"Bucket":       true,
	"MaxKeys":      true,
	"Prefix":       true,
	"Marker":       true,
	"Delimiter":    true,
	"EncodingType": true,
}

func ParseListObjectsRequest(r *http.Request) (args *object.ListObjectsArgs, err error) {
	var input s3.ListObjectsInput
	err = ParseLocation(r, &input, listObjectsSupports)
	if err != nil {
		var er ErrFailedParseValue
		if errors.As(err, &er) && er.Name() == consts.MaxKeys {
			err = ErrMaxKeysInvalid
		}
		return
	}
	args = &object.ListObjectsArgs{
		UserId: contexts.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	args.MaxKeys, err = ValidateMaxKeys(input.MaxKeys)
	if err != nil {
		return
	}
	args.Marker, err = ValidateMarker(input.Marker)
	if err != nil {
		return
	}
	args.Prefix, err = ValidatePrefix(input.Prefix)
	if err != nil {
		return
	}
	err = ValidateMarkerAndPrefixCombination(args.Marker, args.Prefix)
	if err != nil {
		return
	}
	args.Delimiter, err = ValidateDelimiter(input.Delimiter)
	if err != nil {
		return
	}
	args.EncodingType, err = ValidateEncodingType(input.EncodingType)
	return
}

var listObjectsV2Supports = fields{
	"Bucket":            true,
	"MaxKeys":           true,
	"Prefix":            true,
	"ContinuationToken": true,
	"StartAfter":        true,
	"Delimiter":         true,
	"EncodingType":      true,
	"FetchOwner":        true,
}

func ParseListObjectsV2Request(r *http.Request) (args *object.ListObjectsV2Args, err error) {
	var input s3.ListObjectsV2Input
	err = ParseLocation(r, &input, listObjectsV2Supports)
	if err != nil {
		var er ErrFailedParseValue
		if errors.As(err, &er) && er.Name() == "max-keys" {
			err = ErrMaxKeysInvalid
		}
		return
	}
	args = &object.ListObjectsV2Args{
		UserId: contexts.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	args.MaxKeys, err = ValidateMaxKeys(input.MaxKeys)
	if err != nil {
		return
	}
	args.Token, err = ValidateContinuationToken(input.ContinuationToken)
	if err != nil {
		return
	}
	args.After, err = ValidateStartAfter(input.StartAfter)
	if err != nil {
		return
	}
	err = ValidateMarkerAndPrefixCombination(args.Token, args.Prefix)
	if err != nil {
		return
	}
	err = ValidateMarkerAndPrefixCombination(args.After, args.Prefix)
	if err != nil {
		return
	}
	args.Delimiter, err = ValidateDelimiter(input.Delimiter)
	if err != nil {
		return
	}
	args.EncodingType, err = ValidateEncodingType(input.EncodingType)
	if err != nil {
		return
	}
	args.FetchOwner, err = ValidateFetchOwner(input.FetchOwner)
	return
}

var getObjectACLSupports = fields{
	"Bucket": true,
	"Key":    true,
}

func ParseGetObjectACLRequest(r *http.Request) (args *object.GetObjectACLArgs, err error) {
	var input s3.GetObjectAclInput
	err = ParseLocation(r, &input, getObjectACLSupports)
	if err != nil {
		return
	}
	args = &object.GetObjectACLArgs{
		UserId: contexts.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	args.Object, err = ValidateObjectName(input.Key)
	return
}
