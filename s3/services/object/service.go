package object

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/etag"
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
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
