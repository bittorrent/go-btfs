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

type VaultWithdrawCmdRet struct {
	Hash string `json:"hash"`
}

var VaultWithdrawCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Withdraw from vault contract account to beneficiary.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("amount", true, false, "withdraw amount."),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		amount, ok := new(big.Int).SetString(req.Arguments[0], 10)
		if !ok {
			return fmt.Errorf("amount:%s cannot be parsed", req.Arguments[0])
		}
		hash, err := chain.SettleObject.VaultService.Withdraw(context.Background(), amount)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &VaultWithdrawCmdRet{
			Hash: hash.String(),
		})
	},
	Type: &VaultWithdrawCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *VaultWithdrawCmdRet) error {
			_, err := fmt.Fprintf(w, "the hash of transaction: %s", out.Hash)
			return err
		}),
	},
}
