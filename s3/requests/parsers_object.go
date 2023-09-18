package requests

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"net/http"
)

var putObjectSupports = fields{
	"Body":            true,
	"Bucket":          true,
	"Key":             true,
	"ContentLength":   true,
	"ContentEncoding": true,
	"ContentType":     true,
	"Expires":         true,
	"ContentMD5":      true,
	"ChecksumSHA256":  true,
}

func ParsePutObjectRequest(r *http.Request) (args *object.PutObjectArgs, err error) {
	var input s3.PutObjectInput
	err = ParseInput(r, &input, putObjectSupports)
	if err != nil {
		return
	}
	args = &object.PutObjectArgs{
		AccessKey: cctx.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	args.Object, err = ValidateObjectName(input.Key)
	if err != nil {
		return
	}
	args.ContentLength, err = ValidateContentLength(input.ContentLength)
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
	checksumSHA256, err := ValidateCheckSum(input.ChecksumSHA256)
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
	"Bucket":            true,
	"Key":               true,
	"CopySource":        true,
	"ContentEncoding":   true,
	"ContentType":       true,
	"Expires":           true,
	"MetadataDirective": true,
}

func ParseCopyObjectRequest(r *http.Request) (args *object.CopyObjectArgs, err error) {
	var input s3.CopyObjectInput
	err = ParseInput(r, &input, copyObjectSupports)
	if err != nil {
		return
	}
	args = &object.CopyObjectArgs{
		AccessKey: cctx.GetAccessKey(r),
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
	err = ParseInput(r, &input, headObjectSupports)
	if err != nil {
		return
	}
	args = &object.GetObjectArgs{
		AccessKey: cctx.GetAccessKey(r),
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
	err = ParseInput(r, &input, getObjectSupports)
	if err != nil {
		return
	}
	args = &object.GetObjectArgs{
		AccessKey: cctx.GetAccessKey(r),
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
	err = ParseInput(r, &input, deleteObjectSupports)
	if err != nil {
		return
	}
	args = &object.DeleteObjectArgs{
		AccessKey: cctx.GetAccessKey(r),
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
	err = ParseInput(r, &input, deleteObjectsSupports)
	if err != nil {
		return
	}
	args = &object.DeleteObjectsArgs{
		AccessKey: cctx.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
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
	err = ParseInput(r, &input, listObjectsSupports)
	if err != nil {
		var er ErrFailedParseValue
		if errors.As(err, &er) && er.Name() == "max-keys" {
			err = ErrMaxKeysInvalid
		}
		return
	}
	args = &object.ListObjectsArgs{
		AccessKey: cctx.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	args.MaxKeys, err = ValidateMaxKeys(input.MaxKeys)
	if err != nil {
		return
	}
	args.Marker, args.Prefix, err = ValidateMarkerAndPrefix(input.Marker, input.Prefix)
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
