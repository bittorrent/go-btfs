package cheque

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"io"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

type ReceiveTotalCountRet struct {
	Count int `json:"count"`
}

var ReceiveChequesCountCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "send cheque(s) count",
	},
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		fmt.Printf("... token:%+v\n", tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		count, err := chain.SettleObject.SwapService.ReceivedChequeRecordsCount(token)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &ReceiveTotalCountRet{Count: count})
	},
	Type: ReceiveTotalCountRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, c *ReceiveTotalCountRet) error {
			fmt.Println("receive cheque(s) count: ", c.Count)

			return nil
		}),
	},
}
