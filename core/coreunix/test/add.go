package test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/coreapi"
	"github.com/bittorrent/go-btfs/core/coreunix"
	"github.com/bittorrent/go-btfs/repo"

	chunker "github.com/TRON-US/go-btfs-chunker"
	config "github.com/TRON-US/go-btfs-config"
	files "github.com/TRON-US/go-btfs-files"
	ftutil "github.com/TRON-US/go-unixfs/util"
	coreiface "github.com/TRON-US/interface-go-btfs-core"
	"github.com/TRON-US/interface-go-btfs-core/path"
	cid "github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	syncds "github.com/ipfs/go-datastore/sync"
)

const (
	testPeerID = "QmTFauExutTsy4XP6JbMFcw2Wa9645HJt2bTqL6qYDCKfe"

	TestRsDataSize   = 10
	TestRsParitySize = 20
)

// HelpTestMockRepo creates the bare minimum mock repo and returns node
func HelpTestMockRepo(t *testing.T, cfg *config.Config) *core.IpfsNode {
	mc := config.Config{}
	if cfg != nil {
		mc = *cfg
	}
	mc.Identity = config.Identity{
		PeerID: testPeerID, // required by offline node
	}
	r := &repo.Mock{
		C: mc,
		D: syncds.MutexWrap(datastore.NewMapDatastore()),
	}
	node, err := core.NewNode(context.Background(), &core.BuildCfg{Repo: r})
	if err != nil {
		t.Fatal(err)
	}
	return node
}

// HelpTestAddWithReedSolomomonMetadata is both a helper to testing this feature
// and also a helper to add a reed solomon file for other features.
// It returns a mock node, api, and added hash (cid).
func HelpTestAddWithReedSolomonMetadata(t *testing.T) (*core.IpfsNode, coreiface.CoreAPI, cid.Cid, []byte) {
	ctx := context.Background()
	node := HelpTestMockRepo(t, nil)

	out := make(chan interface{})
	adder, err := coreunix.NewAdder(ctx, node.Pinning, node.Blockstore, node.DAG)
	if err != nil {
		t.Fatal(err)
	}
	adder.Out = out
	// Set to default reed solomon for metadata
	dsize, psize, csize := TestRsDataSize, TestRsParitySize, 262144
	adder.Chunker = fmt.Sprintf("reed-solomon-%d-%d-%d", dsize, psize, csize)

	fb := []byte("testfileA")
	fsize := len(fb)
	rfa := files.NewBytesFile(fb)

	go func() {
		defer close(out)
		_, err := adder.AddAllAndPin(ctx, rfa)

		if err != nil {
			t.Fatal(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var addedHash cid.Cid
	select {
	case o := <-out:
		addedHash = o.(*coreiface.AddEvent).Path.Cid()
	case <-ctx.Done():
		t.Fatal("add timed out")
	}

	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		t.Fatal(err)
	}
	// Extract and check metadata
	b, err := coreunix.GetMetaData(ctx, api, path.IpfsPath(addedHash))
	if err != nil {
		t.Fatal(err)
	}
	var rsMeta chunker.RsMetaMap
	b1 := ftutil.GetMetadataElement(b)
	err = json.Unmarshal(b1, &rsMeta)
	if err != nil {
		t.Fatal(err)
	}
	if rsMeta.NumData != uint64(dsize) {
		t.Fatal("reed solomon metadata num data does not match")
	}
	if rsMeta.NumParity != uint64(psize) {
		t.Fatal("reed solomon metadata num parity does not match")
	}
	if rsMeta.FileSize != uint64(fsize) {
		t.Fatal("reed solomon metadata file size does not match")
	}

	return node, api, addedHash, fb
}
