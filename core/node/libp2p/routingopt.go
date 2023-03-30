package libp2p

import (
	"context"

	config "github.com/TRON-US/go-btfs-config"
	irouting "github.com/bittorrent/go-btfs/routing"
	"github.com/ipfs/go-datastore"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dual "github.com/libp2p/go-libp2p-kad-dht/dual"
	record "github.com/libp2p/go-libp2p-record"
	routinghelpers "github.com/libp2p/go-libp2p-routing-helpers"
	host "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	routing "github.com/libp2p/go-libp2p/core/routing"
)

type RoutingOption func(
	context.Context,
	host.Host,
	datastore.Batching,
	record.Validator,
	...peer.AddrInfo,
) (routing.Routing, error)

func constructDHTRouting(mode dht.ModeOpt) func(
	ctx context.Context,
	host host.Host,
	dstore datastore.Batching,
	validator record.Validator,
	bootstrapPeers ...peer.AddrInfo,
) (routing.Routing, error) {
	return func(
		ctx context.Context,
		host host.Host,
		dstore datastore.Batching,
		validator record.Validator,
		bootstrapPeers ...peer.AddrInfo,
	) (routing.Routing, error) {
		return dual.New(
			ctx, host,
			dual.DHTOption(
				dht.Concurrency(10),
				dht.Mode(mode),
				dht.Datastore(dstore),
				// in case of "protocol prefix /ipfs must support the /ipns namespaced Validator" after upgrade libp2p version
				dht.ProtocolPrefix("/btfs"),
				dht.Validator(validator)),
			dual.WanDHTOption(dht.BootstrapPeers(bootstrapPeers...)),
		)
	}
}

func ConstructDelegatedRouting(routers config.Routers, methods config.Methods, peerID string, addrs []string, privKey string) func(
	ctx context.Context,
	host host.Host,
	dstore datastore.Batching,
	validator record.Validator,
	bootstrapPeers ...peer.AddrInfo,
) (routing.Routing, error) {
	return func(
		ctx context.Context,
		host host.Host,
		dstore datastore.Batching,
		validator record.Validator,
		bootstrapPeers ...peer.AddrInfo,
	) (routing.Routing, error) {
		return irouting.Parse(routers, methods,
			&irouting.ExtraDHTParams{
				BootstrapPeers: bootstrapPeers,
				Host:           host,
				Validator:      validator,
				Datastore:      dstore,
				Context:        ctx,
			},
			&irouting.ExtraReframeParams{
				PeerID:     peerID,
				Addrs:      addrs,
				PrivKeyB64: privKey,
			})
	}
}

func constructNilRouting(
	ctx context.Context,
	host host.Host,
	dstore datastore.Batching,
	validator record.Validator,
	bootstrapPeers ...peer.AddrInfo,
) (routing.Routing, error) {
	return routinghelpers.Null{}, nil
}

var (
	DHTOption       RoutingOption = constructDHTRouting(dht.ModeAuto)
	DHTClientOption               = constructDHTRouting(dht.ModeClient)
	DHTServerOption               = constructDHTRouting(dht.ModeServer)
	NilRouterOption               = constructNilRouting
)
