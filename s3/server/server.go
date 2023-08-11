package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

const defaultServerAddress = ":15001"

var (
	ErrServerStarted    = errors.New("server started")
	ErrServerNotStarted = errors.New("server not started")
)

type Server struct {
	handlers Handlerser
	address  string
	shutdown func() error
	mutex    sync.Mutex
}

func NewServer(handlers Handlerser, options ...Option) (s *Server) {
	s = &Server{
		handlers: handlers,
		address:  defaultServerAddress,
		shutdown: nil,
		mutex:    sync.Mutex{},
	}
	for _, option := range options {
		option(s)
	}
	return
}

func (s *Server) Start() (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.shutdown != nil {
		err = ErrServerStarted
		return
	}

	httpSvr := &http.Server{
		Addr:    s.address,
		Handler: s.registerRouter(),
	}

	s.shutdown = func() error {
		return httpSvr.Shutdown(context.TODO())
	}

	go func() {
		fmt.Printf("start s3-compatible-api server\n")
		lErr := httpSvr.ListenAndServe()
		if lErr != nil && !errors.Is(lErr, http.ErrServerClosed) {
			fmt.Printf("start s3-compatible-api server: %v\n", lErr)
		}
	}()

	return
}

func (s *Server) Stop() (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.shutdown == nil {
		err = ErrServerNotStarted
		return
	}
	err = s.shutdown()
	s.shutdown = nil
	fmt.Printf("stoped s3-compatible-api server: %v\n", err)
	return
}

func (s *Server) registerRouter() http.Handler {
	root := mux.NewRouter()

	root.Use(s.handlers.Cors, s.handlers.Sign)

	bucket := root.PathPrefix("/{bucket}").Subrouter()
	bucket.Methods(http.MethodPut).Path("/{object:.+}").HandlerFunc(s.handlers.PutObjectHandler)

	return root
}
