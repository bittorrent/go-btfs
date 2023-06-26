//go:build !windows && !nofuse
// +build !windows,!nofuse

package commands

import (
	"fmt"
	"io"

	cmdenv "github.com/bittorrent/go-btfs/core/commands/cmdenv"
	nodeMount "github.com/bittorrent/go-btfs/fuse/node"

	cmds "github.com/bittorrent/go-btfs-cmds"
	config "github.com/bittorrent/go-btfs-config"
)

const (
	mountIPFSPathOptionName = "btfs-path"
	mountIPNSPathOptionName = "btns-path"
)

var MountCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Mounts BTFS to the filesystem (read-only).",
		ShortDescription: `
Mount BTFS at a read-only mountpoint on the OS (default: /btfs and /btns).
All BTFS objects will be accessible under that directory. Note that the
root will not be listable, as it is virtual. Access known paths directly.

You may have to create /btfs and /btns before using 'btfs mount':

> sudo mkdir /btfs /btns
> sudo chown $(whoami) /btfs /btns
> btfs daemon &
> btfs mount
`,
		LongDescription: `
Mount BTFS at a read-only mountpoint on the OS. The default, /btfs and /btns,
are set in the configuration file, but can be overridden by the options.
All BTFS objects will be accessible under this directory. Note that the
root will not be listable, as it is virtual. Access known paths directly.

You may have to create /btfs and /btns before using 'btfs mount':

> sudo mkdir /btfs /btns
> sudo chown $(whoami) /btfs /btns
> btfs daemon &
> btfs mount

Example:

# setup
> mkdir foo
> echo "baz" > foo/bar
> btfs add -r foo
added QmWLdkp93sNxGRjnFHPaYg8tCQ35NBY3XPn6KiETd3Z4WR foo/bar
added QmSh5e7S6fdcu75LAbXNZAFY2nGyZUJXyLCJDvn2zRkWyC foo
> btfs ls QmSh5e7S6fdcu75LAbXNZAFY2nGyZUJXyLCJDvn2zRkWyC
QmWLdkp93sNxGRjnFHPaYg8tCQ35NBY3XPn6KiETd3Z4WR 12 bar
> btfs cat QmWLdkp93sNxGRjnFHPaYg8tCQ35NBY3XPn6KiETd3Z4WR
baz

# mount
> btfs daemon &
> btfs mount
BTFS mounted at: /btfs
BTNS mounted at: /btns
> cd /btfs/QmSh5e7S6fdcu75LAbXNZAFY2nGyZUJXyLCJDvn2zRkWyC
> ls
bar
> cat bar
baz
> cat /btfs/QmSh5e7S6fdcu75LAbXNZAFY2nGyZUJXyLCJDvn2zRkWyC/bar
baz
> cat /btfs/QmWLdkp93sNxGRjnFHPaYg8tCQ35NBY3XPn6KiETd3Z4WR
baz
`,
	},
	Options: []cmds.Option{
		cmds.StringOption(mountIPFSPathOptionName, "f", "The path where BTFS should be mounted."),
		cmds.StringOption(mountIPNSPathOptionName, "n", "The path where BTNS should be mounted."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}

		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		// error if we aren't running node in online mode
		if !nd.IsOnline {
			return ErrNotOnline
		}

		fsdir, found := req.Options[mountIPFSPathOptionName].(string)
		if !found {
			fsdir = cfg.Mounts.IPFS // use default value
		}

		// get default mount points
		nsdir, found := req.Options[mountIPNSPathOptionName].(string)
		if !found {
			nsdir = cfg.Mounts.IPNS // NB: be sure to not redeclare!
		}

		err = nodeMount.Mount(nd, fsdir, nsdir)
		if err != nil {
			return err
		}

		var output config.Mounts
		output.IPFS = fsdir
		output.IPNS = nsdir
		return cmds.EmitOnce(res, &output)
	},
	Type: config.Mounts{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, mounts *config.Mounts) error {
			fmt.Fprintf(w, "BTFS mounted at: %s\n", mounts.IPFS)
			fmt.Fprintf(w, "BTNS mounted at: %s\n", mounts.IPNS)

			return nil
		}),
	},
}
