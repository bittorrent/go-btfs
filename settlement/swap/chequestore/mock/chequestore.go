package mock

import (
	"context"
	"errors"
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

// LastReceivedCheque returns the last cheque we received from a specific vault.
func (s *Service) LastReceivedCheque(vault common.Address) (*vault.SignedCheque, error) {
	return nil, errors.New("not implemented")
}

// LastReceivedCheques return map[vault]cheque
func (s *Service) LastReceivedCheques() (map[common.Address]*vault.SignedCheque, error) {
	return nil, errors.New("not implemented")
}

// ReceivedChequeRecordsByPeer returns the records we received from a specific vault.
func (s *Service) ReceivedChequeRecordsByPeer(vault common.Address) ([]vault.ChequeRecord, error) {
	return nil, errors.New("not implemented")
}

// ListReceivedChequeRecords returns the records we received from a specific vault.
func (s *Service) ReceivedChequeRecordsAll() (map[common.Address][]vault.ChequeRecord, error) {
	return nil, errors.New("not implemented")
}
func (s *Service) ReceivedStatsHistory(days int) ([]vault.DailyReceivedStats, error) {
	return nil, errors.New("not implemented")
}
func (s *Service) SentStatsHistory(days int) ([]vault.DailySentStats, error) {
	return nil, errors.New("not implemented")
}

// StoreSendChequeRecord store send cheque records.
func (s *Service) StoreSendChequeRecord(vault, beneficiary common.Address, amount *big.Int) error {
	return errors.New("not implemented")
}

// SendChequeRecordsByPeer returns the records we send to a specific vault.
func (s *Service) SendChequeRecordsByPeer(beneficiary common.Address) ([]vault.ChequeRecord, error) {
	return nil, errors.New("not implemented")
}

// SendChequeRecordsAll returns the records we send to a specific vault.
func (s *Service) SendChequeRecordsAll() (map[common.Address][]vault.ChequeRecord, error) {
	return nil, errors.New("not implemented")
}

// Option is the option passed to the mock ChequeStore service
type Option interface {
	apply(*Service)
}

type optionFunc func(*Service)

func (f optionFunc) apply(r *Service) { f(r) }
