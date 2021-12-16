package cheque

import (
	"fmt"
	"io"
	"time"

	cmds "github.com/TRON-US/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

var ChequeSendHistoryListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Display the send cheques from peer.",
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		var listRet ChequeRecords

		records, err := chain.SettleObject.SwapService.SendChequeRecordsAll()
		if err != nil {
			return err
		}
		listRet.Records = records
		listRet.Len = len(records)

		return cmds.EmitOnce(res, &listRet)
	},
	Type: ChequeRecords{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *ChequeRecords) error {
			var tm time.Time
			fmt.Fprintf(w, "\t%-46s\t%-46s\t%-10s\ttimestamp: \n", "beneficiary:", "vault:", "amount:")
			for index := 0; index < out.Len; index++ {
				tm = time.Unix(out.Records[index].ReceiveTime, 0)
				year, mon, day := tm.Date()
				h, m, s := tm.Clock()
				fmt.Fprintf(w, "\t%-46s\t%-46s\t%-10d\t%d-%d-%d %02d:%02d:%02d \n",
					out.Records[index].Beneficiary,
					out.Records[index].Vault,
					out.Records[index].Amount.Uint64(),
					year, mon, day, h, m, s)
			}
			return nil
		}),
	},
}
