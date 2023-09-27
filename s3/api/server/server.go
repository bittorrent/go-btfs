package server

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/api/routers"
	"net/http"
	"sync"
)

var (
	ErrServerStarted    = errors.New("server started")
	ErrServerNotStarted = errors.New("server not started")
)

type Server struct {
	routers routers.Routerser
	address string

	shutdown func() error
	mutex    sync.Mutex
}

func NewServer(routers routers.Routerser, options ...Option) (s *Server) {
	s = &Server{
		routers:  routers,
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
		Handler: s.routers.Register(),
	}

	s.shutdown = func() error {
		return httpSvr.Shutdown(context.TODO())
	}

	go func() {
		_ = httpSvr.ListenAndServe()
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
	return
}
