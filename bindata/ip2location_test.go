package bindata

import (
	"testing"

	"github.com/multiformats/go-multiaddr"
)

func TestCountryCode(t *testing.T) {
	Init()

	tests := []struct {
		Addr   string
		Wanted string
	}{
		{
			"/ip4/8.8.8.8/tcp/4001",
			"US",
		},
		{
			"/ip4/36.112.144.130/tcp/4001",
			"CN",
		},
	}

	for _, test := range tests {
		md, _ := multiaddr.NewMultiaddr(test.Addr)
		code, err := CountryShortCode(md)
		//fmt.Printf("code:%s\n", code)
		if err != nil {
			t.Fatal(err)
		}

		if code != test.Wanted {
			t.Fatalf("test:%+v get code:%s != %s", test, code, test.Wanted)
		}
	}
}
