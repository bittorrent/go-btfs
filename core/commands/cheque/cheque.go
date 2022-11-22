package cheque

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
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
	Token        common.Address
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
	Token       common.Address
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

		"send":                   SendChequeCmd,
		"sendlist":               ListSendChequesCmd,
		"send-history-peer":      ChequeSendHistoryPeerCmd,
		"send-history-list":      ChequeSendHistoryListCmd,
		"send-history-stats":     ChequeSendHistoryStatsCmd,
		"send-history-stats-all": ChequeSendHistoryStatsAllCmd,
		"send-total-count":       SendChequesCountCmd,

		"receive":                   ReceiveChequeCmd,
		"receivelist":               ListReceiveChequeCmd,
		"receive-history-peer":      ChequeReceiveHistoryPeerCmd,
		"receive-history-list":      ChequeReceiveHistoryListCmd,
		"receive-history-stats":     ChequeReceiveHistoryStatsCmd,
		"receive-history-stats-all": ChequeReceiveHistoryStatsAllCmd,
		"receive-total-count":       ReceiveChequesCountCmd,
		"stats":                     ChequeStatsCmd,
		"stats-all":                 ChequeStatsAllCmd,

		"chaininfo":  ChequeChainInfoCmd,
		"bttbalance": ChequeBttBalanceCmd,
	},
}

var StorePriceCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get btfs store price.",
	},
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		// token: parse token option
		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		fmt.Println("... use token = ", tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		totalPrice, err := chain.SettleObject.OracleService.CheckNewPrice(token)
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
	},
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {

		// get the peer id
		peerID := req.Arguments[0]
		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		fmt.Printf("... token:%+v\n", tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

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
