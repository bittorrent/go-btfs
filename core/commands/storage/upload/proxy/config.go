package proxy

import (
	"errors"
	"fmt"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
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

		priceOption := req.Options[ProxyPriceOptionName]
		if priceOption == nil {
			return fmt.Errorf("please specify the price with the --%s option", ProxyPriceOptionName)
		}

		priceInt, ok := priceOption.(int64)
		if !ok {
			return fmt.Errorf("price must be a valid integer")
		}

		tokenStr := "WBTT"
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}
		priceObj, err := chain.SettleObject.OracleService.CurrentPrice(token)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}

		if priceInt <= 0 {
			return fmt.Errorf("price must be greater than 0")
		}

		if priceObj.Uint64() > uint64(priceInt*1000000) {
			return fmt.Errorf("price must be greater than default %d (BTT)", priceObj.Uint64()/1000000)
		}

		err = helper.PutProxyStorageConfig(req.Context, n, &helper.ProxyStorageInfo{
			Price: uint64(priceInt) * 1000000,
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
				Price: nc.StoragePriceDefault / 1000,
			})
		}
		if err != nil {
			return err
		}
		config.Price /= 1000
		return cmds.EmitOnce(res, config)
	},
}
