package upload

import (
	"github.com/bittorrent/go-btfs/core/commands/storage/contracts"

	guardpb "github.com/bittorrent/go-btfs-common/protos/guard"
	"github.com/bittorrent/go-btfs-common/protos/node"
	config "github.com/bittorrent/go-btfs-config"

	"github.com/ipfs/go-datastore"
)

type ICount interface {
	Count(ds datastore.Datastore, peerId string, status guardpb.Contract_ContractState) (int, error)
}

type Count struct{}

type HostManager struct {
	low       int
	high      int
	threshold int64
	count     ICount
}

func NewHostManager(cfg *config.Config) *HostManager {
	return &HostManager{
		low:       cfg.UI.Host.ContractManager.LowWater,
		high:      cfg.UI.Host.ContractManager.HighWater,
		threshold: cfg.UI.Host.ContractManager.Threshold,
		count:     &Count{},
	}
}

func (h *HostManager) AcceptContract(ds datastore.Datastore, peerId string, shardSize int64) (bool, error) {
	count, err := h.count.Count(ds, peerId, guardpb.Contract_READY_CHALLENGE)
	if err != nil {
		log.Debug("err", err)
		return true, nil
	}
	if count <= h.low {
		return true, nil
	} else if count >= h.high {
		return false, nil
	} else {
		return shardSize <= h.threshold, nil
	}
}

func (h *Count) Count(ds datastore.Datastore, peerId string, status guardpb.Contract_ContractState) (int, error) {
	contracts, err := contracts.ListContracts(ds, peerId, node.ContractStat_HOST.String())
	if err != nil {
		return 0, err
	}
	c := 0
	for i := 0; i < len(contracts); i++ {
		if contracts[i].Status == status {
			c++
		}
	}
	return c, nil
}
