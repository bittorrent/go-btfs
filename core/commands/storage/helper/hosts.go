package helper

import (
	"context"
	"fmt"
	"strings"

	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/hub"
	"github.com/bittorrent/go-btfs/repo"

	hubpb "github.com/tron-us/go-btfs-common/protos/hub"
	nodepb "github.com/tron-us/go-btfs-common/protos/node"

	"github.com/dustin/go-humanize"
	"github.com/gogo/protobuf/proto"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	"github.com/shirou/gopsutil/v3/disk"
)

const (
	HostStorePrefix       = "/hosts/"        // from btfs-hub
	HostStorageInfoPrefix = "/host_storage/" // self or from network
)

// GetHostsFromDatastore retrieves `num` hosts from the datastore, if not enough hosts are
// available, return an error instead of partial return.
// When num=0 it means unlimited.
func GetHostsFromDatastore(ctx context.Context, node *core.IpfsNode, mode string, num int) ([]*hubpb.Host, error) {
	// Check valid mode, including all (everything)
	_, mp, err := hub.CheckValidMode(mode, true)
	if err != nil {
		return nil, err
	}

	// get host list from storage
	rds := node.Repo.Datastore()
	qr, err := rds.Query(ctx, query.Query{
		Prefix: HostStorePrefix + mp,
		Orders: []query.Order{query.OrderByKey{}},
	})
	if err != nil {
		return nil, err
	}
	// Add as many hosts as available
	var hosts []*hubpb.Host
	for r := range qr.Next() {
		if r.Error != nil {
			return nil, r.Error
		}
		var h hubpb.Host
		err := proto.Unmarshal(r.Entry.Value, &h)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, &h)
	}
	// we can re-use hosts, but for higher availability, we choose to have the
	// greater than `num` assumption
	if num > 0 && len(hosts) < num {
		return nil, fmt.Errorf("there are not enough locally stored hosts")
	}
	return hosts, nil
}

// GetHostStorageConfigForPeer retrieves locally saved info about peer (including self)
func GetHostStorageConfigForPeer(ctx context.Context, node *core.IpfsNode, peerID string) (*nodepb.Node_Settings, error) {
	rds := node.Repo.Datastore()
	b, err := rds.Get(ctx, GetHostStorageKey(peerID))
	if err != nil {
		return nil, err
	}
	ns := new(nodepb.Node_Settings)
	err = ns.Unmarshal(b)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

// GetHostStorageConfig checks if locally is storing a config, if yes, returns it,
// otherwise, queries hub to retrieve the latest default config.
func GetHostStorageConfig(ctx context.Context, node *core.IpfsNode) (*nodepb.Node_Settings, error) {
	return GetHostStorageConfigHelper(ctx, node, false)
}

// GetHostStorageConfigHelper checks if locally is storing a config, if yes, returns it,
// otherwise, queries hub to retrieve the latest default config.
// If syncHub is on, force a sync from Hub to retrieve latest information.
func GetHostStorageConfigHelper(ctx context.Context, node *core.IpfsNode,
	syncHub bool) (*nodepb.Node_Settings, error) {
	ns, err := GetHostStorageConfigForPeer(ctx, node, node.Identity.Pretty())
	if err != nil && err != ds.ErrNotFound {
		return nil, fmt.Errorf("cannot get current host storage settings: %s", err.Error())
	}
	// Exists
	if err == nil && !syncHub {
		return ns, nil
	}
	cfg, err := node.Repo.Config()
	if err != nil {
		return nil, err
	}
	hns, err := hub.GetHostSettings(ctx, cfg.Services.HubDomain, node.Identity.Pretty())
	if err != nil {
		return nil, err
	}
	// ns aleady exists, so replace with newer settings
	if ns != nil {
		ns.StoragePriceDefault = hns.StoragePriceDefault
		if !ns.CustomizedPricing {
			ns.StoragePriceAsk = hns.StoragePriceAsk
		}
	} else {
		ns = hns
	}
	err = PutHostStorageConfig(ctx, node, ns)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

// PutHostStorageConfig saves an updated storage storage config.
func PutHostStorageConfig(ctx context.Context, node *core.IpfsNode, ns *nodepb.Node_Settings) error {
	rds := node.Repo.Datastore()
	b, err := ns.Marshal()
	if err != nil {
		return fmt.Errorf("cannot put current host storage settings: %s", err.Error())
	}
	return rds.Put(ctx, GetHostStorageKey(node.Identity.Pretty()), b)
}

func GetHostStorageKey(pid string) ds.Key {
	return NewKeyHelper(HostStorageInfoPrefix, pid)
}

func NewKeyHelper(kss ...string) ds.Key {
	return ds.NewKey(strings.Join(kss, ""))
}

// SaveHostsIntoDatastore overwrites (removes all existing) hosts and saves the updated
// hosts according to mode.
func SaveHostsIntoDatastore(ctx context.Context, node *core.IpfsNode, mode string, nodes []*hubpb.Host) error {
	// Check valid mode, including all (everything)
	_, mp, err := hub.CheckValidMode(mode, true)
	if err != nil {
		return err
	}

	rds := node.Repo.Datastore()

	// Dumb strategy right now: remove all existing and add the new ones
	// TODO: Update by timestamp and only overwrite updated
	qr, err := rds.Query(ctx, query.Query{Prefix: HostStorePrefix + mp})
	if err != nil {
		return err
	}

	for r := range qr.Next() {
		if r.Error != nil {
			return r.Error
		}
		err := rds.Delete(ctx, NewKeyHelper(r.Entry.Key))
		if err != nil {
			return err
		}
	}

	for i, ni := range nodes {
		b, err := proto.Marshal(ni)
		if err != nil {
			return err
		}
		err = rds.Put(ctx, NewKeyHelper(HostStorePrefix, mp, "/", fmt.Sprintf("%04d", i), "/", ni.NodeId), b)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckAndValidateHostStorageMax makes sure the current storage max is under the accepted
// disk space max, if not, corrects this value.
// Optionally, this function can take a new max and sets the max to this value.
// Optionally, maxAllowed enables reducing unreasonable settings down to an allowed value.
func CheckAndValidateHostStorageMax(ctx context.Context, cfgRoot string, r repo.Repo,
	newMax *uint64, maxAllowed bool) (uint64, error) {
	cfg, err := r.Config()
	if err != nil {
		return 0, err
	}

	// Check current available space + already used for accurate counting of max storage
	su, err := r.GetStorageUsage()
	if err != nil {
		return 0, err
	}
	du, err := disk.UsageWithContext(ctx, cfgRoot)
	if err != nil {
		return 0, err
	}
	totalAvailable := su + du.Free

	// Setting a new max storage, check if it exceeds available space
	if newMax != nil {
		if *newMax > totalAvailable {
			return 0, fmt.Errorf("new max storage size is invalid (exceeds available space)")
		}
		if *newMax < su {
			return 0, fmt.Errorf("new max storage size is invalid (lower than currently used)")
		}
		cfg.Datastore.StorageMax = humanize.Bytes(*newMax)
		err = r.SetConfig(cfg)
		if err != nil {
			return 0, err
		}
		// Return newly updated size
		return *newMax, nil
	}

	// Grab existing setting and verify
	curMax, err := humanize.ParseBytes(cfg.Datastore.StorageMax)
	if err != nil {
		return 0, err
	}

	if curMax > totalAvailable {
		if !maxAllowed {
			return 0, fmt.Errorf("current max storage size is invalid (exceeds available space)")
		}
		// Reduce current settings down to the max allowed
		cfg.Datastore.StorageMax = humanize.Bytes(totalAvailable)
		err = r.SetConfig(cfg)
		if err != nil {
			return 0, err
		}
		return totalAvailable, nil
	}

	// No change
	return curMax, nil
}
