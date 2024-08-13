package corehttp

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands"
	"github.com/bittorrent/go-btfs/utils"
	ds "github.com/ipfs/go-datastore"
	"net/http"
	"strings"
	"time"
)

const defaultTwoStepDuration = 30 * time.Minute

const firstStepUrl = "/api/v1/dashboard/validate"

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

	err = twoStepCheckInterceptor(r)
	if err != nil {
		return err
	}

	return nil
}

func twoStepCheckInterceptor(r *http.Request) error {
	if !need2StepCheckUrl(r.URL.Path) {
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
	if filterUrl(r) || filterP2pSchema(r) || filterLocalShellApi(r) || filterGatewayUrl(r) {
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
		APIPath + "/dashboard/check": true,
		APIPath + "/dashboard/set":   true,
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

	if need2StepCheckUrl(r.URL.Path) && currentStep == secondStep {
		currentStep = firstStep
	}

	return nil
}

func need2StepCheckUrl(path string) bool {
	if path == "/api/v1/xxx" {
		return true
	}
	return false
}