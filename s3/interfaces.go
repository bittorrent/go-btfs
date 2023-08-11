package s3

import (
	"net/http"
)

type Handlerser interface {
	Cors(handler http.Handler) http.Handler
	PutObjectHandler(w http.ResponseWriter, r *http.Request)
}
