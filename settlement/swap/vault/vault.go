package vault

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"

	conabi "github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/settlement/swap/erc20"
	"github.com/bittorrent/go-btfs/statestore"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// SendChequeFunc is a function to send cheques.
type SendChequeFunc func(cheque *SignedCheque) error

const (
	lastIssuedChequeKeyPrefix = "swap_vault_last_issued_cheque_"
	totalIssuedKey            = "swap_vault_total_issued_"
	totalIssuedCountKey       = "swap_vault_total_issued_count_"
)

var (
	// ErrOutOfFunds is the error when the vault has not enough free funds for a cheque
	ErrOutOfFunds = errors.New("vault out of funds")
	// ErrInsufficientFunds is the error when the vault has not enough free funds for a user action
	ErrInsufficientFunds = errors.New("insufficient token balance")

	vaultABI               = transaction.ParseABIUnchecked(conabi.VaultABI)
	chequeCashedEventType  = vaultABI.Events["ChequeCashed"]
	chequeBouncedEventType = vaultABI.Events["ChequeBounced"]
)

// Service is the main interface for interacting with the nodes vault.
type Service interface {
	// Deposit starts depositing erc20 token into the vault. This returns once the transactions has been broadcast.
	Deposit(ctx context.Context, amount *big.Int) (hash common.Hash, err error)
	// Withdraw starts withdrawing erc20 token from the vault. This returns once the transactions has been broadcast.
	Withdraw(ctx context.Context, amount *big.Int) (hash common.Hash, err error)
	// WaitForDeposit waits for the deposit transaction to confirm and verifies the result.
	WaitForDeposit(ctx context.Context, txHash common.Hash) error
	// TotalBalance returns the token balance of the vault.
	TotalBalance(ctx context.Context) (*big.Int, error)
	// TotalIssuedCount returns total issued count of the vault.
	TotalIssuedCount() (int, error)
	TotalIssued() (*big.Int, error)
	TotalReceivedCount() (int, error)
	TotalReceivedCashedCount() (int, error)
	TotalReceived() (*big.Int, error)
	TotalReceivedCashed() (*big.Int, error)
	TotalDailyReceived() (*big.Int, error)
	TotalDailyReceivedCashed() (*big.Int, error)
	// LiquidBalance returns the token balance of the vault sub stake amount.
	LiquidBalance(ctx context.Context) (*big.Int, error)
	// AvailableBalance returns the token balance of the vault which is not yet used for uncashed cheques.
	AvailableBalance(ctx context.Context) (*big.Int, error)
	// Address returns the address of the used vault contract.
	Address() common.Address
	// Issue a new cheque for the beneficiary with an cumulativePayout amount higher than the last.
	Issue(ctx context.Context, beneficiary common.Address, amount *big.Int, sendChequeFunc SendChequeFunc) (*big.Int, error)
	// LastCheque returns the last cheque we issued for the beneficiary.
	LastCheque(beneficiary common.Address) (*SignedCheque, error)
	// LastCheques returns the last cheques we issued for all beneficiaries.
	LastCheques() (map[common.Address]*SignedCheque, error)
	// GetWithdrawTime returns the time can withdraw
	GetWithdrawTime(ctx context.Context) (*big.Int, error)
	// WbttBalanceOf retrieve the addr balance
	WBTTBalanceOf(ctx context.Context, addr common.Address) (*big.Int, error)
	// BTTBalanceOf retrieve the btt balance of addr
	BTTBalanceOf(ctx context.Context, address common.Address, block *big.Int) (*big.Int, error)
	// TotalPaidOut return total pay out of the vault
	TotalPaidOut(ctx context.Context) (*big.Int, error)
	// CheckBalance
	CheckBalance(amount *big.Int) (err error)
	// UpgradeTo will upgrade vault implementation to `newVaultImpl`
	UpgradeTo(ctx context.Context, newVaultImpl common.Address) (old, new common.Address, err error)
}

type service struct {
	lock               sync.Mutex
	transactionService transaction.Service

	address      common.Address
	contract     *vaultContract
	ownerAddress common.Address

	erc20Service erc20.Service

	store               storage.StateStorer
	chequeSigner        ChequeSigner
	totalIssuedReserved *big.Int
	chequeStore         ChequeStore
}

// New creates a new vault service for the provided vault contract.
func New(transactionService transaction.Service, address, ownerAddress common.Address, store storage.StateStorer,
	chequeSigner ChequeSigner, erc20Service erc20.Service, chequeStore ChequeStore) (Service, error) {
	return &service{
		transactionService:  transactionService,
		address:             address,
		contract:            newVaultContract(address, transactionService),
		ownerAddress:        ownerAddress,
		erc20Service:        erc20Service,
		store:               store,
		chequeSigner:        chequeSigner,
		totalIssuedReserved: big.NewInt(0),
		chequeStore:         chequeStore,
	}, nil
}

// Address returns the address of the used vault contract.
func (s *service) Address() common.Address {
	return s.address
}

// Deposit starts depositing erc20 token into the vault. This returns once the transactions has been broadcast.
func (s *service) Deposit(ctx context.Context, amount *big.Int) (hash common.Hash, err error) {
	balance, err := s.erc20Service.BalanceOf(ctx, s.ownerAddress)
	if err != nil {
		return common.Hash{}, err
	}

	// check we can afford this so we don't waste gas
	if balance.Cmp(amount) < 0 {
		return common.Hash{}, ErrInsufficientFunds
	}

	return s.contract.Deposit(ctx, amount)
}

// Deposit starts depositing erc20 token into the vault. This returns once the transactions has been broadcast.
func (s *service) CheckBalance(amount *big.Int) (err error) {
	balance, err := s.erc20Service.BalanceOf(context.Background(), s.ownerAddress)
	if err != nil {
		return err
	}

	// check we can afford this so we don't waste gas
	if balance.Cmp(amount) < 0 {
		return ErrInsufficientFunds
	}

	return nil
}

// Balance returns the token balance of the vault.
func (s *service) TotalBalance(ctx context.Context) (*big.Int, error) {
	return s.contract.TotalBalance(ctx)
}

// LiquidBalance returns the token balance of the vault sub stake amount.
func (s *service) LiquidBalance(ctx context.Context) (*big.Int, error) {
	return s.contract.LiquidBalance(ctx)
}

// AvailableBalance returns the token balance of the vault which is not yet used for uncashed cheques.
func (s *service) AvailableBalance(ctx context.Context) (*big.Int, error) {
	totalIssued, err := s.totalIssued()
	if err != nil {
		return nil, err
	}

	balance, err := s.TotalBalance(ctx)
	if err != nil {
		return nil, err
	}

	totalPaidOut, err := s.contract.TotalPaidOut(ctx)
	if err != nil {
		return nil, err
	}

	// balance plus totalPaidOut is the total amount ever put into the vault (ignoring deposits and withdrawals which cancelled out)
	// minus the total amount we issued from this vault this gives use the portion of the balance not covered by any cheques
	availableBalance := big.NewInt(0).Add(balance, totalPaidOut)
	availableBalance = availableBalance.Sub(availableBalance, totalIssued)
	return availableBalance, nil
}

// total send cheque count.  returns the token balance of the vault which is not yet used for uncashed cheques.
func (s *service) TotalIssuedCount() (int, error) {
	totalIssuedCount, err := s.totalIssuedCount()
	if err != nil {
		return 0, err
	}

	return totalIssuedCount, nil
}

func (s *service) TotalIssued() (*big.Int, error) {
	return s.totalIssued()
}

// total recevied cheque count.
func (s *service) TotalReceivedCount() (int, error) {
	return s.totalReceivedCount()
}

func (s *service) TotalReceivedCashedCount() (int, error) {
	return s.totalReceivedCashedCount()
}

func (s *service) TotalReceived() (*big.Int, error) {
	return s.totalReceived()
}

func (s *service) TotalReceivedCashed() (*big.Int, error) {
	return s.totalReceivedCashed()
}

func (s *service) TotalDailyReceived() (*big.Int, error) {
	return s.totalDailyReceived()
}

func (s *service) TotalDailyReceivedCashed() (*big.Int, error) {
	return s.totalDailyReceivedCashed()
}

// WaitForDeposit waits for the deposit transaction to confirm and verifies the result.

// WaitForDeposit waits for the deposit transaction to confirm and verifies the result.
func (s *service) WaitForDeposit(ctx context.Context, txHash common.Hash) error {
	receipt, err := s.transactionService.WaitForReceipt(ctx, txHash)
	if err != nil {
		return err
	}
	if receipt.Status != 1 {
		return transaction.ErrTransactionReverted
	}
	return nil
}

// lastIssuedChequeKey computes the key where to store the last cheque for a beneficiary.
func lastIssuedChequeKey(beneficiary common.Address) string {
	return fmt.Sprintf("%s%x", lastIssuedChequeKeyPrefix, beneficiary)
}

func (s *service) reserveTotalIssued(ctx context.Context, amount *big.Int) (*big.Int, error) {
	availableBalance, err := s.AvailableBalance(ctx)
	if err != nil {
		return nil, err
	}

	if amount.Cmp(big.NewInt(0).Sub(availableBalance, s.totalIssuedReserved)) > 0 {
		return nil, ErrOutOfFunds
	}

	s.totalIssuedReserved = s.totalIssuedReserved.Add(s.totalIssuedReserved, amount)
	return big.NewInt(0).Sub(availableBalance, amount), nil
}

func (s *service) unreserveTotalIssued(amount *big.Int) {
	s.totalIssuedReserved = s.totalIssuedReserved.Sub(s.totalIssuedReserved, amount)
}

// Issue issues a new cheque and passes it to sendChequeFunc.
// The cheque is considered sent and saved when sendChequeFunc succeeds.
// The available balance which is available after sending the cheque is passed
// to the caller for it to be communicated over metrics.
func (s *service) Issue(ctx context.Context, beneficiary common.Address, amount *big.Int, sendChequeFunc SendChequeFunc) (*big.Int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	availableBalance, err := s.reserveTotalIssued(ctx, amount)
	if err != nil {
		return nil, err
	}
	defer s.unreserveTotalIssued(amount)

	var cumulativePayout *big.Int
	lastCheque, err := s.LastCheque(beneficiary)
	if err != nil {
		if err != ErrNoCheque {
			return nil, err
		}
		cumulativePayout = big.NewInt(0)
	} else {
		// 1.set last
		cumulativePayout = lastCheque.CumulativePayout

		// 2.get sum of send peer cheque
		records, err := s.chequeStore.SendChequeRecordsByPeer(beneficiary)
		if err != nil {
			return nil, err
		}
		var sendAmount = big.NewInt(0)
		for _, v := range records {
			sendAmount = new(big.Int).Add(sendAmount, v.Amount)
		}

		// 3.if send > last, set send.
		if sendAmount.Cmp(lastCheque.CumulativePayout) > 0 {
			cumulativePayout = sendAmount
		}

		fmt.Println("... issue, send cheque sendAmount, lastCheque.CumulativePayout, diff = ",
			sendAmount,
			lastCheque.CumulativePayout,
			sendAmount.Cmp(lastCheque.CumulativePayout))
	}

	// increase cumulativePayout by amount
	cumulativePayout = cumulativePayout.Add(cumulativePayout, amount)

	// create and sign the new cheque
	cheque := Cheque{
		Vault:            s.address,
		CumulativePayout: cumulativePayout,
		Beneficiary:      beneficiary,
	}

	sig, err := s.chequeSigner.Sign(&Cheque{
		Vault:            s.address,
		CumulativePayout: cumulativePayout,
		Beneficiary:      beneficiary,
	})
	if err != nil {
		return nil, err
	}

	// actually send the check before saving to avoid double payment
	err = sendChequeFunc(&SignedCheque{
		Cheque:    cheque,
		Signature: sig,
	})
	if err != nil {
		return nil, err
	}

	err = s.store.Put(lastIssuedChequeKey(beneficiary), cheque)
	if err != nil {
		return nil, err
	}

	// store the history issued cheque
	err = s.chequeStore.StoreSendChequeRecord(s.address, beneficiary, amount)
	if err != nil {
		return nil, err
	}

	// total issued count
	totalIssuedCount, err := s.totalIssuedCount()
	if err != nil {
		return nil, err
	}
	totalIssuedCount = totalIssuedCount + 1
	err = s.store.Put(totalIssuedCountKey, totalIssuedCount)
	if err != nil {
		return nil, err
	}

	// totalIssued
	totalIssued, err := s.totalIssued()
	if err != nil {
		return nil, err
	}
	totalIssued = totalIssued.Add(totalIssued, amount)
	return availableBalance, s.store.Put(totalIssuedKey, totalIssued)
}

// returns the total amount in cheques issued so far
func (s *service) totalIssued() (totalIssued *big.Int, err error) {
	err = s.store.Get(totalIssuedKey, &totalIssued)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return big.NewInt(0), nil
	}
	return totalIssued, nil
}

// returns the total count in cheques issued so far
func (s *service) totalIssuedCount() (totalIssuedCount int, err error) {
	err = s.store.Get(totalIssuedCountKey, &totalIssuedCount)
	if err != nil {
		if err != storage.ErrNotFound {
			return 0, err
		}
		return 0, nil
	}
	return totalIssuedCount, nil
}

// returns the total amount in cheques recieved so far
func (s *service) totalReceived() (totalReceived *big.Int, err error) {
	err = s.store.Get(statestore.TotalReceivedKey, &totalReceived)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return big.NewInt(0), nil
	}
	return totalReceived, nil
}

func (s *service) totalReceivedCashed() (totalReceived *big.Int, err error) {
	err = s.store.Get(statestore.TotalReceivedCashedKey, &totalReceived)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return big.NewInt(0), nil
	}
	return totalReceived, nil
}

// returns the total amount in cheques recieved so far
func (s *service) totalDailyReceived() (totalReceived *big.Int, err error) {
	var stat DailyReceivedStats
	err = s.store.Get(statestore.GetTodayTotalDailyReceivedKey(), &stat)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return big.NewInt(0), nil
	}
	return stat.Amount, nil
}

func (s *service) totalDailyReceivedCashed() (totalReceived *big.Int, err error) {
	err = s.store.Get(statestore.GetTodayTotalDailyReceivedCashedKey(), &totalReceived)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return big.NewInt(0), nil
	}
	return totalReceived, nil
}

// returns the total count in cheques recieved so far
func (s *service) totalReceivedCount() (totalReceivedCount int, err error) {
	err = s.store.Get(statestore.TotalReceivedCountKey, &totalReceivedCount)
	if err != nil {
		if err != storage.ErrNotFound {
			return 0, err
		}
		return 0, nil
	}
	return totalReceivedCount, nil
}

func (s *service) totalReceivedCashedCount() (totalReceivedCount int, err error) {
	err = s.store.Get(statestore.TotalReceivedCashedCountKey, &totalReceivedCount)
	if err != nil {
		if err != storage.ErrNotFound {
			return 0, err
		}
		return 0, nil
	}
	return totalReceivedCount, nil
}

// LastCheque returns the last cheque we issued for the beneficiary.
func (s *service) LastCheque(beneficiary common.Address) (*SignedCheque, error) {
	var lastCheque *SignedCheque
	err := s.store.Get(lastIssuedChequeKey(beneficiary), &lastCheque)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return nil, ErrNoCheque
	}
	return lastCheque, nil
}

func keyBeneficiary(key []byte, prefix string) (beneficiary common.Address, err error) {
	k := string(key)

	split := strings.SplitAfter(k, prefix)
	if len(split) != 2 {
		return common.Address{}, errors.New("no beneficiary in key")
	}
	return common.HexToAddress(split[1]), nil
}

// LastCheque returns the last cheques for all beneficiaries.
func (s *service) LastCheques() (map[common.Address]*SignedCheque, error) {
	result := make(map[common.Address]*SignedCheque)
	err := s.store.Iterate(lastIssuedChequeKeyPrefix, func(key, val []byte) (stop bool, err error) {
		addr, err := keyBeneficiary(key, lastIssuedChequeKeyPrefix)
		if err != nil {
			return false, fmt.Errorf("parse address from key: %s: %w", string(key), err)
		}

		if _, ok := result[addr]; !ok {

			lastCheque, err := s.LastCheque(addr)
			if err != nil {
				return false, err
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

func (s *service) Withdraw(ctx context.Context, amount *big.Int) (hash common.Hash, err error) {
	availableBalance, err := s.AvailableBalance(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	// check we can afford this so we don't waste gas and don't risk bouncing cheques
	if availableBalance.Cmp(amount) < 0 {
		return common.Hash{}, ErrInsufficientFunds
	}

	callData, err := vaultABI.Pack("withdraw", amount)
	if err != nil {
		return common.Hash{}, err
	}

	request := &transaction.TxRequest{
		To:          &s.address,
		Data:        callData,
		Value:       big.NewInt(0),
		Description: fmt.Sprintf("vault withdrawal of %d WBTT", amount),
	}

	txHash, err := s.transactionService.Send(ctx, request)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (s *service) GetWithdrawTime(ctx context.Context) (*big.Int, error) {
	callData, err := vaultABI.Pack("withdrawTime")
	if err != nil {
		return nil, err
	}

	output, err := s.transactionService.Call(ctx, &transaction.TxRequest{
		To:   &s.address,
		Data: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := vaultABI.Unpack("withdrawTime", output)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, errDecodeABI
	}

	time, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || time == nil {
		return nil, errDecodeABI
	}

	return time, nil
}

func (s *service) WBTTBalanceOf(ctx context.Context, addr common.Address) (*big.Int, error) {
	return s.erc20Service.BalanceOf(ctx, addr)
}

func (s *service) BTTBalanceOf(ctx context.Context, address common.Address, block *big.Int) (*big.Int, error) {
	return s.transactionService.BttBalanceAt(ctx, address, block)
}
func (s *service) TotalPaidOut(ctx context.Context) (*big.Int, error) {
	return s.contract.TotalPaidOut(ctx)
}

// UpgradeTo will upgrade vault implementation to `newVaultImpl`
func (s *service) UpgradeTo(ctx context.Context, newVaultImpl common.Address) (old, new common.Address, err error) {
	empty := common.Address{}
	if newVaultImpl == empty {
		err = errors.New("given vault implementation address is empty")
		return
	}
	oldVaultImpl, err := GetVaultImpl(ctx, s.address, s.transactionService)
	if err != nil {
		return
	}
	if oldVaultImpl == newVaultImpl {
		err = errors.New(fmt.Sprintf("already upgraded to version %s", newVaultImpl))
		return
	}

	err = s.contract.UpgradeTo(ctx, newVaultImpl)
	if err != nil {
		return
	}
	return oldVaultImpl, newVaultImpl, nil
}

// IsVaultImplCompatibleBetween checks whether my vault's impl is compatible with peer's one.
func IsVaultImplCompatibleBetween(ctx context.Context, vault1, vault2 common.Address, trxSvc transaction.Service) (isCompatible bool, err error) {
	notFound := common.Address{}
	impl1, err := GetVaultImpl(ctx, vault1, trxSvc)
	if err != nil || impl1 == notFound {
		return
	}
	impl2, err := GetVaultImpl(ctx, vault2, trxSvc)
	if err != nil || impl2 == notFound {
		return
	}
	isCompatible = impl1 == impl2
	return
}
