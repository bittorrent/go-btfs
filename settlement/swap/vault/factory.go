package vault

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	conabi "github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
)

var (
	ErrInvalidFactory = errors.New(`not a valid factory contract;
	The correct way to upgrade:
		1.Make sure the current version is version 1.6.0 (or other 1.x versions)
		2.Refer to the official documentation to upgrade (https://docs.btfs.io/docs/tutorials-on-upgrading-btfs-v10-to-btfs-v20-mainnet) (Copy this link to open in your browser)`)
	ErrNotDeployedByFactory = errors.New("vault not deployed by factory")
	errDecodeABI            = errors.New("could not decode abi data")

	factoryABI             = transaction.ParseABIUnchecked(conabi.VaultFactoryABI)
	vaultDeployedEventType = factoryABI.Events["VaultDeployed"]
)

// Factory is the main interface for interacting with the vault factory.
type Factory interface {
	// ERC20Address returns the token for which this factory deploys vaults.
	ERC20Address(ctx context.Context) (common.Address, error)
	// Deploy return the existed one if already deployed, otherwise deploy one and return the trxHash
	Deploy(ctx context.Context, issuer common.Address, vaultLogic common.Address, peerId string, tokenAddress common.Address) (vault common.Address, trx common.Hash, err error)
	// WaitDeployed waits for the deployment transaction to confirm and returns the vault address
	WaitDeployed(ctx context.Context, txHash common.Hash) (common.Address, error)
	// VerifyBytecode checks that the factory is valid.
	VerifyBytecode(ctx context.Context) error
	// VerifyVault checks that the supplied vault has been deployed by this factory.
	VerifyVault(ctx context.Context, vault common.Address) error
	// GetPeerVault query peer's vault address deployed by this factory.
	GetPeerVault(ctx context.Context, peerID peer.ID) (vault common.Address, err error)
	// GetPeerVaultWithCache query peer's vault address deployed by this factory. Return cached if cache exists, otherwise query from BTTC.
	GetPeerVaultWithCache(ctx context.Context, peerID peer.ID) (vault common.Address, err error)
	// IsVaultCompatibleBetween checks whether my vault is compatible with the `peerID`'s one.
	IsVaultCompatibleBetween(ctx context.Context, peerID1, peerID2 peer.ID) (isCompatible bool, err error)
}

type factory struct {
	backend            transaction.Backend
	transactionService transaction.Service
	address            common.Address // address of the factory to use for deployments
	peerVaultCache     *cache.Cache
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
	peerVaultCache := cache.New(5*time.Minute, 10*time.Minute)
	return &factory{
		backend:            backend,
		transactionService: transactionService,
		address:            address,
		peerVaultCache:     peerVaultCache,
	}
}

// Deploy return the existed one if already deployed, otherwise deploy one and return the trxHash
func (c *factory) Deploy(
	ctx context.Context,
	issuer common.Address,
	vaultLogic common.Address,
	peerId string,
	erc20Address common.Address,
) (vault common.Address, trx common.Hash, err error) {

	_peerId, err := peer.Decode(peerId)
	if err != nil {
		return
	}

	// check whether vault has already been deployed
	zeroAddr := common.Address{}
	vault, err = c.GetPeerVault(ctx, _peerId)
	if err != nil {
		return
	}
	if vault != zeroAddr {
		var vaultImpl common.Address
		vaultImpl, err = GetVaultImpl(ctx, vault, c.transactionService)
		if err != nil {
			return
		}
		if vaultImpl != vaultLogic {
			fmt.Printf("you already have one vault, you can upgrade it via `btfs vault upgrade` command")
		}
		return
	}

	// deploy one new vault, and return the trxHash
	err = checkBalance(ctx, c.backend, issuer)
	if err != nil {
		return
	}
	trx, err = c.deployVault(ctx, issuer, vaultLogic, peerId, erc20Address)
	return
}

func (c *factory) deployVault(
	ctx context.Context,
	issuer common.Address,
	vaultLogic common.Address,
	peerId string,
	erc20Address common.Address,
) (common.Hash, error) {
	initCallData, err := vaultABI.Pack("init", issuer, erc20Address)
	if err != nil {
		return common.Hash{}, err
	}

	nonceBytes := make([]byte, 32)
	_, err = rand.Read(nonceBytes)
	if err != nil {
		return common.Hash{}, err
	}
	nonce := common.BytesToHash(nonceBytes)

	callData, err := factoryABI.Pack("deployVault", issuer, vaultLogic, nonce, peerId, initCallData)
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

/*
IsVaultCompatibleBetween checks whether my vault is compatible with the `peerID`'s one.
If peer's vaults not compatible, they cannot upload/receive files to/from each other.
*/
func (c *factory) IsVaultCompatibleBetween(ctx context.Context, peerID1, peerID2 peer.ID) (isCompatible bool, err error) {
	notFound := common.Address{}

	// Validate whether peer1 and peer2 are using the same factory
	vault1, err := c.GetPeerVaultWithCache(ctx, peerID1)
	if err != nil || vault1 == notFound {
		return
	}
	vault2, err := c.GetPeerVaultWithCache(ctx, peerID2)
	if err != nil || vault2 == notFound {
		return
	}

	// Validate whether peer1 and peer2 are using the same vault implementation;
	// Because vaults are deployed in proxy mode, different peer may use different vault implementation.
	isCompatible, err = IsVaultImplCompatibleBetween(ctx, vault1, vault2, c.transactionService)
	return
}

/*
GetPeerVaultWithCache query peer's vault address deployed by this factory.
Return cached if cache exists, otherwise query from BTTC.
*/
func (c *factory) GetPeerVaultWithCache(ctx context.Context, peerID peer.ID) (vault common.Address, err error) {
	peerStr := peerID.String()
	vaultAddrIf, found := c.peerVaultCache.Get(peerStr)
	if found {
		return vaultAddrIf.(common.Address), nil
	}

	vault, err = c.GetPeerVault(ctx, peerID)
	if err != nil {
		return vault, err
	}
	c.peerVaultCache.SetDefault(peerStr, vault)
	return vault, nil
}

// GetPeerVault query peer's vault address deployed by this factory .
func (c *factory) GetPeerVault(ctx context.Context, peerID peer.ID) (vault common.Address, err error) {
	callData, err := factoryABI.Pack("peerVaultAddress", peerID.String())
	if err != nil {
		return vault, err
	}
	output, err := c.transactionService.Call(ctx, &transaction.TxRequest{
		To:   &c.address,
		Data: callData,
	})
	if err != nil {
		return vault, err
	}

	results, err := factoryABI.Unpack("peerVaultAddress", output)
	if err != nil {
		return vault, err
	}
	if len(results) != 1 {
		return vault, errDecodeABI
	}
	vaultAddr, ok := abi.ConvertType(results[0], new(common.Address)).(*common.Address)
	if !ok || vaultAddr == nil {
		return vault, errDecodeABI
	}
	return *vaultAddr, nil
}
