package vault

import (
	"context"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"math/big"
	"strings"
	"sync"

	conabi "github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/settlement/swap/erc20"
	"github.com/bittorrent/go-btfs/statestore"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/bittorrent/go-btfs/transaction/storage"
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

	vaultABI                   = transaction.ParseABIUnchecked(conabi.VaultABI)
	vaultABINew                = transaction.ParseABIUnchecked(conabi.MutiVaultABI2)
	chequeCashedEventType      = vaultABI.Events["ChequeCashed"]
	mutiChequeCashedEventType  = vaultABINew.Events["MultiTokenChequeCashed"]
	chequeBouncedEventType     = vaultABI.Events["ChequeBounced"]
	mutiChequeBouncedEventType = vaultABINew.Events["MultiTokenChequeBounced"]
)

// Service is the main interface for interacting with the nodes vault.
type Service interface {
	// Deposit starts depositing erc20 token into the vault. This returns once the transactions has been broadcast.
	Deposit(ctx context.Context, amount *big.Int, token common.Address) (hash common.Hash, err error)
	// Withdraw starts withdrawing erc20 token from the vault. This returns once the transactions has been broadcast.
	Withdraw(ctx context.Context, amount *big.Int, token common.Address) (hash common.Hash, err error)
	// WaitForDeposit waits for the deposit transaction to confirm and verifies the result.
	WaitForDeposit(ctx context.Context, txHash common.Hash) error
	// TotalBalance returns the token balance of the vault.
	TotalBalance(ctx context.Context, token common.Address) (*big.Int, error)
	// TotalIssuedCount returns total issued count of the vault.
	TotalIssuedCount(token common.Address) (int, error)
	TotalIssued(token common.Address) (*big.Int, error)
	TotalReceivedCount(token common.Address) (int, error)
	TotalReceivedCashedCount(token common.Address) (int, error)
	TotalReceived(token common.Address) (*big.Int, error)
	TotalReceivedCashed(token common.Address) (*big.Int, error)
	TotalDailyReceived(token common.Address) (*big.Int, error)
	TotalDailyReceivedCashed(token common.Address) (*big.Int, error)
	// LiquidBalance returns the token balance of the vault sub stake amount. (not use)
	LiquidBalance(ctx context.Context) (*big.Int, error)
	// AvailableBalance returns the token balance of the vault which is not yet used for uncashed cheques.
	AvailableBalance(ctx context.Context, token common.Address) (*big.Int, error)
	// Address returns the address of the used vault contract.
	Address() common.Address
	// Issue a new cheque for the beneficiary with an cumulativePayout amount higher than the last.
	Issue(ctx context.Context, beneficiary common.Address, amount *big.Int, token common.Address, sendChequeFunc SendChequeFunc) (*big.Int, error)
	// LastCheque returns the last cheque we issued for the beneficiary.
	LastCheque(beneficiary common.Address, token common.Address) (*SignedCheque, error)
	// LastCheques returns the last cheques we issued for all beneficiaries.
	LastCheques(token common.Address) (map[common.Address]*SignedCheque, error)
	// WbttBalanceOf retrieve the addr balance
	WBTTBalanceOf(ctx context.Context, addr common.Address) (*big.Int, error)
	// TokenBalanceOf retrieve the addr balance
	TokenBalanceOf(ctx context.Context, addr common.Address, tokenStr string) (*big.Int, error)
	// BTTBalanceOf retrieve the btt balance of addr
	BTTBalanceOf(ctx context.Context, address common.Address, block *big.Int) (*big.Int, error)
	// TotalPaidOut return total pay out of the vault
	TotalPaidOut(ctx context.Context, token common.Address) (*big.Int, error)
	// CheckBalance
	CheckBalance(amount *big.Int) (err error)
	// UpgradeTo will upgrade vault implementation to `newVaultImpl`
	UpgradeTo(ctx context.Context, newVaultImpl common.Address) (old, new common.Address, err error)
}

type service struct {
	lock               sync.Mutex
	transactionService transaction.Service

	address      common.Address
	contract     *vaultContractMuti
	ownerAddress common.Address

	erc20Service   erc20.Service
	mpErc20Service map[string]erc20.Service

	store        storage.StateStorer
	chequeSigner ChequeSigner
	//totalIssuedReserved   *big.Int // replace it with mpTotalIssuedReserved
	mpTotalIssuedReserved map[string]*big.Int
	chequeStore           ChequeStore
}

// New creates a new vault service for the provided vault contract.
func New(transactionService transaction.Service, address, ownerAddress common.Address, store storage.StateStorer,
	chequeSigner ChequeSigner, erc20Service erc20.Service, mpErc20Service map[string]erc20.Service, chequeStore ChequeStore) (Service, error) {
	return &service{
		transactionService: transactionService,
		address:            address,
		contract:           newVaultContractMuti(address, transactionService),
		ownerAddress:       ownerAddress,
		erc20Service:       erc20Service,
		mpErc20Service:     mpErc20Service,
		store:              store,
		chequeSigner:       chequeSigner,
		//totalIssuedReserved:   big.NewInt(0),
		mpTotalIssuedReserved: map[string]*big.Int{},
		chequeStore:           chequeStore,
	}, nil
}

// Address returns the address of the used vault contract.
func (s *service) Address() common.Address {
	return s.address
}

// Deposit starts depositing erc20 token into the vault. This returns once the transactions has been broadcast.
func (s *service) Deposit(ctx context.Context, amount *big.Int, token common.Address) (hash common.Hash, err error) {
	//balance, err := s.erc20Service.BalanceOf(ctx, s.ownerAddress)
	//if err != nil {
	//	return common.Hash{}, err
	//}
	//
	//fmt.Println("Deposit ", balance.String(), amount.String())
	//// check we can afford this so we don't waste gas
	//if balance.Cmp(amount) < 0 {
	//	return common.Hash{}, ErrInsufficientFunds
	//}

	return s.contract.Deposit(ctx, amount, token)
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
func (s *service) TotalBalance(ctx context.Context, token common.Address) (*big.Int, error) {
	return s.contract.TotalBalance(ctx, token)
}

// LiquidBalance returns the token balance of the vault sub stake amount.(not use)
func (s *service) LiquidBalance(ctx context.Context) (*big.Int, error) {
	return s.contract.LiquidBalance(ctx)
}

// AvailableBalance returns the token balance of the vault which is not yet used for uncashed cheques.
func (s *service) AvailableBalance(ctx context.Context, token common.Address) (*big.Int, error) {
	totalIssued, err := s.totalIssued(token)
	if err != nil {
		return nil, err
	}

	balance, err := s.TotalBalance(ctx, token)
	if err != nil {
		return nil, err
	}

	totalPaidOut, err := s.contract.TotalPaidOut(ctx, token)
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
func (s *service) TotalIssuedCount(token common.Address) (int, error) {
	totalIssuedCount, err := s.totalIssuedCount(token)
	if err != nil {
		return 0, err
	}

	return totalIssuedCount, nil
}

func (s *service) TotalIssued(token common.Address) (*big.Int, error) {
	return s.totalIssued(token)
}

// total recevied cheque count.
func (s *service) TotalReceivedCount(token common.Address) (int, error) {
	return s.totalReceivedCount(token)
}

func (s *service) TotalReceivedCashedCount(token common.Address) (int, error) {
	return s.totalReceivedCashedCount(token)
}

func (s *service) TotalReceived(token common.Address) (*big.Int, error) {
	return s.totalReceived(token)
}

func (s *service) TotalReceivedCashed(token common.Address) (*big.Int, error) {
	return s.totalReceivedCashed(token)
}

func (s *service) TotalDailyReceived(token common.Address) (*big.Int, error) {
	return s.totalDailyReceived(token)
}

func (s *service) TotalDailyReceivedCashed(token common.Address) (*big.Int, error) {
	return s.totalDailyReceivedCashed(token)
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
func lastIssuedChequeKey(beneficiary common.Address, token common.Address) string {
	return fmt.Sprintf("%s%x", tokencfg.AddToken(lastIssuedChequeKeyPrefix, token), beneficiary)
}

func (s *service) reserveTotalIssued(ctx context.Context, amount *big.Int, token common.Address) (*big.Int, error) {
	availableBalance, err := s.AvailableBalance(ctx, token)
	if err != nil {
		return nil, err
	}

	tokenString := token.String()
	_, ok := s.mpTotalIssuedReserved[tokenString]
	if !ok {
		s.mpTotalIssuedReserved[tokenString] = big.NewInt(0)
	}

	if amount.Cmp(big.NewInt(0).Sub(availableBalance, s.mpTotalIssuedReserved[tokenString])) > 0 {
		return nil, ErrOutOfFunds
	}

	s.mpTotalIssuedReserved[tokenString] = s.mpTotalIssuedReserved[tokenString].Add(s.mpTotalIssuedReserved[tokenString], amount)
	return big.NewInt(0).Sub(availableBalance, amount), nil
}

func (s *service) unreserveTotalIssued(amount *big.Int, token common.Address) {
	tokenString := token.String()
	_, ok := s.mpTotalIssuedReserved[tokenString]
	if !ok {
		s.mpTotalIssuedReserved[tokenString] = big.NewInt(0)
	}
	s.mpTotalIssuedReserved[tokenString] = s.mpTotalIssuedReserved[tokenString].Sub(s.mpTotalIssuedReserved[tokenString], amount)
}

// Issue issues a new cheque and passes it to sendChequeFunc.
// The cheque is considered sent and saved when sendChequeFunc succeeds.
// The available balance which is available after sending the cheque is passed
// to the caller for it to be communicated over metrics.
func (s *service) Issue(ctx context.Context, beneficiary common.Address, amount *big.Int, token common.Address, sendChequeFunc SendChequeFunc) (*big.Int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	availableBalance, err := s.reserveTotalIssued(ctx, amount, token)
	if err != nil {
		return nil, err
	}
	defer s.unreserveTotalIssued(amount, token)

	var cumulativePayout *big.Int
	lastCheque, err := s.LastCheque(beneficiary, token)
	if err != nil {
		if err != ErrNoCheque {
			return nil, err
		}
		cumulativePayout = big.NewInt(0)
	} else {
		cumulativePayout = lastCheque.CumulativePayout
	}
	// increase cumulativePayout by amount
	cumulativePayout = cumulativePayout.Add(cumulativePayout, amount)

	// create and sign the new cheque
	cheque := Cheque{
		Token:            token,
		Vault:            s.address,
		CumulativePayout: cumulativePayout,
		Beneficiary:      beneficiary,
	}

	sig, err := s.chequeSigner.Sign(&Cheque{
		Token:            token,
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
	err = s.store.Put(lastIssuedChequeKey(beneficiary, token), cheque)
	if err != nil {
		return nil, err
	}

	// store the history issued cheque
	err = s.chequeStore.StoreSendChequeRecord(s.address, beneficiary, amount, token)
	if err != nil {
		return nil, err
	}

	// total issued count
	totalIssuedCount, err := s.totalIssuedCount(token)
	if err != nil {
		return nil, err
	}
	totalIssuedCount = totalIssuedCount + 1
	err = s.store.Put(tokencfg.AddToken(totalIssuedCountKey, token), totalIssuedCount)
	if err != nil {
		return nil, err
	}
	// totalIssued
	totalIssued, err := s.totalIssued(token)
	if err != nil {
		return nil, err
	}
	totalIssued = totalIssued.Add(totalIssued, amount)
	return availableBalance, s.store.Put(tokencfg.AddToken(totalIssuedKey, token), totalIssued)
}

// returns the total amount in cheques issued so far
func (s *service) totalIssued(token common.Address) (totalIssued *big.Int, err error) {
	err = s.store.Get(tokencfg.AddToken(totalIssuedKey, token), &totalIssued)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return big.NewInt(0), nil
	}
	return totalIssued, nil
}

// returns the total count in cheques issued so far
func (s *service) totalIssuedCount(token common.Address) (totalIssuedCount int, err error) {
	err = s.store.Get(tokencfg.AddToken(totalIssuedCountKey, token), &totalIssuedCount)
	if err != nil {
		if err != storage.ErrNotFound {
			return 0, err
		}
		return 0, nil
	}
	return totalIssuedCount, nil
}

// returns the total amount in cheques recieved so far
func (s *service) totalReceived(token common.Address) (totalReceived *big.Int, err error) {
	err = s.store.Get(tokencfg.AddToken(statestore.TotalReceivedKey, token), &totalReceived)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return big.NewInt(0), nil
	}
	return totalReceived, nil
}

func (s *service) totalReceivedCashed(token common.Address) (totalReceived *big.Int, err error) {
	err = s.store.Get(tokencfg.AddToken(statestore.TotalReceivedCashedKey, token), &totalReceived)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return big.NewInt(0), nil
	}
	return totalReceived, nil
}

// returns the total amount in cheques recieved so far
func (s *service) totalDailyReceived(token common.Address) (totalReceived *big.Int, err error) {
	var stat DailyReceivedStats
	err = s.store.Get(statestore.GetTodayTotalDailyReceivedKey(token), &stat)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return big.NewInt(0), nil
	}
	return stat.Amount, nil
}

func (s *service) totalDailyReceivedCashed(token common.Address) (totalReceived *big.Int, err error) {
	err = s.store.Get(statestore.GetTodayTotalDailyReceivedCashedKey(token), &totalReceived)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}
		return big.NewInt(0), nil
	}
	return totalReceived, nil
}

// returns the total count in cheques recieved so far
func (s *service) totalReceivedCount(token common.Address) (totalReceivedCount int, err error) {
	err = s.store.Get(tokencfg.AddToken(statestore.TotalReceivedCountKey, token), &totalReceivedCount)
	if err != nil {
		if err != storage.ErrNotFound {
			return 0, err
		}
		return 0, nil
	}
	return totalReceivedCount, nil
}

func (s *service) totalReceivedCashedCount(token common.Address) (totalReceivedCount int, err error) {
	err = s.store.Get(tokencfg.AddToken(statestore.TotalReceivedCashedCountKey, token), &totalReceivedCount)
	if err != nil {
		if err != storage.ErrNotFound {
			return 0, err
		}
		return 0, nil
	}
	return totalReceivedCount, nil
}

// LastCheque returns the last cheque we issued for the beneficiary.
func (s *service) LastCheque(beneficiary common.Address, token common.Address) (*SignedCheque, error) {
	var lastCheque *SignedCheque
	err := s.store.Get(lastIssuedChequeKey(beneficiary, token), &lastCheque)
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

// LastCheques returns the last cheques for all beneficiaries.
func (s *service) LastCheques(token common.Address) (map[common.Address]*SignedCheque, error) {
	result := make(map[common.Address]*SignedCheque)
	err := s.store.Iterate(tokencfg.AddToken(lastIssuedChequeKeyPrefix, token), func(key, val []byte) (stop bool, err error) {
		addr, err := keyBeneficiary(key, tokencfg.AddToken(lastIssuedChequeKeyPrefix, token))
		//fmt.Println("LastCheques, iterate ", addr, err, tokencfg.AddToken(lastIssuedChequeKeyPrefix, token))
		//fmt.Println("LastCheques, iterate ", key, val, string(key), string(val))

		if err != nil {
			return false, fmt.Errorf("parse address from key: %s: %w", string(key), err)
		}

		if _, ok := result[addr]; !ok {

			lastCheque, err := s.LastCheque(addr, token)
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

// OLD
//func (s *service) Withdraw(ctx context.Context, amount *big.Int, token common.Address) (hash common.Hash, err error) {
//	availableBalance, err := s.AvailableBalance(ctx, token)
//	if err != nil {
//		return common.Hash{}, err
//	}
//
//	// check we can afford this so we don't waste gas and don't risk bouncing cheques
//	if availableBalance.Cmp(amount) < 0 {
//		return common.Hash{}, ErrInsufficientFunds
//	}
//
//	callData, err := vaultABINew.Pack("withdraw", amount)
//	if err != nil {
//		return common.Hash{}, err
//	}
//
//	request := &transaction.TxRequest{
//		To:          &s.address,
//		data:        callData,
//		Value:       big.NewInt(0),
//		Description: fmt.Sprintf("vault withdrawal of %d WBTT", amount),
//	}
//
//	txHash, err := s.transactionService.Send(ctx, request)
//	if err != nil {
//		return common.Hash{}, err
//	}
//
//	return txHash, nil
//}

// Withdraw (2.3.0)
func (s *service) Withdraw(ctx context.Context, amount *big.Int, token common.Address) (hash common.Hash, err error) {
	availableBalance, err := s.AvailableBalance(ctx, token)
	if err != nil {
		return common.Hash{}, err
	}

	// check we can afford this so we don't waste gas and don't risk bouncing cheques
	if availableBalance.Cmp(amount) < 0 {
		return common.Hash{}, ErrInsufficientFunds
	}

	return s.contract.Withdraw(ctx, amount, token)
}

func (s *service) WBTTBalanceOf(ctx context.Context, addr common.Address) (*big.Int, error) {
	return s.erc20Service.BalanceOf(ctx, addr)
}

func (s *service) TokenBalanceOf(ctx context.Context, addr common.Address, tokenStr string) (*big.Int, error) {
	return s.mpErc20Service[tokenStr].BalanceOf(ctx, addr)
}

func (s *service) BTTBalanceOf(ctx context.Context, address common.Address, block *big.Int) (*big.Int, error) {
	return s.transactionService.BttBalanceAt(ctx, address, block)
}
func (s *service) TotalPaidOut(ctx context.Context, token common.Address) (*big.Int, error) {
	return s.contract.TotalPaidOut(ctx, token)
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
