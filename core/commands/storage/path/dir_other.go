//go:build !windows
// +build !windows

package path

import (
	"path/filepath"
	"strings"
)

func isHidden(path string) (bool, error) {
	_, name := filepath.Split(path)
	return strings.HasPrefix(name, "."), nil
}
