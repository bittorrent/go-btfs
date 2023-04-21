package cheque

import (
	"encoding/json"
	"fmt"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"io"
	"math/big"
)

var ChequeStatsAllCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List cheque(s) received from peers, of all tokens",
	},
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		mp := make(map[string]*chequeStats, 0)
		for k, tokenAddr := range tokencfg.MpTokenAddr {
			cs := chequeStats{
				TotalIssued:                big.NewInt(0),
				TotalIssuedCashed:          big.NewInt(0),
				TotalReceived:              big.NewInt(0),
				TotalReceivedUncashed:      big.NewInt(0),
				TotalReceivedDailyUncashed: big.NewInt(0),
			}

			err := GetChequeStatsToken(&cs, tokenAddr)
			if err != nil {
				return err
			}

			mp[k] = &cs
		}

		return cmds.EmitOnce(res, &mp)
	},
	Type: &chequeStats{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *chequeStats) error {
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
