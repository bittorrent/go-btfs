package sign

import (
	"encoding/hex"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/etag"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"net/http"
)

// isReqAuthenticated Verify if request has valid AWS Signature Version '4'.
func (s *service) isReqAuthenticated(r *http.Request, region string, stype serviceType) (ack string, rerr *responses.Error) {
	ack, rerr = s.reqSignatureV4Verify(r, region, stype)
	if rerr != nil {
		return
	}
	clientETag, err := etag.FromContentMD5(r.Header)
	if err != nil {
		rerr = responses.ErrInvalidDigest
		return
	}

	// Extract either 'X-Amz-Content-Sha256' header or 'X-Amz-Content-Sha256' query parameter (if V4 presigned)
	// Do not verify 'X-Amz-Content-Sha256' if skipSHA256.
	var contentSHA256 []byte
	if skipSHA256 := SkipContentSha256Cksum(r); !skipSHA256 && isRequestPresignedSignatureV4(r) {
		if sha256Sum, ok := r.Form[consts.AmzContentSha256]; ok && len(sha256Sum) > 0 {
			contentSHA256, err = hex.DecodeString(sha256Sum[0])
			if err != nil {
				rerr = responses.ErrContentSHA256Mismatch
				return
			}
		}
	} else if _, ok := r.Header[consts.AmzContentSha256]; !skipSHA256 && ok {
		contentSHA256, err = hex.DecodeString(r.Header.Get(consts.AmzContentSha256))
		if err != nil || len(contentSHA256) == 0 {
			rerr = responses.ErrContentSHA256Mismatch
		}
	}

	// Verify 'Content-Md5' and/or 'X-Amz-Content-Sha256' if present.
	// The verification happens implicit during reading.
	reader, err := hash.NewReader(r.Body, -1, clientETag.String(), hex.EncodeToString(contentSHA256), -1)
	if err != nil {
		rerr = responses.ErrInternalError
		return
	}

	r.Body = reader

	return
}

func (s *service) reqSignatureV4Verify(r *http.Request, region string, stype serviceType) (ack string, rerr *responses.Error) {
	sha256sum, err := GetContentSha256Cksum(r, stype)
	if err != nil {
		rerr = responses.ErrInternalError
		return
	}
	switch {
	case IsRequestSignatureV4(r):
		ack, rerr = s.doesSignatureMatch(sha256sum, r, region, stype)
	case isRequestPresignedSignatureV4(r):
		ack, rerr = s.doesPresignedSignatureMatch(sha256sum, r, region, stype)
	default:
		rerr = responses.ErrAccessDenied
	}
	return
}
