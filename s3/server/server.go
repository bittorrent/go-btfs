package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/s3/routers"
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
		fmt.Printf("Start s3-compatible-api server, endpoint-url: http://%s\n", httpSvr.Addr)
		lErr := httpSvr.ListenAndServe()
		fmt.Printf("Stop s3-compatible-api server: %v\n", lErr)
	}()

	return
}

func (s *Server) Stop() (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.shutdown == nil {
		err = ErrServerStarted
		return
	}
	err = s.shutdown()
	s.shutdown = nil
	return
}
