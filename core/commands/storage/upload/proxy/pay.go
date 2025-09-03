package proxy

import (
	"fmt"
	"io"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	cp "github.com/bittorrent/go-btfs-common/crypto"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	"github.com/bittorrent/go-btfs/core/corehttp/remote"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/bittorrent/go-btfs/utils"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

var StorageUploadProxyPayCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Deposit from beneficiary to vault contract account.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("proxy-id", false, false, "proxy peerId."),
		cmds.StringArg("cid", false, false, "cid that need to pay."),
		cmds.StringArg("amount", false, false, "deposit amount of BTT."),
	},
	RunTimeout: 5 * time.Minute,
	Subcommands: map[string]*cmds.Command{
		"balance": StorageUploadProxyPaymentBalanceCmd,
		"history": StorageUploadProxyPaymentHistoryCmd,
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		proxyId := utils.RemoveSpaceAndComma(req.Arguments[0])

		proxyAddr, err := getPublicAddressFromPeerID(proxyId)
		if err != nil {
			return err
		}

		argAmount := utils.RemoveSpaceAndComma(req.Arguments[2])
		amount, ok := new(big.Int).SetString(argAmount, 10)
		if !ok {
			return fmt.Errorf("amount:%s cannot be parsed", req.Arguments[2])
		}

		to := common.HexToAddress(proxyAddr)
		request := &transaction.TxRequest{
			To: &to,
			// convert wei to btt
			Value: new(big.Int).Mul(amount, big.NewInt(1e18)),
		}
		hash, err := chain.ChainObject.TransactionService.Send(req.Context, request)
		if err != nil {
			return err
		}

		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		api, err := cmdenv.GetApi(env, req)
		if err != nil {
			return err
		}

		pId, err := peer.Decode(proxyId)
		if err != nil {
			fmt.Println("invalid peer id:", err)
			return err
		}

		// notify the proxy payment
		_, err = remote.P2PCall(req.Context, node, api, pId, "/storage/upload/proxy/notify-pay",
			hash,
			req.Arguments[1],
		)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &PayCmdRet{
			Hash: hash.String(),
		})
	},
	Type: &PayCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *PayCmdRet) error {
			_, err := fmt.Fprintf(w, "the hash of transaction: %s\n", out.Hash)
			return err
		}),
	},
}

type PayCmdRet struct {
	Hash string `json:"hash"`
}

var StorageUploadProxyPaymentBalanceCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get the balance of deposit from beneficiary to vault contract account.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("recipient", false, false, "proxy account."),
	},
	RunTimeout: 5 * time.Minute,
	Type:       []*BalanceCmdRet{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		recipient := ""

		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		if len(req.Arguments) > 0 {
			recipient = utils.RemoveSpaceAndComma(req.Arguments[0])
			if !common.IsHexAddress(recipient) {
				return fmt.Errorf("invalid bttc address %s", recipient)
			}
		}

		if recipient != "" {
			balance, err := helper.GetBalance(req.Context, node, recipient)
			if err != nil {
				return err
			}
			return cmds.EmitOnce(res, []*BalanceCmdRet{
				{
					Address: recipient,
					Balance: fmt.Sprintf("%d (BTT)", balance),
				},
			})
		}

		balances, err := helper.GetBalanceList(req.Context, node)
		if err != nil {
			return err
		}

		result := make([]*BalanceCmdRet, 0)

		for k, v := range balances {
			ret := &BalanceCmdRet{
				Address: k,
				Balance: fmt.Sprintf("%d (BTT)", v),
			}

			result = append(result, ret)
		}

		return cmds.EmitOnce(res, result)
	},
}

type BalanceCmdRet struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

var StorageUploadProxyPaymentHistoryCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get the history of deposit from beneficiary to vault contract account.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("recipient", false, false, "proxy account."),
	},
	RunTimeout: 5 * time.Minute,
	Type:       []*PaymentHistoryCmdRet{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {

		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		recipient := ""
		if len(req.Arguments) > 0 {
			recipient = utils.RemoveSpaceAndComma(req.Arguments[0])
			if !common.IsHexAddress(recipient) {
				return fmt.Errorf("invalid bttc address %s", recipient)
			}
		}

		payments := make([]*PaymentHistoryCmdRet, 0)
		if recipient != "" {
			ps, err := helper.GetProxyStoragePaymentList(req.Context, node, recipient)
			if err != nil {
				return err
			}
			for _, p := range ps {
				ret := &PaymentHistoryCmdRet{
					Hash:    p.Hash,
					From:    p.From,
					To:      p.To,
					Value:   fmt.Sprintf("%d (BTT)", p.Value),
					PayTime: p.PayTime,
				}
				payments = append(payments, ret)
			}
		} else {
			ps, err := helper.GetProxyStoragePayment(req.Context, node)
			if err != nil {
				return err
			}
			for _, p := range ps {
				ret := &PaymentHistoryCmdRet{
					Hash:    p.Hash,
					From:    p.From,
					To:      p.To,
					Value:   fmt.Sprintf("%d (BTT)", p.Value),
					PayTime: p.PayTime,
				}
				payments = append(payments, ret)
			}

		}

		return cmds.EmitOnce(res, payments)
	},
}

type PaymentHistoryCmdRet struct {
	Hash    string `json:"hash"`
	From    string `json:"from"`
	To      string `json:"to"`
	Value   string `json:"value"`
	PayTime int64  `json:"pay_time"`
}

func getPublicAddressFromPeerID(hostID string) (string, error) {
	peerID, err := peer.Decode(hostID)
	if err != nil {
		return "", fmt.Errorf("failed to decode hostID: %v", err)
	}

	pubKey, err := peerID.ExtractPublicKey()
	if err != nil {
		return "", fmt.Errorf("failed to extract public key: %v", err)
	}

	pkBytes, err := cp.Secp256k1PublicKeyRaw(pubKey)
	if err != nil {
		panic(err)
	}

	ethPk, err := ethCrypto.UnmarshalPubkey(pkBytes)
	if err != nil {
		return "", err
	}

	return ethCrypto.PubkeyToAddress(*ethPk).Hex(), nil

}
