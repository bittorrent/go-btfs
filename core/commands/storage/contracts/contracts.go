package contracts

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/utils"

	cmds "github.com/bittorrent/go-btfs-cmds"
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

	contractsListOrderOptionName  = "order"
	contractsListStatusOptionName = "status"
	contractsListSizeOptionName   = "size"

	contractsKeyPrefix = "/btfs/%s/contracts/"
	hostContractsKey   = contractsKeyPrefix + "host"
	renterContractsKey = contractsKeyPrefix + "renter"
	payoutNotFoundErr  = "rpc error: code = Unknown desc = not found"

	guardTimeout = 360 * time.Second

	guardContractPageSize = 100

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

		purgeOpt, _ := req.Options[contractsSyncPurgeOptionName].(bool)
		if purgeOpt {
			err = sessions.DeleteShardsContracts(n.Repo.Datastore(), n.Identity.String(), role.String())
			if err != nil {
				return err
			}
			err = Save(n.Repo.Datastore(), nil, role.String())
			if err != nil {
				return err
			}
		}
		ctx := context.WithValue(req.Context, contractsSyncVerboseOptionName, req.Options[contractsSyncVerboseOptionName].(bool))
		return SyncContracts(ctx, n, req, env, role.String())
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
		sort.Sort(ByTime(contracts))
		if parts[1] == "" || parts[1] == "desc" {
			// reverse
			for i, j := 0, len(contracts)-1; i < j; i, j = i+1, j-1 {
				contracts[i], contracts[j] = contracts[j], contracts[i]
			}
		}
		result := make([]*nodepb.Contracts_Contract, 0)
		for _, c := range contracts {
			if _, ok := states[c.Status]; !ok {
				continue
			}
			result = append(result, c)
			if len(result) == size {
				break
			}
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
	cs := &contractspb.Contracts{}
	err := sessions.Get(d, getKey(role), cs)
	if err != nil && err != datastore.ErrNotFound {
		return nil, err
	}
	// Because of buggy data in the past, we need to filter out non-host or non-renter contracts
	// It's also possible that user has switched keys manually, so we remove those as well.
	var fcs []*nodepb.Contracts_Contract
	filtered := false
	for _, c := range cs.Contracts {
		if role == nodepb.ContractStat_HOST.String() && c.HostId != peerId {
			filtered = true
			continue
		}
		if role == nodepb.ContractStat_RENTER.String() && c.RenterId != peerId {
			filtered = true
			continue
		}
		fcs = append(fcs, c)
	}
	// No change
	if !filtered {
		return cs.Contracts, nil
	}
	err = Save(d, fcs, role)
	if err != nil {
		return nil, err
	}
	return fcs, nil
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
	var updated []*metadata.Agreement
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
	return nil
}

func GetInvalidContractsForHost(cs []*metadata.Agreement, spId string) ([]*metadata.Agreement, error) {
	var invalid []*metadata.Agreement
	for _, c := range cs {
		if int64(c.Meta.StorageEnd) < time.Now().Unix() && c.Meta.SpId == spId {
			invalid = append(invalid, c)
		}
	}
	return invalid, nil
}

func GetInvalidContractForUser(cs []*metadata.Agreement, peerId string) ([]*metadata.Agreement, error) {
	var invalid []*metadata.Agreement
	for _, c := range cs {
		if c.Meta.CreatorId == peerId && int64(c.Meta.StorageEnd) < time.Now().Unix() {
			// If the contract is expired, we consider it invalid
			invalid = append(invalid, c)
		}
	}
	return invalid, nil
}
