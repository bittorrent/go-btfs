package contracts

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	nodepb "github.com/bittorrent/go-btfs-common/protos/node"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/protos/metadata"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/bittorrent/protobuf/proto"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ipfs/go-datastore"
)

var (
	fileMetaABI      = transaction.ParseABIUnchecked(abi.FileMetaContractABI)
	fileMetaAddEvent = fileMetaABI.Events["FileMetaAdded"].ID
)

// sync contract from bttc chain
func SyncContractFromBttcChain(startBlock, endBlock int64) {
	// chain.ChainObject.Backend.
}

const (
	defaultPollingInterval = 5 * time.Second
	blockPage              = 500
	tailSize               = 0
)

func ScanChainAndSave(d datastore.Datastore, role, identity string, from uint64) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan struct{})
	var wg sync.WaitGroup
	var mu sync.Mutex
	contracts := make([]*nodepb.Contracts_Contract, 0, 100)

	go func() {
		<-quit
		cancel()
	}()

	chainUpdateInterval := defaultPollingInterval
	paged := make(chan struct{}, 1)
	paged <- struct{}{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := syncBlocks(ctx, d, role, identity, from, &contracts, &mu, paged, quit, chainUpdateInterval); err != nil {
			contractsLog.Errorf("BttcListener returned with err: %v", err)
		}
	}()

	wg.Wait()

	if len(contracts) > 0 {
		Save(d, contracts, role)
	}
}

func filterQuery(from, to *big.Int) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		FromBlock: from,
		ToBlock:   to,
		Addresses: []common.Address{chain.ChainObject.Chainconfig.FileMeta2Address},
		Topics: [][]common.Hash{
			{
				fileMetaAddEvent,
			},
		},
	}
}

type FileMetaAddEvent struct {
	Cid      string
	MetaData []byte
}

func processEvent(e types.Log, role, identity string) []*metadata.Contract {
	c := &FileMetaAddEvent{}
	err := transaction.ParseEvent(&fileMetaABI, "FileMetaAdded", c, e)
	if err != nil {
		fmt.Println(err)
	}

	meta := &metadata.FileMetaInfo{}
	err = proto.Unmarshal(c.MetaData, meta)
	if err != nil {
		fmt.Println(err)
	}

	for _, c := range meta.Contracts {
		if role == nodepb.ContractStat_HOST.String() && c.Meta.SpId != identity {
			continue
		}
		if role == nodepb.ContractStat_RENTER.String() && c.Meta.UserId != identity {
			continue
		}
	}

	return meta.Contracts
}

func syncBlocks(ctx context.Context, d datastore.Datastore, role, identity string, from uint64,
	contracts *[]*nodepb.Contracts_Contract, mu *sync.Mutex, paged, quit chan struct{}, interval time.Duration) error {

	for {
		select {
		case <-paged:
		case <-time.After(interval):
		case <-quit:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}

		to, err := chain.ChainObject.Backend.BlockNumber(ctx)
		if err != nil {
			contractsLog.Warnf("BttcListener: could not get block number: %v", err)
			continue
		}

		if to < tailSize {
			contractsLog.Infof("BttcListener: current last block number: %v < tailSize:%v", to, tailSize)
			continue
		}

		to = to - tailSize

		if to < from {
			contractsLog.Warnf("BttcListener: latest block number:%+v < from:%+v", to, from)
			continue
		}

		if to-from >= blockPage {
			select {
			case paged <- struct{}{}:
			default:
				contractsLog.Warnf("Paged channel is full, skipping page signal")
			}
			to = from + blockPage - 1
		} else {
			log.Infof("contracts sync finished...")
			break
		}

		contractsLog.Infof("start sync contracts from:%v, to = %+v", from, to)

		events, err := chain.ChainObject.Backend.FilterLogs(ctx, filterQuery(big.NewInt(int64(from)), big.NewInt(int64(to))))
		if err != nil {
			contractsLog.Warnf("sync contracts could not get logs: %v", err)
			continue
		}

		if len(events) > 0 {
			processEvents(events, role, identity, contracts, mu)
		}

		from = to + 1
	}

	if len(*contracts) > 0 {
		mu.Lock()
		tempContracts := make([]*nodepb.Contracts_Contract, len(*contracts))
		copy(tempContracts, *contracts)
		*contracts = (*contracts)[:0]
		mu.Unlock()

		Save(d, tempContracts, role)
	}

	return nil
}

func processEvents(events []types.Log, role, identity string,
	contracts *[]*nodepb.Contracts_Contract, mu *sync.Mutex) {

	for _, e := range events {
		cs := processEvent(e, role, identity)
		if len(cs) > 0 {
			mu.Lock()
			*contracts = append(*contracts, convertMetadataContractsToNodeContracts(cs)...)
			mu.Unlock()
		}
	}
}
