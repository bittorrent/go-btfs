package vault

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/bittorrent/go-btfs/chain/tokencfg"

	"github.com/bittorrent/go-btfs/statestore"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/bittorrent/go-btfs/transaction/crypto"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/bittorrent/go-btfs/utils"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// prefix for the persistence key
	lastReceivedChequePrefix    = "swap_vault_last_received_cheque"
	receivedChequeHistoryPrefix = "swap_vault_history_received_cheque"
	//receivedChequeHistoryIndexPrefix = "swap_vault_history_received_cheque_index_"
	//180 days
	expireTime = 3600 * 24 * 180

	sendChequeHistoryPrefix = "swap_vault_history_send_cheque"
)

var (
	// ErrNoCheque is the error returned if there is no prior cheque for a vault or beneficiary.
	ErrNoCheque = errors.New("no cheque")
	// ErrNoChequeRecords is the error returned if there is no prior cheque record for a vault or beneficiary.
	ErrNoChequeRecords = errors.New("no cheque records")
	// ErrChequeNotIncreasing is the error returned if the cheque amount is the same or lower.
	ErrChequeNotIncreasing = errors.New("cheque cumulativePayout is not increasing")
	// ErrChequeInvalid is the error returned if the cheque itself is invalid.
	ErrChequeInvalid = errors.New("invalid cheque")
	// ErrWrongBeneficiary is the error returned if the cheque has the wrong beneficiary.
	ErrWrongBeneficiary = errors.New("wrong beneficiary")
	// ErrBouncingCheque is the error returned if the vault is demonstrably illiquid.
	ErrBouncingCheque = errors.New("bouncing cheque")
	ErrTokenCheque    = errors.New("wrong token cheque")
	// ErrChequeValueTooLow is the error returned if the after deduction value of a cheque did not cover 1 accounting credit
	ErrChequeValueTooLow = errors.New("cheque value lower than acceptable")
)

// ChequeStore handles the verification and storage of received cheques
type ChequeStore interface {
	// ReceiveCheque verifies and stores a cheque. It returns the total amount earned.
	ReceiveCheque(ctx context.Context, cheque *SignedCheque, amountCheck *big.Int, token common.Address) (*big.Int, error)
	// LastReceivedCheque returns the last cheque we received from a specific vault.
	LastReceivedCheque(vault common.Address, token common.Address) (*SignedCheque, error)
	// LastReceivedCheques return map[vault]cheque
	LastReceivedCheques(token common.Address) (map[common.Address]*SignedCheque, error)
	// ReceivedChequeRecordsByPeer returns the records we received from a specific vault.
	ReceivedChequeRecordsByPeer(vault common.Address) ([]ChequeRecord, error)
	// ListReceivedChequeRecords returns the records we received from a specific vault.
	ReceivedChequeRecordsAll() (map[common.Address][]ChequeRecord, error)
	ReceivedStatsHistory(days int, token common.Address) ([]DailyReceivedStats, error)
	SentStatsHistory(days int, token common.Address) ([]DailySentStats, error)

	// StoreSendChequeRecord store send cheque records.
	StoreSendChequeRecord(vault, beneficiary common.Address, amount *big.Int, token common.Address) error
	// SendChequeRecordsByPeer returns the records we send to a specific vault.
	SendChequeRecordsByPeer(beneficiary common.Address) ([]ChequeRecord, error)
	// SendChequeRecordsAll returns the records we send to a specific vault.
	SendChequeRecordsAll() (map[common.Address][]ChequeRecord, error)
}

type chequeStore struct {
	lock               sync.Mutex
	store              storage.StateStorer
	factory            Factory
	chaindID           int64
	transactionService transaction.Service
	beneficiary        common.Address // the beneficiary we expect in cheques sent to us
	recoverChequeFunc  RecoverChequeFunc
}

type RecoverChequeFunc func(cheque *SignedCheque, chainID int64) (common.Address, error)

// NewChequeStore creates new ChequeStore
func NewChequeStore(
	store storage.StateStorer,
	factory Factory,
	chainID int64,
	beneficiary common.Address,
	transactionService transaction.Service,
	recoverChequeFunc RecoverChequeFunc) ChequeStore {
	return &chequeStore{
		store:              store,
		factory:            factory,
		chaindID:           chainID,
		transactionService: transactionService,
		beneficiary:        beneficiary,
		recoverChequeFunc:  recoverChequeFunc,
	}
}

// lastReceivedChequeKey computes the key where to store the last cheque received from a vault.
func lastReceivedChequeKey(vault common.Address, token common.Address) string {
	return fmt.Sprintf("%s_%x", tokencfg.AddToken(lastReceivedChequePrefix, token), vault)
}

func historyReceivedChequeIndexKey(vault common.Address) string {
	return fmt.Sprintf("%s_%x", receivedChequeHistoryPrefix, vault)
}

func historyReceivedChequeKey(vault common.Address, index uint64) string {
	vaultStr := vault.String()
	return fmt.Sprintf("%s_%x", vaultStr, index)
}

// LastReceivedCheque map[vault]cheque
func (s *chequeStore) LastReceivedCheque(vault common.Address, token common.Address) (*SignedCheque, error) {
	var cheque *SignedCheque
	err := s.store.Get(lastReceivedChequeKey(vault, token), &cheque)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return nil, ErrNoCheque
	}

	return cheque, nil
}

// ReceiveCheque verifies and stores a cheque. It returns the totam amount earned.
func (s *chequeStore) ReceiveCheque(ctx context.Context, cheque *SignedCheque, amountCheck *big.Int, token common.Address) (*big.Int, error) {
	// verify we are the beneficiary
	if cheque.Beneficiary != s.beneficiary {
		return nil, ErrWrongBeneficiary
	}

	//if token != cheque.Token {
	//	return nil, ErrTokenCheque
	//}

	// don't allow concurrent processing of cheques
	// this would be sufficient on a per vault basis
	s.lock.Lock()
	defer s.lock.Unlock()

	// load the lastCumulativePayout for the cheques vault
	var lastCumulativePayout *big.Int
	var lastReceivedCheque *SignedCheque
	err := s.store.Get(lastReceivedChequeKey(cheque.Vault, token), &lastReceivedCheque)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}

		// if this is the first cheque from this vault, verify with the factory.
		err = s.factory.VerifyVault(ctx, cheque.Vault)
		if err != nil {
			return nil, err
		}

		lastCumulativePayout = big.NewInt(0)
	} else {
		lastCumulativePayout = lastReceivedCheque.CumulativePayout
	}
	// check this cheque is actually increasing in value by local storage
	amount := big.NewInt(0).Sub(cheque.CumulativePayout, lastCumulativePayout)
	if amount.Cmp(big.NewInt(0)) <= 0 {
		fmt.Println("amount < 0, ", ErrChequeNotIncreasing, cheque.CumulativePayout.String(), lastCumulativePayout.String(), token.Hex())
		return nil, ErrChequeNotIncreasing
	}

	//fmt.Println("... ReceiveCheque check, ", amount.String(), amountCheck.String())
	if amount.Cmp(amountCheck) < 0 {
		return nil, errors.New(
			fmt.Sprintf("amount[%v] is less than store amount[%v]. ", amount.String(), amountCheck.String()),
		)
	}

	// blockchain calls below
	contract := newVaultContractMuti(cheque.Vault, s.transactionService)
	// this does not change for the same vault
	expectedIssuer, err := contract.Issuer(ctx)
	if err != nil {
		return nil, err
	}

	// verify the cheque signature, (cheque contains token)
	issuer, err := s.recoverChequeFunc(cheque, s.chaindID)
	if err != nil {
		return nil, err
	}

	if issuer != expectedIssuer {
		return nil, ErrChequeInvalid
	}

	// basic balance check
	// could be omitted as it is not particularly useful
	balance, err := contract.TotalBalance(ctx, token)
	if err != nil {
		return nil, err
	}

	alreadyPaidOut, err := contract.PaidOut(ctx, s.beneficiary, token)
	if err != nil {
		return nil, err
	}

	if balance.Cmp(big.NewInt(0).Sub(cheque.CumulativePayout, alreadyPaidOut)) < 0 {
		return nil, ErrBouncingCheque
	}

	// check this cheque is actually increasing in value by blockchain in case of the host migrate to another machine
	// https://github.com/bittorrent/go-btfs/issues/187
	if cheque.CumulativePayout.Cmp(alreadyPaidOut) <= 0 {
		fmt.Println("CumulativePayout < alreadyPaidOut, ", ErrChequeNotIncreasing, cheque.CumulativePayout.String(), alreadyPaidOut.String(), token.Hex())
		return nil, ErrChequeNotIncreasing
	}

	// store the accepted cheque
	err = s.store.Put(lastReceivedChequeKey(cheque.Vault, token), cheque)
	if err != nil {
		return nil, err
	}

	// store the history cheque
	err = s.storeChequeRecord(cheque.Vault, amount, token)
	if err != nil {
		return nil, err
	}
	return amount, nil
}

// ReceivedChequeRecords returns the records we received from a specific vault.
func (s *chequeStore) ReceivedChequeRecordsByPeer(vault common.Address) ([]ChequeRecord, error) {
	var records []ChequeRecord
	var indexrange IndexRange
	err := s.store.Get(historyReceivedChequeIndexKey(vault), &indexrange)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return nil, ErrNoChequeRecords
	}

	for index := indexrange.MinIndex; index < indexrange.MaxIndex; index++ {
		var record ChequeRecord
		err = s.store.Get(historyReceivedChequeKey(vault, index), &record)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

// store cheque record
// Beneficiary common.Address
func (s *chequeStore) storeChequeRecord(vault common.Address, amount *big.Int, token common.Address) error {
	var indexRange IndexRange
	err := s.store.Get(historyReceivedChequeIndexKey(vault), &indexRange)
	if err != nil {
		if err != storage.ErrNotFound {
			return err
		}
		//not found
		indexRange.MinIndex = 0
		indexRange.MaxIndex = 0
		/*
			err = s.store.Put(historyReceivedChequeIndexKey(vault), indexRange)
			if err != nil {
				fmt.Println("put historyReceivedChequeIndexKey err ", err)
				return err
			}
		*/
	}

	//stroe cheque record with the key: historyReceivedChequeKey(index)
	chequeRecord := ChequeRecord{
		token,
		vault,
		s.beneficiary,
		amount,
		time.Now().Unix(),
	}

	err = s.store.Put(historyReceivedChequeKey(vault, indexRange.MaxIndex), chequeRecord)
	if err != nil {
		return err
	}

	//update Max : add one record
	indexRange.MaxIndex += 1
	//delete records if these record are old (half year)
	minIndex, _ := s.deleteRecordsExpired(vault, indexRange)

	//uopdate Min: add delete count
	indexRange.MinIndex = minIndex

	//update index
	err = s.store.Put(historyReceivedChequeIndexKey(vault), indexRange)
	if err != nil {
		return err
	}

	var stat DailyReceivedStats
	stat.Amount = big.NewInt(0)
	err = s.store.Get(statestore.GetTodayTotalDailyReceivedKey(token), &stat)
	if err != nil {
		if err != storage.ErrNotFound {
			return err
		}
		stat = DailyReceivedStats{
			Amount: big.NewInt(0),
			Count:  0,
		}
		stat.Date = utils.TodayUnix()
	}
	stat.Amount.Add(stat.Amount, amount)
	stat.Count += 1
	err = s.store.Put(statestore.GetTodayTotalDailyReceivedKey(token), stat)
	if err != nil {
		return err
	}

	// update uncashed records key
	uncashed := 0
	err = s.store.Get(statestore.PeerReceivedUncashRecordsCountKey(vault, token), &uncashed)
	if err != nil {
		if err != storage.ErrNotFound {
			return err
		}
	}
	uncashed += 1
	err = s.store.Put(statestore.PeerReceivedUncashRecordsCountKey(vault, token), uncashed)
	if err != nil {
		return err
	}

	return nil
}
func (s *chequeStore) ReceivedStatsHistory(days int, token common.Address) ([]DailyReceivedStats, error) {
	stats := make([]DailyReceivedStats, 0, days)
	y, m, d := time.Now().Date()
	todayStart := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	for i := 0; i < days; i++ {
		stat := DailyReceivedStats{}
		t := todayStart.AddDate(0, 0, -1*i)
		key := statestore.GetTotalDailyReceivedKeyByTime(t.Unix(), token)
		err := s.store.Get(key, &stat)
		if err != nil {
			if err != storage.ErrNotFound {
				return nil, err
			}
			stat = DailyReceivedStats{
				Amount: big.NewInt(0),
				Count:  0,
				Date:   t.Unix(),
			}
		}

		stats = append(stats, stat)
	}
	return stats, nil
}

func (s *chequeStore) SentStatsHistory(days int, token common.Address) ([]DailySentStats, error) {
	stats := make([]DailySentStats, 0, days)
	y, m, d := time.Now().Date()
	todayStart := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	for i := 0; i < days; i++ {
		stat := DailySentStats{}
		t := todayStart.AddDate(0, 0, -1*i)
		key := statestore.GetTotalDailySentKeyByTime(t.Unix(), token)
		err := s.store.Get(key, &stat)
		if err != nil {
			if err != storage.ErrNotFound {
				return nil, err
			}
			stat = DailySentStats{
				Amount: big.NewInt(0),
				Count:  0,
				Date:   t.Unix(),
			}
		}

		stats = append(stats, stat)
	}
	return stats, nil
}

func (s *chequeStore) deleteRecordsExpired(vault common.Address, indexRange IndexRange) (uint64, error) {
	//get the expire time
	expire := time.Now().Unix() - expireTime
	var chequeRecord ChequeRecord
	var endIndex uint64

	//find the last index expired to delete
	for index := indexRange.MinIndex; index < indexRange.MaxIndex; index++ {
		err := s.store.Get(historyReceivedChequeKey(vault, index), &chequeRecord)
		if err != nil {
			return indexRange.MinIndex, err
		}

		if chequeRecord.ReceiveTime >= expire {
			endIndex = index
			break
		}
	}

	//delete [min endIndex) records
	if endIndex <= indexRange.MinIndex {
		return indexRange.MinIndex, nil
	}

	//delete expired records
	for index := indexRange.MinIndex; index < endIndex; index++ {
		err := s.store.Delete(historyReceivedChequeKey(vault, index))
		if err != nil {
			return indexRange.MinIndex, err
		}
		//min++
		indexRange.MinIndex += 1
	}

	return indexRange.MinIndex, nil
}

// RecoverCheque recovers the issuer ethereum address from a signed cheque
func RecoverCheque(cheque *SignedCheque, chaindID int64) (common.Address, error) {
	eip712Data := eip712DataForCheque(&cheque.Cheque, chaindID)

	pubkey, err := crypto.RecoverEIP712(cheque.Signature, eip712Data)
	if err != nil {
		return common.Address{}, err
	}

	ethAddr, err := crypto.NewEthereumAddress(*pubkey)
	if err != nil {
		return common.Address{}, err
	}

	var issuer common.Address
	copy(issuer[:], ethAddr)
	return issuer, nil
}

// keyVault computes the vault a store entry is for.
func keyVault(key []byte, prefix string) (vault common.Address, err error) {
	k := string(key)

	split := strings.SplitAfter(k, prefix)
	if len(split) != 2 {
		return common.Address{}, errors.New("no peer in key")
	}
	return common.HexToAddress(split[1]), nil
}

// LastCheques returns map[vault]cheque
func (s *chequeStore) LastReceivedCheques(token common.Address) (map[common.Address]*SignedCheque, error) {
	result := make(map[common.Address]*SignedCheque)
	err := s.store.Iterate(tokencfg.AddToken(lastReceivedChequePrefix, token), func(key, val []byte) (stop bool, err error) {
		// if token is wbtt, need a wbtt key, other token is not this.

		addr, err := keyVault(key, tokencfg.AddToken(lastReceivedChequePrefix, token)+"_")
		if err != nil {
			return false, fmt.Errorf("parse address from key: %s: %w", string(key), err)
		}

		if _, ok := result[addr]; !ok {
			lastCheque, err := s.LastReceivedCheque(addr, token)
			if err != nil && err != ErrNoCheque {
				return false, err
			} else if err == ErrNoCheque {
				return false, nil
			}

			result[addr] = lastCheque
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListReceivedChequeRecords returns the last received cheques from every known vault.
func (s *chequeStore) ReceivedChequeRecordsAll() (map[common.Address][]ChequeRecord, error) {
	result := make(map[common.Address][]ChequeRecord)
	err := s.store.Iterate(receivedChequeHistoryPrefix, func(key, val []byte) (stop bool, err error) {
		addr, err := keyVault(key, receivedChequeHistoryPrefix+"_")
		if err != nil {
			return false, fmt.Errorf("parse address from key: %s: %w", string(key), err)
		}

		if _, ok := result[addr]; !ok {
			records, err := s.ReceivedChequeRecordsByPeer(addr)
			if err != nil && err != ErrNoCheque && err != ErrNoChequeRecords {
				return false, err
			} else if err == ErrNoCheque || err == ErrNoChequeRecords {
				return false, nil
			}

			result[addr] = records
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func historySendChequeIndexKey(beneficiary common.Address) string {
	return fmt.Sprintf("%s_%x", sendChequeHistoryPrefix, beneficiary)
}

func historySendChequeKey(beneficiary common.Address, index uint64) string {
	beneficiaryStr := beneficiary.String()
	return fmt.Sprintf("%s_%x", beneficiaryStr, index)
}

// store cheque record
// Beneficiary common.Address
func (s *chequeStore) StoreSendChequeRecord(vault, beneficiary common.Address, amount *big.Int, token common.Address) error {
	var indexRange IndexRange
	err := s.store.Get(historySendChequeIndexKey(beneficiary), &indexRange)
	if err != nil {
		if err != storage.ErrNotFound {
			return err
		}
		//not found
		indexRange.MinIndex = 0
		indexRange.MaxIndex = 0
		/*
			err = s.store.Put(historyReceivedChequeIndexKey(vault), indexRange)
			if err != nil {
				fmt.Println("put historyReceivedChequeIndexKey err ", err)
				return err
			}
		*/
	}

	//fmt.Printf("...1 StoreSendChequeRecord, historySendChequeIndexKey=%v, indexRange=%v, token=%v \n",
	//	historySendChequeIndexKey(beneficiary), indexRange, token.Hex())

	//stroe cheque record with the key: historySendChequeKey(index)
	chequeRecord := ChequeRecord{
		token,
		vault,
		beneficiary,
		amount,
		time.Now().Unix(),
	}

	err = s.store.Put(historySendChequeKey(beneficiary, indexRange.MaxIndex), chequeRecord)
	if err != nil {
		return err
	}

	//fmt.Printf("...2 StoreSendChequeRecord, historySendChequeKey=%v, indexRange=%v \n",
	//	historySendChequeKey(beneficiary, indexRange.MaxIndex), indexRange)

	//update Max : add one record
	indexRange.MaxIndex += 1
	//delete records if these record are old (half year)
	minIndex, _ := s.deleteSendRecordsExpired(beneficiary, indexRange)

	//uopdate Min: add delete count
	indexRange.MinIndex = minIndex

	//update index
	err = s.store.Put(historySendChequeIndexKey(beneficiary), indexRange)
	if err != nil {
		return err
	}

	//fmt.Printf("...3 StoreSendChequeRecord, historySendChequeKey=%v, indexRange=%v \n",
	//	historySendChequeKey(beneficiary, indexRange.MaxIndex), indexRange)

	var stat DailySentStats
	err = s.store.Get(statestore.GetTodayTotalDailySentKey(token), &stat)
	if err != nil {
		if err != storage.ErrNotFound {
			return err
		}
		stat = DailySentStats{
			Amount: big.NewInt(0),
			Count:  0,
		}
		stat.Date = utils.TodayUnix()
	}
	stat.Amount.Add(stat.Amount, amount)
	stat.Count += 1
	err = s.store.Put(statestore.GetTodayTotalDailySentKey(token), stat)
	if err != nil {
		return err
	}

	return nil
}

func (s *chequeStore) deleteSendRecordsExpired(beneficiary common.Address, indexRange IndexRange) (uint64, error) {
	//get the expire time
	expire := time.Now().Unix() - expireTime
	var chequeRecord ChequeRecord
	var endIndex uint64

	//find the last index expired to delete
	for index := indexRange.MinIndex; index < indexRange.MaxIndex; index++ {
		err := s.store.Get(historySendChequeKey(beneficiary, index), &chequeRecord)
		if err != nil {
			return indexRange.MinIndex, err
		}

		if chequeRecord.ReceiveTime >= expire {
			endIndex = index
			break
		}
	}

	//delete [min endIndex) records
	if endIndex <= indexRange.MinIndex {
		return indexRange.MinIndex, nil
	}

	//delete expired records
	for index := indexRange.MinIndex; index < endIndex; index++ {
		err := s.store.Delete(historySendChequeKey(beneficiary, index))
		if err != nil {
			return indexRange.MinIndex, err
		}
		//min++
		indexRange.MinIndex += 1
	}

	return indexRange.MinIndex, nil
}

// SendChequeRecordsByPeer returns the records we received from a specific vault.
func (s *chequeStore) SendChequeRecordsByPeer(beneficiary common.Address) ([]ChequeRecord, error) {
	var records []ChequeRecord
	var indexrange IndexRange
	err := s.store.Get(historySendChequeIndexKey(beneficiary), &indexrange)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return nil, ErrNoChequeRecords
	}

	for index := indexrange.MinIndex; index < indexrange.MaxIndex; index++ {
		var record ChequeRecord
		err = s.store.Get(historySendChequeKey(beneficiary, index), &record)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

// SendChequeRecordsAll returns the last send cheques from every known vault.
func (s *chequeStore) SendChequeRecordsAll() (map[common.Address][]ChequeRecord, error) {
	result := make(map[common.Address][]ChequeRecord)
	err := s.store.Iterate(sendChequeHistoryPrefix, func(key, val []byte) (stop bool, err error) {
		addr, err := keyVault(key, sendChequeHistoryPrefix+"_")
		if err != nil {
			return false, fmt.Errorf("parse address from key: %s: %w", string(key), err)
		}

		if _, ok := result[addr]; !ok {
			records, err := s.SendChequeRecordsByPeer(addr)
			if err != nil && err != ErrNoCheque && err != ErrNoChequeRecords {
				return false, err
			} else if err == ErrNoCheque || err == ErrNoChequeRecords {
				return false, nil
			}

			result[addr] = records
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	//fmt.Println("SendChequeRecordsAll ... result = ", result)
	//fmt.Println("SendChequeRecordsAll ... result = ", sendChequeHistoryPrefix, sendChequeHistoryPrefix+"_")

	return result, nil
}
