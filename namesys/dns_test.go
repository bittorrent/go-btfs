package namesys

import (
	"context"
	"fmt"
	"testing"

	opts "github.com/bittorrent/interface-go-btfs-core/options/namesys"
)

type mockDNS struct {
	entries map[string][]string
}

func (m *mockDNS) lookupTXT(ctx context.Context, name string) (txt []string, err error) {
	txt, ok := m.entries[name]
	if !ok {
		return nil, fmt.Errorf("no TXT entry for %s", name)
	}
	return txt, nil
}

func TestDnsEntryParsing(t *testing.T) {

	goodEntries := []string{
		"QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
		"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
		"dnslink=/btns/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
		"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/foo",
		"dnslink=/btns/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/bar",
		"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/foo/bar/baz",
		"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/foo/bar/baz/",
		"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
	}

	badEntries := []string{
		"QmYhE8xgFCjGcz6PHgnvJz5NOTCORRECT",
		"quux=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
		"dnslink=",
		"dnslink=/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/foo",
		"dnslink=btns/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/bar",
	}

	for _, e := range goodEntries {
		_, err := parseEntry(e)
		if err != nil {
			t.Log("expected entry to parse correctly!")
			t.Log(e)
			t.Fatal(err)
		}
	}

	for _, e := range badEntries {
		_, err := parseEntry(e)
		if err == nil {
			t.Log("expected entry parse to fail!")
			t.Fatal(err)
		}
	}
}

func newMockDNS() *mockDNS {
	return &mockDNS{
		entries: map[string][]string{
			"multihash.example.com.": []string{
				"dnslink=QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
			},
			"ipfs.example.com.": []string{
				"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
			},
			"_dnslink.dipfs.example.com.": []string{
				"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
			},
			"dns1.example.com.": []string{
				"dnslink=/btns/ipfs.example.com",
			},
			"dns2.example.com.": []string{
				"dnslink=/btns/dns1.example.com",
			},
			"multi.example.com.": []string{
				"some stuff",
				"dnslink=/btns/dns1.example.com",
				"masked dnslink=/btns/example.invalid",
			},
			"equals.example.com.": []string{
				"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/=equals",
			},
			"loop1.example.com.": []string{
				"dnslink=/btns/loop2.example.com",
			},
			"loop2.example.com.": []string{
				"dnslink=/btns/loop1.example.com",
			},
			"_dnslink.dloop1.example.com.": []string{
				"dnslink=/btns/loop2.example.com",
			},
			"_dnslink.dloop2.example.com.": []string{
				"dnslink=/btns/loop1.example.com",
			},
			"bad.example.com.": []string{
				"dnslink=",
			},
			"withsegment.example.com.": []string{
				"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/sub/segment",
			},
			"withrecsegment.example.com.": []string{
				"dnslink=/btns/withsegment.example.com/subsub",
			},
			"withtrailing.example.com.": []string{
				"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/sub/",
			},
			"withtrailingrec.example.com.": []string{
				"dnslink=/btns/withtrailing.example.com/segment/",
			},
			"double.example.com.": []string{
				"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
			},
			"_dnslink.double.example.com.": []string{
				"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
			},
			"double.conflict.com.": []string{
				"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD",
			},
			"_dnslink.conflict.example.com.": []string{
				"dnslink=/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjE",
			},
			"fqdn.example.com.": []string{
				"dnslink=/btfs/QmYvMB9yrsSf7RKBghkfwmHJkzJhW2ZgVwq3LxBXXPasFr",
			},
			"www.wealdtech.eth.link.": []string{
				"dnslink=/ipns/ipfs.example.com",
			},
		},
	}
}

func TestDNSResolution(t *testing.T) {
	mock := newMockDNS()
	r := &DNSResolver{lookupTXT: mock.lookupTXT}
	testResolution(t, r, "multihash.example.com", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", nil)
	testResolution(t, r, "ipfs.example.com", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", nil)
	testResolution(t, r, "dipfs.example.com", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", nil)
	testResolution(t, r, "dns1.example.com", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", nil)
	testResolution(t, r, "dns1.example.com", 1, "/btns/ipfs.example.com", ErrResolveRecursion)
	testResolution(t, r, "dns2.example.com", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", nil)
	testResolution(t, r, "dns2.example.com", 1, "/btns/dns1.example.com", ErrResolveRecursion)
	testResolution(t, r, "dns2.example.com", 2, "/btns/ipfs.example.com", ErrResolveRecursion)
	testResolution(t, r, "multi.example.com", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", nil)
	testResolution(t, r, "multi.example.com", 1, "/btns/dns1.example.com", ErrResolveRecursion)
	testResolution(t, r, "multi.example.com", 2, "/btns/ipfs.example.com", ErrResolveRecursion)
	testResolution(t, r, "equals.example.com", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/=equals", nil)
	testResolution(t, r, "loop1.example.com", 1, "/btns/loop2.example.com", ErrResolveRecursion)
	testResolution(t, r, "loop1.example.com", 2, "/btns/loop1.example.com", ErrResolveRecursion)
	testResolution(t, r, "loop1.example.com", 3, "/btns/loop2.example.com", ErrResolveRecursion)
	testResolution(t, r, "loop1.example.com", opts.DefaultDepthLimit, "/btns/loop1.example.com", ErrResolveRecursion)
	testResolution(t, r, "dloop1.example.com", 1, "/btns/loop2.example.com", ErrResolveRecursion)
	testResolution(t, r, "dloop1.example.com", 2, "/btns/loop1.example.com", ErrResolveRecursion)
	testResolution(t, r, "dloop1.example.com", 3, "/btns/loop2.example.com", ErrResolveRecursion)
	testResolution(t, r, "dloop1.example.com", opts.DefaultDepthLimit, "/btns/loop1.example.com", ErrResolveRecursion)
	testResolution(t, r, "bad.example.com", opts.DefaultDepthLimit, "", ErrResolveFailed)
	testResolution(t, r, "withsegment.example.com", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/sub/segment", nil)
	testResolution(t, r, "withrecsegment.example.com", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/sub/segment/subsub", nil)
	testResolution(t, r, "withsegment.example.com/test1", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/sub/segment/test1", nil)
	testResolution(t, r, "withrecsegment.example.com/test2", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/sub/segment/subsub/test2", nil)
	testResolution(t, r, "withrecsegment.example.com/test3/", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/sub/segment/subsub/test3/", nil)
	testResolution(t, r, "withtrailingrec.example.com", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD/sub/segment/", nil)
	testResolution(t, r, "double.example.com", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", nil)
	testResolution(t, r, "conflict.example.com", opts.DefaultDepthLimit, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjE", nil)
	testResolution(t, r, "fqdn.example.com.", opts.DefaultDepthLimit, "/btfs/QmYvMB9yrsSf7RKBghkfwmHJkzJhW2ZgVwq3LxBXXPasFr", nil)
	// TODO: add btfs specific domain tests
	//testResolution(t, r, "www.wealdtech.eth", 1, "/btns/btfs.example.com", ErrResolveFailed)
	//testResolution(t, r, "www.wealdtech.eth", 2, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", nil)
	//testResolution(t, r, "www.wealdtech.eth.link", 2, "/btfs/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", nil)
}
