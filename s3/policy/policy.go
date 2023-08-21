package policy

import (
	s3action "github.com/bittorrent/go-btfs/s3/action"
)

const (
	// PublicReadWrite 公开读写，适用于桶ACL和对象ACL
	PublicReadWrite = "public-read-write"

	// PublicRead 公开读，适用于桶ACL和对象ACL
	PublicRead = "public-read"

	// Private 私有，适用于桶ACL和对象ACL
	Private = "private"
)

// 支持匿名公开读写的action集合
var rwActionMap = map[s3action.Action]struct{}{
	s3action.ListObjectsAction:   {},
	s3action.ListObjectsV2Action: {},
	s3action.HeadObjectAction:    {},
	s3action.PutObjectAction:     {},
	s3action.GetObjectAction:     {},
	s3action.CopyObjectAction:    {},
	s3action.DeleteObjectAction:  {},
	s3action.DeleteObjectsAction: {},

	s3action.CreateMultipartUploadAction:   {},
	s3action.AbortMultipartUploadAction:    {},
	s3action.CompleteMultipartUploadAction: {},
	s3action.UploadPartAction:              {},
}

// checkActionInPublicReadWrite - returns whether action is RW or not.
func checkActionInPublicReadWrite(action s3action.Action) bool {
	_, ok := rwActionMap[action]
	return ok
}

// 支持匿名公开读的action集合
var rdActionMap = map[s3action.Action]struct{}{
	s3action.ListObjectsAction:   {},
	s3action.ListObjectsV2Action: {},
	s3action.HeadObjectAction:    {},
	s3action.GetObjectAction:     {},
}

// checkActionInPublicRead - returns whether action is Read or not.
func checkActionInPublicRead(action s3action.Action) bool {
	_, ok := rdActionMap[action]
	return ok
}

func IsAllowed(own bool, acl string, action s3action.Action) (allow bool) {
	// 1.如果是自己，都能操作
	if own {
		return true
	}

	// 2.如果是别人，不能操作bucket
	if action.IsBucketAction() {
		return false
	}

	// 2.如果是别人，区分acl操作object
	if action.IsObjectAction() {
		switch acl {
		case Private:
			return own
		case PublicRead:
			return checkActionInPublicRead(action)
		case PublicReadWrite:
			return checkActionInPublicReadWrite(action)
		}
	}

	return false
}
