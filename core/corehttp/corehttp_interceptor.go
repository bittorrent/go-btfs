package corehttp

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands"
	"github.com/bittorrent/go-btfs/utils"
	ds "github.com/ipfs/go-datastore"
)

const defaultTwoStepDuration = 30 * time.Minute

const firstStepUrl = "dashboard/validate"

var (
	ErrorGatewayCidExits = errors.New("cid exits")
)

func interceptorBeforeReq(r *http.Request, n *core.IpfsNode) error {
	config, err := n.Repo.Config()
	if err != nil {
		return err
	}

	if config.API.EnableTokenAuth {
		err = tokenCheckInterceptor(r, n)
		if err != nil {
			return err
		}

		err = twoStepCheckInterceptor(r)
		if err != nil {
			return err
		}
	}

	exits, err := gatewayCidInterceptor(r, n)
	if err != nil || !exits {
		return nil
	}

	if exits {
		return ErrorGatewayCidExits
	}

	return nil
}

func twoStepCheckInterceptor(r *http.Request) error {
	if !needTwoStepCheckUrl(r.Method, r.URL.Path) {
		return nil
	}
	if currentStep == secondStep {
		return nil
	}

	return errors.New("please validate your password first")
}

func interceptorAfterResp(r *http.Request, w http.ResponseWriter, n *core.IpfsNode) error {
	err := passwordCheckInterceptor(r)
	if err != nil {
		return err
	}
	return nil
}

func tokenCheckInterceptor(r *http.Request, n *core.IpfsNode) error {
	conf, err := n.Repo.Config()
	if err != nil {
		return err
	}
	apiHost := fmt.Sprint(strings.Split(conf.Addresses.API[0], "/")[2], ":", strings.Split(conf.Addresses.API[0], "/")[4])
	if filterNoNeedTokenCheckReq(r, apiHost, conf.Identity.PeerID) {
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

func gatewayCidInterceptor(r *http.Request, n *core.IpfsNode) (bool, error) {
	if filterGatewayUrl(r) {
		sPath := strings.Split(r.URL.Path, "/")
		if len(sPath) < 3 {
			return false, nil
		}
		key := strings.Split(r.URL.Path, "/")[2]
		exits, err := n.Repo.Datastore().Has(r.Context(), ds.NewKey(commands.NewGatewayFilterKey(key)))
		return exits, err
	}
	return false, nil
}

func filterNoNeedTokenCheckReq(r *http.Request, apiHost string, peerId string) bool {
	if filterUrl(r) || filterP2pSchema(r, peerId) || filterLocalShellApi(r, apiHost) || filterGatewayUrl(r) {
		return true
	}
	return false
}

func filterGatewayUrl(r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, "/btfs/") || strings.HasPrefix(r.URL.Path, "/btns/") {
		return true
	}
	return false
}

func filterUrl(r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, "/dashboard") {
		return true
	}
	if strings.HasPrefix(r.URL.Path, "/hostui") {
		return true
	}
	urls := map[string]bool{
		APIPath + "/id":              true,
		APIPath + "/config/show":     true,
		APIPath + "/dashboard/check": true,
		APIPath + "/dashboard/set":   true,
		APIPath + "/dashboard/login": true,
		APIPath + "/dashboard/reset": true,
	}

	return urls[r.URL.Path]
}

const (
	defaultUserAgent = "Go-http-client/1.1"
	cmdUserAgent     = "go-btfs-cmds/http"
)

func filterLocalShellApi(r *http.Request, apiHost string) bool {
	host := r.Host
	ua := r.Header.Get("User-Agent")
	// ua is not Go-http-client
	if host == apiHost && (ua == defaultUserAgent || ua == cmdUserAgent) {
		return true
	}
	return false
}

func filterP2pSchema(r *http.Request, peerId string) bool {
	if r.URL.Scheme == "libp2p" {
		return true
	}
	if r.Host == peerId {
		return true
	}
	return false
}

const (
	_ int = iota
	firstStep
	secondStep
)

var currentStep = firstStep

func passwordCheckInterceptor(r *http.Request) error {
	if r.URL.Path == firstStepUrl && currentStep == firstStep {
		currentStep = secondStep
		// if next step was done after the default duration
		go func() {
			<-time.After(defaultTwoStepDuration)
			currentStep = firstStep
		}()
		return nil
	}

	if needTwoStepCheckUrl(r.Method, r.URL.Path) && currentStep == secondStep {
		currentStep = firstStep
	}

	return nil
}

func needTwoStepCheckUrl(method string, path string) bool {
	if method == http.MethodOptions {
		return false
	}
	urls := map[string]bool{
		APIPath + "/bttc/send-btt-to":   true,
		APIPath + "/bttc/send-wbtt-to":  true,
		APIPath + "/bttc/send-token-to": true,
	}
	if urls[path] {
		return true
	}
	return false
}
