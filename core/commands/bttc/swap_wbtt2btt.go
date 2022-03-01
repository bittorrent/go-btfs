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

type BttcWbtt2BttCmdRet struct {
	Hash string `json:"hash"`
}

var BttcWbtt2BttCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Swap WBTT to BTT at you bttc address",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("amount", true, false, "amount you want to swap"),
	},
	RunTimeout: 5 * time.Minute,
	Type:       &BttcWbtt2BttCmdRet{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) (err error) {
		amountStr := utils.RemoveSpaceAndComma(req.Arguments[0])
		amount, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			return fmt.Errorf("invalid argument amount %s", req.Arguments[0])
		}
		trx, err := chain.SettleObject.BttcService.SwapWbtt2Btt(context.Background(), amount)
		if err != nil {
			return
		}
		return cmds.EmitOnce(res, &BttcWbtt2BttCmdRet{Hash: trx.String()})
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *BttcWbtt2BttCmdRet) error {
			_, err := fmt.Fprintf(w, "the hash of transaction: %s\n", out.Hash)
			return err
		}),
	},
}
