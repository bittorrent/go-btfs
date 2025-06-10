package upload

import (
	"github.com/bittorrent/go-btfs/chain"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/protos/metadata"
	renterpb "github.com/bittorrent/go-btfs/protos/renter"

	"github.com/bittorrent/go-btfs-common/crypto"
	config "github.com/bittorrent/go-btfs-config"

	"github.com/gogo/protobuf/proto"
)

func addFileMetaToBttcChainAndPay(rss *sessions.RenterSession, fileSize int64, offlineSigning bool) error {
	if err := rss.To(sessions.RssToContractEvent); err != nil {
		return err
	}

	as := make([]*metadata.Contract, 0)
	for i, h := range rss.ShardHashes {
		shard, err := sessions.GetUserShard(rss.CtxParams, rss.SsId, h, i)
		if err != nil {
			return err
		}
		contract, err := shard.Contracts()
		if err != nil {
			return err
		}
		as = append(as, contract)
	}
	meta, err := NewFileStatus(as, rss.CtxParams.Cfg, as[0].Meta.UserId, rss.Hash, fileSize)
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
	if err := rss.To(sessions.RssToContractFileMetaSignedEvent); err != nil {
		return err
	}
	// fsStatus.RenterSignature = signBytes
	err = chain.SettleObject.FileMetaService.AddFileMeta(rss.Hash, meta)
	if err != nil {
		return err
	}
	err = rss.To(sessions.RssToContractFileMetaAddedEvent)
	if err != nil {
		return err
	}

	return waitSPSaveFileSuccAndToPay(rss, offlineSigning, meta, false)
}

func NewFileStatus(contracts []*metadata.Contract, configuration *config.Config,
	userId string, fileHash string, fileSize int64) (*metadata.FileMetaInfo, error) {

	return &metadata.FileMetaInfo{
		UserId:     userId,
		FileHash:   fileHash,
		FileSize:   uint64(fileSize),
		ShardCount: uint64(len(contracts)),
		Contracts:  contracts,
	}, nil
}
