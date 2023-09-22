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
	"github.com/bittorrent/go-btfs/s3/api/responses"
	"github.com/bittorrent/go-btfs/s3/consts"
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
func parseCredentialHeader(credElement string, region string) (ch credentialHeader, rerr *responses.Error) {
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
	if credElements[2] != "s3" {
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

// Parses signature version '4' header of the following form.
//
//	Authorization: algorithm Credential=accessKeyID/credScope, \
//	        SignedHeaders=signedHeaders, Signature=signature
func parseSignV4(v4Auth string, region string) (sv signValues, aec *responses.Error) {
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
	// Save credential values.
	signV4Values.Credential, s3Err = parseCredentialHeader(strings.TrimSpace(credElement), region)
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
