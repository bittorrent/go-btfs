package helper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bittorrent/go-btfs/core"

	chunker "github.com/bittorrent/go-btfs-chunker"
	"github.com/bittorrent/go-unixfs"
	coreiface "github.com/bittorrent/interface-go-btfs-core"
	"github.com/bittorrent/interface-go-btfs-core/path"
	cid "github.com/ipfs/go-cid"
)

// CheckAndGetReedSolomonShardHashes checks to see if a root hash is a reed solomon file,
// if ok, returns the list of shard hashes.
func CheckAndGetReedSolomonShardHashes(ctx context.Context, node *core.IpfsNode,
	api coreiface.CoreAPI, rootHash cid.Cid) ([]cid.Cid, int64, error) {
	rootPath := path.IpfsPath(rootHash)
	// check to see if a replicated file using reed-solomon
	mbytes, err := api.Unixfs().GetMetadata(ctx, rootPath)
	if err != nil {
		return nil, 0, fmt.Errorf("file must be reed-solomon encoded: %s", err.Error())
	}
	var rsMeta chunker.RsMetaMap
	err = json.Unmarshal(mbytes, &rsMeta)
	if err != nil {
		return nil, 0, fmt.Errorf("file must be reed-solomon encoded: %s", err.Error())
	}
	if rsMeta.NumData == 0 || rsMeta.NumParity == 0 || rsMeta.FileSize == 0 {
		return nil, 0, fmt.Errorf("file must be reed-solomon encoded: metadata not valid")
	}
	// use unixfs layer helper to grab the raw leaves under the data root node
	// higher level helpers resolve the enrtire DAG instead of resolving the root
	// node as it is required here
	rn, err := api.ResolveNode(ctx, rootPath)
	if err != nil {
		return nil, 0, err
	}
	nodes, err := unixfs.GetChildrenForDagWithMeta(ctx, rn, api.Dag())
	if err != nil {
		return nil, 0, err
	}
	links := nodes.DataNode.Links()
	if len(links) != int(rsMeta.NumData+rsMeta.NumParity) {
		return nil, 0, fmt.Errorf("file must be reed-solomon encoded: encoding scheme mismatch")
	}
	var hashes []cid.Cid
	for _, link := range links {
		hashes = append(hashes, link.Cid)
	}
	return hashes, int64(rsMeta.FileSize), nil
}
