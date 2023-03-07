package tokencfg

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

const (
	TokenTypeName = "token-type"

	WBTT = "WBTT"
	TRX  = "TRX"
	USDD = "USDD"
	USDT = "USDT"
	TST  = "TST"

	// online
	bttcWBTTHex = "0x23181F21DEa5936e24163FFABa4Ea3B316B57f3C"
	bttcTRXHex  = "0xEdf53026aeA60f8F75FcA25f8830b7e2d6200662"
	bttcUSDDHex = "0x17f235fd5974318e4e2a5e37919a209f7c37a6d1"
	bttcUSDTHex = "0xdB28719F7f938507dBfe4f0eAe55668903D34a15"

	// test
	bttcTestWBTTHex = "0x107742eb846b86ceaaf7528d5c85cddcad3e409a"
	bttcTestTRXHex  = "0x8e009872b8a6d469939139be5e3bbd99a731212f"
	bttcTestUSDDHex = "0xa092706717dcb6892b93f0baacc07b902dbd509c"
	bttcTestUSDTHex = "0x7b906030735435422675e0679bc02dae7dfc71da"
	bttcTestTSTHex  = "0xb1cB0B7637C357108E1B72E191Aa41962019c7cc"
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

	if chainID == 199 {
		MpTokenAddr[WBTT] = common.HexToAddress(bttcWBTTHex)
		MpTokenAddr[TRX] = common.HexToAddress(bttcTRXHex)
		MpTokenAddr[USDD] = common.HexToAddress(bttcUSDDHex)
		MpTokenAddr[USDT] = common.HexToAddress(bttcUSDTHex)

		MpTokenStr[common.HexToAddress(bttcWBTTHex)] = WBTT
		MpTokenStr[common.HexToAddress(bttcTRXHex)] = TRX
		MpTokenStr[common.HexToAddress(bttcUSDDHex)] = USDD
		MpTokenStr[common.HexToAddress(bttcUSDTHex)] = USDT
	} else {
		MpTokenAddr[WBTT] = common.HexToAddress(bttcTestWBTTHex)
		MpTokenAddr[TRX] = common.HexToAddress(bttcTestTRXHex)
		MpTokenAddr[USDD] = common.HexToAddress(bttcTestUSDDHex)
		MpTokenAddr[USDT] = common.HexToAddress(bttcTestUSDTHex)
		MpTokenAddr[TST] = common.HexToAddress(bttcTestTSTHex)

		MpTokenStr[common.HexToAddress(bttcTestWBTTHex)] = WBTT
		MpTokenStr[common.HexToAddress(bttcTestTRXHex)] = TRX
		MpTokenStr[common.HexToAddress(bttcTestUSDDHex)] = USDD
		MpTokenStr[common.HexToAddress(bttcTestUSDTHex)] = USDT
		MpTokenStr[common.HexToAddress(bttcTestTSTHex)] = TST
	}

	fmt.Println("InitToken: ", chainIDStore, MpTokenAddr)
}

func GetWbttToken() common.Address {
	//fmt.Println("------ GetWbttToken ", chainIDStore)

	if chainIDStore == 199 {
		return common.HexToAddress(bttcWBTTHex)
	} else {
		return common.HexToAddress(bttcTestWBTTHex)
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
