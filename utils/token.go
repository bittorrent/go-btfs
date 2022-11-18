package utils

import "github.com/ethereum/go-ethereum/common"

const (
	WBTTHex = "0x111"
	TRXHex  = "0x222"
)

var MpTokenAddr map[string]common.Address
var MpTokenStr map[common.Address]string

func init() {
	MpTokenAddr = make(map[string]common.Address)
	MpTokenAddr["WBTT"] = common.HexToAddress(WBTTHex)
	MpTokenAddr["TRX"] = common.HexToAddress(TRXHex)

	MpTokenStr = make(map[common.Address]string)
	MpTokenStr[common.HexToAddress(WBTTHex)] = "WBTT"
	MpTokenStr[common.HexToAddress(TRXHex)] = "TRX"
}

func IsWBTT(addr common.Address) bool {
	return addr == MpTokenAddr["WBTT"]
}
