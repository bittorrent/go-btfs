package cheque

import (
	"encoding/json"
	"fmt"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"io"
)

var ChequeSendHistoryStatsAllCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Display the received cheques from peer, of all tokens",
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		mp := make(map[string][]chequeSentHistoryStats, 0)
		for k, tokenAddr := range tokencfg.MpTokenAddr {
			// now only return 30days cheque sent stats
			const sentStatsDuration = 30
			stats, err := chain.SettleObject.ChequeStore.SentStatsHistory(sentStatsDuration, tokenAddr)
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

			mp[k] = ret
		}

		return cmds.EmitOnce(res, &mp)
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
