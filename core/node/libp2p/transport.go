package libp2p

import (
	"fmt"

	config "github.com/TRON-US/go-btfs-config"

	libp2p "github.com/libp2p/go-libp2p"
	metrics "github.com/libp2p/go-libp2p/core/metrics"

	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	tcp "github.com/libp2p/go-libp2p/p2p/transport/tcp"
	websocket "github.com/libp2p/go-libp2p/p2p/transport/websocket"
	"go.uber.org/fx"
)

func Transports(tptConfig config.Transports) interface{} {
	return func(pnet struct {
		fx.In
		Fprint PNetFingerprint `optional:"true"`
	}) (opts Libp2pOpts, err error) {
		privateNetworkEnabled := pnet.Fprint != nil

		if tptConfig.Network.TCP.WithDefault(true) {
			// TODO(9290): Make WithMetrics configurable
			opts.Opts = append(opts.Opts, libp2p.Transport(tcp.NewTCPTransport, tcp.WithMetrics()))
		}

		if tptConfig.Network.Websocket.WithDefault(true) {
			opts.Opts = append(opts.Opts, libp2p.Transport(websocket.New))
		}

		if tptConfig.Network.QUIC.WithDefault(!privateNetworkEnabled) {
			if privateNetworkEnabled {
				// QUIC was force enabled while the private network was turned on.
				// Fail and tell the user.
				return opts, fmt.Errorf(
					"The QUIC transport does not support private networks. " +
						"Please disable Swarm.Transports.Network.QUIC.",
				)
			}
			// TODO(9290): Make WithMetrics configurable
			opts.Opts = append(opts.Opts, libp2p.Transport(quic.NewTransport, quic.WithMetrics()))
		}

		return opts, nil
	}
}

func BandwidthCounter() (opts Libp2pOpts, reporter *metrics.BandwidthCounter) {
	reporter = metrics.NewBandwidthCounter()
	opts.Opts = append(opts.Opts, libp2p.BandwidthReporter(reporter))
	return opts, reporter
}
