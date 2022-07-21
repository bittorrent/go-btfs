package vault_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"

	conabi "github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/bittorrent/go-btfs/transaction/backendmock"
	transactionmock "github.com/bittorrent/go-btfs/transaction/mock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	factoryABI              = transaction.ParseABIUnchecked(conabi.VaultFactoryABI)
	simpleSwapDeployedEvent = factoryABI.Events["SimpleSwapDeployed"]
)

func TestFactoryERC20Address(t *testing.T) {
	factoryAddress := common.HexToAddress("0xabcd")
	erc20Address := common.HexToAddress("0xeffff")
	factory := vault.NewFactory(
		backendmock.New(),
		transactionmock.New(
			transactionmock.WithABICall(
				&factoryABI,
				factoryAddress,
				erc20Address.Hash().Bytes(),
				"TokenAddress",
			),
		),
		factoryAddress,
		//nil,
	)

	addr, err := factory.ERC20Address(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if addr != erc20Address {
		t.Fatalf("wrong erc20Address. wanted %x, got %x", erc20Address, addr)
	}
}

func backendWithCodeAt(codeMap map[common.Address]string) transaction.Backend {
	return backendmock.New(
		backendmock.WithCodeAtFunc(func(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
			code, ok := codeMap[contract]
			if !ok {
				return nil, fmt.Errorf("called with wrong address. wanted one of %v, got %x", codeMap, contract)
			}
			if blockNumber != nil {
				return nil, errors.New("not called for latest block")
			}
			return common.FromHex(code), nil
		}),
	)
}

func TestFactoryVerifySelf(t *testing.T) {
	factoryAddress := common.HexToAddress("0xabcd")

	t.Run("valid", func(t *testing.T) {
		factory := vault.NewFactory(
			backendWithCodeAt(map[common.Address]string{
				factoryAddress: conabi.FactoryDeployedBin,
			}),
			transactionmock.New(),
			factoryAddress,
		)

		err := factory.VerifyBytecode(context.Background())
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("invalid factory contract", func(t *testing.T) {
		factory := vault.NewFactory(
			backendWithCodeAt(map[common.Address]string{
				factoryAddress: "0xabcd",
			}),
			transactionmock.New(),
			factoryAddress,
			//nil,
		)

		err := factory.VerifyBytecode(context.Background())
		if err == nil {
			t.Fatal("verified invalid factory")
		}
		if !errors.Is(err, vault.ErrInvalidFactory) {
			t.Fatalf("wrong error. wanted %v, got %v", vault.ErrInvalidFactory, err)
		}
	})
}

func TestFactoryVerifyVault(t *testing.T) {
	factoryAddress := common.HexToAddress("0xabcd")
	vaultAddress := common.HexToAddress("0xefff")

	t.Run("valid", func(t *testing.T) {
		factory := vault.NewFactory(
			backendmock.New(),
			transactionmock.New(
				transactionmock.WithABICall(
					&factoryABI,
					factoryAddress,
					common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"),
					"deployedContracts",
					vaultAddress,
				),
			),
			factoryAddress,
		)
		err := factory.VerifyVault(context.Background(), vaultAddress)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("valid legacy", func(t *testing.T) {
		factory := vault.NewFactory(
			backendmock.New(),
			transactionmock.New(
				transactionmock.WithABICallSequence(
					transactionmock.ABICall(
						&factoryABI,
						factoryAddress,
						common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"),
						"deployedContracts",
						vaultAddress,
					),
				)),
			factoryAddress,
			//[]common.Address{legacyFactory1, legacyFactory2},
		)

		err := factory.VerifyVault(context.Background(), vaultAddress)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		factory := vault.NewFactory(
			backendmock.New(),
			transactionmock.New(
				transactionmock.WithABICallSequence(
					transactionmock.ABICall(
						&factoryABI,
						factoryAddress,
						common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
						"deployedContracts",
						vaultAddress,
					),
				)),
			factoryAddress,
			//[]common.Address{legacyFactory1, legacyFactory2},
		)

		err := factory.VerifyVault(context.Background(), vaultAddress)
		if err == nil {
			t.Fatal("verified invalid vault")
		}
		if !errors.Is(err, vault.ErrNotDeployedByFactory) {
			t.Fatalf("wrong error. wanted %v, got %v", vault.ErrNotDeployedByFactory, err)
		}
	})
}

//TODO: FIX ME
// func TestFactoryDeploy(t *testing.T) {
// 	factoryAddress := common.HexToAddress("0xabcd")
// 	issuerAddress := common.HexToAddress("0xefff")
// 	peerId := "16Uiu2HAmCvWwf33iyNWQ2t9kML7BJGHcpdFZSSqYbyvFBsDhXN9Q"
// 	deployTransactionHash := common.HexToHash("0xffff")
// 	deployAddress := common.HexToAddress("0xdddd")
// 	nonce := common.HexToHash("eeff")

// 	factory := vault.NewFactory(
// 		backendmock.New(),
// 		transactionmock.New(
// 			transactionmock.WithABICall(
// 				&factoryABI,
// 				factoryAddress,
// 				deployAddress.Hash().Bytes(),
// 				"deployedContracts",
// 				deployAddress,
// 			),
// 			transactionmock.WithABISend(&factoryABI, deployTransactionHash, factoryAddress, big.NewInt(0), "deploySimpleSwap", issuerAddress, nonce),
// 			transactionmock.WithWaitForReceiptFunc(func(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error) {
// 				if txHash != deployTransactionHash {
// 					t.Fatalf("waiting for wrong transaction. wanted %x, got %x", deployTransactionHash, txHash)
// 				}
// 				logData, err := simpleSwapDeployedEvent.Inputs.NonIndexed().Pack(deployAddress)
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				return &types.Receipt{
// 					Status: 1,
// 					Logs: []*types.Log{
// 						{
// 							Data: logData,
// 						},
// 						{
// 							Address: factoryAddress,
// 							Topics:  []common.Hash{simpleSwapDeployedEvent.ID},
// 							Data:    logData,
// 						},
// 					},
// 				}, nil
// 			},
// 			)),
// 		factoryAddress,
// 		//nil,
// 	)

// 	vaultAddress, txHash, err := factory.Deploy(context.Background(), issuerAddress, deployAddress, peerId, deployAddress)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if txHash != deployTransactionHash {
// 		t.Fatalf("returning wrong transaction hash. wanted %x, got %x", deployTransactionHash, txHash)
// 	}

// 	vaultAddress, err = factory.WaitDeployed(context.Background(), txHash)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if vaultAddress != deployAddress {
// 		t.Fatalf("returning wrong address. wanted %x, got %x", deployAddress, vaultAddress)
// 	}
// }

func TestFactoryDeployReverted(t *testing.T) {
	factoryAddress := common.HexToAddress("0xabcd")
	deployTransactionHash := common.HexToHash("0xffff")
	factory := vault.NewFactory(
		backendmock.New(),
		transactionmock.New(
			transactionmock.WithWaitForReceiptFunc(func(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error) {
				if txHash != deployTransactionHash {
					t.Fatalf("waiting for wrong transaction. wanted %x, got %x", deployTransactionHash, txHash)
				}
				return &types.Receipt{
					Status: 0,
				}, nil
			}),
		),
		factoryAddress,
		//nil,
	)

	_, err := factory.WaitDeployed(context.Background(), deployTransactionHash)
	if err == nil {
		t.Fatal("returned failed vault deployment")
	}
	if !errors.Is(err, transaction.ErrTransactionReverted) {
		t.Fatalf("wrong error. wanted %v, got %v", transaction.ErrTransactionReverted, err)
	}
}
