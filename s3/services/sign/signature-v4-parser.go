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

package sign

import (
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/responses"
	"net/url"
	"strings"
	"time"
)

// credentialHeader data type represents structured form of Credential
// string from authorization header.
type credentialHeader struct {
	accessKey string
	scope     struct {
		date    time.Time
		region  string
		service string
		request string
	}
}

// Return scope string.
func (c credentialHeader) getScope() string {
	return strings.Join([]string{
		c.scope.date.Format(yyyymmdd),
		c.scope.region,
		c.scope.service,
		c.scope.request,
	}, consts.SlashSeparator)
}

// parse credentialHeader string into its structured form.
func parseCredentialHeader(credElement string, region string, stype serviceType) (ch credentialHeader, rerr *responses.Error) {
	creds := strings.SplitN(strings.TrimSpace(credElement), "=", 2)
	if len(creds) != 2 {
		return ch, responses.ErrMissingFields
	}
	if creds[0] != "Credential" {
		return ch, responses.ErrMissingCredTag
	}
	credElements := strings.Split(strings.TrimSpace(creds[1]), consts.SlashSeparator)
	if len(credElements) < 5 {
		return ch, responses.ErrCredMalformed
	}
	accessKey := strings.Join(credElements[:len(credElements)-4], consts.SlashSeparator) // The access key may contain one or more `/`

	// Save access key id.
	cred := credentialHeader{
		accessKey: accessKey,
	}
	credElements = credElements[len(credElements)-4:]
	var e error
	cred.scope.date, e = time.Parse(yyyymmdd, credElements[0])
	if e != nil {
		return ch, responses.ErrAuthorizationHeaderMalformed
	}

	cred.scope.region = credElements[1]
	// Verify if region is valid.
	sRegion := cred.scope.region
	// Region is set to be empty, we use whatever was sent by the
	// request and proceed further. This is a work-around to address
	// an important problem for ListBuckets() getting signed with
	// different regions.
	if region == "" {
		region = sRegion
	}
	// Should validate region, only if region is set.
	if !isValidRegion(sRegion, region) {
		return ch, responses.ErrAuthorizationHeaderMalformed
	}
	if credElements[2] != string(stype) {
		switch stype {
		case ServiceSTS:
			return ch, responses.ErrAuthorizationHeaderMalformed
		}
		return ch, responses.ErrAuthorizationHeaderMalformed
	}
	cred.scope.service = credElements[2]
	if credElements[3] != "aws4_request" {
		return ch, responses.ErrAuthorizationHeaderMalformed
	}
	cred.scope.request = credElements[3]
	return cred, nil
}

// Parse signature from signature tag.
func parseSignature(signElement string) (string, *responses.Error) {
	signFields := strings.Split(strings.TrimSpace(signElement), "=")
	if len(signFields) != 2 {
		return "", responses.ErrMissingFields
	}
	if signFields[0] != "Signature" {
		return "", responses.ErrMissingSignTag
	}
	if signFields[1] == "" {
		return "", responses.ErrMissingFields
	}
	signature := signFields[1]
	return signature, nil
}

// Parse slice of signed headers from signed headers tag.
func parseSignedHeader(signedHdrElement string) ([]string, *responses.Error) {
	signedHdrFields := strings.Split(strings.TrimSpace(signedHdrElement), "=")
	if len(signedHdrFields) != 2 {
		return nil, responses.ErrMissingFields
	}
	if signedHdrFields[0] != "SignedHeaders" {
		return nil, responses.ErrMissingSignHeadersTag
	}
	if signedHdrFields[1] == "" {
		return nil, responses.ErrMissingFields
	}
	signedHeaders := strings.Split(signedHdrFields[1], ";")
	return signedHeaders, nil
}

// signValues data type represents structured form of AWS Signature V4 header.
type signValues struct {
	Credential    credentialHeader
	SignedHeaders []string
	Signature     string
}

// preSignValues data type represents structued form of AWS Signature V4 query string.
type preSignValues struct {
	signValues
	Date    time.Time
	Expires time.Duration
}

// Parses signature version '4' query string of the following form.
//
//	querystring = X-Amz-Algorithm=algorithm
//	querystring += &X-Amz-Credential= urlencode(accessKey + '/' + credential_scope)
//	querystring += &X-Amz-Date=date
//	querystring += &X-Amz-Expires=timeout interval
//	querystring += &X-Amz-SignedHeaders=signed_headers
//	querystring += &X-Amz-Signature=signature
//
// verifies if any of the necessary query params are missing in the presigned request.
func doesV4PresignParamsExist(query url.Values) *responses.Error {
	v4PresignQueryParams := []string{consts.AmzAlgorithm, consts.AmzCredential, consts.AmzSignature, consts.AmzDate, consts.AmzSignedHeaders, consts.AmzExpires}
	for _, v4PresignQueryParam := range v4PresignQueryParams {
		if _, ok := query[v4PresignQueryParam]; !ok {
			return responses.ErrInvalidQueryParams
		}
	}
	return nil
}

// Parses all the presigned signature values into separate elements.
func parsePreSignV4(query url.Values, region string, stype serviceType) (psv preSignValues, rerr *responses.Error) {
	// verify whether the required query params exist.
	rerr = doesV4PresignParamsExist(query)
	if rerr != nil {
		return psv, rerr
	}

	// Verify if the query algorithm is supported or not.
	if query.Get(consts.AmzAlgorithm) != signV4Algorithm {
		return psv, responses.ErrAuthorizationHeaderMalformed
	}

	// Initialize signature version '4' structured header.
	preSignV4Values := preSignValues{}

	// Save credential.
	preSignV4Values.Credential, rerr = parseCredentialHeader("Credential="+query.Get(consts.AmzCredential), region, stype)
	if rerr != nil {
		return psv, rerr
	}

	var e error
	// Save date in native time.Time.
	preSignV4Values.Date, e = time.Parse(iso8601Format, query.Get(consts.AmzDate))
	if e != nil {
		return psv, responses.ErrAuthorizationHeaderMalformed
	}

	// Save expires in native time.Duration.
	preSignV4Values.Expires, e = time.ParseDuration(query.Get(consts.AmzExpires) + "s")
	if e != nil {
		return psv, responses.ErrAuthorizationHeaderMalformed
	}

	if preSignV4Values.Expires < 0 {
		return psv, responses.ErrAuthorizationHeaderMalformed
	}

	// Check if Expiry time is less than 7 days (value in seconds).
	if preSignV4Values.Expires.Seconds() > 604800 {
		return psv, responses.ErrAuthorizationHeaderMalformed
	}

	// Save signed headers.
	preSignV4Values.SignedHeaders, rerr = parseSignedHeader("SignedHeaders=" + query.Get(consts.AmzSignedHeaders))
	if rerr != nil {
		return psv, rerr
	}

	// Save signature.
	preSignV4Values.Signature, rerr = parseSignature("Signature=" + query.Get(consts.AmzSignature))
	if rerr != nil {
		return psv, rerr
	}

	// Return structed form of signature query string.
	return preSignV4Values, nil
}

// Parses signature version '4' header of the following form.
//
//	Authorization: algorithm Credential=accessKeyID/credScope, \
//	        SignedHeaders=signedHeaders, Signature=signature
func parseSignV4(v4Auth string, region string, stype serviceType) (sv signValues, aec *responses.Error) {
	// credElement is fetched first to skip replacing the space in access key.
	credElement := strings.TrimPrefix(strings.Split(strings.TrimSpace(v4Auth), ",")[0], signV4Algorithm)
	// Replace all spaced strings, some clients can send spaced
	// parameters and some won't. So we pro-actively remove any spaces
	// to make parsing easier.
	v4Auth = strings.ReplaceAll(v4Auth, " ", "")
	if v4Auth == "" {
		return sv, responses.ErrAuthHeaderEmpty
	}

	// Verify if the header algorithm is supported or not.
	if !strings.HasPrefix(v4Auth, signV4Algorithm) {
		return sv, responses.ErrSignatureVersionNotSupported
	}

	// Strip off the Algorithm prefix.
	v4Auth = strings.TrimPrefix(v4Auth, signV4Algorithm)
	authFields := strings.Split(strings.TrimSpace(v4Auth), ",")
	if len(authFields) != 3 {
		return sv, responses.ErrMissingFields
	}

	// Initialize signature version '4' structured header.
	signV4Values := signValues{}

	var s3Err *responses.Error
	// Save credentail values.
	signV4Values.Credential, s3Err = parseCredentialHeader(strings.TrimSpace(credElement), region, stype)
	if s3Err != nil {
		return sv, s3Err
	}

	// Save signed headers.
	signV4Values.SignedHeaders, s3Err = parseSignedHeader(authFields[1])
	if s3Err != nil {
		return sv, s3Err
	}

	// Save signature.
	signV4Values.Signature, s3Err = parseSignature(authFields[2])
	if s3Err != nil {
		return sv, s3Err
	}

	// Return the structure here.
	return signV4Values, nil
}
