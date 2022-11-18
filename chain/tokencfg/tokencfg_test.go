package tokencfg

import (
	"fmt"
	"testing"
)

func TestTokenConfig(t *testing.T) {
	InitToken(1029)
	fmt.Println(MpTokenAddr)
	fmt.Println(MpTokenStr)
}
