package cheque

import (
	"math/big"

	cmds "github.com/TRON-US/go-btfs-cmds"
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

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		cs := chequeStats{
			TotalIssued:                big.NewInt(0),
			TotalIssuedCashed:          big.NewInt(0),
			TotalReceived:              big.NewInt(0),
			TotalReceivedUncashed:      big.NewInt(0),
			TotalReceivedDailyUncashed: big.NewInt(0),
		}
		if issued, err := chain.SettleObject.VaultService.TotalIssued(); err == nil {
			cs.TotalIssued = issued
		}
		if issuedCount, err := chain.SettleObject.VaultService.TotalIssuedCount(); err == nil {
			cs.TotalIssuedCount = issuedCount
		}
		if paidOut, err := chain.SettleObject.VaultService.TotalPaidOut(context.Background()); err == nil {
			cs.TotalIssuedCashed = paidOut
		}

		if received, err := chain.SettleObject.VaultService.TotalReceived(); err == nil {
			cs.TotalReceived = received
		}

		if cashed, err := chain.SettleObject.VaultService.TotalReceivedCashed(); err == nil {
			if cs.TotalReceived == nil {
				cs.TotalReceived = big.NewInt(0)
			}
			cs.TotalReceivedUncashed.Sub(cs.TotalReceived, cashed)
		}

		if count, err := chain.SettleObject.VaultService.TotalReceivedCount(); err == nil {
			cs.TotalReceivedCount = count
		}
		if count, err := chain.SettleObject.VaultService.TotalReceivedCashedCount(); err == nil {
			cs.TotalReceivedCashedCount = count
		}
		if dailyReceived, err := chain.SettleObject.VaultService.TotalDailyReceived(); err == nil {
			cs.TotalReceivedDailyUncashed = dailyReceived
		}

		return cmds.EmitOnce(res, &cs)
	},
	Type: &chequeStats{},
	//Encoders: cmds.EncoderMap{
	//	cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *chequeStats) error {
	//		//fmt.Fprintln(w, "cheque status:")
	//		//fmt.Fprintf(w, "\t%-55s\t%-46s\t%-46s\t%-46s\tamount: \n", "peerID:", "vault:", "beneficiary:", "cashout_amount:")
	//		//fmt.Fprintf(w, "\t%-55s\t%-46s\t%-46s\t%d\t%d \n",
	//		//)
	//		return nil
	//	}),
	//},
}
