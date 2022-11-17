package vault_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"

	erc20mock "github.com/bittorrent/go-btfs/settlement/swap/erc20/mock"
	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	storemock "github.com/bittorrent/go-btfs/statestore/mock"
	"github.com/bittorrent/go-btfs/transaction"
	transactionmock "github.com/bittorrent/go-btfs/transaction/mock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestVaultAddress(t *testing.T) {
	address := common.HexToAddress("0xabcd")
	ownerAdress := common.HexToAddress("0xfff")
	vaultService, err := vault.New(
		transactionmock.New(),
		address,
		ownerAdress,
		nil,
		&chequeSignerMock{},
		erc20mock.New(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	if vaultService.Address() != address {
		t.Fatalf("returned wrong address. wanted %x, got %x", address, vaultService.Address())
	}
}

func TestVaultBalance(t *testing.T) {
	address := common.HexToAddress("0xabcd")
	ownerAdress := common.HexToAddress("0xfff")
	balance := big.NewInt(10)
	vaultService, err := vault.New(
		transactionmock.New(
			transactionmock.WithABICall(&vaultABI, address, balance.FillBytes(make([]byte, 32)), "totalbalance"),
		),
		address,
		ownerAdress,
		nil,
		&chequeSignerMock{},
		erc20mock.New(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	returnedBalance, err := vaultService.TotalBalance(context.Background(), _token)
	if err != nil {
		t.Fatal(err)
	}

	if returnedBalance.Cmp(balance) != 0 {
		t.Fatalf("returned wrong balance. wanted %d, got %d", balance, returnedBalance)
	}
}

func TestVaultDeposit(t *testing.T) {
	address := common.HexToAddress("0xabcd")
	ownerAdress := common.HexToAddress("0xfff")
	balance := big.NewInt(30)
	depositAmount := big.NewInt(20)
	txHash := common.HexToHash("0xdddd")

	factoryAddress := common.HexToAddress("0xabcd")
	// issuerAddress := common.HexToAddress("0xefff")
	// deployTransactionHash := common.HexToHash("0xffff")
	// nonce := common.HexToHash("eeff")

	vaultService, err := vault.New(
		transactionmock.New(
			transactionmock.WithABISend(&vaultABI, txHash, factoryAddress, big.NewInt(0), "deposit", depositAmount),
		),
		address,
		ownerAdress,
		nil,
		&chequeSignerMock{},
		erc20mock.New(
			erc20mock.WithBalanceOfFunc(func(ctx context.Context, address common.Address) (*big.Int, error) {
				if address != ownerAdress {
					return nil, errors.New("getting balance of wrong address")
				}
				return balance, nil
			}),
			erc20mock.WithTransferFunc(func(ctx context.Context, to common.Address, value *big.Int) (common.Hash, error) {
				if to != address {
					return common.Hash{}, fmt.Errorf("sending to wrong address. wanted %x, got %x", address, to)
				}
				if depositAmount.Cmp(value) != 0 {
					return common.Hash{}, fmt.Errorf("sending wrong value. wanted %d, got %d", depositAmount, value)
				}
				return txHash, nil
			}),
		),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	returnedTxHash, err := vaultService.Deposit(context.Background(), depositAmount)
	if err != nil {
		t.Fatal(err)
	}

	if txHash != returnedTxHash {
		t.Fatalf("returned wrong transaction hash. wanted %v, got %v", txHash, returnedTxHash)
	}
}

func TestVaultWaitForDeposit(t *testing.T) {
	address := common.HexToAddress("0xabcd")
	ownerAdress := common.HexToAddress("0xfff")
	txHash := common.HexToHash("0xdddd")
	vaultService, err := vault.New(
		transactionmock.New(
			transactionmock.WithWaitForReceiptFunc(func(ctx context.Context, tx common.Hash) (*types.Receipt, error) {
				if tx != txHash {
					t.Fatalf("waiting for wrong transaction. wanted %x, got %x", txHash, tx)
				}
				return &types.Receipt{
					Status: 1,
				}, nil
			}),
		),
		address,
		ownerAdress,
		nil,
		&chequeSignerMock{},
		erc20mock.New(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = vaultService.WaitForDeposit(context.Background(), txHash)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVaultWaitForDepositReverted(t *testing.T) {
	address := common.HexToAddress("0xabcd")
	ownerAdress := common.HexToAddress("0xfff")
	txHash := common.HexToHash("0xdddd")
	vaultService, err := vault.New(
		transactionmock.New(
			transactionmock.WithWaitForReceiptFunc(func(ctx context.Context, tx common.Hash) (*types.Receipt, error) {
				if tx != txHash {
					t.Fatalf("waiting for wrong transaction. wanted %x, got %x", txHash, tx)
				}
				return &types.Receipt{
					Status: 0,
				}, nil
			}),
		),
		address,
		ownerAdress,
		nil,
		&chequeSignerMock{},
		erc20mock.New(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = vaultService.WaitForDeposit(context.Background(), txHash)
	if err == nil {
		t.Fatal("expected reverted error")
	}
	if !errors.Is(err, transaction.ErrTransactionReverted) {
		t.Fatalf("wrong error. wanted %v, got %v", transaction.ErrTransactionReverted, err)
	}
}

// func TestVaultIssue(t *testing.T) {
// 	address := common.HexToAddress("0xabcd")
// 	beneficiary := common.HexToAddress("0xdddd")
// 	ownerAdress := common.HexToAddress("0xfff")
// 	store := storemock.NewStateStore()
// 	amount := big.NewInt(20)
// 	amount2 := big.NewInt(30)
// 	// expectedCumulative := big.NewInt(50)
// 	sig := common.Hex2Bytes("0xffff")
// 	chequeSigner := &chequeSignerMock{}

// 	vaultService, err := vault.New(
// 		transactionmock.New(
// 			transactionmock.WithABICallSequence(
// 				transactionmock.ABICall(&vaultABINew, address, big.NewInt(100).FillBytes(make([]byte, 32)), "totalbalance"),
// 				transactionmock.ABICall(&vaultABINew, address, big.NewInt(0).FillBytes(make([]byte, 32)), "totalPaidOut"),
// 				transactionmock.ABICall(&vaultABINew, address, big.NewInt(0).FillBytes(make([]byte, 32)), "paidOut", beneficiary),
// 				transactionmock.ABICall(&vaultABINew, address, big.NewInt(100).FillBytes(make([]byte, 32)), "totalbalance"),
// 				transactionmock.ABICall(&vaultABINew, address, big.NewInt(0).FillBytes(make([]byte, 32)), "totalPaidOut"),
// 				transactionmock.ABICall(&vaultABINew, address, big.NewInt(0).FillBytes(make([]byte, 32)), "paidOut", beneficiary),
// 				transactionmock.ABICall(&vaultABINew, address, big.NewInt(100).FillBytes(make([]byte, 32)), "totalbalance"),
// 				transactionmock.ABICall(&vaultABINew, address, big.NewInt(0).FillBytes(make([]byte, 32)), "totalPaidOut"),
// 				transactionmock.ABICall(&vaultABINew, address, big.NewInt(0).FillBytes(make([]byte, 32)), "paidOut", ownerAdress),
// 			),
// 		),
// 		address,
// 		ownerAdress,
// 		store,
// 		chequeSigner,
// 		erc20mock.New(),
// 		mockchequestore.NewChequeStore(
// 			mockchequestore.WithStoreSendChequeRecordFunc(func(vault, beneficiary common.Address, amount *big.Int) error {
// 				return nil
// 			}),
// 		),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// issue a cheque
// 	expectedCheque := &vault.SignedCheque{
// 		Cheque: vault.Cheque{
// 			Beneficiary:      beneficiary,
// 			CumulativePayout: amount,
// 			Vault:            address,
// 		},
// 		Signature: sig,
// 	}

// 	chequeSigner.sign = func(cheque *vault.Cheque) ([]byte, error) {
// 		if !cheque.Equal(&expectedCheque.Cheque) {
// 			t.Fatalf("wrong cheque. wanted %v got %v", expectedCheque.Cheque, cheque)
// 		}
// 		return sig, nil
// 	}

// 	_, err = vaultService.Issue(context.Background(), beneficiary, amount, func(cheque *vault.SignedCheque) error {
// 		if !cheque.Equal(expectedCheque) {
// 			t.Fatalf("wrong cheque. wanted %v got %v", expectedCheque, cheque)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	lastCheque, err := vaultService.LastCheque(beneficiary)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !lastCheque.Equal(expectedCheque) {
// 		t.Fatalf("wrong cheque stored. wanted %v got %v", expectedCheque, lastCheque)
// 	}

// 	// issue another cheque for the same beneficiary
// 	expectedCheque = &vault.SignedCheque{
// 		Cheque: vault.Cheque{
// 			Beneficiary:      beneficiary,
// 			CumulativePayout: amount2,
// 			Vault:            address,
// 		},
// 		Signature: sig,
// 	}

// 	chequeSigner.sign = func(cheque *vault.Cheque) ([]byte, error) {
// 		if !cheque.Equal(&expectedCheque.Cheque) {
// 			t.Fatalf("wrong cheque. wanted %v got %v", expectedCheque, cheque)
// 		}
// 		return sig, nil
// 	}

// 	_, err = vaultService.Issue(context.Background(), beneficiary, amount2, func(cheque *vault.SignedCheque) error {
// 		if !cheque.Equal(expectedCheque) {
// 			t.Fatalf("wrong cheque. wanted %v got %v", expectedCheque, cheque)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	lastCheque, err = vaultService.LastCheque(beneficiary)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !lastCheque.Equal(expectedCheque) {
// 		t.Fatalf("wrong cheque stored. wanted %v got %v", expectedCheque, lastCheque)
// 	}

// 	// issue another cheque for the different beneficiary
// 	expectedChequeOwner := &vault.SignedCheque{
// 		Cheque: vault.Cheque{
// 			Beneficiary:      ownerAdress,
// 			CumulativePayout: amount,
// 			Vault:            address,
// 		},
// 		Signature: sig,
// 	}

// 	chequeSigner.sign = func(cheque *vault.Cheque) ([]byte, error) {
// 		if !cheque.Equal(&expectedChequeOwner.Cheque) {
// 			t.Fatalf("wrong cheque. wanted %v got %v", expectedCheque, cheque)
// 		}
// 		return sig, nil
// 	}

// 	_, err = vaultService.Issue(context.Background(), ownerAdress, amount, func(cheque *vault.SignedCheque) error {
// 		if !cheque.Equal(expectedChequeOwner) {
// 			t.Fatalf("wrong cheque. wanted %v got %v", expectedChequeOwner, cheque)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	lastCheque, err = vaultService.LastCheque(ownerAdress)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !lastCheque.Equal(expectedChequeOwner) {
// 		t.Fatalf("wrong cheque stored. wanted %v got %v", expectedChequeOwner, lastCheque)
// 	}

// 	// finally check this did not interfere with the beneficiary cheque
// 	lastCheque, err = vaultService.LastCheque(beneficiary)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !lastCheque.Equal(expectedCheque) {
// 		t.Fatalf("wrong cheque stored. wanted %v got %v", expectedCheque, lastCheque)
// 	}
// }

func TestVaultIssueErrorSend(t *testing.T) {
	address := common.HexToAddress("0xabcd")
	beneficiary := common.HexToAddress("0xdddd")
	ownerAdress := common.HexToAddress("0xfff")
	store := storemock.NewStateStore()
	amount := big.NewInt(20)
	sig := common.Hex2Bytes("0xffff")
	chequeSigner := &chequeSignerMock{}

	vaultService, err := vault.New(
		transactionmock.New(),
		address,
		ownerAdress,
		store,
		chequeSigner,
		erc20mock.New(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	chequeSigner.sign = func(cheque *vault.Cheque) ([]byte, error) {
		return sig, nil
	}

	_, err = vaultService.Issue(context.Background(), beneficiary, amount, func(cheque *vault.SignedCheque) error {
		return errors.New("err")
	})
	if err == nil {
		t.Fatal("expected error")
	}

	// verify the cheque was not saved
	_, err = vaultService.LastCheque(beneficiary)
	if !errors.Is(err, vault.ErrNoCheque) {
		t.Fatalf("wrong error. wanted %v, got %v", vault.ErrNoCheque, err)
	}
}

func TestVaultIssueOutOfFunds(t *testing.T) {
	address := common.HexToAddress("0xabcd")
	beneficiary := common.HexToAddress("0xdddd")
	ownerAdress := common.HexToAddress("0xfff")
	store := storemock.NewStateStore()
	amount := big.NewInt(20)

	vaultService, err := vault.New(
		transactionmock.New(
			transactionmock.WithABICallSequence(
				transactionmock.ABICall(&vaultABI, address, big.NewInt(0).FillBytes(make([]byte, 32)), "totalbalance"),
				transactionmock.ABICall(&vaultABI, address, big.NewInt(0).FillBytes(make([]byte, 32)), "totalPaidOut"),
			),
		),
		address,
		ownerAdress,
		store,
		&chequeSignerMock{},
		erc20mock.New(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = vaultService.Issue(context.Background(), beneficiary, amount, func(cheque *vault.SignedCheque) error {
		return nil
	})
	if !errors.Is(err, vault.ErrOutOfFunds) {
		t.Fatalf("wrong error. wanted %v, got %v", vault.ErrOutOfFunds, err)
	}

	// verify the cheque was not saved
	_, err = vaultService.LastCheque(beneficiary)

	if !errors.Is(err, vault.ErrNoCheque) {
		t.Fatalf("wrong error. wanted %v, got %v", vault.ErrNoCheque, err)
	}
}

func TestVaultWithdraw(t *testing.T) {
	address := common.HexToAddress("0xabcd")
	ownerAdress := common.HexToAddress("0xfff")
	balance := big.NewInt(30)
	withdrawAmount := big.NewInt(20)
	txHash := common.HexToHash("0xdddd")
	store := storemock.NewStateStore()
	vaultService, err := vault.New(
		transactionmock.New(
			transactionmock.WithABICallSequence(
				transactionmock.ABICall(&vaultABI, address, balance.FillBytes(make([]byte, 32)), "totalbalance"),
				transactionmock.ABICall(&vaultABI, address, big.NewInt(0).FillBytes(make([]byte, 32)), "totalPaidOut"),
			),
			transactionmock.WithABISend(&vaultABI, txHash, address, big.NewInt(0), "withdraw", withdrawAmount),
		),
		address,
		ownerAdress,
		store,
		&chequeSignerMock{},
		erc20mock.New(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	returnedTxHash, err := vaultService.Withdraw(context.Background(), withdrawAmount)
	if err != nil {
		t.Fatal(err)
	}

	if txHash != returnedTxHash {
		t.Fatalf("returned wrong transaction hash. wanted %v, got %v", txHash, returnedTxHash)
	}
}

func TestVaultWithdrawInsufficientFunds(t *testing.T) {
	address := common.HexToAddress("0xabcd")
	ownerAdress := common.HexToAddress("0xfff")
	withdrawAmount := big.NewInt(20)
	txHash := common.HexToHash("0xdddd")
	store := storemock.NewStateStore()
	vaultService, err := vault.New(
		transactionmock.New(
			transactionmock.WithABISend(&vaultABI, txHash, address, big.NewInt(0), "withdraw", withdrawAmount),
			transactionmock.WithABICallSequence(
				transactionmock.ABICall(&vaultABI, address, big.NewInt(0).FillBytes(make([]byte, 32)), "totalbalance"),
				transactionmock.ABICall(&vaultABI, address, big.NewInt(0).FillBytes(make([]byte, 32)), "totalPaidOut"),
			),
		),
		address,
		ownerAdress,
		store,
		&chequeSignerMock{},
		erc20mock.New(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = vaultService.Withdraw(context.Background(), withdrawAmount)
	if !errors.Is(err, vault.ErrInsufficientFunds) {
		t.Fatalf("got wrong error. wanted %v, got %v", vault.ErrInsufficientFunds, err)
	}
}

func TestStateStoreKeys(t *testing.T) {
	address := common.HexToAddress("0xabcd")

	expected := "swap_cashout_000000000000000000000000000000000000abcd"
	if vault.CashoutActionKey(address) != expected {
		t.Fatalf("wrong cashout action key. wanted %s, got %s", expected, vault.CashoutActionKey(address))
	}

	expected = "swap_vault_last_issued_cheque_000000000000000000000000000000000000abcd"
	if vault.LastIssuedChequeKey(address) != expected {
		t.Fatalf("wrong last issued cheque key. wanted %s, got %s", expected, vault.LastIssuedChequeKey(address))
	}

	expected = "swap_vault_last_received_cheque_000000000000000000000000000000000000abcd"
	if vault.LastReceivedChequeKey(address) != expected {
		t.Fatalf("wrong last received cheque key. wanted %s, got %s", expected, vault.LastReceivedChequeKey(address))
	}
}
