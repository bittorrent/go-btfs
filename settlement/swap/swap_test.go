package swap_test

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/bittorrent/go-btfs/settlement/swap"
	mockchequestore "github.com/bittorrent/go-btfs/settlement/swap/chequestore/mock"
	"github.com/bittorrent/go-btfs/settlement/swap/swapprotocol"
	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	mockvault "github.com/bittorrent/go-btfs/settlement/swap/vault/mock"
	mockstore "github.com/bittorrent/go-btfs/statestore/mock"
	"github.com/ethereum/go-ethereum/common"

	peerInfo "github.com/libp2p/go-libp2p/core/peer"
)

var TOKEN = common.HexToAddress("0x000")

type swapProtocolMock struct {
	emitCheque func(context.Context, string, *big.Int, string, common.Address, swapprotocol.IssueFunc) (*big.Int, error)
}

func (m *swapProtocolMock) EmitCheque(ctx context.Context, peer string, value *big.Int, contractId string, token common.Address, issueFunc swapprotocol.IssueFunc) (*big.Int, error) {
	if m.emitCheque != nil {
		return m.emitCheque(ctx, peer, value, contractId, token, issueFunc)
	}
	return nil, errors.New("not implemented")
}

type testObserver struct {
	receivedCalled chan notifyPaymentReceivedCall
	sentCalled     chan notifyPaymentSentCall
}

type notifyPaymentReceivedCall struct {
	peer   string
	amount *big.Int
}

type notifyPaymentSentCall struct {
	peer   string
	amount *big.Int
	err    error
}

func newTestObserver() *testObserver {
	return &testObserver{
		receivedCalled: make(chan notifyPaymentReceivedCall, 1),
		sentCalled:     make(chan notifyPaymentSentCall, 1),
	}
}

func (t *testObserver) PeerDebt(peer string) (*big.Int, error) {
	return nil, nil
}

func (t *testObserver) NotifyPaymentReceived(peer string, amount *big.Int, token common.Address) error {
	t.receivedCalled <- notifyPaymentReceivedCall{
		peer:   peer,
		amount: amount,
	}
	return nil
}

func (t *testObserver) NotifyRefreshmentReceived(peer string, amount *big.Int) error {
	return nil
}

func (t *testObserver) Settle(peer string, amount *big.Int, contractId string, token common.Address) error {
	return nil
}

func (t *testObserver) NotifyPaymentSent(peer string, amount *big.Int, err error, token common.Address) {
	t.sentCalled <- notifyPaymentSentCall{
		peer:   peer,
		amount: amount,
		err:    err,
	}
}

func (t *testObserver) Connect(peer string) {

}

func (t *testObserver) Disconnect(peer string) {

}

type addressbookMock struct {
	migratePeer     func(oldPeer, newPeer string) error
	beneficiary     func(peer string) (beneficiary common.Address, known bool, err error)
	vault           func(peer string) (vaultAddress common.Address, known bool, err error)
	beneficiaryPeer func(beneficiary common.Address) (peer string, known bool, err error)
	vaultPeer       func(vault common.Address) (peer string, known bool, err error)
	putBeneficiary  func(peer string, beneficiary common.Address) error
	putVault        func(peer string, vault common.Address) error
}

func (m *addressbookMock) MigratePeer(oldPeer, newPeer string) error {
	return m.migratePeer(oldPeer, newPeer)
}
func (m *addressbookMock) Beneficiary(peer string) (beneficiary common.Address, known bool, err error) {
	return m.beneficiary(peer)
}
func (m *addressbookMock) Vault(peer string) (vaultAddress common.Address, known bool, err error) {
	return m.vault(peer)
}
func (m *addressbookMock) BeneficiaryPeer(beneficiary common.Address) (peer string, known bool, err error) {
	return m.beneficiaryPeer(beneficiary)
}
func (m *addressbookMock) VaultPeer(vault common.Address) (peer string, known bool, err error) {
	return m.vaultPeer(vault)
}
func (m *addressbookMock) PutBeneficiary(peer string, beneficiary common.Address) error {
	return m.putBeneficiary(peer, beneficiary)
}
func (m *addressbookMock) PutVault(peer string, vault common.Address) error {
	return m.putVault(peer, vault)
}

type cashoutMock struct {
	cashCheque              func(ctx context.Context, vault, recipient common.Address, token common.Address) (common.Hash, error)
	cashoutStatus           func(ctx context.Context, vaultAddress common.Address, token common.Address) (*vault.CashoutStatus, error)
	cashoutResults          func() ([]vault.CashOutResult, error)
	hasCashoutAction        func(ctx context.Context, peer common.Address, token common.Address) (bool, error)
	adjustCashCheque        func(ctx context.Context, vaultAddress, recipient common.Address, token common.Address, restartPassFlag bool) (totalCashOutAmount, newCashOutAmount *big.Int, err error)
	restartFixChequeCashOut func()
}

func (m *cashoutMock) CashCheque(ctx context.Context, vault, recipient common.Address, token common.Address) (common.Hash, error) {
	return m.cashCheque(ctx, vault, recipient, token)
}
func (m *cashoutMock) CashoutStatus(ctx context.Context, vaultAddress common.Address, token common.Address) (*vault.CashoutStatus, error) {
	return m.cashoutStatus(ctx, vaultAddress, token)
}
func (m *cashoutMock) CashoutResults() ([]vault.CashOutResult, error) {
	return m.cashoutResults()
}
func (m *cashoutMock) HasCashoutAction(ctx context.Context, peer common.Address, token common.Address) (bool, error) {
	return m.hasCashoutAction(ctx, peer, token)
}
func (m *cashoutMock) AdjustCashCheque(ctx context.Context, vaultAddress, recipient common.Address, token common.Address, restartPassFlag bool) (totalCashOutAmount, newCashOutAmount *big.Int, err error) {
	return m.adjustCashCheque(ctx, vaultAddress, recipient, token, restartPassFlag)
}
func (m *cashoutMock) RestartFixChequeCashOut() {
	m.restartFixChequeCashOut()
}
func TestReceiveCheque(t *testing.T) {
	store := mockstore.NewStateStore()
	vaultService := mockvault.NewVault(
		mockvault.WithTotalReceivedFunc(func(token common.Address) (*big.Int, error) {
			return big.NewInt(0), nil
		}),
		mockvault.WithTotalReceivedCountFunc(func(token common.Address) (int, error) {
			return 0, nil
		}),
	)
	amount := big.NewInt(4)
	exchangeRate := big.NewInt(10)
	vaultAddress := common.HexToAddress("0xcd")

	peer := peerInfo.ID("abcd").String()
	cheque := &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      common.HexToAddress("0xab"),
			CumulativePayout: big.NewInt(10),
			Vault:            vaultAddress,
		},
		Signature: []byte{},
	}

	chequeStore := mockchequestore.NewChequeStore(
		mockchequestore.WithReceiveChequeFunc(func(ctx context.Context, c *vault.SignedCheque, e *big.Int) (*big.Int, error) {
			if !cheque.Equal(c) {
				t.Fatalf("passed wrong cheque to store. wanted %v, got %v", cheque, c)
			}
			if exchangeRate.Cmp(e) != 0 {
				t.Fatalf("passed wrong exchange rate to store. wanted %v, got %v", exchangeRate, e)
			}
			return amount, nil
		}),
	)

	networkID := int64(1)
	addressbook := &addressbookMock{
		vault: func(p string) (common.Address, bool, error) {
			//if !peer.Equal(p) {
			if strings.Compare(peer, p) != 0 {
				t.Fatal("querying vault for wrong peer")
			}
			return vaultAddress, true, nil
		},
		putVault: func(p string, vault common.Address) (err error) {
			//if !peer.Equal(p) {
			if strings.Compare(peer, p) != 0 {
				t.Fatal("storing vault for wrong peer")
			}
			if vault != vaultAddress {
				t.Fatal("storing wrong vault")
			}
			return nil
		},
	}

	observer := newTestObserver()

	swap := swap.New(
		&swapProtocolMock{},
		store,
		vaultService,
		chequeStore,
		addressbook,
		networkID,
		&cashoutMock{},
		observer,
	)

	err := swap.ReceiveCheque(context.Background(), peer, cheque, exchangeRate, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	expectedAmount := big.NewInt(4)

	select {
	case call := <-observer.receivedCalled:
		if call.amount.Cmp(expectedAmount) != 0 {
			t.Fatalf("observer called with wrong amount. got %d, want %d", call.amount, expectedAmount)
		}

		//if !call.peer.Equal(peer) {
		if strings.Compare(call.peer, peer) != 0 {
			t.Fatalf("observer called with wrong peer. got %v, want %v", call.peer, peer)
		}

	case <-time.After(time.Second):
		t.Fatal("expected observer to be called")
	}

}

func TestReceiveChequeReject(t *testing.T) {
	store := mockstore.NewStateStore()
	vaultService := mockvault.NewVault()
	vaultAddress := common.HexToAddress("0xcd")
	exchangeRate := big.NewInt(10)

	peer := peerInfo.ID("abcd").String()
	cheque := &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      common.HexToAddress("0xab"),
			CumulativePayout: big.NewInt(10),
			Vault:            vaultAddress,
		},
		Signature: []byte{},
	}

	var errReject = errors.New("reject")

	chequeStore := mockchequestore.NewChequeStore(
		mockchequestore.WithReceiveChequeFunc(func(ctx context.Context, c *vault.SignedCheque, e *big.Int) (*big.Int, error) {
			return nil, errReject
		}),
	)
	networkID := int64(1)
	addressbook := &addressbookMock{
		vault: func(p string) (common.Address, bool, error) {
			return vaultAddress, true, nil
		},
	}

	observer := newTestObserver()

	swap := swap.New(
		&swapProtocolMock{},
		store,
		vaultService,
		chequeStore,
		addressbook,
		networkID,
		&cashoutMock{},
		observer,
	)

	err := swap.ReceiveCheque(context.Background(), peer, cheque, exchangeRate, TOKEN)
	if err == nil {
		t.Fatal("accepted invalid cheque")
	}
	if !errors.Is(err, errReject) {
		t.Fatalf("wrong error. wanted %v, got %v", errReject, err)
	}

	select {
	case <-observer.receivedCalled:
		t.Fatalf("observer called by error.")
	default:
	}

}

func TestReceiveChequeWrongVault(t *testing.T) {
	store := mockstore.NewStateStore()
	vaultService := mockvault.NewVault()
	vaultAddress := common.HexToAddress("0xcd")
	exchangeRate := big.NewInt(10)

	peer := peerInfo.ID("abcd").String()
	cheque := &vault.SignedCheque{
		Cheque: vault.Cheque{
			Beneficiary:      common.HexToAddress("0xab"),
			CumulativePayout: big.NewInt(10),
			Vault:            vaultAddress,
		},
		Signature: []byte{},
	}

	chequeStore := mockchequestore.NewChequeStore()
	networkID := int64(1)
	addressbook := &addressbookMock{
		vault: func(p string) (common.Address, bool, error) {
			return common.HexToAddress("0xcfff"), true, nil
		},
	}

	observer := newTestObserver()
	swapService := swap.New(
		&swapProtocolMock{},
		store,
		vaultService,
		chequeStore,
		addressbook,
		networkID,
		&cashoutMock{},
		observer,
	)

	err := swapService.ReceiveCheque(context.Background(), peer, cheque, exchangeRate, TOKEN)
	if err == nil {
		t.Fatal("accepted invalid cheque")
	}
	if !errors.Is(err, swap.ErrWrongVault) {
		t.Fatalf("wrong error. wanted %v, got %v", swap.ErrWrongVault, err)
	}

	select {
	case <-observer.receivedCalled:
		t.Fatalf("observer called by error.")
	default:
	}

}

// TODO: FIX ME(checks if pay case is right and workable)
func TestPay(t *testing.T) {
	store := mockstore.NewStateStore()

	amount := big.NewInt(50)
	beneficiary := common.HexToAddress("0xcd")
	peer := peerInfo.ID("abcd").String()

	networkID := int64(1)
	addressbook := &addressbookMock{
		beneficiary: func(p string) (common.Address, bool, error) {
			//if !peer.Equal(p) {
			if strings.Compare(peer, p) != 0 {
				t.Fatal("querying beneficiary for wrong peer")
			}
			return beneficiary, true, nil
		},
	}

	observer := newTestObserver()

	var emitCalled bool
	swap := swap.New(
		&swapProtocolMock{
			emitCheque: func(ctx context.Context, p string, a *big.Int, s string, token common.Address, issueFunc swapprotocol.IssueFunc) (*big.Int, error) {
				//if !peer.Equal(p) {
				if strings.Compare(peer, p) != 0 {
					t.Fatal("sending to wrong peer")
				}
				/*
					if b != beneficiary {
						t.Fatal("issuing for wrong beneficiary")
					}
				*/
				if amount.Cmp(a) != 0 {
					t.Fatal("issuing with wrong amount")
				}
				emitCalled = true
				return amount, nil
			},
		},
		store,
		mockvault.NewVault(),
		mockchequestore.NewChequeStore(),
		addressbook,
		networkID,
		&cashoutMock{},
		observer,
	)

	swap.Pay(context.Background(), peer, amount, "", TOKEN)

	if !emitCalled {
		t.Fatal("swap protocol was not called")
	}
}

func TestPayIssueError(t *testing.T) {
	store := mockstore.NewStateStore()

	amount := big.NewInt(50)
	beneficiary := common.HexToAddress("0xcd")

	peer := peerInfo.ID("abcd").String()
	errReject := errors.New("reject")
	networkID := int64(1)
	addressbook := &addressbookMock{
		beneficiary: func(p string) (common.Address, bool, error) {
			//if !peer.Equal(p) {
			if strings.Compare(peer, p) != 0 {
				t.Fatal("querying beneficiary for wrong peer")
			}
			return beneficiary, true, nil
		},
	}

	swap := swap.New(
		&swapProtocolMock{
			emitCheque: func(c context.Context, a1 string, i *big.Int, s string, token common.Address, issueFunc swapprotocol.IssueFunc) (*big.Int, error) {
				return nil, errReject
			},
		},
		store,
		mockvault.NewVault(),
		mockchequestore.NewChequeStore(),
		addressbook,
		networkID,
		&cashoutMock{},
		nil,
	)

	observer := newTestObserver()
	swap.SetAccounting(observer)

	swap.Pay(context.Background(), peer, amount, "", TOKEN)
	select {
	case call := <-observer.sentCalled:

		//if !call.peer.Equal(peer) {
		if strings.Compare(call.peer, peer) != 0 {
			t.Fatalf("observer called with wrong peer. got %v, want %v", call.peer, peer)
		}
		if !errors.Is(call.err, errReject) {
			t.Fatalf("wrong error. wanted %v, got %v", errReject, call.err)
		}

	case <-time.After(time.Second):
		t.Fatal("expected observer to be called")
	}

}

func TestPayUnknownBeneficiary(t *testing.T) {
	store := mockstore.NewStateStore()

	amount := big.NewInt(50)
	peer := peerInfo.ID("abcd").String()
	networkID := int64(1)
	addressbook := &addressbookMock{
		beneficiary: func(p string) (common.Address, bool, error) {
			//if !peer.Equal(p) {
			if strings.Compare(peer, p) != 0 {
				t.Fatal("querying beneficiary for wrong peer")
			}
			return common.Address{}, false, nil
		},
	}

	observer := newTestObserver()

	swapService := swap.New(
		&swapProtocolMock{
			emitCheque: func(ctx context.Context, p string, a *big.Int, s string, token common.Address, issueFunc swapprotocol.IssueFunc) (*big.Int, error) {
				//if !peer.Equal(p) {
				if strings.Compare(peer, p) != 0 {
					t.Fatal("sending to wrong peer")
				}
				// if b != beneficiary {
				// 	t.Fatal("issuing for wrong beneficiary")
				// }
				if amount.Cmp(a) != 0 {
					t.Fatal("issuing with wrong amount")
				}
				return amount, nil
			},
		},
		store,
		mockvault.NewVault(),
		mockchequestore.NewChequeStore(),
		addressbook,
		networkID,
		&cashoutMock{},
		observer,
	)

	swapService.Pay(context.Background(), peer, amount, "", TOKEN)

	select {
	case call := <-observer.sentCalled:
		//if !call.peer.Equal(peer) {
		if strings.Compare(call.peer, peer) != 0 {
			t.Fatalf("observer called with wrong peer. got %v, want %v", call.peer, peer)
		}
		// ErrUnknownBeneficary has been blocked
		// if !errors.Is(call.err, swap.ErrUnknownBeneficary) {
		// 	t.Fatalf("wrong error. wanted %v, got %v", swap.ErrUnknownBeneficary, call.err)
		// }

	case <-time.After(time.Second):
		t.Fatal("expected observer to be called")
	}
}

func TestCashout(t *testing.T) {
	store := mockstore.NewStateStore()

	theirVaultAddress := common.HexToAddress("ffff")
	ourVaultAddress := common.HexToAddress("fffa")
	peer := peerInfo.ID("abcd").String()

	txHash := common.HexToHash("eeee")
	addressbook := &addressbookMock{
		vault: func(p string) (common.Address, bool, error) {
			//if !peer.Equal(p) {
			if strings.Compare(peer, p) != 0 {
				t.Fatal("querying vault for wrong peer")
			}
			return theirVaultAddress, true, nil
		},
	}

	swapService := swap.New(
		&swapProtocolMock{},
		store,
		mockvault.NewVault(
			mockvault.WithVaultAddressFunc(func() common.Address {
				return ourVaultAddress
			}),
		),
		mockchequestore.NewChequeStore(),
		addressbook,
		int64(1),
		&cashoutMock{
			cashCheque: func(ctx context.Context, c common.Address, r common.Address, token common.Address) (common.Hash, error) {
				if c != theirVaultAddress {
					t.Fatalf("not cashing with the right vault. wanted %v, got %v", theirVaultAddress, c)
				}
				if r != ourVaultAddress {
					t.Fatalf("not cashing with the right recipient. wanted %v, got %v", ourVaultAddress, r)
				}
				return txHash, nil
			},
		},
		nil,
	)

	returnedHash, err := swapService.CashCheque(context.Background(), peer, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	if returnedHash != txHash {
		t.Fatalf("go wrong tx hash. wanted %v, got %v", txHash, returnedHash)
	}
}

func TestCashoutStatus(t *testing.T) {
	store := mockstore.NewStateStore()

	theirVaultAddress := common.HexToAddress("ffff")
	peer := peerInfo.ID("abcd").String()
	addressbook := &addressbookMock{
		vault: func(p string) (common.Address, bool, error) {
			//if !peer.Equal(p) {
			if strings.Compare(peer, p) != 0 {
				t.Fatal("querying vault for wrong peer")
			}
			return theirVaultAddress, true, nil
		},
	}

	expectedStatus := &vault.CashoutStatus{}

	swapService := swap.New(
		&swapProtocolMock{},
		store,
		mockvault.NewVault(),
		mockchequestore.NewChequeStore(),
		addressbook,
		int64(1),
		&cashoutMock{
			cashoutStatus: func(ctx context.Context, c common.Address, token common.Address) (*vault.CashoutStatus, error) {
				if c != theirVaultAddress {
					t.Fatalf("getting status for wrong vault. wanted %v, got %v", theirVaultAddress, c)
				}
				return expectedStatus, nil
			},
		},
		nil,
	)

	returnedStatus, err := swapService.CashoutStatus(context.Background(), peer, TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	if expectedStatus != returnedStatus {
		t.Fatalf("go wrong status. wanted %v, got %v", expectedStatus, returnedStatus)
	}
}

func TestStateStoreKeys(t *testing.T) {
	address := common.HexToAddress("0xabcd")
	swarmAddress := peerInfo.ID("deff").String()

	expected := "swap_vault_peer_3ZqpcV"
	if swap.PeerKey(swarmAddress) != expected {
		t.Fatalf("wrong peer key. wanted %s, got %s", expected, swap.PeerKey(swarmAddress))
	}

	expected = "swap_peer_vault_000000000000000000000000000000000000abcd"
	if swap.VaultPeerKey(address) != expected {
		t.Fatalf("wrong peer key. wanted %s, got %s", expected, swap.VaultPeerKey(address))
	}

	expected = "swap_peer_beneficiary_3ZqpcV"
	if swap.PeerBeneficiaryKey(swarmAddress) != expected {
		t.Fatalf("wrong peer beneficiary key. wanted %s, got %s", expected, swap.PeerBeneficiaryKey(swarmAddress))
	}

	expected = "swap_beneficiary_peer_000000000000000000000000000000000000abcd"
	if swap.BeneficiaryPeerKey(address) != expected {
		t.Fatalf("wrong beneficiary peer key. wanted %s, got %s", expected, swap.BeneficiaryPeerKey(address))
	}

	expected = "swap_deducted_by_peer_3ZqpcV"
	if swap.PeerDeductedByKey(swarmAddress) != expected {
		t.Fatalf("wrong peer deducted by key. wanted %s, got %s", expected, swap.PeerDeductedByKey(swarmAddress))
	}

	expected = "swap_deducted_for_peer_3ZqpcV"
	if swap.PeerDeductedForKey(swarmAddress) != expected {
		t.Fatalf("wrong peer deducted for key. wanted %s, got %s", expected, swap.PeerDeductedForKey(swarmAddress))
	}
}
