package proxy

import (
	"fmt"
	"io"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/bittorrent/go-btfs/utils"
	"github.com/ethereum/go-ethereum/common"
)

var StorageUploadProxyPayCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Deposit from beneficiary to vault contract account.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("recipient", false, false, "proxy account."),
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

		recipient := utils.RemoveSpaceAndComma(req.Arguments[0])
		if !common.IsHexAddress(recipient) {
			return fmt.Errorf("invalid bttc address %s", recipient)
		}
		recipientAddr := common.HexToAddress(recipient)

		argAmount := utils.RemoveSpaceAndComma(req.Arguments[0])
		amount, ok := new(big.Int).SetString(argAmount, 10)
		if !ok {
			return fmt.Errorf("amount:%s cannot be parsed", req.Arguments[0])
		}

		request := &transaction.TxRequest{
			To:    &recipientAddr,
			Value: amount,
		}
		hash, err := chain.ChainObject.TransactionService.Send(req.Context, request)
		if err != nil {
			return err
		}

		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		err = helper.PutProxyStoragePayment(req.Context, node, &helper.ProxyStoragePaymentInfo{
			From:    chain.ChainObject.TransactionService.SenderAddress(req.Context).String(),
			Hash:    hash.String(),
			PayTime: time.Now().Unix(),
			To:      recipient,
			Value:   amount.Uint64(),
		})
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

var StorageUploadProxyPaymentHistoryCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get the history of deposit from beneficiary to vault contract account.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("recipient", false, false, "proxy account."),
	},
	RunTimeout: 5 * time.Minute,
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
