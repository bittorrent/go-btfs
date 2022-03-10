package priceoracle

import (
	"context"
	"errors"
	"math/big"
	"sync"

	conabi "github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	errDecodeABI = errors.New("could not decode abi data")
)

type service struct {
	priceOracleAddress common.Address
	transactionService transaction.Service
}

type Service interface {
	// CurrentPrice CurrentRate CurrentTotalPrice get cached info from memory.
	CurrentPrice() (*big.Int, error)
	CurrentRate() (*big.Int, error)
	CurrentTotalPrice() (*big.Int, error)
	// CheckNewPrice retrieves latest available information from oracle
	CheckNewPrice() (*big.Int, error)
}

var (
	priceOracleABI = transaction.ParseABIUnchecked(conabi.OracleAbi)

	curMutex      sync.Mutex
	curPrice      = new(big.Int)
	curRate       = new(big.Int)
	curTotalPrice = new(big.Int)
)

func New(priceOracleAddress common.Address, transactionService transaction.Service) Service {
	return &service{
		priceOracleAddress: priceOracleAddress,
		transactionService: transactionService,
	}
}

func (s *service) CurrentPrice() (price *big.Int, err error) {
	curMutex.Lock()
	price = big.NewInt(0).Set(curPrice)
	curMutex.Unlock()

	return price, nil
}
func (s *service) CurrentRate() (rate *big.Int, err error) {
	curMutex.Lock()
	rate = big.NewInt(0).Set(curRate)
	curMutex.Unlock()

	return rate, nil
}
func (s *service) CurrentTotalPrice() (totalPrice *big.Int, err error) {
	curMutex.Lock()
	totalPrice = big.NewInt(0).Set(curTotalPrice)
	curMutex.Unlock()

	return totalPrice, nil
}

func (s *service) CheckNewPrice() (*big.Int, error) {
	price, err := s.currentPrice()
	if err != nil {
		return nil, err
	}

	rate, err := s.currentRate()
	if err != nil {
		return nil, err
	}

	curMutex.Lock()
	defer curMutex.Unlock()
	curPrice = price
	curRate = rate
	curTotalPrice = big.NewInt(0).Mul(price, rate)

	return big.NewInt(0).Set(curTotalPrice), nil
}

func (s *service) currentRate() (*big.Int, error) {
	callData, err := priceOracleABI.Pack("getExchangeRate")
	if err != nil {
		return nil, err
	}
	result, err := s.transactionService.Call(context.Background(), &transaction.TxRequest{
		To:   &s.priceOracleAddress,
		Data: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := priceOracleABI.Unpack("getExchangeRate", result)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, errDecodeABI
	}

	rate, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || rate == nil {
		return nil, errDecodeABI
	}

	return rate, nil
}

func (s *service) currentPrice() (*big.Int, error) {
	callData, err := priceOracleABI.Pack("getPrice")
	if err != nil {
		return nil, err
	}
	result, err := s.transactionService.Call(context.Background(), &transaction.TxRequest{
		To:   &s.priceOracleAddress,
		Data: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := priceOracleABI.Unpack("getPrice", result)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, errDecodeABI
	}

	price, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || price == nil {
		return nil, errDecodeABI
	}

	return price, nil
}
