package object

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/etag"
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"github.com/google/uuid"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// CreateMultipartUpload create user specified multipart upload
func (s *service) CreateMultipartUpload(ctx context.Context, user, bucname, objname string, meta map[string]*string) (multipart *Multipart, err error) {
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
	allow := s.checkACL(bucket.Owner, bucket.ACL, user, action.CreateMultipartUploadAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Upload id
	uplid := uuid.NewString()

	// upload key
	uplkey := s.getUploadKey(bucname, objname, uplid)

	// Lock upload
	err = s.lock.Lock(ctx, uplkey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(uplkey)

	// Multipart upload
	multipart = &Multipart{
		Bucket:    bucname,
		Object:    objname,
		UploadID:  uplid,
		MetaData:  meta,
		Initiated: time.Now().UTC(),
	}

	// Put multipart upload
	err = s.providers.StateStore().Put(uplkey, multipart)

	return
}

// UploadPart upload user specified multipart part
func (s *service) UploadPart(ctx context.Context, user, bucname, objname, uplid string, partId int, body *hash.Reader, size int64) (part *Part, err error) {
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
	allow := s.checkACL(bucket.Owner, bucket.ACL, user, action.UploadPartAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Upload key
	uplkey := s.getUploadKey(bucname, objname, uplid)

	// Lock upload
	err = s.lock.Lock(ctx, uplkey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(uplkey)

	// Get multipart upload
	multipart, err := s.getMultipart(uplkey)
	if err != nil {
		return
	}
	if multipart == nil {
		err = ErrUploadNotFound
		return
	}

	// Upload part key
	prtkey := s.getUploadPartKey(uplkey, len(multipart.Parts))

	// Store part body
	cid, err := s.storeBody(ctx, body, prtkey)
	if err != nil {
		return
	}

	// Init a flag to mark if the part body should be removed, this
	// flag will be set to false if the multipart has been successfully put
	var removePartBody = true

	// Try to remove the part body
	defer func() {
		if removePartBody {
			_ = s.removeBody(ctx, cid, prtkey)
		}
	}()

	// Part
	part = &Part{
		Number:  partId,
		ETag:    body.ETag().String(),
		CID:     cid,
		Size:    size,
		ModTime: time.Now().UTC(),
	}

	// Append part
	multipart.Parts = append(multipart.Parts, part)

	// Put multipart upload
	err = s.providers.StateStore().Put(uplkey, multipart)
	if err != nil {
		return
	}

	// Set remove part body flag to false, because this part body has been referenced by the upload
	removePartBody = false

	return
}

// AbortMultipartUpload abort user specified multipart upload
func (s *service) AbortMultipartUpload(ctx context.Context, user, bucname, objname, uplid string) (err error) {
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
	allow := s.checkACL(bucket.Owner, bucket.ACL, user, action.AbortMultipartUploadAction)
	if !allow {
		err = ErrNotAllowed
		return
	}

	// Multipart upload key
	uplkey := s.getUploadKey(bucname, objname, uplid)

	// Lock upload
	err = s.lock.Lock(ctx, uplkey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(uplkey)

	// Get multipart upload
	multipart, err := s.getMultipart(uplkey)
	if err != nil {
		return
	}
	if multipart == nil {
		err = ErrUploadNotFound
		return
	}

	// Delete multipart upload
	err = s.providers.StateStore().Delete(uplkey)
	if err != nil {
		return
	}

	// Try to remove all parts body
	for i, part := range multipart.Parts {
		prtkey := s.getUploadPartKey(uplkey, i)
		_ = s.removeBody(ctx, part.CID, prtkey)
	}

	return
}

// CompleteMultiPartUpload complete user specified multipart upload
func (s *service) CompleteMultiPartUpload(ctx context.Context, user, bucname, objname, uplid string, parts []*CompletePart) (object *Object, err error) {
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
	allow := s.checkACL(bucket.Owner, bucket.ACL, user, action.CompleteMultipartUploadAction)
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

	// Get old object for try to remove the old body
	objectOld, err := s.getObject(objkey)
	if err != nil {
		return
	}

	// Upload key
	uplkey := s.getUploadKey(bucname, objname, uplid)

	// Lock upload
	err = s.lock.Lock(ctx, uplkey)
	if err != nil {
		return
	}
	defer s.lock.Unlock(uplkey)

	// Get multipart upload
	multipart, err := s.getMultipart(uplkey)
	if err != nil {
		return
	}
	if multipart == nil {
		err = ErrUploadNotFound
		return
	}

	// All parts body readers
	var readers []io.Reader

	// Try to close all parts body readers
	defer func() {
		for _, rdr := range readers {
			_ = rdr.(io.ReadCloser).Close()
		}
	}()

	// Total object size
	var size int64

	// Mapping of part number to part index in multipart.Parts
	idxmp := s.partIdxMap(multipart.Parts)

	// Iterate all parts to collect all body readers
	for i, part := range parts {
		// Index in multipart.Parts
		partIndex, ok := idxmp[part.PartNumber]

		// Part not exists in multipart
		if !ok {
			err = s3utils.InvalidPart{
				PartNumber: part.PartNumber,
				GotETag:    part.ETag,
			}
			return
		}

		// Got part in multipart.Parts
		gotPart := multipart.Parts[partIndex]

		// Canonicalize part etag
		part.ETag = s.canonicalizeETag(part.ETag)

		// Check got part etag with part etag
		if gotPart.ETag != part.ETag {
			err = s3utils.InvalidPart{
				PartNumber: part.PartNumber,
				ExpETag:    gotPart.ETag,
				GotETag:    part.ETag,
			}
			return
		}

		// All parts except the last part has to be at least 5MB.
		if (i < len(parts)-1) && !(gotPart.Size >= 0) {
			err = s3utils.PartTooSmall{
				PartNumber: part.PartNumber,
				PartSize:   gotPart.Size,
				PartETag:   part.ETag,
			}
			return
		}

		// Save for total object size.
		size += gotPart.Size

		// Get part body reader
		var rdr io.ReadCloser
		rdr, err = s.providers.FileStore().Cat(gotPart.CID)
		if err != nil {
			return
		}

		// Collect part body reader
		readers = append(readers, rdr)
	}

	// Concat all parts body to one
	body := io.MultiReader(readers...)

	// Store object body
	cid, err := s.storeBody(ctx, body, objkey)
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

	// Calculate multipart etag
	multiEtag, err := s.calcMultiETag(parts)
	if err != nil {
		return
	}

	// Current time
	now := time.Now().UTC()

	// Object
	object = &Object{
		Bucket:           bucname,
		Name:             objname,
		ModTime:          now,
		Size:             size,
		IsDir:            false,
		ETag: multiEtag.String(),
		CID:              cid,
		ACL:              "",
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		ContentType:      "",
		ContentEncoding:  "",
		Expires:          time.Time{},
		AccTime:          time.Time{},
		SuccessorModTime: now,
	}

	//  Set object content type
	ctyp := multipart.MetaData[strings.ToLower(consts.ContentType)]
	if ctyp != nil {
		object.ContentType = *ctyp
	}

	// Set object content encoding
	cecd := multipart.MetaData[strings.ToLower(consts.ContentEncoding)]
	if cecd != nil {
		object.ContentEncoding = *cecd
	}

	// Set object expires
	cexp := multipart.MetaData[strings.ToLower(consts.Expires)]
	if cexp != nil {
		exp, e := time.Parse(http.TimeFormat, *cexp)
		if e != nil {
			object.Expires = exp.UTC()
		}
	}

	// Put object
	err = s.providers.StateStore().Put(objkey, object)
	if err != nil {
		return
	}

	// Set remove object body flag to false, because it has been referenced by the object
	removeObjectBody = false

	// Try to remove old object body if exists, because it has been covered by new one
	if objectOld != nil {
		_ = s.removeBody(ctx, objectOld.CID, objkey)
	}

	// Remove multipart upload
	err = s.providers.StateStore().Delete(uplkey)
	if err != nil {
		return
	}

	// Try to remove all parts body, because they are no longer be referenced
	for i, part := range multipart.Parts {
		prtkey := s.getUploadPartKey(uplkey, i)
		_ = s.removeBody(ctx, part.CID, prtkey)
	}

	return
}

func (s *service) getMultipart(uplkey string) (multipart *Multipart, err error) {
	err = s.providers.StateStore().Get(uplkey, &multipart)
	if errors.Is(err, providers.ErrStateStoreNotFound) {
		err = nil
	}
	return
}

func (s *service) partIdxMap(parts []*Part) map[int]int {
	mp := make(map[int]int)
	for i, part := range parts {
		mp[part.Number] = i
	}
	return mp
}

var etagRegex = regexp.MustCompile("\"*?([^\"]*?)\"*?$")

// canonicalizeETag returns ETag with leading and trailing double-quotes removed,
// if any present
func (s *service) canonicalizeETag(etag string) string {
	return etagRegex.ReplaceAllString(etag, "$1")
}

func (s *service) calcMultiETag(parts []*CompletePart) (multiEtag etag.ETag, err error) {
	var completeETags []etag.ETag
	for _, part := range parts {
		var etg etag.ETag
		etg, err = etag.Parse(part.ETag)
		if err != nil {
			return
		}
		completeETags = append(completeETags, etg)
	}
	multiEtag = etag.Multipart(completeETags...)
	return
}
