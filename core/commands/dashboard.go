package commands

import (
	"encoding/hex"
	"errors"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/utils"
	ds "github.com/ipfs/go-datastore"
)

const DashboardPasswordPrefix = "/dashboard_password"
const TokenExpire = 60 * 60 * 24 * 1

var IsLogin = false

type DashboardResponse struct {
	Success bool
	Text    string
}

var dashboardCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "dashboard password operation",
	},

	Subcommands: map[string]*cmds.Command{
		"check":  checkCmd,
		"set":    setCmd,
		"reset":  resetCmd,
		"change": changeCmd,
		"login":  loginCmd,
		"logout": logoutCmd,
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

		_, err = node.Repo.Datastore().Get(req.Context, ds.NewKey(DashboardPasswordPrefix))

		if err != nil && errors.Is(err, ds.ErrNotFound) {
			return re.Emit(&DashboardResponse{Success: false, Text: "passwd is not set"})
		}

		if err != nil {
			log.Info("check password error", err)
			return err
		}
		return re.Emit(DashboardResponse{Success: true, Text: "password was set"})
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
		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		// check if password has set
		_, err = node.Repo.Datastore().Get(req.Context, ds.NewKey(DashboardPasswordPrefix))
		if err != nil && !errors.Is(err, ds.ErrNotFound) {
			log.Info("set password error", err)
			return err
		}

		if err == nil {
			return re.Emit(&DashboardResponse{
				Success: false,
				Text:    "password has set, if you want to reset your password, please use reset command instead",
			})
		}

		datastore := node.Repo.Datastore()
		key := ds.NewKey(DashboardPasswordPrefix)
		err = datastore.Put(req.Context, key, []byte(req.Arguments[0]))
		if err != nil {
			return err
		}
		return re.Emit(&DashboardResponse{Success: true, Text: "password set success!"})
	},
}

var loginCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "login with passwd and get the token",
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
		if errors.Is(err, ds.ErrNotFound) {
			return re.Emit(&DashboardResponse{Success: false, Text: "password has not set, please set passwd first"})
		}
		if err != nil {
			return err
		}
		if string(value) != req.Arguments[0] {
			log.Info("login password is correct")
			return re.Emit(&DashboardResponse{Success: false, Text: "password is not correct"})
		}

		config, err := node.Repo.Config()

		publicKey := config.Identity.PeerID

		token, err := utils.GenerateToken(publicKey, req.Arguments[0], TokenExpire)
		if err != nil {
			return err
		}

		IsLogin = true

		return re.Emit(&DashboardResponse{
			Success: true,
			Text:    token,
		})
	},
}

var resetCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "reset password",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("privateKey", true, false, "private key"),
		cmds.StringArg("oldPassword", true, false, "old password"),
		cmds.StringArg("newPassword", true, false, "new password"),
	},

	Run: func(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) error {
		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		raw, err := node.PrivateKey.Raw()
		if err != nil {
			return err
		}

		if hex.EncodeToString(raw) != req.Arguments[0] {
			return re.Emit(&DashboardResponse{Success: false, Text: "private key is not correct"})
		}
		datastore := node.Repo.Datastore()

		value, err := datastore.Get(req.Context, ds.NewKey(DashboardPasswordPrefix))
		if err != nil {
			return err
		}
		if string(value) != req.Arguments[1] {
			return re.Emit(&DashboardResponse{Success: false, Text: "the old password is not correct"})
		}
		err = datastore.Put(req.Context, ds.NewKey(DashboardPasswordPrefix), []byte(req.Arguments[2]))
		if err != nil {
			return err
		}
		return re.Emit(&DashboardResponse{Success: true, Text: "password reset success!"})
	},
}

var changeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "change password",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("oldPassword", true, false, "change password"),
		cmds.StringArg("newPassword", true, false, "change password"),
	},
	Run: func(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) error {
		// change password
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
			return re.Emit(&DashboardResponse{Success: false, Text: "the old password is not correct"})
		}
		err = datastore.Put(req.Context, ds.NewKey(DashboardPasswordPrefix), []byte(req.Arguments[1]))
		if err != nil {
			return err
		}
		return re.Emit(&DashboardResponse{Success: true, Text: "password change success!"})
	},
}

var logoutCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "logout",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("token", true, false, "logout"),
	},
	Run: func(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) error {
		// set token expire to 0
		IsLogin = false
		return re.Emit(&DashboardResponse{Success: true, Text: "logout success!"})
	},
}