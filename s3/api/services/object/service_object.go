package object

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/api/providers"
	"github.com/bittorrent/go-btfs/s3/utils"
	"io"
	"strings"
	"time"
)

// PutObject put a user specified object
func (s *service) PutObject(ctx context.Context, args *PutObjectArgs) (object *Object, err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Bucket key
	buckey := s.getBucketKey(args.Bucket)

	// RLock bucket
	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(buckey)

	// Get bucket
	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}
	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	// Check action ACL
	allow := s.checkACL(bucket.Owner, bucket.ACL, args.UserId, action.PutObjectAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Object key
	objkey := s.getObjectKey(args.Bucket, args.Object)

	// Lock object
	err = s.lock.Lock(ctx, objkey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(objkey)

	// Get old object
	objectOld, err := s.getObject(objkey)
	if err != nil {
		return
	}

	// Store object body
	cid, err := s.storeBody(ctx, args.Body, objkey)
	if err != nil {
		return
	}

	// Init a flag to mark if the object body should be removed, this
	// flag will be set to false if the object has been successfully put
	var removeObjectBody = true

	// Try to remove stored body if put object failed
	defer func() {
		if removeObjectBody {
			_ = s.removeBody(ctx, cid, objkey)
		}
	}()

	// now
	now := time.Now().UTC()

	// new object
	object = &Object{
		Bucket:           args.Bucket,
		Name:             args.Object,
		ModTime:          now,
		Size:             args.ContentLength,
		IsDir:            false,
		ETag:             args.Body.ETag().String(),
		CID:              cid,
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		ACL:              "",
		ContentType:      args.ContentType,
		ContentEncoding:  args.ContentEncoding,
		SuccessorModTime: now,
		Expires:          args.Expires,
	}

	// put object
	err = s.putObject(objkey, object)
	if err != nil {
		return
	}

	// Set remove object body flag to false, because it has been referenced by the object
	removeObjectBody = false

	// Try to remove old object body if exists, because it has been covered by new one
	if objectOld != nil {
		_ = s.removeBody(ctx, objectOld.CID, objkey)
	}

	return
}

// CopyObject copy from a user specified source object to a desert object
func (s *service) CopyObject(ctx context.Context, args *CopyObjectArgs) (dstObject *Object, err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Source bucket key
	srcBuckey := s.getBucketKey(args.SrcBucket)

	// RLock source bucket
	err = s.lock.RLock(ctx, srcBuckey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(srcBuckey)

	// Get source bucket
	srcBucket, err := s.getBucket(srcBuckey)
	if err != nil {
		return
	}
	if srcBucket == nil {
		err = ErrBucketNotFound
		return
	}

	// Check source action ACL
	srcAllow := s.checkACL(srcBucket.Owner, srcBucket.ACL, args.UserId, action.GetObjectAction)
	if !srcAllow {
		err = ErrNotAllowed
		return
	}

	// Source object key
	srcObjkey := s.getObjectKey(args.SrcBucket, args.SrcObject)

	// RLock source object
	err = s.lock.RLock(ctx, srcObjkey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(srcObjkey)

	// Get source object
	srcObject, err := s.getObject(srcObjkey)
	if err != nil {
		return
	}
	if srcObject == nil {
		err = ErrObjectNotFound
		return
	}

	// Desert bucket key
	dstBuckey := s.getBucketKey(args.Bucket)

	// RLock destination bucket
	err = s.lock.RLock(ctx, dstBuckey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(dstBuckey)

	// Get destination bucket
	dstBucket, err := s.getBucket(dstBuckey)
	if err != nil {
		return
	}
	if dstBucket == nil {
		err = ErrBucketNotFound
		return
	}

	// Check destination action ACL
	dstAllow := s.checkACL(dstBucket.Owner, dstBucket.ACL, args.UserId, action.PutObjectAction)
	if !dstAllow {
		err = ErrNotAllowed
		return
	}

	// Destination object key
	dstObjkey := s.getObjectKey(args.Bucket, args.Object)

	// Lock Destination object
	err = s.lock.Lock(ctx, dstObjkey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(dstObjkey)

	// Add body Refer
	err = s.addBodyRef(ctx, srcObject.CID, dstObjkey)
	if err != nil {
		return
	}

	// Mark if delete the cid ref
	deleteRef := true

	// If put new object failed, try to delete its reference
	defer func() {
		if deleteRef {
			_ = s.removeBodyRef(ctx, srcObject.CID, dstObjkey)
		}
	}()

	// Old desert object
	oldDstObject, err := s.getObject(dstObjkey)
	if err != nil {
		return
	}

	// now
	now := time.Now().UTC()

	// Destination object
	dstObject = &Object{
		Bucket:           args.Bucket,
		Name:             args.Object,
		ModTime:          now,
		Size:             srcObject.Size,
		IsDir:            false,
		ETag:             srcObject.ETag,
		CID:              srcObject.CID,
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		ContentType:      srcObject.ContentType,
		ContentEncoding:  srcObject.ContentEncoding,
		SuccessorModTime: now,
		Expires:          args.Expires,
	}

	// Replace metadata
	if args.ReplaceMeta {
		dstObject.ContentType = args.ContentType
		dstObject.ContentEncoding = args.ContentEncoding
	}

	// Put destination object
	err = s.putObject(dstObjkey, dstObject)
	if err != nil {
		return
	}

	// Mark the delete ref to false
	deleteRef = false

	// Try to remove the old object body
	if oldDstObject != nil {
		_ = s.removeBody(ctx, oldDstObject.CID, dstObjkey)
	}

	return
}

// GetObject get a user specified object
func (s *service) GetObject(ctx context.Context, args *GetObjectArgs) (object *Object, body io.ReadCloser, err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Unlock-later is a flag mark if the bucket or object will be unlocked later
	// if the flag is true, the bucket and object should not be unlocked as soon as leave the function call
	// they will be automatically unlocked after completely written the object body or write object body timeout
	unlockLater := false

	// bucket key
	buckey := s.getBucketKey(args.Bucket)

	// RLock bucket
	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}
	defer func() {
		if !unlockLater {
			s.lock.RUnlock(buckey)
		}
	}()

	// Get bucket
	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}
	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	// Check action ACL
	allow := s.checkACL(bucket.Owner, bucket.ACL, args.UserId, action.GetObjectAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Object key
	objkey := s.getObjectKey(args.Bucket, args.Object)

	// RLock object
	err = s.lock.RLock(ctx, objkey)
	if err != nil {
		return
	}
	defer func() {
		if !unlockLater {
			s.lock.RUnlock(objkey)
		}
	}()

	// Get object
	object, err = s.getObject(objkey)
	if err != nil {
		return
	}
	if object == nil {
		err = ErrObjectNotFound
		return
	}

	// no need body
	if !args.WithBody {
		return
	}

	// Get object body
	body, err = s.providers.FileStore().Cat(object.CID)
	if err != nil {
		return
	}

	// Set unlock-later flag to true to enable the bucket and object
	// will not be unlocked before completely written the response body
	unlockLater = true

	// Wrap the body with timeout and unlock hooks,
	// this will enable the bucket and object keep rlocked until
	// read timout or read closed. Normally, these locks will
	// be released as soon as leave from the call
	body = WrapCleanReadCloser(body, s.closeBodyTimeout, func() {
		s.lock.RUnlock(objkey) // Note: Release object first
		s.lock.RUnlock(buckey)
	})

	return
}

// DeleteObject delete a user specified object
func (s *service) DeleteObject(ctx context.Context, args *DeleteObjectArgs) (err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Bucket key
	buckey := s.getBucketKey(args.Bucket)

	// RLock bucket
	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(buckey)

	// Get bucket
	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}
	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	// Check action ACL
	allow := s.checkACL(bucket.Owner, bucket.ACL, args.UserId, action.DeleteObjectAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Object key
	objkey := s.getObjectKey(args.Bucket, args.Object)

	// Lock object
	err = s.lock.Lock(ctx, objkey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(objkey)

	// Get object
	object, err := s.getObject(objkey)
	if err != nil {
		return
	}
	if object == nil {
		err = ErrObjectNotFound
		return
	}

	// Delete object
	err = s.deleteObject(objkey)
	if err != nil {
		return
	}

	// Try to delete object body
	_ = s.removeBody(ctx, object.CID, objkey)

	return
}

// DeleteObjects delete multiple user specified objects
func (s *service) DeleteObjects(ctx context.Context, args *DeleteObjectsArgs) (deletes []*DeletedObject, err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Bucket key
	buckey := s.getBucketKey(args.Bucket)

	// RLock bucket
	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(buckey)

	// Get bucket
	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}
	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	// Check action ACL
	allow := s.checkACL(bucket.Owner, bucket.ACL, args.UserId, action.DeleteObjectAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	for _, deleteObj := range args.ToDeleteObjects {
		func(deleteObj *ToDeleteObject) {
			var er error
			// Collection delete result
			defer func() {
				if er != nil || !args.Quite {
					deletes = append(deletes, &DeletedObject{
						Object:    deleteObj.Object,
						DeleteErr: er,
					})
				}
			}()

			// Validate failed
			er = deleteObj.ValidateErr
			if er != nil {
				return
			}

			// Object key
			objkey := s.getObjectKey(args.Bucket, deleteObj.Object)

			// Lock object
			er = s.lock.Lock(ctx, objkey)
			if er != nil {
				return
			}
			defer s.lock.Unlock(objkey)

			// Get object
			object, er := s.getObject(objkey)
			if er != nil {
				return
			}
			if object == nil {
				er = ErrObjectNotFound
				return
			}

			// Delete object
			er = s.deleteObject(objkey)
			if er != nil {
				return
			}

			// Try to delete object body
			_ = s.removeBody(ctx, object.CID, objkey)

		}(deleteObj)
	}

	return
}

// ListObjects list user specified objects
func (s *service) ListObjects(ctx context.Context, args *ListObjectsArgs) (list *ObjectsList, err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Object list
	list = &ObjectsList{
		Args: args,
	}

	// Bucket key
	buckey := s.getBucketKey(args.Bucket)

	// RLock bucket
	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(buckey)

	// Get bucket
	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}
	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	// Check action ACL
	allow := s.checkACL(bucket.Owner, bucket.ACL, args.UserId, action.ListObjectsAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Set objects owner(objects owner is the bucket owner included them)
	list.Owner = bucket.Owner

	// MaxKeys is zero
	if args.MaxKeys == 0 {
		list.IsTruncated = true
		return
	}

	// All bucket objects key prefix
	allObjectsKeyPrefix := s.getAllObjectsKeyPrefix(args.Bucket)

	// List objects key prefix
	listObjectsKeyPrefix := allObjectsKeyPrefix + args.Prefix

	// Accumulate count
	count := int64(0)

	// Flag mark if begin collect, it initialized to true if
	// marker is ""
	begin := args.Marker == ""

	// Seen keys, used to group common keys
	seen := make(map[string]bool)

	// Delimiter length
	dl := len(args.Delimiter)

	// Prefix length
	pl := len(args.Prefix)

	// Iterate all objects with the specified prefix to collect and group specified range items
	err = s.providers.StateStore().Iterate(listObjectsKeyPrefix, func(key, _ []byte) (stop bool, er error) {
		// Object key
		objkey := string(key)

		// Object name
		objname := strings.TrimPrefix(objkey, allObjectsKeyPrefix)

		// Common prefix: if the part of object name without prefix include delimiter
		// it is the string truncated object name after the delimiter, else
		// it is empty string
		commonPrefix := ""
		if dl > 0 {
			di := strings.Index(objname[pl:], args.Delimiter)
			if di >= 0 {
				commonPrefix = objname[:(pl + di + dl)]
			}
		}

		// If collect not begin, check the marker, if it is matched
		// with the common prefix or object name, then begin collection from next iterate
		// and if common prefix matched, mark this common prefix as seen
		if !begin {
			if commonPrefix != "" && args.Marker == commonPrefix {
				seen[commonPrefix] = true
				begin = true
			} else if args.Marker == objname {
				begin = true
			}
			return
		}

		// ToDeleteObjects with same common prefix will be grouped into one
		// note: the objects without common prefix will present only once, so
		// it is not necessary to add these objects names in the seen map

		// ToDeleteObjects with common prefix grouped int one
		if commonPrefix != "" {
			if seen[commonPrefix] {
				return
			}
			list.Prefixes = append(list.Prefixes, commonPrefix)
			list.NextMarker = commonPrefix
			seen[commonPrefix] = true
		} else {
			// object without common prefix
			var object *Object
			er = s.providers.StateStore().Get(objkey, &object)
			if er != nil {
				return
			}
			list.Objects = append(list.Objects, object)
			list.NextMarker = objname
		}

		// Increment collection count
		count++

		// Check the count, if it matched the max, means
		// the collect is complete, but the items may remain, so stop the
		// iteration, and mark the list was truncated
		if count == args.MaxKeys {
			list.IsTruncated = true
			stop = true
		}

		return
	})

	return
}

func (s *service) ListObjectsV2(ctx context.Context, args *ListObjectsV2Args) (list *ObjectsListV2, err error) {
	// Args v1
	v1Args := &ListObjectsArgs{
		UserId:       args.UserId,
		Bucket:       args.Bucket,
		MaxKeys:      args.MaxKeys,
		Marker:       utils.CoalesceStr(args.Token, args.After),
		Prefix:       args.Prefix,
		Delimiter:    args.Delimiter,
		EncodingType: args.EncodingType,
	}

	// Get v1 list
	v1List, err := s.ListObjects(ctx, v1Args)
	if err != nil {
		return
	}

	// List v2
	list = &ObjectsListV2{
		Args:                  args,
		Owner:                 v1List.Owner,
		IsTruncated:           v1List.IsTruncated,
		NextContinuationToken: v1List.NextMarker,
		Objects:               v1List.Objects,
		Prefixes:              v1List.Prefixes,
	}
	return
}

func (s *service) deleteObject(objkey string) (err error) {
	err = s.providers.StateStore().Delete(objkey)
	return
}

func (s *service) putObject(objkey string, object *Object) (err error) {
	err = s.providers.StateStore().Put(objkey, object)
	return
}

func (s *service) getObject(objkey string) (object *Object, err error) {
	err = s.providers.StateStore().Get(objkey, &object)
	if errors.Is(err, providers.ErrStateStoreNotFound) {
		err = nil
	}
	return
}

// GetObjectACL get user specified object ACL(bucket acl)
func (s *service) GetObjectACL(ctx context.Context, args *GetObjectACLArgs) (acl *ACL, err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Bucket key
	buckey := s.getBucketKey(args.Bucket)

	// RLock bucket
	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(buckey)

	// Get bucket
	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}
	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	// Check action ACL
	allow := s.checkACL(bucket.Owner, bucket.ACL, args.UserId, action.GetBucketAclAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Object key
	objkey := s.getObjectKey(args.Bucket, args.Object)

	// RLock object
	err = s.lock.RLock(ctx, objkey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(objkey)

	// Get object
	object, err := s.getObject(objkey)
	if err != nil {
		return
	}
	if object == nil {
		err = ErrObjectNotFound
	}

	// ACL
	acl = &ACL{
		Owner: bucket.Owner,
		ACL:   bucket.ACL,
	}

	return
}
