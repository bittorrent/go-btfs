package cheque

import (
	"fmt"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"github.com/google/martian/log"
	"golang.org/x/net/context"
	"io"
)

var FixChequeCashOutCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List cheque(s) received from peers.",
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		listRet := ListFixChequeRet{}
		listRet.FixCheques = make([]fixCheque, 0)

		for _, tokenAddr := range tokencfg.MpTokenAddr {
			cheques, err := chain.SettleObject.SwapService.LastReceivedCheques(tokenAddr)
			if err != nil {
				return err
			}

			for k, v := range cheques {
				totalCashOutAmount, newCashOutAmount, err := chain.SettleObject.CashoutService.AdjustCashCheque(
					context.Background(), v.Vault, v.Beneficiary, tokenAddr)
				if err != nil {
					return err
				}
				if newCashOutAmount != nil && newCashOutAmount.Uint64() > 0 {
					var record fixCheque
					record.PeerID = k
					record.Token = v.Token.String()
					record.Beneficiary = v.Beneficiary.String()
					record.Vault = v.Vault.String()
					record.TotalCashedAmount = totalCashOutAmount
					record.FixCashedAmount = newCashOutAmount

					listRet.FixCheques = append(listRet.FixCheques, record)
				}
			}
		}
		listRet.Len = len(listRet.FixCheques)

		log.Infof("FixChequeCashOutCmd, listRet = %+v", listRet)

		return cmds.EmitOnce(res, &listRet)
	},
	Type: ListFixChequeRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *ListFixChequeRet) error {
			fmt.Fprintf(w, "fix: \n\t%-55s\t%-46s\t%-46s\t%-46s\tfix_cash_amount: \n", "peerID:", "vault:", "beneficiary:", "total_cash_amount:")
			for iter := 0; iter < out.Len; iter++ {
				fmt.Fprintf(w, "\t%-55s\t%-46s\t%-46s\t%d\t%d \n",
					out.FixCheques[iter].PeerID,
					out.FixCheques[iter].Vault,
					out.FixCheques[iter].Beneficiary,
					out.FixCheques[iter].TotalCashedAmount.Uint64(),
					out.FixCheques[iter].FixCashedAmount.Uint64(),
				)
			}

			return nil
		}),
	},
}
