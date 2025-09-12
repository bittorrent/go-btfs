package proxy

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/challenge"
	proxy "github.com/bittorrent/go-btfs/core/commands/storage/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	ds "github.com/ipfs/go-datastore"

	"github.com/cenkalti/backoff/v4"
	cidlib "github.com/ipfs/go-cid"
)

const (
	storageLengthOptionName = "storage-length"
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
		cmds.StringArg("file-hash", true, false, "Need to uploaded cid."),
	},
	Subcommands: map[string]*cmds.Command{
		"pay":        StorageUploadProxyPayCmd,
		"notify-pay": StorageUploadProxyNotifyPayCmd,
		"notify":     StorageUploadProxyNotifyCmd,
		"config":     StorageUploadProxyConfigCmd,
		"list":       StorageUploadFileListCmd,
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		ctxParams, err := helper.ExtractContextParams(req, env)
		if err != nil {
			return err
		}

		var shardHashes []string
		var fileSize int64
		var shardSize int64

		if !ctxParams.Cfg.Experimental.EnableProxyMode {
			return errors.New("proxy mode is not enabled")
		}

		cid, err := cidlib.Parse(req.Arguments[0])
		if err != nil {
			return err
		}

		shardHashes, fileSize, shardSize, err = helper.GetShardHashes(ctxParams, req.Arguments[0])
		if len(shardHashes) == 0 && fileSize == -1 && shardSize == -1 &&
			strings.HasPrefix(err.Error(), "invalid hash: file must be reed-solomon encoded") {
			shardHashes, fileSize, shardSize, err = helper.GetShardHashesCopy(ctxParams, req.Arguments[0], 0)
			fmt.Printf("copy get, shardHashes:%v fileSize:%v, shardSize:%v, copy:%v err:%v \n",
				shardHashes, fileSize, shardSize, 0, err)
		}
		if err != nil {
			return err
		}
		for _, s := range shardHashes {
			shardCid, err := cidlib.Decode(s)
			if err != nil {
				return err
			}
			scaledRetry := 30 * time.Minute
			err = backoff.Retry(func() error {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
				defer cancel()
				_, err = challenge.NewStorageChallengeResponse(ctx, ctxParams.N, ctxParams.Api, cid, shardCid, "", true, 1758693644)
				return err
			}, helper.DownloadShardBo(scaledRetry))
		}

		if err != nil {
			fmt.Println("download file error")
		}

		ctxParams.Req.Options = map[string]interface{}{
			storageLengthOptionName: 30,
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

		config, err := proxy.GetProxyStorageConfig(req.Context, ctxParams.N)
		if err != nil && errors.Is(err, ds.ErrNotFound) {
			// token: get new price
			priceObj, err := chain.SettleObject.OracleService.CurrentPrice(token)
			if err != nil {
				return err
			}
			config = &proxy.ProxyStorageInfo{}
			config.Price = uint64(priceObj.Int64())
		}
		if err != nil && !errors.Is(err, ds.ErrNotFound) {
			return err
		}
		// token: get new rate
		rate, err := chain.SettleObject.OracleService.CurrentRate(token)
		if err != nil {
			return err
		}
		totalPay, err := helper.TotalPay(shardSize, int64(config.Price), storageLength, rate)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		// save need pay cid and delete it when pay success
		payInfo := &proxy.ProxyNeedPaymentInfo{
			CID:      req.Arguments[0],
			FileSize: fileSize,
			Price:    int64(config.Price),
			NeedBTT:  new(big.Int).SetInt64(totalPay * rate.Int64()),
		}
		err = proxy.PutProxyNeedPaymentCID(ctxParams.Ctx, ctxParams.N, payInfo)
		if err != nil {
			return err
		}
		go func() {
			t := time.NewTimer(proxy.DefaultPayTimeout)
			select {
			case <-t.C:
				_ = proxy.DeleteProxyNeedPaymentCID(ctxParams.Ctx, ctxParams.N, req.Arguments[0])
			}
		}()

		proxyAddress, err := getPublicAddressFromPeerID(ctxParams.N.Identity.String())
		if err != nil {
			return err
		}
		t := new(big.Float).Quo(new(big.Float).SetInt(big.NewInt(totalPay)), big.NewFloat(1e6)).Text('f', 18)
		return res.Emit(map[string]interface{}{
			"proxy_address":   proxyAddress,
			"need_pay_amount": fmt.Sprintf("%s BTT", t), // convert to btt
		})
	},
}

var StorageUploadFileListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List files that uploaded by proxy.",
		ShortDescription: `
This command list files that uploaded by proxy.`,
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		infos, err := proxy.ListProxyUploadedFileInfo(req.Context, n)
		if err != nil {
			return err
		}
		for _, info := range infos {
			info.Price = info.Price / 1e6
		}
		return cmds.EmitOnce(res, infos)
	},
	Type: []*proxy.ProxyUploadFileInfo{},
}
