package mock

import (
	"context"
	"errors"
	"math/big"

	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	"github.com/ethereum/go-ethereum/common"
)

// Service is the mock vault service.
type Service struct {
	checkBalanceFunc             func(bal *big.Int) (err error)
	getWithdrawTimeFunc          func(ctx context.Context) (ti *big.Int, err error)
	liquidBalanceFunc            func(ctx context.Context) (ti *big.Int, err error)
	totalBalanceFunc             func(ctx context.Context, token common.Address) (ti *big.Int, err error)
	totalIssuedCountFunc         func(token common.Address) (ti int, err error)
	totalPaidOutFunc             func(ctx context.Context, token common.Address) (ti *big.Int, err error)
	wbttBalanceFunc              func(ctx context.Context, add common.Address) (bal *big.Int, err error)
	waitForDepositFunc           func(ctx context.Context, txHash common.Hash) error
	totalIssuedFunc              func(common.Address) (*big.Int, error)
	totalReceivedCashedCountFunc func(token common.Address) (int, error)
	totalReceivedCashedFunc      func(token common.Address) (*big.Int, error)
	totalDailyReceivedFunc       func(common.Address) (*big.Int, error)
	totalDailyReceivedCashedFunc func(common.Address) (*big.Int, error)

	vaultBalanceFunc          func(context.Context) (*big.Int, error)
	vaultAvailableBalanceFunc func(context.Context, common.Address) (*big.Int, error)
	vaultAddressFunc          func() common.Address
	vaultIssueFunc            func(ctx context.Context, beneficiary common.Address, amount *big.Int, token common.Address, sendChequeFunc vault.SendChequeFunc) (*big.Int, error)
	vaultWithdrawFunc         func(ctx context.Context, amount *big.Int, token common.Address) (hash common.Hash, err error)
	vaultDepositFunc          func(ctx context.Context, amount *big.Int, token common.Address) (hash common.Hash, err error)
	lastChequeFunc            func(common.Address, common.Address) (*vault.SignedCheque, error)
	lastChequesFunc           func(common.Address) (map[common.Address]*vault.SignedCheque, error)
	bttBalanceFunc            func(context.Context) (*big.Int, error)
	totalReceivedFunc         func(token common.Address) (*big.Int, error)
	totalReceivedCountFunc    func(token common.Address) (int, error)
}

// WithVault*Functions set the mock vault functions
func WithVaultBalanceFunc(f func(ctx context.Context) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) {
		s.vaultBalanceFunc = f
	})
}

func WithVaultAvailableBalanceFunc(f func(ctx context.Context, token common.Address) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) {
		s.vaultAvailableBalanceFunc = f
	})
}

func WithVaultAddressFunc(f func() common.Address) Option {
	return optionFunc(func(s *Service) {
		s.vaultAddressFunc = f
	})
}

func WithVaultDepositFunc(f func(ctx context.Context, amount *big.Int, token common.Address) (hash common.Hash, err error)) Option {
	return optionFunc(func(s *Service) {
		s.vaultDepositFunc = f
	})
}

func WithVaultIssueFunc(f func(ctx context.Context, beneficiary common.Address, amount *big.Int, token common.Address, sendChequeFunc vault.SendChequeFunc) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) {
		s.vaultIssueFunc = f
	})
}

func WithVaultWithdrawFunc(f func(ctx context.Context, amount *big.Int, token common.Address) (hash common.Hash, err error)) Option {
	return optionFunc(func(s *Service) {
		s.vaultWithdrawFunc = f
	})
}

func WithLastChequeFunc(f func(beneficiary common.Address, token common.Address) (*vault.SignedCheque, error)) Option {
	return optionFunc(func(s *Service) {
		s.lastChequeFunc = f
	})
}

func WithLastChequesFunc(f func(token common.Address) (map[common.Address]*vault.SignedCheque, error)) Option {
	return optionFunc(func(s *Service) {
		s.lastChequesFunc = f
	})
}

func WithTotalReceivedFunc(f func(token common.Address) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) {
		s.totalReceivedFunc = f
	})
}

func WithTotalReceivedCountFunc(f func(token common.Address) (int, error)) Option {
	return optionFunc(func(s *Service) {
		s.totalReceivedCountFunc = f
	})
}

func WithCheckBalanceFunc(f func(bal *big.Int) (err error)) Option {
	return optionFunc(func(s *Service) {
		s.checkBalanceFunc = f
	})
}

func WithGetWithdrawTimeFunc(f func(ctx context.Context) (ti *big.Int, err error)) Option {
	return optionFunc(func(s *Service) {
		s.getWithdrawTimeFunc = f
	})
}

func WithLiquidBalanceFunc(f func(ctx context.Context) (ti *big.Int, err error)) Option {
	return optionFunc(func(s *Service) {
		s.liquidBalanceFunc = f
	})
}

func WithTotalBalanceFunc(f func(ctx context.Context, token common.Address) (ti *big.Int, err error)) Option {
	return optionFunc(func(s *Service) {
		s.totalBalanceFunc = f
	})
}

func WithTotalIssuedCountFunc(f func(token common.Address) (int, error)) Option {
	return optionFunc(func(s *Service) {
		s.totalIssuedCountFunc = f
	})
}

func WithTotalPaidOutFunc(f func(ctx context.Context, token common.Address) (ti *big.Int, err error)) Option {
	return optionFunc(func(s *Service) {
		s.totalPaidOutFunc = f
	})
}

func WithWbttBalanceFunc(f func(ctx context.Context, add common.Address) (bal *big.Int, err error)) Option {
	return optionFunc(func(s *Service) {
		s.wbttBalanceFunc = f
	})
}

func WithWaitForDepositFunc(f func(ctx context.Context, txHash common.Hash) error) Option {
	return optionFunc(func(s *Service) {
		s.waitForDepositFunc = f
	})
}

func WithTotalIssuedFunc(f func(token common.Address) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) {
		s.totalIssuedFunc = f
	})
}

func WithTotalReceivedCashedCountFunc(f func(token common.Address) (int, error)) Option {
	return optionFunc(func(s *Service) {
		s.totalReceivedCashedCountFunc = f
	})
}

func WithTotalReceivedCashedFunc(f func(token common.Address) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) {
		s.totalReceivedCashedFunc = f
	})
}

func WithTotalDailyReceivedFunc(f func(token common.Address) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) {
		s.totalDailyReceivedFunc = f
	})
}

func WithTotalDailyReceivedCashedFunc(f func(token common.Address) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) {
		s.totalDailyReceivedCashedFunc = f
	})
}

// NewVault creates the mock vault implementation
func NewVault(opts ...Option) vault.Service {
	mock := new(Service)
	for _, o := range opts {
		o.apply(mock)
	}
	return mock
}

func (s *Service) TokenBalanceOf(ctx context.Context, addr common.Address, tokenStr string) (*big.Int, error) {
	return nil, errors.New("vaultMock.TokenBalanceOf not implemented")
}

// Balance mocks the vault .Balance function
func (s *Service) Balance(ctx context.Context) (bal *big.Int, err error) {
	if s.vaultBalanceFunc != nil {
		return s.vaultBalanceFunc(ctx)
	}
	return nil, errors.New("vaultMock.vaultBalanceFunc not implemented")
}

func (s *Service) BTTBalanceOf(ctx context.Context, add common.Address, block *big.Int) (bal *big.Int, err error) {
	if s.bttBalanceFunc != nil {
		return s.bttBalanceFunc(ctx)
	}
	return nil, errors.New("vaultMock.bttBalanceFunc not implemented")
}

func (s *Service) CheckBalance(bal *big.Int) (err error) {
	if s.checkBalanceFunc != nil {
		return s.checkBalanceFunc(bal)
	}
	return errors.New("vaultMock.checkBalanceFunc not implemented")
}

func (s *Service) GetWithdrawTime(ctx context.Context) (ti *big.Int, err error) {
	if s.getWithdrawTimeFunc != nil {
		return s.getWithdrawTimeFunc(ctx)
	}
	return nil, errors.New("vaultMock.getWithdrawTimeFunc not implemented")
}

func (s *Service) LiquidBalance(ctx context.Context) (ti *big.Int, err error) {
	if s.liquidBalanceFunc != nil {
		return s.liquidBalanceFunc(ctx)
	}
	return nil, errors.New("vaultMock.liquidBalanceFunc not implemented")
}

func (s *Service) TotalBalance(ctx context.Context, token common.Address) (ti *big.Int, err error) {
	if s.totalBalanceFunc != nil {
		return s.totalBalanceFunc(ctx, token)
	}
	return nil, errors.New("vaultMock.totalBalanceFunc not implemented")
}

func (s *Service) TotalIssuedCount(token common.Address) (ti int, err error) {
	if s.totalIssuedCountFunc != nil {
		return s.totalIssuedCountFunc(token)
	}
	return 0, errors.New("vaultMock.totalIssuedCountFunc not implemented")
}

func (s *Service) TotalPaidOut(ctx context.Context, token common.Address) (ti *big.Int, err error) {
	if s.totalPaidOutFunc != nil {
		return s.totalPaidOutFunc(ctx, token)
	}
	return nil, errors.New("vaultMock.totalPaidOutFunc not implemented")
}

func (s *Service) WBTTBalanceOf(ctx context.Context, add common.Address) (bal *big.Int, err error) {
	if s.wbttBalanceFunc != nil {
		return s.wbttBalanceFunc(ctx, add)
	}
	return nil, errors.New("vaultMock.wbttBalanceFunc not implemented")
}

func (s *Service) AvailableBalance(ctx context.Context, token common.Address) (bal *big.Int, err error) {
	if s.vaultAvailableBalanceFunc != nil {
		return s.vaultAvailableBalanceFunc(ctx, token)
	}
	return nil, errors.New("vaultMock.vaultAvailableBalanceFunc not implemented")
}

// Deposit mocks the vault .Deposit function
func (s *Service) Deposit(ctx context.Context, amount *big.Int, token common.Address) (hash common.Hash, err error) {
	if s.vaultDepositFunc != nil {
		return s.vaultDepositFunc(ctx, amount, token)
	}
	return common.Hash{}, errors.New("vaultMock.vaultDepositFunc not implemented")
}

func (s *Service) WaitForDeposit(ctx context.Context, txHash common.Hash) error {
	if s.waitForDepositFunc != nil {
		return s.waitForDepositFunc(ctx, txHash)
	}
	return errors.New("vaultMock.waitForDepositFunc not implemented")
}

// Address mocks the vault .Address function
func (s *Service) Address() common.Address {
	if s.vaultAddressFunc != nil {
		return s.vaultAddressFunc()
	}
	return common.Address{}
}

func (s *Service) Issue(ctx context.Context, beneficiary common.Address, amount *big.Int, token common.Address, sendChequeFunc vault.SendChequeFunc) (*big.Int, error) {
	if s.vaultIssueFunc != nil {
		return s.vaultIssueFunc(ctx, beneficiary, amount, token, sendChequeFunc)
	}
	return nil, errors.New("vaultMock.vaultIssueFunc not implemented")
}

func (s *Service) LastCheque(beneficiary common.Address, token common.Address) (*vault.SignedCheque, error) {
	if s.lastChequeFunc != nil {
		return s.lastChequeFunc(beneficiary, token)
	}
	return nil, errors.New("vaultMock.lastChequeFunc not implemented")
}

func (s *Service) LastCheques(token common.Address) (map[common.Address]*vault.SignedCheque, error) {
	if s.lastChequesFunc != nil {
		return s.lastChequesFunc(token)
	}
	return nil, errors.New("vaultMock.lastChequesFunc not implemented")
}

func (s *Service) Withdraw(ctx context.Context, amount *big.Int, token common.Address) (hash common.Hash, err error) {
	if s.vaultWithdrawFunc != nil {
		return s.vaultWithdrawFunc(ctx, amount, token)
	}
	return common.Hash{}, errors.New("vaultMock.vaultWithdrawFunc not implemented")
}

func (s *Service) TotalIssued(token common.Address) (*big.Int, error) {
	if s.totalIssuedFunc != nil {
		return s.totalIssuedFunc(token)
	}
	return nil, errors.New("vaultMock.totalIssuedFunc not implemented")
}

func (s *Service) TotalReceived(token common.Address) (*big.Int, error) {
	if s.totalReceivedFunc != nil {
		return s.totalReceivedFunc(token)
	}
	return nil, errors.New("vaultMock.totalReceivedFunc not implemented")
}

func (s *Service) TotalReceivedCount(token common.Address) (int, error) {
	if s.totalReceivedCountFunc != nil {
		return s.totalReceivedCountFunc(token)
	}
	return 0, errors.New("vaultMock.totalReceivedCountFunc not implemented")
}

func (s *Service) TotalReceivedCashedCount(token common.Address) (int, error) {
	if s.totalReceivedCashedCountFunc != nil {
		return s.totalReceivedCashedCountFunc(token)
	}
	return 0, errors.New("vaultMock.totalReceivedCashedCountFunc not implemented")
}

func (s *Service) TotalReceivedCashed(token common.Address) (*big.Int, error) {
	if s.totalReceivedCashedFunc != nil {
		return s.totalReceivedCashedFunc(token)
	}
	return nil, errors.New("vaultMock.totalReceivedCashedFunc not implemented")
}

func (s *Service) TotalDailyReceived(token common.Address) (*big.Int, error) {
	if s.totalDailyReceivedFunc != nil {
		return s.totalDailyReceivedFunc(token)
	}
	return nil, errors.New("vaultMock.totalDailyReceivedFunc not implemented")
}
func (s *Service) TotalDailyReceivedCashed(token common.Address) (*big.Int, error) {
	if s.totalDailyReceivedCashedFunc != nil {
		return s.totalDailyReceivedCashedFunc(token)
	}
	return nil, errors.New("vaultMock.totalDailyReceivedCashedFunc not implemented")
}

func (s *Service) UpgradeTo(ctx context.Context, newVaultImpl common.Address) (old, new common.Address, err error) {
	return common.Address{}, common.Address{}, errors.New("vaultMock.UpgradeTo not implemented")
}

// Option is the option passed to the mock Vault service
type Option interface {
	apply(*Service)
}

type optionFunc func(*Service)

func (f optionFunc) apply(r *Service) { f(r) }
