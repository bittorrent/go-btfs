package main

import (
	"context"
	"encoding/json"
	"fmt"
	config "github.com/TRON-US/go-btfs-config"
	cmds "github.com/bittorrent/go-btfs-cmds"
	oldcmds "github.com/bittorrent/go-btfs/commands"
	"github.com/bittorrent/go-btfs/core/commands"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"html/template"
	"net/http"
	"time"
)

type guideInfo struct {
	BtfsVersion string `json:"btfs_version"`
	HostID      string `json:"host_id"`
	BttcAddress string `json:"bttc_address"`
	PrivateKey  string `json:"private_key"`
}

// TODO: replace it
const guidePagePath = "./guide-page/index.html"

func startGuideServer(req *cmds.Request, cctx *oldcmds.Context, info *guideInfo) (closeFunc func()) {
	cfg, err := cctx.GetConfig()
	if err != nil {
		fmt.Printf("Start guide server: get config: %v\n", err)
		return
	}

	// get guide server address
	optionAddr, _ := req.Options[commands.ApiOption].(string)
	guideAddr := getGuideServerAddr(cfg, optionAddr)
	if guideAddr == "" {
		fmt.Println("Start guide server failed: no valid address")
		return
	}

	// guide info api
	guideDataFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]interface{}{
			"ret":  0,
			"desc": "ok",
			"info": info,
		}
		encodeResp, _ := json.Marshal(resp)
		w.Write(encodeResp)
		return
	})

	// guide page
	guidePageFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(guidePagePath)
		if err != nil {
			w.Write([]byte("parse file failed: " + err.Error()))
			fmt.Printf("Guide page parse file failed: %v\n", err)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			fmt.Printf("Guide page template execute failed: %v\n", err)
		}
	})

	// guide server setup
	handler := &http.ServeMux{}
	handler.Handle("/guide-info", guideDataFunc)
	handler.Handle("/hostui", guidePageFunc)
	server := &http.Server{
		Addr:    guideAddr,
		Handler: handler,
	}
	done := make(chan struct{})

	// start server
	fmt.Printf("Start guide server at: %s\n", guideAddr)
	fmt.Printf("Guide: http://%s/hostui\n", guideAddr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Start guide server failed: %v\n", err)
		}
		close(done)
	}()

	// close function, it is idempotent
	closeFunc = func() {
		select {
		case <-done: // if closed, do nothing
		default:
			fmt.Println("Closing guide server...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				fmt.Printf("Close guide server failed: %v\n", err)
			}
			<-done
		}
		fmt.Println("Guide server closed")
	}

	return
}

// getGuideServerAddr find a valid api address for guide server
func getGuideServerAddr(cfg *config.Config, optionAddr string) string {
	apiAddrs := cfg.Addresses.API
	if optionAddr != "" {
		apiAddrs = []string{optionAddr}
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
			return addr
		}
	}
	return ""
}
