package proxy

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	shell "github.com/bittorrent/go-btfs-api"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	"github.com/bittorrent/go-btfs/core/corehttp/remote"
	"github.com/bittorrent/go-btfs/utils"
	coreiface "github.com/bittorrent/interface-go-btfs-core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ds "github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p/core/peer"
)

var StorageUploadProxyNotifyPayCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Notify the proxy that the payment has been made.",
		LongDescription: `
This command is used to notify the proxy that the payment has been made.
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("hash", false, false, "The hash of the storage-upload-proxy-pay command."),
		cmds.StringArg("cid", false, false, "The cid that the transaction paid for"),
		cmds.StringArg("address", false, false, "The address that the payment is paid for"),
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

		var currentBalance *big.Int
		var from string
		var cid string

		hash := req.Arguments[0]
		if hash != "" {
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
			f, err := types.Sender(signer, tx)
			if err != nil {
				return err
			}
			from = f.String()
			// check if the tx has been notified
			d, err := helper.GetProxyStoragePaymentByTxHash(req.Context, nd, from, tx.Hash().Hex())
			if err != nil {
				return nil
			}

			if d != nil && d.Hash == tx.Hash().Hex() {
				return errors.New("the tx hash has been notified")
			}

			// save payment record
			err = helper.PutProxyStoragePayment(req.Context, nd, &helper.ProxyStoragePaymentInfo{
				From:    from,
				Hash:    tx.Hash().Hex(),
				PayTime: tx.Time().Unix(),
				To:      tx.To().Hex(),
				Value:   tx.Value(),
				Balance: currentBalance,
			})

			// charge balance is wei
			currentBalance, err = helper.ChargeBalance(req.Context, nd, from, tx.Value())
			if err != nil {
				return err
			}
		}

		// if cid is empty, just pay
		cid = req.Arguments[1]
		if cid == "" {
			return nil
		}

		if from == "" {
			from = req.Arguments[2]
		}
		currentBalance, err = helper.GetBalance(req.Context, nd, from)
		if err != nil {
			return err
		}

		// check if it is enough to pay
		needPayInfo, err := helper.GetProxyNeedPaymentCID(req.Context, nd, cid)
		if errors.Is(err, ds.ErrNotFound) {
			return fmt.Errorf("you do not need to pay for the cid: {%s} or it has been paid, but you btt has been deposited by the proxy", cid)
		}
		if err != nil {
			return fmt.Errorf("get need pay info error: %v", err)
		}
		if needPayInfo.NeedBTT.Cmp(currentBalance) == 1 {
			return fmt.Errorf("your payment is not enough for the %s to be uploaded by proxy", cid)
		}

		// upload file
		client := shell.NewLocalShell()
		resp, err := client.Request("storage/upload", cid).Send(context.Background())
		if err != nil {
			return err
		}
		defer resp.Close()
		type uploadResponse struct {
			ID string `json:"ID"`
		}

		if resp.Error != nil {
			fmt.Printf("Upload error: %s\n", resp.Error.Message)
			return fmt.Errorf("upload failed: %s", resp.Error.Message)
		}

		var uploadResp uploadResponse
		err = resp.Decode(&uploadResp)
		if err != nil {
			fmt.Printf("Failed to decode response: %v\n", err)
			return err
		}

		go func() {
			// wait for upload success
			ticker := time.NewTicker(time.Second * 10)
			defer ticker.Stop()
			for range ticker.C {
				resp, err := client.Request("storage/upload/status", uploadResp.ID).Send(context.Background())
				if err != nil {
					fmt.Println("proxy get upload status error: ", err)
				}
				if resp.Error != nil {
					fmt.Println("proxy get upload status error: ", resp.Error)
				}
				type StatusRes struct {
					Status   string
					Message  string
					FileHash string
				}
				// parse response
				var statusResp StatusRes
				err = resp.Decode(&statusResp)
				if err != nil {
					fmt.Println("proxy parse upload status error: ", err)
				}
				if statusResp.Status == "complete" {
					fmt.Println("upload file success")
					_ = helper.SubBalance(req.Context, nd, from, needPayInfo.NeedBTT)
					_ = helper.DeleteProxyNeedPaymentCID(req.Context, nd, req.Arguments[1])

					// save proxy upload cid
					ui := &helper.ProxyUploadFileInfo{
						From:     from,
						CID:      req.Arguments[1],
						FileSize: needPayInfo.FileSize,
						Price:    needPayInfo.Price,
						// ExpireAt:  needPayInfo.ExpireAt,
						TotalPay:  needPayInfo.NeedBTT,
						CreatedAt: time.Now().Unix(),
					}
					_ = helper.PutProxyUploadedFileInfo(req.Context, nd, ui)
					return
				}
				if statusResp.Status == "error" {
					fmt.Println("proxy upload file error: ", statusResp.Message)
					return
				}
			}
		}()

		return nil
	},
}

// client used to notify proxy it will call proxy notify pay method

var StorageUploadProxyNotifyCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Notify the proxy to upload file",
		LongDescription: `
This command is used to notify the proxy to upload file to SP.
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("proxy-id", true, false, "The proxy id that will be notified to upload file"),
		cmds.StringArg("cid", true, false, "The cid that need to be uploaded"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		proxyId := req.Arguments[0]
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

		address, err := getPublicAddressFromPeerID(node.Identity.String())
		if err != nil {
			return err
		}
		// notify the proxy payment
		_, err = remote.P2PCall(req.Context, node, api, pId, "/storage/upload/proxy/notify-pay",
			"",
			req.Arguments[1],
			address,
		)

		return err
	},
}
