package cheque

import (
	"math/big"

	cmds "github.com/TRON-US/go-btfs-cmds"
)

type chequeSendHistoryStats struct {
	TotalIssued      big.Int `json:"total_issued"`
	TotalIssuedCount int     `json:"total_issued_count"`
	Date             int64   `json:"date"` //time.now().Unix()
}

var ChequeSendHistoryStatsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Display the received cheques from peer.",
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		return cmds.EmitOnce(res, &[]chequeSendHistoryStats{
			{},
		})
	},
	Type: []chequeSendHistoryStats{},
	//Encoders: cmds.EncoderMap{
	//	cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *[]chequeSendHistoryStats) error {
	//		return nil
	//	}),
	//},
}
