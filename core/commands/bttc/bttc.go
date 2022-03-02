package bttc

import cmds "github.com/bittorrent/go-btfs-cmds"

var BttcCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Interact with bttc related services.",
	},
	Subcommands: map[string]*cmds.Command{
		"btt2wbtt":     BttcBtt2WbttCmd,
		"wbtt2btt":     BttcWbtt2BttCmd,
		"send-btt-to":  BttcSendBttToCmd,
		"send-wbtt-to": BttcSendWbttToCmd,
	},
}
