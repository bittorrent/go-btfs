package handlers

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/s3/api/contexts"
	"github.com/bittorrent/go-btfs/s3/api/responses"
	"github.com/bittorrent/go-btfs/s3/api/services/accesskey"
	"github.com/bittorrent/go-btfs/s3/consts"
	rscors "github.com/rs/cors"
	"net/http"
	"strconv"
	"time"
)

func (h *Handlers) Cors(handler http.Handler) http.Handler {
	headers := h.headers
	cred := headers.Get(consts.AccessControlAllowCredentials) == "true"
	maxAge, _ := strconv.Atoi(headers.Get(consts.AccessControlMaxAge))
	ch := rscors.New(rscors.Options{
		AllowedOrigins:   headers.Values(consts.AccessControlAllowOrigin),
		AllowedMethods:   headers.Values(consts.AccessControlAllowMethods),
		AllowedHeaders:   headers.Values(consts.AccessControlAllowHeaders),
		ExposedHeaders:   headers.Values(consts.AccessControlExposeHeaders),
		MaxAge:           maxAge,
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
		start := time.Now()
		handler.ServeHTTP(w, r)
		hname, args, herr := contexts.GetHandleInf(r)
		end := time.Now()
		ela := end.Sub(start)
		fmt.Printf("s3-api: | %s | <%-4s> | %s | %s | %+v | %v |  %s \n", end.Format(time.RFC3339), r.Method, r.URL, hname, args, herr, ela)
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
				contexts.SetHandleInf(r, "Sign", err, nil)
			}
		}()

		ack, err := h.sigsvc.VerifyRequestSignature(r)
		if err != nil {
			responses.WriteErrorResponse(w, r, err)
			return
		}

		contexts.SetAccessKey(r, ack)

		handler.ServeHTTP(w, r)
	})
}
