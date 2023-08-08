package auth

import (
	"github.com/bittorrent/go-btfs/s3/apierrors"
	"net/http"
)

type Service interface {
	CheckSignatureV4Verify(r *http.Request, region string, stype serviceType) (s3Error apierrors.ErrorCode)
	CheckACL(r *http.Request, region string, stype serviceType) (s3Error apierrors.ErrorCode)
}
