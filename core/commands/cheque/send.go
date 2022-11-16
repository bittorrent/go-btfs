package cheque

import (
	"fmt"
	"io"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

var SendChequeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List cheque send to peers.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("peer-id", true, false, "deposit amount."),
		cmds.StringArg("token", true, false, "token"),
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		var record cheque
		peer_id := req.Arguments[0]
		fmt.Println("SendChequeCmd peer_id = ", peer_id)

		token := req.Arguments[1]
		fmt.Printf("... token:%+v\n", token)

		if len(peer_id) > 0 {
			chequeTmp, err := chain.SettleObject.SwapService.LastSendCheque(peer_id, token)
			if err != nil {
				return err
			}

			record.Beneficiary = chequeTmp.Beneficiary.String()
			record.Vault = chequeTmp.Vault.String()
			record.Payout = chequeTmp.CumulativePayout
			record.PeerID = peer_id
		}

		return cmds.EmitOnce(res, &record)
	},
	Type: cheque{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *cheque) error {
			//fmt.Fprintln(w, "cheque status:")
			fmt.Fprintf(w, "\t%-55s\t%-46s\t%-46s\tamount: \n", "peerID:", "vault:", "beneficiary:")
			fmt.Fprintf(w, "\t%-55s\t%-46s\t%-46s\t%d \n",
				out.PeerID,
				out.Vault,
				out.Beneficiary,
				out.Payout.Uint64(),
			)

			return nil
		}),
	},
}
