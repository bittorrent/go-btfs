package vault

import (
	"context"
	"fmt"
	"math/big"

	"github.com/bittorrent/go-btfs/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type vaultContract struct {
	address            common.Address
	transactionService transaction.Service
}

func newVaultContract(address common.Address, transactionService transaction.Service) *vaultContract {
	return &vaultContract{
		address:            address,
		transactionService: transactionService,
	}
}

func (c *vaultContract) Issuer(ctx context.Context) (common.Address, error) {
	callData, err := vaultABI.Pack("issuer")
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

	results, err := vaultABI.Unpack("issuer", output)
	if err != nil {
		return common.Address{}, err
	}

	return *abi.ConvertType(results[0], new(common.Address)).(*common.Address), nil
}

// TotalBalance returns the token balance of the vault.
func (c *vaultContract) TotalBalance(ctx context.Context) (*big.Int, error) {
	callData, err := vaultABI.Pack("totalbalance")
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

	results, err := vaultABI.Unpack("totalbalance", output)
	if err != nil {
		return nil, err
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

// LiquidBalance returns the token balance of the vault sub stake amount
func (c *vaultContract) LiquidBalance(ctx context.Context) (*big.Int, error) {
	callData, err := vaultABI.Pack("liquidBalance")
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

	results, err := vaultABI.Unpack("liquidBalance", output)
	if err != nil {
		return nil, err
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

func (c *vaultContract) PaidOut(ctx context.Context, beneficiary common.Address) (*big.Int, error) {
	callData, err := vaultABI.Pack("paidOut", beneficiary)
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

	results, err := vaultABI.Unpack("paidOut", output)
	if err != nil {
		return nil, err
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

func (c *vaultContract) TotalPaidOut(ctx context.Context) (*big.Int, error) {
	callData, err := vaultABI.Pack("totalPaidOut")
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

	results, err := vaultABI.Unpack("totalPaidOut", output)
	if err != nil {
		return nil, err
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

func (c *vaultContract) SetReceiver(ctx context.Context, newReceiver common.Address) (common.Hash, error) {
	callData, err := vaultABI.Pack("setReciever", newReceiver)
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

func (c *vaultContract) Deposit(ctx context.Context, amount *big.Int) (common.Hash, error) {
	callData, err := vaultABI.Pack("deposit", amount)
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

// UpgradeTo will upgrade the vault impl to `newImpl`
func (c *vaultContract) UpgradeTo(ctx context.Context, newImpl common.Address) (err error) {
	callData, err := vaultABI.Pack("upgradeTo", newImpl)
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

func (c *vaultContract) Withdraw(ctx context.Context, amount *big.Int) (common.Hash, error) {
	callData, err := vaultABI.Pack("withdraw", amount)
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := c.transactionService.Send(ctx, &transaction.TxRequest{
		To:          &c.address,
		Data:        callData,
		Value:       big.NewInt(0),
		Description: fmt.Sprintf("vault withdrawal of %d WBTT", amount),
	})
	if err != nil {
		return hash, err
	}

	return hash, nil
}

// GetVaultImpl queries the vault implementation used for the proxy
func GetVaultImpl(ctx context.Context, vault common.Address, trxSvc transaction.Service) (vaultImpl common.Address, err error) {
	callData, err := vaultABI.Pack("implementation")
	if err != nil {
		return
	}

	output, err := trxSvc.Call(ctx, &transaction.TxRequest{
		To:   &vault,
		Data: callData,
	})
	if err != nil {
		return
	}

	results, err := vaultABI.Unpack("implementation", output)
	if err != nil {
		return
	}

	vaultImpl = *abi.ConvertType(results[0], new(common.Address)).(*common.Address)
	return
}

func _CashCheque(ctx context.Context, vault, recipient common.Address, cheque *SignedCheque, tS transaction.Service) (common.Hash, error) {
	fmt.Println("_CashCheque ", vault, recipient, cheque.CumulativePayout, cheque.Signature)

	callData, err := vaultABI.Pack("cashChequeBeneficiary", recipient, cheque.CumulativePayout, cheque.Signature)
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

func _PaidOut(ctx context.Context, vault, beneficiary common.Address, tS transaction.Service) (*big.Int, error) {
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
