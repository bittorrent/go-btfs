package auth

import (
	"bytes"
	"context"
	"encoding/hex"
	s3action "github.com/bittorrent/go-btfs/s3/action"
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/policy"
	"github.com/bittorrent/go-btfs/s3/utils/hash"
	"github.com/bittorrent/go-btfs/s3d/store"
	"io"
	"net/http"

	"github.com/bittorrent/go-btfs/s3d/apierrors"
	"github.com/bittorrent/go-btfs/s3d/consts"
	"github.com/bittorrent/go-btfs/s3d/etag"
)

// AuthSys auth and sign system
type AuthSys struct{}

// NewAuthSys new an AuthSys
func NewAuthSys() *AuthSys {
	return &AuthSys{}
}

// CheckRequestAuthTypeCredential Check request auth type verifies the incoming http request
//   - validates the request signature
//   - validates the policy action if anonymous tests bucket policies if any,
//     for authenticated requests validates IAM policies.
//
// returns APIErrorCode if any to be replied to the client.
// Additionally, returns the accessKey used in the request, and if this request is by an admin.
func (s *AuthSys) CheckRequestAuthTypeCredential(ctx context.Context, r *http.Request, action s3action.Action, bucketName string, bmSys *store.BucketMetadataSys) (cred Credentials, err error) {
	//todo 是否需要判断
	if bucketName == "" {
		return cred, handlers.ErrBucketNotFound
	}

	// 1.check signature
	switch GetRequestAuthType(r) {
	case AuthTypeUnknown, AuthTypeStreamingSigned:
		return cred, apierrors.ErrSignatureVersionNotSupported
	case AuthTypePresignedV2, AuthTypeSignedV2:
		return cred, apierrors.ErrSignatureVersionNotSupported
	case AuthTypeSigned, AuthTypePresigned:
		region := ""
		if s3Err = s.IsReqAuthenticated(ctx, r, region, ServiceS3); s3Err != apierrors.ErrNone {
			return cred, s3Err
		}
		cred, s3Err = GetReqAccessKeyV4(r, region, ServiceS3)
	}
	if s3Err != apierrors.ErrNone {
		return cred, s3Err
	}

	// CreateBucketAction
	if action == action.CreateBucketAction {
		// To extract region from XML in request body, get copy of request body.
		payload, err := io.ReadAll(io.LimitReader(r.Body, consts.MaxLocationConstraintSize))
		if err != nil {
			//log.Errorf("ReadAll err:%v", err)
			return cred, apierrors.ErrMalformedXML
		}

		// Populate payload to extract location constraint.
		r.Body = io.NopCloser(bytes.NewReader(payload))
		//todo check HasBucket
		if bmSys.HasBucket(ctx, bucketName) {
			return cred, apierrors.ErrBucketAlreadyExists
		}
	}

	// 2.check acl
	//todo 获取bucket用户信息:owner, acl
	meta, err := bmSys.GetBucketMeta(ctx, bucketName)
	if err != nil {
		return cred, apierrors.ErrAccessDenied
	}

	if policy.IsAllowed(meta.Owner == cred.AccessKey, meta.Acl, action) == false {
		return cred, apierrors.ErrAccessDenied
	}

	return cred, apierrors.ErrNone
}

func (s *AuthSys) ReqSignatureV4Verify(r *http.Request, region string, stype serviceType) (s3Error apierrors.ErrorCode) {
	sha256sum := getContentSha256Cksum(r, stype)
	switch {
	case IsRequestSignatureV4(r):
		return DoesSignatureMatch(sha256sum, r, region, stype)
	case isRequestPresignedSignatureV4(r):
		return DoesPresignedSignatureMatch(sha256sum, r, region, stype)
	default:
		return apierrors.ErrAccessDenied
	}
}

// IsReqAuthenticated Verify if request has valid AWS Signature Version '4'.
func (s *AuthSys) IsReqAuthenticated(ctx context.Context, r *http.Request, region string, stype serviceType) (s3Error apierrors.ErrorCode) {
	if errCode := s.ReqSignatureV4Verify(r, region, stype); errCode != apierrors.ErrNone {
		return errCode
	}
	clientETag, err := etag.FromContentMD5(r.Header)
	if err != nil {
		return apierrors.ErrInvalidDigest
	}

	// Extract either 'X-Amz-Content-Sha256' header or 'X-Amz-Content-Sha256' query parameter (if V4 presigned)
	// Do not verify 'X-Amz-Content-Sha256' if skipSHA256.
	var contentSHA256 []byte
	if skipSHA256 := SkipContentSha256Cksum(r); !skipSHA256 && isRequestPresignedSignatureV4(r) {
		if sha256Sum, ok := r.Form[consts.AmzContentSha256]; ok && len(sha256Sum) > 0 {
			contentSHA256, err = hex.DecodeString(sha256Sum[0])
			if err != nil {
				return apierrors.ErrContentSHA256Mismatch
			}
		}
	} else if _, ok := r.Header[consts.AmzContentSha256]; !skipSHA256 && ok {
		contentSHA256, err = hex.DecodeString(r.Header.Get(consts.AmzContentSha256))
		if err != nil || len(contentSHA256) == 0 {
			return apierrors.ErrContentSHA256Mismatch
		}
	}

	// Verify 'Content-Md5' and/or 'X-Amz-Content-Sha256' if present.
	// The verification happens implicit during reading.
	reader, err := hash.NewReader(r.Body, -1, clientETag.String(), hex.EncodeToString(contentSHA256), -1)
	if err != nil {
		return apierrors.ErrInternalError
	}
	r.Body = reader
	return apierrors.ErrNone
}

//// ValidateAdminSignature validate admin Signature
//func (s *AuthSys) ValidateAdminSignature(ctx context.Context, r *http.Request, region string) (Credentials, map[string]interface{}, bool, apierrors.ErrorCode) {
//	var cred Credentials
//	var owner bool
//	s3Err := apierrors.ErrAccessDenied
//	if _, ok := r.Header[consts.AmzContentSha256]; ok &&
//		GetRequestAuthType(r) == AuthTypeSigned {
//		// We only support admin credentials to access admin APIs.
//		cred, s3Err = GetReqAccessKeyV4(r, region, ServiceS3)
//		if s3Err != apierrors.ErrNone {
//			return cred, nil, owner, s3Err
//		}
//
//		// we only support V4 (no presign) with auth body
//		s3Err = s.IsReqAuthenticated(ctx, r, region, ServiceS3)
//	}
//	if s3Err != apierrors.ErrNone {
//		return cred, nil, owner, s3Err
//	}
//
//	return cred, nil, owner, apierrors.ErrNone
//}
////
//func (s *AuthSys) GetCredential(r *http.Request) (cred auth.Credentials, owner bool, s3Err apierrors.ErrorCode) {
//	switch GetRequestAuthType(r) {
//	case AuthTypeUnknown:
//		s3Err = apierrors.ErrSignatureVersionNotSupported
//	case AuthTypeSignedV2, AuthTypePresignedV2:
//		cred, owner, s3Err = s.getReqAccessKeyV2(r)
//	case AuthTypeStreamingSigned, AuthTypePresigned, AuthTypeSigned:
//		region := ""
//		cred, owner, s3Err = s.GetReqAccessKeyV4(r, region, ServiceS3)
//	}
//	return
//}
