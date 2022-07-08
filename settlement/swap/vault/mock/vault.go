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
	vaultBalanceFunc          func(context.Context) (*big.Int, error)
	vaultAvailableBalanceFunc func(context.Context) (*big.Int, error)
	vaultAddressFunc          func() common.Address
	vaultIssueFunc            func(ctx context.Context, beneficiary common.Address, amount *big.Int, sendChequeFunc vault.SendChequeFunc) (*big.Int, error)
	vaultWithdrawFunc         func(ctx context.Context, amount *big.Int) (hash common.Hash, err error)
	vaultDepositFunc          func(ctx context.Context, amount *big.Int) (hash common.Hash, err error)
	lastChequeFunc            func(common.Address) (*vault.SignedCheque, error)
	lastChequesFunc           func() (map[common.Address]*vault.SignedCheque, error)
	bttBalanceFunc            func(context.Context) (*big.Int, error)
	totalReceivedFunc         func() (*big.Int, error)
	totalReceivedCountFunc    func() (int, error)
}

// WithVault*Functions set the mock vault functions
func WithVaultBalanceFunc(f func(ctx context.Context) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) {
		s.vaultBalanceFunc = f
	})
}

func WithVaultAvailableBalanceFunc(f func(ctx context.Context) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) {
		s.vaultAvailableBalanceFunc = f
	})
}

func WithVaultAddressFunc(f func() common.Address) Option {
	return optionFunc(func(s *Service) {
		s.vaultAddressFunc = f
	})
}

func WithVaultDepositFunc(f func(ctx context.Context, amount *big.Int) (hash common.Hash, err error)) Option {
	return optionFunc(func(s *Service) {
		s.vaultDepositFunc = f
	})
}

func WithVaultIssueFunc(f func(ctx context.Context, beneficiary common.Address, amount *big.Int, sendChequeFunc vault.SendChequeFunc) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) {
		s.vaultIssueFunc = f
	})
}

func WithVaultWithdrawFunc(f func(ctx context.Context, amount *big.Int) (hash common.Hash, err error)) Option {
	return optionFunc(func(s *Service) {
		s.vaultWithdrawFunc = f
	})
}

func WithLastChequeFunc(f func(beneficiary common.Address) (*vault.SignedCheque, error)) Option {
	return optionFunc(func(s *Service) {
		s.lastChequeFunc = f
	})
}

func WithLastChequesFunc(f func() (map[common.Address]*vault.SignedCheque, error)) Option {
	return optionFunc(func(s *Service) {
		s.lastChequesFunc = f
	})
}

func WithTotalReceivedFunc(f func() (*big.Int, error)) Option {
	return optionFunc(func(s *Service) {
		s.totalReceivedFunc = f
	})
}

func WithTotalReceivedCountFunc(f func() (int, error)) Option {
	return optionFunc(func(s *Service) {
		s.totalReceivedCountFunc = f
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

// Balance mocks the vault .Balance function
func (s *Service) Balance(ctx context.Context) (bal *big.Int, err error) {
	if s.vaultBalanceFunc != nil {
		return s.vaultBalanceFunc(ctx)
	}
	return big.NewInt(0), errors.New("Error")
}

func (s *Service) BTTBalanceOf(ctx context.Context, add common.Address, block *big.Int) (bal *big.Int, err error) {
	if s.bttBalanceFunc != nil {
		return s.bttBalanceFunc(ctx)
	}
	return big.NewInt(0), errors.New("Error")
}

func (s *Service) CheckBalance(bal *big.Int) (err error) {
	return nil
}

func (s *Service) GetWithdrawTime(ctx context.Context) (ti *big.Int, err error) {
	return big.NewInt(0), errors.New("Error")
}

func (s *Service) LiquidBalance(ctx context.Context) (ti *big.Int, err error) {
	return big.NewInt(0), errors.New("Error")
}

func (s *Service) TotalBalance(ctx context.Context) (ti *big.Int, err error) {
	return big.NewInt(0), errors.New("Error")
}

func (s *Service) TotalIssuedCount() (ti int, err error) {
	return 0, errors.New("Error")
}

func (s *Service) TotalPaidOut(ctx context.Context) (ti *big.Int, err error) {
	return big.NewInt(0), errors.New("Error")
}

func (s *Service) WBTTBalanceOf(ctx context.Context, add common.Address) (bal *big.Int, err error) {
	if s.bttBalanceFunc != nil {
		return s.bttBalanceFunc(ctx)
	}
	return big.NewInt(0), errors.New("Error")
}

func (s *Service) AvailableBalance(ctx context.Context) (bal *big.Int, err error) {
	if s.vaultAvailableBalanceFunc != nil {
		return s.vaultAvailableBalanceFunc(ctx)
	}
	return big.NewInt(0), errors.New("Error")
}

// Deposit mocks the vault .Deposit function
func (s *Service) Deposit(ctx context.Context, amount *big.Int) (hash common.Hash, err error) {
	if s.vaultDepositFunc != nil {
		return s.vaultDepositFunc(ctx, amount)
	}
	return common.Hash{}, errors.New("Error")
}

// WaitForDeposit mocks the vault .WaitForDeposit function
func (s *Service) WaitForDeposit(ctx context.Context, txHash common.Hash) error {
	return errors.New("Error")
}

// Address mocks the vault .Address function
func (s *Service) Address() common.Address {
	if s.vaultAddressFunc != nil {
		return s.vaultAddressFunc()
	}
	return common.Address{}
}

func (s *Service) Issue(ctx context.Context, beneficiary common.Address, amount *big.Int, sendChequeFunc vault.SendChequeFunc) (*big.Int, error) {
	if s.vaultIssueFunc != nil {
		return s.vaultIssueFunc(ctx, beneficiary, amount, sendChequeFunc)
	}
	return big.NewInt(0), nil
}

func (s *Service) LastCheque(beneficiary common.Address) (*vault.SignedCheque, error) {
	if s.lastChequeFunc != nil {
		return s.lastChequeFunc(beneficiary)
	}
	return nil, errors.New("Error")
}

func (s *Service) LastCheques() (map[common.Address]*vault.SignedCheque, error) {
	if s.lastChequesFunc != nil {
		return s.lastChequesFunc()
	}
	return nil, errors.New("Error")
}

func (s *Service) Withdraw(ctx context.Context, amount *big.Int) (hash common.Hash, err error) {
	return s.vaultWithdrawFunc(ctx, amount)
}

func (s *Service) TotalIssued() (*big.Int, error) {
	return nil, errors.New("not implemented")
}

func (s *Service) TotalReceived() (*big.Int, error) {
	if s.totalReceivedFunc != nil {
		return s.totalReceivedFunc()
	}
	return nil, errors.New("totalReceivedFunc not implemented")
}

func (s *Service) TotalReceivedCount() (int, error) {
	if s.totalReceivedCountFunc != nil {
		return s.totalReceivedCountFunc()
	}
	return 0, errors.New("totalReceivedCountFunc not implemented")
}

func (s *Service) TotalReceivedCashedCount() (int, error) {
	return 0, errors.New("not implemented")
}

func (s *Service) TotalReceivedCashed() (*big.Int, error) {
	return nil, errors.New("not implemented")
}

func (s *Service) TotalDailyReceived() (*big.Int, error) {
	return nil, errors.New("not implemented")
}
func (s *Service) TotalDailyReceivedCashed() (*big.Int, error) {
	return nil, errors.New("not implemented")
}

func (s *Service) UpgradeTo(ctx context.Context, newVaultImpl common.Address) (old, new common.Address, err error) {
	return common.Address{}, common.Address{}, errors.New("not implemented")
}

// Option is the option passed to the mock Vault service
type Option interface {
	apply(*Service)
}

type optionFunc func(*Service)

func (f optionFunc) apply(r *Service) { f(r) }
