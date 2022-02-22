package fsrepo

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

// BestKnownPath returns the best known fsrepo path. If the ENV override is
// present, this function returns that value. Otherwise, it returns the default
// repo path.
func BestKnownPath() (string, error) {
	btfsPath := "~/.btfs"
	if os.Getenv("BTFS_PATH") != "" {
		btfsPath = os.Getenv("BTFS_PATH")
	}
	curPath, err := homedir.Expand(btfsPath)
	if err != nil {
		return "", err
	}
	return curPath, nil
}
