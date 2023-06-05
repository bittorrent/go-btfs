package stats

import (
	"bytes"
	"context"
	"testing"

	unixtest "github.com/bittorrent/go-btfs/core/coreunix/test"

	nodepb "github.com/bittorrent/go-btfs-common/protos/node"

	"github.com/gogo/protobuf/proto"
)

func TestHostStatsSaveGet(t *testing.T) {
	node := unixtest.HelpTestMockRepo(t, nil)

	hs := &nodepb.StorageStat_Host{
		Online: true,
		// Uptime:      0.95,
		// Score:       0.88,
		StorageUsed: 100000,
		StorageCap:  10000000,
	}
	nodeId := node.Identity.Pretty()
	err := SaveHostStatsIntoDatastore(context.Background(), node, nodeId, hs)
	if err != nil {
		t.Fatal(err)
	}
	stored, err := GetHostStatsFromDatastore(context.Background(), node, nodeId)
	if err != nil {
		t.Fatal(err)
	}
	bs1, err := proto.Marshal(hs)
	if err != nil {
		t.Fatal(err)
	}
	bs2, err := proto.Marshal(stored)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(bs1, bs2) != 0 {
		t.Fatal("stored stats do not match saved stats")
	}
}
