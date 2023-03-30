package stats

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/core/corerepo"
	"github.com/bittorrent/go-btfs/core/hub"

	config "github.com/TRON-US/go-btfs-config"
	cmds "github.com/bittorrent/go-btfs-cmds"
	nodepb "github.com/tron-us/go-btfs-common/protos/node"

	ds "github.com/ipfs/go-datastore"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/tron-us/protobuf/proto"
)

const (
	localInfoOnlyOptionName = "local-only"
	versionOptionName       = "version"
)

// Storage Stats
//
// Includes sub-commands: info, sync
var StorageStatsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get node storage stats.",
		ShortDescription: `
This command get node storage stats in the network.`,
	},
	Subcommands: map[string]*cmds.Command{
		"sync": storageStatsSyncCmd,
		"info": storageStatsInfoCmd,
		"list": storageStatsListCmd,
	},
}

// sub-commands: btfs storage stats sync
var storageStatsSyncCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Synchronize node stats.",
		ShortDescription: `
This command synchronize node stats from network(hub) to local node data store.`,
	},
	Arguments: []cmds.Argument{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}

		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		return SyncStats(req.Context, cfg, n, env, true)
	},
}

func SyncStats(ctx context.Context, cfg *config.Config, node *core.IpfsNode, env cmds.Environment, v2 bool) error {
	sr, err := hub.QueryStats(ctx, node, v2)
	if err != nil {
		return err
	}
	stat, err := corerepo.RepoStat(ctx, node)
	if err != nil {
		return err
	}
	cfgRoot, err := cmdenv.GetConfigRoot(env)
	if err != nil {
		return err
	}
	du, err := disk.UsageWithContext(ctx, cfgRoot)
	if err != nil {
		return err
	}
	hs := &nodepb.StorageStat_Host{
		Online:               cfg.Experimental.StorageHostEnabled,
		StorageUsed:          int64(stat.RepoSize),
		StorageCap:           int64(stat.StorageMax),
		StorageDiskTotal:     int64(du.Total),
		StorageDiskAvailable: int64(du.Free),
	}
	hs.StorageStat_HostStats = sr.StorageStat_HostStats
	return SaveHostStatsIntoDatastore(ctx, node, node.Identity.Pretty(), hs)
}

func GetNowStats(ctx context.Context, cfg *config.Config, node *core.IpfsNode, env cmds.Environment, V2 bool) (hs *nodepb.StorageStat_Host, err error) {
	sr, err := hub.QueryStats(ctx, node, V2)
	if err != nil {
		return nil, err
	}
	stat, err := corerepo.RepoStat(ctx, node)
	if err != nil {
		return nil, err
	}
	cfgRoot, err := cmdenv.GetConfigRoot(env)
	if err != nil {
		return nil, err
	}
	du, err := disk.UsageWithContext(ctx, cfgRoot)
	if err != nil {
		return nil, err
	}
	hs = &nodepb.StorageStat_Host{
		Online:               cfg.Experimental.StorageHostEnabled,
		StorageUsed:          int64(stat.RepoSize),
		StorageCap:           int64(stat.StorageMax),
		StorageDiskTotal:     int64(du.Total),
		StorageDiskAvailable: int64(du.Free),
	}
	hs.StorageStat_HostStats = sr.StorageStat_HostStats
	return hs, nil
}

// sub-commands: btfs storage stats info
var storageStatsInfoCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get node stats.",
		ShortDescription: `
This command get node stats in the network from the local node data store.`,
	},
	Arguments: []cmds.Argument{},
	Options: []cmds.Option{
		cmds.BoolOption(localInfoOnlyOptionName, "l", "Return only the locally available disk stats without querying/returning the network stats.").WithDefault(false),
		cmds.IntOption(versionOptionName, "v", "Get new hub score level.").WithDefault(2),
	},
	RunTimeout: 30 * time.Second,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}

		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		var v2Flag bool
		v, _ := req.Options[versionOptionName].(int)
		if v == 1 {
			v2Flag = false
		} else if v == 2 {
			v2Flag = true
		} else {
			return errors.New("version should be 1 or 2, not other. ")
		}

		var hs *nodepb.StorageStat_Host
		hs, err = GetNowStats(req.Context, cfg, n, env, v2Flag)
		if err != nil {
			return err
		}

		// Refresh latest repo stats
		stat, err := corerepo.RepoStat(req.Context, n)
		if err != nil {
			return err
		}

		cfgRoot, err := cmdenv.GetConfigRoot(env)
		if err != nil {
			return err
		}
		du, err := disk.UsageWithContext(req.Context, cfgRoot)
		if err != nil {
			return err
		}

		hs.Online = cfg.Experimental.StorageHostEnabled
		hs.StorageUsed = int64(stat.RepoSize)
		hs.StorageCap = int64(stat.StorageMax)
		hs.StorageDiskTotal = int64(du.Total)
		hs.StorageDiskAvailable = int64(du.Free)

		// Only host stats for now
		return cmds.EmitOnce(res, &nodepb.StorageStat{
			HostStats: *hs,
		})
	},
	Type: nodepb.StorageStat{},
}

// sub-commands: btfs storage stats list
var storageStatsListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List node stats.",
		ShortDescription: `
This command list node stats in the network from the local node data store.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("from", true, false, "list host local stats range from"),
		cmds.StringArg("to", true, false, "list host local stats range to"),
	},
	RunTimeout: 30 * time.Second,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		from, err := strconv.ParseInt(req.Arguments[0], 10, 64)
		if err != nil {
			return err
		}
		to, err := strconv.ParseInt(req.Arguments[1], 10, 64)
		if err != nil {
			return err
		}
		list, err := ListHostStatsFromDatastore(req.Context, n, n.Identity.String(), from, to)
		if err != nil {
			return err
		}

		// Only host stats for now
		return cmds.EmitOnce(res, list)
	},
	Type: []*Stat_HostWithTimeStamp{},
}

const (
	HostStatStorePrefix = "/host_stats/" // from btfs-hub
)

// GetHostStatsFromDatastore retrieves host storage stats based on node id
func GetHostStatsFromDatastore(ctx context.Context, node *core.IpfsNode, nodeId string) (*nodepb.StorageStat_Host, error) {
	rds := node.Repo.Datastore()
	qr, err := rds.Get(ctx, GetHostStatStorageKey(nodeId))
	if err != nil {
		return nil, err
	}
	var hs nodepb.StorageStat_Host
	err = proto.Unmarshal(qr, &hs)
	if err != nil {
		return nil, err
	}
	return &hs, nil
}

type Stat_HostWithTimeStamp struct {
	Stat      nodepb.StorageStat_Host `json:"stat"`
	Timestamp int64                   `json:"timestamp"`
}

// ListHostStatsFromDatastore retrieves host storage stats based on node id
func ListHostStatsFromDatastore(ctx context.Context, node *core.IpfsNode, nodeId string, from int64, to int64) ([]*Stat_HostWithTimeStamp, error) {
	rds := node.Repo.Datastore()
	keys, err := sessions.ListKeys(rds, HostStatStorePrefix+nodeId+"/", "")
	sort.Strings(keys)
	if err != nil {
		return nil, err
	}
	hosts := make([]*Stat_HostWithTimeStamp, 0)
	ly, lm, ld := -1, "", -1
	for _, k := range keys {
		qr, err := rds.Get(ctx, ds.NewKey(k))
		if err != nil {
			continue
		}
		var hs nodepb.StorageStat_Host
		err = proto.Unmarshal(qr, &hs)
		if err != nil {
			continue
		}
		split := strings.Split(k, "/")
		t, err := strconv.ParseInt(split[len(split)-1], 10, 64)
		if err != nil || t < from || t > to {
			continue
		}
		year, month, day := time.Unix(t, 0).Date()
		if ly == year && lm == month.String() && ld == day {
			continue
		}
		ly, lm, ld = year, month.String(), day
		hosts = append(hosts, &Stat_HostWithTimeStamp{
			Stat:      hs,
			Timestamp: t,
		})
	}
	return hosts, nil
}

func GetHostStatStorageKey(pid string) ds.Key {
	return helper.NewKeyHelper(HostStatStorePrefix, pid)
}

func GetHostStatStorageKeyWithTimestamp(pid string) ds.Key {
	return helper.NewKeyHelper(HostStatStorePrefix, pid, "/", strconv.FormatInt(time.Now().Unix(), 10))
}

// SaveHostStatsIntoDatastore overwrites host storage stats based on node id
func SaveHostStatsIntoDatastore(ctx context.Context, node *core.IpfsNode, nodeId string,
	stats *nodepb.StorageStat_Host) error {
	rds := node.Repo.Datastore()
	b, err := proto.Marshal(stats)
	if err != nil {
		return err
	}
	err = rds.Put(ctx, GetHostStatStorageKey(nodeId), b)
	if err != nil {
		return err
	}
	err = rds.Put(ctx, GetHostStatStorageKeyWithTimestamp(nodeId), b)
	if err != nil {
		return err
	}
	return nil
}
