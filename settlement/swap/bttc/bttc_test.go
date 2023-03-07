package bttc_test

import (
	"context"
	"github.com/bittorrent/go-btfs/settlement/swap/erc20"
	"math/big"
	"testing"

	"github.com/bittorrent/go-btfs/settlement/swap/bttc"
	erc20Mock "github.com/bittorrent/go-btfs/settlement/swap/erc20/mock"
	"github.com/bittorrent/go-btfs/transaction"
	trxMock "github.com/bittorrent/go-btfs/transaction/mock"
	"github.com/ethereum/go-ethereum/common"
)

var (
	zeroHash common.Hash = common.Hash{}
)

func TestSwapBtt2WbttSucc(t *testing.T) {
	bttcSvc := bttc.New(
		trxMock.New(
			trxMock.WithMyBttBalance(func(ctx context.Context) (*big.Int, error) {
				return big.NewInt(100), nil
			}),
		),
		erc20Mock.New(
			erc20Mock.WithDepositFunc(func(ctx context.Context, value *big.Int) (common.Hash, error) {
				return common.HexToHash("0x123"), nil
			}),
		),
		make(map[string]erc20.Service),
	)

	trxHash, err := bttcSvc.SwapBtt2Wbtt(context.Background(), big.NewInt(50))
	if trxHash == zeroHash || err != nil {
		t.Fatal("expect swap btt to wbtt succeed")
	}
}

func TestSwapBtt2WbttFail(t *testing.T) {
	bttcSvc := bttc.New(
		trxMock.New(
			trxMock.WithMyBttBalance(func(ctx context.Context) (*big.Int, error) {
				return big.NewInt(100), nil
			}),
		),
		erc20Mock.New(
			erc20Mock.WithDepositFunc(func(ctx context.Context, value *big.Int) (common.Hash, error) {
				return common.HexToHash("0x123"), nil
			}),
		),
		make(map[string]erc20.Service),
	)

	trxHash, err := bttcSvc.SwapBtt2Wbtt(context.Background(), big.NewInt(101))
	if trxHash != zeroHash || err == nil {
		t.Fatal("expect swap btt to wbtt fail")
	}
}

func TestSwapWbtt2BttSucc(t *testing.T) {
	bttcAddr := common.HexToAddress("0xabc")

	bttcSvc := bttc.New(
		trxMock.New(
			trxMock.WithOverlayEthAddress(func(ctx context.Context) (addr common.Address) {
				return bttcAddr
			}),
		),
		erc20Mock.New(
			erc20Mock.WithWithdrawFunc(func(ctx context.Context, value *big.Int) (common.Hash, error) {
				return common.HexToHash("0x123"), nil
			}),
			erc20Mock.WithBalanceOfFunc(func(ctx context.Context, address common.Address) (*big.Int, error) {
				return big.NewInt(100), nil
			}),
		),
		make(map[string]erc20.Service),
	)

	trxHash, err := bttcSvc.SwapWbtt2Btt(context.Background(), big.NewInt(50))
	if trxHash == zeroHash || err != nil {
		t.Fatal("expect swap wbtt to btt succeed")
	}
}

func TestSendWbttToSucc(t *testing.T) {
	bttcAddr := common.HexToAddress("0xabc")
	toAddr := common.HexToAddress("0xdef")

	bttcSvc := bttc.New(
		trxMock.New(
			trxMock.WithOverlayEthAddress(func(ctx context.Context) (addr common.Address) {
				return bttcAddr
			}),
		),
		erc20Mock.New(
			erc20Mock.WithBalanceOfFunc(func(ctx context.Context, address common.Address) (*big.Int, error) {
				return big.NewInt(100), nil
			}),
			erc20Mock.WithTransferFunc(func(ctx context.Context, address common.Address, value *big.Int) (common.Hash, error) {
				return common.HexToHash("0x123"), nil
			}),
		),
		make(map[string]erc20.Service),
	)

	trxHash, err := bttcSvc.SendWbttTo(context.Background(), toAddr, big.NewInt(50))
	if trxHash == zeroHash || err != nil {
		t.Fatal("expect send wbtt to external address succeed")
	}
}

func TestSendBttToSucc(t *testing.T) {
	bttcAddr := common.HexToAddress("0xabc")
	toAddr := common.HexToAddress("0xdef")

	bttcSvc := bttc.New(
		trxMock.New(
			trxMock.WithOverlayEthAddress(func(ctx context.Context) (addr common.Address) {
				return bttcAddr
			}),
			trxMock.WithMyBttBalance(func(ctx context.Context) (*big.Int, error) {
				return big.NewInt(100), nil
			}),
			trxMock.WithSendFunc(func(ctx context.Context, request *transaction.TxRequest) (txHash common.Hash, err error) {
				return common.HexToHash("0x123"), nil
			}),
		),
		erc20Mock.New(),
		make(map[string]erc20.Service),
	)

	trxHash, err := bttcSvc.SendBttTo(context.Background(), toAddr, big.NewInt(50))
	if trxHash == zeroHash || err != nil {
		t.Fatal("expect send btt to external address succeed")
	}
}
