package priceoracle

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/math"

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
	exchangeRate       *big.Int
	timeDivisor        int64
	quitC              chan struct{}
}

type Service interface {
	io.Closer
	// CurrentRates returns the current value of exchange rate
	// according to the latest information from oracle
	CurrentRates() (exchangeRate *big.Int, err error)
	// GetPrice retrieves latest available information from oracle
	GetPrice(ctx context.Context) (*big.Int, error)
	Start()
	SwitchCurrentRates() (exchangeRate *big.Int, err error)
}

var (
	priceOracleABI = transaction.ParseABIUnchecked(conabi.OracleAbi)
)

func New(priceOracleAddress common.Address, transactionService transaction.Service, timeDivisor int64) Service {
	return &service{
		priceOracleAddress: priceOracleAddress,
		transactionService: transactionService,
		exchangeRate:       big.NewInt(0),
		quitC:              make(chan struct{}),
		timeDivisor:        timeDivisor,
	}
}

func (s *service) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		<-s.quitC
	}()

	go func() {
		defer cancel()
		for {
			exchangeRate, err := s.GetPrice(ctx)
			if err != nil {
				log.Errorf("could not get price: %v", err)
			} else {
				log.Infof("updated exchange rate to %d", exchangeRate)
				s.exchangeRate = exchangeRate
			}

			ts := time.Now().Unix()

			// We poll the oracle in every timestamp divisible by constant 300 (timeDivisor)
			// in order to get latest version approximately at the same time on all nodes
			// and to minimize polling frequency
			// If the node gets newer information than what was applicable at last polling point at startup
			// this minimizes the negative scenario to less than 5 minutes
			// during which cheques can not be sent / accepted because of the asymmetric information
			timeUntilNextPoll := time.Duration(s.timeDivisor-ts%s.timeDivisor) * time.Second

			select {
			case <-s.quitC:
				return
			case <-time.After(timeUntilNextPoll):
			}
		}
	}()
}

func (s *service) GetPrice(ctx context.Context) (*big.Int, error) {
	callData, err := priceOracleABI.Pack("getPrice")
	if err != nil {
		return nil, err
	}
	result, err := s.transactionService.Call(ctx, &transaction.TxRequest{
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

	exchangeRate, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || exchangeRate == nil {
		return nil, errDecodeABI
	}

	return exchangeRate, nil
}

func (s *service) CurrentRates() (exchangeRate *big.Int, err error) {
	if s.exchangeRate.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("exchange rate not yet available")
	}
	return s.exchangeRate, nil
}

func (s *service) SwitchCurrentRates() (exchangeRate *big.Int, err error) {
	if s.exchangeRate.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("exchange rate not yet available")
	}

	t := math.Exp(big.NewInt(10), big.NewInt(13))
	fmt.Println("CurrentRates: s.exchangeRate ", s.exchangeRate, t)

	if s.exchangeRate.Cmp(t) > 0 {
		ret := big.NewInt(0).Div(s.exchangeRate, t)

		fmt.Println("SwitchCurrentRates: s.exchangeRate ", ret)
		return ret, nil
	}

	return s.exchangeRate, nil
}

func (s *service) Close() error {
	close(s.quitC)
	return nil
}
