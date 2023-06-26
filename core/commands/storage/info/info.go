package info

import (
	"fmt"
	"github.com/bittorrent/go-btfs/utils"

	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"

	cmds "github.com/bittorrent/go-btfs-cmds"
	nodepb "github.com/bittorrent/go-btfs-common/protos/node"
)

var StorageInfoCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Show storage host information.",
		ShortDescription: `
This command displays host information synchronized from the BTFS network.
By default it shows local host node information.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("peer-id", false, false, "Peer ID to show storage-related information. Default to self."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}
		if len(req.Arguments) > 0 {
			if !cfg.Experimental.StorageClientEnabled {
				return fmt.Errorf("storage client api not enabled")
			}
		}

		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		// Default to self
		var data *nodepb.Node_Settings
		var peerID string
		if len(req.Arguments) > 0 {
			peerID = req.Arguments[0]
			data, err = helper.GetHostStorageConfigForPeer(req.Context, n, peerID)
		} else {
			data, err = helper.GetHostStorageConfig(req.Context, n)
		}
		if err != nil {
			return err
		}
		roles := make([]nodepb.NodeRole, 0)
		if cfg.Experimental.StorageClientEnabled {
			roles = append(roles, nodepb.NodeRole_RENTER)
		}
		if cfg.Experimental.StorageHostEnabled {
			roles = append(roles, nodepb.NodeRole_HOST)
		}
		if cfg.Experimental.HostRepairEnabled {
			roles = append(roles, nodepb.NodeRole_REPAIRER)
		}
		if cfg.Experimental.HostChallengeEnabled {
			roles = append(roles, nodepb.NodeRole_CHALLENGER)
		}
		data.Roles = roles
		return cmds.EmitOnce(res, data)
	},
	Type: nodepb.Node_Settings{},
}
