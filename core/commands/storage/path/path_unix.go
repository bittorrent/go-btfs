//go:build linux || freebsd
// +build linux freebsd

package path

func volumes() ([]*volume, error) {
	var vs []*volume
	vs = append(vs, &volume{
		Name:       "Root",
		MountPoint: "/",
	})
	return vs, nil
}
