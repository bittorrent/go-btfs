package namesys

import (
	"context"
	"fmt"
	"testing"
	"time"

	btns "github.com/bittorrent/go-btns"
	unixfs "github.com/bittorrent/go-unixfs"
	opts "github.com/bittorrent/interface-go-btfs-core/options/namesys"
	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	offroute "github.com/ipfs/go-ipfs-routing/offline"
	path "github.com/ipfs/go-path"
	record "github.com/libp2p/go-libp2p-record"
	ci "github.com/libp2p/go-libp2p/core/crypto"
	peer "github.com/libp2p/go-libp2p/core/peer"
	pstoremem "github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoremem"
)

type mockResolver struct {
	entries map[string]string
}

func testResolution(t *testing.T, resolver Resolver, name string, depth uint, expected string, expError error) {
	t.Helper()
	p, err := resolver.Resolve(context.Background(), name, opts.Depth(depth))
	if err != expError {
		t.Fatal(fmt.Errorf(
			"expected %s with a depth of %d to have a '%s' error, but got '%s'",
			name, depth, expError, err))
	}
	if p.String() != expected {
		t.Fatal(fmt.Errorf(
			"%s with depth %d resolved to %s != %s",
			name, depth, p.String(), expected))
	}
}

func (r *mockResolver) resolveOnceAsync(ctx context.Context, name string, options opts.ResolveOpts) <-chan onceResult {
	p, err := path.ParsePath(r.entries[name])
	out := make(chan onceResult, 1)
	out <- onceResult{value: p, err: err}
	close(out)
	return out
}

func mockResolverOne() *mockResolver {
	return &mockResolver{
		entries: map[string]string{
			"QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy":              "/btfs/Qmcqtw8FfrVSBaRmbWwHxt3AuySBhJLcvmFYi3Lbc4xnwj",
			"QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n":              "/btns/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy",
			"QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD":              "/btns/ipfs.io",
			"QmQ4QZh8nrsczdUEwTyfBope4THUhqxqc1fx6qYhhzZQei":              "/btfs/QmP3ouCnU8NNLsW6261pAx2pNLV2E4dQoisB1sgda12Act",
			"12D3KooWFB51PRY9BxcXSH6khFXw1BZeszeLDy7C8GciskqCTZn5":        "/btns/QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n", // ed25519+identity multihash
			"bafzbeickencdqw37dpz3ha36ewrh4undfjt2do52chtcky4rxkj447qhdm": "/btns/QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n", // cidv1 in base32 with libp2p-key multicodec
		},
	}
}

func mockResolverTwo() *mockResolver {
	return &mockResolver{
		entries: map[string]string{
			"ipfs.io": "/btns/QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n",
		},
	}
}

func TestNamesysResolution(t *testing.T) {
	r := &mpns{
		ipnsResolver: mockResolverOne(),
		dnsResolver:  mockResolverTwo(),
	}

	testResolution(t, r, "Qmcqtw8FfrVSBaRmbWwHxt3AuySBhJLcvmFYi3Lbc4xnwj", opts.DefaultDepthLimit, "/btfs/Qmcqtw8FfrVSBaRmbWwHxt3AuySBhJLcvmFYi3Lbc4xnwj", nil)
	testResolution(t, r, "/btns/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy", opts.DefaultDepthLimit, "/btfs/Qmcqtw8FfrVSBaRmbWwHxt3AuySBhJLcvmFYi3Lbc4xnwj", nil)
	testResolution(t, r, "/btns/QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n", opts.DefaultDepthLimit, "/btfs/Qmcqtw8FfrVSBaRmbWwHxt3AuySBhJLcvmFYi3Lbc4xnwj", nil)
	testResolution(t, r, "/btns/QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n", 1, "/btns/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy", ErrResolveRecursion)
	testResolution(t, r, "/btns/ipfs.io", opts.DefaultDepthLimit, "/btfs/Qmcqtw8FfrVSBaRmbWwHxt3AuySBhJLcvmFYi3Lbc4xnwj", nil)
	testResolution(t, r, "/btns/ipfs.io", 1, "/btns/QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n", ErrResolveRecursion)
	testResolution(t, r, "/btns/ipfs.io", 2, "/btns/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy", ErrResolveRecursion)
	testResolution(t, r, "/btns/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", opts.DefaultDepthLimit, "/btfs/Qmcqtw8FfrVSBaRmbWwHxt3AuySBhJLcvmFYi3Lbc4xnwj", nil)
	testResolution(t, r, "/btns/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", 1, "/btns/ipfs.io", ErrResolveRecursion)
	testResolution(t, r, "/btns/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", 2, "/btns/QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n", ErrResolveRecursion)
	testResolution(t, r, "/btns/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", 3, "/btns/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy", ErrResolveRecursion)
	testResolution(t, r, "/btns/12D3KooWFB51PRY9BxcXSH6khFXw1BZeszeLDy7C8GciskqCTZn5", 1, "/btns/QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n", ErrResolveRecursion)
	testResolution(t, r, "/btns/bafzbeickencdqw37dpz3ha36ewrh4undfjt2do52chtcky4rxkj447qhdm", 1, "/btns/QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n", ErrResolveRecursion)
}

func TestPublishWithCache0(t *testing.T) {
	dst := dssync.MutexWrap(ds.NewMapDatastore())
	priv, _, err := ci.GenerateKeyPair(ci.RSA, 2048)
	if err != nil {
		t.Fatal(err)
	}
	ps, err := pstoremem.NewPeerstore()
	if err != nil {
		t.Fatal(err)
	}
	pid, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	err = ps.AddPrivKey(pid, priv)
	if err != nil {
		t.Fatal(err)
	}

	routing := offroute.NewOfflineRouter(dst, record.NamespacedValidator{
		"btns": btns.Validator{KeyBook: ps},
		"pk":   record.PublicKeyValidator{},
	})

	nsys, err := NewNameSystem(routing, WithDatastore(dst))
	if err != nil {
		t.Fatal(err)
	}
	p, err := path.ParsePath(unixfs.EmptyDirNode().Cid().String())
	if err != nil {
		t.Fatal(err)
	}
	err = nsys.Publish(context.Background(), priv, p)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPublishWithTTL(t *testing.T) {
	dst := dssync.MutexWrap(ds.NewMapDatastore())
	priv, _, err := ci.GenerateKeyPair(ci.RSA, 2048)
	if err != nil {
		t.Fatal(err)
	}
	ps, err := pstoremem.NewPeerstore()
	if err != nil {
		t.Fatal(err)
	}
	pid, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	err = ps.AddPrivKey(pid, priv)
	if err != nil {
		t.Fatal(err)
	}

	routing := offroute.NewOfflineRouter(dst, record.NamespacedValidator{
		"btns": btns.Validator{KeyBook: ps},
		"pk":   record.PublicKeyValidator{},
	})

	nsys, err := NewNameSystem(routing, WithDatastore(dst), WithCache(128))
	if err != nil {
		t.Fatal(err)
	}

	p, err := path.ParsePath(unixfs.EmptyDirNode().Cid().String())
	if err != nil {
		t.Fatal(err)
	}

	ttl := 1 * time.Second
	eol := time.Now().Add(2 * time.Second)

	ctx := context.WithValue(context.Background(), "btns-publish-ttl", ttl)
	err = nsys.Publish(ctx, priv, p)
	if err != nil {
		t.Fatal(err)
	}
	ientry, ok := nsys.(*mpns).cache.Get(string(pid))
	if !ok {
		t.Fatal("cache get failed")
	}
	entry, ok := ientry.(cacheEntry)
	if !ok {
		t.Fatal("bad cache item returned")
	}
	if entry.eol.Sub(eol) > 10*time.Millisecond {
		t.Fatalf("bad cache ttl: expected %s, got %s", eol, entry.eol)
	}
}
