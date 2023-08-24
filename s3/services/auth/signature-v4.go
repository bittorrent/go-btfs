/*
 * The following code tries to reverse engineer the Amazon S3 APIs,
 * and is mostly copied from minio implementation.
 */

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package auth

import (
	"crypto/subtle"
	"errors"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/set"
	"github.com/bittorrent/go-btfs/s3/utils"
)

// AWS Signature Version '4' constants.
const (
	signV4Algorithm = "AWS4-HMAC-SHA256"
	iso8601Format   = "20060102T150405Z"
	yyyymmdd        = "20060102"
)

type serviceType string

const (
	ServiceS3 serviceType = "s3"
	////ServiceSTS STS
	//ServiceSTS serviceType = "sts"
)

// compareSignatureV4 returns true if and only if both signatures
// are equal. The signatures are expected to be HEX encoded strings
// according to the AWS S3 signature V4 spec.
func compareSignatureV4(sig1, sig2 string) bool {
	// The CTC using []byte(str) works because the hex encoding
	// is unique for a sequence of bytes. See also compareSignatureV2.
	return subtle.ConstantTimeCompare([]byte(sig1), []byte(sig2)) == 1
}

// DoesPresignedSignatureMatch - Verify queryString headers with presigned signature
//   - http://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-query-string-auth.html
//
// returns handlers.ErrcodeNone if the signature matches.
func (s *service) doesPresignedSignatureMatch(hashedPayload string, r *http.Request, region string, stype serviceType) (ack *accesskey.AccessKey, err error) {
	// Copy request
	req := *r

	// Parse request query string.
	pSignValues, err := parsePreSignV4(req.Form, region, stype)
	if err != nil {
		return
	}

	// Check accesskey
	ack, err = s.accessKeySvc.Get(pSignValues.Credential.accessKey)
	if errors.Is(err, accesskey.ErrNotFound) {
		err = responses.ErrInvalidAccessKeyID
	}
	if err != nil {
		return
	}
	if !ack.Enable {
		err = responses.ErrAccessKeyDisabled
		return
	}

	// Extract all the signed headers along with its values.
	extractedSignedHeaders, err := extractSignedHeaders(pSignValues.SignedHeaders, r)
	if err != nil {
		return
	}

	// If the host which signed the request is slightly ahead in time (by less than MaxSkewTime) the
	// request should still be allowed.
	if pSignValues.Date.After(time.Now().UTC().Add(consts.MaxSkewTime)) {
		err = responses.ErrRequestNotReadyYet
		return
	}

	if time.Now().UTC().Sub(pSignValues.Date) > pSignValues.Expires {
		err = responses.ErrExpiredPresignRequest
		return
	}

	// Save the date and expires.
	t := pSignValues.Date
	expireSeconds := int(pSignValues.Expires / time.Second)

	// Construct new query.
	query := make(url.Values)
	clntHashedPayload := req.Form.Get(consts.AmzContentSha256)
	if clntHashedPayload != "" {
		query.Set(consts.AmzContentSha256, hashedPayload)
	}

	// not check token?
	//token := req.Form.Get(consts.AmzSecurityToken)
	//if token != "" {
	//	query.Set(consts.AmzSecurityToken, cred.SessionToken)
	//}

	query.Set(consts.AmzAlgorithm, signV4Algorithm)

	// Construct the query.
	query.Set(consts.AmzDate, t.Format(iso8601Format))
	query.Set(consts.AmzExpires, strconv.Itoa(expireSeconds))
	query.Set(consts.AmzSignedHeaders, utils.GetSignedHeaders(extractedSignedHeaders))
	query.Set(consts.AmzCredential, ack.Key+consts.SlashSeparator+pSignValues.Credential.getScope())

	defaultSigParams := set.CreateStringSet(
		consts.AmzContentSha256,
		//consts.AmzSecurityToken,
		consts.AmzAlgorithm,
		consts.AmzDate,
		consts.AmzExpires,
		consts.AmzSignedHeaders,
		consts.AmzCredential,
		consts.AmzSignature,
	)

	// Add missing query parameters if any provided in the request URL
	for k, v := range req.Form {
		if !defaultSigParams.Contains(k) {
			query[k] = v
		}
	}

	// Get the encoded query.
	encodedQuery := query.Encode()

	// Verify if date query is same.
	if req.Form.Get(consts.AmzDate) != query.Get(consts.AmzDate) {
		err = responses.ErrSignatureDoesNotMatch
	}
	// Verify if expires query is same.
	if req.Form.Get(consts.AmzExpires) != query.Get(consts.AmzExpires) {
		err = responses.ErrSignatureDoesNotMatch
		return
	}
	// Verify if signed headers query is same.
	if req.Form.Get(consts.AmzSignedHeaders) != query.Get(consts.AmzSignedHeaders) {
		err = responses.ErrSignatureDoesNotMatch
		return
	}
	// Verify if credential query is same.
	if req.Form.Get(consts.AmzCredential) != query.Get(consts.AmzCredential) {
		err = responses.ErrSignatureDoesNotMatch
		return
	}
	// Verify if sha256 payload query is same.
	if clntHashedPayload != "" && clntHashedPayload != query.Get(consts.AmzContentSha256) {
		err = responses.ErrContentSHA256Mismatch
		return
	}
	// not check SessionToken.
	//// Verify if security token is correct.
	//if token != "" && subtle.ConstantTimeCompare([]byte(token), []byte(cred.SessionToken)) != 1 {
	//	return handlers.ErrInvalidToken
	//}

	// Verify finally if signature is same.

	// Get canonical request.
	presignedCanonicalReq := utils.GetCanonicalRequest(extractedSignedHeaders, hashedPayload, encodedQuery, req.URL.Path, req.Method)

	// Get string to sign from canonical request.
	presignedStringToSign := utils.GetStringToSign(presignedCanonicalReq, t, pSignValues.Credential.getScope())

	// Get hmac presigned signing key.
	presignedSigningKey := utils.GetSigningKey(ack.Secret, pSignValues.Credential.scope.date,
		pSignValues.Credential.scope.region, string(stype))

	// Get new signature.
	newSignature := utils.GetSignature(presignedSigningKey, presignedStringToSign)

	// Verify signature.
	if !compareSignatureV4(req.Form.Get(consts.AmzSignature), newSignature) {
		err = responses.ErrSignatureDoesNotMatch
		return
	}

	return
}

// DoesSignatureMatch - Verify authorization header with calculated header in accordance with
//   - http://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-authenticating-requests.html
func (s *service) doesSignatureMatch(hashedPayload string, r *http.Request, region string, stype serviceType) (ack *accesskey.AccessKey, err error) {
	// Copy request.
	req := *r

	// Save authorization header.
	v4Auth := req.Header.Get(consts.Authorization)

	// Parse signature version '4' header.
	signV4Values, err := parseSignV4(v4Auth, region, stype)
	if err != nil {
		return
	}

	// Extract all the signed headers along with its values.
	extractedSignedHeaders, err := extractSignedHeaders(signV4Values.SignedHeaders, r)
	if err != nil {
		return
	}

	// Check accesskey
	ack, err = s.accessKeySvc.Get(signV4Values.Credential.accessKey)
	if errors.Is(err, accesskey.ErrNotFound) {
		err = responses.ErrInvalidAccessKeyID
	}
	if err != nil {
		return
	}
	if !ack.Enable {
		err = responses.ErrAccessKeyDisabled
		return
	}

	// Extract date, if not present throw error.
	var date string
	if date = req.Header.Get(consts.AmzDate); date == "" {
		if date = r.Header.Get(consts.Date); date == "" {
			err = responses.ErrMissingDateHeader
			return
		}
	}

	// Parse date header.
	t, err := time.Parse(iso8601Format, date)
	if err != nil {
		err = responses.ErrAuthorizationHeaderMalformed
		return
	}

	// Query string.
	queryStr := req.URL.Query().Encode()

	// Get canonical request.
	canonicalRequest := utils.GetCanonicalRequest(extractedSignedHeaders, hashedPayload, queryStr, req.URL.Path, req.Method)

	// Get string to sign from canonical request.
	stringToSign := utils.GetStringToSign(canonicalRequest, t, signV4Values.Credential.getScope())

	// Get hmac signing key.
	signingKey := utils.GetSigningKey(ack.Secret, signV4Values.Credential.scope.date,
		signV4Values.Credential.scope.region, string(stype))

	// Calculate signature.
	newSignature := utils.GetSignature(signingKey, stringToSign)

	// Verify if signature match.
	if !compareSignatureV4(newSignature, signV4Values.Signature) {
		err = responses.ErrSignatureDoesNotMatch
		return
	}

	return
}

//// getScope generate a string of a specific date, an AWS region, and a service.
//func getScope(t time.Time, region string) string {
//	scope := strings.Join([]string{
//		t.Format(yyyymmdd),
//		region,
//		string(ServiceS3),
//		"aws4_request",
//	}, consts.SlashSeparator)
//	return scope
//}
