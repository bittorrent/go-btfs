package handlers

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/consts"
	"net/http"
	"net/textproto"
	"strings"
)

const streamingContentEncoding = "aws-chunked"

// errInvalidArgument means that input argument is invalid.
var errInvalidArgument = errors.New("Invalid arguments specified")

// Supported headers that needs to be extracted.
var supportedHeaders = []string{
	consts.ContentType,
	consts.CacheControl,
	consts.ContentLength,
	consts.ContentEncoding,
	consts.ContentDisposition,
	consts.AmzStorageClass,
	consts.AmzObjectTagging,
	consts.Expires,
	consts.AmzBucketReplicationStatus,
	// Add more supported headers here.
}

// userMetadataKeyPrefixes contains the prefixes of used-defined metadata keys.
// All values stored with a key starting with one of the following prefixes
// must be extracted from the header.
var userMetadataKeyPrefixes = []string{
	"x-amz-meta-",
}

// matches k1 with all keys, returns 'true' if one of them matches
func equals(k1 string, keys ...string) bool {
	for _, k2 := range keys {
		if strings.EqualFold(k1, k2) {
			return true
		}
	}
	return false
}

// extractMetadata extracts metadata from HTTP header and HTTP queryString.
// Note: The key has been converted to lowercase letters
func extractMetadata(ctx context.Context, r *http.Request) (metadata map[string]string, err error) {
	query := r.Form
	header := r.Header
	metadata = make(map[string]string)
	// Extract all query values.
	err = extractMetadataFromMime(ctx, textproto.MIMEHeader(query), metadata)
	if err != nil {
		return nil, err
	}

	// Extract all header values.
	err = extractMetadataFromMime(ctx, textproto.MIMEHeader(header), metadata)
	if err != nil {
		return nil, err
	}

	// Set content-type to default value if it is not set.
	if _, ok := metadata[strings.ToLower(consts.ContentType)]; !ok {
		metadata[strings.ToLower(consts.ContentType)] = "binary/octet-stream"
	}

	// https://github.com/google/security-research/security/advisories/GHSA-76wf-9vgp-pj7w
	for k := range metadata {
		if equals(k, consts.AmzMetaUnencryptedContentLength, consts.AmzMetaUnencryptedContentMD5) {
			delete(metadata, k)
		}
	}

	if contentEncoding, ok := metadata[strings.ToLower(consts.ContentEncoding)]; ok {
		contentEncoding = trimAwsChunkedContentEncoding(contentEncoding)
		if contentEncoding != "" {
			// Make sure to trim and save the content-encoding
			// parameter for a streaming signature which is set
			// to a custom value for example: "aws-chunked,gzip".
			metadata[strings.ToLower(consts.ContentEncoding)] = contentEncoding
		} else {
			// Trimmed content encoding is empty when the header
			// value is set to "aws-chunked" only.

			// Make sure to delete the content-encoding parameter
			// for a streaming signature which is set to value
			// for example: "aws-chunked"
			delete(metadata, strings.ToLower(consts.ContentEncoding))
		}
	}

	// Success.
	return metadata, nil
}

// extractMetadata extracts metadata from map values.
func extractMetadataFromMime(ctx context.Context, v textproto.MIMEHeader, m map[string]string) error {
	if v == nil {
		return errInvalidArgument
	}

	nv := make(textproto.MIMEHeader, len(v))
	for k, kv := range v {
		// Canonicalize all headers, to remove any duplicates.
		nv[strings.ToLower(k)] = kv
	}

	// Save all supported headers.
	for _, supportedHeader := range supportedHeaders {
		value, ok := nv[strings.ToLower(supportedHeader)]
		if ok {
			m[strings.ToLower(supportedHeader)] = strings.Join(value, ",")
		}
	}

	for key := range v {
		lowerKey := strings.ToLower(key)
		for _, prefix := range userMetadataKeyPrefixes {
			if !strings.HasPrefix(lowerKey, strings.ToLower(prefix)) {
				continue
			}
			value, ok := nv[lowerKey]
			if ok {
				m[lowerKey] = strings.Join(value, ",")
				break
			}
		}
	}
	return nil
}

func trimAwsChunkedContentEncoding(contentEnc string) (trimmedContentEnc string) {
	if contentEnc == "" {
		return contentEnc
	}
	var newEncs []string
	for _, enc := range strings.Split(contentEnc, ",") {
		if enc != streamingContentEncoding {
			newEncs = append(newEncs, enc)
		}
	}
	return strings.Join(newEncs, ",")
}
