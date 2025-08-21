package proxy

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/core/commands/storage/challenge"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/cenkalti/backoff/v4"
	cidlib "github.com/ipfs/go-cid"
)

var StorageUploadProxyCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Initialize storage handshake with inquiring client.",
		ShortDescription: `
Storage host opens this endpoint to accept incoming upload/storage requests,
If current host is interested and all validation checks out, host downloads
the shard and replies back to client for the next challenge step.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("proxy-id", true, false, "ProxyId for upload file"),
		cmds.StringArg("file-hash", true, false, "Root file storage node should fetch (the DAG)."),
	},
	Subcommands: map[string]*cmds.Command{
		"pay":        StorageUploadProxyPayCmd,
		"notify-pay": StorageUploadProxyNotifyPayCmd,
		"config":     StorageUploadProxyConfigCmd,
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		fmt.Println("storage proxy do ..............")
		fmt.Println(req.Arguments)
		fmt.Println(req.Options)

		ctxParams, err := helper.ExtractContextParams(req, env)
		if err != nil {
			return err
		}

		var shardHashes []string
		var fileSize int64
		var shardSize int64

		fileHash := req.Arguments[0]

		if !ctxParams.Cfg.Experimental.EnableProxyMode {
			return errors.New("proxy mode is not enabled")
		}
		if req.Arguments[1] == ctxParams.N.Identity.String() {
			// TODO check react as a proxy node
			fileCid, err := cidlib.Parse(req.Arguments[0])
			if err != nil {
				return err
			}
			shardCid, err := cidlib.Decode("")
			if err != nil {
			}

			scaledRetry := 30 * time.Minute
			err = backoff.Retry(func() error {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
				defer cancel()
				_, err = challenge.NewStorageChallengeResponse(ctx, ctxParams.N, ctxParams.Api, fileCid, shardCid, "", false, 0)
				return err
			}, helper.DownloadShardBo(scaledRetry))

			if err != nil {
				fmt.Println("download file error")
			}

		}

		shardHashes, fileSize, shardSize, err = helper.GetShardHashes(ctxParams, fileHash)

		if len(shardHashes) == 0 && fileSize == -1 && shardSize == -1 &&
			strings.HasPrefix(err.Error(), "invalid hash: file must be reed-solomon encoded") {
			if copyNum, ok := req.Options["copy"].(int); ok {
				shardHashes, fileSize, shardSize, err = helper.GetShardHashesCopy(ctxParams, fileHash, copyNum)
				fmt.Printf("copy get, shardHashes:%v fileSize:%v, shardSize:%v, copy:%v err:%v \n",
					shardHashes, fileSize, shardSize, copyNum, err)
			}
		}
		if err != nil {
			return err
		}
		_, storageLength, err := helper.GetPriceAndMinStorageLength(ctxParams)
		if err != nil {
			return err
		}

		tokenStr := "WBTT"
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}
		fmt.Println("token =", token, tokenStr)

		// token: get new price
		priceObj, err := chain.SettleObject.OracleService.CurrentPrice(token)
		if err != nil {
			return err
		}
		price := priceObj.Int64()
		// token: get new rate
		rate, err := chain.SettleObject.OracleService.CurrentRate(token)
		if err != nil {
			return err
		}
		totalPay, err := helper.TotalPay(shardSize, price, storageLength, rate)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		res.Emit(map[string]interface{}{
			"peer_address": ctxParams.N.Identity.String(),
			"total_amount": totalPay,
		})

		return nil
	},
}
