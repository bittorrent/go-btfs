package util

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

const (
	NBitsForKeypairDefault = 2048
	defaultEntropy         = 128
	mnemonicLength         = 12
)

func GenerateKey(importKey string, keyType string, seedPhrase string) (string, string, error) {
	mnemonicLen := len(strings.Split(seedPhrase, ","))
	mnemonic := strings.ReplaceAll(seedPhrase, ",", " ")
	if importKey != "" && keyType != "" && strings.ToLower(keyType) != "secp256k1" {
		return "", "", fmt.Errorf("cannot specify key type and import TRON private key at the same time")
	} else if seedPhrase != "" {
		if mnemonicLen != mnemonicLength {
			return "", "", fmt.Errorf("The seed phrase required to generate TRON private key needs to contain 12 words. Provided mnemonic has %v words.", mnemonicLen)
		}
		if err := !bip39.IsMnemonicValid(mnemonic); err {
			return "", "", fmt.Errorf("Entered seed phrase is not valid")
		}
		fmt.Println("Generating TRON key with BIP39 seed phrase...")
		return GeneratePrivKeyUsingBIP39(mnemonic)
	} else if (keyType == "" && importKey == "") || keyType == "BIP39" {
		fmt.Println("Generating TRON key with BIP39 seed phrase...")
		return GeneratePrivKeyUsingBIP39("")
	} else {
		return importKey, mnemonic, nil
	}

}

func GeneratePrivKeyUsingBIP39(mnemonic string) (string, string, error) {
	if mnemonic == "" {
		entropy, err := bip39.NewEntropy(defaultEntropy)
		if err != nil {
			return "", "", err
		}
		mnemonic, err = bip39.NewMnemonic(entropy)
		if err != nil {
			return "", "", err
		}
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		return "", "", fmt.Errorf("invalid mnemonic: %s", mnemonic)
	}

	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, "")

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", "", err
	}
	publicKey := masterKey.PublicKey()

	childKey, err := masterKey.NewChildKey(44 + bip32.FirstHardenedChild)
	if err != nil {
		return "", "", err
	}
	childKey2, err := childKey.NewChildKey(195 + bip32.FirstHardenedChild)
	if err != nil {
		return "", "", err
	}
	childKey3, err := childKey2.NewChildKey(0 + bip32.FirstHardenedChild)
	if err != nil {
		return "", "", err
	}
	childKey4, err := childKey3.NewChildKey(0)
	if err != nil {
		return "", "", err
	}
	childKey5, err := childKey4.NewChildKey(0)
	if err != nil {
		return "", "", err
	}

	encoding := childKey5.Key
	importKey := hex.EncodeToString(encoding)

	// Display mnemonic and keys
	fmt.Println("Master public key: ", publicKey)

	return importKey, mnemonic, nil
}
