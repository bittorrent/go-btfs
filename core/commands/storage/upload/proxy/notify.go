package proxy

import (
	"errors"
	"math/big"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	"github.com/bittorrent/go-btfs/utils"
	coreiface "github.com/bittorrent/interface-go-btfs-core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var StorageUploadProxyNotifyPayCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Notify the proxy that the payment has been made.",
		LongDescription: `
This command is used to notify the proxy that the payment has been made.
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("hash", true, false, "The hash of the storage-upload-proxy-pay command."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		if !nd.IsOnline {
			return coreiface.ErrOffline
		}
		// check simple mode
		err = utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}
		hash := req.Arguments[0]
		txHash := common.HexToHash(hash)
		tx, _, err := chain.ChainObject.Backend.TransactionByHash(req.Context, txHash)
		if err != nil {
			return err
		}

		conf, err := nd.Repo.Config()
		if err != nil {
			return err
		}

		// check if the tx is for me
		if tx.To().String() != conf.Identity.BttcAddr {
			return nil
		}

		signer := types.NewEIP155Signer(tx.ChainId())
		from, err := types.Sender(signer, tx)
		if err != nil {
			return err
		}
		// check if the tx has been paid
		d, err := helper.GetProxyStoragePaymentByTxHash(req.Context, nd, from.String(), tx.Hash().Hex())
		if err != nil {
			return nil
		}

		if d != nil && d.Hash == tx.Hash().Hex() {
			return errors.New("the tx hash has been notified")
		}

		value := new(big.Int).Div(tx.Value(), big.NewInt(1e18))
		currentBalance, err := helper.ChargeBalance(req.Context, nd, from.String(), value.Uint64())
		if err != nil {
			return err
		}

		err = helper.PutProxyStoragePayment(req.Context, nd, &helper.ProxyStoragePaymentInfo{
			From:    from.String(),
			Hash:    tx.Hash().Hex(),
			PayTime: tx.Time().Unix(),
			To:      tx.To().Hex(),
			Value:   value.Uint64(),
			Balance: currentBalance,
		})

		if err != nil {
			return err
		}

		return nil
	},
}
