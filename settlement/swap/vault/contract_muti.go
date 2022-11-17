package vault

import (
	"context"
	"fmt"
	"math/big"

	"github.com/bittorrent/go-btfs/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type vaultContractMuti struct {
	address            common.Address
	transactionService transaction.Service
	contractWBTT       *vaultContract
}

func newVaultContractMuti(address common.Address, transactionService transaction.Service) *vaultContractMuti {
	return &vaultContractMuti{
		address:            address,
		transactionService: transactionService,
		contractWBTT:       newVaultContract(address, transactionService),
	}
}

// Issuer (all the same)
func (c *vaultContractMuti) Issuer(ctx context.Context) (common.Address, error) {
	callData, err := vaultABINew.Pack("issuer")
	if err != nil {
		return common.Address{}, err
	}

	output, err := c.transactionService.Call(ctx, &transaction.TxRequest{
		To:   &c.address,
		Data: callData,
	})
	if err != nil {
		return common.Address{}, err
	}

	results, err := vaultABINew.Unpack("issuer", output)
	if err != nil {
		return common.Address{}, err
	}

	return *abi.ConvertType(results[0], new(common.Address)).(*common.Address), nil
}

// TotalBalance returns the token balance of the vault.
func (c *vaultContractMuti) TotalBalance(ctx context.Context, token string) (*big.Int, error) {
	if IsWbtt(token) {
		return c.contractWBTT.TotalBalance(ctx)
	}

	callData, err := vaultABINew.Pack("totalbalance")
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

	results, err := vaultABINew.Unpack("totalbalance", output)
	if err != nil {
		return nil, err
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

// LiquidBalance returns the token balance of the vault sub stake amount (not use)
func (c *vaultContractMuti) LiquidBalance(ctx context.Context) (*big.Int, error) {
	callData, err := vaultABINew.Pack("liquidBalance")
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

	results, err := vaultABINew.Unpack("liquidBalance", output)
	if err != nil {
		return nil, err
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

func (c *vaultContractMuti) PaidOut(ctx context.Context, address common.Address, token string) (*big.Int, error) {
	if IsWbtt(token) {
		return c.contractWBTT.TotalBalance(ctx)
	}

	callData, err := vaultABINew.Pack("paidOut", address)
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

	results, err := vaultABINew.Unpack("paidOut", output)
	if err != nil {
		return nil, err
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

func (c *vaultContractMuti) TotalPaidOut(ctx context.Context, token string) (*big.Int, error) {
	if IsWbtt(token) {
		return c.contractWBTT.TotalBalance(ctx)
	}

	callData, err := vaultABINew.Pack("totalPaidOut")
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

	results, err := vaultABINew.Unpack("totalPaidOut", output)
	if err != nil {
		return nil, err
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

// SetReceiver (not use)
func (c *vaultContractMuti) SetReceiver(ctx context.Context, newReceiver common.Address) (common.Hash, error) {
	callData, err := vaultABINew.Pack("setReciever", newReceiver)
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := c.transactionService.Send(ctx, &transaction.TxRequest{
		To:   &c.address,
		Data: callData,
	})
	if err != nil {
		return hash, err
	}

	return hash, nil
}

func (c *vaultContractMuti) Deposit(ctx context.Context, amount *big.Int, token string) (common.Hash, error) {
	if IsWbtt(token) {
		return c.contractWBTT.Deposit(ctx, amount)
	}

	callData, err := vaultABINew.Pack("deposit", amount)
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := c.transactionService.Send(ctx, &transaction.TxRequest{
		To:   &c.address,
		Data: callData,
	})
	if err != nil {
		return hash, err
	}

	return hash, nil
}

func (c *vaultContractMuti) Withdraw(ctx context.Context, amount *big.Int, token string) (common.Hash, error) {
	if IsWbtt(token) {
		return c.contractWBTT.Withdraw(ctx, amount)
	}

	callData, err := vaultABINew.Pack("withdraw", amount)
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := c.transactionService.Send(ctx, &transaction.TxRequest{
		To:          &c.address,
		Data:        callData,
		Value:       big.NewInt(0),
		Description: fmt.Sprintf("vault withdrawal of %d [%s]", amount, token),
	})
	if err != nil {
		return hash, err
	}

	return hash, nil
}

// UpgradeTo will upgrade the vault impl to `newImpl` (all the same)
func (c *vaultContractMuti) UpgradeTo(ctx context.Context, newImpl common.Address) (err error) {
	callData, err := vaultABINew.Pack("upgradeTo", newImpl)
	if err != nil {
		return
	}
	txHash, err := c.transactionService.Send(ctx, &transaction.TxRequest{
		To:   &c.address,
		Data: callData,
	})
	if err != nil {
		return
	}

	// wait tx finish
	receipt, err := c.transactionService.WaitForReceipt(ctx, txHash)
	if err != nil {
		return
	}
	if receipt.Status != 1 {
		return transaction.ErrTransactionReverted
	}
	return
}

func _CashChequeMuti(ctx context.Context, vault, recipient common.Address, cheque *SignedCheque, tS transaction.Service, token string) (common.Hash, error) {
	if IsWbtt(token) {
		return _CashCheque(ctx, vault, recipient, cheque, tS)
	}

	callData, err := vaultABINew.Pack("cashChequeBeneficiary", recipient, cheque.CumulativePayout, cheque.Signature)
	if err != nil {
		return common.Hash{}, err
	}
	request := &transaction.TxRequest{
		To:          &vault,
		Data:        callData,
		Value:       big.NewInt(0),
		Description: "cheque cashout",
	}

	txHash, err := tS.Send(ctx, request)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func _PaidOutMuti(ctx context.Context, vault, beneficiary common.Address, tS transaction.Service, token string) (*big.Int, error) {
	if IsWbtt(token) {
		return _PaidOut(ctx, vault, beneficiary, tS)
	}
	
	callData, err := vaultABINew.Pack("paidOut", beneficiary)
	if err != nil {
		return nil, err
	}

	output, err := tS.Call(ctx, &transaction.TxRequest{
		To:   &vault,
		Data: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := vaultABINew.Unpack("paidOut", output)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, errDecodeABI
	}

	paidOut, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || paidOut == nil {
		return nil, errDecodeABI
	}

	return paidOut, nil
}
