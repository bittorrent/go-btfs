package commands

import (
	"fmt"
	"io"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

const (
	SizeOptionName = "size"
)

const (
	FilterKeyPrefix = "/gateway/filter/cid"
	DefaultSize     = 50
)

var CidStoreCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Manage cid stored in this node but don't want to be get by gateway api.",
		ShortDescription: "Commands for generate, update, get and list access-keys stored in this node.",
	},
	Subcommands: map[string]*cmds.Command{
		"add":  addCidCmd,
		"del":  delCidCmd,
		"get":  getCidCmd,
		"has":  hasCidCmd,
		"list": listCidCmd,
	},
	NoLocal: true,
}

var addCidCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Add cid to store.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("cid", true, false, "cid to add to store"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		err = nd.Repo.Datastore().Put(req.Context, datastore.NewKey(NewGatewayFilterKey(req.Arguments[0])),
			[]byte(req.Arguments[0]))
		if err != nil {
			return cmds.EmitOnce(res, err.Error())
		}
		return cmds.EmitOnce(res, "Add ok.")
	},
}

var getCidCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get cid from store.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("cid", true, false, "cid to add to store"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		v, err := nd.Repo.Datastore().Get(req.Context, datastore.NewKey(NewGatewayFilterKey(req.Arguments[0])))
		if err != nil {
			return cmds.EmitOnce(res, err.Error())
		}
		return cmds.EmitOnce(res, string(v))
	},
}

var delCidCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Delete cid from store.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("cid", true, false, "cid to add to store"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		err = nd.Repo.Datastore().Delete(req.Context, datastore.NewKey(NewGatewayFilterKey(req.Arguments[0])))
		if err != nil {
			return cmds.EmitOnce(res, err.Error())
		}
		return cmds.EmitOnce(res, "Del ok.")
	},
}

var hasCidCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Check cid exits in store",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("cid", true, false, "cid to add to store"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		exits, err := nd.Repo.Datastore().Has(req.Context, datastore.NewKey(NewGatewayFilterKey(req.Arguments[0])))
		if err != nil {
			return err
		}
		if !exits {
			return cmds.EmitOnce(res, "Cid not exits")
		}
		return cmds.EmitOnce(res, "Cid exits")
	},
}

var listCidCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List all cids in store",
	},
	Options: []cmds.Option{
		cmds.IntOption(SizeOptionName, "s", "Number of cids to return.").WithDefault(DefaultSize),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		size := req.Options[SizeOptionName].(int)
		results, err := nd.Repo.Datastore().Query(req.Context, query.Query{
			Prefix: FilterKeyPrefix,
			Limit:  size,
		})
		if err != nil {
			return err
		}
		var resStr []string
		for v := range results.Next() {
			resStr = append(resStr, string(v.Value))
		}
		return cmds.EmitOnce(res, resStr)
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, cids []string) error {
			for _, v := range cids {
				_, err := w.Write([]byte(v + "\n"))
				if err != nil {
					return err
				}
			}
			return nil
		}),
	},
	Type: []string{},
}

func NewGatewayFilterKey(key string) string {
	return fmt.Sprintf("%s/%s", FilterKeyPrefix, key)
}
