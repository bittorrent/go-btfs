package upload

import (
	"context"
	"fmt"

	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/guard"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/protos/metadata"
	renterpb "github.com/bittorrent/go-btfs/protos/renter"

	"github.com/bittorrent/go-btfs-common/crypto"
	guardpb "github.com/bittorrent/go-btfs-common/protos/guard"
	cgrpc "github.com/bittorrent/go-btfs-common/utils/grpc"
	config "github.com/bittorrent/go-btfs-config"

	"github.com/gogo/protobuf/proto"
)

func doAgreementAndPay(rss *sessions.RenterSession, fileSize int64, offlineSigning bool) error {
	// 事件驱动状态机流转
	if err := rss.To(sessions.RssToGuardEvent); err != nil {
		return err
	}

	as := make([]*metadata.Agreement, 0)
	selectedSPs := make([]string, 0)

	for i, h := range rss.ShardHashes {
		shard, err := sessions.GetRenterShard(rss.CtxParams, rss.SsId, h, i)
		if err != nil {
			return err
		}
		// 获取签名了的agreement
		agreement, err := shard.Contracts()
		if err != nil {
			return err
		}
		selectedSPs = append(selectedSPs, agreement.Meta.SpId)
		as = append(as, agreement)
	}
	// TODO 数据结构要调整
	meta, err := NewFileStatus(as, rss.CtxParams.Cfg, as[0].Meta.CreatorId, rss.Hash, fileSize)
	if err != nil {
		return err
	}
	cb := make(chan []byte)
	uh.FileMetaChanMaps.Set(rss.SsId, cb)
	// 离线签名
	if offlineSigning {
		raw, err := proto.Marshal(meta)
		if err != nil {
			return err
		}
		err = rss.SaveOfflineSigning(&renterpb.OfflineSigning{
			Raw: raw,
		})
		if err != nil {
			return err
		}
	} else {
		go func() {
			if sig, err := func() ([]byte, error) {
				payerPrivKey, err := rss.CtxParams.Cfg.Identity.DecodePrivateKey("")
				if err != nil {
					return nil, err
				}
				sig, err := crypto.Sign(payerPrivKey, meta)
				if err != nil {
					return nil, err
				}
				return sig, nil
			}(); err != nil {
				_ = rss.To(sessions.RssToErrorEvent, err)
				return
			} else {
				cb <- sig
			}
		}()
	}
	signBytes := <-cb
	fmt.Println("doAgreementAndPay, signBytes: ", signBytes)
	uh.FileMetaChanMaps.Remove(rss.SsId)
	if err := rss.To(sessions.RssToGuardFileMetaSignedEvent); err != nil {
		return err
	}
	// fsStatus.RenterSignature = signBytes
	err = chain.SettleObject.FileMetaService.AddFileMeta(rss.Hash, meta)
	if err != nil {
		return err
	}
	// fsStatus, err = submitFileMetaHelper(rss.Ctx, rss.CtxParams.Cfg, fsStatus, signBytes)
	// if err != nil {
	// 	return err
	// }
	// qs, err := guard.PrepFileChallengeQuestions(rss, fsStatus, rss.Hash, offlineSigning, fsStatus.RenterPid)
	// if err != nil {
	// 	return err
	// }
	err = rss.To(sessions.RssToGuardQuestionsSignedEvent)

	// fcid, err := cidlib.Parse(rss.Hash)
	if err != nil {
		return err
	}
	// err = guard.SendChallengeQuestions(rss.Ctx, rss.CtxParams.Cfg, fcid, qs)
	// if err != nil {
	// 	return fmt.Errorf("failed to send challenge questions to guard: [%v]", err)
	// }
	return waitUpload(rss, offlineSigning, meta, false)
}

func NewFileStatus(contracts []*metadata.Agreement, configuration *config.Config,
	renterId string, fileHash string, fileSize int64) (*metadata.FileMetaInfo, error) {

	return &metadata.FileMetaInfo{
		CreatorId:  renterId,
		FileHash:   fileHash,
		FileSize:   uint64(fileSize),
		ShardCount: uint64(len(contracts)),
		Agreements: contracts,
	}, nil
}

func submitFileMetaHelper(ctx context.Context, configuration *config.Config,
	fileStatus *guardpb.FileStoreStatus, sign []byte) (*guardpb.FileStoreStatus, error) {
	if fileStatus.PreparerPid == fileStatus.RenterPid {
		fileStatus.RenterSignature = sign
	} else {
		fileStatus.RenterSignature = sign
		fileStatus.PreparerSignature = sign
	}

	err := submitFileStatus(ctx, configuration, fileStatus)
	if err != nil {
		return nil, err
	}

	return fileStatus, nil
}

func submitFileStatus(ctx context.Context, cfg *config.Config,
	fileStatus *guardpb.FileStoreStatus) error {
	cb := cgrpc.GuardClient(cfg.Services.GuardDomain)
	cb.Timeout(guard.GuardTimeout)
	return cb.WithContext(ctx, func(ctx context.Context, client guardpb.GuardServiceClient) error {
		res, err := client.SubmitFileStoreMeta(ctx, fileStatus)
		if err != nil {
			return err
		}
		if res.Code != guardpb.ResponseCode_SUCCESS {
			return fmt.Errorf("failed to execute submit file status to guard: %v", res.Message)
		}
		return nil
	})
}
