package cheque

import (
	"math/big"
	"time"

	cmds "github.com/TRON-US/go-btfs-cmds"
)

type ChequeCashListRet struct {
	TxHash   string  `json:"tx_hash"`
	PeerID   string  `json:"peer_id"`
	Vault    string  `json:"vault"`
	Amount   big.Int `json:"amount"`
	CashTime int64   `json:"cash_time"`
	Status   string  `json:"status"`
}

var ChequeCashListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get cash status by peerID.",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {

		return cmds.EmitOnce(res, &[]ChequeCashListRet{{Status: CashoutStatusSuccess}})
	},
	Type: &[]CashOutStatusRet{},
}
