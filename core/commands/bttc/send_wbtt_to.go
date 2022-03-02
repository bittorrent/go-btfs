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
	"github.com/ethereum/go-ethereum/common"
)

type BttcSendWbttToCmdRet struct {
	Hash string `json:"hash"`
}

var BttcSendWbttToCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Transfer your WBTT to other bttc address",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("addr", true, false, "target bttc address"),
		cmds.StringArg("amount", true, false, "amount you want to send"),
	},
	RunTimeout: 5 * time.Minute,
	Type:       &BttcSendWbttToCmdRet{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) (err error) {
		addressStr := req.Arguments[0]
		if !common.IsHexAddress(addressStr) {
			return fmt.Errorf("invalid bttc address %s", addressStr)
		}
		amount, ok := new(big.Int).SetString(utils.RemoveSpaceAndComma(req.Arguments[1]), 10)
		if !ok {
			return fmt.Errorf("invalid argument amount %s", req.Arguments[1])
		}
		trx, err := chain.SettleObject.BttcService.SendWbttTo(context.Background(), common.HexToAddress(addressStr), amount)
		if err != nil {
			return
		}
		return cmds.EmitOnce(res, &BttcSendWbttToCmdRet{Hash: trx.String()})
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *BttcSendWbttToCmdRet) error {
			_, err := fmt.Fprintf(w, "the hash of transaction: %s\n", out.Hash)
			return err
		}),
	},
}
