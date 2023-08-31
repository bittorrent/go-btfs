package object

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"io"
	"net/http"
	"strings"
	"time"
)

// PutObject put a user specified object
func (s *service) PutObject(ctx context.Context, user string, bucname, objname string, reader *hash.Reader, size int64, meta map[string]string) (object *Object, err error) {
	// operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// bucket key
	buckey := s.getBucketKey(bucname)

	// rlock bucket
	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(buckey)

	// get bucket
	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}
	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	// check acl
	allow := s.checkAcl(bucket.Owner, bucket.Acl, user, action.PutObjectAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// object key
	objkey := s.getObjectKey(bucname, objname)

	// lock object
	err = s.lock.Lock(ctx, objkey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(objkey)

	// get old object
	oldObject, err := s.getObject(objkey)
	if err != nil {
		return
	}

	// remove old file, if old object exists and put new object successfully
	defer func() {
		if oldObject != nil && err == nil {
			_ = s.providers.FileStore().Remove(oldObject.Cid)
			// todo: log this remove error
		}
	}()

	// store file
	cid, err := s.providers.FileStore().Store(reader)
	if err != nil {
		return
	}

	// now
	now := time.Now()

	// new object
	object = &Object{
		Bucket:           bucname,
		Name:             objname,
		ModTime:          now.UTC(),
		Size:             size,
		IsDir:            false,
		ETag:             reader.ETag().String(),
		Cid:              cid,
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		Acl:              meta[consts.AmzACL],
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

	return
}

// CopyObject copy from a user specified source object to a desert object
func (s *service) CopyObject(ctx context.Context, user string, srcBucname, srcObjname, dstBucname, dstObjname string, meta map[string]string) (dstObject *Object, err error) {
	// operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// source bucket key
	srcBuckey := s.getBucketKey(srcBucname)

	// rlock source bucket
	err = s.lock.RLock(ctx, srcBuckey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(srcBuckey)

	// get source bucket
	srcBucket, err := s.getBucket(srcBuckey)
	if err != nil {
		return
	}
	if srcBucket == nil {
		err = ErrBucketNotFound
		return
	}

	// check source acl
	srcAllow := s.checkAcl(srcBucket.Owner, srcBucket.Acl, user, action.GetObjectAction)
	if !srcAllow {
		err = ErrNotAllowed
		return
	}

	// source object key
	srcObjkey := s.getObjectKey(srcBucname, srcObjname)

	// rlock source object
	err = s.lock.RLock(ctx, srcObjkey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(srcObjkey)

	// get source object
	srcObject, err := s.getObject(srcObjkey)
	if err != nil {
		return
	}
	if srcObject == nil {
		err = ErrObjectNotFound
		return
	}

	// desert bucket key
	dstBuckey := s.getBucketKey(dstBucname)

	// rlock desert bucket
	err = s.lock.RLock(ctx, dstBuckey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(dstBuckey)

	// get desert bucket
	dstBucket, err := s.getBucket(dstBuckey)
	if err != nil {
		return
	}
	if dstBucket == nil {
		err = ErrBucketNotFound
		return
	}

	// check desert acl
	dstAllow := s.checkAcl(dstBucket.Owner, dstBucket.Acl, user, action.PutObjectAction)
	if !dstAllow {
		err = ErrNotAllowed
		return
	}

	// desert object key
	dstObjkey := s.getObjectKey(dstBucname, dstObjname)

	// lock desert object
	err = s.lock.Lock(ctx, dstObjkey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(dstObjkey)

	// now
	now := time.Now()

	// desert object
	dstObject = &Object{
		Bucket:           dstBucname,
		Name:             dstObjname,
		ModTime:          now.UTC(),
		Size:             srcObject.Size,
		IsDir:            false,
		ETag:             srcObject.ETag,
		Cid:              srcObject.Cid,
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		ContentType:      meta[strings.ToLower(consts.ContentType)],
		ContentEncoding:  meta[strings.ToLower(consts.ContentEncoding)],
		SuccessorModTime: now.UTC(),
	}

	// set object desert expires
	exp, er := time.Parse(http.TimeFormat, strings.ToLower(consts.Expires))
	if er != nil {
		dstObject.Expires = exp.UTC()
	}

	// put desert object
	err = s.providers.StateStore().Put(dstObjkey, dstObject)

	return
}

// GetObject get an object for the specified user
func (s *service) GetObject(ctx context.Context, user, bucname, objname string) (object *Object, body io.ReadCloser, err error) {
	// operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// bucket key
	buckey := s.getBucketKey(bucname)

	// rlock bucket
	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}
	defer func() {
		// rUnlock bucket just if getting failed
		if err != nil {
			s.lock.RUnlock(buckey)
		}
	}()

	// get bucket
	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}
	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	// check acl
	allow := s.checkAcl(bucket.Owner, bucket.Acl, user, action.GetObjectAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// object key
	objkey := s.getObjectKey(bucname, objname)

	// rlock object
	err = s.lock.RLock(ctx, objkey)
	if err != nil {
		return
	}
	defer func() {
		// rUnlock object just if getting failed
		if err != nil {
			s.lock.RUnlock(objkey)
		}
	}()

	// get object
	object, err = s.getObject(objkey)
	if err != nil {
		return
	}
	if object == nil {
		err = ErrObjectNotFound
		return
	}

	// get object body
	body, err = s.providers.FileStore().Cat(object.Cid)
	if err != nil {
		return
	}

	// wrap the body with timeout and unlock hooks
	// this will enable the bucket and object keep rlocked until
	// read timout or read closed. Normally, these locks will
	// be released as soon as leave from the call
	body = WrapCleanReadCloser(
		body,
		s.readObjectTimeout,
		func() {
			s.lock.RUnlock(objkey) // note: release object first
			s.lock.RUnlock(buckey)
		},
	)

	return
}

// DeleteObject delete a user specified object
func (s *service) DeleteObject(ctx context.Context, user, bucname, objname string) (err error) {
	// operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// bucket key
	buckey := s.getBucketKey(bucname)

	// rlock bucket
	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(buckey)

	// get bucket
	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}
	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	// check acl
	allow := s.checkAcl(bucket.Owner, bucket.Acl, user, action.DeleteObjectAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// object key
	objkey := s.getObjectKey(bucname, objname)

	// lock object
	err = s.lock.Lock(ctx, objkey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(objkey)

	// get object
	object, err := s.getObject(objkey)
	if err != nil {
		return
	}
	if object == nil {
		err = ErrObjectNotFound
		return
	}

	// delete object body
	err = s.providers.FileStore().Remove(object.Cid)
	if err != nil {
		return
	}

	// delete object
	err = s.providers.StateStore().Delete(objkey)

	return
}

// ListObjects list user specified objects
func (s *service) ListObjects(ctx context.Context, user, bucname, prefix, delimiter, marker string, max int) (list *ObjectsList, err error) {
	// operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// bucket key
	buckey := s.getBucketKey(bucname)

	// rlock bucket
	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(buckey)

	// get bucket
	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}
	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	// check acl
	allow := s.checkAcl(bucket.Owner, bucket.Acl, user, action.ListObjectsAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// object key prefix
	objkeyPrefix := s.getObjectKeyPrefix(bucname)

	// objects key prefix
	objskeyPrefix := objkeyPrefix + prefix

	// accumulate count
	count := 0

	// begin collect
	begin := marker == ""

	// seen keys
	seen := make(map[string]bool)

	// iterate all objects with the specified prefix to collect and group specified range items
	err = s.providers.StateStore().Iterate(objskeyPrefix, func(key, _ []byte) (stop bool, er error) {
		// object key
		objkey := string(key)

		// object name
		objname := objkey[len(objkeyPrefix):]

		// common prefix: if the part of object name without prefix include delimiter
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

		// if collect not begin, check the marker, if it is matched
		// with the common prefix, then begin collection from next iterate turn
		// and mark this common prefix as seen
		// note: common prefix also can be object name, so when marker is
		// an object name, the check will be also done correctly
		if !begin && marker == commonPrefix {
			begin = true
			seen[commonPrefix] = true
			return
		}

		// no begin, jump the item
		if !begin {
			return
		}

		// objects with same common prefix will be grouped into one
		// note: the objects without common prefix will present only once, so
		// it is not necessary to add these objects names in the seen map
		if seen[commonPrefix] {
			return
		}

		// objects with common prefix grouped int one
		if commonPrefix != objname {
			list.Prefixes = append(list.Prefixes, commonPrefix)
			list.NextMarker = commonPrefix
			seen[commonPrefix] = true
		} else {
			// object without common prefix
			var object *Object
			er = s.providers.StateStore().Get(objkey, object)
			if er != nil {
				return
			}
			list.Objects = append(list.Objects, object)
			list.NextMarker = objname
		}

		// increment collection count
		count++

		// check the count, if it matched the max, means
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

// EmptyBucket check if the user specified bucked is empty
func (s *service) EmptyBucket(ctx context.Context, user, bucname string) (empty bool, err error) {
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// bucket key
	buckey := s.getBucketKey(bucname)

	// rlock bucket
	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}
	defer s.lock.RUnlock(buckey)

	// get bucket
	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}
	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	// check acl
	allow := s.checkAcl(bucket.Owner, bucket.Acl, user, action.HeadBucketAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// object key prefix
	objkeyPrefix := s.getObjectKeyPrefix(bucname)

	// initially set empty to true
	empty = true

	// iterate the bucket objects, if no item, empty keep true
	// if at least one, set empty to false, and stop iterate
	err = s.providers.StateStore().Iterate(objkeyPrefix, func(_, _ []byte) (stop bool, er error) {
		empty = false
		stop = true
		return
	})

	return
}
