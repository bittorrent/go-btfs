package vault

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	conabi "github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/net/context"
)

var (
	ErrInvalidFactory       = errors.New("not a valid factory contract")
	ErrNotDeployedByFactory = errors.New("vault not deployed by factory")
	errDecodeABI            = errors.New("could not decode abi data")

	factoryABI             = transaction.ParseABIUnchecked(conabi.VaultFactoryABI)
	vaultDeployedEventType = factoryABI.Events["VaultDeployed"]
)

// Factory is the main interface for interacting with the vault factory.
type Factory interface {
	// ERC20Address returns the token for which this factory deploys vaults.
	ERC20Address(ctx context.Context) (common.Address, error)
	// Deploy deploys a new vault and returns once the transaction has been submitted.
	Deploy(ctx context.Context, issuer common.Address, nonce common.Hash, peerId string) (common.Hash, error)
	// WaitDeployed waits for the deployment transaction to confirm and returns the vault address
	WaitDeployed(ctx context.Context, txHash common.Hash) (common.Address, error)
	// VerifyBytecode checks that the factory is valid.
	VerifyBytecode(ctx context.Context) error
	// VerifyVault checks that the supplied vault has been deployed by this factory.
	VerifyVault(ctx context.Context, vault common.Address) error
}

type factory struct {
	backend            transaction.Backend
	transactionService transaction.Service
	address            common.Address // address of the factory to use for deployments
}

type vaultDeployedEvent struct {
	Issuer          common.Address
	ContractAddress common.Address
	Id              string
}

// the bytecode of factories which can be used for deployment
var currentDeployVersion []byte = common.FromHex(conabi.FactoryDeployedBin)

// NewFactory creates a new factory service for the provided factory contract.
func NewFactory(backend transaction.Backend, transactionService transaction.Service, address common.Address) Factory {
	return &factory{
		backend:            backend,
		transactionService: transactionService,
		address:            address,
	}
}

// Deploy deploys a new vault and returns once the transaction has been submitted.
func (c *factory) Deploy(
	ctx context.Context,
	issuer common.Address,
	nonce common.Hash,
	peerId string,
) (common.Hash, error) {
	callData, err := factoryABI.Pack("deployVault", issuer, nonce, peerId)
	if err != nil {
		return common.Hash{}, err
	}

	request := &transaction.TxRequest{
		To:          &c.address,
		Data:        callData,
		Value:       big.NewInt(0),
		Description: "vault deployment",
	}

	txHash, err := c.transactionService.Send(ctx, request)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

// WaitDeployed waits for the deployment transaction to confirm and returns the vault address
func (c *factory) WaitDeployed(ctx context.Context, txHash common.Hash) (common.Address, error) {
	receipt, err := c.transactionService.WaitForReceipt(ctx, txHash)
	if err != nil {
		return common.Address{}, err
	}

	var event vaultDeployedEvent
	err = transaction.FindSingleEvent(&factoryABI, receipt, c.address, vaultDeployedEventType, &event)
	if err != nil {
		return common.Address{}, fmt.Errorf("contract deployment failed: %w", err)
	}

	return event.ContractAddress, nil
}

// VerifyBytecode checks that the factory is valid.
func (c *factory) VerifyBytecode(ctx context.Context) (err error) {
	code, err := c.backend.CodeAt(ctx, c.address, nil)
	if err != nil {
		return err
	}

	if !bytes.Equal(code, currentDeployVersion) {
		return ErrInvalidFactory
	}

	return nil
}

func (c *factory) verifyVaultAgainstFactory(ctx context.Context, factory, vault common.Address) (bool, error) {
	callData, err := factoryABI.Pack("deployedContracts", vault)
	if err != nil {
		return false, err
	}

	output, err := c.transactionService.Call(ctx, &transaction.TxRequest{
		To:   &factory,
		Data: callData,
	})
	if err != nil {
		return false, err
	}

	results, err := factoryABI.Unpack("deployedContracts", output)
	if err != nil {
		return false, err
	}

	if len(results) != 1 {
		return false, errDecodeABI
	}

	deployed, ok := abi.ConvertType(results[0], new(bool)).(*bool)
	if !ok || deployed == nil {
		return false, errDecodeABI
	}
	if !*deployed {
		return false, nil
	}
	return true, nil
}

// VerifyVault checks that the supplied vault has been deployed by a supported factory.
func (c *factory) VerifyVault(ctx context.Context, vault common.Address) error {
	deployed, err := c.verifyVaultAgainstFactory(ctx, c.address, vault)
	if err != nil {
		return err
	}
	if deployed {
		return nil
	}

	return ErrNotDeployedByFactory
}

// ERC20Address returns the token for which this factory deploys vaults.
func (c *factory) ERC20Address(ctx context.Context) (common.Address, error) {
	callData, err := factoryABI.Pack("TokenAddress")
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

	results, err := factoryABI.Unpack("TokenAddress", output)
	if err != nil {
		return common.Address{}, err
	}

	if len(results) != 1 {
		return common.Address{}, errDecodeABI
	}

	erc20Address, ok := abi.ConvertType(results[0], new(common.Address)).(*common.Address)
	if !ok || erc20Address == nil {
		return common.Address{}, errDecodeABI
	}
	return *erc20Address, nil
}
