package upload

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
)

func payInCheque(rss *sessions.RenterSession) error {
	for i, hash := range rss.ShardHashes {
		shard, err := sessions.GetUserShard(rss.CtxParams, rss.SsId, hash, i)
		if err != nil {
			return err
		}
		c, err := shard.Contracts()
		if err != nil {
			return err
		}

		// token: get real amount
		// realAmount, err := getRealAmount(c.SignedGuardContract.Amount)
		realAmount, err := getRealAmount(int64(c.Meta.Amount), rss.Token)
		if err != nil {
			return err
		}

		host := c.Meta.SpId
		contractId := c.Meta.ContractId
		fmt.Printf("send cheque: paying...  host:%v, amount:%v, contractId:%v, token:%v. \n", host, realAmount.String(), contractId, rss.Token.String())

		err = chain.SettleObject.SwapService.Settle(host, realAmount, contractId, rss.Token)
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func getRealAmount(amount int64, token common.Address) (*big.Int, error) {
	// this is price's rate [Compatible with older versions]
	rateObj, err := chain.SettleObject.OracleService.CurrentRate(token)
	if err != nil {
		return nil, err
	}

	realAmount := big.NewInt(0).Mul(big.NewInt(amount), rateObj)
	return realAmount, nil
}
