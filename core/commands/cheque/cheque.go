package cheque

import (
	"fmt"
	"io"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/ethereum/go-ethereum/common"
)

type StorePriceRet struct {
	Price *big.Int `json:"price"`
}

type CashChequeRet struct {
	TxHash string
}

type cheque struct {
	PeerID       string
	Beneficiary  string
	Vault        string
	Payout       *big.Int
	CashedAmount *big.Int
}

type ListChequeRet struct {
	Cheques []cheque
	Len     int
}

type ChequeRecords struct {
	Records []chequeRecordRet
	Len     int
}

type chequeRecordRet struct {
	PeerId      string
	Vault       common.Address
	Beneficiary common.Address
	Amount      *big.Int
	Time        int64 //time.now().Unix()
}

var ChequeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Interact with vault services on BTFS.",
		ShortDescription: `
Vault services include issue cheque to peer, receive cheque and store operations.`,
	},
	Subcommands: map[string]*cmds.Command{
		"cash":       CashChequeCmd,
		"cashstatus": ChequeCashStatusCmd,
		"cashlist":   ChequeCashListCmd,
		"price":      StorePriceCmd,

		"send":               SendChequeCmd,
		"sendlist":           ListSendChequesCmd,
		"send-history-peer":  ChequeSendHistoryPeerCmd,
		"send-history-list":  ChequeSendHistoryListCmd,
		"send-history-stats": ChequeSendHistoryStatsCmd,
		"send-total-count":   SendChequesCountCmd,

		"receive":               ReceiveChequeCmd,
		"receivelist":           ListReceiveChequeCmd,
		"receive-history-peer":  ChequeReceiveHistoryPeerCmd,
		"receive-history-list":  ChequeReceiveHistoryListCmd,
		"receive-history-stats": ChequeReceiveHistoryStatsCmd,
		"receive-total-count":   ReceiveChequesCountCmd,
		"stats":                 ChequeStatsCmd,

		"chaininfo":  ChequeChainInfoCmd,
		"bttbalance": ChequeBttBalanceCmd,
	},
}

var StorePriceCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get btfs store price.",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		totalPrice, err := chain.SettleObject.OracleService.CheckNewPrice()
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &StorePriceRet{
			Price: totalPrice,
		})
	},
	Type: StorePriceRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *StorePriceRet) error {
			_, err := fmt.Fprintf(w, "the btfs store price: %v\n", out.Price)
			return err
		}),
	},
}

var CashChequeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Cash a cheque by peerID.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("peer-id", true, false, "Peer id tobe cashed."),
		cmds.StringArg("token", true, false, "token"),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {

		// get the peer id
		peerID := req.Arguments[0]
		token := req.Arguments[1]
		fmt.Printf("... token:%+v\n", token)
		tx_hash, err := chain.SettleObject.SwapService.CashCheque(req.Context, peerID, token)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &CashChequeRet{
			TxHash: tx_hash.String(),
		})
	},
	Type: CashChequeRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *CashChequeRet) error {
			_, err := fmt.Fprintf(w, "the hash of transaction: %s", out.TxHash)
			return err
		}),
	},
}
