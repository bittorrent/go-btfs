package gateway

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	cid "github.com/ipfs/go-cid"
	path "github.com/ipfs/go-path"
	"github.com/stretchr/testify/assert"
)

func TestToSubdomainURL(t *testing.T) {
	gwAPI, _ := newMockAPI(t)
	testCID, err := cid.Decode("bafkqaglimvwgy3zakrsxg5cun5jxkyten5wwc2lokvjeycq")
	assert.Nil(t, err)

	gwAPI.namesys["/btns/dnslink.long-name.example.com"] = path.FromString(testCID.String())
	gwAPI.namesys["/btns/dnslink.too-long.f1siqrebi3vir8sab33hu5vcy008djegvay6atmz91ojesyjs8lx350b7y7i1nvyw2haytfukfyu2f2x4tocdrfa0zgij6p4zpl4u5o.example.com"] = path.FromString(testCID.String())
	httpRequest := httptest.NewRequest("GET", "http://127.0.0.1:8080", nil)
	httpsRequest := httptest.NewRequest("GET", "https://https-request-stub.example.com", nil)
	httpsProxiedRequest := httptest.NewRequest("GET", "http://proxied-https-request-stub.example.com", nil)
	httpsProxiedRequest.Header.Set("X-Forwarded-Proto", "https")

	for _, test := range []struct {
		// in:
		request       *http.Request
		gwHostname    string
		inlineDNSLink bool
		path          string
		// out:
		url string
		err error
	}{

		// DNSLink
		{httpRequest, "localhost", false, "/btns/dnslink.io", "http://dnslink.io.btns.localhost/", nil},
		// Hostname with port
		{httpRequest, "localhost:8080", false, "/btns/dnslink.io", "http://dnslink.io.btns.localhost:8080/", nil},
		// CIDv0 → CIDv1base32
		{httpRequest, "localhost", false, "/btfs/QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n", "http://bafybeif7a7gdklt6hodwdrmwmxnhksctcuav6lfxlcyfz4khzl3qfmvcgu.btfs.localhost/", nil},
		// CIDv1 with long sha512
		{httpRequest, "localhost", false, "/btfs/bafkrgqe3ohjcjplc6n4f3fwunlj6upltggn7xqujbsvnvyw764srszz4u4rshq6ztos4chl4plgg4ffyyxnayrtdi5oc4xb2332g645433aeg", "", errors.New("CID incompatible with DNS label length limit of 63: kf1siqrebi3vir8sab33hu5vcy008djegvay6atmz91ojesyjs8lx350b7y7i1nvyw2haytfukfyu2f2x4tocdrfa0zgij6p4zpl4u5oj")},
		// PeerID as CIDv1 needs to have libp2p-key multicodec
		{httpRequest, "localhost", false, "/btns/QmY3hE8xgFCjGcz6PHgnvJz5HZi1BaKRfPkn1ghZUcYMjD", "http://k2k4r8n0flx3ra0y5dr8fmyvwbzy3eiztmtq6th694k5a3rznayp3e4o.btns.localhost/", nil},
		{httpRequest, "localhost", false, "/btns/bafybeickencdqw37dpz3ha36ewrh4undfjt2do52chtcky4rxkj447qhdm", "http://k2k4r8l9ja7hkzynavdqup76ou46tnvuaqegbd04a4o1mpbsey0meucb.btns.localhost/", nil},
		// PeerID: ed25519+identity multihash → CIDv1Base36
		{httpRequest, "localhost", false, "/btns/12D3KooWFB51PRY9BxcXSH6khFXw1BZeszeLDy7C8GciskqCTZn5", "http://k51qzi5uqu5di608geewp3nqkg0bpujoasmka7ftkyxgcm3fh1aroup0gsdrna.btns.localhost/", nil},
		{httpRequest, "sub.localhost", false, "/btfs/QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n", "http://bafybeif7a7gdklt6hodwdrmwmxnhksctcuav6lfxlcyfz4khzl3qfmvcgu.btfs.sub.localhost/", nil},
		// HTTPS requires DNSLink name to fit in a single DNS label – see "Option C" from https://github.com/btfs/in-web-browsers/issues/169
		{httpRequest, "dweb.link", false, "/btns/dnslink.long-name.example.com", "http://dnslink.long-name.example.com.btns.dweb.link/", nil},
		{httpsRequest, "dweb.link", false, "/btns/dnslink.long-name.example.com", "https://dnslink-long--name-example-com.btns.dweb.link/", nil},
		{httpsProxiedRequest, "dweb.link", false, "/btns/dnslink.long-name.example.com", "https://dnslink-long--name-example-com.btns.dweb.link/", nil},
		// HTTP requests can also be converted to fit into a single DNS label - https://github.com/btfs/kubo/issues/9243
		{httpRequest, "localhost", true, "/btns/dnslink.long-name.example.com", "http://dnslink-long--name-example-com.btns.localhost/", nil},
		{httpRequest, "dweb.link", true, "/btns/dnslink.long-name.example.com", "http://dnslink-long--name-example-com.btns.dweb.link/", nil},
	} {
		testName := fmt.Sprintf("%s, %v, %s", test.gwHostname, test.inlineDNSLink, test.path)
		t.Run(testName, func(t *testing.T) {
			url, err := toSubdomainURL(test.gwHostname, test.path, test.request, test.inlineDNSLink, gwAPI)
			assert.Equal(t, test.url, url)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestToDNSLinkDNSLabel(t *testing.T) {
	for _, test := range []struct {
		in  string
		out string
		err error
	}{
		{"dnslink.long-name.example.com", "dnslink-long--name-example-com", nil},
		{"dnslink.too-long.f1siqrebi3vir8sab33hu5vcy008djegvay6atmz91ojesyjs8lx350b7y7i1nvyw2haytfukfyu2f2x4tocdrfa0zgij6p4zpl4u5o.example.com", "", errors.New("DNSLink representation incompatible with DNS label length limit of 63: dnslink-too--long-f1siqrebi3vir8sab33hu5vcy008djegvay6atmz91ojesyjs8lx350b7y7i1nvyw2haytfukfyu2f2x4tocdrfa0zgij6p4zpl4u5o-example-com")},
	} {
		t.Run(test.in, func(t *testing.T) {
			out, err := toDNSLinkDNSLabel(test.in)
			assert.Equal(t, test.out, out)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestToDNSLinkFQDN(t *testing.T) {
	for _, test := range []struct {
		in  string
		out string
	}{
		{"singlelabel", "singlelabel"},
		{"docs-btfs-tech", "docs.btfs.tech"},
		{"dnslink-long--name-example-com", "dnslink.long-name.example.com"},
	} {
		t.Run(test.in, func(t *testing.T) {
			out := toDNSLinkFQDN(test.in)
			assert.Equal(t, test.out, out)
		})
	}
}

func TestIsHTTPSRequest(t *testing.T) {
	httpRequest := httptest.NewRequest("GET", "http://127.0.0.1:8080", nil)
	httpsRequest := httptest.NewRequest("GET", "https://https-request-stub.example.com", nil)
	httpsProxiedRequest := httptest.NewRequest("GET", "http://proxied-https-request-stub.example.com", nil)
	httpsProxiedRequest.Header.Set("X-Forwarded-Proto", "https")
	httpProxiedRequest := httptest.NewRequest("GET", "http://proxied-http-request-stub.example.com", nil)
	httpProxiedRequest.Header.Set("X-Forwarded-Proto", "http")
	oddballRequest := httptest.NewRequest("GET", "foo://127.0.0.1:8080", nil)
	for _, test := range []struct {
		in  *http.Request
		out bool
	}{
		{httpRequest, false},
		{httpsRequest, true},
		{httpsProxiedRequest, true},
		{httpProxiedRequest, false},
		{oddballRequest, false},
	} {
		testName := fmt.Sprintf("%+v", test.in)
		t.Run(testName, func(t *testing.T) {
			out := isHTTPSRequest(test.in)
			assert.Equal(t, test.out, out)
		})
	}
}

func TestHasPrefix(t *testing.T) {
	for _, test := range []struct {
		prefixes []string
		path     string
		out      bool
	}{
		{[]string{"/btfs"}, "/btfs/cid", true},
		{[]string{"/btfs/"}, "/btfs/cid", true},
		{[]string{"/version/"}, "/version", true},
		{[]string{"/version"}, "/version", true},
	} {
		testName := fmt.Sprintf("%+v, %s", test.prefixes, test.path)
		t.Run(testName, func(t *testing.T) {
			out := hasPrefix(test.path, test.prefixes...)
			assert.Equal(t, test.out, out)
		})
	}
}

func TestIsDomainNameAndNotPeerID(t *testing.T) {
	for _, test := range []struct {
		hostname string
		out      bool
	}{
		{"", false},
		{"example.com", true},
		{"non-icann.something", true},
		{"..", false},
		{"12D3KooWFB51PRY9BxcXSH6khFXw1BZeszeLDy7C8GciskqCTZn5", false},           // valid peerid
		{"k51qzi5uqu5di608geewp3nqkg0bpujoasmka7ftkyxgcm3fh1aroup0gsdrna", false}, // valid peerid
	} {
		t.Run(test.hostname, func(t *testing.T) {
			out := isDomainNameAndNotPeerID(test.hostname)
			assert.Equal(t, test.out, out)
		})
	}
}

func TestPortStripping(t *testing.T) {
	for _, test := range []struct {
		in  string
		out string
	}{
		{"localhost:8080", "localhost"},
		{"bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am.btfs.localhost:8080", "bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am.btfs.localhost"},
		{"example.com:443", "example.com"},
		{"example.com", "example.com"},
		{"foo-dweb.btfs.pvt.k12.ma.us:8080", "foo-dweb.btfs.pvt.k12.ma.us"},
		{"localhost", "localhost"},
		{"[::1]:8080", "::1"},
	} {
		t.Run(test.in, func(t *testing.T) {
			out := stripPort(test.in)
			assert.Equal(t, test.out, out)
		})
	}
}

func TestToDNSLabel(t *testing.T) {
	for _, test := range []struct {
		in  string
		out string
		err error
	}{
		// <= 63
		{"QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n", "QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n", nil},
		{"bafybeickencdqw37dpz3ha36ewrh4undfjt2do52chtcky4rxkj447qhdm", "bafybeickencdqw37dpz3ha36ewrh4undfjt2do52chtcky4rxkj447qhdm", nil},
		// > 63
		// PeerID: ed25519+identity multihash → CIDv1Base36
		{"bafzaajaiaejca4syrpdu6gdx4wsdnokxkprgzxf4wrstuc34gxw5k5jrag2so5gk", "k51qzi5uqu5dj16qyiq0tajolkojyl9qdkr254920wxv7ghtuwcz593tp69z9m", nil},
		// CIDv1 with long sha512 → error
		{"bafkrgqe3ohjcjplc6n4f3fwunlj6upltggn7xqujbsvnvyw764srszz4u4rshq6ztos4chl4plgg4ffyyxnayrtdi5oc4xb2332g645433aeg", "", errors.New("CID incompatible with DNS label length limit of 63: kf1siqrebi3vir8sab33hu5vcy008djegvay6atmz91ojesyjs8lx350b7y7i1nvyw2haytfukfyu2f2x4tocdrfa0zgij6p4zpl4u5oj")},
	} {
		t.Run(test.in, func(t *testing.T) {
			inCID, _ := cid.Decode(test.in)
			out, err := toDNSLabel(test.in, inCID)
			assert.Equal(t, test.out, out)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestKnownSubdomainDetails(t *testing.T) {
	gwLocalhost := &Specification{Paths: []string{"/btfs", "/btns", "/api"}, UseSubdomains: true}
	gwDweb := &Specification{Paths: []string{"/btfs", "/btns", "/api"}, UseSubdomains: true}
	gwLong := &Specification{Paths: []string{"/btfs", "/btns", "/api"}, UseSubdomains: true}
	gwWildcard1 := &Specification{Paths: []string{"/btfs", "/btns", "/api"}, UseSubdomains: true}
	gwWildcard2 := &Specification{Paths: []string{"/btfs", "/btns", "/api"}, UseSubdomains: true}

	gateways := prepareHostnameGateways(map[string]*Specification{
		"localhost":               gwLocalhost,
		"dweb.link":               gwDweb,
		"devgateway.dweb.link":    gwDweb,
		"dweb.btfs.pvt.k12.ma.us": gwLong, // note the sneaky ".btfs." ;-)
		"*.wildcard1.tld":         gwWildcard1,
		"*.*.wildcard2.tld":       gwWildcard2,
	})

	for _, test := range []struct {
		// in:
		hostHeader string
		// out:
		gw       *Specification
		hostname string
		ns       string
		rootID   string
		ok       bool
	}{
		// no subdomain
		{"127.0.0.1:8080", nil, "", "", "", false},
		{"[::1]:8080", nil, "", "", "", false},
		{"hey.look.example.com", nil, "", "", "", false},
		{"dweb.link", nil, "", "", "", false},
		// malformed Host header
		{".....dweb.link", nil, "", "", "", false},
		{"link", nil, "", "", "", false},
		{"8080:dweb.link", nil, "", "", "", false},
		{" ", nil, "", "", "", false},
		{"", nil, "", "", "", false},
		// unknown gateway host
		{"bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am.btfs.unknown.example.com", nil, "", "", "", false},
		// cid in subdomain, known gateway
		{"bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am.btfs.localhost:8080", gwLocalhost, "localhost:8080", "btfs", "bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am", true},
		{"bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am.btfs.dweb.link", gwDweb, "dweb.link", "btfs", "bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am", true},
		{"bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am.btfs.devgateway.dweb.link", gwDweb, "devgateway.dweb.link", "btfs", "bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am", true},
		// capture everything before .btfs.
		{"foo.bar.boo-buzz.btfs.dweb.link", gwDweb, "dweb.link", "btfs", "foo.bar.boo-buzz", true},
		// btns
		{"bafzbeihe35nmjqar22thmxsnlsgxppd66pseq6tscs4mo25y55juhh6bju.btns.localhost:8080", gwLocalhost, "localhost:8080", "btns", "bafzbeihe35nmjqar22thmxsnlsgxppd66pseq6tscs4mo25y55juhh6bju", true},
		{"bafzbeihe35nmjqar22thmxsnlsgxppd66pseq6tscs4mo25y55juhh6bju.btns.dweb.link", gwDweb, "dweb.link", "btns", "bafzbeihe35nmjqar22thmxsnlsgxppd66pseq6tscs4mo25y55juhh6bju", true},
		// edge case check: public gateway under long TLD (see: https://publicsuffix.org)
		{"bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am.btfs.dweb.btfs.pvt.k12.ma.us", gwLong, "dweb.btfs.pvt.k12.ma.us", "btfs", "bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am", true},
		{"bafzbeihe35nmjqar22thmxsnlsgxppd66pseq6tscs4mo25y55juhh6bju.btns.dweb.btfs.pvt.k12.ma.us", gwLong, "dweb.btfs.pvt.k12.ma.us", "btns", "bafzbeihe35nmjqar22thmxsnlsgxppd66pseq6tscs4mo25y55juhh6bju", true},
		// dnslink in subdomain
		{"en.wikipedia-on-btfs.org.btns.localhost:8080", gwLocalhost, "localhost:8080", "btns", "en.wikipedia-on-btfs.org", true},
		{"en.wikipedia-on-btfs.org.btns.localhost", gwLocalhost, "localhost", "btns", "en.wikipedia-on-btfs.org", true},
		{"dist.btfs.tech.btns.localhost:8080", gwLocalhost, "localhost:8080", "btns", "dist.btfs.tech", true},
		{"en.wikipedia-on-btfs.org.btns.dweb.link", gwDweb, "dweb.link", "btns", "en.wikipedia-on-btfs.org", true},
		// edge case check: public gateway under long TLD (see: https://publicsuffix.org)
		{"foo.dweb.btfs.pvt.k12.ma.us", nil, "", "", "", false},
		{"bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am.btfs.dweb.btfs.pvt.k12.ma.us", gwLong, "dweb.btfs.pvt.k12.ma.us", "btfs", "bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am", true},
		{"bafzbeihe35nmjqar22thmxsnlsgxppd66pseq6tscs4mo25y55juhh6bju.btns.dweb.btfs.pvt.k12.ma.us", gwLong, "dweb.btfs.pvt.k12.ma.us", "btns", "bafzbeihe35nmjqar22thmxsnlsgxppd66pseq6tscs4mo25y55juhh6bju", true},
		// other namespaces
		{"api.localhost", nil, "", "", "", false},
		{"peerid.p2p.localhost", gwLocalhost, "localhost", "p2p", "peerid", true},
		// wildcards
		{"wildcard1.tld", nil, "", "", "", false},
		{".wildcard1.tld", nil, "", "", "", false},
		{"bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am.btfs.wildcard1.tld", nil, "", "", "", false},
		{"bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am.btfs.sub.wildcard1.tld", gwWildcard1, "sub.wildcard1.tld", "btfs", "bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am", true},
		{"bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am.btfs.sub1.sub2.wildcard1.tld", nil, "", "", "", false},
		{"bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am.btfs.sub1.sub2.wildcard2.tld", gwWildcard2, "sub1.sub2.wildcard2.tld", "btfs", "bafkreicysg23kiwv34eg2d7qweipxwosdo2py4ldv42nbauguluen5v6am", true},
	} {
		t.Run(test.hostHeader, func(t *testing.T) {
			gw, hostname, ns, rootID, ok := gateways.knownSubdomainDetails(test.hostHeader)
			assert.Equal(t, test.ok, ok)
			assert.Equal(t, test.rootID, rootID)
			assert.Equal(t, test.ns, ns)
			assert.Equal(t, test.hostname, hostname)
			assert.Equal(t, test.gw, gw)
		})
	}
}
