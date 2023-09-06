package etag

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ETag is a single S3 ETag.
//
// An S3 ETag sometimes corresponds to the MD5 of
// the S3 object content. However, when an object
// is encrypted, compressed or uploaded using
// the S3 multipart API then its ETag is not
// necessarily the MD5 of the object content.
//
// For a more detailed description of S3 ETags
// take a look at the package documentation.
type ETag []byte

// String returns the string representation of the ETag.
//
// The returned string is a hex representation of the
// binary ETag with an optional '-<part-number>' suffix.
func (e ETag) String() string {
	if e.IsMultipart() {
		return hex.EncodeToString(e[:16]) + string(e[16:])
	}
	return hex.EncodeToString(e)
}

// IsEncrypted reports whether the ETag is encrypted.
func (e ETag) IsEncrypted() bool {
	return len(e) > 16 && !bytes.ContainsRune(e, '-')
}

// IsMultipart reports whether the ETag belongs to an
// object that has been uploaded using the S3 multipart
// API.
// An S3 multipart ETag has a -<part-number> suffix.
func (e ETag) IsMultipart() bool {
	return len(e) > 16 && bytes.ContainsRune(e, '-')
}

// Parts returns the number of object parts that are
// referenced by this ETag. It returns 1 if the object
// has been uploaded using the S3 singlepart API.
//
// Parts may panic if the ETag is an invalid multipart
// ETag.
func (e ETag) Parts() int {
	if !e.IsMultipart() {
		return 1
	}

	n := bytes.IndexRune(e, '-')
	parts, err := strconv.Atoi(string(e[n+1:]))
	if err != nil {
		panic(err) // malformed ETag
	}
	return parts
}

var _ Tagger = ETag{} // compiler check

// ETag returns the ETag itself.
//
// By providing this method ETag implements
// the Tagger interface.
func (e ETag) ETag() ETag { return e }

// FromContentMD5 decodes and returns the Content-MD5
// as ETag, if set. If no Content-MD5 header is set
// it returns an empty ETag and no error.
func FromContentMD5(h http.Header) (ETag, error) {
	v, ok := h["Content-Md5"]
	if !ok {
		return nil, nil
	}
	if v[0] == "" {
		return nil, errors.New("etag: content-md5 is set but contains no value")
	}
	b, err := base64.StdEncoding.Strict().DecodeString(v[0])
	if err != nil {
		return nil, err
	}
	if len(b) != md5.Size {
		return nil, errors.New("etag: invalid content-md5")
	}
	return ETag(b), nil
}

// Multipart computes an S3 multipart ETag given a list of
// S3 singlepart ETags. It returns nil if the list of
// ETags is empty.
//
// Any encrypted or multipart ETag will be ignored and not
// used to compute the returned ETag.
func Multipart(etags ...ETag) ETag {
	if len(etags) == 0 {
		return nil
	}

	var n int64
	h := md5.New()
	for _, etag := range etags {
		if !etag.IsMultipart() && !etag.IsEncrypted() {
			h.Write(etag)
			n++
		}
	}
	etag := append(h.Sum(nil), '-')
	return strconv.AppendInt(etag, n, 10)
}

// Equal returns true if and only if the two ETags are
// identical.
func Equal(a, b ETag) bool { return bytes.Equal(a, b) }

// Parse parses s as an S3 ETag, returning the result.
// The string can be an encrypted, singlepart
// or multipart S3 ETag. It returns an error if s is
// not a valid textual representation of an ETag.
func Parse(s string) (ETag, error) {
	const strict = false
	return parse(s, strict)
}

// parse parse s as an S3 ETag, returning the result.
// It operates in one of two modes:
//   - strict
//   - non-strict
//
// In strict mode, parse only accepts ETags that
// are AWS S3 compatible. In particular, an AWS
// S3 ETag always consists of a 128 bit checksum
// value and an optional -<part-number> suffix.
// Therefore, s must have the following form in
// strict mode:  <32-hex-characters>[-<integer>]
//
// In non-strict mode, parse also accepts ETags
// that are not AWS S3 compatible - e.g. encrypted
// ETags.
func parse(s string, strict bool) (ETag, error) {
	// An S3 ETag may be a double-quoted string.
	// Therefore, we remove double quotes at the
	// start and end, if any.
	if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
		s = s[1 : len(s)-1]
	}

	// An S3 ETag may be a multipart ETag that
	// contains a '-' followed by a number.
	// If the ETag does not a '-' is is either
	// a singlepart or encrypted ETag.
	n := strings.IndexRune(s, '-')
	if n == -1 {
		etag, err := hex.DecodeString(s)
		if err != nil {
			return nil, err
		}
		if strict && len(etag) != 16 { // AWS S3 ETags are always 128 bit long
			return nil, fmt.Errorf("etag: invalid length %d", len(etag))
		}
		return ETag(etag), nil
	}

	prefix, suffix := s[:n], s[n:]
	if len(prefix) != 32 {
		return nil, fmt.Errorf("etag: invalid prefix length %d", len(prefix))
	}
	if len(suffix) <= 1 {
		return nil, errors.New("etag: suffix is not a part number")
	}

	etag, err := hex.DecodeString(prefix)
	if err != nil {
		return nil, err
	}
	partNumber, err := strconv.Atoi(suffix[1:]) // suffix[0] == '-' Therefore, we start parsing at suffix[1]
	if err != nil {
		return nil, err
	}
	if strict && (partNumber == 0 || partNumber > 10000) {
		return nil, fmt.Errorf("etag: invalid part number %d", partNumber)
	}
	return ETag(append(etag, suffix...)), nil
}
