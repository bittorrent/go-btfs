package contracts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/utils"

	cmds "github.com/bittorrent/go-btfs-cmds"
	guardpb "github.com/bittorrent/go-btfs-common/protos/guard"
	nodepb "github.com/bittorrent/go-btfs-common/protos/node"
	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/rm"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/logger"
	contractspb "github.com/bittorrent/go-btfs/protos/contracts"
	"github.com/bittorrent/go-btfs/protos/metadata"

	"github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log"
)

var contractsLog = logging.Logger("storage/contracts")

var log = logger.InitLogger("contracts.log").Sugar()

const (
	contractsSyncPurgeOptionName   = "purge"
	contractsSyncVerboseOptionName = "verbose"
	contractsSyncBlockHeight       = "height"

	contractsListOrderOptionName  = "order"
	contractsListStatusOptionName = "status"
	contractsListSizeOptionName   = "size"

	contractsKeyPrefix = "/btfs/%s/contracts/"
	hostContractsKey   = contractsKeyPrefix + "host"
	renterContractsKey = contractsKeyPrefix + "renter"
	payoutNotFoundErr  = "rpc error: code = Unknown desc = not found"

	notSupportErr = "only host and renter contract sync are supported currently"
)

// Storage contracts
// Includes sub-commands: sync, stat, list

var StorageContractsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get node storage contracts info.",
		ShortDescription: `
This command get node storage contracts info respect to different roles.`,
	},
	Subcommands: map[string]*cmds.Command{
		"sync": storageContractsSyncCmd,
		"stat": storageContractsStatCmd,
		"list": storageContractsListCmd,
	},
}

// checkContractStatRole checks role argument strings against valid roles
// and returns the role type
func checkContractStatRole(roleArg string) (nodepb.ContractStat_Role, error) {
	if cr, ok := nodepb.ContractStat_Role_value[strings.ToUpper(roleArg)]; ok {
		return nodepb.ContractStat_Role(cr), nil
	}
	return 0, fmt.Errorf("invalid role: %s", roleArg)
}

// sub-commands: btfs storage contracts sync
var storageContractsSyncCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Synchronize contracts stats based on role.",
		ShortDescription: `
This command contracts stats based on role from network(hub) to local node data storage.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("role", true, false, "Role in BTFS storage network [host|renter|reserved]."),
	},
	Options: []cmds.Option{
		cmds.BoolOption(contractsSyncPurgeOptionName, "p", "Purge local contracts cache and sync from the beginning.").WithDefault(false),
		cmds.BoolOption(contractsSyncVerboseOptionName, "v", "Make the operation more talkative.").WithDefault(false),
		cmds.Int64Option(contractsSyncBlockHeight, "sh", "Start block height to sync contracts.").WithDefault(53849058),
	},
	RunTimeout: 10 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		role, err := checkContractStatRole(req.Arguments[0])
		if err != nil {
			return err
		}
		if role != nodepb.ContractStat_HOST && role != nodepb.ContractStat_RENTER {
			return fmt.Errorf(notSupportErr)
		}

		ctx := context.WithValue(req.Context, contractsSyncVerboseOptionName, req.Options[contractsSyncVerboseOptionName].(bool))

		purgeOpt, _ := req.Options[contractsSyncPurgeOptionName].(bool)
		if purgeOpt {
			err = sessions.DeleteShardsContracts(n.Repo.Datastore(), n.Identity.String(), role.String())
			if err != nil {
				return err
			}
			go func() {
				ScanChainAndSave(n.Repo.Datastore(), role.String(), n.Identity.String(), uint64(req.Options[contractsSyncBlockHeight].(int64)))
				SyncContracts(ctx, n, req, env, role.String())
			}()
			// err = Save(n.Repo.Datastore(), nil, role.String())
			if err != nil {
				return err
			}
			return nil
		}
		SyncContracts(ctx, n, req, env, role.String())
		return err
	},
}

// sub-commands: btfs storage contracts stat
var storageContractsStatCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get contracts stats based on role.",
		ShortDescription: `
This command get contracts stats based on role from the local node data store.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("role", true, false, "Role in BTFS storage network [host|renter|reserved]."),
	},
	RunTimeout: 3 * time.Second,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		cr, err := checkContractStatRole(req.Arguments[0])
		if err != nil {
			return err
		}
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		contracts, err := ListContracts(n.Repo.Datastore(), n.Identity.String(), cr.String())
		if err != nil {
			return err
		}
		activeStates := helper.ContractFilterMap["active"]
		invalidStates := helper.ContractFilterMap["invalid"]
		var activeCount, totalPaid, totalUnpaid int64
		var first, last time.Time
		for _, c := range contracts {
			if _, ok := activeStates[c.Status]; ok {
				activeCount++
				// Count outstanding on only active ones
				totalUnpaid += c.CompensationOutstanding
			}
			// Count all paid on all contracts
			totalPaid += c.CompensationPaid
			// Count start/end for all non-invalid ones
			if _, ok := invalidStates[c.Status]; !ok {
				if (first == time.Time{}) || c.StartTime.Before(first) {
					first = c.StartTime
				}
				if (last == time.Time{}) || c.EndTime.After(last) {
					last = c.EndTime
				}
			}
		}
		data := &nodepb.ContractStat{
			ActiveContractNum:       activeCount,
			CompensationPaid:        totalPaid,
			CompensationOutstanding: totalUnpaid,
			FirstContractStart:      first,
			LastContractEnd:         last,
			Role:                    cr,
		}
		return cmds.EmitOnce(res, data)
	},
	Type: nodepb.ContractStat{},
}

var (
	contractOrderList = []string{"escrow_time"}
)

type ByTime []*nodepb.Contracts_Contract

func (a ByTime) Len() int      { return len(a) }
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool {
	return a[i].NextEscrowTime.UnixNano() < a[j].NextEscrowTime.UnixNano()
}

// sub-commands: btfs storage contracts list
var storageContractsListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get contracts list based on role.",
		ShortDescription: `
This command get contracts list based on role from the local node data store.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("role", true, false, "Role in BTFS storage network [host|renter|reserved]."),
	},
	Options: []cmds.Option{
		cmds.StringOption(contractsListOrderOptionName, "o", "Order to return the list of contracts.").WithDefault("escrow_time,asc"),
		cmds.StringOption(contractsListStatusOptionName, "st", "Filter the returned list by contract status [active|finished|invalid|all].").WithDefault("active"),
		cmds.IntOption(contractsListSizeOptionName, "s", "Number of contracts to return.").WithDefault(20),
	},
	RunTimeout: 3 * time.Second,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		cr, err := checkContractStatRole(req.Arguments[0])
		if err != nil {
			return err
		}
		orderOpt, _ := req.Options[contractsListOrderOptionName].(string)
		parts := strings.Split(orderOpt, ",")
		if len(parts) != 2 {
			return fmt.Errorf(`bad order format "<order-name>,<order-direction>"`)
		}
		var order string
		for _, o := range contractOrderList {
			if o == parts[0] {
				order = o
				break
			}
		}
		if order == "" {
			return fmt.Errorf("bad order name: %s", parts[0])
		}
		if parts[1] != "asc" && parts[1] != "desc" {
			return fmt.Errorf("bad order direction: %s", parts[1])
		}
		filterOpt, _ := req.Options[contractsListStatusOptionName].(string)
		states, ok := helper.ContractFilterMap[filterOpt]
		if !ok {
			return fmt.Errorf("invalid filter option: %s", filterOpt)
		}
		size, _ := req.Options[contractsListSizeOptionName].(int)
		contracts, err := ListContracts(n.Repo.Datastore(), n.Identity.String(), cr.String())
		if err != nil {
			return err
		}

		// First filter by status
		filteredContracts := make([]*nodepb.Contracts_Contract, 0)
		for _, c := range contracts {
			if _, ok := states[c.Status]; !ok {
				continue
			}

			filteredContracts = append(filteredContracts, c)
		}

		// Sort by escrow time
		sort.Sort(ByTime(filteredContracts))
		if parts[1] == "" || parts[1] == "desc" {
			// reverse
			for i, j := 0, len(filteredContracts)-1; i < j; i, j = i+1, j-1 {
				filteredContracts[i], filteredContracts[j] = filteredContracts[j], filteredContracts[i]
			}
		}

		// Apply size limit
		result := filteredContracts
		if len(result) > size {
			result = result[:size]
		}

		return cmds.EmitOnce(res, &nodepb.Contracts{Contracts: result})
	},
	Type: nodepb.Contracts{},
}

func getKey(role string) string {
	var k string
	if role == nodepb.ContractStat_HOST.String() {
		k = hostContractsKey
	} else if role == nodepb.ContractStat_RENTER.String() {
		k = renterContractsKey
	} else {
		return "reserved"
	}
	return k
}

func Save(d datastore.Datastore, cs []*nodepb.Contracts_Contract, role string) error {
	return sessions.Save(d, getKey(role), &contractspb.Contracts{
		Contracts: cs,
	})
}

func ListContracts(d datastore.Datastore, peerId, role string) ([]*nodepb.Contracts_Contract, error) {
	metadataContracts, err := sessions.ListShardsContracts(d, peerId, role)
	if err != nil {
		return nil, fmt.Errorf("failed to list shard contracts: %v", err)
	}

	nodeContracts := convertMetadataContractsToNodeContracts(metadataContracts)

	cs := &contractspb.Contracts{}
	err = sessions.Get(d, getKey(role), cs)
	if err != nil && err != datastore.ErrNotFound {
		return nodeContracts, nil
	}

	existingContracts := make(map[string]bool)
	for _, c := range nodeContracts {
		existingContracts[c.ContractId] = true
	}

	for _, c := range cs.Contracts {
		if existingContracts[c.ContractId] {
			continue
		}

		if role == nodepb.ContractStat_HOST.String() && c.HostId != peerId {
			continue
		}
		if role == nodepb.ContractStat_RENTER.String() && c.RenterId != peerId {
			continue
		}

		nodeContracts = append(nodeContracts, c)
		existingContracts[c.ContractId] = true
	}

	if len(nodeContracts) > 0 {
		if err := Save(d, nodeContracts, role); err != nil {
			log.Warnf("Failed to save contracts to legacy path: %v", err)
		}
	}

	return nodeContracts, nil
}

func SyncContracts(ctx context.Context, n *core.IpfsNode, req *cmds.Request, env cmds.Environment, role string) error {
	cs, err := sessions.ListShardsContracts(n.Repo.Datastore(), n.Identity.String(), role)
	if err != nil {
		return err
	}
	var latest *time.Time
	for _, c := range cs {
		createTime := time.Unix(int64(c.CreateTime), 0)
		if latest == nil || createTime.After(*latest) {
			latest = &createTime
		}
	}
	var updated []*metadata.Contract
	switch role {
	case nodepb.ContractStat_HOST.String():
		updated, err = GetInvalidContractsForHost(cs, n.Identity.String())
		if err != nil {
			return err
		}
	case nodepb.ContractStat_RENTER.String():
		updated, err = GetInvalidContractForUser(cs, n.Identity.String())
		if err != nil {
			return err
		}
	default:
		return errors.New(notSupportErr)
	}

	if len(updated) > 0 {
		// save and retrieve updated signed contracts
		// cs, stale, err = sessions.SaveShardsContract(n.Repo.Datastore(), cs, updated, n.Identity.String(), role)
		stales, err := sessions.RefreshLocalContracts(ctx, n.Repo.Datastore(), cs, updated, n.Identity.String(), role)
		if err != nil {
			return err
		}
		if role == nodepb.ContractStat_HOST.String() {
			go func() {
				// Use a new context that can clean up in the background
				_, err := rm.RmDag(context.Background(), stales, n, req, env, true)
				if err != nil {
					// may have been cleaned up already, ignore
					contractsLog.Error("stale contracts clean up error:", err)
				}
			}()
		}
	}

	if len(cs) > 0 {
		results := convertMetadataContractsToNodeContracts(cs)

		if b := ctx.Value(contractsSyncVerboseOptionName); b != nil && b.(bool) {
			go func() {
				for _, ct := range results {
					if bs, err := json.Marshal(ct); err == nil {
						log.Info(string(bs))
					}
				}
			}()
		}
		return Save(n.Repo.Datastore(), results, role)
	}
	return nil
}

func convertMetadataContractsToNodeContracts(contracts []*metadata.Contract) []*nodepb.Contracts_Contract {
	results := make([]*nodepb.Contracts_Contract, 0, len(contracts))

	for _, c := range contracts {
		if c.Meta == nil {
			continue
		}

		var status guardpb.Contract_ContractState
		switch c.Status {
		case metadata.Contract_COMPLETED:
			status = guardpb.Contract_UPLOADED
		case metadata.Contract_INIT:
			status = guardpb.Contract_SIGNED
		default:
			status = guardpb.Contract_SIGNED
		}

		endTime := time.Unix(int64(c.Meta.StorageEnd), 0)
		if time.Now().Unix() > endTime.Unix() {
			status = guardpb.Contract_CLOSED
		}

		nc := &nodepb.Contracts_Contract{
			ContractId: c.Meta.ContractId,
			HostId:     c.Meta.SpId,
			RenterId:   c.Meta.UserId,
			Status:     status,
			StartTime:  time.Unix(int64(c.Meta.StorageStart), 0),
			EndTime:    endTime,
			UnitPrice:  int64(c.Meta.Price),
			ShardSize:  int64(c.Meta.ShardSize),
			ShardHash:  c.Meta.ShardHash,
		}

		results = append(results, nc)
	}

	return results
}

func GetInvalidContractsForHost(cs []*metadata.Contract, spId string) ([]*metadata.Contract, error) {
	var invalid []*metadata.Contract
	for _, c := range cs {
		if int64(c.Meta.StorageEnd) < time.Now().Unix() && c.Meta.SpId == spId {
			invalid = append(invalid, c)
		}

		if c.Meta.SpId == spId && (c.Status == metadata.Contract_CLOSED || c.Status == metadata.Contract_INVALID) {
			invalid = append(invalid, c)
		}

		if c.Meta.SpId == spId && c.Status == metadata.Contract_INIT && c.CreateTime > uint64(time.Now().Unix())-uint64(time.Hour) {
			invalid = append(invalid, c)
		}

	}
	return invalid, nil
}

func GetInvalidContractForUser(cs []*metadata.Contract, peerId string) ([]*metadata.Contract, error) {
	var invalid []*metadata.Contract
	for _, c := range cs {
		if c.Meta.UserId == peerId && int64(c.Meta.StorageEnd) < time.Now().Unix() {
			// If the contract is expired, we consider it invalid
			invalid = append(invalid, c)
		}
	}
	return invalid, nil
}
