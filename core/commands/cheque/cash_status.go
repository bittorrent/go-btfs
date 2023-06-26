package cheque

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

const (
	CashoutStatusSuccess  = "success"
	CashoutStatusFail     = "fail"
	CashoutStatusPending  = "pending"
	CashoutStatusNotFound = "not_found"
)

type CashOutStatusRet struct {
	Status         string   `json:"status"` // pending,fail,success,not_found
	TotalPayout    *big.Int `json:"total_payout"`
	UncashedAmount *big.Int `json:"uncashed_amount"` // amount not yet cashed out
}

var ChequeCashStatusCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get cash status by peerID.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("peer-id", true, false, "Peer id tobe cashed."),
	},
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		// get the peer id
		peerID := req.Arguments[0]
		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		fmt.Printf("peerID:%+v, token:%+v\n", peerID, tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		cashStatus, err := chain.SettleObject.SwapService.CashoutStatus(req.Context, peerID, token)
		if err != nil {
			return err
		}

		status := CashoutStatusSuccess
		totalPayout := big.NewInt(0)
		if cashStatus.Last == nil {
			status = CashoutStatusNotFound
		} else if cashStatus.Last.Reverted {
			status = CashoutStatusFail
		} else if cashStatus.Last.Result == nil {
			status = CashoutStatusPending
		} else {
			totalPayout = cashStatus.Last.Result.TotalPayout
		}

		return cmds.EmitOnce(res, &CashOutStatusRet{
			UncashedAmount: cashStatus.UncashedAmount,
			Status:         status,
			TotalPayout:    totalPayout,
		})
	},
	Type: &CashOutStatusRet{},
}
