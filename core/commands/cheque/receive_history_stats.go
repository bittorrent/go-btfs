package cheque

import (
	"math/big"

	cmds "github.com/TRON-US/go-btfs-cmds"
)

type chequeReceivedHistoryStats struct {
	TotalReceived      big.Int `json:"total_received"`
	TotalReceivedCount int     `json:"total_received_count"`
	Date               int64   `json:"date"` //time.now().Unix()
}

var ChequeReceiveHistoryStatsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Display the received cheques from peer.",
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		return cmds.EmitOnce(res, &[]chequeReceivedHistoryStats{
			{},
		})
	},
	Type: []chequeReceivedHistoryStats{},
	//Encoders: cmds.EncoderMap{
	//	cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *[]chequeReceivedHistoryStats) error {
	//		return nil
	//	}),
	//},
}
