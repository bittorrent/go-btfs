package upload

import (
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
)

var StorageUploadSupportTokensCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "support cheque, return tokens.",
		ShortDescription: `support cheque, return tokens.`,
	},
	Arguments:  []cmds.Argument{},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}
		if !ctxParams.Cfg.Experimental.StorageHostEnabled {
			return fmt.Errorf("storage host api not enabled")
		}

		return cmds.EmitOnce(res, &tokencfg.MpTokenAddr)
	},
}
