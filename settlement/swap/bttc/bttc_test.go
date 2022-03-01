package bttc

import (
	"fmt"
	"strings"
	"testing"
)

// RemoveSpaceAndComma remove white space and comma
func RemoveSpaceAndComma(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), ",", "")
}

func Test123(t *testing.T) {
	s := "abc,def, hig fsl  flsdjf ll   "
	s1 := RemoveSpaceAndComma(s)
	fmt.Println(s1)
}
