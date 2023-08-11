package auth

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/services"
	s3action "github.com/bittorrent/go-btfs/s3d/action"
	"github.com/bittorrent/go-btfs/s3d/apierrors"
	"github.com/bittorrent/go-btfs/s3d/store"
	"net/http"
)

var _ handlers.SignatureService = (*Signature)(nil)

type Signature struct {
	providers    services.Providerser
	accesskeySvc handlers.AccessKeyService
	au           *AuthSys
	bmSys        *store.BucketMetadataSys
}

func NewSignature(providers services.Providerser, accesskeySvc handlers.AccessKeyService, options ...Option) (signature *Signature) {
	signature = &Signature{
		providers:    providers,
		accesskeySvc: accesskeySvc,
	}
	for _, option := range options {

	}
	return
}

func (s *service) CheckSignatureAndAcl(ctx context.Context, r *http.Request, action s3action.Action, bucketName string) (
	cred Credentials, s3Error apierrors.ErrorCode) {

	return s.au.CheckRequestAuthTypeCredential(ctx, r, action, bucketName, s.bmSys)
}
