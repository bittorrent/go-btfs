package server

import "net/http"

type Routerser interface {
	Register() http.Handler
}
