package spin

import (
	"context"
	"time"

	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/wallet"

	cmds "github.com/TRON-US/go-btfs-cmds"
)

const (
	bttTxSyncPeriod  = 2 * 60 * time.Minute
	bttTxSyncTimeout = 1 * time.Minute
)

func BttTransactions(node *core.IpfsNode, env cmds.Environment) {
	cfg, err := node.Repo.Config()
	if err != nil {
		log.Errorf("Failed to get configuration %s", err)
		return
	}

	go periodicSync(bttTxSyncPeriod, bttTxSyncTimeout, "btt-tx",
		func(ctx context.Context) error {
			_, err := wallet.SyncTxFromTronGrid(ctx, cfg, node.Repo.Datastore())
			return err
		})
}
