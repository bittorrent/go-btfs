package vault

import (
	"fmt"
	"github.com/bittorrent/go-btfs/utils"
	"io"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

type VaultAddrCmdRet struct {
	Addr string `json:"addr"`
}

var VaultAddrCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get vault address.",
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		addr := chain.SettleObject.VaultService.Address()

		return cmds.EmitOnce(res, &VaultAddrCmdRet{
			Addr: addr.String(),
		})
	},
	Type: &VaultAddrCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *VaultAddrCmdRet) error {
			_, err := fmt.Fprintf(w, "the vault addr: %s", out.Addr)
			return err
		}),
	},
}
