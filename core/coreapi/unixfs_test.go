package coreapi

import (
	"testing"

	"github.com/bittorrent/interface-go-btfs-core/path"

	"github.com/stretchr/testify/assert"
)

func TestIsPath(t *testing.T) {
	var pathResult = []struct {
		in  string
		out bool
	}{
		{".", true},
		{"/", true},
		{"~", true},
		{"|", false},
		{"", false},
	}
	for _, tt := range pathResult {
		t.Run(tt.in, func(t *testing.T) {
			r := isPath(tt.in)
			assert.Equal(t, r, tt.out)
		})
	}
}

func TestSplitPath(t *testing.T) {
	var empty []string
	var r = []struct {
		in        string
		hash      path.Path
		paths     []string
		errResult bool
	}{
		{"/btfs/Qm", nil, empty, true},
		{"/btfs/Qm/a/b/c", path.New("/btfs/Qm"), []string{"a", "b", "c"}, false},
		{"/btns/Qm", nil, empty, true},
		{"/btns/Qm/a/b", nil, empty, true},
		{"", nil, empty, true},
		{"", nil, empty, true},
	}
	for _, tt := range r {
		t.Run(tt.in, func(t *testing.T) {
			h, ps, e := splitPath(path.New(tt.in))
			assert.Equal(t, tt.errResult, e != nil)
			assert.Equal(t, tt.hash, h)
			assert.Equal(t, tt.paths, ps)
		})
	}
}
