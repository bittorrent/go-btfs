package upload

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/protos/metadata"
	"github.com/ethereum/go-ethereum/common"

	"github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/core/corehttp/remote"

	"github.com/cenkalti/backoff/v4"
	"github.com/libp2p/go-libp2p/core/peer"
)

type ShardUploadContext struct {
	Rss            *sessions.RenterSession
	HostsProvider  helper.IHostsProvider
	Price          int64
	Token          common.Address
	ShardSize      int64
	StorageLength  int
	OfflineSigning bool
	RenterId       peer.ID
	FileSize       int64
	ShardIndexes   []int
	RepairParams   *RepairParams
}

func UploadShard(ctx *ShardUploadContext) error {
	if err := checkAndPreparePayment(ctx); err != nil {
		return err
	}
	for i, shardHash := range ctx.Rss.ShardHashes {
		go uploadSingleShard(ctx, ctx.ShardIndexes[i], shardHash)
	}

	complete, err := waitForAllShards(ctx)
	if err != nil {
		return err
	}
	if complete {
		err := Submit(ctx.Rss, ctx.FileSize, ctx.OfflineSigning)
		if err != nil {
			_ = ctx.Rss.To(sessions.RssToErrorEvent, err)
			return err
		}
	}
	return nil
}

func checkAndPreparePayment(ctx *ShardUploadContext) error {
	rate, err := chain.SettleObject.OracleService.CurrentRate(ctx.Token)
	if err != nil {
		return err
	}
	expectOnePay, err := helper.TotalPay(ctx.ShardSize, ctx.Price, ctx.StorageLength, rate)
	if err != nil {
		return err
	}
	expectTotalPay := expectOnePay * int64(len(ctx.Rss.ShardHashes))
	return checkAvailableBalance(ctx.Rss.Ctx, expectTotalPay, ctx.Token)
}

func uploadSingleShard(ctx *ShardUploadContext, shardIndex int, shardHash string) {
	err := backoff.Retry(func() error {
		if err := handleSingleShard(ctx, shardIndex, shardHash); err != nil {
			return err
		}
		return nil
	}, helper.HandleShardBo)
	if err != nil {
		_ = ctx.Rss.To(
			sessions.RssToErrorEvent,
			errors.New("timeout: failed to setup contract in "+helper.HandleShardBo.MaxElapsedTime.String()),
		)
	}
}

func handleSingleShard(ctx *ShardUploadContext, shardIndex int, shardHash string) error {
	select {
	case <-ctx.Rss.Ctx.Done():
		return nil
	default:
	}
	host, err := ctx.HostsProvider.NextValidHost()
	if err != nil {
		terr := ctx.Rss.To(sessions.RssToErrorEvent, err)
		if terr != nil {
			log.Debugf("original err: %s, transition err: %s", err.Error(), terr.Error())
		}
		return nil
	}
	hostPid, err := peer.Decode(host)
	if err != nil {
		log.Errorf("shard %s decodes host_pid error: %s", shardHash, err.Error())
		return err
	}
	if err := checkHostTokenSupport(ctx, hostPid); err != nil {
		return err
	}
	return processShardAgreementAndInit(ctx, host, hostPid, shardIndex, shardHash)
}

func checkHostTokenSupport(ctx *ShardUploadContext, hostPid peer.ID) error {
	c, cancel := context.WithTimeout(ctx.Rss.Ctx, 60*time.Second)
	defer cancel()
	output, err := remote.P2PCall(c, ctx.Rss.CtxParams.N, ctx.Rss.CtxParams.Api, hostPid, "/storage/upload/supporttokens")
	if err != nil {
		fmt.Printf("uploadShard, remote.P2PCall(supporttokens) timeout, hostPid = %v, will try again. \n", hostPid)
		return err
	}
	var mpToken map[string]common.Address
	err = json.Unmarshal(output, &mpToken)
	if err != nil {
		return err
	}
	for _, v := range mpToken {
		if ctx.Token == v {
			return nil
		}
	}
	return errors.New("host does not support token")
}

func processShardAgreementAndInit(ctx *ShardUploadContext, host string, hostPid peer.ID, shardIndex int, shardHash string) error {
	agreementID := helper.NewAgreementID(ctx.Rss.SsId)
	cb := make(chan error)
	ShardErrChanMap.Set(agreementID, cb)
	errChan := make(chan error, 2)
	var agreementBytes []byte
	go func() {
		errChan <- func() error {
			var err error
			agreementBytes, err = GetCreatorAgreement(
				ctx.Rss,
				&metadata.AgreementMeta{
					AgreementId:  agreementID,
					CreatorId:    ctx.RenterId.String(),
					SpId:         host,
					ShardIndex:   uint64(shardIndex),
					ShardHash:    shardHash,
					ShardSize:    uint64(ctx.ShardSize),
					Token:        ctx.Token.String(),
					StorageStart: uint64(time.Now().Unix()),
					StorageEnd:   uint64(time.Now().Add(time.Duration(ctx.StorageLength) * time.Second).Unix()),
					Price:        uint64(ctx.Price),
					Amount:       0, // expectOnePay 已在主流程校验
				},
				ctx.OfflineSigning,
				ctx.RepairParams,
				ctx.Token.String(),
			)
			if err != nil {
				log.Errorf("shard %s signs guard_contract error: %s", shardHash, err.Error())
				return err
			}
			return nil
		}()
	}()
	c := 0
	for err := range errChan {
		c++
		if err != nil {
			return err
		}
		if c >= 1 {
			break
		}
	}
	go func() {
		c, cancel := context.WithTimeout(ctx.Rss.Ctx, 10*time.Second)
		defer cancel()
		_, err := remote.P2PCall(
			c, ctx.Rss.CtxParams.N, ctx.Rss.CtxParams.Api, hostPid, "/storage/upload/init",
			ctx.Rss.SsId, ctx.Rss.Hash, shardHash, ctx.Price, agreementBytes, ctx.StorageLength, ctx.ShardSize, shardIndex, ctx.RenterId,
		)
		if err != nil {
			cb <- err
		}
	}()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	select {
	case err := <-cb:
		ShardErrChanMap.Remove(agreementID)
		return err
	case <-ticker.C:
		return errors.New("host timeout")
	}
}

// waitForAllShards 只负责等待分片完成并返回状态
// 返回 complete=true 表示全部分片完成，complete=false 表示有分片出错
func waitForAllShards(ctx *ShardUploadContext) (complete bool, err error) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			completeNum, errorNum, err := ctx.Rss.GetCompleteShardsNum()
			if err != nil {
				continue
			}
			log.Info("session", ctx.Rss.SsId, "agreementNum", completeNum, "errorNum", errorNum)
			if completeNum == len(ctx.Rss.ShardHashes) {
				return true, nil
			} else if errorNum > 0 {
				_ = ctx.Rss.To(sessions.RssToErrorEvent, errors.New("there are some error shards"))
				log.Error("session:", ctx.Rss.SsId, ",errorNum:", errorNum)
				return false, errors.New("there are some error shards")
			}
		case <-ctx.Rss.Ctx.Done():
			log.Infof("session %s done", ctx.Rss.SsId)
			return false, errors.New("session context done")
		}
	}
}
