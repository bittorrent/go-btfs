package spin

import (
	"context"
	"strings"
	"time"
)

func periodicSync(period, timeout time.Duration, msg string, syncFunc func(context.Context) error) {
	tick := time.NewTicker(period)
	defer tick.Stop()
	// Force tick on immediate start
	for ; true; <-tick.C {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err := syncFunc(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "Cannot query to get stats, please wait for host validation") {
				log.Infof("Currently syncing the BTFS network, " +
					"which may last several hours depending on network conditions. Please be patient.")
			} else if strings.Contains(err.Error(), "Cannot query to get hosts") {
				log.Errorf("Cannot query to get hosts. " +
					"If you are not a renter, i.e. don't need to upload files, you can ignore it.")
			}
		}
	}
}
