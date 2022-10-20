//go:build !cgo && !noplugin && (linux || darwin)
// +build !cgo
// +build !noplugin
// +build linux darwin

package loader

import (
	"errors"

	iplugin "github.com/bittorrent/go-btfs/plugin"
)

func init() {
	loadPluginFunc = nocgoLoadPlugin
}

func nocgoLoadPlugin(fi string) ([]iplugin.Plugin, error) {
	return nil, errors.New("not built with cgo support")
}
