package erc20

import (
	"context"
	"errors"
	"math/big"

	conabi "github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	erc20ABI     = transaction.ParseABIUnchecked(conabi.Erc20ABI)
	errDecodeABI = errors.New("could not decode abi data")
)

type Service interface {
	Address(ctx context.Context) common.Address
	BalanceOf(ctx context.Context, address common.Address) (*big.Int, error)
	Deposit(ctx context.Context, value *big.Int) (common.Hash, error)
	Withdraw(ctx context.Context, value *big.Int) (common.Hash, error)
	Transfer(ctx context.Context, address common.Address, value *big.Int) (common.Hash, error)
	Allowance(ctx context.Context, issuer common.Address, vault common.Address) (*big.Int, error)
	Approve(ctx context.Context, address common.Address, value *big.Int) (common.Hash, error)
	TransferFrom(ctx context.Context, issuer common.Address, vault common.Address, value *big.Int) (common.Hash, error)
}

type erc20Service struct {
	backend            transaction.Backend
	transactionService transaction.Service
	address            common.Address
}

func New(backend transaction.Backend, transactionService transaction.Service, address common.Address) Service {
	return &erc20Service{
		backend:            backend,
		transactionService: transactionService,
		address:            address,
	}
}

func (c *erc20Service) Address(ctx context.Context) common.Address {
	return c.address
}

func (c *erc20Service) BalanceOf(ctx context.Context, address common.Address) (*big.Int, error) {
	callData, err := erc20ABI.Pack("balanceOf", address)
	if err != nil {
		return nil, err
	}

	output, err := c.transactionService.Call(ctx, &transaction.TxRequest{
		To:   &c.address,
		Data: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := erc20ABI.Unpack("balanceOf", output)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, errDecodeABI
	}

	balance, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || balance == nil {
		return nil, errDecodeABI
	}
	return balance, nil
}

func (c *erc20Service) Deposit(ctx context.Context, value *big.Int) (trx common.Hash, err error) {
	callData, err := erc20ABI.Pack("deposit")
	if err != nil {
		return
	}
	req := &transaction.TxRequest{
		To:          &c.address,
		Data:        callData,
		Value:       value,
		Description: "deposit wbtt",
	}
	trx, err = c.transactionService.Send(ctx, req)
	return
}

func (c *erc20Service) Withdraw(ctx context.Context, value *big.Int) (trx common.Hash, err error) {
	callData, err := erc20ABI.Pack("withdraw", value)
	if err != nil {
		return
	}
	req := &transaction.TxRequest{
		To:          &c.address,
		Data:        callData,
		Value:       big.NewInt(0),
		Description: "withdraw wbtt",
	}
	trx, err = c.transactionService.Send(ctx, req)
	return
}

func (c *erc20Service) Transfer(ctx context.Context, address common.Address, value *big.Int) (common.Hash, error) {
	callData, err := erc20ABI.Pack("transfer", address, value)
	if err != nil {
		return common.Hash{}, err
	}

	request := &transaction.TxRequest{
		To:          &c.address,
		Data:        callData,
		Value:       big.NewInt(0),
		Description: "token transfer",
	}

	txHash, err := c.transactionService.Send(ctx, request)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

// first, approve to vault
func (c *erc20Service) TransferFrom(ctx context.Context, issuer common.Address, vault common.Address, value *big.Int) (common.Hash, error) {
	callData, err := erc20ABI.Pack("transferFrom", issuer, vault, value)
	if err != nil {
		return common.Hash{}, err
	}

	request := &transaction.TxRequest{
		To:          &c.address,
		Data:        callData,
		Value:       big.NewInt(0),
		Description: "token transfer",
	}

	txHash, err := c.transactionService.Send(ctx, request)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (c *erc20Service) Allowance(ctx context.Context, issuer common.Address, vault common.Address) (*big.Int, error) {
	callData, err := erc20ABI.Pack("allowance", issuer, vault)
	if err != nil {
		return nil, err
	}

	output, err := c.transactionService.Call(ctx, &transaction.TxRequest{
		To:   &c.address,
		Data: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := erc20ABI.Unpack("allowance", output)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, errDecodeABI
	}

	allowance, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || allowance == nil {
		return nil, errDecodeABI
	}
	return allowance, nil
}

func (c *erc20Service) Approve(ctx context.Context, address common.Address, value *big.Int) (common.Hash, error) {
	callData, err := erc20ABI.Pack("approve", address, value)
	if err != nil {
		return common.Hash{}, err
	}

	request := &transaction.TxRequest{
		To:          &c.address,
		Data:        callData,
		Value:       big.NewInt(0),
		Description: "approve",
	}

	txHash, err := c.transactionService.Send(ctx, request)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}
