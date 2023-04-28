package cheque

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

type ChequeBttBalanceCmdRet struct {
	Balance *big.Int `json:"balance"`
}

var ChequeBttBalanceCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get btt balance by addr.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("addr", true, false, "bttc account address"),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		addr := req.Arguments[0]
		balance, err := chain.SettleObject.VaultService.BTTBalanceOf(context.Background(), common.HexToAddress(addr), nil)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &ChequeBttBalanceCmdRet{
			Balance: balance,
		})
	},
	Type: &ChequeBttBalanceCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *ChequeBttBalanceCmdRet) error {
			_, err := fmt.Fprintf(w, "the balance is: %v\n", out.Balance)
			return err
		}),
	},
}
