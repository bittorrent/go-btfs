module github.com/ipfs/go-ipfs/examples/go-ipfs-as-a-library

go 1.14

require (
	github.com/bittorrent/go-btfs v0.0.0-20230626064024-58978cbfe949
	github.com/bittorrent/go-btfs-config v0.12.3
	github.com/bittorrent/go-btfs-files v0.3.1
	github.com/bittorrent/interface-go-btfs-core v0.8.2
	github.com/klauspost/cpuid v1.2.4 // indirect
	github.com/libp2p/go-libp2p v0.24.2
	github.com/multiformats/go-multiaddr v0.8.0
)

replace github.com/ipfs/go-ipfs => ./../../..
