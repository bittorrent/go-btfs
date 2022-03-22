package wallet

import (
	"testing"

	coremock "github.com/bittorrent/go-btfs/core/mock"

	"github.com/stretchr/testify/assert"
)

const (
	expectedPrivKeyBase64 = "CAISIFmxIUg17m/CM3nAeRAsjKHMb7pgkVmCCYfEFyES9Jkx"
	expectedPrivKeyHex    = "59b1214835ee6fc23379c079102c8ca1cc6fba609159820987c4172112f49931"
	expectedMnemonic      = "absurd adapt skin skin settle smart other table school toss give reform"
)

func TestImportPrivateKey(t *testing.T) {
	var testCases = []struct {
		privKey          string // input
		mnemonic         string // input
		expectedPrivKey  string // expected
		expectedMnemonic string // expected
		returnErr        bool   // expected
	}{
		{"CAISIFmxIUg17m/CM3nAeRAsjKHMb7pgkVmCCYfEFyES9Jkx", "", expectedPrivKeyBase64, "", false},
		{"59b1214835ee6fc23379c079102c8ca1cc6fba609159820987c4172112f49931", "", expectedPrivKeyBase64, "", false},
		{"", "absurd adapt skin skin settle smart other table school toss give reform", expectedPrivKeyBase64, expectedMnemonic, false},
		{"CAISIFmxIUg17m/CM3nAeRAsjKHMb7pgkVmCCYfEFyES9Jkx", "absurd adapt skin skin settle smart other table school toss give reform", expectedPrivKeyBase64, expectedMnemonic, false},
		{"errorPrivKey", "", "", "", true},
		{"", "errorMnemonic", "", "", true},
		{"errorPrivKey", "errorMnemonic", "", "", true},
	}

	n, err := coremock.NewMockNode()
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := n.Repo.Config()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range testCases {
		err := SetKeys(n, tc.privKey, tc.mnemonic)
		assert.Equal(t, tc.returnErr, err != nil)
		if err == nil {
			assert.Equal(t, tc.expectedPrivKey, cfg.Identity.PrivKey)
			assert.Equal(t, tc.expectedMnemonic, cfg.Identity.Mnemonic)
		}
	}
}
