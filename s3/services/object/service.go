package object

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/etag"
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
)

const (
	// bigFileThreshold is the point where we add readahead to put operations.
	bigFileThreshold = 64 * humanize.MiByte
	// equals unixfsChunkSize
	chunkSize int = 1 << 20

	objectKeyFormat        = "obj/%s/%s"
	allObjectPrefixFormat  = "obj/%s/%s"
	allObjectSeekKeyFormat = "obj/%s/%s"

	uploadKeyFormat        = "uploadObj/%s/%s/%s"
	allUploadPrefixFormat  = "uploadObj/%s/%s"
	allUploadSeekKeyFormat = "uploadObj/%s/%s/%s"

	deleteKeyFormat       = "delObj/%s"
	allDeletePrefixFormat = "delObj/"

	globalOperationTimeout = 5 * time.Minute
	deleteOperationTimeout = 1 * time.Minute

	maxCpuPercent        = 60
	maxUsedMemoryPercent = 80
)

var etagRegex = regexp.MustCompile("\"*?([^\"]*?)\"*?$")

var _ Service = (*service)(nil)

// service captures all bucket metadata for a given cluster.
type service struct {
	providers providers.Providerser
}

// NewService - creates new policy system.
func NewService(providers providers.Providerser, options ...Option) Service {
	s := &service{
		providers: providers,
	}
	for _, option := range options {
		option(s)
	}
	return s
}

func getObjectKey(bucname, objname string) string {
	return fmt.Sprintf(objectKeyFormat, bucname, objname)
}

func getUploadKey(bucname, objname, uploadID string) string {
	return fmt.Sprintf(uploadKeyFormat, bucname, objname, uploadID)
}

func (s *service) PutObject(ctx context.Context, bucname, objname string, reader *hash.Reader, size int64, meta map[string]string) (obj Object, err error) {
	cid, err := s.providers.GetFileStore().Store(reader)
	if err != nil {
		return
	}

	obj = Object{
		Bucket:           bucname,
		Name:             objname,
		ModTime:          time.Now().UTC(),
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
		SuccessorModTime: time.Now().UTC(),
	}

	// Update expires
	if exp, ok := meta[strings.ToLower(consts.Expires)]; ok {
		if t, e := time.Parse(http.TimeFormat, exp); e == nil {
			obj.Expires = t.UTC()
		}
	}

	err = s.providers.GetStateStore().Put(getObjectKey(bucname, objname), obj)
	if err != nil {
		return
	}

	return
}

// CopyObject store object
func (s *service) CopyObject(ctx context.Context, bucket, object string, info Object, size int64, meta map[string]string) (Object, error) {
	obj := Object{
		Bucket:           bucket,
		Name:             object,
		ModTime:          time.Now().UTC(),
		Size:             size,
		IsDir:            false,
		ETag:             info.ETag,
		Cid:              info.Cid,
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		ContentType:      meta[strings.ToLower(consts.ContentType)],
		ContentEncoding:  meta[strings.ToLower(consts.ContentEncoding)],
		SuccessorModTime: time.Now().UTC(),
	}
	// Update expires
	if exp, ok := meta[strings.ToLower(consts.Expires)]; ok {
		if t, e := time.Parse(http.TimeFormat, exp); e == nil {
			obj.Expires = t.UTC()
		}
	}

	err := s.providers.GetStateStore().Put(getObjectKey(bucket, object), obj)
	if err != nil {
		return Object{}, err
	}
	return obj, nil
}

// GetObject Get object
func (s *service) GetObject(ctx context.Context, bucket, object string) (Object, io.ReadCloser, error) {
	var obj Object
	err := s.providers.GetStateStore().Get(getObjectKey(bucket, object), &obj)
	if errors.Is(err, providers.ErrStateStoreNotFound) {
		err = ErrObjectNotFound
		return Object{}, nil, err
	}

	reader, err := s.providers.GetFileStore().Cat(obj.Cid)
	if err != nil {
		return Object{}, nil, err
	}

	return obj, reader, nil
}

// GetObjectInfo Get object info
func (s *service) GetObjectInfo(ctx context.Context, bucket, object string) (Object, error) {
	var obj Object
	err := s.providers.GetStateStore().Get(getObjectKey(bucket, object), &obj)
	if errors.Is(err, providers.ErrStateStoreNotFound) {
		err = ErrObjectNotFound
		return Object{}, err
	}

	return obj, nil
}

// DeleteObject delete object
func (s *service) DeleteObject(ctx context.Context, bucket, object string) error {
	var obj Object
	err := s.providers.GetStateStore().Get(getObjectKey(bucket, object), &obj)
	if errors.Is(err, providers.ErrStateStoreNotFound) {
		err = ErrObjectNotFound
		return err
	}

	if err = s.providers.GetStateStore().Delete(getObjectKey(bucket, object)); err != nil {
		return err
	}

	//todo 是否先进性unpin，然后remove？
	if err := s.providers.GetFileStore().Remove(obj.Cid); err != nil {
		errMsg := fmt.Sprintf("mark Objet to delete error, bucket:%s, object:%s, cid:%s, error:%v \n", bucket, object, obj.Cid, err)
		return errors.New(errMsg)
	}
	return nil
}

func (s *service) CleanObjectsInBucket(ctx context.Context, bucket string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	prefixKey := fmt.Sprintf(allObjectPrefixFormat, bucket, "")
	err := s.providers.GetStateStore().Iterate(prefixKey, func(key, _ []byte) (stop bool, er error) {
		record := &Object{}
		er = s.providers.GetStateStore().Get(string(key), record)
		if er != nil {
			return
		}

		if err := s.DeleteObject(ctx, bucket, record.Name); err != nil {
			return
		}
		return
	})

	return err
}

// ListObjectsInfo - container for list objects.
type ListObjectsInfo struct {
	// Indicates whether the returned list objects response is truncated. A
	// value of true indicates that the list was truncated. The list can be truncated
	// if the number of objects exceeds the limit allowed or specified
	// by max keys.
	IsTruncated bool

	// When response is truncated (the IsTruncated element value in the response is true),
	// you can use the key name in this field as marker in the subsequent
	// request to get next set of objects.
	//
	// NOTE: AWS S3 returns NextMarker only if you have delimiter request parameter specified,
	NextMarker string

	// List of objects info for this request.
	Objects []Object

	// List of prefixes for this request.
	Prefixes []string
}

// ListObjects list user object
// TODO use more params
func (s *service) ListObjects(ctx context.Context, bucket string, prefix string, marker string, delimiter string, maxKeys int) (loi ListObjectsInfo, err error) {
	if maxKeys == 0 {
		return loi, nil
	}

	if len(prefix) > 0 && maxKeys == 1 && delimiter == "" && marker == "" {
		// Optimization for certain applications like
		// - Cohesity
		// - Actifio, Splunk etc.
		// which send ListObjects requests where the actual object
		// itself is the prefix and max-keys=1 in such scenarios
		// we can simply verify locally if such an object exists
		// to avoid the need for ListObjects().
		var obj Object
		err = s.providers.GetStateStore().Get(getObjectKey(bucket, prefix), &obj)
		if err == nil {
			loi.Objects = append(loi.Objects, obj)
			return loi, nil
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	seekKey := ""
	if marker != "" {
		seekKey = fmt.Sprintf(allObjectSeekKeyFormat, bucket, marker)
	}
	prefixKey := fmt.Sprintf(allObjectPrefixFormat, bucket, prefix)

	begin := false
	index := 0
	err = s.providers.GetStateStore().Iterate(prefixKey, func(key, _ []byte) (stop bool, er error) {
		record := &Object{}
		er = s.providers.GetStateStore().Get(string(key), record)
		if er != nil {
			return
		}
		if seekKey == string(key) {
			begin = true
		}

		if begin {
			loi.Objects = append(loi.Objects, *record)
			index++
		}

		if index == maxKeys {
			loi.IsTruncated = true
			begin = false
			return
		}

		return
	})

	if loi.IsTruncated {
		loi.NextMarker = loi.Objects[len(loi.Objects)-1].Name
	}

	return loi, nil
}

func (s *service) EmptyBucket(ctx context.Context, bucket string) (bool, error) {
	loi, err := s.ListObjects(ctx, bucket, "", "", "", 1)
	if err != nil {
		return false, err
	}
	return len(loi.Objects) == 0, nil
}

// ListObjectsV2Info - container for list objects version 2.
type ListObjectsV2Info struct {
	// Indicates whether the returned list objects response is truncated. A
	// value of true indicates that the list was truncated. The list can be truncated
	// if the number of objects exceeds the limit allowed or specified
	// by max keys.
	IsTruncated bool

	// When response is truncated (the IsTruncated element value in the response
	// is true), you can use the key name in this field as marker in the subsequent
	// request to get next set of objects.
	//
	// NOTE: This element is returned only if you have delimiter request parameter
	// specified.
	ContinuationToken     string
	NextContinuationToken string

	// List of objects info for this request.
	Objects []Object

	// List of prefixes for this request.
	Prefixes []string
}

// ListObjectsV2 list objects
func (s *service) ListObjectsV2(ctx context.Context, bucket string, prefix string, continuationToken string, delimiter string, maxKeys int, owner bool, startAfter string) (ListObjectsV2Info, error) {
	marker := continuationToken
	if marker == "" {
		marker = startAfter
	}
	loi, err := s.ListObjects(ctx, bucket, prefix, marker, delimiter, maxKeys)
	if err != nil {
		return ListObjectsV2Info{}, err
	}
	listV2Info := ListObjectsV2Info{
		IsTruncated:           loi.IsTruncated,
		ContinuationToken:     continuationToken,
		NextContinuationToken: loi.NextMarker,
		Objects:               loi.Objects,
		Prefixes:              loi.Prefixes,
	}
	return listV2Info, nil
}

/*---------------------------------------------------*/

func (s *service) CreateMultipartUpload(ctx context.Context, bucname string, objname string, meta map[string]string) (mtp Multipart, err error) {
	uploadId := uuid.NewString()
	mtp = Multipart{
		Bucket:    bucname,
		Object:    objname,
		UploadID:  uploadId,
		MetaData:  meta,
		Initiated: time.Now().UTC(),
	}

	err = s.providers.GetStateStore().Put(getUploadKey(bucname, objname, uploadId), mtp)
	if err != nil {
		return
	}

	return
}

func (s *service) UploadPart(ctx context.Context, bucname string, objname string, uploadID string, partID int, reader *hash.Reader, size int64, meta map[string]string) (part ObjectPart, err error) {
	cid, err := s.providers.GetFileStore().Store(reader)
	if err != nil {
		return
	}

	part = ObjectPart{
		Number:  partID,
		ETag:    reader.ETag().String(),
		Cid:     cid,
		Size:    size,
		ModTime: time.Now().UTC(),
	}

	mtp, err := s.getMultipart(ctx, bucname, objname, uploadID)
	if err != nil {
		return
	}

	mtp.Parts = append(mtp.Parts, part)
	err = s.providers.GetStateStore().Put(getUploadKey(bucname, objname, uploadID), mtp)
	if err != nil {
		return part, err
	}

	return
}

func (s *service) AbortMultipartUpload(ctx context.Context, bucname string, objname string, uploadID string) (err error) {
	mtp, err := s.getMultipart(ctx, bucname, objname, uploadID)
	if err != nil {
		return
	}

	for _, part := range mtp.Parts {
		err = s.providers.GetFileStore().Remove(part.Cid)
		if err != nil {
			return
		}
	}

	err = s.removeMultipart(ctx, bucname, objname, uploadID)
	if err != nil {
		return
	}

	return
}

func (s *service) CompleteMultiPartUpload(ctx context.Context, bucname string, objname string, uploadID string, parts []CompletePart) (obj Object, err error) {
	mi, err := s.getMultipart(ctx, bucname, objname, uploadID)
	if err != nil {
		return
	}

	var (
		readers    []io.Reader
		objectSize int64
	)

	defer func() {
		for _, rdr := range readers {
			_ = rdr.(io.ReadCloser).Close()
		}
	}()

	idxMap := objectPartIndexMap(mi.Parts)
	for i, part := range parts {
		partIndex, ok := idxMap[part.PartNumber]
		if !ok {
			err = s3utils.InvalidPart{
				PartNumber: part.PartNumber,
				GotETag:    part.ETag,
			}
			return
		}

		gotPart := mi.Parts[partIndex]

		part.ETag = canonicalizeETag(part.ETag)
		if gotPart.ETag != part.ETag {
			err = s3utils.InvalidPart{
				PartNumber: part.PartNumber,
				ExpETag:    gotPart.ETag,
				GotETag:    part.ETag,
			}
			return
		}

		// All parts except the last part has to be at least 5MB.
		if (i < len(parts)-1) && !(gotPart.Size >= consts.MinPartSize) {
			err = s3utils.PartTooSmall{
				PartNumber: part.PartNumber,
				PartSize:   gotPart.Size,
				PartETag:   part.ETag,
			}
			return
		}

		// Save for total objname size.
		objectSize += gotPart.Size

		var rdr io.ReadCloser
		rdr, err = s.providers.GetFileStore().Cat(gotPart.Cid)
		if err != nil {
			return
		}

		readers = append(readers, rdr)
	}

	cid, err := s.providers.GetFileStore().Store(io.MultiReader(readers...))
	if err != nil {
		return
	}

	obj = Object{
		Bucket:           bucname,
		Name:             objname,
		ModTime:          time.Now().UTC(),
		Size:             objectSize,
		IsDir:            false,
		ETag:             computeCompleteMultipartMD5(parts),
		Cid:              cid,
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		ContentType:      mi.MetaData[strings.ToLower(consts.ContentType)],
		ContentEncoding:  mi.MetaData[strings.ToLower(consts.ContentEncoding)],
		SuccessorModTime: time.Now().UTC(),
	}

	if exp, ok := mi.MetaData[strings.ToLower(consts.Expires)]; ok {
		if t, e := time.Parse(http.TimeFormat, exp); e == nil {
			obj.Expires = t.UTC()
		}
	}

	err = s.providers.GetStateStore().Put(getObjectKey(bucname, objname), obj)
	if err != nil {
		return
	}

	err = s.removeMultipartInfo(ctx, bucname, objname, uploadID)
	if err != nil {
		return
	}

	return
}

func (s *service) GetMultipart(ctx context.Context, bucname string, objname string, uploadID string) (mtp Multipart, err error) {
	return s.getMultipart(ctx, bucname, objname, uploadID)
}

func (s *service) getMultipart(ctx context.Context, bucname string, objname string, uploadID string) (mtp Multipart, err error) {
	err = s.providers.GetStateStore().Get(getUploadKey(bucname, objname, uploadID), &mtp)
	if errors.Is(err, providers.ErrStateStoreNotFound) {
		err = ErrUploadNotFound
		return
	}
	return
}

func (s *service) removeMultipart(ctx context.Context, bucname string, objname string, uploadID string) (err error) {
	err = s.providers.GetStateStore().Delete(getUploadKey(bucname, objname, uploadID))
	if errors.Is(err, providers.ErrStateStoreNotFound) {
		err = ErrUploadNotFound
		return
	}
	return
}

func (s *service) removeMultipartInfo(ctx context.Context, bucname string, objname string, uploadID string) (err error) {
	err = s.providers.GetStateStore().Delete(getUploadKey(bucname, objname, uploadID))
	if errors.Is(err, providers.ErrStateStoreNotFound) {
		err = ErrUploadNotFound
		return
	}
	return
}

func objectPartIndexMap(parts []ObjectPart) map[int]int {
	mp := make(map[int]int)
	for i, part := range parts {
		mp[part.Number] = i
	}
	return mp
}

// canonicalizeETag returns ETag with leading and trailing double-quotes removed,
// if any present
func canonicalizeETag(etag string) string {
	return etagRegex.ReplaceAllString(etag, "$1")
}

func computeCompleteMultipartMD5(parts []CompletePart) string {
	var finalMD5Bytes []byte
	for _, part := range parts {
		md5Bytes, err := hex.DecodeString(canonicalizeETag(part.ETag))
		if err != nil {
			finalMD5Bytes = append(finalMD5Bytes, []byte(part.ETag)...)
		} else {
			finalMD5Bytes = append(finalMD5Bytes, md5Bytes...)
		}
	}
	s3MD5 := fmt.Sprintf("%s-%d", etag.Multipart(finalMD5Bytes), len(parts))
	return s3MD5
}
