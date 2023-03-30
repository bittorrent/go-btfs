package helper

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	unixtest "github.com/bittorrent/go-btfs/core/coreunix/test"
	"github.com/bittorrent/go-btfs/repo"

	hubpb "github.com/tron-us/go-btfs-common/protos/hub"
	nodepb "github.com/tron-us/go-btfs-common/protos/node"

	config "github.com/TRON-US/go-btfs-config"
	"github.com/alecthomas/units"
	"github.com/gogo/protobuf/proto"
)

func TestHostsSaveGet(t *testing.T) {
	node := unixtest.HelpTestMockRepo(t, nil)

	// test all possible modes
	for _, mode := range hubpb.HostsReq_Mode_name {
		var nodes []*hubpb.Host
		for i := 0; i < 100; i++ {
			ni := &hubpb.Host{NodeId: fmt.Sprintf("%s:node:%d", mode, i)}
			nodes = append(nodes, ni)
		}
		err := SaveHostsIntoDatastore(context.Background(), node, mode, nodes)
		if err != nil {
			t.Fatal(err)
		}
		stored, err := GetHostsFromDatastore(context.Background(), node, mode, 100)
		if err != nil {
			t.Fatal(err)
		}
		for i, sn := range stored {
			bs1, err := proto.Marshal(sn)
			if err != nil {
				t.Fatal(err)
			}
			bs2, err := proto.Marshal(nodes[i])
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Compare(bs1, bs2) != 0 {
				t.Fatal("stored nodes do not match saved nodes")
			}
		}
	}
}

func TestHostStorageConfigPutGet(t *testing.T) {
	node := unixtest.HelpTestMockRepo(t, nil)

	ns := &nodepb.Node_Settings{
		StoragePriceAsk:   111111,
		BandwidthPriceAsk: 111,
		StorageTimeMin:    15,
		BandwidthLimit:    999,
		CollateralStake:   1000,
	}
	ctx := context.TODO()
	err := PutHostStorageConfig(ctx, node, ns)
	if err != nil {
		t.Fatal(err)
	}
	stored, err := GetHostStorageConfig(ctx, node)
	if err != nil {
		t.Fatal(err)
	}
	bs1, err := proto.Marshal(ns)
	if err != nil {
		t.Fatal(err)
	}
	bs2, err := proto.Marshal(stored)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(bs1, bs2) != 0 {
		t.Fatal("stored settings do not match saved settings")
	}
}

func newSmCfg(max string) *config.Config {
	cfg := &config.Config{}
	cfg.Datastore.StorageMax = max
	return cfg
}

func TestCheckAndValidateHostStorageMax(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// test normal case of checking for current storage
	node := unixtest.HelpTestMockRepo(t, newSmCfg("1GB"))
	m, err := CheckAndValidateHostStorageMax(context.Background(), pwd, node.Repo, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if m != uint64(units.GB) {
		t.Fatalf("current max is incorrect, got %v, want %v", m, uint64(units.GB))
	}

	// test oversetting case of checking for current storage
	node = unixtest.HelpTestMockRepo(t, newSmCfg("1PB"))
	_, err = CheckAndValidateHostStorageMax(context.Background(), pwd, node.Repo, nil, false)
	if err == nil {
		t.Fatal("should return an error for current max")
	}

	// test oversetting case of checking for current storage, then correct to max allowed
	node = unixtest.HelpTestMockRepo(t, newSmCfg("1PB"))
	m, err = CheckAndValidateHostStorageMax(context.Background(), pwd, node.Repo, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if m == 0 {
		t.Fatal("should return max allowed for current max")
	}

	// test normal case of setting a new reasonable max
	node = unixtest.HelpTestMockRepo(t, newSmCfg("500MB"))
	nm := uint64(units.GB)
	m, err = CheckAndValidateHostStorageMax(context.Background(), pwd, node.Repo, &nm, false)
	if err != nil {
		t.Fatal(err)
	}
	if m != uint64(units.GB) {
		t.Fatalf("new max is incorrect, got %v, want %v", m, uint64(units.GB))
	}

	// test normal case of setting a new unreasonable max
	node = unixtest.HelpTestMockRepo(t, newSmCfg("500MB"))
	nm = uint64(units.PB)
	_, err = CheckAndValidateHostStorageMax(context.Background(), pwd, node.Repo, &nm, false)
	if err == nil {
		t.Fatal("should return error for too large of a new max")
	}
	_, err = CheckAndValidateHostStorageMax(context.Background(), pwd, node.Repo, &nm, true)
	if err == nil {
		t.Fatal("should return error for too large of a new max (even on max allowed flag)")
	}

	// test normal case of setting a new unreasonable max below used size
	node = unixtest.HelpTestMockRepo(t, newSmCfg("500MB"))
	// mock used size
	mn := node.Repo.(*repo.Mock)
	mn.Used = 10000
	nm = uint64(1000)
	_, err = CheckAndValidateHostStorageMax(context.Background(), pwd, node.Repo, &nm, false)
	if err == nil {
		t.Fatal("should return error for too small of a new max")
	}
	_, err = CheckAndValidateHostStorageMax(context.Background(), pwd, node.Repo, &nm, true)
	if err == nil {
		t.Fatal("should return error for too small of a new max (even on max allowed flag)")
	}
}
