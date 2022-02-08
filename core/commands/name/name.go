package name

import (
	"github.com/bittorrent/go-btfs-cmds"
)

type IpnsEntry struct {
	Name  string
	Value string
}

var NameCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Publish and resolve BTNS names.",
		ShortDescription: `
BTNS is a PKI namespace, where names are the hashes of public keys, and
the private key enables publishing new (signed) values. In both publish
and resolve, the default name used is the node's own PeerID,
which is the hash of its public key.
`,
		LongDescription: `
BTNS is a PKI namespace, where names are the hashes of public keys, and
the private key enables publishing new (signed) values. In both publish
and resolve, the default name used is the node's own PeerID,
which is the hash of its public key.

You can use the 'btfs key' commands to list and generate more names and their
respective keys.

Examples:

Publish an <btfs-path> with your default name:

  > btfs name publish /btfs/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy
  Published to QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n: /btfs/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy

Publish an <btfs-path> with another name, added by an 'btfs key' command:

  > btfs key gen --type=rsa --size=2048 mykey
  > btfs name publish --key=mykey /btfs/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy
  Published to QmSrPmbaUKA3ZodhzPWZnpFgcPMFWF4QsxXbkWfEptTBJd: /btfs/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy

Resolve the value of your name:

  > btfs name resolve
  /btfs/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy

Resolve the value of another name:

  > btfs name resolve QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ
  /btfs/QmSiTko9JZyabH56y2fussEt1A5oDqsFXB3CkvAqraFryz

Resolve the value of a dnslink:

  > btfs name resolve ipfs.io
  /btfs/QmaBvfZooxWkrv7D3r8LS9moNjzD2o525XMZze69hhoxf5

`,
	},

	Subcommands: map[string]*cmds.Command{
		"publish": PublishCmd,
		"resolve": IpnsCmd,
		"pubsub":  IpnsPubsubCmd,
	},
}
