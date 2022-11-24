//go:build darwin
// +build darwin

package path

import "io/ioutil"

func volumes() ([]*volume, error) {
	var vs []*volume
	fs, err := ioutil.ReadDir("/Volumes")
	if err != nil {
		return vs, err
	}
	for _, f := range fs {
		vs = append(vs, &volume{
			Name:       f.Name(),
			MountPoint: "/Volumes/" + f.Name(),
		})
	}
	return vs, nil
}
