module github.com/ipfs/go-ipfs/examples/go-ipfs-as-a-library

go 1.14

require (
	github.com/TRON-US/go-btfs v1.5.0
	github.com/TRON-US/go-btfs-config v0.11.5
	github.com/TRON-US/go-btfs-files v0.2.0
	github.com/TRON-US/interface-go-btfs-core v0.7.0
	github.com/ipfs/go-ipfs v0.7.0
	github.com/ipfs/go-ipfs-config v0.9.0
	github.com/ipfs/go-ipfs-files v0.0.8
	github.com/TRON-US/interface-go-btfs-core v0.4.0
	github.com/libp2p/go-libp2p-core v0.9.0
	github.com/libp2p/go-libp2p-peerstore v0.2.7
	github.com/multiformats/go-multiaddr v0.3.3
)

replace github.com/ipfs/go-ipfs => ./../../..
