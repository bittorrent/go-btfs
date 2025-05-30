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

func UploadShard(rss *sessions.RenterSession, hp helper.IHostsProvider, price int64, token common.Address, shardSize int64,
	storageLength int,
	offlineSigning bool, renterId peer.ID, fileSize int64, shardIndexes []int, rp *RepairParams) error {

	// token: get new rate
	rate, err := chain.SettleObject.OracleService.CurrentRate(token)
	if err != nil {
		return err
	}
	expectOnePay, err := helper.TotalPay(shardSize, price, storageLength, rate)
	if err != nil {
		return err
	}
	expectTotalPay := expectOnePay * int64(len(rss.ShardHashes))
	err = checkAvailableBalance(rss.Ctx, expectTotalPay, token)
	if err != nil {
		return err
	}

	// 正常只有一个，如果是使用了reed-solomon-4-2-1048576算法就可能是多个了
	// 如果是指定了副本数量的话也是多个, 但是每个都一样的
	for index, shardHash := range rss.ShardHashes {
		// 每个shard使用一个协程去处理
		go func(i int, h string) {
			err := backoff.Retry(func() error {
				select {
				case <-rss.Ctx.Done():
					return nil
				default:
					break
				}
				host, err := hp.NextValidHost()
				if err != nil {
					terr := rss.To(sessions.RssToErrorEvent, err)
					if terr != nil {
						// Ignore err, just print error log
						log.Debugf("original err: %s, transition err: %s", err.Error(), terr.Error())
					}
					return nil
				}

				hostPid, err := peer.Decode(host)
				if err != nil {
					log.Errorf("shard %s decodes host_pid error: %s", h, err.Error())
					return err
				}

				//token: check host tokens
				{
					ctx, _ := context.WithTimeout(rss.Ctx, 60*time.Second)
					// TODO sp节点需要支持这个
					output, err := remote.P2PCall(ctx, rss.CtxParams.N, rss.CtxParams.Api, hostPid, "/storage/upload/supporttokens")
					if err != nil {
						fmt.Printf("uploadShard, remote.P2PCall(supporttokens) timeout, hostPid = %v, will try again. \n", hostPid)
						return err
					}

					var mpToken map[string]common.Address
					err = json.Unmarshal(output, &mpToken)
					if err != nil {
						return err
					}

					ok := false
					for _, v := range mpToken {
						if token == v {
							ok = true
						}
					}
					if !ok {
						return nil
					}
				}

				agreementID := helper.NewAgreementID(rss.SsId)
				cb := make(chan error)
				ShardErrChanMap.Set(agreementID, cb)

				errChan := make(chan error, 2)
				var agreementBytes []byte

				// 对shard进行签名
				go func() {
					tmp := func() error {
						agreementBytes, err = GetCreatorAgreement(
							rss,
							&metadata.AgreementMeta{
								AgreementId:  agreementID,
								CreatorId:    renterId.String(),
								SpId:         host,
								ShardIndex:   uint64(i),
								ShardHash:    h,
								ShardSize:    uint64(shardSize),
								Token:        token.String(),
								StorageStart: uint64(time.Now().Unix()),
								StorageEnd:   uint64(time.Now().Add(time.Duration(storageLength) * time.Second).Unix()),
								Price:        uint64(price),
								Amount:       uint64(expectOnePay),
							},
							offlineSigning,
							rp,
							token.String(),
						)
						if err != nil {
							log.Errorf("shard %s signs guard_contract error: %s", h, err.Error())
							return err
						}
						return nil
					}()
					errChan <- tmp
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

				// 发送合同 给 host
				go func() {
					ctx, _ := context.WithTimeout(rss.Ctx, 10*time.Second)
					// TODO 这里调用SP的init接口，SP需要调整
					_, err := remote.P2PCall(ctx, rss.CtxParams.N, rss.CtxParams.Api, hostPid, "/storage/upload/init",
						rss.SsId,
						rss.Hash,
						h,
						price,
						nil,
						agreementBytes,
						storageLength,
						shardSize,
						i,
						renterId,
					)
					if err != nil {
						cb <- err
					}
				}()
				// host needs to send recv in 30 seconds, or the contract will be invalid.
				tick := time.Tick(30 * time.Second)
				select {
				case err = <-cb:
					ShardErrChanMap.Remove(agreementID)
					return err
				case <-tick:
					return errors.New("host timeout")
				}
			}, helper.HandleShardBo)
			if err != nil {
				_ = rss.To(sessions.RssToErrorEvent,
					errors.New("timeout: failed to setup contract in "+helper.HandleShardBo.MaxElapsedTime.String()))
			}
		}(shardIndexes[index], shardHash)
	}
	// waiting for contracts of 30(n) shards
	go func(rss *sessions.RenterSession, numShards int) {
		tick := time.Tick(5 * time.Second)
		for true {
			select {
			case <-tick:
				completeNum, errorNum, err := rss.GetCompleteShardsNum()
				if err != nil {
					continue
				}
				log.Info("session", rss.SsId, "contractNum", completeNum, "errorNum", errorNum)
				if completeNum == numShards {
					// while all shards upload completely, submit its.
					err := Submit(rss, fileSize, offlineSigning)
					if err != nil {
						_ = rss.To(sessions.RssToErrorEvent, err)
					}
					return
				} else if errorNum > 0 {
					_ = rss.To(sessions.RssToErrorEvent, errors.New("there are some error shards"))
					log.Error("session:", rss.SsId, ",errorNum:", errorNum)
					return
				}
			case <-rss.Ctx.Done():
				log.Infof("session %s done", rss.SsId)
				return
			}
		}
	}(rss, len(rss.ShardHashes))

	return nil
}
