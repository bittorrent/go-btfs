package stake

import (
	cmds "github.com/TRON-US/go-btfs-cmds"
)

var StakeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Interact with stake services on BTFS.",
		ShortDescription: `
Stake services include stake, unstake, view stake info operations.`,
	},
	Subcommands: map[string]*cmds.Command{
		"stake":     AddStakeCmd,
		"unstake":   RmStakeCmd,
		"stakeinfo": StakeInfoCmd,
		"approve":   ApproveCmd,
	},
}
