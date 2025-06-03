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

	hshInitStatus      = "init"
	hshAgreementStatus = "agreement"
	hshPayStatus       = "paid"
	hshCompleteStatus  = "complete"
	hshErrorStatus     = "error"

	hshToAgreementEvent = "to-agreement"
	hshToPayEvent       = "to-pay"
	hshToCompleteEvent  = "to-complete"
	hshToErrorEvent     = "to-error"
)

var (
	spShardFsmEvents = fsm.Events{
		{Name: hshToAgreementEvent, Src: []string{hshInitStatus}, Dst: hshAgreementStatus},
		{Name: hshToPayEvent, Src: []string{hshAgreementStatus}, Dst: hshPayStatus},
		{Name: hshToCompleteEvent, Src: []string{hshPayStatus}, Dst: hshCompleteStatus},
		{Name: hshToErrorEvent, Src: []string{hshInitStatus, hshAgreementStatus}, Dst: hshToErrorEvent},
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
	case hshAgreementStatus:
		hs.doContract(e.Args[0].(*metadata.Agreement))
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

func (hs *SPShard) doContract(signedGuardContract *metadata.Agreement) error {
	status := &shardpb.Status{
		Status: hshAgreementStatus,
	}
	return Batch(hs.ds, []string{
		fmt.Sprintf(hostShardStatusKey, hs.peerId, hs.contractId),
		fmt.Sprintf(hostShardContractsKey, hs.peerId, hs.contractId),
	}, []proto.Message{
		status, signedGuardContract,
	})
}

func (hs *SPShard) Contract(signedGuardContract *metadata.Agreement) error {
	return hs.fsm.Event(hshToAgreementEvent, signedGuardContract)
}

func (hs *SPShard) IsPayStatus() bool {
	fmt.Printf("IsPayStatus Current:%v,  hshPayStatus:%v \n", hs.fsm.Current(), hshPayStatus)
	return hs.fsm.Current() == hshPayStatus
}
func (hs *SPShard) IsContractStatus() bool {
	fmt.Printf("IsContractStatus Current:%v,  hshAgreementStatus:%v \n", hs.fsm.Current(), hshAgreementStatus)
	return hs.fsm.Current() == hshAgreementStatus
}

func (hs *SPShard) ReceivePayCheque() error {
	fmt.Printf("ReceivePayCheque cur=%+v \n", hs.fsm.Current())
	return hs.fsm.Event(hshToPayEvent)
}

func (hs *SPShard) Complete() error {
	return hs.fsm.Event(hshToCompleteEvent)
}

func (hs *SPShard) contracts() (*shardpb.SignedContracts, error) {
	contracts := &shardpb.SignedContracts{}
	err := Get(hs.ds, fmt.Sprintf(hostShardContractsKey, hs.peerId, hs.contractId), contracts)
	if err == datastore.ErrNotFound {
		return contracts, nil
	}
	return contracts, err
}
