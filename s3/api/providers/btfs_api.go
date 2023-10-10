package providers

import (
	shell "github.com/bittorrent/go-btfs-api"
	"github.com/mitchellh/go-homedir"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var _ FileStorer = (*BtfsAPI)(nil)

type BtfsAPI struct {
	shell        *shell.Shell
	headerTimout time.Duration
	timeout      time.Duration
	endpointUrl  string
}

func NewBtfsAPI(options ...BtfsAPIOption) (api *BtfsAPI, err error) {
	api = &BtfsAPI{
		headerTimout: defaultBtfsAPIResponseHeaderTimeout,
		timeout:      defaultBtfsAPITimeout,
		endpointUrl:  defaultBtfsAPIEndpointUrl,
	}
	for _, option := range options {
		option(api)
	}

	if api.endpointUrl == "" {
		api.endpointUrl, err = api.getLocalUrl()
		if err != nil {
			return
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: api.headerTimout,
		},
		Timeout: api.timeout,
	}

	api.shell = shell.NewShellWithClient(
		api.endpointUrl, client,
	)

	return
}

func (api *BtfsAPI) Store(r io.Reader) (id string, err error) {
	id, err = api.shell.Add(r, shell.Pin(true))
	return
}

func (api *BtfsAPI) Remove(id string) (err error) {
	err = api.shell.Unpin(id)
	return
}

func (api *BtfsAPI) Cat(id string) (rc io.ReadCloser, err error) {
	rc, err = api.shell.Cat(id)
	return
}

func (api *BtfsAPI) getLocalUrl() (url string, err error) {
	baseDir := os.Getenv(shell.EnvDir)
	if baseDir == "" {
		baseDir = shell.DefaultPathRoot
	}

	baseDir, err = homedir.Expand(baseDir)
	if err != nil {
		return
	}

	apiFile := path.Join(baseDir, shell.DefaultApiFile)

	_, err = os.Stat(apiFile)
	if err != nil {
		return
	}

	bs, err := os.ReadFile(apiFile)
	if err != nil {
		return
	}

	url = strings.TrimSpace(string(bs))
	return

}
