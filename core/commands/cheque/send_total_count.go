package cheque

import (
	"encoding/json"
	"fmt"
	"io"

	cmds "github.com/TRON-US/go-btfs-cmds"
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
			marshaled, err := json.MarshalIndent(c, "", "\t")
			if err != nil {
				return err
			}
			marshaled = append(marshaled, byte('\n'))
			fmt.Fprintln(w, string(marshaled))

			return nil
		}),
	},
}
