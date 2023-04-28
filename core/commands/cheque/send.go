package cheque

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
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
		fmt.Println("SendChequeCmd peer_id = ", peer_id)

		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		//fmt.Printf("... token:%+v\n", tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

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
