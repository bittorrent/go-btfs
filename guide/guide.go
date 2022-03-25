package guide

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/markbates/pkger"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

const (
	pageFilePath       = "/hostui/"
	pagePath           = "/hostui/"
	infoPath           = "/api/v1/guide-info/"
	serverCloseTimeout = 5 * time.Second
)

const (
	BalanceStatusNotOk = iota
	BalanceStatusOK
)

var (
	infoVal            *Info
	serverAddr         string
	shutdownServerFunc func()
	balanceStatus      int
)

type Info struct {
	BtfsVersion string `json:"btfs_version"`
	HostID      string `json:"host_id"`
	BttcAddress string `json:"bttc_address"`
	PrivateKey  string `json:"private_key"`
}

func SetInfoVal(val *Info) {
	infoVal = val
}

func SetBalanceStatusOK() {
	balanceStatus = BalanceStatusOK
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
	static := http.StripPrefix(pageFilePath, http.FileServer(pkger.Dir(pageFilePath)))
	mux.HandleFunc(pagePath, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pagePath {
			indexPath := path.Join(pagePath, "index.html")
			f, err := pkger.Open(indexPath)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			http.ServeContent(w, r, "index.html", time.Now(), f)
			return
		}
		static.ServeHTTP(w, r)
	})

	// api
	mux.HandleFunc(infoPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]interface{}{
			"info":           infoVal,
			"balance_status": balanceStatus,
		}
		encodeResp, _ := json.Marshal(resp)
		_, _ = w.Write(encodeResp)
		return
	})

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
	return
}

func TryShutdownServer() {
	if shutdownServerFunc == nil {
		return
	}
	shutdownServerFunc()
}
