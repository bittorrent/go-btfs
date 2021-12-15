package spin

import (
	"os"
	"os/signal"
	"time"

	"github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/wallet"
)

const (
	period = 15 * time.Minute
)

type walletWrap struct {
	*helper.ContextParams
}

func NewWalletWrap(params *helper.ContextParams) *walletWrap {
	return &walletWrap{params}
}

func (wt *walletWrap) UpdateStatus() {
	ex := make(chan os.Signal, 1)
	signal.Notify(ex, os.Interrupt)
	go func() {
		tick := time.NewTicker(period)
		defer tick.Stop()
		for {
			sc, ec, err := wallet.UpdatePendingTransactions(wt.Ctx, wt.N.Repo.Datastore(), wt.Cfg, wt.N.Identity.String())
			log.Debugf("update pending tx, success: %v, error: %v, err: %v", sc, ec, err)
			if err := wallet.CloseLedgerChannel(wt.Ctx, wt.N.Repo.Datastore(), wt.Cfg); err != nil {
				log.Debugf("close ledger channel failed:%v", err)
			}
			select {
			case <-wt.Ctx.Done():
				return
			case <-tick.C:
				continue
			case <-ex:
				break
			}
		}
	}()
}
