package cheque

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	cmds "github.com/TRON-US/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

type chequeReceivedHistoryStats struct {
	TotalReceived      *big.Int `json:"total_received"`
	TotalReceivedCount int      `json:"total_received_count"`
	Date               int64    `json:"date"` //time.now().Unix()
}

var ChequeReceiveHistoryStatsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Display the received cheques from peer.",
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		stats, err := chain.SettleObject.ChequeStore.ReceivedStatsHistory(30)
		if err != nil {
			return err
		}

		ret := make([]chequeReceivedHistoryStats, 0, len(stats))
		for _, stat := range stats {
			ret = append(ret, chequeReceivedHistoryStats{
				TotalReceived:      stat.Amount,
				TotalReceivedCount: stat.Count,
				Date:               stat.Date,
			})
		}
		return cmds.EmitOnce(res, &ret)
	},
	Type: []chequeReceivedHistoryStats{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *[]chequeReceivedHistoryStats) error {
			marshaled, err := json.MarshalIndent(out, "", "\t")
			if err != nil {
				return err
			}
			marshaled = append(marshaled, byte('\n'))
			fmt.Fprintln(w, string(marshaled))
			return nil
		}),
	},
}
