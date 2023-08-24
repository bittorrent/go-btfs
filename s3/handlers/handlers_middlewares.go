package handlers

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/s3/cctx"
	"github.com/bittorrent/go-btfs/s3/responses"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	rscors "github.com/rs/cors"
	"net/http"
)

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
