package cheque

import (
	"fmt"
	"strconv"

	cmds "github.com/TRON-US/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

type chequeReceivedHistoryListRet struct {
	Records []chequeRecordRet `json:"records"`
	Total   int               `json:"total"`
}

var ChequeReceiveHistoryListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Display the received cheques from peer.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("from", true, false, "page offset"),
		cmds.StringArg("limit", true, false, "page limit."),
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		from, err := strconv.Atoi(req.Arguments[0])
		if err != nil {
			return fmt.Errorf("parse from:%v failed", req.Arguments[0])
		}
		limit, err := strconv.Atoi(req.Arguments[1])
		if err != nil {
			return fmt.Errorf("parse limit:%v failed", req.Arguments[1])
		}

		var listRet chequeReceivedHistoryListRet
		records, err := chain.SettleObject.SwapService.ReceivedChequeRecordsAll()
		if err != nil {
			return err
		}
		listRet.Total = len(records)
		ret := make([]chequeRecordRet, 0, limit)
		if from < len(records) {
			if (from + limit) <= len(records) {
				records = records[from : from+limit]
			} else {
				records = records[from:]
			}
			for _, result := range records {
				peer, known, err := chain.SettleObject.SwapService.VaultPeer(result.Vault)
				if err == nil {
					if !known {
						continue
					}
					r := chequeRecordRet{
						PeerId:      peer,
						Vault:       result.Vault,
						Beneficiary: result.Beneficiary,
						Amount:      result.Amount,
						ReceiveTime: result.ReceiveTime,
					}
					ret = append(ret, r)
				}
			}
		}
		listRet.Records = ret

		return cmds.EmitOnce(res, &listRet)
	},
	Type: chequeReceivedHistoryListRet{},
	//Encoders: cmds.EncoderMap{
	//	cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *ChequeRecords) error {
	//		var tm time.Time
	//		fmt.Fprintf(w, "\t%-46s\t%-46s\t%-10s\ttimestamp: \n", "beneficiary:", "vault:", "amount:")
	//		for index := 0; index < out.Len; index++ {
	//			tm = time.Unix(out.Records[index].ReceiveTime, 0)
	//			year, mon, day := tm.Date()
	//			h, m, s := tm.Clock()
	//			fmt.Fprintf(w, "\t%-46s\t%-46s\t%-10d\t%d-%d-%d %02d:%02d:%02d \n",
	//				out.Records[index].Beneficiary,
	//				out.Records[index].Vault,
	//				out.Records[index].Amount.Uint64(),
	//				year, mon, day, h, m, s)
	//		}
	//		return nil
	//	}),
	//},
}
