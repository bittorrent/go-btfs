package main

import (
	"testing"

	utilmain "github.com/bittorrent/go-btfs/cmd/btfs/util"
)

const privateKey = "db6bf5da1c2225d76356df2f59ebd3207b91228bcf429f0c791f94e9095a1f8e"
const mnemonic = "record access aerobic glow retreat language distance stamp cattle arrive defy movie"

func TestSeedsPhrase(t *testing.T) {
	finalImportKey, _, err := utilmain.GenerateKey(privateKey, "secp256k1", "")
	if err != nil || finalImportKey != privateKey {
		t.Error("ImportKey generated not matching")
	}

	finalImportKey, m, err := utilmain.GenerateKey("", "",
		"record,access,aerobic,glow,retreat,language,distance,stamp,cattle,arrive,defy,movie")
	if err != nil || finalImportKey != privateKey || m != mnemonic {
		t.Error("Generating from seed phrases failed")
	}

	// one less word
	finalImportKey, m, err = utilmain.GenerateKey("", "",
		"record,access,aerobic,glow,retreat,language,distance,stamp,cattle,arrive,defy")
	if err.Error() != "The seed phrase required to generate TRON private key needs to contain 12 words. Provided mnemonic has 11 words." {
		t.Error("Parameter check failed")
	}
	// word not in dictionary
	finalImportKey, m, err = utilmain.GenerateKey("", "",
		"record,access,aerobic,glow,retreat,language,distance,stamp,cattle,arrive,defy,mov")
	if err.Error() != "Entered seed phrase is not valid" {
		t.Error("Parameter check failed")
	}

}
