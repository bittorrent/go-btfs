package store

import "errors"

var ErrBucketNotEmpty = errors.New("bucket not empty")

// BucketPolicyNotFound - no bucket policy found.
type BucketPolicyNotFound struct {
	Bucket string
	Err    error
}

func (e BucketPolicyNotFound) Error() string {
	return "No bucket policy configuration found for bucket: " + e.Bucket
}

// BucketNotFound - no bucket found.
type BucketNotFound struct {
	Bucket string
	Err    error
}

func (e BucketNotFound) Error() string {
	return "Not found for bucket: " + e.Bucket
}

type BucketTaggingNotFound struct {
	Bucket string
	Err    error
}

func (e BucketTaggingNotFound) Error() string {
	return "No bucket tagging configuration found for bucket: " + e.Bucket
}
