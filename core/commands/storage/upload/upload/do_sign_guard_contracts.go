package upload

import (
	"fmt"
	"time"

	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"

	"github.com/bittorrent/go-btfs-common/crypto"
	guardpb "github.com/bittorrent/go-btfs-common/protos/guard"
	config "github.com/bittorrent/go-btfs-config"
	"github.com/bittorrent/protobuf/proto"

	"github.com/libp2p/go-libp2p/core/peer"
)

type ContractParams struct {
	ContractId    string
	RenterPid     string
	HostPid       string
	ShardIndex    int32
	ShardHash     string
	ShardSize     int64
	FileHash      string
	StartTime     time.Time
	StorageLength int64
	Price         int64
	TotalPay      int64
}

type RepairParams struct {
	RenterStart time.Time
	RenterEnd   time.Time
}

func RenterSignGuardContract(rss *sessions.RenterSession, params *ContractParams, offlineSigning bool,
	rp *RepairParams, token string) ([]byte,
	error) {
	guardPid, escrowPid, err := getGuardAndEscrowPid(rss.CtxParams.Cfg)
	if err != nil {
		return nil, err
	}
	gm := &guardpb.ContractMeta{
		ContractId:    params.ContractId,
		RenterPid:     params.RenterPid,
		HostPid:       params.HostPid,
		ShardHash:     params.ShardHash,
		ShardIndex:    params.ShardIndex,
		ShardFileSize: params.ShardSize,
		FileHash:      params.FileHash,
		RentStart:     params.StartTime,
		RentEnd:       params.StartTime.Add(time.Duration(params.StorageLength*24) * time.Hour),
		GuardPid:      guardPid.Pretty(),
		EscrowPid:     escrowPid.Pretty(),
		Price:         params.Price,
		Amount:        params.TotalPay,
	}
	cont := &guardpb.Contract{
		ContractMeta:   *gm,
		LastModifyTime: time.Now(),
	}
	if rp != nil {
		cont.State = guardpb.Contract_RENEWED
		cont.RentStart = rp.RenterStart
		cont.RentEnd = rp.RenterEnd
	}
	cont.RenterPid = params.RenterPid
	cont.PreparerPid = params.RenterPid
	bc := make(chan []byte)
	shardId := sessions.GetShardId(rss.SsId, gm.ShardHash, int(gm.ShardIndex))
	uh.GuardChanMaps.Set(shardId, bc)
	bytes, err := proto.Marshal(gm)
	if err != nil {
		return nil, err
	}
	uh.GuardContractMaps.Set(shardId, bytes)
	if !offlineSigning {
		go func() {
			sign, err := crypto.Sign(rss.CtxParams.N.PrivateKey, gm)
			if err != nil {
				_ = rss.To(sessions.RssToErrorEvent, err)
				return
			}
			bc <- sign
		}()
	}
	signedBytes := <-bc
	uh.GuardChanMaps.Remove(shardId)
	uh.GuardContractMaps.Remove(shardId)
	cont.RenterSignature = signedBytes
	cont.Token = token
	return proto.Marshal(cont)
}

func getGuardAndEscrowPid(configuration *config.Config) (peer.ID, peer.ID, error) {
	escrowPubKeys := configuration.Services.EscrowPubKeys
	if len(escrowPubKeys) == 0 {
		return "", "", fmt.Errorf("missing escrow public key in config")
	}
	guardPubKeys := configuration.Services.GuardPubKeys
	if len(guardPubKeys) == 0 {
		return "", "", fmt.Errorf("missing guard public key in config")
	}
	escrowPid, err := helper.PidFromString(escrowPubKeys[0])
	if err != nil {
		log.Error("parse escrow config failed", escrowPubKeys[0])
		return "", "", err
	}
	guardPid, err := helper.PidFromString(guardPubKeys[0])
	if err != nil {
		log.Error("parse guard config failed", guardPubKeys[1])
		return "", "", err
	}
	return guardPid, escrowPid, err
}
