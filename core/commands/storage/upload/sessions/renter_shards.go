package sessions

import (
	"context"
	"fmt"
	"strings"
	"time"

	renterpb "github.com/bittorrent/go-btfs/protos/renter"

	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	shardpb "github.com/bittorrent/go-btfs/protos/shard"

	guardpb "github.com/bittorrent/go-btfs-common/protos/guard"
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

	rshInitStatus      = "init"
	rshContractStatus  = "contract"
	rshErrorStatus     = "error"
	rshToContractEvent = "to-contract"
)

var log = logging.Logger("sessions")

var (
	renterShardFsmEvents = fsm.Events{
		{Name: rshToContractEvent, Src: []string{rshInitStatus}, Dst: rshContractStatus},
	}
	renterShardsInMem = cmap.New()
)

type RenterShard struct {
	peerId string
	ssId   string
	hash   string
	index  int
	fsm    *fsm.FSM
	ctx    context.Context
	ds     datastore.Datastore
}

func GetRenterShard(ctxParams *uh.ContextParams, ssId string, hash string, index int) (*RenterShard, error) {
	shardId := GetShardId(ssId, hash, index)
	k := fmt.Sprintf(renterShardsInMemKey, ctxParams.N.Identity.Pretty(), shardId)
	var rs *RenterShard
	if tmp, ok := renterShardsInMem.Get(k); ok {
		rs = tmp.(*RenterShard)
	} else {
		ctx, _ := helper.NewGoContext(ctxParams.Ctx)
		rs = &RenterShard{
			peerId: ctxParams.N.Identity.Pretty(),
			ssId:   ssId,
			hash:   hash,
			index:  index,
			ctx:    ctx,
			ds:     ctxParams.N.Repo.Datastore(),
		}
		renterShardsInMem.Set(k, rs)
	}
	status, err := rs.Status()
	if err != nil {
		return nil, err
	}
	if rs.fsm == nil && status.Status == rshInitStatus {
		rs.fsm = fsm.NewFSM(status.Status, renterShardFsmEvents, fsm.Callbacks{
			"enter_state": rs.enterState,
		})
	}
	return rs, nil
}

func (rs *RenterShard) enterState(e *fsm.Event) {
	log.Infof("shard: %s:%s enter status: %s", rs.ssId, rs.hash, e.Dst)
	switch e.Dst {
	case rshContractStatus:
		rs.doContract(e.Args[0].([]byte), e.Args[1].(*guardpb.Contract))
	}
}

func (rs *RenterShard) Status() (*shardpb.Status, error) {
	status := new(shardpb.Status)
	shardId := GetShardId(rs.ssId, rs.hash, rs.index)
	k := fmt.Sprintf(renterShardStatusKey, rs.peerId, shardId)
	err := Get(rs.ds, k, status)
	if err == datastore.ErrNotFound {
		status = &shardpb.Status{
			Status: rshInitStatus,
		}
		//ignore error
		_ = Save(rs.ds, k, status)
	} else if err != nil {
		return nil, err
	}
	return status, nil
}

func GetShardId(ssId string, shardHash string, index int) (contractId string) {
	return fmt.Sprintf("%s:%s:%d", ssId, shardHash, index)
}

// ExtractSessionIDFromContractID takes the first segment separated by ","
// and returns it as the session id. If in an old format, i.e. did not
// have a session id, return an error.
func extractSessionIDFromContractID(contractID string) (string, error) {
	ids := strings.Split(contractID, ":")
	if len(ids) != 2 {
		return "", fmt.Errorf("bad contract id: fewer than 2 segments")
	}
	if len(ids[0]) != 36 {
		return "", fmt.Errorf("invalid session id within contract id")
	}
	return ids[0], nil
}

func (rs *RenterShard) doContract(signedEscrowContract []byte, signedGuardContract *guardpb.Contract) error {
	status := &shardpb.Status{
		Status: rshContractStatus,
	}
	signedContracts := &shardpb.SignedContracts{
		SignedEscrowContract: signedEscrowContract,
		SignedGuardContract:  signedGuardContract,
	}
	shardId := GetShardId(rs.ssId, rs.hash, rs.index)
	return Batch(rs.ds, []string{
		fmt.Sprintf(renterShardStatusKey, rs.peerId, shardId),
		fmt.Sprintf(renterShardContractsKey, rs.peerId, shardId),
	}, []proto.Message{
		status, signedContracts,
	})
}

func (rs *RenterShard) Contract(signedEscrowContract []byte, signedGuardContract *guardpb.Contract) error {
	return rs.fsm.Event(rshToContractEvent, signedEscrowContract, signedGuardContract)
}

func (rs *RenterShard) Contracts() (*shardpb.SignedContracts, error) {
	contracts := &shardpb.SignedContracts{}
	err := Get(rs.ds, fmt.Sprintf(renterShardContractsKey, rs.peerId, GetShardId(rs.ssId, rs.hash, rs.index)), contracts)
	if err == datastore.ErrNotFound {
		return contracts, nil
	}
	return contracts, err
}

func ListShardsContracts(d datastore.Datastore, peerId string, role string) ([]*shardpb.SignedContracts, error) {
	var k string
	if k = fmt.Sprintf(renterShardPrefix, peerId); role == nodepb.ContractStat_HOST.String() {
		k = fmt.Sprintf(hostShardPrefix, peerId)
	}
	vs, err := List(d, k, "/contracts")
	if err != nil {
		return nil, err
	}
	contracts := make([]*shardpb.SignedContracts, 0)
	for _, v := range vs {
		sc := &shardpb.SignedContracts{}
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

// SaveShardsContracts persists updated guard contracts from upstream, if an existing entry
// is not available, then an empty signed escrow contract is inserted along with the
// new guard contract.
func SaveShardsContracts(ds datastore.Datastore, scs []*shardpb.SignedContracts,
	gcs []*guardpb.Contract, peerID, role string) ([]*shardpb.SignedContracts, []string, error) {
	var ks []string
	var vs []proto.Message
	gmap := map[string]*guardpb.Contract{}
	for _, g := range gcs {
		gmap[g.ContractId] = g
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
		if gc, ok := gmap[c.SignedGuardContract.ContractId]; ok {
			ks = append(ks, fmt.Sprintf(key, peerID, c.SignedGuardContract.ContractId))
			// update
			c.SignedGuardContract = gc
			vs = append(vs, c)
			delete(gmap, c.SignedGuardContract.ContractId)

			// mark stale files if no longer active (must be synced to become inactive)
			if _, ok := helper.ContractFilterMap["active"][gc.State]; !ok {
				invalidShards[gc.ShardHash] = append(invalidShards[gc.ShardHash], gc.FileHash)
			}
		} else {
			activeShards[c.SignedGuardContract.ShardHash] = true
			activeFiles[c.SignedGuardContract.FileHash] = true
		}
	}
	// append what's left in guard map as new contracts
	for contractID, gc := range gmap {
		ks = append(ks, fmt.Sprintf(key, peerID, contractID))
		// add a new (guard contract only) signed contracts
		c := &shardpb.SignedContracts{SignedGuardContract: gc}
		scs = append(scs, c)
		vs = append(vs, c)

		// mark stale files if no longer active (must be synced to become inactive)
		if _, ok := helper.ContractFilterMap["active"][gc.State]; !ok {
			invalidShards[gc.ShardHash] = append(invalidShards[gc.ShardHash], gc.FileHash)
		} else {
			activeShards[gc.ShardHash] = true
			activeFiles[gc.FileHash] = true
		}
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
		//staleHashes = append(staleHashes, ish)
	}
	return scs, staleHashes, nil
}

func (rs *RenterShard) UpdateAdditionalInfo(info string) error {
	shardId := GetShardId(rs.ssId, rs.hash, rs.index)
	return Save(rs.ds, fmt.Sprintf(renterShardAdditionalInfoKey, rs.peerId, shardId),
		&renterpb.RenterSessionAdditionalInfo{
			Info:        info,
			LastUpdated: time.Now(),
		})
}

func (rs *RenterShard) GetAdditionalInfo() (*shardpb.AdditionalInfo, error) {
	pb := &shardpb.AdditionalInfo{}
	shardId := GetShardId(rs.ssId, rs.hash, rs.index)
	err := Get(rs.ds, fmt.Sprintf(renterShardAdditionalInfoKey, rs.peerId, shardId), pb)
	return pb, err
}
