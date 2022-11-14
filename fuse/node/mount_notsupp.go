//go:build (!nofuse && openbsd) || (!nofuse && netbsd)
// +build !nofuse,openbsd !nofuse,netbsd

package node

import (
	"errors"

	core "github.com/bittorrent/go-btfs/core"
)

func Mount(node *core.IpfsNode, fsdir, nsdir string) error {
	return errors.New("FUSE not supported on OpenBSD or NetBSD.")
}
