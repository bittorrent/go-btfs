package sign

import (
	"github.com/bittorrent/go-btfs/s3/api/responses"
	"github.com/bittorrent/go-btfs/s3/consts"
	"net/http"
	"time"
)

// AWS Signature Version '4' constants.
const (
	signV2Algorithm = "AWS"
	signV4Algorithm = "AWS4-HMAC-SHA256"
	iso8601Format   = "20060102T150405Z"
	yyyymmdd        = "20060102"
)

func (s *service) reqSignatureV4Verify(r *http.Request, region string) (ack string, rerr *responses.Error) {
	sha256sum, err := GetContentSHA256Checksum(r)
	if err != nil {
		rerr = responses.ErrInternalError
		return
	}
	ack, rerr = s.doesSignatureMatch(sha256sum, r, region)
	return
}

// doesSignatureMatch - Verify authorization header with calculated header in accordance with
//   - http://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-authenticating-requests.html
//
// returns nil if signature matches.
func (s *service) doesSignatureMatch(hashedPayload string, r *http.Request, region string) (ack string, rerr *responses.Error) {
	// Copy request.
	req := *r

	// Save authorization header.
	v4Auth := req.Header.Get(consts.Authorization)

	// Parse signature version '4' header.
	signV4Values, rerr := parseSignV4(v4Auth, region)
	if rerr != nil {
		return
	}

	// Extract all the signed headers along with its values.
	extractedSignedHeaders, rerr := extractSignedHeaders(signV4Values.SignedHeaders, r)
	if rerr != nil {
		return
	}

	ack = signV4Values.Credential.accessKey
	secret, rerr := s.checkKeyValid(ack)
	if rerr != nil {
		return
	}

	// Extract date, if not present throw error.
	var date string
	if date = req.Header.Get(consts.AmzDate); date == "" {
		if date = r.Header.Get(consts.Date); date == "" {
			rerr = responses.ErrMissingDateHeader
			return
		}
	}

	// Parse date header.
	t, err := time.Parse(iso8601Format, date)
	if err != nil {
		rerr = responses.ErrAuthorizationHeaderMalformed
		return
	}

	// Query string.
	queryStr := req.URL.Query().Encode()

	// Get canonical request.
	canonicalRequest := GetCanonicalRequest(extractedSignedHeaders, hashedPayload, queryStr, req.URL.Path, req.Method)

	// Get string to sign from canonical request.
	stringToSign := GetStringToSign(canonicalRequest, t, signV4Values.Credential.getScope())

	// Get hmac signing key.
	signingKey := GetSigningKey(secret, signV4Values.Credential.scope.date,
		signV4Values.Credential.scope.region)

	// Calculate signature.
	newSignature := GetSignature(signingKey, stringToSign)

	// Verify if signature match.
	if !compareSignatureV4(newSignature, signV4Values.Signature) {
		rerr = responses.ErrSignatureDoesNotMatch
		return
	}

	// Return error none.
	return
}

// check if the access key is valid and recognized, additionally
// also returns if the access key is owner/admin.
func (s *service) checkKeyValid(ack string) (secret string, rerr *responses.Error) {
	secret, exists, enable, err := s.getSecret(ack)
	if err != nil {
		rerr = responses.ErrInternalError
		return
	}

	if !exists {
		rerr = responses.ErrInvalidAccessKeyID
		return
	}

	if !enable {
		rerr = responses.ErrAccessKeyDisabled
		return
	}

	return
}
