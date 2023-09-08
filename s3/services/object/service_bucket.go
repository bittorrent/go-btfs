package object

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/policy"
	"github.com/bittorrent/go-btfs/s3/providers"
	"time"

	"github.com/bittorrent/go-btfs/s3/action"
)

// CreateBucket create a new bucket for the specified user
func (s *service) CreateBucket(ctx context.Context, user, bucname, region, acl string) (bucket *Bucket, err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Bucket key
	buckey := s.getBucketKey(bucname)

	// Lock bucket
	err = s.lock.Lock(ctx, buckey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(buckey)

	// Get old bucket
	bucketOld, err := s.getBucket(buckey)
	if err != nil {
		return
	}
	if bucketOld != nil {
		err = ErrBucketAlreadyExists
		return
	}

	// Check action ACL
	allow := s.checkACL(user, policy.Private, user, action.CreateBucketAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Bucket
	bucket = &Bucket{
		Name:    bucname,
		Region:  region,
		Owner:   user,
		ACL:     acl,
		Created: time.Now().UTC(),
	}

	// Put bucket
	err = s.providers.StateStore().Put(buckey, bucket)

	return
}

// GetBucket get a user specified bucket
func (s *service) GetBucket(ctx context.Context, user, bucname string) (bucket *Bucket, err error) {
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
	bucket, err = s.getBucket(buckey)
	if err != nil {
		return
	}
	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	// Check action ACL
	allow := s.checkACL(bucket.Owner, bucket.ACL, user, action.HeadBucketAction)
	if !allow {
		err = ErrNotAllowed
	}

	return
}

// DeleteBucket delete a user specified bucket and clear all bucket objects and uploads
func (s *service) DeleteBucket(ctx context.Context, user, bucname string) (err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Bucket key
	buckey := s.getBucketKey(bucname)

	// Lock bucket
	err = s.lock.Lock(ctx, buckey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(buckey)

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
	allow := s.checkACL(bucket.Owner, bucket.ACL, user, action.DeleteBucketAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Check if bucket is empty
	empty, err := s.isBucketEmpty(bucname)
	if err != nil {
		return
	}
	if !empty {
		err = ErrBucketeNotEmpty
		return
	}

	// Delete bucket
	err = s.providers.StateStore().Delete(buckey)

	return
}

// GetAllBuckets get all buckets of the specified user
func (s *service) GetAllBuckets(ctx context.Context, user string) (list []*Bucket, err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Check action ACL
	allow := s.checkACL(user, policy.Private, user, action.ListBucketAction)
	if !allow {
		err = ErrNotAllowed
		return
	}
	// All buckets prefix
	bucketsPrefix := s.getAllBucketsKeyPrefix()

	// Collect user's buckets from all buckets
	err = s.providers.StateStore().Iterate(bucketsPrefix, func(key, _ []byte) (stop bool, er error) {
		// Stop the iteration if error occurred
		defer func() {
			if er != nil {
				stop = true
			}
		}()

		// Bucket key
		buckey := string(key)

		// Get Bucket
		bucket, er := s.getBucket(buckey)
		if er != nil {
			return
		}

		// Bucket has been deleted
		if bucket == nil {
			return
		}

		// Collect user's bucket
		if bucket.Owner == user {
			list = append(list, bucket)
		}

		return
	})

	return
}

// PutBucketACL update user specified bucket's ACL field value
func (s *service) PutBucketACL(ctx context.Context, user, bucname, acl string) (err error) {
	// Operation context
	ctx, cancel := s.opctx(ctx)
	defer cancel()

	// Bucket key
	buckey := s.getBucketKey(bucname)

	// Lock bucket
	err = s.lock.Lock(ctx, buckey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(buckey)

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
	allow := s.checkACL(bucket.Owner, bucket.ACL, user, action.PutBucketAclAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Update bucket ACL
	bucket.ACL = acl

	// Put bucket
	err = s.providers.StateStore().Put(buckey, bucket)

	return
}

// GetBucketACL get user specified bucket ACL field value
func (s *service) GetBucketACL(ctx context.Context, user, bucname string) (acl string, err error) {
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
	allow := s.checkACL(bucket.Owner, bucket.ACL, user, action.GetBucketAclAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Get ACL field value
	acl = bucket.ACL

	return
}

// EmptyBucket check if the user specified bucked is empty
func (s *service) isBucketEmpty(bucname string) (empty bool, err error) {
	// All bucket objects prefix
	objectsPrefix := s.getAllObjectsKeyPrefix(bucname)

	// Initially set empty to true
	empty = true

	// Iterate the bucket objects, if no item, empty keep true
	// if at least one, set empty to false, and stop iterate
	err = s.providers.StateStore().Iterate(objectsPrefix, func(_, _ []byte) (stop bool, er error) {
		empty = false
		stop = true
		return
	})

	// If bucket have at least one object, return not empty, else check if bucket
	// have at least one upload
	if !empty {
		return
	}

	// All bucket uploads prefix
	uploadsPrefix := s.getAllUploadsKeyPrefix(bucname)

	// Set empty to false if bucket has at least one upload
	err = s.providers.StateStore().Iterate(uploadsPrefix, func(_, _ []byte) (stop bool, er error) {
		empty = false
		stop = true
		return
	})

	return
}

func (s *service) getBucket(buckey string) (bucket *Bucket, err error) {
	err = s.providers.StateStore().Get(buckey, &bucket)
	if errors.Is(err, providers.ErrStateStoreNotFound) {
		err = nil
	}
	return
}
