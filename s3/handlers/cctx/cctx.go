package cctx

import (
	"context"
	"net/http"
)

type key *struct{}

var (
	keyOfAccessKey = new(struct{})
	keyOfHandleInf = new(struct{})
)

func set(r *http.Request, k key, v any) {
	ctx := context.WithValue(r.Context(), k, v)
	nr := r.WithContext(ctx)
	*r = *nr
	return
}

func get(r *http.Request, k key) (v any) {
	v = r.Context().Value(k)
	return
}
