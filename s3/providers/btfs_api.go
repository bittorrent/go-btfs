package providers

import (
	"errors"
	shell "github.com/bittorrent/go-btfs-api"
	"io"
)

var _ FileStorer = (*BtfsAPI)(nil)

type BtfsAPI struct {
	shell *shell.Shell
}

func NewBtfsAPI(endpointUrl string) (api *BtfsAPI) {
	api = &BtfsAPI{}
	if endpointUrl == "" {
		api.shell = shell.NewLocalShell()
	} else {
		api.shell = shell.NewShell(endpointUrl)
	}
	return
}

func (api *BtfsAPI) Store(r io.Reader) (id string, err error) {
	id, err = api.shell.Add(r, shell.Pin(true))
	return
}

func (api *BtfsAPI) Remove(id string) (err error) {
	ok := api.shell.Remove(id)
	if !ok {
		err = errors.New("not removed")
	}
	return
}

func (api *BtfsAPI) Cat(id string) (rc io.ReadCloser, err error) {
	rc, err = api.shell.Cat(id)
	return
}
