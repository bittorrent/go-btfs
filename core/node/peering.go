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
	"time"
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
	maxNLastConn    = 10
	maxTryLimit     = 100
	maxTimeDuration = 20 * time.Second
	connTimeout     = 3 * time.Second
)

func loadConnPeers(host host.Host, cfg *config.Config) map[peer.ID]bool {
	peerIds := host.Peerstore().PeersWithAddrs()

	bootstrap, err := cfg.BootstrapPeers()
	if err != nil {
		logger.Warn("failed to parse bootstrap peers from config")
	}

	filter := make(map[peer.ID]bool, len(bootstrap))
	for _, id := range bootstrap {
		filter[id.ID] = true
	}

	canConnect := make(map[peer.ID]bool)
	for _, id := range peerIds {
		if host.Network().Connectedness(id) == network.Connected || id == host.ID() || filter[id] {
			continue
		}
		canConnect[id] = true
	}
	return canConnect
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

func clearTriedPeer(peers map[peer.ID]bool, triedPeer map[peer.ID]bool) map[peer.ID]bool {
	for k := range triedPeer {
		peers[k] = false
	}
	return peers
}

func doConcurrentConn(ctx context.Context, host host.Host, peers map[peer.ID]bool) int32 {
	if len(peers) < 1 {
		return 0
	}

	g := errgroup.Group{}
	g.SetLimit(len(peers))

	connected := int32(0)

	for id := range peers {
		peerId := id
		g.Go(func() error {
			if err := host.Connect(ctx, host.Peerstore().PeerInfo(peerId)); err != nil {
				logger.Debugf("connect to last connection peer %s, error %v", peerId, err)
				return nil
			}
			atomic.AddInt32(&connected, 1)
			return nil
		})
	}
	_ = g.Wait()

	return connected
}

func tryConn(host host.Host, peers map[peer.ID]bool) {

	success := make(chan struct{})
	useOut := make(chan struct{})
	maxTry := make(chan struct{})
	timeout := make(chan struct{})

	timer := time.NewTimer(maxTimeDuration)

	needPeerCount := maxNLastConn
	canTryPeerCount := maxTryLimit

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, connTimeout)
	defer cancel()

	go func() {
		for {
			select {
			case <-timer.C:
				timeout <- struct{}{}
				return
			default:
				conns := randomSubsetOfPeers(peers, needPeerCount)
				peers = clearTriedPeer(peers, conns)

				connCount := doConcurrentConn(ctx, host, conns)

				needPeerCount -= int(connCount)
				canTryPeerCount -= len(conns)

				if len(conns) <= 0 {
					useOut <- struct{}{}
					return
				}

				if needPeerCount <= 0 {
					success <- struct{}{}
					return
				}

				if canTryPeerCount <= 0 {
					maxTry <- struct{}{}
					return
				}
			}
		}
	}()

	select {
	case <-timeout:
		logger.Debugf("connect to last connection timeout")
		return
	case <-success:
		logger.Debugf("connect to last connection success")
		return
	case <-useOut:
		logger.Debugf("connect to last connection use out")
		return
	case <-maxTry:
		logger.Debugf("connect to last connection try limited")
		return
	}
}

// PeerWithLastConn tryConn to connect to last peers
func PeerWithLastConn() fx.Option {
	return fx.Invoke(func(host host.Host, cfg *config.Config) {
		peers := loadConnPeers(host, cfg)
		tryConn(host, peers)
	})
}
