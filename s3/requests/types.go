package requests

// PutBucketRequest .
type PutBucketRequest struct {
	Bucket string
	ACL    string
	Region string
}

// HeadBucketRequest .
type HeadBucketRequest struct {
	Bucket string
}
