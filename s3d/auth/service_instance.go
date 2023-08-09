package auth

import (
	"github.com/bittorrent/go-btfs/s3d/apierrors"
	"net/http"
)

type service struct {
}

func newService() (svc *service, err error) {
	svc = &service{}
	return
}

func (s *service) CheckSignatureV4Verify(r *http.Request, region string, stype serviceType) (s3Error apierrors.ErrorCode) {
	sha256sum := getContentSha256Cksum(r, stype)
	switch {
	case isRequestSignatureV4(r):
		return DoesSignatureMatch(sha256sum, r, region, stype)
	case isRequestPresignedSignatureV4(r):
		return DoesPresignedSignatureMatch(sha256sum, r, region, stype)
	default:
		return apierrors.ErrAccessDenied
	}
}

func (s *service) CheckACL(r *http.Request, region string, stype serviceType) (s3Error apierrors.ErrorCode) {
	return
}
