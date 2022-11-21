package upload

import (
	"encoding/json"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/ethereum/go-ethereum/common"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
)

var StorageUploadCheckTokensCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "receive upload cheque, do with cheque, and return it.",
		ShortDescription: `receive upload cheque, deal it and return it.`,
	},
	Arguments:  []cmds.Argument{},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}
		if !ctxParams.Cfg.Experimental.StorageHostEnabled {
			return fmt.Errorf("storage host api not enabled")
		}

		tokens := make([]common.Address, 0)
		for _, v := range tokencfg.MpTokenAddr {
			tokens = append(tokens, v)
		}

		fmt.Println("receive check tokens, ", tokens)
		output, err := json.Marshal(tokens)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &output)
	},
}
