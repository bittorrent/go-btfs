package commands

import (
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/rm"

	cmds "github.com/TRON-US/go-btfs-cmds"
)

const (
	rmForceOptionName = "force"
)

var RmCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Remove files or directories from a local btfs node.",
		ShortDescription: `Removes all blocks under <hash> recursively from a local btfs node.`,
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("hash", true, true, "The hash(es) of the file(s)/directory(s) to be removed from the local btfs node.").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.BoolOption(rmForceOptionName, "f", "Forcibly remove the object, even if there are constraints (such as host-stored file).").WithDefault(false),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		force, _ := req.Options[rmForceOptionName].(bool)

		results, err := rm.RmDag(req.Context, req.Arguments, n, req, env, force)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &stringList{Strings: results})
	},
	Type: stringList{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(stringListEncoder),
	},
}
