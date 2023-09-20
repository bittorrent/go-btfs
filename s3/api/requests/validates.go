package requests

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bittorrent/go-btfs/s3/api/services/object"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/etag"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	validBucketName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9\.\-]{1,61}[A-Za-z0-9]$`)
	ipAddress       = regexp.MustCompile(`^(\d+\.){3}\d+$`)
)

func ValidateBucketName(bucketName *string) (val string, err error) {
	if bucketName == nil {
		return
	}
	if *bucketName == "" {
		err = ErrBucketNameInvalid
		return
	}
	val = *bucketName
	if ipAddress.MatchString(val) ||
		!validBucketName.MatchString(val) ||
		strings.Contains(val, "..") ||
		strings.Contains(val, ".-") ||
		strings.Contains(val, "-.") {
		err = ErrBucketNameInvalid
	}
	return
}

func ValidateBucketACL(acl *string) (val string, err error) {
	if acl == nil || *acl == "" {
		val = consts.DefaultBucketACL
	} else {
		val = *acl
	}
	if !consts.SupportedBucketACLs[val] {
		err = ErrACLUnsupported
	}
	return
}

func ValidateCreateBucketConfiguration(configuration *s3.CreateBucketConfiguration) (val string, err error) {
	if configuration == nil || configuration.LocationConstraint == nil || *configuration.LocationConstraint == "" {
		val = consts.DefaultBucketRegion
	}
	if !consts.SupportedBucketRegions[val] {
		err = ErrRegionUnsupported
	}
	return
}

func ValidateObjectName(objectName *string) (val string, err error) {
	if objectName == nil {
		return
	}
	if *objectName == "" {
		err = ErrObjectNameInvalid
		return
	}
	val, err = url.PathUnescape(*objectName)
	if err != nil {
		err = ErrObjectNameInvalid
		return
	}
	if len(val) > 1024 {
		err = ErrObjectNameTooLong
		return
	}
	if strings.HasPrefix(val, "/") {
		err = ErrObjectNamePrefixSlash
		return
	}
	if !utf8.ValidString(val) || strings.Contains(val, `//`) {
		err = ErrObjectNameInvalid
	}
	for _, p := range strings.Split(val, "/") {
		switch strings.TrimSpace(p) {
		case "..", ".":
			err = ErrObjectNameInvalid
			return
		}
	}
	return
}

func ValidateContentMD5(contentMD5 *string) (val string, err error) {
	if contentMD5 == nil {
		return
	}
	if *contentMD5 == "" {
		err = ErrContentMd5Invalid
		return
	}
	b, err := base64.StdEncoding.Strict().DecodeString(*contentMD5)
	if err != nil || len(b) != md5.Size {
		err = ErrContentMd5Invalid
		return
	}
	val = etag.ETag(b).String()
	return
}

func ValidateChecksumSHA256(checksumSHA256 *string) (val string, err error) {
	if checksumSHA256 == nil || *checksumSHA256 == "" {
		return
	}
	if *checksumSHA256 == consts.UnsignedSHA256 {
		return
	}
	b, err := hex.DecodeString(*checksumSHA256)
	if err != nil || len(b) == 0 {
		err = ErrChecksumSha256Invalid
		return
	}
	val = hex.EncodeToString(b)
	return
}

func ValidateContentLength(contentLength *int64, max int64) (val int64, err error) {
	if contentLength == nil {
		return
	}
	if *contentLength == -1 {
		err = ErrContentLengthMissing
		return
	}
	if *contentLength < 1 {
		err = ErrContentLengthTooSmall
		return
	}
	if *contentLength > max {
		err = ErrContentLengthTooLarge
		return
	}
	val = *contentLength
	return
}

func ValidateContentType(contentType *string) (val string, err error) {
	if contentType == nil || *contentType == "" {
		val = consts.DefaultContentType
		return
	}
	val = *contentType
	return
}

func ValidateContentEncoding(contentEncoding *string) (val string, err error) {
	if contentEncoding == nil || *contentEncoding == "" {
		return
	}
	encs := make([]string, 0)
	for _, enc := range strings.Split(*contentEncoding, ",") {
		if enc != consts.StreamingContentEncoding {
			encs = append(encs, enc)
		}
	}
	val = strings.Join(encs, ",")
	return
}

func ValidateExpires(expires *time.Time) (val time.Time, err error) {
	if expires == nil {
		return
	}
	val = *expires
	return
}

func ValidateCopySource(copySource *string) (val1, val2 string, err error) {
	if copySource == nil {
		return
	}
	src, err := url.PathUnescape(*copySource)
	if err != nil {
		err = ErrCopySrcInvalid
		return
	}
	src = strings.TrimPrefix(src, consts.SlashSeparator)
	idx := strings.Index(src, consts.SlashSeparator)
	if idx < 0 {
		err = ErrCopySrcInvalid
		return
	}
	val1 = src[:idx]
	val2 = src[idx+len(consts.SlashSeparator):]
	if val1 == "" || val2 == "" {
		err = ErrCopySrcInvalid
		return
	}
	val1, err = ValidateBucketName(&val1)
	if err != nil {
		return
	}
	val2, err = ValidateObjectName(&val2)
	return
}

func ValidateMetadataDirective(metadataDirective *string) (val bool, err error) {
	if metadataDirective == nil {
		return
	}
	if *metadataDirective == "REPLACE" {
		val = true
	}
	return
}

func ValidateObjectsDelete(delete *s3.Delete) (vals []*object.ToDeleteObject, quite bool, err error) {
	if delete == nil {
		err = ErrFailedDecodeXML{errors.New("delete is nil")}
		return
	}
	if len(delete.Objects) < 1 || len(delete.Objects) > consts.MaxDeleteList {
		err = ErrDeletesCountInvalid
		return
	}
	if delete.Quiet != nil && *delete.Quiet == true {
		quite = true
	}
	for _, obj := range delete.Objects {
		deleteObj := &object.ToDeleteObject{}
		deleteObj.Object, deleteObj.ValidateErr = ValidateObjectName(obj.Key)
		vals = append(vals, deleteObj)
	}
	return
}

func ValidateMaxKeys(maxKeys *int64) (val int64, err error) {
	if maxKeys == nil || *maxKeys > consts.MaxObjectList {
		val = consts.MaxObjectList
		return
	}
	if *maxKeys < 0 {
		err = ErrMaxKeysInvalid
		return
	}
	val = *maxKeys
	return
}

func ValidateMarker(marker *string) (val string, err error) {
	if marker == nil || *marker == "" {
		return
	}
	val, err = ValidateObjectName(marker)
	if err != nil {
		err = ErrMarkerInvalid
	}
	return
}

func ValidatePrefix(prefix *string) (val string, err error) {
	if prefix == nil || *prefix == "" {
		return
	}
	val, err = ValidateObjectName(prefix)
	if err != nil {
		err = ErrPrefixInvalid
	}
	return
}

func ValidateMarkerAndPrefixCombination(marker, prefix string) (err error) {
	if marker != "" && !strings.HasPrefix(marker, prefix) {
		err = ErrMarkerPrefixCombinationInvalid
	}
	return
}

func ValidateDelimiter(delimiter *string) (val string, err error) {
	if delimiter == nil {
		return
	}
	val = *delimiter
	return
}

func ValidateEncodingType(encodingType *string) (val string, err error) {
	if encodingType == nil || *encodingType == "" {
		return
	}
	if !strings.EqualFold(*encodingType, consts.DefaultEncodingType) {
		err = ErrEncodingTypeInvalid
		return
	}
	val = consts.DefaultEncodingType
	return
}

func ValidateContinuationToken(continuationToken *string) (val string, err error) {
	if continuationToken == nil || *continuationToken == "" {
		return
	}
	token, err := base64.StdEncoding.DecodeString(*continuationToken)
	if err != nil {
		err = ErrContinuationTokenInvalid
		return
	}
	tokenStr := string(token)
	val, err = ValidateObjectName(&tokenStr)
	if err != nil {
		err = ErrContinuationTokenInvalid
	}
	return
}

func ValidateStartAfter(startAfter *string) (val string, err error) {
	if startAfter == nil || *startAfter == "" {
		return
	}
	val, err = ValidateObjectName(startAfter)
	if err != nil {
		err = ErrStartAfterInvalid
	}
	return
}

func ValidateFetchOwner(fetchOwner *bool) (val bool, err error) {
	if fetchOwner == nil {
		return
	}
	val = *fetchOwner
	return
}

func ValidateUploadId(uploadId *string) (val string, err error) {
	if uploadId == nil {
		return
	}
	val = *uploadId
	return
}

func ValidatePartNumber(partNumber *int64) (val int64, err error) {
	if partNumber == nil {
		return
	}
	if *partNumber < consts.MinPartNumber || *partNumber > consts.MaxPartNumber {
		err = ErrPartNumberInvalid
		return
	}
	val = *partNumber
	return
}

func ValidateCompletedMultipartUpload(upload *s3.CompletedMultipartUpload) (val object.CompletedParts, err error) {
	if upload == nil {
		err = ErrFailedDecodeXML{errors.New("complete-upload is nil")}
		return
	}
	if len(upload.Parts) < 1 || len(upload.Parts) > consts.MaxPartNumber {
		err = ErrPartsCountInvalid
		return
	}
	for _, part := range upload.Parts {
		if part.PartNumber == nil || part.ETag == nil {
			err = ErrPartInvalid
			return
		}

		opart := &object.CompletePart{
			ETag: *part.ETag,
		}
		opart.PartNumber, err = ValidatePartNumber(part.PartNumber)
		if err != nil {
			return
		}
		val = append(val, opart)
	}
	if !sort.IsSorted(val) {
		err = ErrPartOrderInvalid
	}
	return
}
