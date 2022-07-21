package reportstatus

import (
	config "github.com/TRON-US/go-btfs-config"
	"github.com/bittorrent/go-btfs/chain"
)

// CheckExistLastOnline sync conf and lastOnlineInfo
func CheckExistLastOnline(cfg *config.Config) error {
	lastOnline, err := chain.GetLastOnline()
	if err != nil {
		return err
	}

	err = SyncOnlineConfig()
	if err != nil {
		return err
	}
	return err

	if lastOnline == nil {
		err = SyncOnlineConfig()
		if err != nil {
			return err
		}

		err = serv.checkLastOnlineInfo(cfg.Identity.PeerID, cfg.Identity.BttcAddr)
		if err != nil {
			return err
		}
	}
	return nil
}

func SyncOnlineConfig() error {

	return nil
}
