// Package handlers is an implementation of Handlerser
package handlers

import (
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"github.com/bittorrent/go-btfs/s3/services/sign"
	"net/url"
	"runtime"
	"strconv"
)

var _ Handlerser = (*Handlers)(nil)

type Handlers struct {
	headers map[string][]string
	acksvc  accesskey.Service
	sigsvc  sign.Service
	objsvc  object.Service
}

func NewHandlers(acksvc accesskey.Service, sigsvc sign.Service, objsvc object.Service, options ...Option) (handlers *Handlers) {
	handlers = &Handlers{
		headers: defaultHeaders,
		acksvc:  acksvc,
		sigsvc:  sigsvc,
		objsvc:  objsvc,
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

// Parse object url queries
func (h *Handlers) getObjectResources(values url.Values) (uploadId string, partNumberMarker, maxParts int, encodingType string, rerr *responses.Error) {
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

	uploadId = values.Get("uploadId")
	encodingType = values.Get("encoding-type")
	return
}
