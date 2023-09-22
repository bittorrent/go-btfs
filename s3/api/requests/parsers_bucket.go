package requests

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/api/contexts"
	"github.com/bittorrent/go-btfs/s3/api/services/object"
	"net/http"
)

var createBucketSupports = fields{
	"ACL":                       true,
	"Bucket":                    true,
	"CreateBucketConfiguration": true,
}

func ParseCreateBucketRequest(r *http.Request) (args *object.CreateBucketArgs, err error) {
	var input s3.CreateBucketInput
	err = ParseLocation(r, &input, createBucketSupports)
	if err != nil {
		return
	}
	args = &object.CreateBucketArgs{
		UserId: contexts.GetAccessKey(r),
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
	err = ParseLocation(r, &input, headBucketSupports)
	if err != nil {
		return
	}
	args = &object.GetBucketArgs{
		UserId: contexts.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	return
}

var deleteBucketSupports = fields{
	"Bucket": true,
}

func ParseDeleteBucketRequest(r *http.Request) (args *object.DeleteBucketArgs, err error) {
	var input s3.DeleteBucketInput
	err = ParseLocation(r, &input, deleteBucketSupports)
	if err != nil {
		return
	}
	args = &object.DeleteBucketArgs{
		UserId: contexts.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	return
}

var listBucketsSupports = fields{}

func ParseListBucketsRequest(r *http.Request) (args *object.ListBucketsArgs, err error) {
	var input s3.ListBucketsInput
	err = ParseLocation(r, input, listBucketsSupports)
	if err != nil {
		return
	}
	args = &object.ListBucketsArgs{
		UserId: contexts.GetAccessKey(r),
	}
	return
}

var putBucketACLSupports = fields{
	"ACL":    true,
	"Bucket": true,
}

func ParsePutBucketAclRequest(r *http.Request) (args *object.PutBucketACLArgs, err error) {
	var input s3.PutBucketAclInput
	err = ParseLocation(r, &input, putBucketACLSupports)
	if err != nil {
		return
	}
	args = &object.PutBucketACLArgs{
		UserId: contexts.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	if err != nil {
		return
	}
	args.ACL, err = ValidateBucketACL(input.ACL)
	return
}

var getBucketACLSupports = fields{
	"Bucket": true,
}

func ParseGetBucketACLRequest(r *http.Request) (args *object.GetBucketACLArgs, err error) {
	var input s3.GetBucketAclInput
	err = ParseLocation(r, &input, getBucketACLSupports)
	if err != nil {
		return
	}
	args = &object.GetBucketACLArgs{
		UserId: contexts.GetAccessKey(r),
	}
	args.Bucket, err = ValidateBucketName(input.Bucket)
	return
}
