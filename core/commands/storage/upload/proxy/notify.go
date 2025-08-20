package proxy

import (
	"fmt"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/ethereum/go-ethereum/common"
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
		tx, isPending, err := chain.ChainObject.Backend.TransactionByHash(req.Context, txHash)
		if err != nil {
			return err
		}

		fmt.Println("tx hash:", tx.Hash().Hex())
		fmt.Println("is pending:", isPending)
		fmt.Println("To:", tx.To().Hex())
		fmt.Println("Value:", tx.Value().String())
		fmt.Println("Gas Limit:", tx.Gas())
		fmt.Println("Gas Price:", tx.GasPrice())

		return nil

	},
}
