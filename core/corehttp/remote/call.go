package remote

import (
	"context"

	"github.com/bittorrent/go-btfs-api"
	logging "github.com/ipfs/go-log"
)

type Call interface {
	CallGet(context.Context, string, []string) ([]byte, error)
	CallPost()
}

var log = logging.Logger("corehttp/remote")

// Must exist here to avoid circular dependency from
// core/commands -> core/corehttp/remote -> core/commands
const apiPrefix = "/api/" + shell.API_VERSION
