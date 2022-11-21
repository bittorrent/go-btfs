package name

import (
	"errors"
	"fmt"
	"io"
	"time"

	cmdenv "github.com/bittorrent/go-btfs/core/commands/cmdenv"
	ke "github.com/bittorrent/go-btfs/core/commands/keyencode"

	iface "github.com/TRON-US/interface-go-btfs-core"
	options "github.com/TRON-US/interface-go-btfs-core/options"
	path "github.com/TRON-US/interface-go-btfs-core/path"
	cmds "github.com/bittorrent/go-btfs-cmds"

	peer "github.com/libp2p/go-libp2p/core/peer"
)

var (
	errAllowOffline = errors.New("can't publish while offline: pass `--allow-offline` to override")
)

const (
	ipfsPathOptionName     = "btfs-path"
	resolveOptionName      = "resolve"
	allowOfflineOptionName = "allow-offline"
	lifeTimeOptionName     = "lifetime"
	ttlOptionName          = "ttl"
	keyOptionName          = "key"
	quieterOptionName      = "quieter"
)

var PublishCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Publish BTNS names.",
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

Alternatively, publish an <btfs-path> using a valid PeerID (as listed by 
'btfs key list -l'):

 > btfs name publish --key=QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n /btfs/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy
  Published to QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n: /btfs/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy

`,
	},

	Arguments: []cmds.Argument{
		cmds.StringArg(ipfsPathOptionName, true, false, "btfs path of the object to be published.").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.BoolOption(resolveOptionName, "Check if the given path can be resolved before publishing.").WithDefault(true),
		cmds.StringOption(lifeTimeOptionName, "t",
			`Time duration that the record will be valid for. <<default>>
    This accepts durations such as "300s", "1.5h" or "2h45m". Valid time units are
    "ns", "us" (or "µs"), "ms", "s", "m", "h".`).WithDefault("24h"),
		cmds.BoolOption(allowOfflineOptionName, "When offline, save the BTNS record to the the local datastore without broadcasting to the network instead of simply failing."),
		cmds.StringOption(ttlOptionName, "Time duration this record should be cached for. Uses the same syntax as the lifetime option. (caution: experimental)"),
		cmds.StringOption(keyOptionName, "k", "Name of the key to be used or a valid PeerID, as listed by 'btfs key list -l'.").WithDefault("self"),
		cmds.BoolOption(quieterOptionName, "Q", "Write only final hash."),
		ke.OptionIPNSBase,
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		api, err := cmdenv.GetApi(env, req)
		if err != nil {
			return err
		}
		keyEnc, err := ke.KeyEncoderFromString(req.Options[ke.OptionIPNSBase.Name()].(string))
		if err != nil {
			return err
		}

		allowOffline, _ := req.Options[allowOfflineOptionName].(bool)
		kname, _ := req.Options[keyOptionName].(string)

		validTimeOpt, _ := req.Options[lifeTimeOptionName].(string)
		validTime, err := time.ParseDuration(validTimeOpt)
		if err != nil {
			return fmt.Errorf("error parsing lifetime option: %s", err)
		}

		opts := []options.NamePublishOption{
			options.Name.AllowOffline(allowOffline),
			options.Name.Key(kname),
			options.Name.ValidTime(validTime),
		}

		if ttl, found := req.Options[ttlOptionName].(string); found {
			d, err := time.ParseDuration(ttl)
			if err != nil {
				return err
			}

			opts = append(opts, options.Name.TTL(d))
		}

		p := path.New(req.Arguments[0])

		if verifyExists, _ := req.Options[resolveOptionName].(bool); verifyExists {
			_, err := api.ResolveNode(req.Context, p)
			if err != nil {
				return err
			}
		}

		out, err := api.Name().Publish(req.Context, p, opts...)
		if err != nil {
			if err == iface.ErrOffline {
				err = errAllowOffline
			}
			return err
		}

		// parse path, extract cid, re-base cid, reconstruct path
		pid, err := peer.Decode(out.Name())
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &IpnsEntry{
			Name:  keyEnc.FormatID(pid),
			Value: out.Value().String(),
		})
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, ie *IpnsEntry) error {
			var err error
			quieter, _ := req.Options[quieterOptionName].(bool)
			if quieter {
				_, err = fmt.Fprintln(w, ie.Name)
			} else {
				_, err = fmt.Fprintf(w, "Published to %s: %s\n", ie.Name, ie.Value)
			}
			return err
		}),
	},
	Type: IpnsEntry{},
}
