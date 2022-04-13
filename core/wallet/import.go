package wallet

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/bittorrent/go-btfs/cmd/btfs/util"
	"github.com/bittorrent/go-btfs/core"

	config "github.com/TRON-US/go-btfs-config"
)

func SetKeys(n *core.IpfsNode, privKey string, mnemonic string) (err error) {
	var privK, m string
	if mnemonic != "" {
		mnemonic = strings.ReplaceAll(mnemonic, " ", ",")
		privK, m, err = util.GenerateKey("", "BIP39", mnemonic)
		if err != nil {
			return err
		}
	} else if privKey != "" {
		privKey, err = privKeyToHex(privKey)
		if err != nil {
			return err
		}
		privK, m, err = util.GenerateKey(privKey, "Secp256k1", "")
		if err != nil {
			return err
		}
	}
	identity, err := config.IdentityConfig(os.Stdout, util.NBitsForKeypairDefault, "Secp256k1", privK, m)
	if err != nil {
		return err
	}
	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}
	cfg.Identity = identity
	cfg.UI.Wallet.Initialized = false
	err = n.Repo.SetConfig(cfg)
	if err != nil {
		return err
	}
	return nil
}

func privKeyToHex(input string) (string, error) {
	isHex := true
	for _, v := range input {
		// 0-9 || A-F || a-f
		if !(v >= 48 && v <= 57 || v >= 65 && v <= 70 || v >= 97 && v <= 102) {
			isHex = false
			break
		}
	}
	if !isHex {
		bytes, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			return "", err
		}
		if len(bytes) != 36 || !strings.HasPrefix(input, "CAISI") {
			return "", fmt.Errorf("invalid privKey: %s", input)
		}
		return hex.EncodeToString(bytes[4:]), nil
	}
	_, err := hex.DecodeString(input)
	if err != nil {
		return "", err
	}
	return input, nil
}
