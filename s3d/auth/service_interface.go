package auth

import (
	"github.com/bittorrent/go-btfs/s3d/apierrors"
	"net/http"
)

type Service interface {
	CheckSignatureAndAcl(r *http.Request, region string, stype serviceType) (s3Error apierrors.ErrorCode)
}
