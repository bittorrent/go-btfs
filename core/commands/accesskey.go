package commands

import (
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/s3/accesskey"
)

const ()

var AccessKeyCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "",
		ShortDescription: `
`,
	},

	Subcommands: map[string]*cmds.Command{
		"generate": accessKeyGenerateCmd,
		"enable":   accessKeyEnableCmd,
		"disable":  accessKeyDisableCmd,
		"reset":    accessKeyResetCmd,
		"delete":   accessKeyDeleteCmd,
		"get":      accessKeyGetCmd,
		"list":     accessKeyListCmd,
	},
}

var accessKeyGenerateCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "",
		ShortDescription: `
`,
	},
	Arguments: []cmds.Argument{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		ack, err := accesskey.Generate()
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, ack)
	},
}

var accessKeyEnableCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "",
		ShortDescription: `
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("key", true, true, "The key").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		key := req.Arguments[0]
		err := accesskey.Enable(key)
		return err
	},
}

var accessKeyDisableCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "",
		ShortDescription: `
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("key", true, true, "The key").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		key := req.Arguments[0]
		err := accesskey.Disable(key)
		return err
	},
}

var accessKeyResetCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "",
		ShortDescription: `
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("key", true, true, "The key").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		key := req.Arguments[0]
		err := accesskey.Reset(key)
		return err
	},
}

var accessKeyDeleteCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "",
		ShortDescription: `
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("key", true, true, "The key").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		key := req.Arguments[0]
		err := accesskey.Delete(key)
		return err
	},
}

var accessKeyGetCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "",
		ShortDescription: `
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("key", true, true, "The key").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		key := req.Arguments[0]
		ack, err := accesskey.Get(key)
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, ack)
	},
}

var accessKeyListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "",
		ShortDescription: `
`,
	},
	Arguments: []cmds.Argument{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		list, err := accesskey.List()
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, list)
	},
}
