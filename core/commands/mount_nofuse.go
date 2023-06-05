//go:build !windows && nofuse
// +build !windows,nofuse

package commands

import (
	cmds "github.com/bittorrent/go-btfs-cmds"
)

var MountCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Mounts btfs to the filesystem (disabled).",
		ShortDescription: `
This version of btfs is compiled without fuse support, which is required
for mounting. If you'd like to be able to mount, please use a version of
btfs compiled with fuse.

For the latest instructions, please check the project's repository:
  http://github.com/bittorrent/go-btfs
`,
	},
}
