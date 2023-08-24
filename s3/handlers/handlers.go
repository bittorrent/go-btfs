// Package handlers is an implementation of s3.Handlerser
package handlers

import (
	"context"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/ctxmu"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"github.com/bittorrent/go-btfs/s3/services/bucket"
	"github.com/bittorrent/go-btfs/s3/services/cors"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"github.com/bittorrent/go-btfs/s3/services/sign"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
)

const lockPrefix = "s3:lock/"

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

func (h *Handlers) name() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func (h *Handlers) rlock(ctx context.Context, key string, w http.ResponseWriter, r *http.Request) (runlock func(), err error) {
	key = lockPrefix + key
	ctx, cancel := context.WithTimeout(ctx, lockWaitTimeout)
	err = h.nslock.RLock(ctx, key)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		cancel()
		return
	}
	runlock = func() {
		h.nslock.RUnlock(key)
		cancel()
	}
	return
}

func (h *Handlers) lock(ctx context.Context, key string, w http.ResponseWriter, r *http.Request) (unlock func(), err error) {
	key = lockPrefix + key
	ctx, cancel := context.WithTimeout(ctx, lockWaitTimeout)
	err = h.nslock.Lock(ctx, key)
	if err != nil {
		responses.WriteErrorResponse(w, r, err)
		cancel()
		return
	}
	unlock = func() {
		h.nslock.Unlock(key)
		cancel()
	}
	return
}

// Parse object url queries
func (h *Handlers) getObjectResources(values url.Values) (uploadID string, partNumberMarker, maxParts int, encodingType string, rerr *responses.Error) {
	var err error
	if values.Get("max-parts") != "" {
		if maxParts, err = strconv.Atoi(values.Get("max-parts")); err != nil {
			rerr = responses.ErrInvalidMaxParts
			return
		}
	} else {
		maxParts = consts.MaxPartsList
	}

	if values.Get("part-number-marker") != "" {
		if partNumberMarker, err = strconv.Atoi(values.Get("part-number-marker")); err != nil {
			rerr = responses.ErrInvalidPartNumberMarker
			return
		}
	}

	uploadID = values.Get("uploadId")
	encodingType = values.Get("encoding-type")
	return
}
