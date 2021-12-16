package mock

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Service struct {
	rate *big.Int
}

func New(rate *big.Int) Service {
	return Service{
		rate: rate,
	}
}

func (s Service) Start() {
}

func (s Service) GetPrice(ctx context.Context) (*big.Int, error) {
	return s.rate, nil
}

func (s Service) CurrentRates() (exchangeRate *big.Int, err error) {
	return s.rate, nil
}

func (s Service) Close() error {
	return nil
}

func DiscoverPriceOracleAddress(chainID int64) (priceOracleAddress common.Address, found bool) {
	return common.Address{}, false
}

func (s Service) SetValues(rate *big.Int) {
	s.rate = rate
}
