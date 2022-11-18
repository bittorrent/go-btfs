package tokencfg

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

const (
	bttcTestTokenWBTTHex = "0x111"
	bttcTestTokenTRXHex  = "0x222"

	bttcTokenWBTTHex = ""
	bttcTokenTRXHex  = ""
)

var MpTokenAddr map[string]common.Address
var MpTokenStr map[common.Address]string

func init() {
	MpTokenAddr = make(map[string]common.Address)
	MpTokenStr = make(map[common.Address]string)
}

func InitToken(chainID, bttcTestChainID, bttcChainID int64) {
	switch chainID {
	case bttcTestChainID:
		MpTokenAddr["WBTT"] = common.HexToAddress(bttcTestTokenWBTTHex)
		MpTokenAddr["TRX"] = common.HexToAddress(bttcTestTokenTRXHex)

		MpTokenStr[common.HexToAddress(bttcTestTokenTRXHex)] = "WBTT"
		MpTokenStr[common.HexToAddress(bttcTestTokenTRXHex)] = "TRX"

	case bttcChainID:
		MpTokenAddr["WBTT"] = common.HexToAddress(bttcTokenWBTTHex)
		MpTokenAddr["TRX"] = common.HexToAddress(bttcTokenTRXHex)

		MpTokenStr[common.HexToAddress(bttcTokenWBTTHex)] = "WBTT"
		MpTokenStr[common.HexToAddress(bttcTokenTRXHex)] = "TRX"
	}
}

func IsWBTT(token common.Address) bool {
	return token == MpTokenAddr["WBTT"]
}

func AddToken(s string, token common.Address) string {
	if token == MpTokenAddr["WBTT"] {
		return s
	}
	return fmt.Sprintf("%s_%s", s, token.String())
}
