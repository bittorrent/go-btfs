package spin

import (
	"context"
	"time"

	"github.com/bittorrent/go-btfs/core/commands/storage/contracts"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/core"
)

const (
	hostContractsSyncPeriod  = 24 * 60 * time.Minute
	hostContractsSyncTimeout = 10 * time.Minute
)

func Contracts(n *core.IpfsNode, req *cmds.Request, env cmds.Environment, role string) {
	cfg, err := n.Repo.Config()
	if err != nil {
		log.Errorf("Failed to get configuration %s", err)
		return
	}
	if cfg.Experimental.StorageHostEnabled {
		go periodicSync(hostContractsSyncPeriod, hostContractsSyncTimeout, role+" contracts",
			func(ctx context.Context) error {
				return contracts.SyncContracts(ctx, n, req, env, role)
			})
	}
}
