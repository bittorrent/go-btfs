package proxy

import (
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
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
		hash := req.Arguments[0]
		txHash := common.HexToHash(hash)
		tx, _, err := chain.ChainObject.Backend.TransactionByHash(req.Context, txHash)
		if err != nil {
			return err
		}

		tx.ChainId()
		signer := types.NewEIP155Signer(tx.ChainId())
		from, err := types.Sender(signer, tx)
		if err != nil {
			return err
		}

		currentBalance, err := helper.GetBalance(req.Context, env.(*core.IpfsNode), from.String())
		if err != nil {
			return err
		}

		err = helper.PutProxyStoragePayment(req.Context, env.(*core.IpfsNode), &helper.ProxyStoragePaymentInfo{
			From:    from.String(),
			Hash:    tx.Hash().Hex(),
			PayTime: tx.Time().Unix(),
			To:      tx.To().Hex(),
			Value:   tx.Value().Uint64(),
			Balance: currentBalance + tx.Value().Uint64(),
		})

		if err != nil {
			return err
		}

		return nil
	},
}
