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
func (s *service) CreateBucket(ctx context.Context, user, bucname, region, acl string) (err error) {
	buckey := s.getBucketKey(bucname)

	ctx, cancel := s.opctx(ctx)
	defer cancel()

	err = s.lock.Lock(ctx, buckey)
	if err != nil {
		return
	}

	defer s.lock.Unlock(buckey)

	allow := s.checkAcl(user, acl, user, action.CreateBucketAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	bucket, err := s.getBucket(buckey)
	if err == nil {
		return
	}

	if bucket != nil {
		err = ErrBucketAlreadyExists
		return
	}

	err = s.providers.StateStore().Put(
		buckey,
		&Bucket{
			Name:    bucname,
			Region:  region,
			Owner:   user,
			Acl:     acl,
			Created: time.Now().UTC(),
		},
	)

	return
}

// GetBucket get a bucket for the specified user
func (s *service) GetBucket(ctx context.Context, user, bucname string) (bucket *Bucket, err error) {
	buckey := s.getBucketKey(bucname)

	ctx, cancel := s.opctx(ctx)
	defer cancel()

	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}

	defer s.lock.RUnlock(buckey)

	bucket, err = s.getBucket(buckey)
	if err != nil {
		return
	}

	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	allow := s.checkAcl(bucket.Owner, bucket.Acl, user, action.HeadBucketAction)
	if !allow {
		err = ErrNotAllowed
	}

	return
}

// DeleteBucket delete the specified user bucket and all the bucket's objects
func (s *service) DeleteBucket(ctx context.Context, user, bucname string) (err error) {
	buckey := s.getBucketKey(bucname)

	ctx, cancel := s.opctx(ctx)
	defer cancel()

	err = s.lock.Lock(ctx, buckey)
	if err != nil {
		return
	}

	defer s.lock.Unlock(buckey)

	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}

	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	allow := s.checkAcl(bucket.Owner, bucket.Acl, user, action.DeleteBucketAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	err = s.providers.StateStore().Delete(buckey)

	// todo: delete all objects below to this bucket

	return
}

// GetAllBuckets get all buckets of the specified user
func (s *service) GetAllBuckets(ctx context.Context, user string) (list []*Bucket, err error) {
	bucprefix := s.getBucketKeyPrefix()

	ctx, cancel := s.opctx(ctx)
	defer cancel()

	allow := s.checkAcl(user, policy.Private, user, action.ListBucketAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	err = s.providers.StateStore().Iterate(bucprefix, func(key, _ []byte) (stop bool, er error) {
		defer func() {
			if er != nil {
				stop = true
			}
		}()

		er = ctx.Err()
		if er != nil {
			return
		}

		var bucket *Bucket

		er = s.providers.StateStore().Get(string(key), bucket)
		if er != nil {
			return
		}

		if bucket.Owner == user {
			list = append(list, bucket)
		}

		return
	})

	return
}

// PutBucketAcl update the acl field value of the specified user's bucket
func (s *service) PutBucketAcl(ctx context.Context, user, bucname, acl string) (err error) {
	buckey := s.getBucketKey(bucname)

	ctx, cancel := s.opctx(ctx)
	defer cancel()

	err = s.lock.Lock(ctx, buckey)
	if err != nil {
		return
	}

	defer s.lock.Unlock(buckey)

	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}

	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	allow := s.checkAcl(bucket.Owner, bucket.Acl, user, action.PutBucketAclAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	bucket.Acl = acl

	err = s.providers.StateStore().Put(buckey, bucket)

	return
}

// GetBucketAcl get the acl field value of the specified user's bucket
func (s *service) GetBucketAcl(ctx context.Context, user, bucname string) (acl string, err error) {
	buckey := s.getBucketKey(bucname)

	ctx, cancel := s.opctx(ctx)
	defer cancel()

	err = s.lock.RLock(ctx, buckey)
	if err != nil {
		return
	}

	defer s.lock.RUnlock(buckey)

	bucket, err := s.getBucket(buckey)
	if err != nil {
		return
	}

	if bucket == nil {
		err = ErrBucketNotFound
		return
	}

	allow := s.checkAcl(bucket.Owner, bucket.Acl, user, action.GetBucketAclAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	acl = bucket.Acl

	return
}

func (s *service) getBucket(buckey string) (bucket *Bucket, err error) {
	err = s.providers.StateStore().Get(buckey, bucket)
	if errors.Is(err, providers.ErrStateStoreNotFound) {
		err = nil
	}
	return
}
