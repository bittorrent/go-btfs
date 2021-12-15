package cheque

import (
	"context"
	"fmt"
	"io"
	"math/big"

	cmds "github.com/TRON-US/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

var ReceiveChequeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List cheque(s) received from peers.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("peer-id", true, false, "deposit amount."),
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {

		var record cheque
		peer_id := req.Arguments[0]
		fmt.Println("ReceiveChequeCmd peer_id = ", peer_id)

		if len(peer_id) > 0 {
			chequeTmp, err := chain.SettleObject.SwapService.LastReceivedCheque(peer_id)
			if err != nil {
				return err
			}

			record.Beneficiary = chequeTmp.Beneficiary.String()
			record.Vault = chequeTmp.Vault.String()
			record.Payout = chequeTmp.CumulativePayout
			record.PeerID = peer_id

			cashStatus, err := chain.SettleObject.CashoutService.CashoutStatus(context.Background(), chequeTmp.Vault)
			if err != nil {
				return err
			}
			if cashStatus.UncashedAmount != nil {
				record.CashedAmount = big.NewInt(0).Sub(chequeTmp.CumulativePayout, cashStatus.UncashedAmount)
			}
		}

		return cmds.EmitOnce(res, &record)
	},
	Type: cheque{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *cheque) error {
			//fmt.Fprintln(w, "cheque status:")
			fmt.Fprintf(w, "\t%-55s\t%-46s\t%-46s\t%-46s\tamount: \n", "peerID:", "vault:", "beneficiary:", "cashout_amount:")
			fmt.Fprintf(w, "\t%-55s\t%-46s\t%-46s\t%d\t%d \n",
				out.PeerID,
				out.Beneficiary,
				out.Vault,
				out.Payout.Uint64(),
				out.CashedAmount.Uint64(),
			)

			return nil
		}),
	},
}
