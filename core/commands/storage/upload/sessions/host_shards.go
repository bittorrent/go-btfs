package sessions

import (
	"context"
	"fmt"
	"math/big"

	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	shardpb "github.com/bittorrent/go-btfs/protos/shard"

	guardpb "github.com/bittorrent/go-btfs-common/protos/guard"
	"github.com/bittorrent/protobuf/proto"

	"github.com/ipfs/go-datastore"
	"github.com/looplab/fsm"
	cmap "github.com/orcaman/concurrent-map"
)

const (
	hostShardPrefix       = "/btfs/%s/host/shards/"
	hostShardKey          = hostShardPrefix + "%s/"
	hostShardsInMemKey    = hostShardKey
	hostShardStatusKey    = hostShardKey + "status"
	hostShardContractsKey = hostShardKey + "contracts"

	hshInitStatus     = "init"
	hshContractStatus = "contract"
	hshPayStatus      = "paid"
	hshCompleteStatus = "complete"
	hshErrorStatus    = "error"

	hshToContractEvent = "to-contract"
	hshToPayEvent      = "to-pay"
	hshToCompleteEvent = "to-complete"
	hshToErrorEvent    = "to-error"
)

var (
	hostShardFsmEvents = fsm.Events{
		{Name: hshToContractEvent, Src: []string{hshInitStatus}, Dst: hshContractStatus},
		{Name: hshToPayEvent, Src: []string{hshContractStatus}, Dst: hshPayStatus},
		{Name: hshToCompleteEvent, Src: []string{hshPayStatus}, Dst: hshCompleteStatus},
		{Name: hshToErrorEvent, Src: []string{hshInitStatus, hshContractStatus}, Dst: hshToErrorEvent},
	}
	hostShardsInMem = cmap.New()
)

type HostShard struct {
	peerId     string
	contractId string
	fsm        *fsm.FSM
	ctx        context.Context
	ds         datastore.Datastore
	inputPrice int64
	amount     int64
	rate       *big.Int
}

func GetHostShard(ctxParams *uh.ContextParams, contractId string, inputPrice int64, amount int64, rate *big.Int) (*HostShard, error) {
	k := fmt.Sprintf(hostShardsInMemKey, ctxParams.N.Identity.Pretty(), contractId)
	var hs *HostShard
	if tmp, ok := hostShardsInMem.Get(k); ok {
		hs = tmp.(*HostShard)
	} else {
		ctx, _ := helper.NewGoContext(ctxParams.Ctx)
		hs = &HostShard{
			peerId:     ctxParams.N.Identity.Pretty(),
			contractId: contractId,
			ctx:        ctx,
			ds:         ctxParams.N.Repo.Datastore(),
			inputPrice: inputPrice,
			amount:     amount,
			rate:       rate,
		}
		hostShardsInMem.Set(k, hs)
	}
	status, err := hs.status()
	if err != nil {
		return nil, err
	}
	if hs.fsm == nil && status.Status == hshInitStatus {
		hs.fsm = fsm.NewFSM(status.Status, hostShardFsmEvents, fsm.Callbacks{
			"enter_state": hs.enterState,
		})
	}
	return hs, nil
}

func (hs *HostShard) GetInputPrice() int64 {
	return hs.inputPrice
}
func (hs *HostShard) GetInputAmount() int64 {
	return hs.amount
}
func (hs *HostShard) GetInputRate() *big.Int {
	return hs.rate
}

func (hs *HostShard) enterState(e *fsm.Event) {
	log.Infof("shard: %s enter status: %s\n", hs.contractId, e.Dst)
	switch e.Dst {
	case hshContractStatus:
		hs.doContract(e.Args[0].([]byte), e.Args[1].(*guardpb.Contract))
	}
}

func (hs *HostShard) status() (*shardpb.Status, error) {
	status := new(shardpb.Status)
	k := fmt.Sprintf(hostShardStatusKey, hs.peerId, hs.contractId)
	err := Get(hs.ds, k, status)
	if err == datastore.ErrNotFound {
		status = &shardpb.Status{
			Status: hshInitStatus,
		}
		//ignore error
		_ = Save(hs.ds, k, status)
	} else if err != nil {
		return nil, err
	}
	return status, nil
}

func (hs *HostShard) doContract(signedEscrowContract []byte, signedGuardContract *guardpb.Contract) error {
	status := &shardpb.Status{
		Status: hshContractStatus,
	}
	signedContracts := &shardpb.SignedContracts{
		SignedEscrowContract: signedEscrowContract,
		SignedGuardContract:  signedGuardContract,
	}
	return Batch(hs.ds, []string{
		fmt.Sprintf(hostShardStatusKey, hs.peerId, hs.contractId),
		fmt.Sprintf(hostShardContractsKey, hs.peerId, hs.contractId),
	}, []proto.Message{
		status, signedContracts,
	})
}

func (hs *HostShard) Contract(signedEscrowContract []byte, signedGuardContract *guardpb.Contract) error {
	return hs.fsm.Event(hshToContractEvent, signedEscrowContract, signedGuardContract)
}

func (hs *HostShard) IsPayStatus() bool {
	fmt.Printf("IsPayStatus Current:%v,  hshPayStatus:%v \n", hs.fsm.Current(), hshPayStatus)
	return hs.fsm.Current() == hshPayStatus
}
func (hs *HostShard) IsContractStatus() bool {
	fmt.Printf("IsContractStatus Current:%v,  hshContractStatus:%v \n", hs.fsm.Current(), hshContractStatus)
	return hs.fsm.Current() == hshContractStatus
}

func (hs *HostShard) ReceivePayCheque() error {
	fmt.Printf("ReceivePayCheque cur=%+v \n", hs.fsm.Current())
	return hs.fsm.Event(hshToPayEvent)
}

func (hs *HostShard) Complete() error {
	return hs.fsm.Event(hshToCompleteEvent)
}

func (hs *HostShard) contracts() (*shardpb.SignedContracts, error) {
	contracts := &shardpb.SignedContracts{}
	err := Get(hs.ds, fmt.Sprintf(hostShardContractsKey, hs.peerId, hs.contractId), contracts)
	if err == datastore.ErrNotFound {
		return contracts, nil
	}
	return contracts, err
}
