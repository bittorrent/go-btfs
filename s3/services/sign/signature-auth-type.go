package sign

import (
	"github.com/bittorrent/go-btfs/s3/consts"
	"net/http"
	"net/url"
	"strings"
)

// Verify if request has JWT.
func isRequestJWT(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Authorization"), "Bearer")
}

// IsRequestSignatureV4 Verify if request has AWS Signature Version '4'.
func IsRequestSignatureV4(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Authorization"), signV4Algorithm)
}

// Verify if request has AWS Signature Version '2'.
func isRequestSignatureV2(r *http.Request) bool {
	return !strings.HasPrefix(r.Header.Get("Authorization"), signV4Algorithm) &&
		strings.HasPrefix(r.Header.Get("Authorization"), signV2Algorithm)
}

// Verify if request has AWS PreSign Version '4'.
func isRequestPresignedSignatureV4(r *http.Request) bool {
	_, ok := r.URL.Query()["X-Amz-Credential"]
	return ok
}

// Verify request has AWS PreSign Version '2'.
func isRequestPresignedSignatureV2(r *http.Request) bool {
	_, ok := r.URL.Query()["AWSAccessKeyId"]
	return ok
}

// Verify if request has AWS Post policy Signature Version '4'.
func isRequestPostPolicySignatureV4(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") &&
		r.Method == http.MethodPost
}

// Verify if the request has AWS Streaming Signature Version '4'. This is only valid for 'PUT' operation.
func isRequestSignStreamingV4(r *http.Request) bool {
	return r.Header.Get("x-amz-content-sha256") == consts.StreamingContentSHA256 &&
		r.Method == http.MethodPut
}

// AuthType Authorization type.
type AuthType int

// List of all supported auth types.
const (
	AuthTypeUnknown AuthType = iota
	AuthTypeAnonymous
	AuthTypePresigned
	AuthTypePresignedV2
	AuthTypePostPolicy
	AuthTypeStreamingSigned
	AuthTypeSigned
	AuthTypeSignedV2
	AuthTypeJWT
	AuthTypeSTS
)

// GetRequestAuthType Get request authentication type.
func GetRequestAuthType(r *http.Request) AuthType {
	if r.URL != nil {
		var err error
		r.Form, err = url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			return AuthTypeUnknown
		}
	}
	if isRequestSignatureV2(r) {
		return AuthTypeSignedV2
	} else if isRequestPresignedSignatureV2(r) {
		return AuthTypePresignedV2
	} else if isRequestSignStreamingV4(r) {
		return AuthTypeStreamingSigned
	} else if IsRequestSignatureV4(r) {
		return AuthTypeSigned
	} else if isRequestPresignedSignatureV4(r) {
		return AuthTypePresigned
	} else if isRequestJWT(r) {
		return AuthTypeJWT
	} else if isRequestPostPolicySignatureV4(r) {
		return AuthTypePostPolicy
	} else if _, ok := r.Form[consts.StsAction]; ok {
		return AuthTypeSTS
	} else if _, ok := r.Header[consts.Authorization]; !ok {
		return AuthTypeAnonymous
	}
	return AuthTypeUnknown
}
