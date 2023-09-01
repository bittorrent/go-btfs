package requests

// PutBucketRequest .
type PutBucketRequest struct {
	User   string
	Bucket string
	ACL    string
	Region string
}

// HeadBucketRequest .
type HeadBucketRequest struct {
	Bucket string
}
