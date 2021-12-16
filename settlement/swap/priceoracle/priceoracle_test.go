package priceoracle_test

import (
	"context"
	"math/big"
	"testing"

	conabi "github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/settlement/swap/priceoracle"
	"github.com/bittorrent/go-btfs/transaction"
	transactionmock "github.com/bittorrent/go-btfs/transaction/mock"
	"github.com/ethereum/go-ethereum/common"
)

var (
	priceOracleABI = transaction.ParseABIUnchecked(conabi.OracleAbi)
)

func TestExchangeGetPrice(t *testing.T) {
	priceOracleAddress := common.HexToAddress("0xabcd")

	expectedPrice := big.NewInt(100)
	expectedDeduce := big.NewInt(200)

	result := make([]byte, 64)
	expectedPrice.FillBytes(result[0:32])
	expectedDeduce.FillBytes(result[32:64])

	ex := priceoracle.New(
		priceOracleAddress,
		transactionmock.New(
			transactionmock.WithABICall(
				&priceOracleABI,
				priceOracleAddress,
				result,
				"getPrice",
			),
		),
		1,
	)

	price, err := ex.GetPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if expectedPrice.Cmp(price) != 0 {
		t.Fatalf("got wrong price. wanted %d, got %d", expectedPrice, price)
	}
}
