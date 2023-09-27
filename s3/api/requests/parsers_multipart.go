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

var createMultipartUploadSupports = fields{
	"Bucket":          true,
	"Key":             true,
	"CacheControl":    true,
	"ContentLength":   true,
	"ContentEncoding": true,
	"ContentType":     true,
	"Expires":         true,
}

func ParseCreateMultipartUploadRequest(r *http.Request) (args *object.CreateMultipartUploadArgs, err error) {
	var input s3.CreateMultipartUploadInput
	err = ParseLocation(r, &input, createMultipartUploadSupports)
	if err != nil {
		return
	}
	args = &object.CreateMultipartUploadArgs{
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

var uploadPartSupports = fields{
	"Body":           true,
	"Bucket":         true,
	"Key":            true,
	"UploadId":       true,
	"PartNumber":     true,
	"ContentLength":  true,
	"ContentMD5":     true,
	"ChecksumSHA256": true,
}

func ParseUploadPartRequest(r *http.Request) (args *object.UploadPartArgs, err error) {
	var input s3.UploadPartInput
	err = ParseLocation(r, &input, uploadPartSupports)
	if err != nil {
		var er ErrFailedParseValue
		if errors.As(err, &er) && er.Name() == consts.PartNumber {
			err = ErrPartNumberInvalid
		}
		return
	}
	args = &object.UploadPartArgs{
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
	args.UploadId, err = ValidateUploadId(input.UploadId)
	if err != nil {
		return
	}
	args.PartNumber, err = ValidatePartNumber(input.PartNumber)
	if err != nil {
		return
	}
	args.ContentLength, err = ValidateContentLength(input.ContentLength, consts.MaxPartSize)
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

var abortMultipartUploadSupports = fields{
	"Bucket":   true,
	"Key":      true,
	"UploadId": true,
}

func ParseAbortMultipartUploadRequest(r *http.Request) (args *object.AbortMultipartUploadArgs, err error) {
	var input s3.AbortMultipartUploadInput
	err = ParseLocation(r, &input, abortMultipartUploadSupports)
	if err != nil {
		return
	}
	args = &object.AbortMultipartUploadArgs{
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
	args.UploadId, err = ValidateUploadId(input.UploadId)
	return
}

var completeMultipartUploadSupports = fields{
	"Bucket":          true,
	"Key":             true,
	"UploadId":        true,
	"MultipartUpload": true,
	"ChecksumSHA256":  true,
}

func ParseCompleteMultipartUploadRequest(r *http.Request) (args *object.CompleteMultipartUploadArgs, err error) {
	var input s3.CompleteMultipartUploadInput
	err = ParseLocation(r, &input, completeMultipartUploadSupports)
	if err != nil {
		return
	}
	args = &object.CompleteMultipartUploadArgs{
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
	args.UploadId, err = ValidateUploadId(input.UploadId)
	if err != nil {
		return
	}
	size, err := ValidateContentLength(&r.ContentLength, consts.MaxXMLBodySize)
	if err != nil {
		return
	}
	checksumSHA256, err := ValidateChecksumSHA256(input.ChecksumSHA256)
	if err != nil {
		return
	}
	r.Body, err = hash.NewReader(r.Body, size, "", checksumSHA256, size)
	if err != nil {
		return
	}
	err = ParseXMLBody(r, &input)
	if err != nil {
		return
	}
	args.CompletedParts, err = ValidateCompletedMultipartUpload(input.MultipartUpload)
	return
}
