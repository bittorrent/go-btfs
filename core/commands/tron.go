package commands

import (
	"encoding/hex"
	"strconv"

	cmds "github.com/TRON-US/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/wallet"
	"github.com/gogo/protobuf/proto"
)

var TronCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "BTFS tron",
		ShortDescription: `'btfs tron' is a set of commands to forward requests to trongrid node.`,
		LongDescription:  `'btfs tron' is a set of commands to forward requests to trongrid node.`,
	},

	Subcommands: map[string]*cmds.Command{
		"prepare": prepareCmd,
		"send":    sendCmd,
		"status":  statusCmd,
		"balance": balanceCmd,
	},
}

var prepareCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "prepare transaction",
		ShortDescription: "prepare transaction",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("from", true, false, "address for the token transfer from"),
		cmds.StringArg("to", true, false, "address for the token transfer to"),
		cmds.StringArg("amount", true, false, "amount of µBTT (=0.000001BTT) to transfer"),
		cmds.StringArg("memo", false, false, "memo"),
	},
	Options: []cmds.Option{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}
		amount, err := strconv.ParseInt(req.Arguments[2], 10, 64)
		if err != nil {
			return err
		}
		memo := ""
		if len(req.Arguments) >= 4 {
			memo = req.Arguments[3]
		}
		tx, err := wallet.PrepareTx(req.Context, cfg, req.Arguments[0], req.Arguments[1], amount, memo)
		if err != nil {
			return err
		}
		rawBytes, err := proto.Marshal(tx.Transaction.RawData)
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, &TxResult{
			TxId: hex.EncodeToString(tx.Txid),
			Raw:  hex.EncodeToString(rawBytes),
		})
	},
	Type: &TxResult{},
}

type TxResult struct {
	TxId string
	Raw  string
}

var sendCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "send raw transaction",
		ShortDescription: "send raw transaction",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("raw", true, false, "raw transaction in hex"),
		cmds.StringArg("sig", true, false, "sig of raw transaction in hex"),
	},
	Options: []cmds.Option{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}
		rawBytes, err := hex.DecodeString(req.Arguments[0])
		if err != nil {
			return err
		}
		sigBytes, err := hex.DecodeString(req.Arguments[1])
		if err != nil {
			return err
		}
		return wallet.SendRawTransaction(req.Context, cfg.Services.FullnodeDomain, rawBytes, sigBytes)
	},
}

var statusCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "get status of tx",
		ShortDescription: "get status of tx",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("txId", true, false, "transaction id"),
	},
	Options: []cmds.Option{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}
		status, err := wallet.GetStatus(req.Context, cfg.Services.SolidityDomain, req.Arguments[0])
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, map[string]string{"status": status})
	},
}

var balanceCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "get balance of pubaddress",
		ShortDescription: "get status of tx",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("pubAddress", true, false, "wallet address"),
	},
	Options: []cmds.Option{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cfg, err := n.Repo.Config()
		if err != nil {
			return err
		}
		amountMap, err := wallet.GetBalanceByWalletAddress(req.Context, cfg.Services.SolidityDomain, req.Arguments[0])
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, amountMap)
	},
}
