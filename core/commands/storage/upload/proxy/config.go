package proxy

import (
	"errors"
	"fmt"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	"github.com/bittorrent/go-btfs/utils"
	ds "github.com/ipfs/go-datastore"
)

const (
	ProxyPriceOptionName = "proxy-price"
)

var StorageUploadProxyConfigCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Set storage upload proxy config.",
		ShortDescription: `
This command set storage upload proxy config such as price, the unit of price is BTT.`,
	},

	Options: []cmds.Option{
		cmds.Int64Option(ProxyPriceOptionName, "the price of proxy storage"),
	},
	Subcommands: map[string]*cmds.Command{
		"show": StorageUploadProxyConfigShowCmd,
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {

		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}
		if !cfg.Experimental.StorageClientEnabled {
			return fmt.Errorf("storage client api not enabled")
		}

		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		priceInt := req.Options[ProxyPriceOptionName].(int64)

		nc, err := helper.GetHostStorageConfig(req.Context, n)
		if err != nil {
			return err
		}
		if nc.GetStoragePriceDefault() > uint64(priceInt*1000) {
			return fmt.Errorf("price must be greater than %d", nc.StoragePriceDefault)
		}

		err = helper.PutProxyStorageConfig(req.Context, n, &helper.ProxyStorageInfo{
			Price: uint64(priceInt),
		})

		return err
	},
}

var StorageUploadProxyConfigShowCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Show storage upload proxy config.",
		ShortDescription: `
This command show storage upload proxy config such as price. The price is in BTT.`,
	},
	Type: helper.ProxyStorageInfo{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		config, err := helper.GetProxyStorageConfig(req.Context, n)
		if errors.Is(err, ds.ErrNotFound) {
			nc, err := helper.GetHostStorageConfig(req.Context, n)
			if err != nil {
				return err
			}
			return cmds.EmitOnce(res, &helper.ProxyStorageInfo{
				Price: nc.StoragePriceDefault,
			})
		}
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, config)
	},
}
