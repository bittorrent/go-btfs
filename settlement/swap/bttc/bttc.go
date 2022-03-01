package bttc

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/bittorrent/go-btfs/settlement/swap/erc20"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/ethereum/go-ethereum/common"
)

var (
	zeroBig             = big.NewInt(0)
	zeroHash            = common.Hash{}
	zeroAddr            = common.Address{}
	ErrInsufficientBTT  = errors.New("insufficient BTT")
	ErrInsufficientWBTT = errors.New("insufficient WBTT")
)

type Service interface {
	// SwapBtt2Wbtt swaps BTT to WBTT in your BTTC address
	SwapBtt2Wbtt(ctx context.Context, amount *big.Int) (trx common.Hash, err error)
	// SwapWbtt2Btt swaps WBTT to BTT in your BTTC address
	SwapWbtt2Btt(ctx context.Context, amount *big.Int) (trx common.Hash, err error)
	// SendBttTo transfers `amount` of BTT to the given address `to`
	SendBttTo(ctx context.Context, to common.Address, amount *big.Int) (trx common.Hash, err error)
	// SendWbttTo transfers `amount` of WBTT to the given address `to`
	SendWbttTo(ctx context.Context, to common.Address, amount *big.Int) (trx common.Hash, err error)
}

type service struct {
	trxService   transaction.Service
	erc20Service erc20.Service
}

func New(trxSvc transaction.Service, erc20Svc erc20.Service) *service {
	return &service{
		trxService:   trxSvc,
		erc20Service: erc20Svc,
	}
}

func (svc *service) SwapBtt2Wbtt(ctx context.Context, amount *big.Int) (trx common.Hash, err error) {
	if amount.Cmp(zeroBig) <= 0 {
		return zeroHash, errors.New("amount should bigger than zero")
	}
	balance, err := svc.trxService.MyBttBalance(ctx)
	if err != nil {
		return
	}
	fmt.Printf("your btt balance is %d, will swap %d btt to wbtt\n", balance, amount)
	if balance.Cmp(amount) < 0 {
		return zeroHash, ErrInsufficientBTT
	}
	trx, err = svc.erc20Service.Deposit(ctx, amount)
	return
}

func (svc *service) SwapWbtt2Btt(ctx context.Context, amount *big.Int) (trx common.Hash, err error) {
	if amount.Cmp(zeroBig) <= 0 {
		return zeroHash, errors.New("amount should bigger than zero")
	}
	myAddr := svc.trxService.OverlayEthAddress(ctx)
	balance, err := svc.erc20Service.BalanceOf(ctx, myAddr)
	if err != nil {
		return
	}
	fmt.Printf("your wbtt balance is %d, will swap %d wbtt to btt\n", balance, amount)
	if balance.Cmp(amount) < 0 {
		return zeroHash, ErrInsufficientWBTT
	}
	trx, err = svc.erc20Service.Withdraw(ctx, amount)
	return
}

func (svc *service) SendWbttTo(ctx context.Context, to common.Address, amount *big.Int) (trx common.Hash, err error) {
	myAddr := svc.trxService.OverlayEthAddress(ctx)
	err = validateBeforeTransfer(ctx, myAddr, to, amount)
	if err != nil {
		return zeroHash, err
	}
	balance, err := svc.erc20Service.BalanceOf(ctx, myAddr)
	if err != nil {
		return
	}
	fmt.Printf("your wbtt balance is %d, will transfer %d wbtt to address %s\n", balance, amount, to)
	if balance.Cmp(amount) < 0 {
		return zeroHash, ErrInsufficientWBTT
	}
	trx, err = svc.erc20Service.Transfer(ctx, to, amount)
	return
}

func (svc *service) SendBttTo(ctx context.Context, to common.Address, amount *big.Int) (trx common.Hash, err error) {
	myAddr := svc.trxService.OverlayEthAddress(ctx)
	err = validateBeforeTransfer(ctx, myAddr, to, amount)
	if err != nil {
		return zeroHash, err
	}
	balance, err := svc.trxService.MyBttBalance(ctx)
	if err != nil {
		return
	}
	fmt.Printf("your btt balance is %d, will transfer %d btt to address %s\n", balance, amount, to)
	if balance.Cmp(amount) < 0 {
		return zeroHash, ErrInsufficientBTT
	}

	var callData []byte
	req := &transaction.TxRequest{
		To:          &to,
		Data:        callData,
		Value:       amount,
		Description: fmt.Sprintf("send %d btt to %s", amount, to),
	}
	trx, err = svc.trxService.Send(ctx, req)
	return
}

func validateBeforeTransfer(ctx context.Context, from, to common.Address, amount *big.Int) error {
	if amount.Cmp(zeroBig) <= 0 {
		return errors.New("amount should bigger than zero")
	}
	if to == zeroAddr {
		return errors.New("please input the bttc address")
	}
	if to == from {
		return errors.New("target address is your bttc address, please input another one")
	}
	return nil
}
