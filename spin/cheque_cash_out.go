package spin

import (
	"github.com/bittorrent/go-btfs/chain"
)

func RestartFixChequeCashOut() {
	chain.SettleObject.CashoutService.RestartFixChequeCashOut()
}
