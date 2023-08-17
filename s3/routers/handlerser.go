package routers

import (
	"net/http"
)

type Handlerser interface {
	Cors(handler http.Handler) http.Handler
	Sign(handler http.Handler) http.Handler

	PutBucketHandler(w http.ResponseWriter, r *http.Request)

	//PutObjectHandler(w http.ResponseWriter, r *http.Request)
}
