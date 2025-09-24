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

	creatorShardContractKey = "/btfs/%s/creator/shard-contracts/%s"
	userFileShard           = "/btfs/%s/shards/file/%s"

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
		rs.saveShardStatusAndContract(e.Args[0].(*metadata.Contract))
		rs.saveUserShard(e.Args[0].(*metadata.Contract).Meta.ContractId)
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

func (rs *UserShard) saveShardStatusAndContract(signedGuardContract *metadata.Contract) error {
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

func (rs *UserShard) UpdateContractsStatus() error {
	contract := &metadata.Contract{}
	err := Get(rs.ds, fmt.Sprintf(renterShardContractsKey, rs.peerId, GetShardId(rs.ssId, rs.hash, rs.index)), contract)
	if err != nil {
		return err
	}
	contract.Status = metadata.Contract_COMPLETED
	return Save(rs.ds, fmt.Sprintf(renterShardContractsKey, rs.peerId, GetShardId(rs.ssId, rs.hash, rs.index)), contract)
}

func (rs *UserShard) saveUserShard(contractId string) {
	err := rs.ds.Put(context.Background(), datastore.NewKey(fmt.Sprintf(userFileShard, rs.peerId, contractId)), []byte(rs.hash))
	if err != nil {
		return
	}
}

func (rs *UserShard) UpdateShardToContractStatus(signedContract *metadata.Contract) error {
	return rs.fsm.Event(rshToContractEvent, signedContract)
}

func (rs *UserShard) Contracts() (*metadata.Contract, error) {
	contract := &metadata.Contract{}
	err := Get(rs.ds, fmt.Sprintf(renterShardContractsKey, rs.peerId, GetShardId(rs.ssId, rs.hash, rs.index)), contract)
	if errors.Is(err, datastore.ErrNotFound) {
		return contract, nil
	}
	return contract, err
}

func ListShardsContracts(d datastore.Datastore, peerId string, role string) ([]*metadata.Contract, error) {
	var k string
	if k = fmt.Sprintf(renterShardPrefix, peerId); role == nodepb.ContractStat_HOST.String() {
		k = fmt.Sprintf(hostShardPrefix, peerId)
	}
	vs, err := List(d, k, "/contracts")
	if err != nil {
		if err == datastore.ErrNotFound {
			return make([]*metadata.Contract, 0), nil
		}
		return nil, err
	}

	contracts := make([]*metadata.Contract, 0)
	contractMap := make(map[string]*metadata.Contract) // Use map to deduplicate contracts

	for _, v := range vs {
		sc := &metadata.Contract{}
		err := proto.Unmarshal(v, sc)
		if err != nil {
			log.Error("Failed to unmarshal contract:", err)
			continue
		}

		// Validate contract
		if sc.Meta == nil || sc.Meta.ContractId == "" {
			// log.Error("Invalid contract: missing metadata or contract ID")
			continue
		}

		// Check if we already have this contract
		if existing, exists := contractMap[sc.Meta.ContractId]; exists {
			// If existing contract is older, replace it
			if sc.CreateTime > existing.CreateTime {
				contractMap[sc.Meta.ContractId] = sc
			}
		} else {
			contractMap[sc.Meta.ContractId] = sc
		}
	}

	// Convert map back to slice
	for _, contract := range contractMap {
		contracts = append(contracts, contract)
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
func SaveShardsContract(ds datastore.Datastore, scs []*metadata.Contract,
	gcs []*metadata.Contract, peerID, role string) ([]*metadata.Contract, []string, error) {
	var ks []string
	var vs []proto.Message
	gmap := map[string]*metadata.Contract{}
	for _, g := range gcs {
		gmap[g.Meta.ContractId] = g
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
		if gc, ok := gmap[c.Meta.ContractId]; ok {
			ks = append(ks, fmt.Sprintf(key, peerID, c.Meta.ContractId))
			// update
			c = gc
			vs = append(vs, c)
			delete(gmap, c.Meta.ContractId)

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

func UpdateShardContract(ds datastore.Datastore, sc *metadata.Contract, key string) error {
	return Save(ds, key, sc)
}

func GetUserShardContract(ds datastore.Datastore, peerID string, role string, contractID string) (string, *metadata.Contract, error) {
	key := ""
	if role == nodepb.ContractStat_HOST.String() {
		key = hostShardPrefix
	} else {
		key = renterShardPrefix
	}
	values, keys, err := ListWithKeys(ds, fmt.Sprintf(key, peerID), "contract")
	if err != nil {
		return "", nil, err
	}
	for i, v := range values {
		sc := &metadata.Contract{}
		err = proto.Unmarshal(v, sc)
		if err != nil {
			fmt.Println("get contract error", err)
			continue
		}
		if sc.Meta.ContractId == contractID {
			return keys[i], sc, nil
		}
	}

	return "", nil, datastore.ErrNotFound
}

func RefreshLocalContracts(ctx context.Context, ds datastore.Datastore, all []*metadata.Contract, outdated []*metadata.Contract, peerID, role string) ([]string, error) {
	newKeys := make([]string, 0)
	newValues := make([]proto.Message, 0)
	outedFileCIDs := make(map[string]bool)
	updatedContracts := make(map[string]bool)

	key := ""
	if role == nodepb.ContractStat_HOST.String() {
		key = hostShardContractsKey
	} else {
		key = renterShardContractsKey
	}

	// Mark outdated contracts
	outdatedMap := make(map[string]bool)
	for _, o := range outdated {
		if o.Meta != nil {
			outdatedMap[o.Meta.ContractId] = true
		}
	}

	// Process all contracts, excluding outdated ones
	for _, a := range all {
		if a.Meta == nil || a.Meta.ContractId == "" {
			continue
		}

		// Skip if contract is outdated
		if outdatedMap[a.Meta.ContractId] {
			continue
		}

		// Add to new contracts list
		newKeys = append(newKeys, fmt.Sprintf(key, peerID, a.Meta.ContractId))
		newValues = append(newValues, a)
		updatedContracts[a.Meta.ContractId] = true
	}

	// Get file CIDs for outdated contracts
	for _, o := range outdated {
		if o.Meta == nil {
			continue
		}
		cid, err := ds.Get(ctx, datastore.NewKey(fmt.Sprintf(userFileShard, peerID, o.Meta.ContractId)))
		if err != nil {
			if errors.Is(err, datastore.ErrNotFound) {
				continue
			}
			log.Error("Failed to get CID for outdated contract:", err)
			continue
		}
		outedFileCIDs[string(cid)] = true
	}

	// Save updated contracts
	if len(newKeys) > 0 {
		err := Batch(ds, newKeys, newValues)
		if err != nil {
			return nil, fmt.Errorf("failed to save updated contracts: %v", err)
		}
	}

	// Return list of stale CIDs
	staled := make([]string, 0)
	for k := range outedFileCIDs {
		staled = append(staled, k)
	}

	return staled, nil
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
