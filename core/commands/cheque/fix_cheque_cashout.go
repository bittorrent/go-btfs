package cheque

import (
	"fmt"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"golang.org/x/net/context"
	"io"
)

var FixChequeCashOutCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List cheque(s) received from peers.",
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		fmt.Println("FixChequeCashOutCmd ... ")

		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		for _, tokenAddr := range tokencfg.MpTokenAddr {
			fmt.Println("FixChequeCashOutCmd ... 2")
			cheques, err := chain.SettleObject.SwapService.LastReceivedCheques(tokenAddr)
			fmt.Println("FixChequeCashOutCmd ... 3", cheques)
			if err != nil {
				return err
			}
			for _, v := range cheques {
				err := chain.SettleObject.CashoutService.AdjustCashCheque(
					context.Background(), v.Vault, v.Beneficiary, tokenAddr)
				if err != nil {
					return err
				}
			}
		}

		return cmds.EmitOnce(res, nil)
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
