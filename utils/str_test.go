package utils

import (
	"testing"
)

func TestRemoveSpaceAndComma(t *testing.T) {
	type Test struct {
		s   string
		out string
	}
	testCases := []Test{
		{"abc de", "abcde"},
		{" abcde", "abcde"},
		{"abcde ", "abcde"},
		{" ab cde ", "abcde"},
		{",abcde", "abcde"},
		{"ab,cde", "abcde"},
		{"abcde,", "abcde"},
		{"ab cde,", "abcde"},
	}
	for _, test := range testCases {
		actual := RemoveSpaceAndComma(test.s)
		if actual != test.out {
			t.Errorf("RemoveSpaceAndComma(%q) = %v; want %v", test.s, actual, test.out)
		}
	}
}
