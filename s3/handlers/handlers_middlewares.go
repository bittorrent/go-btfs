package handlers

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	rscors "github.com/rs/cors"
	"net/http"
)

func (h *Handlers) Cors(handler http.Handler) http.Handler {
	headers := h.headers
	cred := len(headers[consts.AccessControlAllowCredentials]) > 0 &&
		headers[consts.AccessControlAllowCredentials][0] == "true"
	ch := rscors.New(rscors.Options{
		AllowedOrigins:   headers[consts.AccessControlAllowOrigin],
		AllowedMethods:   headers[consts.AccessControlAllowMethods],
		AllowedHeaders:   headers[consts.AccessControlExposeHeaders],
		ExposedHeaders:   headers[consts.AccessControlAllowHeaders],
		AllowCredentials: cred,
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// add all user headers
		for k, v := range h.headers {
			w.Header()[k] = v
		}
		// next
		ch.Handler(handler).ServeHTTP(w, r)
	})
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