package cheque

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"io"
	"math/big"
	"strconv"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"go4.org/sort"
	"golang.org/x/net/context"
)

var ListReceiveChequeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List cheque(s) received from peers.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("offset", true, false, "page offset"),
		cmds.StringArg("limit", true, false, "page limit."),
	},
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
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
		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		//fmt.Printf("... token:%+v\n", tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		var listRet ListChequeRet
		cheques, err := chain.SettleObject.SwapService.LastReceivedCheques(token)

		if err != nil {
			return err
		}
		peerIds := make([]string, 0, len(cheques))
		for key := range cheques {
			peerIds = append(peerIds, key)
		}
		sort.Strings(peerIds)
		//[offset:offset+limit]
		if len(peerIds) < offset+1 {
			peerIds = peerIds[0:0]
		} else {
			peerIds = peerIds[offset:]
		}

		if len(peerIds) > limit {
			peerIds = peerIds[:limit]
		}

		for _, k := range peerIds {
			v := cheques[k]
			var record cheque
			record.PeerID = k
			record.Token = v.Token.String()
			record.Beneficiary = v.Beneficiary.String()
			record.Vault = v.Vault.String()
			record.Payout = v.CumulativePayout

			cashStatus, err := chain.SettleObject.CashoutService.CashoutStatus(context.Background(), v.Vault, token)
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
