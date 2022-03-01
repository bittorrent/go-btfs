package bttc

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/utils"
)

type BttcBtt2WbttCmdRet struct {
	Hash string `json:"hash"`
}

var BttcBtt2WbttCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Swap BTT to WBTT at you bttc address",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("amount", true, false, "amount you want to swap"),
	},
	RunTimeout: 5 * time.Minute,
	Type:       &BttcBtt2WbttCmdRet{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) (err error) {
		amountStr := utils.RemoveSpaceAndComma(req.Arguments[0])
		amount, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			return fmt.Errorf("invalid argument amount %s", req.Arguments[0])
		}
		trx, err := chain.SettleObject.BttcService.SwapBtt2Wbtt(context.Background(), amount)
		if err != nil {
			return
		}
		return cmds.EmitOnce(res, &BttcBtt2WbttCmdRet{Hash: trx.String()})
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *BttcBtt2WbttCmdRet) error {
			_, err := fmt.Fprintf(w, "the hash of transaction: %s\n", out.Hash)
			return err
		}),
	},
}
