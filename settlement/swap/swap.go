package swap

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"math/big"
	"sync"

	"github.com/bittorrent/go-btfs/settlement"
	"github.com/bittorrent/go-btfs/settlement/swap/swapprotocol"
	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	"github.com/bittorrent/go-btfs/statestore"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/ethereum/go-ethereum/common"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("swap")

var (
	// ErrWrongVault is the error if a peer uses a different vault from before.
	ErrWrongVault = errors.New("wrong vault")
	// ErrUnknownBeneficary is the error if a peer has never announced a beneficiary.
	ErrUnknownBeneficary = errors.New("unknown beneficiary for peer")
	// ErrChequeValueTooLow is the error a peer issued a cheque not covering 1 accounting credit
	ErrChequeValueTooLow = errors.New("cheque value too low")
)

type Interface interface {
	settlement.Interface

	// LastReceivedCheque returns the last received cheque for the peer
	LastReceivedCheque(peer string) (*vault.SignedCheque, error)
	// LastReceivedCheques returns the list of last received cheques for all peers
	LastReceivedCheques() (map[string]*vault.SignedCheque, error)
	// ReceivedChequeRecordsByPeer gets the records of the cheque for the peer's vault
	ReceivedChequeRecordsByPeer(peer string) ([]vault.ChequeRecord, error)
	// ReceivedChequeRecordsAll gets the records of the cheque for all vault
	ReceivedChequeRecordsAll() ([]vault.ChequeRecord, error)
	// ReceivedChequeRecordsCount
	ReceivedChequeRecordsCount() (int, error)

	// LastSendCheque returns the last received cheque for the peer
	LastSendCheque(peer string) (*vault.SignedCheque, error)
	// LastSendCheques returns the list of last send cheques for this peer
	LastSendCheques() (map[string]*vault.SignedCheque, error)
	// SendChequeRecordsByPeer gets the records of the cheque for the peer's vault
	SendChequeRecordsByPeer(peer string) ([]vault.ChequeRecord, error)
	// SendChequeRecordsAll gets the records of the cheque for all vault
	SendChequeRecordsAll() ([]vault.ChequeRecord, error)

	// CashCheque sends a cashing transaction for the last cheque of the peer
	CashCheque(ctx context.Context, peer string) (common.Hash, error)
	// CashoutStatus gets the status of the latest cashout transaction for the peers vault
	CashoutStatus(ctx context.Context, peer string) (*vault.CashoutStatus, error)
	// HasCashoutAction
	HasCashoutAction(ctx context.Context, peer string) (bool, error)
}

// Service is the implementation of the swap settlement layer.
type Service struct {
	lock        sync.Mutex
	proto       swapprotocol.Interface
	store       storage.StateStorer
	accounting  settlement.Accounting
	metrics     metrics
	vault       vault.Service
	chequeStore vault.ChequeStore
	cashout     vault.CashoutService
	addressbook Addressbook
	chainID     int64
}

// New creates a new swap Service.
func New(proto swapprotocol.Interface, store storage.StateStorer, vault vault.Service, chequeStore vault.ChequeStore, addressbook Addressbook, chainID int64, cashout vault.CashoutService, accounting settlement.Accounting) *Service {
	return &Service{
		proto:       proto,
		store:       store,
		metrics:     newMetrics(),
		vault:       vault,
		chequeStore: chequeStore,
		addressbook: addressbook,
		chainID:     chainID,
		cashout:     cashout,
		accounting:  accounting,
	}
}

func (s *Service) GetProtocols() swapprotocol.Interface {
	return s.proto
}

// ReceiveCheque is called by the swap protocol if a cheque is received.
func (s *Service) ReceiveCheque(ctx context.Context, peer string, cheque *vault.SignedCheque, realAmount *big.Int, token common.Address) (err error) {
	// check this is the same vault for this peer as previously
	expectedVault, known, err := s.addressbook.Vault(peer)
	if err != nil {
		return err
	}

	if known && expectedVault != cheque.Vault {
		return ErrWrongVault
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	receivedAmount, err := s.chequeStore.ReceiveCheque(ctx, cheque, realAmount, token)
	if err != nil {
		s.metrics.ChequesRejected.Inc()
		return fmt.Errorf("rejecting cheque: %w", err)
	}

	if !known {
		err = s.addressbook.PutVault(peer, cheque.Vault)
		if err != nil {
			return err
		}
	}

	tot, _ := big.NewFloat(0).SetInt(receivedAmount).Float64()
	s.metrics.TotalReceived.Add(tot)
	s.metrics.ChequesReceived.Inc()

	// total received count
	totalReceivedCount, err := s.vault.TotalReceivedCount(token)
	if err != nil {
		return err
	}
	totalReceivedCount = totalReceivedCount + 1
	err = s.store.Put(tokencfg.AddToken(statestore.TotalReceivedCountKey, token), totalReceivedCount)
	if err != nil {
		return err
	}

	// totalReceived
	totalReceived, err := s.vault.TotalReceived(token)
	if err != nil {
		return err
	}
	totalReceived = totalReceived.Add(totalReceived, receivedAmount)
	err = s.store.Put(tokencfg.AddToken(statestore.TotalReceivedKey, token), totalReceived)
	if err != nil {
		return err
	}

	return s.accounting.NotifyPaymentReceived(peer, receivedAmount, token)
}

// Pay initiates a payment to the given peer
func (s *Service) Pay(ctx context.Context, peer string, amount *big.Int, contractId string, token common.Address) {
	var err error
	defer func() {
		if err != nil {
			s.accounting.NotifyPaymentSent(peer, amount, err, token)
		}
	}()

	/*
		beneficiary, known, err := s.addressbook.Beneficiary(peer)
		if err != nil {
			return
		}
		if !known {
			err = ErrUnknownBeneficary
			return
		}
	*/
	balance, err := s.proto.EmitCheque(ctx, peer, amount, contractId, token, s.vault.Issue)

	if err != nil {
		return
	}

	bal, _ := big.NewFloat(0).SetInt(balance).Float64()
	s.metrics.AvailableBalance.Set(bal)
	s.accounting.NotifyPaymentSent(peer, amount, nil, token)
	amountFloat, _ := big.NewFloat(0).SetInt(amount).Float64()
	s.metrics.TotalSent.Add(amountFloat)
	s.metrics.ChequesSent.Inc()
}

func (s *Service) SetAccounting(accounting settlement.Accounting) {
	s.accounting = accounting
}

// TotalSent returns the total amount sent to a peer
func (s *Service) TotalSent(peer string, token common.Address) (totalSent *big.Int, err error) {
	beneficiary, known, err := s.addressbook.Beneficiary(peer)
	if err != nil {
		return nil, err
	}
	if !known {
		return nil, settlement.ErrPeerNoSettlements
	}
	cheque, err := s.vault.LastCheque(beneficiary, token)
	if err != nil {
		if err == vault.ErrNoCheque {
			return nil, settlement.ErrPeerNoSettlements
		}
		return nil, err
	}
	return cheque.CumulativePayout, nil
}

// TotalReceived returns the total amount received from a peer
func (s *Service) TotalReceived(peer string, token common.Address) (totalReceived *big.Int, err error) {
	vaultAddress, known, err := s.addressbook.Vault(peer)
	if err != nil {
		return nil, err
	}
	if !known {
		return nil, settlement.ErrPeerNoSettlements
	}

	cheque, err := s.chequeStore.LastReceivedCheque(vaultAddress, token)
	if err != nil {
		if err == vault.ErrNoCheque {
			return nil, settlement.ErrPeerNoSettlements
		}
		return nil, err
	}
	return cheque.CumulativePayout, nil
}

// SettlementsSent returns sent settlements for each individual known peer
func (s *Service) SettlementsSent(token common.Address) (map[string]*big.Int, error) {
	result := make(map[string]*big.Int)
	cheques, err := s.vault.LastCheques(token)
	if err != nil {
		return nil, err
	}

	for beneficiary, cheque := range cheques {
		peer, known, err := s.addressbook.BeneficiaryPeer(beneficiary)
		if err != nil {
			return nil, err
		}
		if !known {
			continue
		}
		result[peer] = cheque.CumulativePayout
	}

	return result, nil
}

// SettlementsReceived returns received settlements for each individual known peer.
func (s *Service) SettlementsReceived(token common.Address) (map[string]*big.Int, error) {
	result := make(map[string]*big.Int)
	cheques, err := s.chequeStore.LastReceivedCheques(token)
	if err != nil {
		return nil, err
	}

	for vault, cheque := range cheques {
		peer, known, err := s.addressbook.VaultPeer(vault)
		if err != nil {
			return nil, err
		}
		if !known {
			continue
		}
		result[peer] = cheque.CumulativePayout
	}
	return result, err
}

//// Handshake is called by the swap protocol when a handshake is received.
//func (s *Service) Handshake(peer string, beneficiary common.Address) error {
//	oldPeer, known, err := s.addressbook.BeneficiaryPeer(beneficiary)
//	fmt.Println("Handshake 1, peer, oldPeer, beneficiary: ", peer, oldPeer, beneficiary)
//
//	if err != nil {
//		return err
//	}
//	if known && strings.Compare(peer, oldPeer) != 0 {
//		log.Debugf("migrating swap addresses from peer %s to %s", oldPeer, peer)
//
//		fmt.Println("Handshake 2")
//
//		return s.addressbook.MigratePeer(oldPeer, peer)
//	}
//
//	_, known, err = s.addressbook.Beneficiary(peer)
//	fmt.Println("Handshake 3", known, err)
//
//	if err != nil {
//		return err
//	}
//	if !known {
//		log.Infof("initial swap handshake peer: %v beneficiary: %x", peer, beneficiary)
//		return s.addressbook.PutBeneficiary(peer, beneficiary)
//	}
//
//	return nil
//}

// LastReceivedCheque returns the last received cheque for the peer
func (s *Service) LastReceivedCheque(peer string, token common.Address) (*vault.SignedCheque, error) {

	common, known, err := s.addressbook.Vault(peer)

	if err != nil {
		return nil, err
	}

	if !known {
		return nil, vault.ErrNoCheque
	}

	return s.chequeStore.LastReceivedCheque(common, token)
}

// LastReceivedCheques returns map[peer]cheque
func (s *Service) LastReceivedCheques(token common.Address) (map[string]*vault.SignedCheque, error) {
	lastcheques, err := s.chequeStore.LastReceivedCheques(token)
	if err != nil {
		return nil, err
	}

	resultmap := make(map[string]*vault.SignedCheque, len(lastcheques))

	for i, j := range lastcheques {
		addr, known, err := s.addressbook.VaultPeer(i)
		if err == nil && known {
			resultmap[addr] = j
		}
	}

	return resultmap, nil
}

// LastReceivedCheque returns the last received cheque for the peer
func (s *Service) ReceivedChequeRecordsByPeer(peer string) ([]vault.ChequeRecord, error) {
	common, known, err := s.addressbook.Vault(peer)
	if err != nil {
		return nil, err
	}

	if !known {
		return nil, vault.ErrNoCheque
	}

	return s.chequeStore.ReceivedChequeRecordsByPeer(common)
}

// ReceivedChequeRecordsAll returns the last received cheque for the peer
func (s *Service) ReceivedChequeRecordsAll() ([]vault.ChequeRecord, error) {
	mp, err := s.chequeStore.ReceivedChequeRecordsAll()
	if err != nil {
		return nil, err
	}

	records := make([]vault.ChequeRecord, 0)
	for comm, _ := range mp {
		l, err := s.chequeStore.ReceivedChequeRecordsByPeer(comm)
		if err != nil {
			return nil, err
		}

		for _, a := range l {
			records = append(records, a)
		}
	}

	return records, nil
}

// ReceivedChequeRecordsCount returns the last received cheque for the peer
func (s *Service) ReceivedChequeRecordsCount() (int, error) {
	mp, err := s.chequeStore.ReceivedChequeRecordsAll()
	if err != nil {
		return 0, err
	}

	count := 0
	for comm, _ := range mp {
		l, err := s.chequeStore.ReceivedChequeRecordsByPeer(comm)
		if err != nil {
			return 0, err
		}

		count += len(l)
	}

	return count, nil
}

// LastReceivedCheque returns the last received cheque for the peer
func (s *Service) LastSendCheque(peer string, token common.Address) (*vault.SignedCheque, error) {
	comm, known, err := s.addressbook.Vault(peer)
	if err != nil {
		return nil, err
	}

	if !known {
		return nil, vault.ErrNoCheque
	}

	return s.vault.LastCheque(comm, token)
}

// LastReceivedCheques returns the list of last received cheques for all peers
func (s *Service) LastSendCheques(token common.Address) (map[string]*vault.SignedCheque, error) {
	lastcheques, err := s.vault.LastCheques(token)
	fmt.Println("LastSendCheques ", lastcheques, err)

	if err != nil {
		return nil, err
	}

	resultmap := make(map[string]*vault.SignedCheque, len(lastcheques))

	for i, j := range lastcheques {
		addr, known, err := s.addressbook.BeneficiaryPeer(i)
		if err == nil && known {
			resultmap[addr] = j
		}
	}

	return resultmap, nil
}

// SendChequeRecordsByPeer returns the last received cheque for the peer
func (s *Service) SendChequeRecordsByPeer(peer string) ([]vault.ChequeRecord, error) {
	common, known, err := s.addressbook.Beneficiary(peer)
	if err != nil {
		return nil, err
	}

	if !known {
		return nil, vault.ErrNoCheque
	}

	return s.chequeStore.SendChequeRecordsByPeer(common)
}

// SendChequeRecordsAll returns the last received cheque for all peer
func (s *Service) SendChequeRecordsAll() ([]vault.ChequeRecord, error) {
	mp, err := s.chequeStore.SendChequeRecordsAll()
	if err != nil {
		return nil, err
	}

	records := make([]vault.ChequeRecord, 0)
	for comm, _ := range mp {
		l, err := s.chequeStore.SendChequeRecordsByPeer(comm)
		if err != nil {
			return nil, err
		}

		for _, a := range l {
			records = append(records, a)
		}
	}

	return records, nil
}

// CashCheque sends a cashing transaction for the last cheque of the peer
func (s *Service) CashCheque(ctx context.Context, peer string, token common.Address) (common.Hash, error) {
	vaultAddress, known, err := s.addressbook.Vault(peer)
	if err != nil {
		return common.Hash{}, err
	}
	if !known {
		return common.Hash{}, vault.ErrNoCheque
	}

	fmt.Println("...CashCheque, vaultAddress, s.vault.Address(), token ", vaultAddress, s.vault.Address(), token)
	return s.cashout.CashCheque(ctx, vaultAddress, s.vault.Address(), token)
}

// CashoutStatus gets the status of the latest cashout transaction for the peers vault
func (s *Service) CashoutStatus(ctx context.Context, peer string, token common.Address) (*vault.CashoutStatus, error) {
	vaultAddress, known, err := s.addressbook.Vault(peer)
	if err != nil {
		return nil, err
	}
	if !known {
		return nil, vault.ErrNoCheque
	}
	return s.cashout.CashoutStatus(ctx, vaultAddress, token)
}

func (s *Service) GetChainid() int64 {
	return s.chainID
}

func (s *Service) Settle(toPeer string, paymentAmount *big.Int, contractId string, token common.Address) error {
	return s.accounting.Settle(toPeer, paymentAmount, contractId, token)
}

func (s *Service) PutBeneficiary(peer string, beneficiary common.Address) (common.Address, error) {
	last_address, known, err := s.addressbook.Beneficiary(peer)
	if err != nil {
		return beneficiary, err
	}
	if !known {
		return beneficiary, s.addressbook.PutBeneficiary(peer, beneficiary)
	}

	//compare last_address and beneficiary
	if !bytes.Equal(last_address.Bytes(), beneficiary.Bytes()) {
		s.addressbook.MigratePeer(last_address.String(), beneficiary.String())
	}

	return beneficiary, nil
}

func (s *Service) Beneficiary(peer string) (beneficiary common.Address, known bool, err error) {
	return s.addressbook.Beneficiary(peer)
}

func (s *Service) VaultPeer(vault common.Address) (peer string, known bool, err error) {
	return s.addressbook.VaultPeer(vault)
}

func (s *Service) BeneficiaryPeer(beneficiary common.Address) (peer string, known bool, err error) {
	return s.addressbook.BeneficiaryPeer(beneficiary)
}

func (s *Service) HasCashoutAction(ctx context.Context, peer string, token common.Address) (bool, error) {
	vault, known, err := s.addressbook.Vault(peer)
	if err != nil {
		return false, err
	}
	if !known {
		return false, fmt.Errorf("unkonw peer")
	}
	return s.cashout.HasCashoutAction(ctx, vault, token)
}
