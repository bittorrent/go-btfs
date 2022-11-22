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

	lastReceivedChequeFunc          func(vault common.Address) (*vault.SignedCheque, error)
	lastReceivedChequesFunc         func() (map[common.Address]*vault.SignedCheque, error)
	receivedChequeRecordsByPeerFunc func(vault common.Address) ([]vault.ChequeRecord, error)
	receivedChequeRecordsAllFunc    func() (map[common.Address][]vault.ChequeRecord, error)
	receivedStatsHistoryFunc        func(days int) ([]vault.DailyReceivedStats, error)
	sentStatsHistoryFunc            func(days int) ([]vault.DailySentStats, error)
	storeSendChequeRecordFunc       func(vault, beneficiary common.Address, amount *big.Int) error
	sendChequeRecordsByPeerFunc     func(beneficiary common.Address) ([]vault.ChequeRecord, error)
	sendChequeRecordsAllFunc        func() (map[common.Address][]vault.ChequeRecord, error)
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

func WithLastReceivedChequeFunc(f func(vault common.Address) (*vault.SignedCheque, error)) Option {
	return optionFunc(func(s *Service) {
		s.lastReceivedChequeFunc = f
	})
}

func WithLastReceivedChequesFunc(f func() (map[common.Address]*vault.SignedCheque, error)) Option {
	return optionFunc(func(s *Service) {
		s.lastReceivedChequesFunc = f
	})
}

func WithReceivedChequeRecordsByPeerFunc(f func(vault common.Address) ([]vault.ChequeRecord, error)) Option {
	return optionFunc(func(s *Service) {
		s.receivedChequeRecordsByPeerFunc = f
	})
}

func WithReceivedChequeRecordsAllFunc(f func() (map[common.Address][]vault.ChequeRecord, error)) Option {
	return optionFunc(func(s *Service) {
		s.receivedChequeRecordsAllFunc = f
	})
}

func WithReceivedStatsHistoryFunc(f func(days int) ([]vault.DailyReceivedStats, error)) Option {
	return optionFunc(func(s *Service) {
		s.receivedStatsHistoryFunc = f
	})
}

func WithSentStatsHistoryFunc(f func(days int) ([]vault.DailySentStats, error)) Option {
	return optionFunc(func(s *Service) {
		s.sentStatsHistoryFunc = f
	})
}

func WithStoreSendChequeRecordFunc(f func(vault, beneficiary common.Address, amount *big.Int) error) Option {
	return optionFunc(func(s *Service) {
		s.storeSendChequeRecordFunc = f
	})
}

func WithSendChequeRecordsByPeerFunc(f func(beneficiary common.Address) ([]vault.ChequeRecord, error)) Option {
	return optionFunc(func(s *Service) {
		s.sendChequeRecordsByPeerFunc = f
	})
}

func WithSendChequeRecordsAllFunc(f func() (map[common.Address][]vault.ChequeRecord, error)) Option {
	return optionFunc(func(s *Service) {
		s.sendChequeRecordsAllFunc = f
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

func (s *Service) ReceiveCheque(ctx context.Context, cheque *vault.SignedCheque, exchangeRate *big.Int, token common.Address) (*big.Int, error) {
	if s.receiveCheque != nil {
		return s.receiveCheque(ctx, cheque, exchangeRate)
	}
	return nil, errors.New("checkstoreMock.receiveCheque not implemented")
}

func (s *Service) LastCheque(vault common.Address) (*vault.SignedCheque, error) {
	if s.lastCheque != nil {
		return s.lastCheque(vault)
	}
	return nil, errors.New("checkstoreMock.lastCheque not implemented")
}

func (s *Service) LastCheques() (map[common.Address]*vault.SignedCheque, error) {
	if s.lastCheques != nil {
		return s.lastCheques()
	}
	return nil, errors.New("checkstoreMock.lastCheques not implemented")
}

// LastReceivedCheque returns the last cheque we received from a specific vault.
func (s *Service) LastReceivedCheque(vault common.Address, token common.Address) (*vault.SignedCheque, error) {
	if s.lastReceivedChequeFunc != nil {
		return s.lastReceivedChequeFunc(vault)
	}
	return nil, errors.New("checkstoreMock.lastReceivedChequeFunc not implemented")
}

// LastReceivedCheques return map[vault]cheque
func (s *Service) LastReceivedCheques(token common.Address) (map[common.Address]*vault.SignedCheque, error) {
	if s.lastReceivedChequesFunc != nil {
		return s.lastReceivedChequesFunc()
	}
	return nil, errors.New("checkstoreMock.lastReceivedChequesFunc not implemented")
}

// ReceivedChequeRecordsByPeer returns the records we received from a specific vault.
func (s *Service) ReceivedChequeRecordsByPeer(vault common.Address, token common.Address) ([]vault.ChequeRecord, error) {
	if s.receivedChequeRecordsByPeerFunc != nil {
		return s.receivedChequeRecordsByPeerFunc(vault)
	}
	return nil, errors.New("checkstoreMock.receivedChequeRecordsByPeerFunc not implemented")
}

// ListReceivedChequeRecords returns the records we received from a specific vault.
func (s *Service) ReceivedChequeRecordsAll(token common.Address) (map[common.Address][]vault.ChequeRecord, error) {
	if s.receivedChequeRecordsAllFunc != nil {
		return s.receivedChequeRecordsAllFunc()
	}
	return nil, errors.New("checkstoreMock.receivedChequeRecordsAllFunc not implemented")
}

func (s *Service) ReceivedStatsHistory(days int, token common.Address) ([]vault.DailyReceivedStats, error) {
	if s.receivedStatsHistoryFunc != nil {
		return s.receivedStatsHistoryFunc(days)
	}
	return nil, errors.New("checkstoreMock.receivedStatsHistoryFunc not implemented")
}

func (s *Service) SentStatsHistory(days int, token common.Address) ([]vault.DailySentStats, error) {
	if s.sentStatsHistoryFunc != nil {
		return s.sentStatsHistoryFunc(days)
	}
	return nil, errors.New("checkstoreMock.sentStatsHistoryFunc not implemented")
}

// StoreSendChequeRecord store send cheque records.
func (s *Service) StoreSendChequeRecord(vault, beneficiary common.Address, amount *big.Int, token common.Address) error {
	if s.storeSendChequeRecordFunc != nil {
		return s.storeSendChequeRecordFunc(vault, beneficiary, amount)
	}
	return errors.New("checkstoreMock.storeSendChequeRecordFunc not implemented")
}

// SendChequeRecordsByPeer returns the records we send to a specific vault.
func (s *Service) SendChequeRecordsByPeer(beneficiary common.Address, token common.Address) ([]vault.ChequeRecord, error) {
	if s.sendChequeRecordsByPeerFunc != nil {
		return s.sendChequeRecordsByPeerFunc(beneficiary)
	}
	return nil, errors.New("checkstoreMock.sendChequeRecordsByPeerFunc not implemented")
}

// SendChequeRecordsAll returns the records we send to a specific vault.
func (s *Service) SendChequeRecordsAll(token common.Address) (map[common.Address][]vault.ChequeRecord, error) {
	if s.sendChequeRecordsAllFunc != nil {
		return s.sendChequeRecordsAllFunc()
	}
	return nil, errors.New("checkstoreMock.sendChequeRecordsAllFunc not implemented")
}

// Option is the option passed to the mock ChequeStore service
type Option interface {
	apply(*Service)
}

type optionFunc func(*Service)

func (f optionFunc) apply(r *Service) { f(r) }
