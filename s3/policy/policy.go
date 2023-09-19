package policy

import (
	s3action "github.com/bittorrent/go-btfs/s3/action"
)

const (
	PublicReadWrite = "public-read-write"
	PublicRead      = "public-read"
	Private         = "private"
)

var rwActionMap = map[s3action.Action]struct{}{
	s3action.ListObjectsAction:             {},
	s3action.ListObjectsV2Action:           {},
	s3action.HeadObjectAction:              {},
	s3action.PutObjectAction:               {},
	s3action.GetObjectAction:               {},
	s3action.CopyObjectAction:              {},
	s3action.DeleteObjectAction:            {},
	s3action.DeleteObjectsAction:           {},
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
	if own {
		return true
	}

	if action.IsBucketAction() {
		return false
	}

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
