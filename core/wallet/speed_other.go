// +build !windows,!darwin

package wallet

import "errors"

const (
	portPath = ""
)

func validateOs() error {
	return errors.New("not support os type")
}
