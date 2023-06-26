package helper

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	unixtest "github.com/bittorrent/go-btfs/core/coreunix/test"

	uio "github.com/bittorrent/go-unixfs/io"
	"github.com/bittorrent/interface-go-btfs-core/path"
	rs "github.com/klauspost/reedsolomon"
)

func TestReedSolomonCheckGet(t *testing.T) {
	node, api, root, full := unixtest.HelpTestAddWithReedSolomonMetadata(t)
	hashes, _, err := CheckAndGetReedSolomonShardHashes(context.Background(), node, api, root)
	if err != nil {
		t.Fatal(err)
	}
	// check the full file
	enc, err := rs.New(int(unixtest.TestRsDataSize), int(unixtest.TestRsParitySize))
	if err != nil {
		t.Fatal("unable to create reference reedsolomon object", err)
	}
	shards, err := enc.Split(full)
	if err != nil {
		t.Fatal("unable to split reference reedsolomon shards", err)
	}
	err = enc.Encode(shards)
	if err != nil {
		t.Fatal("unable to encode reference reedsolomon shards", err)
	}
	// check if the hashes match the split shards
	for i, h := range hashes {
		hn, err := api.ResolveNode(context.Background(), path.IpfsPath(h))
		if err != nil {
			t.Fatal(err)
		}
		hr, err := uio.NewDagReader(context.Background(), hn, api.Dag())
		if err != nil {
			t.Fatal(err)
		}
		hash, err := ioutil.ReadAll(hr)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(hash, shards[i]) != 0 {
			t.Fatalf("shard [%d] does not match split hash contents", i)
		}
	}
}
