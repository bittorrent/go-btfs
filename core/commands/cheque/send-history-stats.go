package cheque

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

type chequeSentHistoryStats struct {
	TotalIssued      *big.Int `json:"total_issued"`
	TotalIssuedCount int      `json:"total_issued_count"`
	Date             int64    `json:"date"` //time.now().Unix()
}

var ChequeSendHistoryStatsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Display the received cheques from peer.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("token", true, false, "token"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		token := req.Arguments[0]
		fmt.Printf("... token:%+v\n", token)

		// now only return 30days cheque sent stats
		const sentStatsDuration = 30
		stats, err := chain.SettleObject.ChequeStore.SentStatsHistory(sentStatsDuration, token)
		if err != nil {
			return err
		}

		ret := make([]chequeSentHistoryStats, 0, len(stats))
		for _, stat := range stats {
			ret = append(ret, chequeSentHistoryStats{
				TotalIssued:      stat.Amount,
				TotalIssuedCount: stat.Count,
				Date:             stat.Date,
			})
		}
		return cmds.EmitOnce(res, &ret)
	},
	Type: []chequeSentHistoryStats{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *[]chequeSentHistoryStats) error {
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
