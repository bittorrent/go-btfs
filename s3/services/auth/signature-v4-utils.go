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
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/handlers"
)

// http Header "x-amz-content-sha256" == "UNSIGNED-PAYLOAD" indicates that the
// client did not calculate sha256 of the payload.
const unsignedPayload = "UNSIGNED-PAYLOAD"

// isValidRegion - verify if incoming region value is valid with configured Region.
func isValidRegion(reqRegion string, confRegion string) bool {
	if confRegion == "" {
		return true
	}
	if confRegion == "US" {
		confRegion = consts.DefaultRegion
	}
	// Some older s3 clients set region as "US" instead of
	// globalDefaultRegion, handle it.
	if reqRegion == "US" {
		reqRegion = consts.DefaultRegion
	}
	return reqRegion == confRegion
}

func contains(slice interface{}, elem interface{}) bool {
	v := reflect.ValueOf(slice)
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			if v.Index(i).Interface() == elem {
				return true
			}
		}
	}
	return false
}

// extractSignedHeaders extract signed headers from Authorization header
func extractSignedHeaders(signedHeaders []string, r *http.Request) (http.Header, handlers.ErrorCode) {
	reqHeaders := r.Header
	reqQueries := r.Form
	// find whether "host" is part of list of signed headers.
	// if not return ErrUnsignedHeaders. "host" is mandatory.
	if !contains(signedHeaders, "host") {
		return nil, handlers.ErrUnsignedHeaders
	}
	extractedSignedHeaders := make(http.Header)
	for _, header := range signedHeaders {
		// `host` will not be found in the headers, can be found in r.Host.
		// but its alway necessary that the list of signed headers containing host in it.
		val, ok := reqHeaders[http.CanonicalHeaderKey(header)]
		if !ok {
			// try to set headers from Query String
			val, ok = reqQueries[header]
		}
		if ok {
			extractedSignedHeaders[http.CanonicalHeaderKey(header)] = val
			continue
		}
		switch header {
		case "expect":
			// Golang http server strips off 'Expect' header, if the
			// client sent this as part of signed headers we need to
			// handle otherwise we would see a signature mismatch.
			// `aws-cli` sets this as part of signed headers.
			//
			// According to
			// http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.20
			// Expect header is always of form:
			//
			//   Expect       =  "Expect" ":" 1#expectation
			//   expectation  =  "100-continue" | expectation-extension
			//
			// So it safe to assume that '100-continue' is what would
			// be sent, for the time being keep this work around.
			// Adding a *TODO* to remove this later when Golang server
			// doesn't filter out the 'Expect' header.
			extractedSignedHeaders.Set(header, "100-continue")
		case "host":
			// Go http server removes "host" from Request.Header

			//extractedSignedHeaders.Set(header, r.Host)
			// todo use r.Host, or filedag-web deal with
			//value := strings.Split(r.Host, ":")
			extractedSignedHeaders.Set(header, r.Host)
		case "transfer-encoding":
			// Go http server removes "host" from Request.Header
			extractedSignedHeaders[http.CanonicalHeaderKey(header)] = r.TransferEncoding
		case "content-length":
			// Signature-V4 spec excludes Content-Length from signed headers list for signature calculation.
			// But some clients deviate from this rule. Hence we consider Content-Length for signature
			// calculation to be compatible with such clients.
			extractedSignedHeaders.Set(header, strconv.FormatInt(r.ContentLength, 10))
		default:
			return nil, handlers.ErrUnsignedHeaders
		}
	}
	return extractedSignedHeaders, handlers.ErrNone
}

// Returns SHA256 for calculating canonical-request.
func getContentSha256Cksum(r *http.Request, stype serviceType) string {
	//if stype == ServiceSTS {
	//	payload, err := ioutil.ReadAll(io.LimitReader(r.Body, consts.StsRequestBodyLimit))
	//	if err != nil {
	//		//log.Errorf("ServiceSTS ReadAll err:%v", err)
	//	}
	//	sum256 := sha256.Sum256(payload)
	//	r.Body = ioutil.NopCloser(bytes.NewReader(payload))
	//	return hex.EncodeToString(sum256[:])
	//}

	var (
		defaultSha256Cksum string
		v                  []string
		ok                 bool
	)

	// For a presigned request we look at the query param for sha256.
	if isRequestPresignedSignatureV4(r) {
		// X-Amz-Content-Sha256, if not set in presigned requests, checksum
		// will default to 'UNSIGNED-PAYLOAD'.
		defaultSha256Cksum = unsignedPayload
		v, ok = r.Form[consts.AmzContentSha256]
		if !ok {
			v, ok = r.Header[consts.AmzContentSha256]
		}
	} else {
		// X-Amz-Content-Sha256, if not set in signed requests, checksum
		// will default to sha256([]byte("")).
		defaultSha256Cksum = consts.EmptySHA256
		v, ok = r.Header[consts.AmzContentSha256]
	}

	// We found 'X-Amz-Content-Sha256' return the captured value.
	if ok {
		return v[0]
	}

	// We couldn't find 'X-Amz-Content-Sha256'.
	return defaultSha256Cksum
}

// isRequestSignatureV4 Verify if request has AWS Signature Version '4'.
func isRequestSignatureV4(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Authorization"), signV4Algorithm)
}

// Verify if request has AWS PreSign Version '4'.
func isRequestPresignedSignatureV4(r *http.Request) bool {
	_, ok := r.URL.Query()["X-Amz-Credential"]
	return ok
}

// SkipContentSha256Cksum returns true if caller needs to skip
// payload checksum, false if not.
func SkipContentSha256Cksum(r *http.Request) bool {
	var (
		v  []string
		ok bool
	)

	if isRequestPresignedSignatureV4(r) {
		v, ok = r.Form[consts.AmzContentSha256]
		if !ok {
			v, ok = r.Header[consts.AmzContentSha256]
		}
	} else {
		v, ok = r.Header[consts.AmzContentSha256]
	}

	// Skip if no header was set.
	if !ok {
		return true
	}

	// If x-amz-content-sha256 is set and the value is not
	// 'UNSIGNED-PAYLOAD' we should validate the content sha256.
	switch v[0] {
	case unsignedPayload:
		return true
	case consts.EmptySHA256:
		// some broken clients set empty-sha256
		// with > 0 content-length in the body,
		// we should skip such clients and allow
		// blindly such insecure clients only if
		// S3 strict compatibility is disabled.
		if r.ContentLength > 0 {
			// We return true only in situations when
			// deployment has asked MinIO to allow for
			// such broken clients and content-length > 0.
			return true
		}
	}
	return false
}
