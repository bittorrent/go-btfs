package action

import "github.com/bittorrent/go-btfs/s3d/set"

type Action string

// ActionSet - set of actions.
// https://docs.aws.amazon.com/service-authorization/latest/reference/list_amazons3.html#amazons3-actions-as-permissions
const (
	//--- bucket

	// CreateBucketAction - CreateBucket Rest API action.
	CreateBucketAction = "s3d:CreateBucket"

	// HeadBucketAction - HeadBucket Rest API action.
	HeadBucketAction = "s3d:HeadBucket"

	// ListBucketAction - ListBucket Rest API action.
	ListBucketAction = "s3d:ListBucket"

	// DeleteBucketAction - DeleteBucket Rest API action.
	DeleteBucketAction = "s3d:DeleteBucket"

	// PutBucketAclAction - PutBucketAcl Rest API action.
	PutBucketAclAction = "s3d:PutBucketAcl"

	// GetBucketAclAction - GetBucketAcl Rest API action.
	GetBucketAclAction = "s3d:GetBucketAcl"

	//--- object

	// ListObjectsAction - ListObjects Rest API action.
	ListObjectsAction = "s3d:ListObjects"

	// ListObjectsV2Action - ListObjectsV2 Rest API action.
	ListObjectsV2Action = "s3d:ListObjectsV2"

	// HeadObjectAction - HeadObject Rest API action.
	HeadObjectAction = "s3d:HeadObject"

	// PutObjectAction - PutObject Rest API action.
	PutObjectAction = "s3d:PutObject"

	// GetObjectAction - GetObject Rest API action.
	GetObjectAction = "s3d:GetObject"

	// CopyObjectAction - CopyObject Rest API action.
	CopyObjectAction = "s3d:CopyObject"

	// DeleteObjectAction - DeleteObject Rest API action.
	DeleteObjectAction = "s3d:DeleteObject"

	// DeleteObjectsAction - DeleteObjects Rest API action.
	DeleteObjectsAction = "s3d:DeleteObjects"

	//--- multipart upload

	// CreateMultipartUploadAction - CreateMultipartUpload Rest API action.
	CreateMultipartUploadAction Action = "s3d:CreateMultipartUpload"

	// AbortMultipartUploadAction - AbortMultipartUpload Rest API action.
	AbortMultipartUploadAction Action = "s3d:AbortMultipartUpload"

	// CompleteMultipartUploadAction - CompleteMultipartUpload Rest API action.
	CompleteMultipartUploadAction Action = "s3d:CompleteMultipartUpload"

	// UploadPartAction - UploadPartUpload Rest API action.
	UploadPartAction Action = "s3d:UploadPartUpload"
)

// SupportedActions List of all supported actions.
var SupportedActions = map[Action]struct{}{
	CreateBucketAction: {},
	HeadBucketAction:   {},
	ListBucketAction:   {},
	DeleteBucketAction: {},
	PutBucketAclAction: {},
	GetBucketAclAction: {},

	ListObjectsAction:   {},
	ListObjectsV2Action: {},
	HeadObjectAction:    {},
	PutObjectAction:     {},
	GetObjectAction:     {},
	CopyObjectAction:    {},
	DeleteObjectAction:  {},
	DeleteObjectsAction: {},

	CreateMultipartUploadAction:   {},
	AbortMultipartUploadAction:    {},
	CompleteMultipartUploadAction: {},
	UploadPartAction:              {},
}

// IsValid - checks if action is valid or not.
func (action Action) IsValid() bool {
	for supAction := range SupportedActions {
		if action.Match(supAction) {
			return true
		}
	}
	return false
}

// Match - matches action name with action patter.
func (action Action) Match(a Action) bool {
	return set.Match(string(action), string(a))
	//return true
}

// List of all supported object actions.
var supportedBucketActions = map[Action]struct{}{
	CreateBucketAction: {},
	HeadBucketAction:   {},
	ListBucketAction:   {},
	DeleteBucketAction: {},
	PutBucketAclAction: {},
	GetBucketAclAction: {},
}

// IsBucketAction - returns whether action is bucket type or not.
func (action Action) IsBucketAction() bool {
	_, ok := supportedBucketActions[action]
	return ok
}

// List of all supported object actions.
var supportedObjectActions = map[Action]struct{}{
	ListObjectsAction:   {},
	ListObjectsV2Action: {},
	HeadObjectAction:    {},
	PutObjectAction:     {},
	GetObjectAction:     {},
	CopyObjectAction:    {},
	DeleteObjectAction:  {},
	DeleteObjectsAction: {},

	CreateMultipartUploadAction:   {},
	AbortMultipartUploadAction:    {},
	CompleteMultipartUploadAction: {},
	UploadPartAction:              {},
}

// IsObjectAction - returns whether action is object type or not.
func (action Action) IsObjectAction() bool {
	_, ok := supportedObjectActions[action]
	return ok
}
