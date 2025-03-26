package commands

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
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
	"github.com/mr-tron/base58"
)

const (
	UnitWei    = "Wei"
	UnitKwei   = "KWei"
	UnitMwei   = "MWei"
	UnitGwei   = "GWei"
	UnitSzabo  = "Szabo"
	UnitFinney = "Finney"
	UnitBTT    = "BTT"
	UnitKBTT   = "KBTT"
	UnitMBTT   = "MBTT"
	UnitGBTT   = "GBTT"
	UnitTBTT   = "TBTT"
)

const (
	unitOptionName    = "unit"
	addressOptionName = "address"
	modOptionName     = "mode"
)

var StakeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Manage BTFS node staking",
		ShortDescription: "Staking commands for managing BTFS node staking operations, including stake, unlock, withdraw and query stakes.",
	},

	Subcommands: map[string]*cmds.Command{
		"info":     infoCmd,
		"unlock":   unStakeCmd,
		"withdraw": withdrawCmd,
		"query":    queryCmd,
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("amount", true, false, "the amount you want to stake"),
	},

	Options: []cmds.Option{
		cmds.StringOption(unitOptionName, "u", "the unit of amount, default is BTT").WithDefault(UnitBTT),
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		amount := req.Arguments[0]
		unit := req.Options[unitOptionName].(string)
		lockAmount, err := parseAmount(amount, unit)
		if err != nil {
			return err
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
		if err != nil {
			return err
		}
		opts, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(cfg.ChainInfo.ChainId))
		if err != nil {
			return err
		}
		opts.Value = lockAmount

		tx, err := sc.Stake(opts)
		if err != nil {
			return err
		}

		return res.Emit(map[string]string{
			"txHash": tx.Hash().Hex(),
			"status": "success",
		})
	},
}

type ContractInfo struct {
	MinStakeAmount string `json:"min_stake_amount"`
	UnlockDuration string `json:"unlock_duration"`
	Address        string `json:"address"`
}

var infoCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Query stake contract information",
		ShortDescription: `
Query stake contract information, including minimum stake amount and unlock duration.
Example: btfs stake info
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

		minStake, err := sc.MinStakeAmount(nil)
		if err != nil {
			return err
		}

		minStakeBTT := new(big.Float).Quo(
			new(big.Float).SetInt(minStake),
			new(big.Float).SetFloat64(1e18),
		)

		unlockDuration, err := sc.UnlockPeriod(nil)
		if err != nil {
			return err
		}

		durationStr, err := parseUnlockDuration(*unlockDuration)
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}

		return res.Emit(&ContractInfo{
			MinStakeAmount: minStakeBTT.Text('f', 0) + " BTT",
			UnlockDuration: durationStr,
			Address:        currChainCfg.StakeAddress.String(),
		})
	},
	Type: ContractInfo{},
}

var unStakeCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Unlock part of stake",
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
		unit := req.Options[unitOptionName].(string)
		unlockAmount, err := convert2BTT(amount, unit)
		if err != nil {
			return err
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
		if err != nil {
			return err
		}
		opts, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(cfg.ChainInfo.ChainId))
		if err != nil {
			return err
		}

		ua, ok := new(big.Int).SetString(unlockAmount, 10)
		if !ok {
			return fmt.Errorf("invalid amount")
		}
		tx, err := sc.Unstake(opts, ua)
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
		if err != nil {
			return err
		}
		opts, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(cfg.ChainInfo.ChainId))
		if err != nil {
			return err
		}

		tx, err := sc.Withdraw(opts)
		if err != nil {
			return err
		}

		return res.Emit(map[string]string{
			"status": "success",
			"txHash": tx.Hash().Hex(),
		})
	},
}

type StakeInfo struct {
	Amount       string `json:"amount"`        // Stake amount
	UnlockAmount string `json:"unlock_amount"` // Stake start time
	UnlockTime   string `json:"unlock_time"`
}

type StakeGlobalInfo struct {
	Balance       string `json:"balance"`
	TotalStaked   string `json:"total_staked"`
	TotalUnlocked string `json:"total_unlocked"`
}

var queryCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Query stake info",
		ShortDescription: `
Query stake information in different modes:
- total: query global stake statistics
- self: query your own stake info
- address: query stake info for specific address

Examples:
  btfs stake query                     # query total stats
  btfs stake query --mode self         # query your own stake
  btfs stake query --mode address --address <address>  # query specific address
`,
	},
	Options: []cmds.Option{
		cmds.StringOption(modOptionName, "m", "query mode: total/self/address").WithDefault("total"),
		cmds.StringOption(addressOptionName, "a", "address to query (required when mode is 'address')"),
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		mode := req.Options[modOptionName].(string)
		address, _ := req.Options[addressOptionName].(string)

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

		if err != nil {
			return err
		}

		switch mode {
		case "total":
			tx, err := sc.GetGlobalStats(nil)
			if err != nil {
				return err
			}
			return res.Emit(&StakeGlobalInfo{
				TotalStaked:   convertWei2BTT(tx.TotalStaked.String()),
				TotalUnlocked: convertWei2BTT(tx.TotalUnlocked.String()),
				Balance:       convertWei2BTT(tx.ContractBalance.String()),
			})
		case "self":
			pkOri, err := base64.StdEncoding.DecodeString(cfg.Identity.PrivKey)
			if err != nil {
				return err
			}
			pk, err := ethCrypto.ToECDSA(pkOri[4:])
			if err != nil {
				return err
			}

			nodeAddress := ethCrypto.PubkeyToAddress(pk.PublicKey)
			tx, err := sc.GetUserStake(nil, nodeAddress)
			if err != nil {
				return err
			}
			return res.Emit(&StakeInfo{
				Amount:       convertWei2BTT(tx.StakedAmount.String()),
				UnlockAmount: convertWei2BTT(tx.UnlockedAmount.String()),
				UnlockTime:   time.Unix(tx.UnlockTime.Int64(), 0).Format(time.RFC3339),
			})
		case "address":
			if address == "" {
				return fmt.Errorf("address is required when mode is 'address'")
			}
			// Convert ETH address to TRON address if needed
			var queryAddr common.Address
			if strings.HasPrefix(address, "T") {
				// Convert TRON address to ETH address
				decoded, err := base58.Decode(address)
				if err != nil || len(decoded) < 21 {
					return fmt.Errorf("invalid TRON address: %s", address)
				}
				// Remove TRON address version prefix (0x41) and take the next 20 bytes
				queryAddr = common.BytesToAddress(decoded[1:21])
			} else {
				queryAddr = common.HexToAddress(address)
			}

			tx, err := sc.GetUserStake(nil, queryAddr)
			if err != nil {
				return err
			}
			return res.Emit(&StakeInfo{
				Amount:       convertWei2BTT(tx.StakedAmount.String()),
				UnlockAmount: convertWei2BTT(tx.UnlockedAmount.String()),
				UnlockTime:   time.Unix(tx.UnlockTime.Int64(), 0).Format(time.RFC3339),
			})
		default:
			return fmt.Errorf("invalid mode: %s", mode)
		}
	},
	Type: StakeInfo{},
}

func convertToWei(amount string, unit string) (string, error) {
	units := map[string]string{
		UnitWei:    "",
		UnitKwei:   "000",
		UnitMwei:   "000000",
		UnitGwei:   "000000000",
		UnitSzabo:  "000000000000",
		UnitFinney: "000000000000000",
		UnitBTT:    "000000000000000000",
		UnitKBTT:   "000000000000000000000",
		UnitMBTT:   "000000000000000000000000",
		UnitGBTT:   "000000000000000000000000000",
		UnitTBTT:   "000000000000000000000000000000",
	}

	suffix, ok := units[unit]
	if !ok {
		return "", fmt.Errorf("invalid unit")
	}

	return fmt.Sprintf("%s%s", amount, suffix), nil
}

func convert2BTT(amount string, unit string) (string, error) {
	amountInWei, err := convertToWei(amount, unit)
	if err != nil {
		return "", err
	}

	return convertWei2BTT(amountInWei), nil
}

func parseAmount(amount, unit string) (*big.Int, error) {
	amount, err := convertToWei(amount, unit)
	if err != nil {
		return big.NewInt(0), err
	}

	am, ok := new(big.Int).SetString(strings.Replace(amount, ",", "", -1), 10)
	if !ok {
		return big.NewInt(0), fmt.Errorf("invalid amount: %s", amount)
	}
	return am, nil
}

func formatWithCommas(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}

	var result []byte
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, s[i])
	}
	return string(result)
}

func convertWei2BTT(amount string) string {
	am, ok := new(big.Int).SetString(strings.Replace(amount, ",", "", -1), 10)
	if !ok {
		return "0"
	}

	a := new(big.Float).Quo(
		new(big.Float).SetInt(am),
		new(big.Float).SetFloat64(1e18),
	)
	return formatWithCommas(a.Text('f', 0))
}

type DurationUnit struct {
	seconds int64
	name    string
}

var durationUnits = []DurationUnit{
	{2592000, "m"},
	{86400, "d"},
	{3600, "h"},
	{60, "m"},
	{1, "s"},
}

func parseUnlockDuration(duration big.Int) (string, error) {
	if duration.Cmp(big.NewInt(0)) == 0 {
		return "0 s", nil
	}

	for _, unit := range durationUnits {
		threshold := big.NewInt(unit.seconds)
		if duration.Cmp(threshold) >= 0 {
			result := new(big.Int).Div(new(big.Int).Set(&duration), threshold)
			return fmt.Sprintf("%d %s", result.Int64(), unit.name), nil
		}
	}

	return fmt.Sprintf("%d s", duration.Int64()), nil
}
