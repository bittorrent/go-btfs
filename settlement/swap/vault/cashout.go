package vault

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/TRON-US/go-btfs/transaction"
	"github.com/TRON-US/go-btfs/transaction/storage"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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
	CashCheque(ctx context.Context, vault, recipient common.Address) (common.Hash, error)
	// CashoutStatus gets the status of the latest cashout transaction for the vault
	CashoutStatus(ctx context.Context, vaultAddress common.Address) (*CashoutStatus, error)
	HasCashoutAction(ctx context.Context, peer common.Address) (bool, error)
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
type chequeCashedEvent struct {
	Beneficiary      common.Address
	Recipient        common.Address
	Caller           common.Address
	TotalPayout      *big.Int
	CumulativePayout *big.Int
	CallerPayout     *big.Int
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
func cashoutActionKey(vault common.Address) string {
	return fmt.Sprintf("swap_cashout_%x", vault)
}

func (s *cashoutService) paidOut(ctx context.Context, vault, beneficiary common.Address) (*big.Int, error) {
	callData, err := vaultABI.Pack("paidOut", beneficiary)
	if err != nil {
		return nil, err
	}

	output, err := s.transactionService.Call(ctx, &transaction.TxRequest{
		To:   &vault,
		Data: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := vaultABI.Unpack("paidOut", output)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, errDecodeABI
	}

	paidOut, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || paidOut == nil {
		return nil, errDecodeABI
	}

	return paidOut, nil
}

// CashCheque sends a cashout transaction for the last cheque of the vault
func (s *cashoutService) CashCheque(ctx context.Context, vault, recipient common.Address) (common.Hash, error) {
	cheque, err := s.chequeStore.LastReceivedCheque(vault)
	if err != nil {
		return common.Hash{}, err
	}

	callData, err := vaultABI.Pack("cashChequeBeneficiary", recipient, cheque.CumulativePayout, cheque.Signature)
	if err != nil {
		return common.Hash{}, err
	}
	request := &transaction.TxRequest{
		To:          &vault,
		Data:        callData,
		Value:       big.NewInt(0),
		Description: "cheque cashout",
	}

	txHash, err := s.transactionService.Send(ctx, request)
	if err != nil {
		return common.Hash{}, err
	}

	err = s.store.Put(cashoutActionKey(vault), &cashoutAction{
		TxHash: txHash,
		Cheque: *cheque,
	})
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

// CashoutStatus gets the status of the latest cashout transaction for the vault
func (s *cashoutService) CashoutStatus(ctx context.Context, vaultAddress common.Address) (*CashoutStatus, error) {
	cheque, err := s.chequeStore.LastReceivedCheque(vaultAddress)
	if err != nil {
		return nil, err
	}

	var action cashoutAction
	err = s.store.Get(cashoutActionKey(vaultAddress), &action)
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
		paidOut, err := s.paidOut(ctx, vaultAddress, cheque.Beneficiary)
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

	result, err := s.parseCashChequeBeneficiaryReceipt(vaultAddress, receipt)
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

func (s *cashoutService) HasCashoutAction(ctx context.Context, peer common.Address) (bool, error) {
	var action cashoutAction
	err := s.store.Get(cashoutActionKey(peer), &action)

	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}