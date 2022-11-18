package tokencfg

import (
	"fmt"
	"testing"
)

func TestTokenConfig(t *testing.T) {
	InitToken(1029, 1029, 199)
	fmt.Println(MpTokenAddr)
	fmt.Println(MpTokenStr)
}
