package spin

import (
	"context"
	"fmt"
	"time"

	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/hosts"
	"github.com/bittorrent/go-btfs/core/commands/storage/stats"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/core"
)

const (
	hostSyncPeriod         = 2 * 60 * time.Minute
	hostStatsSyncPeriod    = 2 * 60 * time.Minute
	hostSettingsSyncPeriod = 2 * 60 * time.Minute
	hostSyncTimeout        = 30 * time.Second
	hostSortTimeout        = 5 * time.Minute
)

func Hosts(node *core.IpfsNode, env cmds.Environment) {
	cfg, err := node.Repo.Config()
	if err != nil {
		log.Errorf("Failed to get configuration %s", err)
		return
	}

	if cfg.Experimental.HostsSyncEnabled {
		m := cfg.Experimental.HostsSyncMode
		fmt.Printf("Storage host info will be synced at [%s] mode\n", m)
		go periodicSync(hostSyncPeriod, hostSyncTimeout+hostSortTimeout, "hosts",
			func(ctx context.Context) error {
				_, err := hosts.SyncHostsMixture(ctx, node, m)
				return err
			})
	}
	if cfg.Experimental.StorageHostEnabled {
		fmt.Println("Current host stats will be synced")
		go periodicSync(hostStatsSyncPeriod, hostSyncTimeout, "host stats",
			func(ctx context.Context) error {
				return stats.SyncStats(ctx, cfg, node, env, true)
			})
		fmt.Println("Current host settings will be synced")
		go periodicSync(hostSettingsSyncPeriod, hostSyncTimeout, "host settings",
			func(ctx context.Context) error {
				_, err = helper.GetHostStorageConfigHelper(ctx, node, true)
				return err
			})
	}
}
