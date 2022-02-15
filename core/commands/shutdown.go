package commands

import (
	cmdenv "github.com/bittorrent/go-btfs/core/commands/cmdenv"

	"github.com/bittorrent/go-btfs-cmds"
)

var daemonShutdownCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Shut down the btfs daemon",
	},
	Run: func(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		if !nd.IsDaemon {
			return cmds.Errorf(cmds.ErrClient, "daemon not running")
		}

		if err := nd.Close(); err != nil {
			log.Error("error while shutting down btfs daemon:", err)
		}

		return nil
	},
}
