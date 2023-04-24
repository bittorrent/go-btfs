package utils

import (
	"errors"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
)

func CheckSimpleMode(env cmds.Environment) error {
	conf, err := cmdenv.GetConfig(env)
	if err != nil {
		return err
	}

	//fmt.Println("CheckSimpleMode ... ", conf.SimpleMode)

	if conf.SimpleMode {
		return errors.New("this api is not support in simple mode, please check the node's simple mode! ")
	}

	return nil
}
