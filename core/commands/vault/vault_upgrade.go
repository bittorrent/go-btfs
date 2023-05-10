package vault

import (
	"context"
	"fmt"
	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/utils"
	"io"
	"time"
)

type VaultUpgradeCmdRet struct {
	Upgraded     bool   `json:"upgraded"`
	Description  string `json:"Description"`
	OldVaultImpl string `json:"OldVaultImpl"`
	NewVaultImpl string `json:"NewVaultImpl"`
}

var VaultUpgradeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Upgrade vault contract to the latest version",
		ShortDescription: `
If peers are using different vault versions, they can't upload/receive files to/from each other.
So we highly recommend you to upgrade your vault to the latest version. Finally, upgrading only
upgrade your vault contract's version, and won't modify any data in your vault.`,
	},
	RunTimeout: 5 * time.Minute,
	Type:       &VaultUpgradeCmdRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *VaultUpgradeCmdRet) error {
			if out.Upgraded {
				fmt.Fprintf(w, "vault version upgraded from %s to %s successfully\n", out.OldVaultImpl, out.NewVaultImpl)
			} else {
				fmt.Fprintf(w, "%s\n", out.Description)
			}
			return nil
		}),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		oldImpl, newImpl, err := chain.SettleObject.VaultService.UpgradeTo(context.Background(), chain.ChainObject.Chainconfig.VaultLogicAddress)
		upgraded := true
		description := ""
		if err != nil {
			upgraded = false
			description = fmt.Sprintf("%v", err)
		}
		return cmds.EmitOnce(res, &VaultUpgradeCmdRet{
			Upgraded:     upgraded,
			Description:  description,
			OldVaultImpl: fmt.Sprintf("%s", oldImpl),
			NewVaultImpl: fmt.Sprintf("%s", newImpl),
		})
	},
}
