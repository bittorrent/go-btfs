package cheque

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"io"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

type SendTotalCountRet struct {
	Count int `json:"count"`
}

var SendChequesCountCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "send cheque(s) count",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("token", true, false, "token"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		tokenStr := req.Arguments[0]
		fmt.Printf("... token:%+v\n", tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		count, err := chain.SettleObject.VaultService.TotalIssuedCount(token)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &SendTotalCountRet{Count: count})
	},
	Type: SendTotalCountRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, c *SendTotalCountRet) error {
			fmt.Println("send cheque(s) count: ", c.Count)

			return nil
		}),
	},
}
