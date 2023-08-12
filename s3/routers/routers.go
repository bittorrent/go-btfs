package routers

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Routers struct {
	handlers Handlerser
}

func NewRouters(handlers Handlerser, options ...Option) (routers *Routers) {
	routers = &Routers{
		handlers: handlers,
	}
	for _, option := range options {
		option(routers)
	}
	return
}

func (routers *Routers) Register() http.Handler {
	root := mux.NewRouter()

	root.Use(routers.handlers.Cors, routers.handlers.Sign)

	bucket := root.PathPrefix("/{bucket}").Subrouter()
	bucket.Methods(http.MethodPut).Path("/{object:.+}").HandlerFunc(routers.handlers.PutObjectHandler)

	return root
}
