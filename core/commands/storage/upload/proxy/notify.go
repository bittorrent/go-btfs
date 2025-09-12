package proxy

import (
	"errors"
	"fmt"
	"time"

	shell "github.com/bittorrent/go-btfs-api"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	"github.com/bittorrent/go-btfs/utils"
	coreiface "github.com/bittorrent/interface-go-btfs-core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ds "github.com/ipfs/go-datastore"
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
		cmds.StringArg("cid", true, false, "The cid that the transaction paid for"),
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

		// balance is wei
		currentBalance, err := helper.ChargeBalance(req.Context, nd, from.String(), tx.Value())
		if err != nil {
			return err
		}

		err = helper.PutProxyStoragePayment(req.Context, nd, &helper.ProxyStoragePaymentInfo{
			From:    from.String(),
			Hash:    tx.Hash().Hex(),
			PayTime: tx.Time().Unix(),
			To:      tx.To().Hex(),
			Value:   tx.Value(),
			Balance: currentBalance,
		})
		if err != nil {
			return err
		}

		// check if it is enough to pay
		needPayInfo, err := helper.GetProxyNeedPaymentCID(req.Context, nd, req.Arguments[1])
		if errors.Is(err, ds.ErrNotFound) {
			return fmt.Errorf("you do not need to pay for the cid: {%s} or it has been paid, but you btt has been deposited by the proxy", req.Arguments[1])
		}
		if err != nil {
			return fmt.Errorf("get need pay info error: %v", err)
		}
		if needPayInfo.NeedBTT.Cmp(currentBalance) == 1 {
			return fmt.Errorf("your payment is not enough for the %s to be uploaded by proxy", req.Arguments[1])
		}

		// upload file
		client := shell.NewLocalShell()
		_, err = client.Request("storage/upload", req.Arguments[1]).Send(req.Context)
		if err != nil {
			return err
		}

		_ = helper.SubBalance(req.Context, nd, from.String(), needPayInfo.NeedBTT)
		_ = helper.DeleteProxyNeedPaymentCID(req.Context, nd, req.Arguments[1])

		// save proxy upload cid
		ui := &helper.ProxyUploadFileInfo{
			From:      from.String(),
			CID:       req.Arguments[1],
			FileSize:  needPayInfo.FileSize,
			Price:     needPayInfo.Price,
			ExpireAt:  needPayInfo.ExpireAt,
			TotalPay:  needPayInfo.NeedBTT,
			CreatedAt: time.Now().Unix(),
		}
		_ = helper.PutProxyUploadedFileInfo(req.Context, nd, ui)

		return nil
	},
}
