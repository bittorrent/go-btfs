package sessions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/protos/metadata"
	renterpb "github.com/bittorrent/go-btfs/protos/renter"
	shardpb "github.com/bittorrent/go-btfs/protos/shard"

	nodepb "github.com/bittorrent/go-btfs-common/protos/node"
	"github.com/bittorrent/protobuf/proto"

	"github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log"
	"github.com/looplab/fsm"
	cmap "github.com/orcaman/concurrent-map"
)

const (
	renterShardPrefix            = "/btfs/%s/renter/shards/"
	renterShardKey               = renterShardPrefix + "%s/"
	renterShardsInMemKey         = renterShardKey
	renterShardStatusKey         = renterShardKey + "status"
	renterShardContractsKey      = renterShardKey + "contracts"
	renterShardAdditionalInfoKey = renterShardKey + "additional-info"

	creatorShardAgreementKey = "/btfs/%s/creator/shard-agreements/%s"
	userFileShard            = "/btfs/%s/shards/file/%s"

	// status
	rshInitStatus     = "init"
	rshContractStatus = "contract"
	rshErrorStatus    = "error"

	// event
	rshToContractEvent = "to-contract"
)

var log = logging.Logger("sessions")

var (
	renterShardFsmEvents = fsm.Events{
		{Name: rshToContractEvent, Src: []string{rshInitStatus}, Dst: rshContractStatus},
	}
	renterShardsInMem = cmap.New()
)

type UserShard struct {
	peerId string
	ssId   string
	hash   string
	index  int
	fsm    *fsm.FSM
	ctx    context.Context
	ds     datastore.Datastore
}

func GetUserShard(ctxParams *uh.ContextParams, ssId string, hash string, index int) (*UserShard, error) {
	shardId := GetShardId(ssId, hash, index)
	k := fmt.Sprintf(renterShardsInMemKey, ctxParams.N.Identity.String(), shardId)
	var us *UserShard
	if tmp, ok := renterShardsInMem.Get(k); ok {
		us = tmp.(*UserShard)
	} else {
		ctx, _ := helper.NewGoContext(ctxParams.Ctx)
		us = &UserShard{
			peerId: ctxParams.N.Identity.String(),
			ssId:   ssId,
			hash:   hash,
			index:  index,
			ctx:    ctx,
			ds:     ctxParams.N.Repo.Datastore(),
		}
		renterShardsInMem.Set(k, us)
	}
	status, err := us.GetShardStatus()
	if err != nil {
		return nil, err
	}
	if us.fsm == nil && status.Status == rshInitStatus {
		us.fsm = fsm.NewFSM(status.Status, renterShardFsmEvents, fsm.Callbacks{
			"enter_state": us.enterState,
		})
	}
	return us, nil
}

func (rs *UserShard) enterState(e *fsm.Event) {
	log.Infof("shard: %s:%s enter status: %s", rs.ssId, rs.hash, e.Dst)
	switch e.Dst {
	case rshContractStatus:
		rs.saveShardStatusAndContract(e.Args[0].(*metadata.Agreement))
		rs.saveUserShard(e.Args[0].(*metadata.Agreement).Meta.AgreementId)
	}
}

func (rs *UserShard) GetShardStatus() (*shardpb.Status, error) {
	status := new(shardpb.Status)
	shardId := GetShardId(rs.ssId, rs.hash, rs.index)
	k := fmt.Sprintf(renterShardStatusKey, rs.peerId, shardId)
	err := Get(rs.ds, k, status)
	if err == datastore.ErrNotFound {
		status = &shardpb.Status{
			Status: rshInitStatus,
		}
		// ignore error
		_ = Save(rs.ds, k, status)
	} else if err != nil {
		return nil, err
	}
	return status, nil
}

func GetShardId(ssId string, shardHash string, index int) (contractId string) {
	return fmt.Sprintf("%s:%s:%d", ssId, shardHash, index)
}

func (rs *UserShard) saveShardStatusAndContract(signedGuardContract *metadata.Agreement) error {
	status := &shardpb.Status{
		Status: rshContractStatus,
	}
	shardId := GetShardId(rs.ssId, rs.hash, rs.index)
	return Batch(rs.ds, []string{
		fmt.Sprintf(renterShardStatusKey, rs.peerId, shardId),
		fmt.Sprintf(renterShardContractsKey, rs.peerId, shardId),
	}, []proto.Message{
		status,
		signedGuardContract,
	})
}

func (rs *UserShard) saveUserShard(contractId string) {
	err := rs.ds.Put(context.Background(), datastore.NewKey(fmt.Sprintf(userFileShard, rs.peerId, contractId)), []byte(rs.hash))
	if err != nil {
		return
	}
}

func (rs *UserShard) UpdateShardToContractStatus(signedAgreement *metadata.Agreement) error {
	return rs.fsm.Event(rshToContractEvent, signedAgreement)
}

func (rs *UserShard) Agreements() (*metadata.Agreement, error) {
	agreement := &metadata.Agreement{}
	err := Get(rs.ds, fmt.Sprintf(renterShardContractsKey, rs.peerId, GetShardId(rs.ssId, rs.hash, rs.index)), agreement)
	if errors.Is(err, datastore.ErrNotFound) {
		return agreement, nil
	}
	return agreement, err
}

func ListShardsContracts(d datastore.Datastore, peerId string, role string) ([]*metadata.Agreement, error) {
	var k string
	if k = fmt.Sprintf(renterShardPrefix, peerId); role == nodepb.ContractStat_HOST.String() {
		k = fmt.Sprintf(hostShardPrefix, peerId)
	}
	vs, err := List(d, k, "/contracts")
	if err != nil {
		return nil, err
	}
	contracts := make([]*metadata.Agreement, 0)
	for _, v := range vs {
		sc := &metadata.Agreement{}
		err := proto.Unmarshal(v, sc)
		if err != nil {
			log.Error(err)
			continue
		}
		contracts = append(contracts, sc)
	}
	return contracts, nil
}

func DeleteShardsContracts(d datastore.Datastore, peerId string, role string) error {
	var k string
	if k = fmt.Sprintf(renterShardPrefix, peerId); role == nodepb.ContractStat_HOST.String() {
		k = fmt.Sprintf(hostShardPrefix, peerId)
	}
	ks, err := ListKeys(d, k, "/contracts")
	if err != nil {
		return err
	}
	vs := make([]proto.Message, len(ks))
	for range ks {
		vs = append(vs, nil)
	}
	return Batch(d, ks, vs)
}

// SaveShardsContract persists updated guard contracts from upstream, if an existing entry
// is not available, then an empty signed escrow contract is inserted along with the
// new guard contract.
func SaveShardsContract(ds datastore.Datastore, scs []*metadata.Agreement,
	gcs []*metadata.Agreement, peerID, role string) ([]*metadata.Agreement, []string, error) {
	var ks []string
	var vs []proto.Message
	gmap := map[string]*metadata.Agreement{}
	for _, g := range gcs {
		gmap[g.Meta.AgreementId] = g
	}
	activeShards := map[string]bool{}      // active shard hash -> has one file hash (bool)
	activeFiles := map[string]bool{}       // active file hash -> has one shard hash (bool)
	invalidShards := map[string][]string{} // invalid shard hash -> (maybe) invalid file hash list
	var key string
	if role == nodepb.ContractStat_HOST.String() {
		key = hostShardContractsKey
	} else {
		key = renterShardContractsKey
	}
	for _, c := range scs {
		// only append the updated contracts
		if gc, ok := gmap[c.Meta.AgreementId]; ok {
			ks = append(ks, fmt.Sprintf(key, peerID, c.Meta.AgreementId))
			// update
			c = gc
			vs = append(vs, c)
			delete(gmap, c.Meta.AgreementId)

			// mark stale files if no longer active (must be synced to become inactive)
			invalidShards[gc.Meta.ShardHash] = append(invalidShards[gc.Meta.ShardHash], gc.Meta.ShardHash)
		} else {
			activeShards[c.Meta.ShardHash] = true
			activeFiles[c.Meta.ShardHash] = true
		}
	}
	// append what's left in guard map as new contracts
	for contractID, gc := range gmap {
		ks = append(ks, fmt.Sprintf(key, peerID, contractID))
		// add a new (guard contract only) signed contracts
		scs = append(scs, gc)
		vs = append(vs, gc)

		// mark stale files if no longer active (must be synced to become inactive)
		activeShards[gc.Meta.ShardHash] = true
		activeFiles[gc.Meta.ShardHash] = true
	}
	if len(ks) > 0 {
		err := Batch(ds, ks, vs)
		if err != nil {
			return nil, nil, err
		}
	}
	var staleHashes []string
	// compute what's stale
	for ish, fhs := range invalidShards {
		if _, ok := activeShards[ish]; ok {
			// other files are referring to this hash, skip
			continue
		}
		for _, fh := range fhs {
			if _, ok := activeFiles[fh]; !ok {
				// file does not have other active shards
				staleHashes = append(staleHashes, fh)
			}
		}
		// TODO: Cannot prematurally remove shard because it's indirectly pinned
		// Need a way to disassociated indirect pins from parent...
		// remove hash anyway even if no file is getting removed
		// staleHashes = append(staleHashes, ish)
	}
	return scs, staleHashes, nil
}

func RefreshLocalContracts(ctx context.Context, ds datastore.Datastore, all []*metadata.Agreement, outdated []*metadata.Agreement, peerID, role string) ([]string, error) {
	newKeys := make([]string, 0)
	newValues := make([]proto.Message, 0)
	outedFileCIDs := make(map[string]bool)

	key := ""
	if role == nodepb.ContractStat_HOST.String() {
		key = hostShardContractsKey
	} else {
		key = renterShardContractsKey
	}

	for _, o := range outdated {
		for _, a := range all {
			if a.Meta.AgreementId == o.Meta.AgreementId {
				continue
			}
			newKeys = append(newKeys, fmt.Sprintf(key, peerID, a.Meta.AgreementId))
			newValues = append(newValues, a)
		}
	}

	for _, o := range outdated {
		cid, err := ds.Get(ctx, datastore.NewKey(fmt.Sprintf(userFileShard, peerID, o.Meta.AgreementId)))
		if err != nil {
			log.Error(err)
			continue
		}
		outedFileCIDs[string(cid)] = true
	}

	err := Batch(ds, newKeys, newValues)
	staled := make([]string, 0)
	for k := range outedFileCIDs {
		staled = append(staled, k)
	}

	return staled, err
}

func (rs *UserShard) UpdateAdditionalInfo(info string) error {
	shardId := GetShardId(rs.ssId, rs.hash, rs.index)
	return Save(rs.ds, fmt.Sprintf(renterShardAdditionalInfoKey, rs.peerId, shardId),
		&renterpb.RenterSessionAdditionalInfo{
			Info:        info,
			LastUpdated: time.Now(),
		})
}

func (rs *UserShard) GetAdditionalInfo() (*shardpb.AdditionalInfo, error) {
	pb := &shardpb.AdditionalInfo{}
	shardId := GetShardId(rs.ssId, rs.hash, rs.index)
	err := Get(rs.ds, fmt.Sprintf(renterShardAdditionalInfoKey, rs.peerId, shardId), pb)
	return pb, err
}
