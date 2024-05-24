package node

import (
	"context"
	config "github.com/bittorrent/go-btfs-config"
	"github.com/bittorrent/go-btfs/peering"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"math/rand"
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

const (
	maxNLastConn = 10
	maxTryLimit  = 1000
)

// PeerWithLastConn try to connect to last peers
func PeerWithLastConn() fx.Option {
	return fx.Invoke(func(host host.Host, cfg *config.Config) {
		peerIds := host.Peerstore().PeersWithAddrs()

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
			if host.Network().Connectedness(id) == network.Connected || id == host.ID() || filter[id] {
				continue
			}
			connection[id] = true
		}

		g := errgroup.Group{}
		g.SetLimit(maxNLastConn)
		needConnect := int32(maxNLastConn)
		tryCount := 0

		for {
			if tryCount >= maxTryLimit {
				logger.Infof("max try count limited.")
				break
			}

			if needConnect <= 0 {
				break
			}

			randomSubSet := randomSubsetOfPeers(connection, int(needConnect))
			tryCount += len(randomSubSet)

			if len(randomSubSet) == 0 {
				break
			}

			for id, _ := range randomSubSet {
				connection[id] = false
				peerId := id
				g.Go(func() error {
					if err = host.Connect(context.Background(), host.Peerstore().PeerInfo(peerId)); err != nil {
						logger.Debugf("connect to last connection peer %s, error %v", peerId, err)
						return nil
					}
					atomic.AddInt32(&needConnect, -1)
					return nil
				})
			}
			err = g.Wait()
			if err != nil {
				logger.Debugf("connect to last connection error %v", err)
				return
			}
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
