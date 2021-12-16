package upload

import (
	"fmt"
	"math/big"
	"time"

	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
)

func payInCheque(rss *sessions.RenterSession) error {
	for i, hash := range rss.ShardHashes {
		shard, err := sessions.GetRenterShard(rss.CtxParams, rss.SsId, hash, i)
		if err != nil {
			return err
		}
		c, err := shard.Contracts()
		if err != nil {
			return err
		}

		amount := c.SignedGuardContract.Amount
		host := c.SignedGuardContract.HostPid
		contractId := c.SignedGuardContract.ContractId

		fmt.Printf("send cheque: paying...  host:%v, amount:%v, contractId:%v. \n", host, amount, contractId)
		err = chain.SettleObject.SwapService.Settle(host, big.NewInt(amount), contractId)
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}
