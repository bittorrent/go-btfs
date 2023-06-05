package challenge

import (
	"context"
	"testing"

	unixtest "github.com/bittorrent/go-btfs/core/coreunix/test"

	unixfs "github.com/bittorrent/go-unixfs"
	path "github.com/bittorrent/interface-go-btfs-core/path"
)

func TestGenAndSolveChallenge(t *testing.T) {
	node, api, root, _ := unixtest.HelpTestAddWithReedSolomonMetadata(t)

	// Resolve file root node
	rn, err := api.ResolveNode(context.Background(), path.IpfsPath(root))
	if err != nil {
		t.Fatal(err)
	}

	// Get sharded root node (strip away metadata)
	nodes, err := unixfs.GetChildrenForDagWithMeta(context.Background(), rn, api.Dag())
	if err != nil {
		t.Fatal(err)
	}

	// Check all shards under reed solomon
	for _, link := range nodes.DataNode.Links() {
		sid := link.Cid
		// simulate client challenge gen
		sc, err := NewStorageChallenge(context.Background(), node, api, root, sid)
		if err != nil {
			t.Fatal(err)
		}
		// simulate host challenge preparation
		scr, err := NewStorageChallengeResponse(context.Background(), node, api, root, sid, sc.ID, false, 0)
		if err != nil {
			t.Fatal(err)
		}

		// simulate 100 different challenges per shard
		chs := map[string]bool{}
		for i := 0; i < 100; i++ {
			err := sc.GenChallenge()
			if err != nil {
				t.Fatal(err)
			}

			// should not re-generate duplicate challenges
			if _, ok := chs[sc.Hash]; ok {
				t.Fatal("Duplicate challenge generated")
			}
			chs[sc.Hash] = true

			err = scr.SolveChallenge(sc.CIndex, sc.Nonce)
			if err != nil {
				t.Fatal(err)
			}

			if scr.Hash != sc.Hash {
				t.Fatal("Challenge is not solved, not matching original hash")
			}
		}
	}
}
