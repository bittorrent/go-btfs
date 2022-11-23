package vault

import (
	cmds "github.com/bittorrent/go-btfs-cmds"
)

var VaultCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Interact with vault services on BTFS.",
		ShortDescription: `
Vault services include balance, address, withdraw, deposit operations.`,
	},
	Subcommands: map[string]*cmds.Command{
		"balance":     VaultBalanceCmd,
		"balance_all": VaultBalanceAllCmd,
		"address":     VaultAddrCmd,
		"withdraw":    VaultWithdrawCmd,
		"deposit":     VaultDepositCmd,
		"wbttbalance": VaultWbttBalanceCmd,
		"upgrade":     VaultUpgradeCmd,
	},
}
