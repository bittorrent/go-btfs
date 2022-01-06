package stake

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"time"

	cmds "github.com/TRON-US/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

type AddStakeCmdRet struct {
	Hash string `json:"hash"`
}

var AddStakeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "add stake.",
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("amount", true, false, "stake amount."),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		amount, ok := new(big.Int).SetString(req.Arguments[0], 10)
		if !ok {
			return fmt.Errorf("amount:%s cannot be parsed", req.Arguments[0])
		}
		hash, err := chain.SettleObject.StakeService.AddStake(context.Background(), amount)

		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &AddStakeCmdRet{
			Hash: hash.String(),
		})
	},
	Type: &AddStakeCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *AddStakeCmdRet) error {
			_, err := fmt.Fprintf(w, "the tx is: %s", out.Hash)
			return err
		}),
	},
}
