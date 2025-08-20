package proxy

import (
	"fmt"
	"io"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
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
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		// if just notify the transaction hash

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
