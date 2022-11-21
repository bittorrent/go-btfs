package remote

import (
	"context"

	"github.com/bittorrent/go-btfs/core"

	cmds "github.com/bittorrent/go-btfs-cmds"
	cmdsHttp "github.com/bittorrent/go-btfs-cmds/http"

	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

// GetStreamRequestRemotePeerID checks to see if current request is part of a streamedd
// libp2p connection, if yes, return the remote peer id, otherwise return false.
func GetStreamRequestRemotePeerID(req *cmds.Request, node *core.IpfsNode) (peer.ID, bool) {
	remoteAddr, ok := cmdsHttp.GetRequestRemoteAddr(req.Context)
	if !ok {
		return "", false
	}
	return node.P2P.Streams.GetStreamRemotePeerID(remoteAddr)
}

// FindPeer decodes a string-based peer id and tries to find it in the current routing
// table (if not connected, will retry).
func FindPeer(ctx context.Context, n *core.IpfsNode, pid string) (*peer.AddrInfo, error) {
	id, err := peer.Decode(pid)
	if err != nil {
		return nil, err
	}
	pinfo, err := n.Routing.FindPeer(ctx, id)
	if err != nil {
		return nil, err
	}
	addr, err := ma.NewMultiaddr("/p2p-circuit/btfs/" + pid)
	if err != nil {
		return nil, err
	}
	pinfo.Addrs = append(pinfo.Addrs, addr)
	return &pinfo, nil
}
