package coreapi

import (
	"context"
	"fmt"
	gopath "path"

	"github.com/bittorrent/go-btfs/namesys/resolve"

	uio "github.com/TRON-US/go-unixfs/io"
	coreiface "github.com/TRON-US/interface-go-btfs-core"
	path "github.com/TRON-US/interface-go-btfs-core/path"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	ipfspath "github.com/ipfs/go-path"
	"github.com/ipfs/go-path/resolver"
)

// ResolveNode resolves the path `p` using Unixfs resolver, gets and returns the
// resolved Node.
func (api *CoreAPI) ResolveNode(ctx context.Context, p path.Path) (ipld.Node, error) {
	rp, err := api.ResolvePath(ctx, p)
	if err != nil {
		return nil, err
	}

	node, err := api.dag.Get(ctx, rp.Cid())
	if err != nil {
		return nil, err
	}
	return node, nil
}

// ResolvePath resolves the path `p` using Unixfs resolver, returns the
// resolved path.
func (api *CoreAPI) ResolvePath(ctx context.Context, p path.Path) (path.Resolved, error) {
	if _, ok := p.(path.Resolved); ok {
		return p.(path.Resolved), nil
	}

	ipath, err := api.ResolveIpnsPath(ctx, p)
	if err != nil {
		return nil, err
	}

	var resolveOnce resolver.ResolveOnce

	switch ipath.Segments()[0] {
	case "btfs":
		resolveOnce = uio.ResolveUnixfsOnce
	case "ipld":
		resolveOnce = resolver.ResolveSingle
	default:
		return nil, fmt.Errorf("unsupported path namespace: %s", p.Namespace())
	}

	r := &resolver.Resolver{
		DAG:         api.dag,
		ResolveOnce: resolveOnce,
	}

	node, rest, err := r.ResolveToLastNode(ctx, *ipath)
	if err != nil {
		return nil, err
	}

	root, err := cid.Parse(ipath.Segments()[1])
	if err != nil {
		return nil, err
	}

	return path.NewResolvedPath(*ipath, node, root, gopath.Join(rest...)), nil
}

// ResolveIpnsPath resolves only the IPNS path `p` using the Unixfs resolver and returns the
// resolved path.
func (api *CoreAPI) ResolveIpnsPath(ctx context.Context, p path.Path) (*ipfspath.Path, error) {
	if err := p.IsValid(); err != nil {
		return nil, err
	}

	ipath := ipfspath.Path(p.String())
	ipath, err := resolve.ResolveIPNS(ctx, api.namesys, ipath)
	if err == resolve.ErrNoNamesys {
		return nil, coreiface.ErrOffline
	} else if err != nil {
		return nil, err
	}
	return &ipath, nil
}
