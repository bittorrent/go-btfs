package vault_test

import (
	"context"
	"math/big"
	"testing"

	conabi "github.com/bittorrent/go-btfs/chain/abi"
	chequestoremock "github.com/bittorrent/go-btfs/settlement/swap/chequestore/mock"
	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	storemock "github.com/bittorrent/go-btfs/statestore/mock"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/bittorrent/go-btfs/transaction/backendmock"
	transactionmock "github.com/bittorrent/go-btfs/transaction/mock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	vaultABI               = transaction.ParseABIUnchecked(conabi.VaultABI)
	chequeCashedEventType  = vaultABI.Events["ChequeCashed"]
	chequeBouncedEventType = vaultABI.Events["ChequeBounced"]
)

func TestCashout(t *testing.T) {
	vault.RestartFixCashOutStatusLock = false
	vaultAddress := common.HexToAddress("abcd")
	recipientAddress := common.HexToAddress("efff")
	txHash := common.HexToHash("dddd")
	totalPayout := big.NewInt(100)
	cumulativePayout := big.NewInt(500)

	cheque := &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      common.HexToAddress("aaaa"),
			CumulativePayout: cumulativePayout,
			Vault:            vaultAddress,
		},
		Signature: []byte{},
	}

	store := storemock.NewStateStore()
	cashoutService := vault.NewCashoutService(
		store,
		backendmock.New(
			backendmock.WithTransactionByHashFunc(func(_ context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
				if hash != txHash {
					t.Fatalf("fetching wrong transaction. wanted %v, got %v", txHash, hash)
				}
				return nil, false, nil
			}),
			backendmock.WithTransactionReceiptFunc(func(_ context.Context, hash common.Hash) (*types.Receipt, error) {
				if hash != txHash {
					t.Fatalf("fetching receipt for transaction. wanted %v, got %v", txHash, hash)
				}

				logData, err := chequeCashedEventType.Inputs.NonIndexed().Pack(totalPayout, cumulativePayout, big.NewInt(0))
				if err != nil {
					t.Fatal(err)
				}

				return &types.Receipt{
					Status: types.ReceiptStatusSuccessful,
					Logs: []*types.Log{
						{
							Address: vaultAddress,
							Topics:  []common.Hash{chequeCashedEventType.ID, cheque.Beneficiary.Hash(), recipientAddress.Hash(), cheque.Beneficiary.Hash()},
							Data:    logData,
						},
					},
				}, nil
			}),
		),
		transactionmock.New(
			transactionmock.WithABISend(&vaultABI, txHash, vaultAddress, big.NewInt(0), "cashChequeBeneficiary", recipientAddress, cheque.CumulativePayout, cheque.Signature),
		),
		chequestoremock.NewChequeStore(
			chequestoremock.WithLastChequeFunc(func(c common.Address) (*vault.SignedCheque, error) {
				if c != vaultAddress {
					t.Fatalf("using wrong vault. wanted %v, got %v", vaultAddress, c)
				}
				return cheque, nil
			}),
			chequestoremock.WithLastReceivedChequeFunc(func(vault common.Address) (*vault.SignedCheque, error) {
				if vault != vaultAddress {
					t.Fatalf("using wrong vault. wanted %v, got %v", vaultAddress, vault)
				}
				return cheque, nil
			}),
		),
	)

	returnedTxHash, err := cashoutService.CashCheque(context.Background(), vaultAddress, recipientAddress, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	if returnedTxHash != txHash {
		t.Fatalf("returned wrong transaction hash. wanted %v, got %v", txHash, returnedTxHash)
	}

	status, err := cashoutService.CashoutStatus(context.Background(), vaultAddress, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	verifyStatus(t, status, vault.CashoutStatus{
		Last: &vault.LastCashout{
			TxHash: txHash,
			Cheque: *cheque,
			Result: &vault.CashChequeResult{
				Beneficiary:      cheque.Beneficiary,
				Recipient:        recipientAddress,
				Caller:           cheque.Beneficiary,
				TotalPayout:      totalPayout,
				CumulativePayout: cumulativePayout,
				CallerPayout:     big.NewInt(0),
				Bounced:          false,
			},
			Reverted: false,
		},
		UncashedAmount: big.NewInt(0),
	})
}

func TestCashoutBounced(t *testing.T) {
	vault.RestartFixCashOutStatusLock = false
	vaultAddress := common.HexToAddress("abcd")
	recipientAddress := common.HexToAddress("efff")
	txHash := common.HexToHash("dddd")
	totalPayout := big.NewInt(100)
	cumulativePayout := big.NewInt(500)

	cheque := &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      common.HexToAddress("aaaa"),
			CumulativePayout: cumulativePayout,
			Vault:            vaultAddress,
		},
		Signature: []byte{},
	}

	store := storemock.NewStateStore()
	cashoutService := vault.NewCashoutService(
		store,
		backendmock.New(
			backendmock.WithTransactionByHashFunc(func(_ context.Context, hash common.Hash) (*types.Transaction, bool, error) {
				if hash != txHash {
					t.Fatalf("fetching wrong transaction. wanted %v, got %v", txHash, hash)
				}
				return nil, false, nil
			}),
			backendmock.WithTransactionReceiptFunc(func(_ context.Context, hash common.Hash) (*types.Receipt, error) {
				if hash != txHash {
					t.Fatalf("fetching receipt for transaction. wanted %v, got %v", txHash, hash)
				}

				chequeCashedLogData, err := chequeCashedEventType.Inputs.NonIndexed().Pack(totalPayout, cumulativePayout, big.NewInt(0))
				if err != nil {
					t.Fatal(err)
				}

				return &types.Receipt{
					Status: types.ReceiptStatusSuccessful,
					Logs: []*types.Log{
						{
							Address: vaultAddress,
							Topics:  []common.Hash{chequeCashedEventType.ID, cheque.Beneficiary.Hash(), recipientAddress.Hash(), cheque.Beneficiary.Hash()},
							Data:    chequeCashedLogData,
						},
						{
							Address: vaultAddress,
							Topics:  []common.Hash{chequeBouncedEventType.ID},
						},
					},
				}, nil
			}),
		),
		transactionmock.New(
			transactionmock.WithABISend(&vaultABI, txHash, vaultAddress, big.NewInt(0), "cashChequeBeneficiary", recipientAddress, cheque.CumulativePayout, cheque.Signature),
		),
		chequestoremock.NewChequeStore(
			chequestoremock.WithLastChequeFunc(func(c common.Address) (*vault.SignedCheque, error) {
				if c != vaultAddress {
					t.Fatalf("using wrong vault. wanted %v, got %v", vaultAddress, c)
				}
				return cheque, nil
			}),
			chequestoremock.WithLastReceivedChequeFunc(func(vault common.Address) (*vault.SignedCheque, error) {
				if vault != vaultAddress {
					t.Fatalf("using wrong vault. wanted %v, got %v", vaultAddress, vault)
				}
				return cheque, nil
			}),
		),
	)

	returnedTxHash, err := cashoutService.CashCheque(context.Background(), vaultAddress, recipientAddress, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	if returnedTxHash != txHash {
		t.Fatalf("returned wrong transaction hash. wanted %v, got %v", txHash, returnedTxHash)
	}

	status, err := cashoutService.CashoutStatus(context.Background(), vaultAddress, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	verifyStatus(t, status, vault.CashoutStatus{
		Last: &vault.LastCashout{
			TxHash: txHash,
			Cheque: *cheque,
			Result: &vault.CashChequeResult{
				Beneficiary:      cheque.Beneficiary,
				Recipient:        recipientAddress,
				Caller:           cheque.Beneficiary,
				TotalPayout:      totalPayout,
				CumulativePayout: cumulativePayout,
				CallerPayout:     big.NewInt(0),
				Bounced:          true,
			},
			Reverted: false,
		},
		UncashedAmount: big.NewInt(0),
	})
}

func TestCashoutStatusReverted(t *testing.T) {
	vault.RestartFixCashOutStatusLock = false
	vaultAddress := common.HexToAddress("abcd")
	recipientAddress := common.HexToAddress("efff")
	txHash := common.HexToHash("dddd")
	cumulativePayout := big.NewInt(500)
	onChainPaidOut := big.NewInt(100)
	beneficiary := common.HexToAddress("aaaa")

	cheque := &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      beneficiary,
			CumulativePayout: cumulativePayout,
			Vault:            vaultAddress,
		},
		Signature: []byte{},
	}

	store := storemock.NewStateStore()
	cashoutService := vault.NewCashoutService(
		store,
		backendmock.New(
			backendmock.WithTransactionByHashFunc(func(_ context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
				if hash != txHash {
					t.Fatalf("fetching wrong transaction. wanted %v, got %v", txHash, hash)
				}
				return nil, false, nil
			}),
			backendmock.WithTransactionReceiptFunc(func(_ context.Context, hash common.Hash) (*types.Receipt, error) {
				if hash != txHash {
					t.Fatalf("fetching receipt for transaction. wanted %v, got %v", txHash, hash)
				}
				return &types.Receipt{
					Status: types.ReceiptStatusFailed,
				}, nil
			}),
		),
		transactionmock.New(
			transactionmock.WithABISend(&vaultABI, txHash, vaultAddress, big.NewInt(0), "cashChequeBeneficiary", recipientAddress, cheque.CumulativePayout, cheque.Signature),
			transactionmock.WithABICall(&vaultABI, vaultAddress, onChainPaidOut.FillBytes(make([]byte, 32)), "paidOut", beneficiary),
		),
		chequestoremock.NewChequeStore(
			chequestoremock.WithLastChequeFunc(func(c common.Address) (*vault.SignedCheque, error) {
				if c != vaultAddress {
					t.Fatalf("using wrong vault. wanted %v, got %v", vaultAddress, c)
				}
				return cheque, nil
			}),
			chequestoremock.WithLastReceivedChequeFunc(func(vault common.Address) (*vault.SignedCheque, error) {
				if vault != vaultAddress {
					t.Fatalf("using wrong vault. wanted %v, got %v", vaultAddress, vault)
				}
				return cheque, nil
			}),
		),
	)

	returnedTxHash, err := cashoutService.CashCheque(context.Background(), vaultAddress, recipientAddress, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	if returnedTxHash != txHash {
		t.Fatalf("returned wrong transaction hash. wanted %v, got %v", txHash, returnedTxHash)
	}

	status, err := cashoutService.CashoutStatus(context.Background(), vaultAddress, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	verifyStatus(t, status, vault.CashoutStatus{
		Last: &vault.LastCashout{
			Reverted: true,
			TxHash:   txHash,
			Cheque:   *cheque,
		},
		UncashedAmount: new(big.Int).Sub(cheque.CumulativePayout, onChainPaidOut),
	})
}

func TestCashoutStatusPending(t *testing.T) {
	vault.RestartFixCashOutStatusLock = false
	vaultAddress := common.HexToAddress("abcd")
	recipientAddress := common.HexToAddress("efff")
	txHash := common.HexToHash("dddd")
	cumulativePayout := big.NewInt(500)

	cheque := &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      common.HexToAddress("aaaa"),
			CumulativePayout: cumulativePayout,
			Vault:            vaultAddress,
		},
		Signature: []byte{},
	}

	store := storemock.NewStateStore()
	cashoutService := vault.NewCashoutService(
		store,
		backendmock.New(
			backendmock.WithTransactionByHashFunc(func(_ context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
				if hash != txHash {
					t.Fatalf("fetching wrong transaction. wanted %v, got %v", txHash, hash)
				}
				return nil, true, nil
			}),
		),
		transactionmock.New(
			transactionmock.WithABISend(&vaultABI, txHash, vaultAddress, big.NewInt(0), "cashChequeBeneficiary", recipientAddress, cheque.CumulativePayout, cheque.Signature),
		),
		chequestoremock.NewChequeStore(
			chequestoremock.WithLastChequeFunc(func(c common.Address) (*vault.SignedCheque, error) {
				if c != vaultAddress {
					t.Fatalf("using wrong vault. wanted %v, got %v", vaultAddress, c)
				}
				return cheque, nil
			}),
			chequestoremock.WithLastReceivedChequeFunc(func(vault common.Address) (*vault.SignedCheque, error) {
				if vault != vaultAddress {
					t.Fatalf("using wrong vault. wanted %v, got %v", vaultAddress, vault)
				}
				return cheque, nil
			}),
		),
	)

	returnedTxHash, err := cashoutService.CashCheque(context.Background(), vaultAddress, recipientAddress, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	if returnedTxHash != txHash {
		t.Fatalf("returned wrong transaction hash. wanted %v, got %v", txHash, returnedTxHash)
	}

	status, err := cashoutService.CashoutStatus(context.Background(), vaultAddress, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	verifyStatus(t, status, vault.CashoutStatus{
		Last: &vault.LastCashout{
			Reverted: false,
			TxHash:   txHash,
			Cheque:   *cheque,
			Result:   nil,
		},
		UncashedAmount: big.NewInt(0),
	})

}

func verifyStatus(t *testing.T, status *vault.CashoutStatus, expected vault.CashoutStatus) {
	if expected.Last == nil {
		if status.Last != nil {
			t.Fatal("unexpected last cashout")
		}
	} else {
		if status.Last == nil {
			t.Fatal("no last cashout")
		}
		if status.Last.Reverted != expected.Last.Reverted {
			t.Fatalf("wrong reverted value. wanted %v, got %v", expected.Last.Reverted, status.Last.Reverted)
		}
		if status.Last.TxHash != expected.Last.TxHash {
			t.Fatalf("wrong transaction hash. wanted %v, got %v", expected.Last.TxHash, status.Last.TxHash)
		}
		if !status.Last.Cheque.Equal(&expected.Last.Cheque) {
			t.Fatalf("wrong cheque in status. wanted %v, got %v", expected.Last.Cheque, status.Last.Cheque)
		}

		if expected.Last.Result != nil {
			if !expected.Last.Result.Equal(status.Last.Result) {
				t.Fatalf("wrong result. wanted %v, got %v", expected.Last.Result, status.Last.Result)
			}
		}
	}

	if status.UncashedAmount.Cmp(expected.UncashedAmount) != 0 {
		t.Fatalf("wrong uncashed amount. wanted %d, got %d", expected.UncashedAmount, status.UncashedAmount)
	}
}
