package commands

import (
	"testing"

	"github.com/bittorrent/go-btfs/namesys"

	ipns "github.com/bittorrent/go-btns"
	"github.com/libp2p/go-libp2p/core/test"
)

func TestKeyTranslation(t *testing.T) {
	pid := test.RandPeerIDFatal(t)
	pkname := namesys.PkKeyForID(pid)
	ipnsname := ipns.RecordKey(pid)

	pkk, err := escapeDhtKey("/pk/" + pid.String())
	if err != nil {
		t.Fatal(err)
	}

	ipnsk, err := escapeDhtKey("/btns/" + pid.String())
	if err != nil {
		t.Fatal(err)
	}

	if pkk != pkname {
		t.Fatal("keys didnt match!")
	}

	if ipnsk != ipnsname {
		t.Fatal("keys didnt match!")
	}
}
