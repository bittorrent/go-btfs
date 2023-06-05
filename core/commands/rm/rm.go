package rm

import (
	"context"
	"fmt"

	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"

	cmds "github.com/bittorrent/go-btfs-cmds"
	coreiface "github.com/bittorrent/interface-go-btfs-core"
	"github.com/bittorrent/interface-go-btfs-core/options"
	"github.com/bittorrent/interface-go-btfs-core/path"
	ipld "github.com/ipfs/go-ipld-format"
)

func RmDag(ctx context.Context, hashes []string, n *core.IpfsNode, req *cmds.Request, env cmds.Environment,
	force bool) ([]string, error) {
	// do not perform online operations for local delete
	api, err := cmdenv.GetApi(env, req, options.Api.Offline(true), options.Api.FetchBlocks(false))
	if err != nil {
		return nil, err
	}

	var results []string
	for _, b := range hashes {
		// Make sure node exists
		p := path.New(b)
		node, err := api.ResolveNode(ctx, p)
		if err != nil {
			results = append(results, fmt.Sprintf("Error resolving root %s: %v", b, err))
			continue
		}

		_, pinned, err := n.Pinning.IsPinned(ctx, node.Cid())
		if err != nil {
			return nil, err
		}

		if pinned {
			// Since we are removing a file, we need to set recursive flag to true
			err = api.Pin().Rm(ctx, p, options.Pin.RmRecursive(true), options.Pin.RmForce(force))
			if err != nil {
				results = append(results, fmt.Sprintf("Error removing root %s pin: %v", b, err))
				continue
			}
		}

		// Rm all child links
		err = rmAllDags(ctx, api, node, &results, map[string]bool{})
		if err != nil {
			results = append(results, fmt.Sprintf("Error removing root %s child objects: %v", b, err))
			continue
		}
	}

	return results, nil
}

func rmAllDags(ctx context.Context, api coreiface.CoreAPI, node ipld.Node, res *[]string,
	removed map[string]bool) error {
	for _, nl := range node.Links() {
		// Skipped just removed nodes
		if _, ok := removed[nl.Cid.String()]; ok {
			continue
		}
		// Resolve, recurse, then finally remove
		rn, err := api.ResolveNode(ctx, path.IpfsPath(nl.Cid))
		if err != nil {
			*res = append(*res, fmt.Sprintf("Error resolving object %s: %v", nl.Cid, err))
			continue // continue with other nodes
		}
		if err := rmAllDags(ctx, api, rn, res, removed); err != nil {
			continue // continue with other nodes
		}
	}
	ncid := node.Cid()
	if err := api.Dag().Remove(ctx, ncid); err != nil {
		*res = append(*res, fmt.Sprintf("Error removing object %s: %v", ncid, err))
		return err
	} else {
		*res = append(*res, fmt.Sprintf("Removed %s", ncid))
		removed[ncid.String()] = true
	}
	return nil
}
