package commands

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	cmdenv "github.com/bittorrent/go-btfs/core/commands/cmdenv"
	ncmd "github.com/bittorrent/go-btfs/core/commands/name"
	ns "github.com/bittorrent/go-btfs/namesys"

	cmds "github.com/bittorrent/go-btfs-cmds"
	options "github.com/bittorrent/interface-go-btfs-core/options"
	nsopts "github.com/bittorrent/interface-go-btfs-core/options/namesys"
	path "github.com/bittorrent/interface-go-btfs-core/path"
	cidenc "github.com/ipfs/go-cidutil/cidenc"
	ipfspath "github.com/ipfs/go-path"
)

const (
	resolveRecursiveOptionName      = "recursive"
	resolveDhtRecordCountOptionName = "dht-record-count"
	resolveDhtTimeoutOptionName     = "dht-timeout"
)

var ResolveCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Resolve the value of names to BTFS.",
		ShortDescription: `
There are a number of mutable name protocols that can link among
themselves and into BTNS. This command accepts any of these
identifiers and resolves them to the referenced item.
`,
		LongDescription: `
There are a number of mutable name protocols that can link among
themselves and into BTNS. For example BTNS references can (currently)
point at an BTFS object, and DNS links can point at other DNS links, BTNS
entries, or BTFS objects. This command accepts any of these
identifiers and resolves them to the referenced item.

EXAMPLES

Resolve the value of your identity:

  $ btfs resolve /btns/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy
  /btfs/Qmcqtw8FfrVSBaRmbWwHxt3AuySBhJLcvmFYi3Lbc4xnwj

Resolve the value of another name:

  $ btfs resolve /btns/QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n
  /btns/QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy

Resolve the value of another name recursively:

  $ btfs resolve -r /btns/QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n
  /btfs/Qmcqtw8FfrVSBaRmbWwHxt3AuySBhJLcvmFYi3Lbc4xnwj

Resolve the value of an BTFS DAG path:

  $ btfs resolve /btfs/QmeZy1fGbwgVSrqbfh9fKQrAWgeyRnj7h8fsHS1oy3k99x/beep/boop
  /btfs/QmYRMjyvAiHKN9UTi8Bzt1HUspmSRD8T8DwxfSMzLgBon1

`,
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("name", true, false, "The name to resolve.").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.BoolOption(resolveRecursiveOptionName, "r", "Resolve until the result is an BTFS name.").WithDefault(true),
		cmds.IntOption(resolveDhtRecordCountOptionName, "dhtrc", "Number of records to request for DHT resolution."),
		cmds.StringOption(resolveDhtTimeoutOptionName, "dhtt", "Max time to collect values during DHT resolution eg \"30s\". Pass 0 for no timeout."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		api, err := cmdenv.GetApi(env, req)
		if err != nil {
			return err
		}

		name := req.Arguments[0]
		recursive, _ := req.Options[resolveRecursiveOptionName].(bool)

		// the case when ipns is resolved step by step
		if strings.HasPrefix(name, "/btns/") && !recursive {
			rc, rcok := req.Options[resolveDhtRecordCountOptionName].(uint)
			dhtt, dhttok := req.Options[resolveDhtTimeoutOptionName].(string)
			ropts := []options.NameResolveOption{
				options.Name.ResolveOption(nsopts.Depth(1)),
			}

			if rcok {
				ropts = append(ropts, options.Name.ResolveOption(nsopts.DhtRecordCount(rc)))
			}
			if dhttok {
				d, err := time.ParseDuration(dhtt)
				if err != nil {
					return err
				}
				if d < 0 {
					return errors.New("DHT timeout value must be >= 0")
				}
				ropts = append(ropts, options.Name.ResolveOption(nsopts.DhtTimeout(d)))
			}
			p, err := api.Name().Resolve(req.Context, name, ropts...)
			// ErrResolveRecursion is fine
			if err != nil && err != ns.ErrResolveRecursion {
				return err
			}
			return cmds.EmitOnce(res, &ncmd.ResolvedPath{Path: ipfspath.Path(p.String())})
		}

		var enc cidenc.Encoder
		switch {
		case !cmdenv.CidBaseDefined(req) && !strings.HasPrefix(name, "/ipns/"):
			// Not specified, check the path.
			enc, err = cmdenv.CidEncoderFromPath(name)
			if err == nil {
				break
			}
			// Nope, fallback on the default.
			fallthrough
		default:
			enc, err = cmdenv.GetCidEncoder(req)
			if err != nil {
				return err
			}
		}

		// else, btfs path or ipns with recursive flag
		rp, err := api.ResolvePath(req.Context, path.New(name))
		if err != nil {
			return err
		}

		encoded := "/" + rp.Namespace() + "/" + enc.Encode(rp.Cid())
		if remainder := rp.Remainder(); remainder != "" {
			encoded += "/" + remainder
		}

		return cmds.EmitOnce(res, &ncmd.ResolvedPath{Path: ipfspath.Path(encoded)})
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, rp *ncmd.ResolvedPath) error {
			fmt.Fprintln(w, rp.Path.String())
			return nil
		}),
	},
	Type: ncmd.ResolvedPath{},
}
