package auth

import (
	"net/http"
	"net/url"
	"strings"
)

// IsRequestSignatureV4 Verify if request has AWS Signature Version '4'.
func IsRequestSignatureV4(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Authorization"), signV4Algorithm)
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
			//log.Infof("parse query failed, query: %s, error: %v", r.URL.RawQuery, err)
			return AuthTypeUnknown
		}
	}
	if IsRequestSignatureV4(r) {
		return AuthTypeSigned
	} else if isRequestPresignedSignatureV4(r) {
		return AuthTypePresigned
	}
	return AuthTypeUnknown
}
