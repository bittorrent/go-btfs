package commands

import (
	"errors"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/s3/api/services/accesskey"
)

var AccessKeyCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Manage S3-Compatible-API access-keys.",
		ShortDescription: "Commands for generate, update, get and list access-keys stored in this node.",
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
	NoLocal: true,
}

func checkDaemon(env cmds.Environment) (err error) {
	node, err := cmdenv.GetNode(env)
	if err != nil {
		return
	}
	if !node.IsDaemon {
		err = errors.New("please start the node first")
	}
	return
}

var accessKeyGenerateCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Generate a new access-key record.",
		ShortDescription: "Outputs the new created access-key record.",
	},
	Arguments: []cmds.Argument{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) (err error) {
		err = checkDaemon(env)
		if err != nil {
			return
		}
		ack, err := accesskey.Generate()
		if err != nil {
			return
		}
		err = cmds.EmitOnce(res, ack)
		return
	},
}

var accessKeyEnableCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Set status of the specified access-key to enable.",
		ShortDescription: "Outputs empty if the access-key has been set to enable or it was already enabled.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("key", true, true, "The key").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) (err error) {
		err = checkDaemon(env)
		if err != nil {
			return
		}
		key := req.Arguments[0]
		err = accesskey.Enable(key)
		return
	},
}

var accessKeyDisableCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Set status of the specified access-key to enable.",
		ShortDescription: "Outputs empty if the access-key has been set to disable or it was already disabled.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("key", true, true, "The key").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) (err error) {
		err = checkDaemon(env)
		if err != nil {
			return
		}
		key := req.Arguments[0]
		err = accesskey.Disable(key)
		return
	},
}

var accessKeyResetCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Rest secret of the specified access-key.",
		ShortDescription: "Outputs the updated access-key record if it's secret has been reset.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("key", true, true, "The key").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) (err error) {
		err = checkDaemon(env)
		if err != nil {
			return
		}
		key := req.Arguments[0]
		err = accesskey.Reset(key)
		return
	},
}

var accessKeyDeleteCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Delete the specified access-key",
		ShortDescription: "Outputs empty if access-key record has been deleted.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("key", true, true, "The key").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) (err error) {
		err = checkDaemon(env)
		if err != nil {
			return
		}
		key := req.Arguments[0]
		err = accesskey.Delete(key)
		return
	},
}

var accessKeyGetCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Get an access-key detail info.",
		ShortDescription: "Outputs access-key record for the specified key.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("key", true, true, "The key").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) (err error) {
		err = checkDaemon(env)
		if err != nil {
			return
		}
		key := req.Arguments[0]
		ack, err := accesskey.Get(key)
		if err != nil {
			return
		}
		err = cmds.EmitOnce(res, ack)
		return
	},
}

var accessKeyListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "List all access-keys.",
		ShortDescription: "Outputs all non-deleted access-keys stored in current node.",
	},
	Arguments: []cmds.Argument{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) (err error) {
		err = checkDaemon(env)
		if err != nil {
			return
		}
		list, err := accesskey.List()
		if err != nil {
			return
		}
		err = cmds.EmitOnce(res, list)
		return
	},
}
