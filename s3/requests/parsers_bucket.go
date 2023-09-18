package requests

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"net/http"
)

var createBucketSupports = fields{
	"ACL":                       true,
	"Bucket":                    true,
	"CreateBucketConfiguration": true,
}

func ParseCreateBucketRequest(r *http.Request) (args *object.CreateBucketArgs, err error) {
	var input s3.CreateBucketInput
	err = ParseInput(r, &input, createBucketSupports)
	if err != nil {
		return
	}
	args = &object.CreateBucketArgs{
		AccessKey: cctx.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	args.ACL, err = ValidateBucketACL(input.ACL)
	if err != nil {
		return
	}
	args.Region, err = ValidateCreateBucketConfiguration(input.CreateBucketConfiguration)
	return
}

var headBucketSupports = fields{
	"Bucket": true,
}

func ParseHeadBucketRequest(r *http.Request) (args *object.GetBucketArgs, err error) {
	var input s3.HeadBucketInput
	err = ParseInput(r, &input, headBucketSupports)
	if err != nil {
		return
	}
	args = &object.GetBucketArgs{
		AccessKey: cctx.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	return
}

var deleteBucketSupports = fields{
	"Bucket": true,
}

func ParseDeleteBucketRequest(r *http.Request) (args *object.DeleteBucketArgs, err error) {
	var input s3.DeleteBucketInput
	err = ParseInput(r, &input, deleteBucketSupports)
	if err != nil {
		return
	}
	args = &object.DeleteBucketArgs{
		AccessKey: cctx.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	return
}

var listBucketsSupports = fields{}

func ParseListBucketsRequest(r *http.Request) (args *object.ListBucketsArgs, err error) {
	var input s3.ListBucketsInput
	err = ParseInput(r, input, listBucketsSupports)
	if err != nil {
		return
	}
	args = &object.ListBucketsArgs{
		AccessKey: cctx.GetAccessKey(r),
	}
	return
}

var putBucketACLSupports = fields{
	"ACL":    true,
	"Bucket": true,
}

func ParsePutBucketAclRequest(r *http.Request) (args *object.PutBucketACLArgs, err error) {
	var input s3.PutBucketAclInput
	err = ParseInput(r, &input, putBucketACLSupports)
	if err != nil {
		return
	}
	args = &object.PutBucketACLArgs{
		AccessKey: cctx.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	args.Bucket, err = ValidateBucketACL(input.ACL)
	return
}

var getBucketACLSupports = fields{
	"Bucket": true,
}

func ParseGetBucketACLRequest(r *http.Request) (args *object.GetBucketACLArgs, err error) {
	var input s3.GetBucketAclInput
	err = ParseInput(r, &input, getBucketACLSupports)
	if err != nil {
		return
	}
	args = &object.GetBucketACLArgs{
		AccessKey: cctx.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	return
}
