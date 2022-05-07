package guide

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/markbates/pkger"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

const (
	pageFilePath       = "/hostui"
	pagePath           = "/hostui"
	infoPath           = "/api/v1/guide-info"
	serverCloseTimeout = 5 * time.Second
)

var (
	info               *Info
	serverAddr         string
	shutdownServerFunc func()
)

type Info struct {
	BtfsVersion string `json:"btfs_version"`
	HostID      string `json:"host_id"`
	BttcAddress string `json:"bttc_address"`
	PrivateKey  string `json:"private_key"`
}

func SetInfo(vinfo *Info) {
	info = vinfo
}

func SetServerAddr(cfgAddrs []string, optAddr string) {
	apiAddrs := cfgAddrs
	if optAddr != "" {
		apiAddrs = []string{optAddr}
	}
	for _, apiAddr := range apiAddrs {
		maddr, err := ma.NewMultiaddr(apiAddr)
		if err != nil {
			continue
		}
		network, addr, err := manet.DialArgs(maddr)
		if err != nil {
			continue
		}
		switch network {
		case "tcp", "tcp4", "tcp6":
			serverAddr = addr
			return
		}
	}
}

func newServer() *http.Server {
	mux := http.NewServeMux()

	// page
	page := http.StripPrefix(pageFilePath, http.FileServer(pkger.Dir(pageFilePath)))

	// info
	info := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]interface{}{
			"info": info,
		}
		encodeResp, _ := json.Marshal(resp)
		_, _ = w.Write(encodeResp)
	})

	// router
	mux.Handle(pagePath+"/", page)
	mux.Handle(infoPath+"/", info)

	return &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}
}

func StartServer() {
	server := newServer()
	done := make(chan struct{})
	go func() {
		fmt.Printf("Guide: http://%s%s\n", server.Addr, pagePath)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Start guide server failed: %v\n", err)
		}
		close(done)
	}()

	shutdownServerFunc = func() {
		select {
		case <-done:
			return // if the server has been closed, just return
		default:
			ctx, cancel := context.WithTimeout(context.Background(), serverCloseTimeout)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				fmt.Printf("Close guide server failed: %v\n", err)
			}
			<-done
		}
	}
}

func TryShutdownServer() {
	if shutdownServerFunc != nil {
		shutdownServerFunc()
	}
}
