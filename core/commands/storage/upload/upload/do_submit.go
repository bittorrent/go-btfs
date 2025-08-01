package upload

import (
	"context"
	"fmt"

	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	"github.com/ethereum/go-ethereum/common"

	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
)

func SubmitToChain(rss *sessions.RenterSession, fileSize int64, offlineSigning bool) error {
	err := preCheck(rss)
	if err != nil {
		return err
	}
	return addFileMetaToBttcChainAndPay(rss, fileSize, offlineSigning)
}

func prepareAmount(rss *sessions.RenterSession, shardHashes []string) (int64, error) {
	var totalPrice int64
	for i, hash := range shardHashes {
		shard, err := sessions.GetUserShard(rss.CtxParams, rss.SsId, hash, i)
		if err != nil {
			return 0, err
		}
		c, err := shard.Contracts()
		if err != nil {
			return 0, err
		}
		totalPrice += int64(c.Meta.Amount)
	}
	return totalPrice, nil
}

func preCheck(rss *sessions.RenterSession) error {
	amount, err := prepareAmount(rss, rss.ShardHashes)
	if err != nil {
		return err
	}

	err = checkAvailableBalance(rss.Ctx, amount, rss.Token)
	if err != nil {
		return err
	}

	return nil
}

func checkAvailableBalance(ctx context.Context, amount int64, token common.Address) error {
	realAmount, err := getRealAmount(amount, token)
	if err != nil {
		return err
	}

	// token: get available balance of token.
	// AvailableBalance, err := chain.SettleObject.VaultService.AvailableBalance(ctx, token)
	AvailableBalance, err := chain.SettleObject.VaultService.AvailableBalance(ctx, token)
	if err != nil {
		return err
	}

	fmt.Printf("check,  balance=%v, realAmount=%v \n", AvailableBalance, realAmount)
	if AvailableBalance.Cmp(realAmount) < 0 {
		fmt.Println("check, err: ", vault.ErrInsufficientFunds)
		return vault.ErrInsufficientFunds
	}
	return nil
}
