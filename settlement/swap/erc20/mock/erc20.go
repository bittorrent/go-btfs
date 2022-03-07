package mock

import (
	"context"
	"errors"
	"math/big"

	"github.com/bittorrent/go-btfs/settlement/swap/erc20"
	"github.com/ethereum/go-ethereum/common"
)

type Service struct {
	addressFunc      func(ctx context.Context) common.Address
	depositFunc      func(ctx context.Context, value *big.Int) (common.Hash, error)
	withdrawFunc     func(ctx context.Context, value *big.Int) (common.Hash, error)
	balanceOfFunc    func(ctx context.Context, address common.Address) (*big.Int, error)
	transferFunc     func(ctx context.Context, address common.Address, value *big.Int) (common.Hash, error)
	allowanceFunc    func(ctx context.Context, issuer common.Address, vault common.Address) (*big.Int, error)
	approveFunc      func(ctx context.Context, address common.Address, value *big.Int) (common.Hash, error)
	transferFromFunc func(ctx context.Context, issuer common.Address, vault common.Address, value *big.Int) (common.Hash, error)
}

func WithAddressFunc(f func(ctx context.Context) common.Address) Option {
	return optionFunc(func(s *Service) { s.addressFunc = f })
}

func WithDepositFunc(f func(ctx context.Context, value *big.Int) (common.Hash, error)) Option {
	return optionFunc(func(s *Service) { s.depositFunc = f })
}

func WithWithdrawFunc(f func(ctx context.Context, value *big.Int) (common.Hash, error)) Option {
	return optionFunc(func(s *Service) { s.withdrawFunc = f })
}

func WithAllowanceFunc(f func(ctx context.Context, issuer common.Address, vault common.Address) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) { s.allowanceFunc = f })
}

func WithApproveFunc(f func(ctx context.Context, address common.Address, value *big.Int) (common.Hash, error)) Option {
	return optionFunc(func(s *Service) { s.approveFunc = f })
}

func WithTransferFromFunc(f func(ctx context.Context, issuer common.Address, vault common.Address, value *big.Int) (common.Hash, error)) Option {
	return optionFunc(func(s *Service) { s.transferFromFunc = f })
}

func WithBalanceOfFunc(f func(ctx context.Context, address common.Address) (*big.Int, error)) Option {
	return optionFunc(func(s *Service) { s.balanceOfFunc = f })
}

func WithTransferFunc(f func(ctx context.Context, address common.Address, value *big.Int) (common.Hash, error)) Option {
	return optionFunc(func(s *Service) { s.transferFunc = f })
}

func New(opts ...Option) erc20.Service {
	mock := new(Service)
	for _, o := range opts {
		o.apply(mock)
	}
	return mock
}

func (s *Service) Address(ctx context.Context) common.Address {
	if s.addressFunc != nil {
		return s.addressFunc(ctx)
	}
	return common.Address{}
}

func (s *Service) Deposit(ctx context.Context, value *big.Int) (common.Hash, error) {
	if s.depositFunc != nil {
		return s.depositFunc(ctx, value)
	}
	return common.Hash{}, errors.New("error")
}

func (s *Service) Withdraw(ctx context.Context, value *big.Int) (common.Hash, error) {
	if s.withdrawFunc != nil {
		return s.withdrawFunc(ctx, value)
	}
	return common.Hash{}, errors.New("Error")
}

func (s *Service) Allowance(ctx context.Context, issuer common.Address, vault common.Address) (*big.Int, error) {
	if s.allowanceFunc != nil {
		return s.allowanceFunc(ctx, issuer, vault)
	}
	return big.NewInt(0), errors.New("Error")
}

func (s *Service) Approve(ctx context.Context, address common.Address, value *big.Int) (common.Hash, error) {
	if s.approveFunc != nil {
		return s.approveFunc(ctx, address, value)
	}
	return common.Hash{}, errors.New("Error")
}

func (s *Service) TransferFrom(ctx context.Context, issuer common.Address, vault common.Address, value *big.Int) (common.Hash, error) {
	if s.transferFromFunc != nil {
		return s.transferFromFunc(ctx, issuer, vault, value)
	}
	return common.Hash{}, errors.New("Error")
}

func (s *Service) BalanceOf(ctx context.Context, address common.Address) (*big.Int, error) {
	if s.balanceOfFunc != nil {
		return s.balanceOfFunc(ctx, address)
	}
	return big.NewInt(0), errors.New("Error")
}

func (s *Service) Transfer(ctx context.Context, address common.Address, value *big.Int) (common.Hash, error) {
	if s.transferFunc != nil {
		return s.transferFunc(ctx, address, value)
	}
	return common.Hash{}, errors.New("Error")
}

// Option is the option passed to the mock Chequebook service
type Option interface {
	apply(*Service)
}

type optionFunc func(*Service)

func (f optionFunc) apply(r *Service) { f(r) }
