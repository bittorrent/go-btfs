package announce

import (
	"fmt"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"

	cmds "github.com/bittorrent/go-btfs-cmds"

	"github.com/alecthomas/units"
)

const (
	hostStoragePriceOptionName             = "host-storage-price"
	hostBandwidthPriceOptionName           = "host-bandwidth-price"
	hostCollateralPriceOptionName          = "host-collateral-price"
	hostBandwidthLimitOptionName           = "host-bandwidth-limit"
	hostStorageTimeMinOptionName           = "host-storage-time-min"
	hostStorageMaxOptionName               = "host-storage-max"
	hostStorageEnableOptionName            = "enable-host-mode"
	hostStorageCustomizedPricingOptionName = "host-storage-customized-pricing"

	repairHostEnabledOptionName       = "repair-host-enabled"
	repairPriceDefaultOptionName      = "repair-price-default"
	repairPriceCustomizedOptionName   = "repair-price-customized"
	repairCustomizedPricingOptionName = "repair-customized-pricing"

	challengeHostEnabledOptionName       = "challenge-host-enabled"
	challengePriceDefaultOptionName      = "challenge-price-default"
	challengePriceCustomizedOptionName   = "challenge-price-customized"
	challengeCustomizedPricingOptionName = "challenge-customized-pricing"

	bttTotalSupply uint64 = 990_000_000_000
)

var StorageAnnounceCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Update and announce storage host information.",
		ShortDescription: `
This command updates host information and broadcasts to the BTFS network. 

Examples

To set the min price per GiB to 1000000 µBTT (1 BTT):
$ btfs storage announce --host-storage-price=1000000`,
	},
	Options: []cmds.Option{
		cmds.Uint64Option(hostStoragePriceOptionName, "s", "Min price per GiB of storage per day in µBTT."),
		cmds.Uint64Option(hostBandwidthPriceOptionName, "b", "Min price per MiB of bandwidth in µBTT."),
		cmds.Uint64Option(hostCollateralPriceOptionName, "cl", "Max collateral stake per hour per GiB in µBTT."),
		cmds.FloatOption(hostBandwidthLimitOptionName, "l", "Max bandwidth limit per MB/s."),
		cmds.Uint64Option(hostStorageTimeMinOptionName, "d", "Min number of days for storage."),
		cmds.Uint64Option(hostStorageMaxOptionName, "m", "Max number of GB this host provides for storage."),
		cmds.BoolOption(hostStorageEnableOptionName, "hm", "Enable/disable host storage mode. By default no mode change is made. When specified, toggles between enable/disable host mode."),
		cmds.BoolOption(hostStorageCustomizedPricingOptionName, "scp", fmt.Sprintf("Control customized pricing feature. Set false to disable and use network default price instead. Can only be enabled by explicitly setting %s.", hostStoragePriceOptionName)),
		cmds.BoolOption(repairHostEnabledOptionName, "rm", "Enable/disable repair mode. By default no mode change is made. When specified, toggles between enable/disable repair mode."),
		cmds.BoolOption(challengeHostEnabledOptionName, "cm", "Enable/disable challenge mode. By default no mode change is made. When specified, toggles between enable/disable challenge mode."),
		cmds.BoolOption(repairCustomizedPricingOptionName, "rc", "Options of repair price, true for customized price and false means default price."),
		cmds.BoolOption(challengeCustomizedPricingOptionName, "cc", "Options of challenge price, true for customized price and false means default price."),
		cmds.Uint64Option(repairPriceDefaultOptionName, "rpd", "Host repair default price refer to market."),
		cmds.Uint64Option(repairPriceCustomizedOptionName, "rpc", "Customized repair price provides by enabled Host."),
		cmds.Uint64Option(challengePriceDefaultOptionName, "cpd", "Host challenge default price refer to market."),
		cmds.Uint64Option(challengePriceCustomizedOptionName, "cpc", "Customized challenge price provides by enabled Host."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		cfg, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}

		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		hm, hmFound := req.Options[hostStorageEnableOptionName].(bool)
		if hmFound {
			// New value is different from stored, then
			if hm != cfg.Experimental.StorageHostEnabled {
				cfg.Experimental.StorageHostEnabled = hm
				err = n.Repo.SetConfig(cfg)
				if err != nil {
					return err
				}
			}
		}

		rm, rmFound := req.Options[repairHostEnabledOptionName].(bool)
		if rmFound {
			if rm != cfg.Experimental.HostRepairEnabled {
				cfg.Experimental.HostRepairEnabled = rm
				err = n.Repo.SetConfig(cfg)
				if err != nil {
					return err
				}
			}
		}

		cm, cmFound := req.Options[challengeHostEnabledOptionName].(bool)
		if cmFound {
			if cm != cfg.Experimental.HostChallengeEnabled {
				cfg.Experimental.HostChallengeEnabled = cm
				err = n.Repo.SetConfig(cfg)
				if err != nil {
					return err
				}
			}
		}

		sp, spFound := req.Options[hostStoragePriceOptionName].(uint64)
		bp, bpFound := req.Options[hostBandwidthPriceOptionName].(uint64)
		cp, cpFound := req.Options[hostCollateralPriceOptionName].(uint64)
		bl, blFound := req.Options[hostBandwidthLimitOptionName].(float64)
		stm, stmFound := req.Options[hostStorageTimeMinOptionName].(uint64)
		sm, smFound := req.Options[hostStorageMaxOptionName].(uint64)
		scp, scpFound := req.Options[hostStorageCustomizedPricingOptionName].(bool)

		rc, rcFound := req.Options[repairCustomizedPricingOptionName].(bool)
		cc, ccFound := req.Options[challengeCustomizedPricingOptionName].(bool)
		rpd, rpdFound := req.Options[repairPriceDefaultOptionName].(uint64)
		rpc, rpcFound := req.Options[repairPriceCustomizedOptionName].(uint64)
		cpd, cpdFound := req.Options[challengePriceDefaultOptionName].(uint64)
		cpc, cpcFound := req.Options[challengePriceCustomizedOptionName].(uint64)

		if sp > bttTotalSupply || cp > bttTotalSupply || bp > bttTotalSupply {
			return fmt.Errorf("maximum price is %d", bttTotalSupply)
		}

		ns, err := helper.GetHostStorageConfig(req.Context, n)
		if err != nil {
			return err
		}

		// Update fields if set
		if spFound {
			ns.StoragePriceAsk = sp
			ns.CustomizedPricing = true // turns on since we've set a price
		} else if scpFound && !scp {
			// Can only disable if no conflict with set price
			ns.StoragePriceAsk = ns.StoragePriceDefault
			ns.CustomizedPricing = false
		}
		if bpFound {
			ns.BandwidthPriceAsk = bp
		}
		if cpFound {
			ns.CollateralStake = cp
		}
		if blFound {
			ns.BandwidthLimit = bl
		}
		if stmFound {
			ns.StorageTimeMin = stm
		}
		// Storage size max is set in config instead of dynamic store
		if smFound {
			cfgRoot, err := cmdenv.GetConfigRoot(env)
			if err != nil {
				return err
			}
			sm = sm * uint64(units.GB)
			_, err = helper.CheckAndValidateHostStorageMax(req.Context, cfgRoot, n.Repo, &sm, false)
			if err != nil {
				return err
			}
		}
		if rcFound {
			ns.RepairCustomizedPricing = rc
		}
		if rpdFound {
			ns.RepairPriceDefault = rpd
		}
		if rpcFound {
			ns.RepairPriceCustomized = rpc
		}
		if ccFound {
			ns.ChallengeCustomizedPricing = cc
		}
		if cpdFound {
			ns.ChallengePriceDefault = cpd
		}
		if cpcFound {
			ns.ChallengePriceCustomized = cpc
		}

		err = helper.PutHostStorageConfig(req.Context, n, ns)
		if err != nil {
			return err
		}

		return nil
	},
}
