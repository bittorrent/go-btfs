package cheque

import (
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"io"
	"math/big"
	"sort"
	"strconv"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"golang.org/x/net/context"
)

var ListReceiveChequeAllCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List cheque(s) received from peers.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("offset", true, false, "page offset"),
		cmds.StringArg("limit", true, false, "page limit."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		offset, err := strconv.Atoi(req.Arguments[0])
		if err != nil {
			return fmt.Errorf("parse offset:%v failed", req.Arguments[0])
		}
		limit, err := strconv.Atoi(req.Arguments[1])
		if err != nil {
			return fmt.Errorf("parse limit:%v failed", req.Arguments[1])
		}

		listCheques := make([]ReceiveCheque, 0)
		for _, tokenAddr := range tokencfg.MpTokenAddr {
			cheques, err := chain.SettleObject.SwapService.LastReceivedCheques(tokenAddr)
			if err != nil {
				return err
			}
			for k, v := range cheques {
				var record ReceiveCheque
				record.PeerID = k
				record.Token = tokenAddr
				record.Vault = v.Vault
				record.Beneficiary = v.Beneficiary
				record.CumulativePayout = v.CumulativePayout

				listCheques = append(listCheques, record)
			}
		}

		sort.Slice(listCheques, func(i, j int) bool {
			return listCheques[i].PeerID < listCheques[j].PeerID
		})

		//[offset:offset+limit]
		if len(listCheques) < offset+1 {
			listCheques = listCheques[0:0]
		} else {
			listCheques = listCheques[offset:]
		}

		if len(listCheques) > limit {
			listCheques = listCheques[:limit]
		}

		var listRet ListChequeRet
		for _, v := range listCheques {
			k := v.PeerID
			var record cheque
			record.PeerID = k
			record.Token = v.Token.String()
			record.Beneficiary = v.Beneficiary.String()
			record.Vault = v.Vault.String()
			record.Payout = v.CumulativePayout

			cashStatus, err := chain.SettleObject.CashoutService.CashoutStatus(context.Background(), v.Vault, v.Token)
			if err != nil {
				return err
			}
			if cashStatus.UncashedAmount != nil {
				record.CashedAmount = big.NewInt(0).Sub(v.CumulativePayout, cashStatus.UncashedAmount)
			}

			listRet.Cheques = append(listRet.Cheques, record)
		}

		listRet.Len = len(listRet.Cheques)
		return cmds.EmitOnce(res, &listRet)
	},
	Type: ListChequeRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *ListChequeRet) error {
			fmt.Fprintf(w, "\t%-55s\t%-46s\t%-46s\t%-46s\tamount: \n", "peerID:", "vault:", "beneficiary:", "cashout_amount:")
			for iter := 0; iter < out.Len; iter++ {
				fmt.Fprintf(w, "\t%-55s\t%-46s\t%-46s\t%d\t%d \n",
					out.Cheques[iter].PeerID,
					out.Cheques[iter].Beneficiary,
					out.Cheques[iter].Vault,
					out.Cheques[iter].Payout.Uint64(),
					out.Cheques[iter].CashedAmount.Uint64(),
				)
			}

			return nil
		}),
	},
}
