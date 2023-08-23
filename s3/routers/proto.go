package routers

import "net/http"

type Routerser interface {
	Register() http.Handler
}
