package commands

import (
	"encoding/base64"
	"fmt"
	"math/big"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/chain/abi"
	chainconfig "github.com/bittorrent/go-btfs/chain/config"
	oldcmds "github.com/bittorrent/go-btfs/commands"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

var StakeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Manage BTFS node staking",
		ShortDescription: "Staking commands for managing BTFS node staking operations, including create, remove, and query stakes.",
	},

	Subcommands: map[string]*cmds.Command{
		"unlock":   unStakeCmd,  // Unlock part of stake
		"withdraw": withdrawCmd, // Withdraw all stake
		"query":    queryCmd,    // Query user stakes
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("amount", true, false, "the amount you want to stake (unit: BTT)"),
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		amount := req.Arguments[0]

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

		sc, err := abi.NewStakeContract(contractAddress, chain.ChainObject.Backend)
		if err != nil {
			return err
		}

		pkOri, err := base64.StdEncoding.DecodeString(cfg.Identity.PrivKey)
		if err != nil {
			return err
		}

		pk, err := ethCrypto.ToECDSA(pkOri[4:])
		opts, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(cfg.ChainInfo.ChainId))
		if err != nil {
			return err
		}
		if opts.Value, ok = new(big.Int).SetString(amount, 10); !ok {
			return fmt.Errorf("invalid amount: %s", amount)
		}

		tx, err := sc.Stake(opts)
		if err != nil {
			return err
		}

		fmt.Println("Stake success! Transaction hash is: ", tx.Hash().Hex())

		return res.Emit(map[string]string{
			"status": "success",
		})
	},

	NoLocal: true,
}

var unStakeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Remove stake",
		ShortDescription: `
Remove specified stake. Note: Can only remove expired stakes.
Example: btfs stake remove <stake_id>
`,
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("amount", true, false, "amount you want to unStake"),
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		amount, _ := req.Options["amount"].(uint64)

		cctx := env.(*oldcmds.Context)
		cfg, err := cctx.GetConfig()

		currChainCfg, ok := chainconfig.GetChainConfig(cfg.ChainInfo.ChainId)
		if !ok {
			return fmt.Errorf("chain %d is not supported yet", cfg.ChainInfo.ChainId)
		}
		contractAddress := currChainCfg.StakeAddress
		sc, err := abi.NewStakeContract(contractAddress, chain.ChainObject.Backend)
		if err != nil {
			return err
		}

		pkOri, err := base64.StdEncoding.DecodeString(cfg.Identity.PrivKey)
		if err != nil {
			return err
		}

		pk, err := ethCrypto.ToECDSA(pkOri[4:])
		opts, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(cfg.ChainInfo.ChainId))
		if err != nil {
			return err
		}

		tx, err := sc.Unstake(opts, new(big.Int).SetUint64(amount))
		if err != nil {
			return err
		}

		fmt.Println("UnStake success! Transaction hash is: ", tx.Hash().Hex())

		return res.Emit(map[string]string{
			"status": "success",
		})
	},
	Type: map[string]string{},
}

var withdrawCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Withdraw all stake",
		ShortDescription: `
Withdraw all stake.
Example: btfs stake withdraw
`,
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
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

		sc, err := abi.NewStakeContract(contractAddress, chain.ChainObject.Backend)
		if err != nil {
			return err
		}

		pkOri, err := base64.StdEncoding.DecodeString(cfg.Identity.PrivKey)
		if err != nil {
			return err
		}

		pk, err := ethCrypto.ToECDSA(pkOri[4:])
		opts, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(cfg.ChainInfo.ChainId))
		if err != nil {
			return err
		}

		tx, err := sc.Withdraw(opts)
		if err != nil {
			return err
		}

		fmt.Println("Withdraw success! Transaction hash is: ", tx.Hash().Hex())

		return res.Emit(map[string]string{
			"status": "success",
		})
	},
}

type StakeInfo struct {
	Amount       string `json:"amount"`        // Stake amount
	UnlockAmount string `json:"unlock_amount"` // Stake start time
	UnlockTime   string `json:"unlock_time"`
}

var queryCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Remove stake",
		ShortDescription: `
Remove specified stake. Note: Can only remove expired stakes.
Example: btfs stake remove <stake_id>
`,
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("address", true, false, "address you want to query"),
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		address := req.Arguments[0]

		cctx := env.(*oldcmds.Context)
		cfg, err := cctx.GetConfig()

		currChainCfg, ok := chainconfig.GetChainConfig(cfg.ChainInfo.ChainId)
		if !ok {
			return fmt.Errorf("chain %d is not supported yet", cfg.ChainInfo.ChainId)
		}
		contractAddress := currChainCfg.StakeAddress
		sc, err := abi.NewStakeContract(contractAddress, chain.ChainObject.Backend)
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}

		opts := &bind.CallOpts{}
		if err != nil {
			return err
		}

		tx, err := sc.GetUserStake(opts, common.HexToAddress(address))
		if err != nil {
			return err
		}

		return res.Emit(&StakeInfo{
			Amount:       tx.StakedAmount.String(),
			UnlockAmount: tx.UnlockedAmount.String(),
			UnlockTime:   tx.UnlockTime.String(),
		})

	},
	Type: StakeInfo{},
}
