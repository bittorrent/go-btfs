package cheque

import (
	"fmt"
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

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		count, err := chain.SettleObject.VaultService.TotalIssuedCount()

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
