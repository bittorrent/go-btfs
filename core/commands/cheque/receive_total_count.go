package cheque

import (
	"fmt"
	"github.com/bittorrent/go-btfs/utils"
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

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		count, err := chain.SettleObject.SwapService.ReceivedChequeRecordsCount()
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
