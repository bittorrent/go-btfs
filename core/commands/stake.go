package commands

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/chain/abi"
	chainconfig "github.com/bittorrent/go-btfs/chain/config"
	oldcmds "github.com/bittorrent/go-btfs/commands"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var StakeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Manage BTFS node staking",
		ShortDescription: "Staking commands for managing BTFS node staking operations, including stake, unlock, withdraw and query stakes.",
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
		lockAmount, ok := new(big.Int).SetString(amount, 10)
		if !ok {
			return fmt.Errorf("invalid amount: %s", amount)
		}

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

		cli := chain.ChainObject.Backend
		if cli == nil {
			cli, err = ethclient.Dial(cfg.ChainInfo.Endpoint)
			if err != nil {
				return err
			}
		}

		sc, err := abi.NewStakeContract(contractAddress, cli)
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
		opts.Value = lockAmount

		tx, err := sc.Stake(opts)
		if err != nil {
			return err
		}

		fmt.Println("Stake success! Transaction hash is: ", tx.Hash().Hex())

		return res.Emit(map[string]string{
			"status": "success",
		})
	},
}

var unStakeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Unlock part of stake (unit: wei)",
		ShortDescription: `
Unlock part of stake.
Example: btfs stake unlock <amount>
`,
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("amount", true, false, "amount you want to unStake"),
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		amount := req.Arguments[0]
		unlockAmount, ok := new(big.Int).SetString(amount, 10)
		if !ok {
			return fmt.Errorf("invalid amount: %s", amount)
		}

		cctx := env.(*oldcmds.Context)
		cfg, err := cctx.GetConfig()

		currChainCfg, ok := chainconfig.GetChainConfig(cfg.ChainInfo.ChainId)
		if !ok {
			return fmt.Errorf("chain %d is not supported yet", cfg.ChainInfo.ChainId)
		}

		cli := chain.ChainObject.Backend
		if cli == nil {
			cli, err = ethclient.Dial(cfg.ChainInfo.Endpoint)
			if err != nil {
				return err
			}
		}

		sc, err := abi.NewStakeContract(currChainCfg.StakeAddress, cli)
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

		tx, err := sc.Unstake(opts, unlockAmount)
		if err != nil {
			return err
		}

		return res.Emit(map[string]string{
			"status": "success",
			"txHash": tx.Hash().Hex(),
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

		cli := chain.ChainObject.Backend
		if cli == nil {
			cli, err = ethclient.Dial(cfg.ChainInfo.Endpoint)
			if err != nil {
				return err
			}
		}
		sc, err := abi.NewStakeContract(currChainCfg.StakeAddress, cli)
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
		Tagline: "Query stake info by address",
		ShortDescription: `
Query stake info by address.
Example: btfs stake query <address>
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
		cli := chain.ChainObject.Backend
		if cli == nil {
			cli, err = ethclient.Dial(cfg.ChainInfo.Endpoint)
			if err != nil {
				return err
			}
		}
		sc, err := abi.NewStakeContract(currChainCfg.StakeAddress, cli)
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}

		tx, err := sc.GetUserStake(nil, common.HexToAddress(address))
		if err != nil {
			return err
		}

		return res.Emit(&StakeInfo{
			Amount:       tx.StakedAmount.String(),
			UnlockAmount: tx.UnlockedAmount.String(),
			UnlockTime:   time.Unix(tx.UnlockTime.Int64(), 0).Format(time.RFC3339),
		})

	},
	Type: StakeInfo{},
}
