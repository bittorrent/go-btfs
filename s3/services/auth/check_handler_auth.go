package auth

import (
	"context"
	"encoding/hex"
	"github.com/bittorrent/go-btfs/s3/services"
	"net/http"

	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/etag"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
)

// CheckRequestAuthTypeCredential Check request auth type verifies the incoming http request
//   - validates the request signature
//   - validates the policy action if anonymous tests bucket policies if any,
//     for authenticated requests validates IAM policies.
//
// returns APIErrorcode if any to be replied to the client.
// Additionally, returns the accessKey used in the request, and if this request is by an admin.
func (s *Service) CheckRequestAuthTypeCredential(ctx context.Context, r *http.Request) (cred *services.AccessKey, err error) {
	// check signature
	switch GetRequestAuthType(r) {
	case AuthTypeSigned, AuthTypePresigned:
		region := ""
		if err = s.IsReqAuthenticated(ctx, r, region, ServiceS3); err != nil {
			return
		}
		cred, err = s.getReqAccessKeyV4(r, region, ServiceS3)
	default:
		err = services.ErrSignatureVersionNotSupported
		return
	}

	return
}

func (s *Service) ReqSignatureV4Verify(r *http.Request, region string, stype serviceType) error {
	sha256sum := getContentSha256Cksum(r, stype)
	switch {
	case IsRequestSignatureV4(r):
		return s.doesSignatureMatch(sha256sum, r, region, stype)
	case isRequestPresignedSignatureV4(r):
		return s.doesPresignedSignatureMatch(sha256sum, r, region, stype)
	default:
		return services.ErrAccessDenied
	}
}

// IsReqAuthenticated Verify if request has valid AWS Signature Version '4'.
func (s *Service) IsReqAuthenticated(ctx context.Context, r *http.Request, region string, stype serviceType) (err error) {
	if err = s.ReqSignatureV4Verify(r, region, stype); err != nil {
		return
	}
	clientETag, err := etag.FromContentMD5(r.Header)
	if err != nil {
		err = services.ErrInvalidDigest
		return
	}

	// Extract either 'X-Amz-Content-Sha256' header or 'X-Amz-Content-Sha256' query parameter (if V4 presigned)
	// Do not verify 'X-Amz-Content-Sha256' if skipSHA256.
	var contentSHA256 []byte
	if skipSHA256 := SkipContentSha256Cksum(r); !skipSHA256 && isRequestPresignedSignatureV4(r) {
		if sha256Sum, ok := r.Form[consts.AmzContentSha256]; ok && len(sha256Sum) > 0 {
			contentSHA256, err = hex.DecodeString(sha256Sum[0])
			if err != nil {
				err = services.ErrContentSHA256Mismatch
				return
			}
		}
	} else if _, ok := r.Header[consts.AmzContentSha256]; !skipSHA256 && ok {
		contentSHA256, err = hex.DecodeString(r.Header.Get(consts.AmzContentSha256))
		if err != nil || len(contentSHA256) == 0 {
			err = services.ErrContentSHA256Mismatch
			return
		}
	}

	// Verify 'Content-Md5' and/or 'X-Amz-Content-Sha256' if present.
	// The verification happens implicit during reading.
	reader, err := hash.NewReader(r.Body, -1, clientETag.String(), hex.EncodeToString(contentSHA256), -1)
	if err != nil {
		err = services.ErrInternalError
		return
	}
	r.Body = reader
	return
}

//// ValidateAdminSignature validate admin Signature
//func (s *Service) ValidateAdminSignature(ctx context.Context, r *http.Request, region string) (Credentials, map[string]interface{}, bool, handlers.Errorcode) {
//	var cred Credentials
//	var owner bool
//	s3Err := handlers.ErrcodeAccessDenied
//	if _, ok := r.Header[consts.AmzContentSha256]; ok &&
//		GetRequestAuthType(r) == AuthTypeSigned {
//		// We only support admin credentials to access admin APIs.
//		cred, s3Err = GetReqAccessKeyV4(r, region, ServiceS3)
//		if s3Err != handlers.ErrcodeNone {
//			return cred, nil, owner, s3Err
//		}
//
//		// we only support V4 (no presign) with auth body
//		s3Err = s.IsReqAuthenticated(ctx, r, region, ServiceS3)
//	}
//	if s3Err != handlers.ErrcodeNone {
//		return cred, nil, owner, s3Err
//	}
//
//	return cred, nil, owner, handlers.ErrcodeNone
//}
////
//func (s *Service) GetCredential(r *http.Request) (cred auth.Credentials, owner bool, s3Err handlers.Errorcode) {
//	switch GetRequestAuthType(r) {
//	case AuthTypeUnknown:
//		s3Err = handlers.ErrcodeSignatureVersionNotSupported
//	case AuthTypeSignedV2, AuthTypePresignedV2:
//		cred, owner, s3Err = s.getReqAccessKeyV2(r)
//	case AuthTypeStreamingSigned, AuthTypePresigned, AuthTypeSigned:
//		region := ""
//		cred, owner, s3Err = s.GetReqAccessKeyV4(r, region, ServiceS3)
//	}
//	return
//}
