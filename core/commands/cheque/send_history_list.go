package cheque

import (
	"encoding/json"
	"fmt"
	"github.com/bittorrent/go-btfs/utils"
	"io"
	"sort"
	"strconv"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

var ChequeSendHistoryListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Display the send cheques from peer.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("from", true, false, "page offset"),
		cmds.StringArg("limit", true, false, "page limit."),
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		from, err := strconv.Atoi(req.Arguments[0])
		if err != nil {
			return fmt.Errorf("parse from:%v failed", req.Arguments[0])
		}
		limit, err := strconv.Atoi(req.Arguments[1])
		if err != nil {
			return fmt.Errorf("parse limit:%v failed", req.Arguments[1])
		}
		if from < 0 {
			return fmt.Errorf("invalid from: %d", from)
		}
		if limit < 0 {
			return fmt.Errorf("invalid limit: %d", limit)
		}

		var listRet chequeReceivedHistoryListRet
		records, err := chain.SettleObject.SwapService.SendChequeRecordsAll()
		if err != nil {
			return err
		}
		sort.Slice(records, func(i, j int) bool {
			return records[i].ReceiveTime > records[j].ReceiveTime
		})
		listRet.Total = len(records)
		ret := make([]chequeRecordRet, 0, limit)
		if from < len(records) {
			if (from + limit) <= len(records) {
				records = records[from : from+limit]
			} else {
				records = records[from:]
			}
			for _, result := range records {
				peer, known, err := chain.SettleObject.SwapService.BeneficiaryPeer(result.Beneficiary)
				if err == nil {
					if !known {
						peer = "unknown"
					}
					r := chequeRecordRet{
						PeerId:      peer,
						Token:       result.Token,
						Vault:       result.Vault,
						Beneficiary: result.Beneficiary,
						Amount:      result.Amount,
						Time:        result.ReceiveTime,
					}
					ret = append(ret, r)
				}
			}
		}
		listRet.Records = ret

		return cmds.EmitOnce(res, &listRet)
	},
	Type: chequeReceivedHistoryListRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *chequeReceivedHistoryListRet) error {
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
