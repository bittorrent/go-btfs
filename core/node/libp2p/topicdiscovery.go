package libp2p

import (
	"math/rand"
	"time"

	"github.com/libp2p/go-libp2p/core/discovery"
	"github.com/libp2p/go-libp2p/core/host"
	backoff "github.com/libp2p/go-libp2p/p2p/discovery/backoff"
	routing "github.com/libp2p/go-libp2p/p2p/discovery/routing"

	"github.com/bittorrent/go-btfs/core/node/helpers"
	"go.uber.org/fx"
)

func TopicDiscovery() interface{} {
	return func(mctx helpers.MetricsCtx, lc fx.Lifecycle, host host.Host, cr BaseIpfsRouting) (service discovery.Discovery, err error) {
		baseDisc := routing.NewRoutingDiscovery(cr)
		minBackoff, maxBackoff := time.Second*60, time.Hour
		rng := rand.New(rand.NewSource(rand.Int63()))
		d, err := backoff.NewBackoffDiscovery(
			baseDisc,
			backoff.NewExponentialBackoff(minBackoff, maxBackoff, backoff.FullJitter, time.Second, 5.0, 0, rng),
		)

		if err != nil {
			return nil, err
		}

		return d, nil
	}
}
