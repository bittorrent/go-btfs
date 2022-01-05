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

		//this is old price's rate [Compatible with older versions]
		rateObj, err := chain.SettleObject.OracleService.CurrentRate()
		if err != nil {
			return err
		}
		amount := c.SignedGuardContract.Amount
		realAmount := big.NewInt(0).Mul(big.NewInt(amount), rateObj)

		host := c.SignedGuardContract.HostPid
		contractId := c.SignedGuardContract.ContractId
		fmt.Printf("send cheque: paying...  host:%v, amount:%v, contractId:%v. \n", host, realAmount.String(), contractId)

		err = chain.SettleObject.SwapService.Settle(host, realAmount, contractId)
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}
