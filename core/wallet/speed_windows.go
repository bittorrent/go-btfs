// +build windows

package wallet

import "os"

var portPath = func() string {
	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home + "\\AppData\\Local\\BitTorrentHelper\\"
}()

func validateOs() error {
	return nil
}
