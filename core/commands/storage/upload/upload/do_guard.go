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
	if err := rss.To(sessions.RssToAgreementEvent); err != nil {
		return err
	}

	as := make([]*metadata.Agreement, 0)
	for i, h := range rss.ShardHashes {
		shard, err := sessions.GetUserShard(rss.CtxParams, rss.SsId, h, i)
		if err != nil {
			return err
		}
		agreement, err := shard.Agreements()
		if err != nil {
			return err
		}
		as = append(as, agreement)
	}
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

	return waitSPSaveFileSuccAndToPay(rss, offlineSigning, meta, false)
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
