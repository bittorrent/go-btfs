package cheque

import (
	"math/big"

	cmds "github.com/TRON-US/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

type chequeSentHistoryStats struct {
	TotalIssued      big.Int `json:"total_issued"`
	TotalIssuedCount int     `json:"total_issued_count"`
	Date             int64   `json:"date"` //time.now().Unix()
}

var ChequeSendHistoryStatsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Display the received cheques from peer.",
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		stats, err := chain.SettleObject.ChequeStore.SentStatsHistory(30)
		if err != nil {
			return err
		}

		ret := make([]chequeSentHistoryStats, 0, len(stats))
		for _, stat := range stats {
			ret = append(ret, chequeSentHistoryStats{
				TotalIssued:      *stat.Amount,
				TotalIssuedCount: stat.Count,
				Date:             stat.Date,
			})
		}
		return cmds.EmitOnce(res, &ret)
	},
	Type: []chequeSentHistoryStats{},
	//Encoders: cmds.EncoderMap{
	//	cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *[]chequeSendHistoryStats) error {
	//		return nil
	//	}),
	//},
}
