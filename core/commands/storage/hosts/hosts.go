package hosts

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	"github.com/bittorrent/go-btfs/core/hub"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	hubpb "github.com/tron-us/go-btfs-common/protos/hub"

	logging "github.com/ipfs/go-log"
)

var hostsLog = logging.Logger("storage/hosts")

const (
	hostInfoModeOptionName = "host-info-mode"
	hostSyncModeOptionName = "host-sync-mode"
)

var StorageHostsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Interact with information on hosts.",
		ShortDescription: `Allows interaction with information on hosts. Host information is synchronized from btfs-hub and saved in local datastore.`,
	},
	Subcommands: map[string]*cmds.Command{
		"info": storageHostsInfoCmd,
		"sync": storageHostsSyncCmd,
	},
}

var storageHostsInfoCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Display saved host information.",
		ShortDescription: `
This command displays saved information from btfs-hub under multiple modes.
Each mode ranks hosts based on its criteria and is randomized based on current node location.

Mode options include:` + hub.AllModeHelpText,
	},
	Options: []cmds.Option{
		cmds.StringOption(hostInfoModeOptionName, "m", "Hosts info showing mode. Default: mode set in config option Experimental.HostsSyncMode."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}
		if !cfg.Experimental.StorageClientEnabled {
			return fmt.Errorf("storage client api not enabled")
		}

		mode, ok := req.Options[hostInfoModeOptionName].(string)
		if !ok {
			mode = cfg.Experimental.HostsSyncMode
		}

		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		nodes, err := helper.GetHostsFromDatastore(req.Context, n, mode, 0)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &HostInfoRes{nodes})
	},
	Type: HostInfoRes{},
}

type HostInfoRes struct {
	Nodes []*hubpb.Host
}

var storageHostsSyncCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Synchronize host information from btfs-hub.",
		ShortDescription: `
This command synchronizes information from btfs-hub using multiple modes.
Each mode ranks hosts based on its criteria and is randomized based on current node location.

Mode options include:` + hub.AllModeHelpText,
	},
	Options: []cmds.Option{
		cmds.StringOption(hostSyncModeOptionName, "m", "Hosts syncing mode. Default: mode set in config option Experimental.HostsSyncMode."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}
		if !cfg.Experimental.StorageClientEnabled {
			return fmt.Errorf("storage client api not enabled")
		}

		mode, ok := req.Options[hostSyncModeOptionName].(string)
		if !ok {
			mode = cfg.Experimental.HostsSyncMode
		}

		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		_, err = SyncHostsMixture(req.Context, n, mode)
		return err
	},
}

func SyncHosts(ctx context.Context, node *core.IpfsNode, mode string) ([]*hubpb.Host, error) {
	nodes, err := hub.QueryHosts(ctx, node, mode)
	if err != nil {
		return nil, err
	}
	err = helper.SaveHostsIntoDatastore(ctx, node, mode, nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func SyncHostsMixture(ctx context.Context, node *core.IpfsNode, mode string) ([]*hubpb.Host, error) {
	if !json.Valid([]byte(mode)) {
		return SyncHosts(ctx, node, mode)
	}
	modes := map[string]int{}
	if err := json.Unmarshal([]byte(mode), &modes); err != nil {
		return nil, errors.Wrap(err, "invalid mode")
	}
	c := 0
	preMixing := map[string][]*hubpb.Host{}
	for k, v := range modes {
		if hosts, err := SyncHosts(ctx, node, k); err != nil {
			log.Error(err)
			continue
		} else {
			preMixing[k] = hosts
			c += v
		}
	}

	result := make([]*hubpb.Host, 0)
	for k, v := range preMixing {
		r := modes[k] * 500 / c
		if r > len(v) {
			r = len(v)
		}
		result = append(result, v[0:r]...)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })
	if err := helper.SaveHostsIntoDatastore(ctx, node, mode, result); err != nil {
		log.Error(err)
	}
	return result, nil
}
