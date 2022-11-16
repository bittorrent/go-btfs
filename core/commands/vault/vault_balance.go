package vault

import (
	"fmt"
	"io"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
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
	Arguments: []cmds.Argument{
		cmds.StringArg("token", true, false, "token"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		token := req.Arguments[0]
		fmt.Printf("... token:%+v\n", token)

		balance, err := chain.SettleObject.VaultService.AvailableBalance(context.Background(), token)
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
