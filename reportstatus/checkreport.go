package reportstatus

import (
	config "github.com/TRON-US/go-btfs-config"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands"
)

// CheckExistLastOnline sync conf and lastOnlineInfo
func CheckExistLastOnline(cfg *config.Config, configRoot string, chainId int64) error {
	lastOnline, err := chain.GetLastOnline()
	if err != nil {
		return err
	}

	// if nil, set config online status config
	if lastOnline == nil {
		var reportOnline bool
		var reportStatusContract bool
		if cfg.Experimental.StorageHostEnabled {
			reportOnline = true
			reportStatusContract = true
		}

		var onlineServerDomain string
		if chainId == 199 {
			onlineServerDomain = config.DefaultServicesConfig().OnlineServerDomain
		} else {
			onlineServerDomain = config.DefaultServicesConfigTestnet().OnlineServerDomain
		}

		err = commands.SyncConfigOnlineCfg(configRoot, onlineServerDomain, reportOnline, reportStatusContract)
		if err != nil {
			return err
		}
	}

	// if nil, set last online info
	if lastOnline == nil {
		err = serv.checkLastOnlineInfo(cfg.Identity.PeerID, cfg.Identity.BttcAddr)
		if err != nil {
			return err
		}
	}
	return nil
}
