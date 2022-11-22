package tokencfg

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func TestTokenConfig(t *testing.T) {
	InitToken(1029)
	fmt.Println(MpTokenAddr)
	fmt.Println(MpTokenStr)

	fmt.Println("zero address, ", common.HexToAddress(""))
}
