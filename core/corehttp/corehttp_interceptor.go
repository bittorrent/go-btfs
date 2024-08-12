package corehttp

import (
	"fmt"
	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands"
	"github.com/bittorrent/go-btfs/utils"
	ds "github.com/ipfs/go-datastore"
	"net/http"
)

func interceptorBeforeReq(r *http.Request, n *core.IpfsNode) error {
	config, err := n.Repo.Config()
	if err != nil {
		return err
	}

	if config.API.EnableTokenAuth {
		err := tokenCheckInterceptor(r, n)
		if err != nil {
			return err
		}
	}

	return nil
}

func tokenCheckInterceptor(r *http.Request, n *core.IpfsNode) error {
	if filterNoNeedTokenCheckReq(r) {
		return nil
	}
	if !commands.IsLogin {
		return fmt.Errorf("please login")
	}
	args := r.URL.Query()
	token := args.Get("token")
	password, err := n.Repo.Datastore().Get(r.Context(), ds.NewKey(commands.DashboardPasswordPrefix))
	if err != nil {
		return err
	}
	claims, err := utils.VerifyToken(token, string(password))
	if err != nil {
		return err
	}
	if claims.PeerId != n.Identity.String() {
		return fmt.Errorf("token is invalid")
	}

	return nil
}

func filterNoNeedTokenCheckReq(r *http.Request) bool {
	if filterUrl(r) || filterP2pSchema(r) || filterLocalShellApi(r) {
		return true
	}
	return false
}

func filterUrl(r *http.Request) bool {
	urls := map[string]bool{
		// local
		"/dashboard": true,
		"/hostui":    true,
		// no need url
		APIPath + "/id":              true,
		APIPath + "/dashboard/check": true,
		APIPath + "/dashboard/login": true,
		APIPath + "/dashboard/reset": true,
	}

	return urls[r.URL.Path]
}

const defaultUserAgent = "Go-http-client/1.1"

func filterLocalShellApi(r *http.Request) bool {
	host := r.Host
	ua := r.Header.Get("User-Agent")
	// ua is not Go-http-client
	if host == "127.0.0.1:5001" && ua == defaultUserAgent {
		return true
	}
	return false
}

func filterP2pSchema(r *http.Request) bool {
	if r.URL.Scheme == "libp2p" {
		return true
	}
	return false
}