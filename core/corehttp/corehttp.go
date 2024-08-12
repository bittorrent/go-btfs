/*
Package corehttp provides utilities for the webui, gateways, and other
high-level HTTP interfaces to IPFS.
*/
package corehttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	core "github.com/bittorrent/go-btfs/core"
	logging "github.com/ipfs/go-log"
	"github.com/jbenet/goprocess"
	periodicproc "github.com/jbenet/goprocess/periodic"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

var log = logging.Logger("core/server")

// shutdownTimeout is the timeout after which we'll stop waiting for hung
// commands to return on shutdown.
const shutdownTimeout = 30 * time.Second

// ServeOption registers any HTTP handlers it provides on the given mux.
// It returns the mux to expose to future options, which may be a new mux if it
// is interested in mediating requests to future options, or the same mux
// initially passed in if not.
type ServeOption func(*core.IpfsNode, net.Listener, *http.ServeMux) (*http.ServeMux, error)

// makeHandler turns a list of ServeOptions into a http.Handler that implements
// all of the given options, in order.
func makeHandler(n *core.IpfsNode, l net.Listener, options ...ServeOption) (http.Handler, error) {
	topMux := http.NewServeMux()
	mux := topMux
	for _, option := range options {
		var err error
		mux, err = option(n, l, mux)
		if err != nil {
			return nil, err
		}
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ServeMux does not support requests with CONNECT method,
		// so we need to handle them separately
		// https://golang.org/src/net/http/request.go#L111
		if r.Method == http.MethodConnect {
			w.WriteHeader(http.StatusOK)
			return
		}

		err := interceptorBeforeReq(r, n)
		if err != nil {
			// set allow origin
			w.Header().Set("Access-Control-Allow-Origin", "*")
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Headers", "X-Stream-Output, X-Chunked-Output, X-Content-Length")
				w.WriteHeader(http.StatusOK)
				return
			}
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		topMux.ServeHTTP(w, r)

		err = interceptorAfterResp(r, w, n)
		if err != nil {
			return
		}

	})
	return handler, nil
}

// ListenAndServe runs an HTTP server listening at |listeningMultiAddr| with
// the given serve options. The address must be provided in multiaddr format.
//
// TODO intelligently parse address strings in other formats so long as they
// unambiguously map to a valid multiaddr. e.g. for convenience, ":8080" should
// map to "/ip4/0.0.0.0/tcp/8080".
func ListenAndServe(n *core.IpfsNode, listeningMultiAddr string, options ...ServeOption) error {
	addr, err := ma.NewMultiaddr(listeningMultiAddr)
	if err != nil {
		return err
	}

	list, err := manet.Listen(addr)
	if err != nil {
		return err
	}

	// we might have listened to /tcp/0 - let's see what we are listing on
	addr = list.Multiaddr()
	fmt.Printf("API server listening on %s\n", addr)

	return Serve(n, manet.NetListener(list), options...)
}

// Serve accepts incoming HTTP connections on the listener and pass them
// to ServeOption handlers.
func Serve(node *core.IpfsNode, lis net.Listener, options ...ServeOption) error {
	// make sure we close this no matter what.
	defer lis.Close()

	handler, err := makeHandler(node, lis, options...)
	if err != nil {
		return err
	}

	addr, err := manet.FromNetAddr(lis.Addr())
	if err != nil {
		return err
	}

	select {
	case <-node.Process.Closing():
		return fmt.Errorf("failed to start server, process closing")
	default:
	}

	server := &http.Server{
		Handler: handler,
	}

	var serverError error
	serverProc := node.Process.Go(func(p goprocess.Process) {
		serverError = server.Serve(lis)
	})

	// wait for server to exit.
	select {
	case <-serverProc.Closed():
	// if node being closed before server exits, close server
	case <-node.Process.Closing():
		log.Infof("server at %s terminating...", addr)

		warnProc := periodicproc.Tick(5*time.Second, func(_ goprocess.Process) {
			log.Infof("waiting for server at %s to terminate...", addr)
		})

		// This timeout shouldn't be necessary if all of our commands
		// are obeying their contexts but we should have *some* timeout.
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		err := server.Shutdown(ctx)

		// Should have already closed but we still need to wait for it
		// to set the error.
		<-serverProc.Closed()
		serverError = err

		warnProc.Close()
	}

	log.Infof("server at %s terminated", addr)
	return serverError
}