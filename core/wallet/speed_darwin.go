// +build darwin

package wallet

import (
	"os"
	"path/filepath"
)

var portPath = func() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "~/"
	}
	return filepath.Join(dir, "Library/Application Support/uTorrent Web/BitTorrentHelper/")
}()

func validateOs() error {
	return nil
}
