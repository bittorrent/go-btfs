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

	size := r.ContentLength

	if size == -1 {
		rerr = responses.ErrMissingContentLength
		return
	}

	if size > consts.MaxObjectSize {
		rerr = responses.ErrEntityTooLarge
		return
	}

	md5Hex, sha256Hex, rerr := s.getClientCheckSum(r)
	if rerr != nil {
		return
	}

	reader, err := hash.NewReader(r.Body, size, md5Hex, sha256Hex, size)
	if err != nil {
		rerr = responses.ErrInternalError
		return
	}

	r.Body = reader

	return
}

func (s *service) getClientCheckSum(r *http.Request) (md5TagStr, sha256SumStr string, rerr *responses.Error) {
	eTag, err := etag.FromContentMD5(r.Header)
	if err != nil {
		rerr = responses.ErrInvalidDigest
		return
	}
	md5TagStr = eTag.String()

	skipSHA256 := SkipContentSha256Cksum(r)
	if skipSHA256 {
		return
	}

	var (
		contentSHA256 []byte
		sha256Sum     []string
	)

	if isRequestPresignedSignatureV4(r) {
		sha256Sum = r.Form[consts.AmzContentSha256]
	} else {
		sha256Sum = r.Header[consts.AmzContentSha256]
	}

	if len(sha256Sum) > 0 {
		contentSHA256, err = hex.DecodeString(sha256Sum[0])
		if err != nil || len(contentSHA256) == 0 {
			rerr = responses.ErrContentSHA256Mismatch
			return
		}
		sha256SumStr = hex.EncodeToString(contentSHA256)
	}

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
