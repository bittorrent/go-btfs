package libp2p

import (
	"github.com/libp2p/go-libp2p"
	// github.com/libp2p/go-libp2p-circuit has migrate to github.com/libp2p/go-libp2p/p2p/protocol/internal/circuitv1-deprecated,
	// and this internal package we can't access
	// relay "github.com/libp2p/go-libp2p-circuit"
)

func Relay(enableRelay, enableHop bool) func() (opts Libp2pOpts, err error) {
	return func() (opts Libp2pOpts, err error) {
		if enableRelay {
			// relayOpts := []relay.RelayOpt{}
			// if enableHop {
			// 	relayOpts = append(relayOpts, relay.OptHop)
			// }
			opts.Opts = append(opts.Opts, libp2p.EnableRelay())
		} else {
			opts.Opts = append(opts.Opts, libp2p.DisableRelay())
		}
		return
	}
}

var AutoRelay = simpleOpt(libp2p.ChainOptions(libp2p.EnableAutoRelay(), libp2p.DefaultStaticRelays()))
