package upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/protos/metadata"

	"github.com/alecthomas/units"
	"github.com/cenkalti/backoff/v4"
)

const (
	thresholdContractsNums = 20
)

func getSuccessThreshold(totalShards int) int {
	return int(math.Min(float64(totalShards), thresholdContractsNums))
}

func ResumeWaitUploadOnSigning(rss *sessions.RenterSession) error {
	// return waitUpload(rss, false, &guardpb.FileStoreStatus{
	// FileStoreMeta: guardpb.FileStoreMeta{
	// RenterPid: rss.CtxParams.N.Identity.String(),
	// FileSize:  math.MaxInt64,
	// },
	// }, true)
	return nil
}

func waitSPSaveFileSuccAndToPay(rss *sessions.RenterSession, offlineSigning bool, fsStatus *metadata.FileMetaInfo, resume bool) error {
	threshold := getSuccessThreshold(len(rss.ShardHashes))
	if !resume {
		if err := rss.To(sessions.RssToWaitUploadEvent); err != nil {
			return err
		}
	}
	// req := &guardpb.CheckFileStoreMetaRequest{
	// 	FileHash:     rss.Hash,
	// 	RenterPid:    fsStatus.RenterPid,
	// 	RequesterPid: fsStatus.RenterPid,
	// 	RequestTime:  time.Now().UTC(),
	// }
	// payerPrivKey, err := rss.CtxParams.Cfg.Identity.DecodePrivateKey("")
	// if err != nil {
	// return err
	// }
	cb := make(chan []byte)
	helper.WaitUploadChanMap.Set(rss.SsId, cb)
	// if offlineSigning {
	// 	raw, err := proto.Marshal(req)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	err = rss.SaveOfflineSigning(&renterpb.OfflineSigning{
	// 		Raw: raw,
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// } else {
	// 	go func() {
	// 		sign, err := crypto.Sign(payerPrivKey, req)
	// 		if err != nil {
	// 			_ = rss.To(sessions.RssToErrorEvent, err)
	// 			return
	// 		}
	// 		cb <- sign
	// 	}()
	// }
	// sign := <-cb
	helper.WaitUploadChanMap.Remove(rss.SsId)
	if !resume {
		if err := rss.To(sessions.RssToWaitUploadReqSignedEvent); err != nil {
			return err
		}
	}
	// req.Signature = sign
	lowRetry := 30 * time.Minute
	highRetry := 24 * time.Hour
	scaledRetry := time.Duration(float64(fsStatus.FileSize) / float64(units.GiB) * float64(highRetry))
	if scaledRetry < lowRetry {
		scaledRetry = lowRetry
	} else if scaledRetry > highRetry {
		scaledRetry = highRetry
	}
	var contracts []string
	for _, c := range fsStatus.Contracts {
		contracts = append(contracts, c.Meta.ContractId)
	}
	err := backoff.Retry(func() error {
		meta, err := chain.SettleObject.FileMetaService.GetFileMeta(rss.Hash, contracts)
		if err != nil {
			return err
		}
		num := 0
		m := make(map[string]int)
		for _, c := range meta.Contracts {
			m[c.Status.String()]++
			switch c.Status {
			case metadata.Contract_COMPLETED:
				num++
			}
			shard, err := sessions.GetUserShard(rss.CtxParams, rss.SsId, c.Meta.ShardHash, int(c.Meta.ShardIndex))
			if err != nil {
				return err
			}
			err = shard.UpdateAdditionalInfo(c.Status.String())
			if err != nil {
				return err
			}
			err = shard.UpdateContractsStatus()
			if err != nil {
				return err
			}
		}
		bytes, err := json.Marshal(m)
		if err == nil {
			_ = rss.UpdateAdditionalInfo(string(bytes))
		}
		log.Infof("%d shards uploaded.", num)
		if num >= threshold {
			return nil
		}
		return errors.New("uploading")
	}, helper.WaitUploadBo(highRetry))
	if err != nil {
		return err
	}

	// pay in cheque
	if err := rss.To(sessions.RssToPayEvent); err != nil {
		return err
	}
	var errC = make(chan error)
	go func() {
		err = func() error {
			return payInCheque(rss)
		}()
		if err != nil {
			fmt.Println("payInCheque error:", err)
		}
		fmt.Println("payInCheque done")
		errC <- err
	}()
	err = <-errC
	if err != nil {
		if fsmErr := rss.To(sessions.RssToErrorEvent); fsmErr != nil {
			log.Errorf("fsm transfer error:%v", fsmErr)
		}
		log.Errorf("payInCheque error:%v", err)
		return err
	}
	// Complete
	if err := rss.To(sessions.RssToCompleteEvent); err != nil {
		return err
	}
	return nil
}
