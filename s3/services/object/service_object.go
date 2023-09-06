package object

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"io"
	"net/http"
	"strings"
	"time"
)

// PutObject put a user specified object
func (s *service) PutObject(ctx context.Context, user, bucname, objname string, body *hash.Reader, size int64, meta map[string]string) (object *Object, err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Bucket key
	buckey := s.getBucketKey(bucname)

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
	allow := s.checkACL(bucket.Owner, bucket.ACL, user, action.PutObjectAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Object key
	objkey := s.getObjectKey(bucname, objname)

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
	cid, err := s.providers.FileStore().Store(body)
	if err != nil {
		return
	}

	// Init a flag to mark if the object body should be removed, this
	// flag will be set to false if the object has been successfully put
	var removeObjectBody = true

	// Try to remove stored body if put object failed
	defer func() {
		if removeObjectBody {
			_ = s.providers.FileStore().Remove(cid)
		}
	}()

	// now
	now := time.Now()

	// new object
	object = &Object{
		Bucket:           bucname,
		Name:             objname,
		ModTime:          now.UTC(),
		Size:             size,
		IsDir:            false,
		ETag:             body.ETag().String(),
		CID:              cid,
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		ACL:              meta[consts.AmzACL],
		ContentType:      meta[strings.ToLower(consts.ContentType)],
		ContentEncoding:  meta[strings.ToLower(consts.ContentEncoding)],
		SuccessorModTime: now.UTC(),
	}

	// set object expires
	exp, er := time.Parse(http.TimeFormat, meta[strings.ToLower(consts.Expires)])
	if er == nil {
		object.Expires = exp.UTC()
	}

	// put object
	err = s.providers.StateStore().Put(objkey, object)
	if err != nil {
		return
	}

	// Set remove object body flag to false, because it has been referenced by the object
	removeObjectBody = false

	// Try to remove old object body if exists, because it has been covered by new one
	if objectOld != nil {
		_ = s.providers.FileStore().Remove(objectOld.CID)
	}

	return
}

// CopyObject copy from a user specified source object to a desert object
func (s *service) CopyObject(ctx context.Context, user, srcBucname, srcObjname, dstBucname, dstObjname string, meta map[string]string) (dstObject *Object, err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Source bucket key
	srcBuckey := s.getBucketKey(srcBucname)

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
	srcAllow := s.checkACL(srcBucket.Owner, srcBucket.ACL, user, action.GetObjectAction)
	if !srcAllow {
		err = ErrNotAllowed
		return
	}

	// Source object key
	srcObjkey := s.getObjectKey(srcBucname, srcObjname)

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
	dstBuckey := s.getBucketKey(dstBucname)

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
	dstAllow := s.checkACL(dstBucket.Owner, dstBucket.ACL, user, action.PutObjectAction)
	if !dstAllow {
		err = ErrNotAllowed
		return
	}

	// Destination object key
	dstObjkey := s.getObjectKey(dstBucname, dstObjname)

	// Lock Destination object
	err = s.lock.Lock(ctx, dstObjkey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(dstObjkey)

	// now
	now := time.Now()

	// Destination object
	dstObject = &Object{
		Bucket:           dstBucname,
		Name:             dstObjname,
		ModTime:          now.UTC(),
		Size:             srcObject.Size,
		IsDir:            false,
		ETag:             srcObject.ETag,
		CID:              srcObject.CID,
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		ContentType:     srcObject.ContentType,
		ContentEncoding: srcObject.ContentEncoding,
		SuccessorModTime: now.UTC(),
		Expires:         srcObject.Expires,
	}

	// Set destination object metadata
	val, ok := meta[consts.ContentType]
	if ok {
		dstObject.ContentType = val
	}
	val, ok = meta[consts.ContentEncoding]
	if ok {
		dstObject.ContentEncoding = val
	}
	val, ok = meta[strings.ToLower(consts.Expires)]
	if ok {
		exp, er := time.Parse(http.TimeFormat, val)
		if er != nil {
			dstObject.Expires = exp.UTC()
		}
	}

	// Put destination object
	err = s.providers.StateStore().Put(dstObjkey, dstObject)

	return
}

// GetObject get a user specified object
func (s *service) GetObject(ctx context.Context, user, bucname, objname string, withBody bool) (object *Object, body io.ReadCloser, err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// bucket key
	buckey := s.getBucketKey(bucname)

	// RLock bucket
	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}
	defer func() {
		// RUnlock bucket just if getting failed
		if err != nil {
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
	allow := s.checkACL(bucket.Owner, bucket.ACL, user, action.GetObjectAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Object key
	objkey := s.getObjectKey(bucname, objname)

	// RLock object
	err = s.lock.RLock(ctx, objkey)
	if err != nil {
		return
	}
	defer func() {
		// RUnlock object just if getting failed
		if err != nil {
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
	if !withBody {
		return
	}

	// Get object body
	body, err = s.providers.FileStore().Cat(object.CID)
	if err != nil {
		return
	}

	// Wrap the body with timeout and unlock hooks,
	// this will enable the bucket and object keep rlocked until
	// read timout or read closed. Normally, these locks will
	// be released as soon as leave from the call
	body = WrapCleanReadCloser(
		body,
		s.closeBodyTimeout,
		func() {
			s.lock.RUnlock(objkey) // Note: Release object first
			s.lock.RUnlock(buckey)
		},
	)

	return
}

// DeleteObject delete a user specified object
func (s *service) DeleteObject(ctx context.Context, user, bucname, objname string) (err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Bucket key
	buckey := s.getBucketKey(bucname)

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
	allow := s.checkACL(bucket.Owner, bucket.ACL, user, action.DeleteObjectAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Object key
	objkey := s.getObjectKey(bucname, objname)

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
	err = s.providers.StateStore().Delete(objkey)
	if err != nil {
		return
	}

	// Try to delete object body
	_ = s.providers.FileStore().Remove(object.CID)

	return
}

// ListObjects list user specified objects
func (s *service) ListObjects(ctx context.Context, user, bucname, prefix, delimiter, marker string, max int64) (list *ObjectsList, err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Bucket key
	buckey := s.getBucketKey(bucname)

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
	allow := s.checkACL(bucket.Owner, bucket.ACL, user, action.ListObjectsAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	list = &ObjectsList{}

	// All bucket objects key prefix
	allObjectsKeyPrefix := s.getAllObjectsKeyPrefix(bucname)

	// List objects key prefix
	listObjectsKeyPrefix := allObjectsKeyPrefix + prefix

	// Accumulate count
	count := int64(0)

	// Flag mark if begin collect, it initialized to true if
	// marker is ""
	begin := marker == ""

	// Seen keys, used to group common keys
	seen := make(map[string]bool)

	// Iterate all objects with the specified prefix to collect and group specified range items
	err = s.providers.StateStore().Iterate(listObjectsKeyPrefix, func(key, _ []byte) (stop bool, er error) {
		// Object key
		objkey := string(key)

		// Object name
		objname := strings.TrimPrefix(objkey, allObjectsKeyPrefix)

		// Common prefix: if the part of object name without prefix include delimiter
		// it is the string truncated object name after the delimiter, else
		// it is the bucket name itself
		commonPrefix := objname
		if delimiter != "" {
			dl := len(delimiter)
			pl := len(prefix)
			di := strings.Index(objname[pl:], delimiter)
			if di >= 0 {
				commonPrefix = objname[:(pl + di + dl)]
			}
		}

		// If collect not begin, check the marker, if it is matched
		// with the common prefix, then begin collection from next iterate turn
		// and mark this common prefix as seen
		// note: common prefix also can be object name, so when marker is
		// an object name, the check will be also done correctly
		if !begin && marker == commonPrefix {
			begin = true
			seen[commonPrefix] = true
			return
		}

		// Not begin, jump the item
		if !begin {
			return
		}

		// Objects with same common prefix will be grouped into one
		// note: the objects without common prefix will present only once, so
		// it is not necessary to add these objects names in the seen map
		if seen[commonPrefix] {
			return
		}

		// Objects with common prefix grouped int one
		if commonPrefix != objname {
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
		if count == max {
			list.IsTruncated = true
			stop = true
		}

		return
	})

	return
}

func (s *service) getObject(objkey string) (object *Object, err error) {
	err = s.providers.StateStore().Get(objkey, &object)
	if errors.Is(err, providers.ErrStateStoreNotFound) {
		err = nil
	}
	return
}

// deleteObjectsByPrefix try to delete all objects with the specified common prefix
func (s *service) deleteObjectsByPrefix(objectsPrefix string) (err error) {
	err = s.providers.StateStore().Iterate(objectsPrefix, func(key, _ []byte) (stop bool, er error) {
		objkey := string(key)
		var object *Object
		er = s.providers.StateStore().Get(objkey, object)
		if er != nil {
			return
		}
		er = s.providers.StateStore().Delete(objkey)
		if er != nil {
			return
		}
		_ = s.providers.FileStore().Remove(object.CID)
		return
	})

	return
}
