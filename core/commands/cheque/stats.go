package cheque

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"golang.org/x/net/context"
)

type chequeStats struct {
	TotalIssuedCount  int      `json:"total_issued_count"`
	TotalIssued       *big.Int `json:"total_issued"`
	TotalIssuedCashed *big.Int `json:"total_issued_cashed"`

	TotalReceived              *big.Int `json:"total_received"`
	TotalReceivedUncashed      *big.Int `json:"total_received_uncashed"`
	TotalReceivedCount         int      `json:"total_received_count"`
	TotalReceivedCashedCount   int      `json:"total_received_cashed_count"`
	TotalReceivedDailyUncashed *big.Int `json:"total_received_daily_uncashed"`
}

var ChequeStatsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List cheque(s) received from peers.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("token", true, false, "token"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		cs := chequeStats{
			TotalIssued:                big.NewInt(0),
			TotalIssuedCashed:          big.NewInt(0),
			TotalReceived:              big.NewInt(0),
			TotalReceivedUncashed:      big.NewInt(0),
			TotalReceivedDailyUncashed: big.NewInt(0),
		}

		token := req.Arguments[0]
		fmt.Printf("... token:%+v\n", token)

		issued, err := chain.SettleObject.VaultService.TotalIssued(token)
		if err != nil {
			return err
		}
		cs.TotalIssued = issued

		issuedCount, err := chain.SettleObject.VaultService.TotalIssuedCount(token)
		if err != nil {
			return err
		}
		cs.TotalIssuedCount = issuedCount

		paidOut, err := chain.SettleObject.VaultService.TotalPaidOut(context.Background(), token)
		if err != nil {
			return err
		}
		cs.TotalIssuedCashed = paidOut

		received, err := chain.SettleObject.VaultService.TotalReceived(token)
		if err != nil {
			return err
		}
		cs.TotalReceived = received

		cashed, err := chain.SettleObject.VaultService.TotalReceivedCashed(token)
		if err != nil {
			return err
		}
		cs.TotalReceivedUncashed.Sub(cs.TotalReceived, cashed)

		count, err := chain.SettleObject.VaultService.TotalReceivedCount(token)
		if err != nil {
			return err
		}
		cs.TotalReceivedCount = count

		receivedCount, err := chain.SettleObject.VaultService.TotalReceivedCashedCount(token)
		if err != nil {
			return err
		}
		cs.TotalReceivedCashedCount = receivedCount

		dailyReceived, err := chain.SettleObject.VaultService.TotalDailyReceived(token)
		if err != nil {
			return err
		}
		cs.TotalReceivedDailyUncashed = dailyReceived

		return cmds.EmitOnce(res, &cs)
	},
	Type: &chequeStats{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *chequeStats) error {
			marshaled, err := json.MarshalIndent(out, "", "\t")
			if err != nil {
				return err
			}
			marshaled = append(marshaled, byte('\n'))
			fmt.Fprintln(w, string(marshaled))
			return nil
		}),
	},
}
