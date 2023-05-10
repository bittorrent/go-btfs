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

var ChequeReceiveHistoryStatsAllCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Display the received cheques from peer, of all tokens.",
	},
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		mp := make(map[string][]chequeReceivedHistoryStats, 0)
		for k, tokenAddr := range tokencfg.MpTokenAddr {

			// now only return 30days cheque received stats
			const receivedStatsDuration = 30
			stats, err := chain.SettleObject.ChequeStore.ReceivedStatsHistory(receivedStatsDuration, tokenAddr)
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

			mp[k] = ret
		}

		return cmds.EmitOnce(res, &mp)
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
