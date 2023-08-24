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
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/iam/set"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/utils"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	//ServiceSTS STS
	ServiceSTS serviceType = "sts"
)

// compareSignatureV4 returns true if and only if both signatures
// are equal. The signatures are expected to be HEX encoded strings
// according to the AWS S3 signature V4 spec.
func compareSignatureV4(sig1, sig2 string) bool {
	// The CTC using []byte(str) works because the hex encoding
	// is unique for a sequence of bytes. See also compareSignatureV2.
	return subtle.ConstantTimeCompare([]byte(sig1), []byte(sig2)) == 1
}

// doesPresignedSignatureMatch - Verify query headers with presigned signature
//   - http://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-query-string-auth.html
//
// returns nil if the signature matches.
func (s *service) doesPresignedSignatureMatch(hashedPayload string, r *http.Request, region string, stype serviceType) *responses.Error {
	// Copy request
	req := *r

	// Parse request query string.
	pSignValues, err := parsePreSignV4(req.Form, region, stype)
	if err != nil {
		return err
	}

	cred, _, s3Err := s.checkKeyValid(r, pSignValues.Credential.accessKey)
	if s3Err != nil {
		return s3Err
	}

	// Extract all the signed headers along with its values.
	extractedSignedHeaders, errCode := extractSignedHeaders(pSignValues.SignedHeaders, r)
	if errCode != nil {
		return errCode
	}

	// If the host which signed the request is slightly ahead in time (by less than MaxSkewTime) the
	// request should still be allowed.
	if pSignValues.Date.After(time.Now().UTC().Add(consts.MaxSkewTime)) {
		return responses.ErrRequestNotReadyYet
	}

	if time.Now().UTC().Sub(pSignValues.Date) > pSignValues.Expires {
		return responses.ErrExpiredPresignRequest
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

	token := req.Form.Get(consts.AmzSecurityToken)
	if token != "" {
		query.Set(consts.AmzSecurityToken, cred.SessionToken)
	}

	query.Set(consts.AmzAlgorithm, signV4Algorithm)

	// Construct the query.
	query.Set(consts.AmzDate, t.Format(iso8601Format))
	query.Set(consts.AmzExpires, strconv.Itoa(expireSeconds))
	query.Set(consts.AmzSignedHeaders, utils.GetSignedHeaders(extractedSignedHeaders))
	query.Set(consts.AmzCredential, cred.AccessKey+consts.SlashSeparator+pSignValues.Credential.getScope())

	defaultSigParams := set.CreateStringSet(
		consts.AmzContentSha256,
		consts.AmzSecurityToken,
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
		return responses.ErrSignatureDoesNotMatch
	}
	// Verify if expires query is same.
	if req.Form.Get(consts.AmzExpires) != query.Get(consts.AmzExpires) {
		return responses.ErrSignatureDoesNotMatch
	}
	// Verify if signed headers query is same.
	if req.Form.Get(consts.AmzSignedHeaders) != query.Get(consts.AmzSignedHeaders) {
		return responses.ErrSignatureDoesNotMatch
	}
	// Verify if credential query is same.
	if req.Form.Get(consts.AmzCredential) != query.Get(consts.AmzCredential) {
		return responses.ErrSignatureDoesNotMatch
	}
	// Verify if sha256 payload query is same.
	if clntHashedPayload != "" && clntHashedPayload != query.Get(consts.AmzContentSha256) {
		return responses.ErrContentSHA256Mismatch
	}
	// Verify if security token is correct.
	if token != "" && subtle.ConstantTimeCompare([]byte(token), []byte(cred.SessionToken)) != 1 {
		return responses.ErrInvalidToken
	}

	// Verify finally if signature is same.

	// Get canonical request.
	presignedCanonicalReq := utils.GetCanonicalRequest(extractedSignedHeaders, hashedPayload, encodedQuery, req.URL.Path, req.Method)

	// Get string to sign from canonical request.
	presignedStringToSign := utils.GetStringToSign(presignedCanonicalReq, t, pSignValues.Credential.getScope())

	// Get hmac presigned signing key.
	presignedSigningKey := utils.GetSigningKey(cred.SecretKey, pSignValues.Credential.scope.date,
		pSignValues.Credential.scope.region, string(stype))

	// Get new signature.
	newSignature := utils.GetSignature(presignedSigningKey, presignedStringToSign)

	// Verify signature.
	if !compareSignatureV4(req.Form.Get(consts.AmzSignature), newSignature) {
		return responses.ErrSignatureDoesNotMatch
	}
	return nil
}

// doesSignatureMatch - Verify authorization header with calculated header in accordance with
//   - http://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-authenticating-requests.html
//
// returns nil if signature matches.
func (s *service) doesSignatureMatch(hashedPayload string, r *http.Request, region string, stype serviceType) *responses.Error {
	// Copy request.
	req := *r

	// Save authorization header.
	v4Auth := req.Header.Get(consts.Authorization)

	// Parse signature version '4' header.
	signV4Values, err := parseSignV4(v4Auth, region, stype)
	if err != nil {
		return err
	}

	// Extract all the signed headers along with its values.
	extractedSignedHeaders, errCode := extractSignedHeaders(signV4Values.SignedHeaders, r)
	if errCode != nil {
		return errCode
	}

	cred, _, s3Err := s.checkKeyValid(r, signV4Values.Credential.accessKey)
	if s3Err != nil {
		return s3Err
	}

	// Extract date, if not present throw error.
	var date string
	if date = req.Header.Get(consts.AmzDate); date == "" {
		if date = r.Header.Get(consts.Date); date == "" {
			return responses.ErrMissingDateHeader
		}
	}

	// Parse date header.
	t, e := time.Parse(iso8601Format, date)
	if e != nil {
		return responses.ErrAuthorizationHeaderMalformed
	}

	// Query string.
	queryStr := req.URL.Query().Encode()

	// Get canonical request.
	canonicalRequest := utils.GetCanonicalRequest(extractedSignedHeaders, hashedPayload, queryStr, req.URL.Path, req.Method)

	// Get string to sign from canonical request.
	stringToSign := utils.GetStringToSign(canonicalRequest, t, signV4Values.Credential.getScope())

	// Get hmac signing key.
	signingKey := utils.GetSigningKey(cred.SecretKey, signV4Values.Credential.scope.date,
		signV4Values.Credential.scope.region, string(stype))

	// Calculate signature.
	newSignature := utils.GetSignature(signingKey, stringToSign)

	// Verify if signature match.
	if !compareSignatureV4(newSignature, signV4Values.Signature) {
		return responses.ErrSignatureDoesNotMatch
	}

	// Return error none.
	return nil
}

// getScope generate a string of a specific date, an AWS region, and a service.
func getScope(t time.Time, region string) string {
	scope := strings.Join([]string{
		t.Format(yyyymmdd),
		region,
		string(ServiceS3),
		"aws4_request",
	}, consts.SlashSeparator)
	return scope
}
