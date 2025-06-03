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
	if err := rss.To(sessions.RssToAgreementEvent); err != nil {
		return err
	}

	as := make([]*metadata.Agreement, 0)
	for i, h := range rss.ShardHashes {
		shard, err := sessions.GetRenterShard(rss.CtxParams, rss.SsId, h, i)
		if err != nil {
			return err
		}
		agreement, err := shard.Contracts()
		if err != nil {
			return err
		}
		as = append(as, agreement)
	}
	// TODO 数据结构要调整
	meta, err := NewFileStatus(as, rss.CtxParams.Cfg, as[0].Meta.CreatorId, rss.Hash, fileSize)
	if err != nil {
		return err
	}
	cb := make(chan []byte)
	uh.FileMetaChanMaps.Set(rss.SsId, cb)
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
	<-cb
	uh.FileMetaChanMaps.Remove(rss.SsId)
	if err := rss.To(sessions.RssToGuardFileMetaSignedEvent); err != nil {
		return err
	}
	// fsStatus.RenterSignature = signBytes
	err = chain.SettleObject.FileMetaService.AddFileMeta(rss.Hash, meta)
	if err != nil {
		return err
	}
	err = rss.To(sessions.RssToGuardQuestionsSignedEvent)
	if err != nil {
		return err
	}

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
