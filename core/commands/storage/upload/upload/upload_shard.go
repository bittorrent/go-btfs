package upload

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/renewal"
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
	AutoRenewal    bool
	RepairParams   *RepairParams
	TotalPay       *big.Int
}

func UploadShard(ctx *ShardUploadContext) error {
	expectOnePay, err := checkAndPreparePayment(ctx)
	if err != nil {
		return err
	}
	for i, shardHash := range ctx.Rss.ShardHashes {
		h := shardHash
		index := i
		go sendShardContractToHost(ctx, ctx.ShardIndexes[index], h, expectOnePay)
	}

	go func() {
		isComplete, err := waitForAllShardsComplete(ctx)
		if err != nil {
			log.Errorf("wait for all shards complete error: %s", err.Error())
			return
		}

		if isComplete {
			// set rss status from init to submit
			if err := ctx.Rss.To(sessions.RssToSubmitEvent); err != nil {
				log.Errorf("set rss status from init to submit error: %s", err.Error())
				return
			}
			err := SubmitToChain(ctx.Rss, ctx.FileSize, ctx.OfflineSigning)
			if err != nil {
				_ = ctx.Rss.To(sessions.RssToErrorEvent, err)
				return
			}

			// save auto-renewal info
			shardsInfo := make([]*renewal.RenewalShardInfo, 0)
			for i, shard := range ctx.Rss.ShardHashes {
				shards, err := sessions.GetUserShard(ctx.Rss.CtxParams, ctx.Rss.SsId, shard, i)
				if err != nil {
					log.Errorf("get user shard error: %s", err.Error())
					continue
				}
				contracts, err := shards.Contracts()
				if err != nil {
					log.Errorf("get contracts error: %s", err.Error())
					continue
				}
				si := &renewal.RenewalShardInfo{
					SPId:       contracts.Meta.SpId,
					ShardId:    contracts.Meta.ShardHash,
					ShardSize:  int(contracts.Meta.ShardSize),
					ContractID: contracts.Meta.ContractId,
				}
				shardsInfo = append(shardsInfo, si)
			}
			info := &renewal.RenewalInfo{
				CID:             ctx.Rss.Hash,
				RenewalDuration: ctx.StorageLength,
				Token:           ctx.Token,
				Price:           ctx.Price,
				Enabled:         ctx.AutoRenewal,
				CreatedAt:       time.Now(),
				LastRenewalAt:   nil,
				NextRenewalAt:   time.Now().Add(time.Duration(ctx.StorageLength) * 24 * time.Hour),
				ShardsInfo:      shardsInfo,
				TotalPay:        ctx.TotalPay,
			}
			err = renewal.StoreRenewalInfo(ctx.Rss.CtxParams, info, renewal.RenewTypeAuto)
			if err != nil {
				log.Errorf("Failed to store auto-renewal config: %v", err)
				return
			}
		}
	}()
	return nil
}

func checkAndPreparePayment(ctx *ShardUploadContext) (int64, error) {
	rate, err := chain.SettleObject.OracleService.CurrentRate(ctx.Token)
	if err != nil {
		return 0, err
	}
	expectOnePay, err := helper.TotalPay(ctx.ShardSize, ctx.Price, ctx.StorageLength, rate)
	if err != nil {
		return 0, err
	}
	expectTotalPay := expectOnePay * int64(len(ctx.Rss.ShardHashes))
	return expectOnePay, checkAvailableBalance(ctx.Rss.Ctx, expectTotalPay, ctx.Token)
}

func sendShardContractToHost(ctx *ShardUploadContext, shardIndex int, shardHash string, amount int64) {
	err := backoff.Retry(func() error {
		if err := sendShardContract(ctx, shardIndex, shardHash, amount); err != nil {
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

func sendShardContract(ctx *ShardUploadContext, shardIndex int, shardHash string, amount int64) error {
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
	return signShardContractAndSendToSP(ctx, host, hostPid, shardIndex, shardHash, amount)
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

func signShardContractAndSendToSP(ctx *ShardUploadContext, host string, hostPid peer.ID, shardIndex int, shardHash string, amount int64) error {
	contractID := helper.NewContractID(ctx.Rss.SsId)
	cb := make(chan error)
	ShardErrChanMap.Set(contractID, cb)
	errChan := make(chan error, 2)
	var signedContractBytes []byte
	go func() {
		errChan <- func() error {
			var err error
			signedContractBytes, err = SignUserContract(
				ctx.Rss,
				&metadata.ContractMeta{
					ContractId:   contractID,
					UserId:       ctx.RenterId.String(),
					SpId:         host,
					ShardIndex:   uint64(shardIndex),
					ShardHash:    shardHash,
					ShardSize:    uint64(ctx.ShardSize),
					Token:        ctx.Token.String(),
					StorageStart: uint64(time.Now().Unix()),
					StorageEnd:   uint64(time.Now().Add(time.Duration(ctx.StorageLength) * 24 * time.Hour).Unix()),
					Price:        uint64(ctx.Price),
					Amount:       uint64(amount),
					AutoRenewal:  ctx.AutoRenewal,
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
			ctx.Rss.SsId, ctx.Rss.Hash, shardHash, ctx.Price, signedContractBytes, ctx.StorageLength, ctx.ShardSize, shardIndex, ctx.RenterId,
		)
		if err != nil {
			cb <- err
		}
	}()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	select {
	case err := <-cb:
		ShardErrChanMap.Remove(contractID)
		return err
	case <-ticker.C:
		return errors.New("host timeout")
	}
}

func waitForAllShardsComplete(ctx *ShardUploadContext) (complete bool, err error) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			completeNum, errorNum, err := ctx.Rss.GetCompleteShardsNum()
			if err != nil {
				continue
			}
			log.Info("session", ctx.Rss.SsId, "contractNum", completeNum, "errorNum", errorNum)
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
