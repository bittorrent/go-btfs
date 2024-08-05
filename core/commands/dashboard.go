package commands

import (
	"bytes"
	"errors"
	"fmt"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/utils"
	ds "github.com/ipfs/go-datastore"
)

const DashboardPasswordPrefix = "/dashboard_password"

const TokenOption = "token"

var dashboardCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "dashboard password operation",
	},

	Subcommands: map[string]*cmds.Command{
		"check": checkCmd,
		"set":   setCmd,
		"login": loginCmd,
		"reset": resetCmd,
	},
}

var checkCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "check if password is set",
	},
	Run: func(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) error {
		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		value, err := node.Repo.Datastore().Get(req.Context, ds.NewKey(DashboardPasswordPrefix))
		if err != nil {
			log.Info("check password error", err)
			return errors.New("password is not set")
		}
		fmt.Println("password..............", string(value))
		return re.Emit(bytes.NewReader([]byte("check")))
	},
}

var setCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "set password",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("password", true, false, "set password"),
	},
	Run: func(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) error {
		// 写入leveldb
		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		datastore := node.Repo.Datastore()
		key := ds.NewKey(DashboardPasswordPrefix)
		err = datastore.Put(req.Context, key, []byte(req.Arguments[0]))
		if err != nil {
			return err
		}
		return re.Emit(bytes.NewReader([]byte("set")))
	},
}

var loginCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "login password",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("password", true, false, "set password"),
	},
	Run: func(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) error {
		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		value, err := node.Repo.Datastore().Get(req.Context, ds.NewKey(DashboardPasswordPrefix))
		if err != nil {
			return err
		}
		if string(value) != req.Arguments[0] {
			return errors.New("password is not correct")
		}
		log.Info("login password is correct")

		config, err := node.Repo.Config()

		publicKey := config.Identity.PeerID

		token, err := utils.GenerateToken(publicKey, req.Arguments[0], 60*60)
		if err != nil {
			return err
		}

		return re.Emit(bytes.NewReader([]byte(token)))
	},
}

var resetCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "reset password",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("password", true, false, "reset password"),
	},
	Run: func(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) error {
		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		datastore := node.Repo.Datastore()
		value, err := datastore.Get(req.Context, ds.NewKey(DashboardPasswordPrefix))
		if err != nil {
			return err
		}
		if string(value) != req.Arguments[0] {
			return errors.New("password is not correct")
		}
		err = datastore.Put(req.Context, ds.NewKey(DashboardPasswordPrefix), []byte(""))
		if err != nil {
			return err
		}
		return re.Emit(bytes.NewReader([]byte("reset")))
	},
}
