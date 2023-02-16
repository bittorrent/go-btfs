package libp2p

import (
	"context"
	"time"

	"github.com/bittorrent/go-btfs/core/node/helpers"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"go.uber.org/fx"
)

const discoveryConnTimeout = time.Second * 30

type discoveryHandler struct {
	ctx  context.Context
	host host.Host
}

func (dh *discoveryHandler) HandlePeerFound(p peer.AddrInfo) {
	log.Info("connecting to discovered peer: ", p)
	ctx, cancel := context.WithTimeout(dh.ctx, discoveryConnTimeout)
	defer cancel()
	if err := dh.host.Connect(ctx, p); err != nil {
		log.Warnf("failed to connect to peer %s found by discovery: %s", p.ID, err)
	}
}

func DiscoveryHandler(mctx helpers.MetricsCtx, lc fx.Lifecycle, host host.Host) *discoveryHandler {
	return &discoveryHandler{
		ctx:  helpers.LifecycleCtx(mctx, lc),
		host: host,
	}
}

func SetupDiscovery(mdns bool, mdnsInterval int) func(helpers.MetricsCtx, fx.Lifecycle, host.Host, *discoveryHandler) error {
	return func(mctx helpers.MetricsCtx, lc fx.Lifecycle, host host.Host, handler *discoveryHandler) error {
		if mdns {
			if mdnsInterval == 0 {
				mdnsInterval = 5
			}
			service := discovery.NewMdnsService(host, discovery.ServiceName, handler)
			if err := service.Start(); err != nil {
				log.Error("error starting mdns service: ", err)
				return nil
			}
			// if err != nil {
			// 	log.Error("mdns error: ", err)
			// 	return nil
			// }
			// service.RegisterNotifee(handler)
		}
		return nil
	}
}
