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

type ApproveCmdRet struct {
	Hash string `json:"hash"`
}

var ApproveCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "approve to stake contract.",
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("amount", true, false, "approve amount."),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		amount, ok := new(big.Int).SetString(req.Arguments[0], 10)
		if !ok {
			return fmt.Errorf("amount:%s cannot be parsed", req.Arguments[0])
		}
		hash, err := chain.SettleObject.StakeService.Approve(context.Background(), amount)

		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &ApproveCmdRet{
			Hash: hash.String(),
		})
	},
	Type: &ApproveCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *ApproveCmdRet) error {
			_, err := fmt.Fprintf(w, "the tx is: %s", out.Hash)
			return err
		}),
	},
}
