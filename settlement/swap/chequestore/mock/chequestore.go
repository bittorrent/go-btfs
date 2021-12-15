package mock

import (
	"context"
	"math/big"

	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	"github.com/ethereum/go-ethereum/common"
)

// Service is the mock chequeStore service.
type Service struct {
	receiveCheque func(ctx context.Context, cheque *vault.SignedCheque, exchangeRate *big.Int) (*big.Int, error)
	lastCheque    func(vault common.Address) (*vault.SignedCheque, error)
	lastCheques   func() (map[common.Address]*vault.SignedCheque, error)
}

func WithReceiveChequeFunc(f func(ctx context.Context, cheque *vault.SignedCheque, exchangeRate *big.Int) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) {
		s.receiveCheque = f
	})
}

func WithLastChequeFunc(f func(vault common.Address) (*vault.SignedCheque, error)) Option {
	return optionFunc(func(s *Service) {
		s.lastCheque = f
	})
}

func WithLastChequesFunc(f func() (map[common.Address]*vault.SignedCheque, error)) Option {
	return optionFunc(func(s *Service) {
		s.lastCheques = f
	})
}

// NewChequeStore creates the mock chequeStore implementation
func NewChequeStore(opts ...Option) vault.ChequeStore {
	mock := new(Service)
	for _, o := range opts {
		o.apply(mock)
	}
	return mock
}

func (s *Service) ReceiveCheque(ctx context.Context, cheque *vault.SignedCheque, exchangeRate *big.Int) (*big.Int, error) {
	return s.receiveCheque(ctx, cheque, exchangeRate)
}

func (s *Service) LastCheque(vault common.Address) (*vault.SignedCheque, error) {
	return s.lastCheque(vault)
}

func (s *Service) LastCheques() (map[common.Address]*vault.SignedCheque, error) {
	return s.lastCheques()
}

// Option is the option passed to the mock ChequeStore service
type Option interface {
	apply(*Service)
}

type optionFunc func(*Service)

func (f optionFunc) apply(r *Service) { f(r) }
