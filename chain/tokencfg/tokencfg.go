package tokencfg

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

const (
	TokenTypeName = "token-type"

	bttcTestTokenWBTTHex = "0x107742eb846b86ceaaf7528d5c85cddcad3e409a"
	bttcTestTokenTRXHex  = "0xb1cB0B7637C357108E1B72E191Aa41962019c7cc"

	bttcTokenWBTTHex = "0x23181F21DEa5936e24163FFABa4Ea3B316B57f3C"
	bttcTokenTRXHex  = "0xEdf53026aeA60f8F75FcA25f8830b7e2d6200662"
)

var chainIDStore int64

var MpTokenAddr map[string]common.Address
var MpTokenStr map[common.Address]string

func init() {
	MpTokenAddr = make(map[string]common.Address)
	MpTokenStr = make(map[common.Address]string)
}

func InitToken(chainID int64) {
	chainIDStore = chainID
	fmt.Println("------ InitToken ", chainIDStore)

	if chainID == 199 {
		MpTokenAddr["WBTT"] = common.HexToAddress(bttcTokenWBTTHex)
		MpTokenAddr["TRX"] = common.HexToAddress(bttcTokenTRXHex)

		MpTokenStr[common.HexToAddress(bttcTokenWBTTHex)] = "WBTT"
		MpTokenStr[common.HexToAddress(bttcTokenTRXHex)] = "TRX"
	} else {
		MpTokenAddr["WBTT"] = common.HexToAddress(bttcTestTokenWBTTHex)
		MpTokenAddr["TRX"] = common.HexToAddress(bttcTestTokenTRXHex)

		MpTokenStr[common.HexToAddress(bttcTestTokenWBTTHex)] = "WBTT"
		MpTokenStr[common.HexToAddress(bttcTestTokenTRXHex)] = "TRX"
	}
}

func GetWbttToken() common.Address {
	fmt.Println("------ GetWbttToken ", chainIDStore)

	if chainIDStore == 199 {
		return common.HexToAddress(bttcTokenWBTTHex)
	} else {
		return common.HexToAddress(bttcTestTokenWBTTHex)
	}
}

func IsWBTT(token common.Address) bool {
	return token == MpTokenAddr["WBTT"]
}

func AddToken(s string, token common.Address) string {
	if token == MpTokenAddr["WBTT"] {
		return s
	}
	return fmt.Sprintf("%s_%s", token.String(), s)
}
