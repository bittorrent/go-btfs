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

var ListSendChequesCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List cheque(s) send to peers.",
	},
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		//fmt.Printf("... token:%+v\n", tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		listRet := ListChequeRet{}
		listRet.Cheques = make([]cheque, 0, 0)
		cheques, err := chain.SettleObject.SwapService.LastSendCheques(token)

		if err != nil {
			return err
		}
		for k, v := range cheques {
			var record cheque
			record.PeerID = k
			record.Token = v.Token.String()
			record.Beneficiary = v.Beneficiary.String()
			record.Vault = v.Vault.String()
			record.Payout = v.CumulativePayout

			listRet.Cheques = append(listRet.Cheques, record)
		}

		listRet.Len = len(listRet.Cheques)

		return cmds.EmitOnce(res, &listRet)
	},
	Type: ListChequeRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *ListChequeRet) error {
			fmt.Fprintf(w, "\t%-55s\t%-46s\t%-46s\tamount: \n", "peerID:", "vault:", "beneficiary:")
			for iter := 0; iter < out.Len; iter++ {
				fmt.Fprintf(w, "\t%-55s\t%-46s\t%-46s\t%d \n",
					out.Cheques[iter].PeerID,
					out.Cheques[iter].Vault,
					out.Cheques[iter].Beneficiary,
					out.Cheques[iter].Payout.Uint64(),
				)
			}

			return nil
		}),
	},
}
