// Package handlers is an implementation of s3.Handlerser
package handlers

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/ctxmu"
	"github.com/bittorrent/go-btfs/s3/requests"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"github.com/bittorrent/go-btfs/s3/services/bucket"
	"github.com/bittorrent/go-btfs/s3/services/cors"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"github.com/bittorrent/go-btfs/s3/services/sign"
	"net/http"
	"runtime"

	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/s3utils"
	rscors "github.com/rs/cors"
)

var _ Handlerser = (*Handlers)(nil)

type Handlers struct {
	corsvc cors.Service
	acksvc accesskey.Service
	sigsvc sign.Service
	bucsvc bucket.Service
	objsvc object.Service
	nslock ctxmu.MultiCtxRWLocker
}

func NewHandlers(
	corsvc cors.Service,
	acksvc accesskey.Service,
	sigsvc sign.Service,
	bucsvc bucket.Service,
	objsvc object.Service,
	options ...Option,
) (handlers *Handlers) {
	handlers = &Handlers{
		corsvc: corsvc,
		acksvc: acksvc,
		sigsvc: sigsvc,
		bucsvc: bucsvc,
		objsvc: objsvc,
		nslock: ctxmu.NewDefaultMultiCtxRWMutex(),
	}
	for _, option := range options {
		option(handlers)
	}
	return
}

func (h *Handlers) Cors(handler http.Handler) http.Handler {
	return rscors.New(rscors.Options{
		AllowedOrigins:   h.corsvc.GetAllowOrigins(),
		AllowedMethods:   h.corsvc.GetAllowMethods(),
		AllowedHeaders:   h.corsvc.GetAllowHeaders(),
		ExposedHeaders:   h.corsvc.GetAllowHeaders(),
		AllowCredentials: true,
	}).Handler(handler)
}

func (h *Handlers) Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[REQ] <%4s> | %s\n", r.Method, r.URL)
		handler.ServeHTTP(w, r)
		hname, herr := cctx.GetHandleInf(r)
		fmt.Printf("[RSP] <%4s> | %s | %s | %v\n", r.Method, r.URL, hname, herr)
	})
}

func (h *Handlers) Sign(handler http.Handler) http.Handler {
	h.sigsvc.SetSecretGetter(func(key string) (secret string, exists, enable bool, err error) {
		ack, err := h.acksvc.Get(key)
		if errors.Is(err, accesskey.ErrNotFound) {
			exists = false
			enable = true
			err = nil
			return
		}
		if err != nil {
			return
		}
		exists = true
		secret = ack.Secret
		enable = ack.Enable
		return
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err *responses.Error
		defer func() {
			if err != nil {
				cctx.SetHandleInf(r, h.name(), err)
			}
		}()

		ack, err := h.sigsvc.VerifyRequestSignature(r)
		if err != nil {
			responses.WriteErrorResponse(w, r, err)
			return
		}

		cctx.SetAccessKey(r, ack)

		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		cctx.SetHandleInf(r, h.name(), err)
	}()

	req, err := requests.ParsePutBucketRequest(r)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestBody)
		return
	}

	// issue: lock for check
	ctx := r.Context()
	ack := cctx.GetAccessKey(r)

	if err = s3utils.CheckValidBucketNameStrict(req.Bucket); err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInvalidBucketName)
		return
	}

	if !requests.CheckAclPermissionType(&req.ACL) {
		err = responses.ErrNotImplemented
		responses.WriteErrorResponse(w, r, responses.ErrNotImplemented)
		return
	}

	if ok := h.bucsvc.HasBucket(ctx, req.Bucket); ok {
		err = responses.ErrBucketAlreadyExists
		responses.WriteErrorResponseHeadersOnly(w, r, responses.ErrBucketAlreadyExists)
		return
	}

	err = h.bucsvc.CreateBucket(ctx, req.Bucket, req.Region, ack, req.ACL)
	if err != nil {
		responses.WriteErrorResponse(w, r, responses.ErrInternalError)
		return
	}

	// Make sure to add Location information here only for bucket
	if cp := requests.PathClean(r.URL.Path); cp != "" {
		w.Header().Set(consts.Location, cp) // Clean any trailing slashes.
	}

	responses.WritePutBucketResponse(w, r)

	return
}

func (h *Handlers) name() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
