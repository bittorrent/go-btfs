package auth

import (
	"context"
	s3action "github.com/bittorrent/go-btfs/s3d/action"
	"github.com/bittorrent/go-btfs/s3d/apierrors"
	"github.com/bittorrent/go-btfs/s3d/store"
	"net/http"
)

type service struct {
	au    *AuthSys
	bmSys *store.BucketMetadataSys
}

func newService(bmSys *store.BucketMetadataSys) (svc *service, err error) {
	svc = &service{
		au:    NewAuthSys(),
		bmSys: bmSys,
	}
	return
}

func (s *service) CheckSignatureAndAcl(ctx context.Context, r *http.Request, action s3action.Action, bucketName string) (
	cred Credentials, s3Error apierrors.ErrorCode) {

	return s.au.CheckRequestAuthTypeCredential(ctx, r, action, bucketName, s.bmSys)
}
