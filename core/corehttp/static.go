package corehttp

import (
	"net"
	"net/http"

	"github.com/bittorrent/go-btfs/core"

	"github.com/markbates/pkger"
)

func StaticOption(path, dir string) ServeOption {
	return func(n *core.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		mux.Handle("/"+path+"/", http.StripPrefix(dir, http.FileServer(pkger.Dir(dir))))
		return mux, nil
	}
}
