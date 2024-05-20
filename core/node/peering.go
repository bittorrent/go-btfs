package node

import (
	"context"
	config "github.com/bittorrent/go-btfs-config"
	"github.com/bittorrent/go-btfs/peering"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"go.uber.org/fx"
	"math/rand"
	"sync"
	"sync/atomic"
)

// Peering constructs the peering service and hooks it into fx's lifetime
// management system.
func Peering(lc fx.Lifecycle, host host.Host) *peering.PeeringService {
	ps := peering.NewPeeringService(host)
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return ps.Start()
		},
		OnStop: func(context.Context) error {
			return ps.Stop()
		},
	})
	return ps
}

// PeerWith configures the peering service to peer with the specified peers.
func PeerWith(peers ...peer.AddrInfo) fx.Option {
	return fx.Invoke(func(ps *peering.PeeringService) {
		for _, ai := range peers {
			ps.AddPeer(ai)
		}
	})
}

const maxNLastConn = 10

// PeerWithLastConn try to connect to last peers
func PeerWithLastConn() fx.Option {
	return fx.Invoke(func(host host.Host, cfg *config.Config) {
		peerIds := host.Peerstore().Peers()

		bootstrap, err := cfg.BootstrapPeers()
		if err != nil {
			logger.Warn("failed to parse bootstrap peers from config")
		}

		filter := make(map[peer.ID]bool, len(bootstrap))
		for _, id := range bootstrap {
			filter[id.ID] = true
		}

		connection := make(map[peer.ID]bool)
		for _, id := range peerIds {
			if host.Network().Connectedness(id) != network.Connected || id != host.ID() || !filter[id] {
				connection[id] = true
			}
		}

		wg := sync.WaitGroup{}
		needConnect := int32(maxNLastConn)
		for {
			if needConnect <= 0 {
				break
			}
			randomSubSet := randomSubsetOfPeers(connection, int(needConnect))
			if len(randomSubSet) == 0 {
				break
			}
			for id, _ := range randomSubSet {
				connection[id] = false
				wg.Add(1)
				go func(peerId peer.ID) {
					defer wg.Done()
					if err = host.Connect(context.Background(), host.Peerstore().PeerInfo(peerId)); err != nil {
						logger.Debugf("connect to last connection peer %s, error %v", peerId, err)
						return
					}
					atomic.AddInt32(&needConnect, -1)
				}(id)
			}
			wg.Wait()
		}
	})
}

func randomSubsetOfPeers(in map[peer.ID]bool, max int) map[peer.ID]bool {
	c := 0
	for _, v := range in {
		if v {
			c++
		}
	}

	if max > c {
		max = c
	}

	out := make(map[peer.ID]bool, max)

	tem := make([]peer.ID, 0)
	for k, v := range in {
		if v {
			tem = append(tem, k)
		}
	}

	for _, val := range rand.Perm(max) {
		out[tem[val]] = true
	}

	return out
}
