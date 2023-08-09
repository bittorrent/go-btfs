package action

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
	//return set.Match(string(action), string(a))
	return true
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

//func createActionConditionKeyMap() map[Action]condition.KeySet {
//	commonKeys := []condition.Key{}
//	for _, keyName := range condition.CommonKeys {
//		commonKeys = append(commonKeys, keyName.ToKey())
//	}
//
//	return map[Action]condition.KeySet{
//		AbortMultipartUploadAction: condition.NewKeySet(commonKeys...),
//
//		CreateBucketAction: condition.NewKeySet(commonKeys...),
//
//		DeleteObjectAction: condition.NewKeySet(commonKeys...),
//
//		GetBucketLocationAction: condition.NewKeySet(commonKeys...),
//
//		GetBucketPolicyStatusAction: condition.NewKeySet(commonKeys...),
//
//		GetObjectAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3XAmzServerSideEncryption.ToKey(),
//				condition.S3XAmzServerSideEncryptionCustomerAlgorithm.ToKey(),
//			}, commonKeys...)...),
//
//		HeadBucketAction: condition.NewKeySet(commonKeys...),
//
//		ListAllMyBucketsAction: condition.NewKeySet(commonKeys...),
//
//		ListBucketAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3Prefix.ToKey(),
//				condition.S3Delimiter.ToKey(),
//				condition.S3MaxKeys.ToKey(),
//			}, commonKeys...)...),
//
//		ListBucketVersionsAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3Prefix.ToKey(),
//				condition.S3Delimiter.ToKey(),
//				condition.S3MaxKeys.ToKey(),
//			}, commonKeys...)...),
//
//		ListBucketMultipartUploadsAction: condition.NewKeySet(commonKeys...),
//
//		ListenNotificationAction: condition.NewKeySet(commonKeys...),
//
//		ListenBucketNotificationAction: condition.NewKeySet(commonKeys...),
//
//		ListMultipartUploadPartsAction: condition.NewKeySet(commonKeys...),
//
//		PutObjectAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3XAmzCopySource.ToKey(),
//				condition.S3XAmzServerSideEncryption.ToKey(),
//				condition.S3XAmzServerSideEncryptionCustomerAlgorithm.ToKey(),
//				condition.S3XAmzMetadataDirective.ToKey(),
//				condition.S3XAmzStorageClass.ToKey(),
//				condition.S3ObjectLockRetainUntilDate.ToKey(),
//				condition.S3ObjectLockMode.ToKey(),
//				condition.S3ObjectLockLegalHold.ToKey(),
//				condition.S3RequestObjectTagKeys.ToKey(),
//				condition.S3RequestObjectTag.ToKey(),
//			}, commonKeys...)...),
//
//		// https://docs.aws.amazon.com/AmazonS3/latest/dev/list_amazons3.html
//		// LockLegalHold is not supported with PutObjectRetentionAction
//		PutObjectRetentionAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3ObjectLockRemainingRetentionDays.ToKey(),
//				condition.S3ObjectLockRetainUntilDate.ToKey(),
//				condition.S3ObjectLockMode.ToKey(),
//			}, commonKeys...)...),
//
//		GetObjectRetentionAction: condition.NewKeySet(commonKeys...),
//		PutObjectLegalHoldAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3ObjectLockLegalHold.ToKey(),
//			}, commonKeys...)...),
//		GetObjectLegalHoldAction: condition.NewKeySet(commonKeys...),
//
//		// https://docs.aws.amazon.com/AmazonS3/latest/dev/list_amazons3.html
//		BypassGovernanceRetentionAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3ObjectLockRemainingRetentionDays.ToKey(),
//				condition.S3ObjectLockRetainUntilDate.ToKey(),
//				condition.S3ObjectLockMode.ToKey(),
//				condition.S3ObjectLockLegalHold.ToKey(),
//			}, commonKeys...)...),
//
//		GetBucketObjectLockConfigurationAction: condition.NewKeySet(commonKeys...),
//		PutBucketObjectLockConfigurationAction: condition.NewKeySet(commonKeys...),
//		GetBucketTaggingAction:                 condition.NewKeySet(commonKeys...),
//		PutBucketTaggingAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3RequestObjectTagKeys.ToKey(),
//				condition.S3RequestObjectTag.ToKey(),
//			}, commonKeys...)...),
//		PutObjectTaggingAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3RequestObjectTagKeys.ToKey(),
//				condition.S3RequestObjectTag.ToKey(),
//			}, commonKeys...)...),
//		GetObjectTaggingAction: condition.NewKeySet(commonKeys...),
//		DeleteObjectTaggingAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3RequestObjectTagKeys.ToKey(),
//				condition.S3RequestObjectTag.ToKey(),
//			}, commonKeys...)...),
//		PutObjectVersionTaggingAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3VersionID.ToKey(),
//				condition.S3RequestObjectTagKeys.ToKey(),
//				condition.S3RequestObjectTag.ToKey(),
//			}, commonKeys...)...),
//		GetObjectVersionAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3VersionID.ToKey(),
//			}, commonKeys...)...),
//		GetObjectVersionTaggingAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3VersionID.ToKey(),
//			}, commonKeys...)...),
//		DeleteObjectVersionAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3VersionID.ToKey(),
//			}, commonKeys...)...),
//		DeleteObjectVersionTaggingAction: condition.NewKeySet(
//			append([]condition.Key{
//				condition.S3VersionID.ToKey(),
//				condition.S3RequestObjectTagKeys.ToKey(),
//				condition.S3RequestObjectTag.ToKey(),
//			}, commonKeys...)...),
//		GetReplicationConfigurationAction:    condition.NewKeySet(commonKeys...),
//		PutReplicationConfigurationAction:    condition.NewKeySet(commonKeys...),
//		ReplicateObjectAction:                condition.NewKeySet(commonKeys...),
//		ReplicateDeleteAction:                condition.NewKeySet(commonKeys...),
//		ReplicateTagsAction:                  condition.NewKeySet(commonKeys...),
//		GetObjectVersionForReplicationAction: condition.NewKeySet(commonKeys...),
//		RestoreObjectAction:                  condition.NewKeySet(commonKeys...),
//	}
//}
//
//// ActionConditionKeyMap - holds mapping of supported condition key for an action.
//var ActionConditionKeyMap = createActionConditionKeyMap()
