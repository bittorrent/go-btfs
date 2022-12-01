package vault_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	storemock "github.com/bittorrent/go-btfs/statestore/mock"
	transactionmock "github.com/bittorrent/go-btfs/transaction/mock"
	"github.com/ethereum/go-ethereum/common"
)

func TestReceiveCheque(t *testing.T) {
	store := storemock.NewStateStore()
	beneficiary := common.HexToAddress("0xffff")
	issuer := common.HexToAddress("0xbeee")
	cumulativePayout := big.NewInt(101)
	cumulativePayout2 := big.NewInt(201)
	vaultAddress := common.HexToAddress("0xeeee")
	sig := make([]byte, 65)
	chainID := int64(1)
	exchangeRate := big.NewInt(10)

	cheque := &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      beneficiary,
			CumulativePayout: cumulativePayout,
			Vault:            vaultAddress,
		},
		Signature: sig,
	}

	var verifiedWithFactory bool
	factory := &factoryMock{
		verifyVault: func(ctx context.Context, address common.Address) error {
			if address != vaultAddress {
				t.Fatal("verifying wrong vault")
			}
			verifiedWithFactory = true
			return nil
		},
	}

	chequestore := vault.NewChequeStore(
		store,
		factory,
		chainID,
		beneficiary,
		transactionmock.New(
			transactionmock.WithABICallSequence(
				transactionmock.ABICall(&vaultABI, vaultAddress, issuer.Hash().Bytes(), "issuer"),
				transactionmock.ABICall(&vaultABI, vaultAddress, cumulativePayout2.FillBytes(make([]byte, 32)), "totalbalance"),
				transactionmock.ABICall(&vaultABI, vaultAddress, big.NewInt(0).FillBytes(make([]byte, 32)), "paidOut", beneficiary),
				transactionmock.ABICall(&vaultABI, vaultAddress, issuer.Hash().Bytes(), "issuer"),
				transactionmock.ABICall(&vaultABI, vaultAddress, cumulativePayout2.FillBytes(make([]byte, 32)), "totalbalance"),
				transactionmock.ABICall(&vaultABI, vaultAddress, big.NewInt(0).FillBytes(make([]byte, 32)), "paidOut", beneficiary),
			),
		),
		func(c *vault.SignedCheque, cid int64) (common.Address, error) {
			if cid != chainID {
				t.Fatalf("recovery with wrong chain id. wanted %d, got %d", chainID, cid)
			}
			if !cheque.Equal(c) {
				t.Fatalf("recovery with wrong cheque. wanted %v, got %v", cheque, c)
			}
			return issuer, nil
		})

	received, err := chequestore.ReceiveCheque(context.Background(), cheque, exchangeRate, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	if !verifiedWithFactory {
		t.Fatal("did not verify with factory")
	}

	if received.Cmp(cumulativePayout) != 0 {
		t.Fatalf("calculated wrong received cumulativePayout. wanted %d, got %d", cumulativePayout, received)
	}

	lastCheque, err := chequestore.LastReceivedCheque(vaultAddress, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	if !cheque.Equal(lastCheque) {
		t.Fatalf("stored wrong cheque. wanted %v, got %v", cheque, lastCheque)
	}

	cheque = &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      beneficiary,
			CumulativePayout: cumulativePayout2,
			Vault:            vaultAddress,
		},
		Signature: sig,
	}

	verifiedWithFactory = false
	received, err = chequestore.ReceiveCheque(context.Background(), cheque, exchangeRate, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	if verifiedWithFactory {
		t.Fatal("needlessly verify with factory")
	}

	expectedReceived := big.NewInt(0).Sub(cumulativePayout2, cumulativePayout)
	if received.Cmp(expectedReceived) != 0 {
		t.Fatalf("calculated wrong received cumulativePayout. wanted %d, got %d", expectedReceived, received)
	}
}

func TestReceiveChequeInvalidBeneficiary(t *testing.T) {
	store := storemock.NewStateStore()
	beneficiary := common.HexToAddress("0xffff")
	issuer := common.HexToAddress("0xbeee")
	cumulativePayout := big.NewInt(10)
	vaultAddress := common.HexToAddress("0xeeee")
	sig := make([]byte, 65)
	chainID := int64(1)

	cheque := &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      issuer,
			CumulativePayout: cumulativePayout,
			Vault:            vaultAddress,
		},
		Signature: sig,
	}

	chequestore := vault.NewChequeStore(
		store,
		&factoryMock{},
		chainID,
		beneficiary,
		transactionmock.New(),
		nil,
	)

	_, err := chequestore.ReceiveCheque(context.Background(), cheque, cumulativePayout, TOKEN)
	if err == nil {
		t.Fatal("accepted cheque with wrong beneficiary")
	}
	if !errors.Is(err, vault.ErrWrongBeneficiary) {
		t.Fatalf("wrong error. wanted %v, got %v", vault.ErrWrongBeneficiary, err)
	}
}

func TestReceiveChequeInvalidAmount(t *testing.T) {
	store := storemock.NewStateStore()
	beneficiary := common.HexToAddress("0xffff")
	issuer := common.HexToAddress("0xbeee")
	cumulativePayout := big.NewInt(10)
	cumulativePayoutLower := big.NewInt(5)
	vaultAddress := common.HexToAddress("0xeeee")
	sig := make([]byte, 65)
	chainID := int64(1)

	chequestore := vault.NewChequeStore(
		store,
		&factoryMock{
			verifyVault: func(ctx context.Context, address common.Address) error {
				return nil
			},
		},
		chainID,
		beneficiary,
		transactionmock.New(
			transactionmock.WithABICallSequence(
				transactionmock.ABICall(&vaultABI, vaultAddress, issuer.Hash().Bytes(), "issuer"),
				transactionmock.ABICall(&vaultABI, vaultAddress, cumulativePayout.FillBytes(make([]byte, 32)), "totalbalance"),
				transactionmock.ABICall(&vaultABI, vaultAddress, big.NewInt(0).FillBytes(make([]byte, 32)), "paidOut", beneficiary),
			),
		),
		func(c *vault.SignedCheque, cid int64) (common.Address, error) {
			return issuer, nil
		})

	_, err := chequestore.ReceiveCheque(context.Background(), &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      beneficiary,
			CumulativePayout: cumulativePayout,
			Vault:            vaultAddress,
		},
		Signature: sig,
	}, cumulativePayout, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	_, err = chequestore.ReceiveCheque(context.Background(), &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      beneficiary,
			CumulativePayout: cumulativePayoutLower,
			Vault:            vaultAddress,
		},
		Signature: sig,
	}, cumulativePayout, TOKEN)
	if err == nil {
		t.Fatal("accepted lower amount cheque")
	}
	if !errors.Is(err, vault.ErrChequeNotIncreasing) {
		t.Fatalf("wrong error. wanted %v, got %v", vault.ErrChequeNotIncreasing, err)
	}
}

func TestReceiveChequeInvalidVault(t *testing.T) {
	store := storemock.NewStateStore()
	beneficiary := common.HexToAddress("0xffff")
	issuer := common.HexToAddress("0xbeee")
	cumulativePayout := big.NewInt(10)
	vaultAddress := common.HexToAddress("0xeeee")
	sig := make([]byte, 65)
	chainID := int64(1)

	chequestore := vault.NewChequeStore(
		store,
		&factoryMock{
			verifyVault: func(ctx context.Context, address common.Address) error {
				return vault.ErrNotDeployedByFactory
			},
		},
		chainID,
		beneficiary,
		transactionmock.New(
			transactionmock.WithABICallSequence(
				transactionmock.ABICall(&vaultABI, vaultAddress, issuer.Bytes(), "issuer"),
				transactionmock.ABICall(&vaultABI, vaultAddress, cumulativePayout.FillBytes(make([]byte, 32)), "totalbalance"),
			),
		),
		func(c *vault.SignedCheque, cid int64) (common.Address, error) {
			return issuer, nil
		})

	_, err := chequestore.ReceiveCheque(context.Background(), &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      beneficiary,
			CumulativePayout: cumulativePayout,
			Vault:            vaultAddress,
		},
		Signature: sig,
	}, cumulativePayout, TOKEN)
	if !errors.Is(err, vault.ErrNotDeployedByFactory) {
		t.Fatalf("wrong error. wanted %v, got %v", vault.ErrNotDeployedByFactory, err)
	}
}

func TestReceiveChequeInvalidSignature(t *testing.T) {
	store := storemock.NewStateStore()
	beneficiary := common.HexToAddress("0xffff")
	issuer := common.HexToAddress("0xbeee")
	cumulativePayout := big.NewInt(10)
	vaultAddress := common.HexToAddress("0xeeee")
	sig := make([]byte, 65)
	chainID := int64(1)

	chequestore := vault.NewChequeStore(
		store,
		&factoryMock{
			verifyVault: func(ctx context.Context, address common.Address) error {
				return nil
			},
		},
		chainID,
		beneficiary,
		transactionmock.New(
			transactionmock.WithABICallSequence(
				transactionmock.ABICall(&vaultABI, vaultAddress, issuer.Hash().Bytes(), "issuer"),
			),
		),
		func(c *vault.SignedCheque, cid int64) (common.Address, error) {
			return common.Address{}, nil
		})

	_, err := chequestore.ReceiveCheque(context.Background(), &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      beneficiary,
			CumulativePayout: cumulativePayout,
			Vault:            vaultAddress,
		},
		Signature: sig,
	}, cumulativePayout, TOKEN)
	if !errors.Is(err, vault.ErrChequeInvalid) {
		t.Fatalf("wrong error. wanted %v, got %v", vault.ErrChequeInvalid, err)
	}
}

func TestReceiveChequeInsufficientBalance(t *testing.T) {
	store := storemock.NewStateStore()
	beneficiary := common.HexToAddress("0xffff")
	issuer := common.HexToAddress("0xbeee")
	cumulativePayout := big.NewInt(10)
	vaultAddress := common.HexToAddress("0xeeee")
	sig := make([]byte, 65)
	chainID := int64(1)

	chequestore := vault.NewChequeStore(
		store,
		&factoryMock{
			verifyVault: func(ctx context.Context, address common.Address) error {
				return nil
			},
		},
		chainID,
		beneficiary,
		transactionmock.New(
			transactionmock.WithABICallSequence(
				transactionmock.ABICall(&vaultABI, vaultAddress, issuer.Hash().Bytes(), "issuer"),
				transactionmock.ABICall(&vaultABI, vaultAddress, new(big.Int).Sub(cumulativePayout, big.NewInt(1)).FillBytes(make([]byte, 32)), "totalbalance"),
				transactionmock.ABICall(&vaultABI, vaultAddress, big.NewInt(0).FillBytes(make([]byte, 32)), "paidOut", beneficiary),
			),
		),
		func(c *vault.SignedCheque, cid int64) (common.Address, error) {
			return issuer, nil
		})

	_, err := chequestore.ReceiveCheque(context.Background(), &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      beneficiary,
			CumulativePayout: cumulativePayout,
			Vault:            vaultAddress,
		},
		Signature: sig,
	}, cumulativePayout, TOKEN)
	if !errors.Is(err, vault.ErrBouncingCheque) {
		t.Fatalf("wrong error. wanted %v, got %v", vault.ErrBouncingCheque, err)
	}
}

func TestReceiveChequeSufficientBalancePaidOut(t *testing.T) {
	store := storemock.NewStateStore()
	beneficiary := common.HexToAddress("0xffff")
	issuer := common.HexToAddress("0xbeee")
	cumulativePayout := big.NewInt(10)
	vaultAddress := common.HexToAddress("0xeeee")
	sig := make([]byte, 65)
	chainID := int64(1)

	chequestore := vault.NewChequeStore(
		store,
		&factoryMock{
			verifyVault: func(ctx context.Context, address common.Address) error {
				return nil
			},
		},
		chainID,
		beneficiary,
		transactionmock.New(
			transactionmock.WithABICallSequence(
				transactionmock.ABICall(&vaultABI, vaultAddress, issuer.Hash().Bytes(), "issuer"),
				transactionmock.ABICall(&vaultABI, vaultAddress, new(big.Int).Sub(cumulativePayout, big.NewInt(100)).FillBytes(make([]byte, 32)), "totalbalance"),
				transactionmock.ABICall(&vaultABI, vaultAddress, big.NewInt(0).FillBytes(make([]byte, 32)), "paidOut", beneficiary),
			),
		),
		func(c *vault.SignedCheque, cid int64) (common.Address, error) {
			return issuer, nil
		})

	_, err := chequestore.ReceiveCheque(context.Background(), &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      beneficiary,
			CumulativePayout: cumulativePayout,
			Vault:            vaultAddress,
		},
		Signature: sig,
	}, cumulativePayout, TOKEN)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReceiveChequeNotEnoughValue(t *testing.T) {
	store := storemock.NewStateStore()
	beneficiary := common.HexToAddress("0xffff")
	issuer := common.HexToAddress("0xbeee")
	cumulativePayout := big.NewInt(100)
	vaultAddress := common.HexToAddress("0xeeee")
	sig := make([]byte, 65)
	chainID := int64(1)
	amountCheck := big.NewInt(100)

	cheque := &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      beneficiary,
			CumulativePayout: cumulativePayout,
			Vault:            vaultAddress,
		},
		Signature: sig,
	}

	factory := &factoryMock{
		verifyVault: func(ctx context.Context, address common.Address) error {
			if address != vaultAddress {
				t.Fatal("verifying wrong vault")
			}
			return nil
		},
	}

	chequestore := vault.NewChequeStore(
		store,
		factory,
		chainID,
		beneficiary,
		transactionmock.New(
			transactionmock.WithABICallSequence(
				transactionmock.ABICall(&vaultABI, vaultAddress, issuer.Hash().Bytes(), "issuer"),
				transactionmock.ABICall(&vaultABI, vaultAddress, cumulativePayout.FillBytes(make([]byte, 32)), "totalbalance"),
				transactionmock.ABICall(&vaultABI, vaultAddress, big.NewInt(0).FillBytes(make([]byte, 32)), "paidOut", beneficiary),
			),
		),
		func(c *vault.SignedCheque, cid int64) (common.Address, error) {
			if cid != chainID {
				t.Fatalf("recovery with wrong chain id. wanted %d, got %d", chainID, cid)
			}
			if !cheque.Equal(c) {
				t.Fatalf("recovery with wrong cheque. wanted %v, got %v", cheque, c)
			}
			return issuer, nil
		})

	_, err := chequestore.ReceiveCheque(context.Background(), cheque, amountCheck, TOKEN)
	if err != nil {
		t.Fatalf("got wrong error. wanted nil, got %v", err)
	}
}

func TestReceiveChequeNotEnoughValue2(t *testing.T) {
	store := storemock.NewStateStore()
	beneficiary := common.HexToAddress("0xffff")
	issuer := common.HexToAddress("0xbeee")
	cumulativePayout := big.NewInt(100)
	vaultAddress := common.HexToAddress("0xeeee")
	sig := make([]byte, 65)
	chainID := int64(1)

	// cheque needs to cover initial deduction (if applicable) plus one times the exchange rate
	// in order to amount to at least 1 accounting credit and be accepted
	// in this test cheque amount is just not enough to cover that therefore we expect

	amountCheck := big.NewInt(100)

	cheque := &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      beneficiary,
			CumulativePayout: cumulativePayout,
			Vault:            vaultAddress,
		},
		Signature: sig,
	}

	factory := &factoryMock{
		verifyVault: func(ctx context.Context, address common.Address) error {
			if address != vaultAddress {
				t.Fatal("verifying wrong vault")
			}
			return nil
		},
	}

	chequestore := vault.NewChequeStore(
		store,
		factory,
		chainID,
		beneficiary,
		transactionmock.New(
			transactionmock.WithABICallSequence(
				transactionmock.ABICall(&vaultABI, vaultAddress, issuer.Hash().Bytes(), "issuer"),
				transactionmock.ABICall(&vaultABI, vaultAddress, cumulativePayout.FillBytes(make([]byte, 32)), "totalbalance"),
				transactionmock.ABICall(&vaultABI, vaultAddress, big.NewInt(0).FillBytes(make([]byte, 32)), "paidOut", beneficiary),
			),
		),
		func(c *vault.SignedCheque, cid int64) (common.Address, error) {
			if cid != chainID {
				t.Fatalf("recovery with wrong chain id. wanted %d, got %d", chainID, cid)
			}
			if !cheque.Equal(c) {
				t.Fatalf("recovery with wrong cheque. wanted %v, got %v", cheque, c)
			}
			return issuer, nil
		})

	_, err := chequestore.ReceiveCheque(context.Background(), cheque, amountCheck, TOKEN)
	if err != nil {
		t.Fatalf("got wrong error. wanted nil, got %v", err)
	}
}
