package contexts

import (
	"context"
	"net/http"
)

type key string

const (
	keyOfAccessKey   key = "ctx-access-key"
	keyOfHandleInf   key = "ctx-handle-inf"
	keyOfRequestArgs key = "ctx-request-args"
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
