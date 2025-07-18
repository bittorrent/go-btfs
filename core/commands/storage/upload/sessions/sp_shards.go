package sessions

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/protos/metadata"
	shardpb "github.com/bittorrent/go-btfs/protos/shard"

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

	// status
	hshInitStatus     = "init"
	hshContractStatus = "contract"
	hshPayStatus      = "paid"
	hshCompleteStatus = "complete"
	hshErrorStatus    = "error"

	// event
	hshToContractEvent = "to-contract"
	hshToPayEvent      = "to-pay"
	hshToCompleteEvent = "to-complete"
	hshToErrorEvent    = "to-error"
)

var (
	spShardFsmEvents = fsm.Events{
		{Name: hshToContractEvent, Src: []string{hshInitStatus}, Dst: hshContractStatus},
		{Name: hshToPayEvent, Src: []string{hshContractStatus}, Dst: hshPayStatus},
		{Name: hshToCompleteEvent, Src: []string{hshPayStatus}, Dst: hshCompleteStatus},
		{Name: hshToErrorEvent, Src: []string{hshInitStatus, hshContractStatus}, Dst: hshToErrorEvent},
	}
	hostShardsInMem = cmap.New()
)

type SPShard struct {
	peerId     string
	contractId string
	fsm        *fsm.FSM
	ctx        context.Context
	ds         datastore.Datastore
	inputPrice int64
	amount     int64
	rate       *big.Int
}

func GetSPShard(ctxParams *uh.ContextParams, contractId string, inputPrice int64, amount int64, rate *big.Int) (*SPShard, error) {
	k := fmt.Sprintf(hostShardsInMemKey, ctxParams.N.Identity.String(), contractId)
	var hs *SPShard
	if tmp, ok := hostShardsInMem.Get(k); ok {
		hs = tmp.(*SPShard)
	} else {
		ctx, _ := helper.NewGoContext(ctxParams.Ctx)
		hs = &SPShard{
			peerId:     ctxParams.N.Identity.String(),
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
		hs.fsm = fsm.NewFSM(status.Status, spShardFsmEvents, fsm.Callbacks{
			"enter_state": hs.enterState,
		})
	}
	return hs, nil
}

func (hs *SPShard) GetInputPrice() int64 {
	return hs.inputPrice
}
func (hs *SPShard) GetInputAmount() int64 {
	return hs.amount
}
func (hs *SPShard) GetInputRate() *big.Int {
	return hs.rate
}

func (hs *SPShard) enterState(e *fsm.Event) {
	log.Infof("shard: %s enter status: %s\n", hs.contractId, e.Dst)
	switch e.Dst {
	case hshContractStatus:
		hs.saveContract(e.Args[0].(*metadata.Contract))
		hs.saveUserShard(e.Args[0].(*metadata.Contract).Meta.ContractId, e.Args[0].(*metadata.Contract).Meta.ShardHash)
	}
}

func (hs *SPShard) saveContract(signedGuardContract *metadata.Contract) error {
	status := &shardpb.Status{
		Status: hshContractStatus,
	}
	return Batch(hs.ds, []string{
		fmt.Sprintf(hostShardStatusKey, hs.peerId, hs.contractId),
		fmt.Sprintf(hostShardContractsKey, hs.peerId, hs.contractId),
	}, []proto.Message{
		status, signedGuardContract,
	})
}

func (hs *SPShard) saveUserShard(contractId string, shardHash string) {
	err := hs.ds.Put(context.Background(), datastore.NewKey(fmt.Sprintf(userFileShard, hs.peerId, contractId)), []byte(shardHash))
	if err != nil {
		return
	}
}

func (hs *SPShard) status() (*shardpb.Status, error) {
	status := new(shardpb.Status)
	k := fmt.Sprintf(hostShardStatusKey, hs.peerId, hs.contractId)
	err := Get(hs.ds, k, status)
	if errors.Is(err, datastore.ErrNotFound) {
		status = &shardpb.Status{
			Status: hshInitStatus,
		}
		// ignore error
		_ = Save(hs.ds, k, status)
	} else if err != nil {
		return nil, err
	}
	return status, nil
}

func (hs *SPShard) IsPayStatus() bool {
	fmt.Printf("IsPayStatus Current:%v,  hshPayStatus:%v \n", hs.fsm.Current(), hshPayStatus)
	return hs.fsm.Current() == hshPayStatus
}
func (hs *SPShard) IsContractStatus() bool {
	fmt.Printf("IsContractStatus Current:%v,  hshContractStatus:%v \n", hs.fsm.Current(), hshContractStatus)
	return hs.fsm.Current() == hshContractStatus
}

func (hs *SPShard) UpdateToContractStatus(signedGuardContract *metadata.Contract) error {
	return hs.fsm.Event(hshToContractEvent, signedGuardContract)
}

func (hs *SPShard) ReceivePayCheque() error {
	fmt.Printf("ReceivePayCheque cur=%+v \n", hs.fsm.Current())
	return hs.fsm.Event(hshToPayEvent)
}

func (hs *SPShard) Complete() error {
	return hs.fsm.Event(hshToCompleteEvent)
}

func (hs *SPShard) UpdateContractStatus() error {
	meta := &metadata.Contract{}
	err := Get(hs.ds, fmt.Sprintf(hostShardContractsKey, hs.peerId, hs.contractId), meta)
	if err != nil {
		return err
	}
	meta.Status = metadata.Contract_COMPLETED
	return Save(hs.ds, fmt.Sprintf(hostShardContractsKey, hs.peerId, hs.contractId), meta)
}
