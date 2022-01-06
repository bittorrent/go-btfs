package priceoracle

import (
	"context"
	"errors"
	"math/big"

	conabi "github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("priceoracle")

var (
	errDecodeABI = errors.New("could not decode abi data")
)

type service struct {
	priceOracleAddress common.Address
	transactionService transaction.Service
}

type Service interface {
	// CurrentRates returns the current value of exchange rate
	// according to the latest information from oracle
	CurrentRate() (*big.Int, error)
	CurrentPrice() (*big.Int, error)
	// GetPrice retrieves latest available information from oracle
	GetPrice(ctx context.Context) (*big.Int, error)
}

var (
	priceOracleABI = transaction.ParseABIUnchecked(conabi.OracleAbi)
)

func New(priceOracleAddress common.Address, transactionService transaction.Service) Service {
	return &service{
		priceOracleAddress: priceOracleAddress,
		transactionService: transactionService,
	}
}

func (s *service) GetPrice(ctx context.Context) (*big.Int, error) {
	price, err := s.CurrentPrice()
	if err != nil {
		return nil, err
	}

	rate, err := s.CurrentRate()
	if err != nil {
		return nil, err
	}

	return price.Mul(price, rate), nil
}

func (s *service) CurrentRate() (*big.Int, error) {
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

func (s *service) CurrentPrice() (*big.Int, error) {
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
