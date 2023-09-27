package vault

import (
	"context"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/status-im/keycard-go/hexutils"
	"math/big"
	"time"

	"github.com/bittorrent/go-btfs/statestore"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	// ErrNoCashout is the error if there has not been any cashout action for the vault
	ErrNoCashout = errors.New("no prior cashout")
)

// CashoutService is the service responsible for managing cashout actions
type CashoutService interface {
	// CashCheque sends a cashing transaction for the last cheque of the vault
	CashCheque(ctx context.Context, vault, recipient common.Address, token common.Address) (common.Hash, error)
	// CashoutStatus gets the status of the latest cashout transaction for the vault
	CashoutStatus(ctx context.Context, vaultAddress common.Address, token common.Address) (*CashoutStatus, error)
	AdjustCashCheque(ctx context.Context, vaultAddress, recipient common.Address, token common.Address) (totalCashOutAmount, newCashOutAmount *big.Int, err error)
	AdjustCashChequeTxHash(ctx context.Context, vaultAddress, recipient common.Address, token common.Address, txHash common.Hash, cumulativePayout *big.Int) (totalCashOutAmount, newCashOutAmount *big.Int, err error)
	HasCashoutAction(ctx context.Context, peer common.Address, token common.Address) (bool, error)
	CashoutResults() ([]CashOutResult, error)
	RestartFixChequeCashOut()
}

type cashoutService struct {
	store              storage.StateStorer
	backend            transaction.Backend
	transactionService transaction.Service
	chequeStore        ChequeStore
}

// LastCashout contains information about the last cashout
type LastCashout struct {
	TxHash   common.Hash
	Cheque   SignedCheque // the cheque that was used to cashout which may be different from the latest cheque
	Result   *CashChequeResult
	Reverted bool
}

// CashoutStatus is information about the last cashout and uncashed amounts
type CashoutStatus struct {
	Last           *LastCashout // last cashout for a vault
	UncashedAmount *big.Int     // amount not yet cashed out
}

// CashChequeResult summarizes the result of a CashCheque or CashChequeBeneficiary call
type CashChequeResult struct {
	Beneficiary      common.Address // beneficiary of the cheque
	Recipient        common.Address // address which received the funds
	Caller           common.Address // caller of cashCheque
	TotalPayout      *big.Int       // total amount that was paid out in this call
	CumulativePayout *big.Int       // cumulative payout of the cheque that was cashed
	CallerPayout     *big.Int       // payout for the caller of cashCheque
	Bounced          bool           // indicates wether parts of the cheque bounced
}

// cashoutAction is the data we store for a cashout
type cashoutAction struct {
	TxHash common.Hash
	Cheque SignedCheque // the cheque that was used to cashout which may be different from the latest cheque
}

type CashOutResult struct {
	TxHash   common.Hash
	Token    common.Address
	Vault    common.Address
	Amount   *big.Int
	CashTime int64
	Status   string
}

type chequeCashedEvent struct {
	Beneficiary      common.Address
	Recipient        common.Address
	Caller           common.Address
	TotalPayout      *big.Int
	CumulativePayout *big.Int
	CallerPayout     *big.Int
}

type mutiChequeCashedEvent struct {
	Token            common.Address
	Beneficiary      common.Address
	Recipient        common.Address
	Caller           common.Address
	TotalPayout      *big.Int
	CumulativePayout *big.Int
	CallerPayout     *big.Int
}

type mutiChequeBouncedEvent struct {
	Token common.Address
}

// NewCashoutService creates a new CashoutService
func NewCashoutService(
	store storage.StateStorer,
	backend transaction.Backend,
	transactionService transaction.Service,
	chequeStore ChequeStore,
) CashoutService {
	return &cashoutService{
		store:              store,
		backend:            backend,
		transactionService: transactionService,
		chequeStore:        chequeStore,
	}
}

// cashoutActionKey computes the store key for the last cashout action for the vault
func cashoutActionKey(vault common.Address, token common.Address) string {
	return tokencfg.AddToken(fmt.Sprintf("swap_cashout_%x", vault), token)
}

// paidOut (dropped 2.3.0)
//func (s *cashoutService) paidOut(ctx context.Context, vault, beneficiary common.Address) (*big.Int, error) {
//	callData, err := vaultABI.Pack("paidOut", beneficiary)
//	if err != nil {
//		return nil, err
//	}
//
//	output, err := s.transactionService.Call(ctx, &transaction.TxRequest{
//		To:   &vault,
//		Data: callData,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	results, err := vaultABI.Unpack("paidOut", output)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(results) != 1 {
//		return nil, errDecodeABI
//	}
//
//	paidOut, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
//	if !ok || paidOut == nil {
//		return nil, errDecodeABI
//	}
//
//	return paidOut, nil
//}

// paidOutMuti (2.3.0 import)
func (s *cashoutService) paidOutMuti(ctx context.Context, vault, beneficiary common.Address, token common.Address) (*big.Int, error) {
	return _PaidOutMuti(ctx, vault, beneficiary, s.transactionService, token)
}

func (s *cashoutService) CashoutResults() ([]CashOutResult, error) {
	result := make([]CashOutResult, 0, 0)
	err := s.store.Iterate(statestore.CashoutResultPrefixKey(), func(key, val []byte) (stop bool, err error) {
		cashOutResult := CashOutResult{}
		err = s.store.Get(string(key), &cashOutResult)
		if err != nil {
			return false, err
		}
		result = append(result, cashOutResult)
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CashCheque sends a cashout transaction for the last cheque of the vault
func (s *cashoutService) CashCheque(ctx context.Context, vault, recipient common.Address, token common.Address) (common.Hash, error) {
	if RestartFixCashOutStatusLock {
		return common.Hash{}, errors.New("Just started, it can not cash cheque for processing the cash cheque out status last time, you will wait for about 40s to do it. ")
	}

	cheque, err := s.chequeStore.LastReceivedCheque(vault, token)
	if err != nil {
		return common.Hash{}, err
	}

	fmt.Println("_CashCheque ", vault, recipient, cheque.CumulativePayout, hexutils.BytesToHex(cheque.Signature), cheque.Token)

	//callData, err := vaultABI.Pack("cashChequeBeneficiary", recipient, cheque.CumulativePayout, cheque.Signature)
	//if err != nil {
	//	return common.Hash{}, err
	//}
	//request := &transaction.TxRequest{
	//	To:          &vault,
	//	Data:        callData,
	//	Value:       big.NewInt(0),
	//	Description: "cheque cashout",
	//}
	//
	//txHash, err := s.transactionService.Send(ctx, request)
	//if err != nil {
	//	return common.Hash{}, err
	//}

	// 2.3.0 import
	txHash, err := _CashChequeMuti(ctx, vault, recipient, cheque, s.transactionService, token)
	if err != nil {
		return common.Hash{}, err
	}

	err = s.store.Put(cashoutActionKey(vault, token), &cashoutAction{
		TxHash: txHash,
		Cheque: *cheque,
	})
	if err != nil {
		return common.Hash{}, err
	}

	// 1.add cash out status
	cashOutStateInfo := CashOutStatusStoreInfo{
		Token:            token,
		Vault:            cheque.Vault,
		Beneficiary:      cheque.Beneficiary,
		CumulativePayout: cheque.CumulativePayout,
		TxHash:           txHash.String(),
	}
	err = s.AddCashOutStatusStore(cashOutStateInfo)
	if err != nil {
		return common.Hash{}, err
	}

	// WaitForReceipt takes long time
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("storeCashResult recovered:%+v", err)
			}
		}()
		s.storeCashResult(context.Background(), vault, txHash, cheque, token)

		// 2.delete cash out status
		err = s.DeleteCashOutStatusStore(cashOutStateInfo)
		if err != nil {
			log.Errorf("delete cashout status, err = %v", err)
			return
		}
	}()
	return txHash, nil
}

func (s *cashoutService) storeCashResult(ctx context.Context, vault common.Address, txHash common.Hash, cheque *SignedCheque, token common.Address) error {
	cashResult := CashOutResult{
		TxHash:   txHash,
		Vault:    vault,
		Token:    token,
		Amount:   cheque.CumulativePayout,
		CashTime: time.Now().Unix(),
		Status:   "fail",
	}

	_, err := s.transactionService.WaitForReceipt(ctx, txHash)
	if err != nil {
		log.Infof("storeCashResult err:%+v", err)
	} else {
		cs, err := s.CashoutStatus(ctx, vault, token)
		if err != nil {
			log.Infof("CashOutStats:get cashout status err:%+v", err)
			if cs.UncashedAmount != nil {
				cashResult.Amount = cs.UncashedAmount
			}
		} else {
			// update totalReceivedCashed
			totalPaidOut := big.NewInt(0)
			if cs.Last != nil && cs.Last.Result != nil && cs.Last.Result.TotalPayout != nil {
				totalPaidOut = cs.Last.Result.TotalPayout
			}
			if cs.Last != nil && !cs.Last.Reverted {
				cashResult.Status = "success"
			}
			cashResult.Amount = totalPaidOut
			totalReceivedCashed := big.NewInt(0)
			if err = s.store.Get(tokencfg.AddToken(statestore.TotalReceivedCashedKey, token), &totalReceivedCashed); err == nil || err == storage.ErrNotFound {
				totalReceivedCashed = totalReceivedCashed.Add(totalReceivedCashed, totalPaidOut)
				err := s.store.Put(tokencfg.AddToken(statestore.TotalReceivedCashedKey, token), totalReceivedCashed)
				if err != nil {
					log.Infof("CashOutStats:put totalReceivedCashdKey err:%+v", err)
				}
			}

			totalDailyReceivedCashed := big.NewInt(0)
			if err = s.store.Get(statestore.GetTodayTotalDailyReceivedCashedKey(token), &totalDailyReceivedCashed); err == nil || err == storage.ErrNotFound {
				totalDailyReceivedCashed = totalDailyReceivedCashed.Add(totalDailyReceivedCashed, totalPaidOut)
				err := s.store.Put(statestore.GetTodayTotalDailyReceivedCashedKey(token), totalDailyReceivedCashed)
				if err != nil {
					log.Infof("CashOutStats:put totalReceivedDailyCashdKey err:%+v", err)
				}
			}

			// update TotalReceivedCountCashed
			uncashed := 0
			err := s.store.Get(statestore.PeerReceivedUncashRecordsCountKey(vault, token), &uncashed)
			if err != nil {
				log.Infof("CashOutStats:put totalReceivedCountCashed err:%+v", err)
			} else {
				cashedCount := 0
				err := s.store.Get(tokencfg.AddToken(statestore.TotalReceivedCashedCountKey, token), &cashedCount)
				if err == nil || err == storage.ErrNotFound {
					err := s.store.Put(tokencfg.AddToken(statestore.TotalReceivedCashedCountKey, token), cashedCount+uncashed)
					if err != nil {
						log.Infof("CashOutStats:put totalReceivedCashedConuntKey err:%+v", err)
					} else {
						err := s.store.Put(statestore.PeerReceivedUncashRecordsCountKey(vault, token), 0)
						if err != nil {
							log.Infof("CashOutStats:put totalReceivedCashedConuntKey err:%+v", err)
						}
					}
				}
			}
		}
	}
	err = s.store.Put(statestore.CashoutResultKey(vault), &cashResult)
	if err != nil {
		log.Infof("CashOutStats:put cashoutResultKey err:%+v", err)
	}
	return nil
}

// AdjustCashCheque .
func (s *cashoutService) AdjustCashCheque(ctx context.Context, vaultAddress, recipient common.Address, token common.Address) (totalCashOutAmount, newCashOutAmount *big.Int, err error) {
	if RestartFixCashOutStatusLock {
		return nil, nil, errors.New("Just started, it can not fix cash out info for processing the cash cheque out status last time, you will wait for about 40s to do it. ")
	}

	// 1.totalReceivedCashed
	totalReceivedCashed := big.NewInt(0)
	err = s.store.Get(tokencfg.AddToken(statestore.TotalReceivedCashedKey, token), &totalReceivedCashed)
	if err != nil && err != storage.ErrNotFound {
		return nil, nil, err
	}

	// 2.alreadyPaidOut in renter contract
	// blockchain calls below
	contract := newVaultContractMuti(vaultAddress, s.transactionService)
	alreadyPaidOutOnline, err := contract.PaidOut(ctx, recipient, token)
	if err != nil {
		return nil, nil, err
	}

	// 3.compare it to fix.
	diff := big.NewInt(0).Sub(alreadyPaidOutOnline, totalReceivedCashed)
	log.Infof("AdjustCashCheque: diff > 0, vault=%s, recipient=%s, online=%s, local=%s, diff=%s",
		vaultAddress.String(), recipient.String(),
		alreadyPaidOutOnline.String(), totalReceivedCashed.String(), diff.String(),
	)

	if diff.Cmp(big.NewInt(0)) > 0 {
		txHash := common.Hash{} //fix txHash: 0x0000...
		cashResult, err := s.fixStoreCashResult(vaultAddress, diff, token, txHash)
		if err != nil {
			return nil, nil, err
		}
		newCashOutAmount = cashResult.Amount
	}

	return alreadyPaidOutOnline, newCashOutAmount, nil
}

// AdjustCashChequeTxHash . when node starts.
func (s *cashoutService) AdjustCashChequeTxHash(ctx context.Context, vaultAddress, recipient common.Address, token common.Address, txHash common.Hash, cumulativePayout *big.Int) (totalCashOutAmount, newCashOutAmount *big.Int, err error) {
	// 1.totalReceivedCashed
	totalReceivedCashed := big.NewInt(0)
	err = s.store.Get(tokencfg.AddToken(statestore.TotalReceivedCashedKey, token), &totalReceivedCashed)
	if err != nil && err != storage.ErrNotFound {
		return nil, nil, err
	}

	// 2.alreadyPaidOut in renter contract
	// blockchain calls below
	contract := newVaultContractMuti(vaultAddress, s.transactionService)
	alreadyPaidOutOnline, err := contract.PaidOut(ctx, recipient, token)
	if err != nil {
		return nil, nil, err
	}

	// 3.compare it to fix.
	diffReceived := big.NewInt(0).Sub(alreadyPaidOutOnline, totalReceivedCashed)
	log.Infof("AdjustCashCheque: diffReceived, vault=%s, recipient=%s, txHash=%s, online=%s, local=%s, diffReceived=%s",
		vaultAddress.String(), recipient.String(), txHash.String(),
		alreadyPaidOutOnline.String(), totalReceivedCashed.String(), diffReceived.String(),
	)

	if diffReceived.Cmp(big.NewInt(0)) > 0 {
		diffCheque := big.NewInt(0).Sub(alreadyPaidOutOnline, cumulativePayout)
		log.Infof("AdjustCashCheque: diffCheque, vault=%s, recipient=%s, txHash=%s, online=%s, cumulativePayout=%s, diffCheque=%s",
			vaultAddress.String(), recipient.String(), txHash.String(),
			alreadyPaidOutOnline.String(), cumulativePayout.String(), diffCheque.String(),
		)
		if diffCheque.Cmp(big.NewInt(0)) > 0 {
			txHash = common.Hash{} //fix txHash: 0x0000...
		}
		cashResult, err := s.fixStoreCashResult(vaultAddress, diffReceived, token, txHash)
		if err != nil {
			return nil, nil, err
		}
		newCashOutAmount = cashResult.Amount
	}

	return alreadyPaidOutOnline, newCashOutAmount, nil
}

func (s *cashoutService) RestartFixChequeCashOut() {
	if RestartFixCashOutStatusLock {
		list, err := s.GetAllCashOutStatusStore()
		if err != nil {
			log.Infof("RestartFixChequeCashOut: GetAllCashOutStatusStore err = %v", err)
			return
		}

		if len(list) > 0 {
			log.Infof("wait 30s, for fixing cash out status")

			// wait 30s,  for online cashing out ok.
			time.Sleep(time.Second * RestartWaitCashOutOnlineTime)

			for _, v := range list {
				txHash := common.HexToHash(v.TxHash)

				// 1.check txHash is ok
				receipt, err := s.backend.TransactionReceipt(context.Background(), txHash)
				if err != nil {
					log.Infof("RestartFixChequeCashOut: TransactionReceipt err = %v, info = %+v", err, v)
					continue
				}
				if receipt.Status == types.ReceiptStatusFailed {
					log.Infof("RestartFixChequeCashOut: TransactionReceipt err = the txHash is failed, info = %+v", v)
					continue
				}

				// 2.adjust cash cheque info
				_, _, err = s.AdjustCashChequeTxHash(context.Background(), v.Vault, v.Beneficiary, v.Token, txHash, v.CumulativePayout)
				if err != nil {
					log.Infof("RestartFixChequeCashOut: AdjustCashCheque err = %v, info = %+v", err, v)
					continue
				}

				// 3.delete cash out status
				err = s.DeleteCashOutStatusStore(v)
				if err != nil {
					log.Infof("RestartFixChequeCashOut: DeleteCashOutStatusStore err = %v, info = %+v", err, v)
					continue
				}
			}
		}
		RestartFixCashOutStatusLock = false
	}
	return
}

func (s *cashoutService) fixStoreCashResult(vault common.Address, shouldPaidOut *big.Int, token common.Address, txHash common.Hash) (cashResult *CashOutResult, err error) {
	//txHash := common.Hash{} //fix txHash: 0x0000...
	cashResult = &CashOutResult{
		TxHash:   txHash,
		Vault:    vault,
		Token:    token,
		Amount:   shouldPaidOut,
		CashTime: time.Now().Unix(),
		Status:   "success",
	}

	totalReceivedCashed := big.NewInt(0)
	if err = s.store.Get(tokencfg.AddToken(statestore.TotalReceivedCashedKey, token), &totalReceivedCashed); err == nil || err == storage.ErrNotFound {
		totalReceivedCashed = totalReceivedCashed.Add(totalReceivedCashed, shouldPaidOut)
		err := s.store.Put(tokencfg.AddToken(statestore.TotalReceivedCashedKey, token), totalReceivedCashed)
		if err != nil {
			log.Infof("fixStoreCashResult:put totalReceivedCashdKey err:%+v", err)
		}
	}

	totalDailyReceivedCashed := big.NewInt(0)
	if err = s.store.Get(statestore.GetTodayTotalDailyReceivedCashedKey(token), &totalDailyReceivedCashed); err == nil || err == storage.ErrNotFound {
		totalDailyReceivedCashed = totalDailyReceivedCashed.Add(totalDailyReceivedCashed, shouldPaidOut)
		err := s.store.Put(statestore.GetTodayTotalDailyReceivedCashedKey(token), totalDailyReceivedCashed)
		if err != nil {
			log.Infof("fixStoreCashResult:put totalReceivedDailyCashdKey err:%+v", err)
		}
	}

	// update TotalReceivedCountCashed
	uncashed := 0
	err = s.store.Get(statestore.PeerReceivedUncashRecordsCountKey(vault, token), &uncashed)
	if err != nil {
		log.Infof("fixStoreCashResult:put totalReceivedCountCashed err:%+v", err)
	} else {
		cashedCount := 0
		err := s.store.Get(tokencfg.AddToken(statestore.TotalReceivedCashedCountKey, token), &cashedCount)
		if err == nil || err == storage.ErrNotFound {
			err := s.store.Put(tokencfg.AddToken(statestore.TotalReceivedCashedCountKey, token), cashedCount+uncashed)
			if err != nil {
				log.Infof("fixStoreCashResult:put totalReceivedCashedConuntKey err:%+v", err)
			} else {
				err := s.store.Put(statestore.PeerReceivedUncashRecordsCountKey(vault, token), 0)
				if err != nil {
					log.Infof("fixStoreCashResult:put totalReceivedCashedConuntKey err:%+v", err)
				}
			}
		}
	}

	err = s.store.Put(statestore.CashoutResultKey(vault), &cashResult)
	if err != nil {
		log.Infof("fixStoreCashResult:put cashoutResultKey err:%+v", err)
	}
	return
}

// CashoutStatus gets the status of the latest cashout transaction for the vault
func (s *cashoutService) CashoutStatus(ctx context.Context, vaultAddress common.Address, token common.Address) (*CashoutStatus, error) {
	cheque, err := s.chequeStore.LastReceivedCheque(vaultAddress, token)
	if err != nil {
		return nil, err
	}

	var action cashoutAction
	err = s.store.Get(cashoutActionKey(vaultAddress, token), &action)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &CashoutStatus{
				Last:           nil,
				UncashedAmount: cheque.CumulativePayout, // if we never cashed out, assume everything is uncashed
			}, nil
		}
		return nil, err
	}

	_, pending, err := s.backend.TransactionByHash(ctx, action.TxHash)
	if err != nil {
		// treat not found as pending
		if !errors.Is(err, ethereum.NotFound) {
			return nil, err
		}
		pending = true
	}

	if pending {
		return &CashoutStatus{
			Last: &LastCashout{
				TxHash:   action.TxHash,
				Cheque:   action.Cheque,
				Result:   nil,
				Reverted: false,
			},
			// uncashed is the difference since the last sent cashout. we assume that the entire cheque will clear in the pending transaction.
			UncashedAmount: new(big.Int).Sub(cheque.CumulativePayout, action.Cheque.CumulativePayout),
		}, nil
	}

	receipt, err := s.backend.TransactionReceipt(ctx, action.TxHash)
	if err != nil {
		return nil, err
	}

	if receipt.Status == types.ReceiptStatusFailed {
		// if a tx failed (should be almost impossible in practice) we no longer have the necessary information to compute uncashed locally
		// assume there are no pending transactions and that the on-chain paidOut is the last cashout action
		paidOut, err := s.paidOutMuti(ctx, vaultAddress, cheque.Beneficiary, token)
		if err != nil {
			return nil, err
		}

		return &CashoutStatus{
			Last: &LastCashout{
				TxHash:   action.TxHash,
				Cheque:   action.Cheque,
				Result:   nil,
				Reverted: true,
			},
			UncashedAmount: new(big.Int).Sub(cheque.CumulativePayout, paidOut),
		}, nil
	}

	result, err := s.parseCashChequeBeneficiaryReceiptMuti(vaultAddress, receipt, token)
	if err != nil {
		return nil, err
	}

	return &CashoutStatus{
		Last: &LastCashout{
			TxHash:   action.TxHash,
			Cheque:   action.Cheque,
			Result:   result,
			Reverted: false,
		},
		// uncashed is the difference since the last sent (and confirmed) cashout.
		UncashedAmount: new(big.Int).Sub(cheque.CumulativePayout, result.CumulativePayout),
	}, nil
}

// parseCashChequeBeneficiaryReceipt processes the receipt from a CashChequeBeneficiary transaction
func (s *cashoutService) parseCashChequeBeneficiaryReceipt(vaultAddress common.Address, receipt *types.Receipt) (*CashChequeResult, error) {
	result := &CashChequeResult{
		Bounced: false,
	}

	var cashedEvent chequeCashedEvent
	err := transaction.FindSingleEvent(&vaultABI, receipt, vaultAddress, chequeCashedEventType, &cashedEvent)
	if err != nil {
		return nil, err
	}

	result.Beneficiary = cashedEvent.Beneficiary
	result.Caller = cashedEvent.Caller
	result.CallerPayout = cashedEvent.CallerPayout
	result.TotalPayout = cashedEvent.TotalPayout
	result.CumulativePayout = cashedEvent.CumulativePayout
	result.Recipient = cashedEvent.Recipient

	err = transaction.FindSingleEvent(&vaultABI, receipt, vaultAddress, chequeBouncedEventType, nil)
	if err == nil {
		result.Bounced = true
	} else if !errors.Is(err, transaction.ErrEventNotFound) {
		return nil, err
	}

	return result, nil
}

// parseCashChequeBeneficiaryReceiptMuti processes the receipt from a CashChequeBeneficiary transaction
func (s *cashoutService) parseCashChequeBeneficiaryReceiptMuti(vaultAddress common.Address, receipt *types.Receipt, token common.Address) (*CashChequeResult, error) {
	if tokencfg.IsWBTT(token) {
		return s.parseCashChequeBeneficiaryReceipt(vaultAddress, receipt)
	}

	result := &CashChequeResult{
		Bounced: false,
	}

	var mtCashedEvent mutiChequeCashedEvent
	err := transaction.FindSingleEvent(&vaultABINew, receipt, vaultAddress, mutiChequeCashedEventType, &mtCashedEvent)
	if err != nil {
		return nil, err
	}

	result.Beneficiary = mtCashedEvent.Beneficiary
	result.Caller = mtCashedEvent.Caller
	result.CallerPayout = mtCashedEvent.CallerPayout
	result.TotalPayout = mtCashedEvent.TotalPayout
	result.CumulativePayout = mtCashedEvent.CumulativePayout
	result.Recipient = mtCashedEvent.Recipient

	//var mtBouncedEvent mutiChequeBouncedEvent
	err = transaction.FindSingleEvent(&vaultABINew, receipt, vaultAddress, mutiChequeBouncedEventType, nil)
	if err == nil {
		result.Bounced = true
	} else if !errors.Is(err, transaction.ErrEventNotFound) {
		return nil, err
	}

	return result, nil
}

// Equal compares to CashChequeResults
func (r *CashChequeResult) Equal(o *CashChequeResult) bool {
	if r.Beneficiary != o.Beneficiary {
		return false
	}
	if r.Bounced != o.Bounced {
		return false
	}
	if r.Caller != o.Caller {
		return false
	}
	if r.CallerPayout.Cmp(o.CallerPayout) != 0 {
		return false
	}
	if r.CumulativePayout.Cmp(o.CumulativePayout) != 0 {
		return false
	}
	if r.Recipient != o.Recipient {
		return false
	}
	if r.TotalPayout.Cmp(o.TotalPayout) != 0 {
		return false
	}
	return true
}

func (s *cashoutService) HasCashoutAction(ctx context.Context, peer common.Address, token common.Address) (bool, error) {
	var action cashoutAction
	err := s.store.Get(cashoutActionKey(peer, token), &action)

	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
