package vault

import (
	"fmt"
	"io"
	"math/big"
	"time"

	cmds "github.com/TRON-US/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"golang.org/x/net/context"
)

type VaultBalanceCmdRet struct {
	Balance *big.Int `json:"balance"`
}

var VaultBalanceCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get vault balance.",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		balance, err := chain.SettleObject.VaultService.AvailableBalance(context.Background())
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &VaultBalanceCmdRet{
			Balance: balance,
		})
	},
	Type: &VaultBalanceCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *VaultBalanceCmdRet) error {
			_, err := fmt.Fprintf(w, "the vault available balance: %v\n", out.Balance)
			return err
		}),
	},
}
