package utils

import "strings"

// RemoveSpaceAndComma remove white space and comma
func RemoveSpaceAndComma(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), ",", "")
}
