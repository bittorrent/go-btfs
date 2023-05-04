package cheque

import (
	"context"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"io"
	"math/big"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

var ReceiveChequeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List cheque(s) received from peers.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("peer-id", true, false, "deposit amount."),
	},
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		var record cheque
		peer_id := req.Arguments[0]
		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		fmt.Printf("ReceiveChequeCmd peer_id:%+v, token:%+v\n", peer_id, tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		if len(peer_id) > 0 {
			chequeTmp, err := chain.SettleObject.SwapService.LastReceivedCheque(peer_id, token)
			if err != nil {
				return err
			}

			record.Beneficiary = chequeTmp.Beneficiary.String()
			record.Vault = chequeTmp.Vault.String()
			record.Payout = chequeTmp.CumulativePayout
			record.PeerID = peer_id

			cashStatus, err := chain.SettleObject.CashoutService.CashoutStatus(context.Background(), chequeTmp.Vault, token)
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
