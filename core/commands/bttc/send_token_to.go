package bttc

import (
	"context"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"io"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/utils"
	"github.com/ethereum/go-ethereum/common"
)

type BttcSendTokenToCmdRet struct {
	Hash string `json:"hash"`
}

var BttcSendTokenToCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Transfer your WBTT to other bttc address",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("addr", true, false, "target bttc address"),
		cmds.StringArg("amount", true, false, "amount you want to send"),
	},
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	RunTimeout: 5 * time.Minute,
	Type:       &BttcSendTokenToCmdRet{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) (err error) {
		err = utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		addressStr := req.Arguments[0]
		if !common.IsHexAddress(addressStr) {
			return fmt.Errorf("invalid bttc address %s", addressStr)
		}

		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		//fmt.Printf("token:%+v\n", tokenStr)
		_, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		amount, ok := new(big.Int).SetString(utils.RemoveSpaceAndComma(req.Arguments[1]), 10)
		if !ok {
			return fmt.Errorf("invalid argument amount %s", req.Arguments[1])
		}

		trx, err := chain.SettleObject.BttcService.SendTokenTo(context.Background(), common.HexToAddress(addressStr), amount, tokenStr)
		if err != nil {
			return
		}

		return cmds.EmitOnce(res, &BttcSendTokenToCmdRet{Hash: trx.String()})
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *BttcSendTokenToCmdRet) error {
			_, err := fmt.Fprintf(w, "the hash of transaction: %s\n", out.Hash)
			return err
		}),
	},
}
