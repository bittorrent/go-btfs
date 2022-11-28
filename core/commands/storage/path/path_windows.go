//go:build windows
// +build windows

package path

import (
	"os"
)

func volumes() ([]*volume, error) {
	vs := make([]*volume, 0)
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		d := string(drive) + ":\\"
		f, err := os.Open(string(drive) + ":\\")
		if err == nil {
			vs = append(vs, &volume{
				Name:       d,
				MountPoint: d,
			})
			f.Close()
		}
	}
	return vs, nil
}
