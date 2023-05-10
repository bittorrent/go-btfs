package vault

import (
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"io"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"golang.org/x/net/context"
)

var VaultBalanceAllCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get vault balance.",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		mp := make(map[string]*big.Int, 0)
		for k, tokenAddr := range tokencfg.MpTokenAddr {
			balance, err := chain.SettleObject.VaultService.AvailableBalance(context.Background(), tokenAddr)
			if err != nil {
				return err
			}

			mp[k] = balance
		}

		return cmds.EmitOnce(res, &mp)
	},
	Type: &VaultBalanceCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *VaultBalanceCmdRet) error {
			_, err := fmt.Fprintf(w, "the vault available balance: %v\n", out.Balance)
			return err
		}),
	},
}
