package vault

import (
	"fmt"
	"github.com/bittorrent/go-btfs/utils"
	"io"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/net/context"
)

type VaultWbttBalanceCmdRet struct {
	Balance *big.Int `json:"balance"`
}

var VaultWbttBalanceCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get wbtt balance by addr.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("addr", true, false, "wbtt token address"),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		addr := req.Arguments[0]
		balance, err := chain.SettleObject.VaultService.WBTTBalanceOf(context.Background(), common.HexToAddress(addr))
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &VaultBalanceCmdRet{
			Balance: balance,
		})
	},
	Type: &VaultWbttBalanceCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *VaultWbttBalanceCmdRet) error {
			_, err := fmt.Fprintf(w, "the balance is: %v\n", out.Balance)
			return err
		}),
	},
}
