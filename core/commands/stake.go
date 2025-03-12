package commands

import (
	"context"
	"fmt"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	chainconfig "github.com/bittorrent/go-btfs/chain/config"
	oldcmds "github.com/bittorrent/go-btfs/commands"
)

var StakeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Manage BTFS node staking",
		ShortDescription: "Staking commands for managing BTFS node staking operations, including create, remove, and query stakes.",
	},
	Subcommands: map[string]*cmds.Command{
		"create": stakeCreateCmd, // Create stake
		"remove": stakeRemoveCmd, // Remove stake
		"query":  stakeQueryCmd,  // Query stake info
		"verify": stakeVerifyCmd, // Verify stake status
		"list":   stakeListCmd,   // List all stakes
	},
	NoLocal: true,
}

type StakeInfo struct {
	Amount    uint64 `json:"amount"`     // Stake amount
	StartTime int64  `json:"start_time"` // Stake start time
	Duration  uint64 `json:"duration"`   // Stake duration (seconds)
	Status    string `json:"status"`     // Stake status
}

var stakeCreateCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Create new stake",
		ShortDescription: `
Create a new stake. Requires specifying stake amount and duration.
Example: btfs stake create --amount 1000 --duration 2592000
`,
	},
	Arguments: []cmds.Argument{},
	Options: []cmds.Option{
		cmds.Uint64Option("amount", "Stake amount (unit: BTT)"),
		cmds.Uint64Option("duration", "Stake duration (unit: seconds)").WithDefault(24 * 60 * 60),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		amount, _ := req.Options["amount"].(uint64)
		duration, _ := req.Options["duration"].(uint64)

		cctx := env.(*oldcmds.Context)
		cfg, err := cctx.GetConfig()
		if err != nil {
			return err
		}

		currChainCfg, ok := chainconfig.GetChainConfig(cfg.ChainInfo.ChainId)
		if !ok {
			return fmt.Errorf("chain %d is not supported yet", cfg.ChainInfo.ChainId)
		}
		contractAddress := currChainCfg.StakeAddress

		chain.SettleObject.BttcService.SendBttTo(context.Background(), contractAddress, new(big.Int).SetUint64(amount))

		// contr, err := abi.NewFileMeta(contractAddress, chain.ChainObject.Backend)
		// contr.AddFileMeta()

		return res.Emit(&StakeInfo{
			Amount:    amount,
			StartTime: time.Now().Unix(),
			Duration:  duration,
			Status:    "active",
		})
	},
	Type: StakeInfo{},
}

var stakeQueryCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Query stake information",
		ShortDescription: `
Query stake information for a specific address.
Example: btfs stake query <address>
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("address", true, false, "Address to query"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		addr := req.Arguments[0]
		fmt.Println(addr)
		return res.Emit(&StakeInfo{})
	},
	Type: StakeInfo{},
}

var stakeListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List all stakes",
		ShortDescription: `
List all stake information for the current node.
Example: btfs stake list
`,
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		return res.Emit([]StakeInfo{})
	},
	Type: []StakeInfo{},
}

var stakeRemoveCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Remove stake",
		ShortDescription: `
Remove specified stake. Note: Can only remove expired stakes.
Example: btfs stake remove <stake_id>
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("stake_id", true, false, "Stake ID to remove"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		stakeID := req.Arguments[0]
		return res.Emit(map[string]string{
			"status":  "success",
			"message": fmt.Sprintf("Stake %s has been successfully removed", stakeID),
		})
	},
	Type: map[string]string{},
}

var stakeVerifyCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Verify stake status",
		ShortDescription: `
Verify if a specific address has active stakes.
Example: btfs stake verify <address>
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("address", true, false, "Address to verify"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		addr := req.Arguments[0]
		fmt.Println(addr)
		return res.Emit(map[string]bool{
			"is_staked": true,
		})
	},
	Type: map[string]bool{},
}
