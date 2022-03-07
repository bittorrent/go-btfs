package mock

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/bittorrent/go-btfs/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type transactionServiceMock struct {
	senderAddress        func(ctx context.Context) (addr common.Address)
	ethereumAddress      func(ctx context.Context) (addr common.Address)
	overlayEthAddress    func(ctx context.Context) (addr common.Address)
	send                 func(ctx context.Context, request *transaction.TxRequest) (txHash common.Hash, err error)
	call                 func(ctx context.Context, request *transaction.TxRequest) (result []byte, err error)
	waitForReceipt       func(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error)
	watchSentTransaction func(txHash common.Hash) (chan types.Receipt, chan error, error)
	storedTransaction    func(txHash common.Hash) (*transaction.StoredTransaction, error)
	pendingTransactions  func() ([]common.Hash, error)
	resendTransaction    func(ctx context.Context, txHash common.Hash) error
	cancelTransaction    func(ctx context.Context, originalTxHash common.Hash) (common.Hash, error)
	bttBalanceAt         func(ctx context.Context, address common.Address, block *big.Int) (*big.Int, error)
	myBttBalance         func(ctx context.Context) (*big.Int, error)
}

func (m *transactionServiceMock) SenderAddress(ctx context.Context) (addr common.Address) {
	if m.senderAddress != nil {
		return m.SenderAddress(ctx)
	}
	return common.Address{}
}

func (m *transactionServiceMock) EthereumAddress(ctx context.Context) (addr common.Address) {
	if m.ethereumAddress != nil {
		return m.ethereumAddress(ctx)
	}
	return common.Address{}
}

func (m *transactionServiceMock) OverlayEthAddress(ctx context.Context) (addr common.Address) {
	if m.overlayEthAddress != nil {
		return m.overlayEthAddress(ctx)
	}
	return common.Address{}
}

func (m *transactionServiceMock) Send(ctx context.Context, request *transaction.TxRequest) (txHash common.Hash, err error) {
	if m.send != nil {
		return m.send(ctx, request)
	}
	return common.Hash{}, errors.New("not implemented")
}

func (m *transactionServiceMock) WaitForReceipt(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error) {
	if m.waitForReceipt != nil {
		return m.waitForReceipt(ctx, txHash)
	}
	return nil, errors.New("not implemented")
}

func (m *transactionServiceMock) WatchSentTransaction(txHash common.Hash) (<-chan types.Receipt, <-chan error, error) {
	if m.watchSentTransaction != nil {
		return m.watchSentTransaction(txHash)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *transactionServiceMock) Call(ctx context.Context, request *transaction.TxRequest) (result []byte, err error) {
	if m.call != nil {
		return m.call(ctx, request)
	}
	return nil, errors.New("not implemented")
}

func (m *transactionServiceMock) PendingTransactions() ([]common.Hash, error) {
	if m.pendingTransactions != nil {
		return m.pendingTransactions()
	}
	return nil, errors.New("not implemented")
}

func (m *transactionServiceMock) ResendTransaction(ctx context.Context, txHash common.Hash) error {
	if m.resendTransaction != nil {
		return m.resendTransaction(ctx, txHash)
	}
	return errors.New("not implemented")
}

func (m *transactionServiceMock) StoredTransaction(txHash common.Hash) (*transaction.StoredTransaction, error) {
	if m.storedTransaction != nil {
		return m.storedTransaction(txHash)
	}
	return nil, errors.New("not implemented")
}

func (m *transactionServiceMock) CancelTransaction(ctx context.Context, originalTxHash common.Hash) (common.Hash, error) {
	if m.cancelTransaction != nil {
		return m.cancelTransaction(ctx, originalTxHash)
	}
	return common.Hash{}, errors.New("not implemented")
}

func (m *transactionServiceMock) Close() error {
	return nil
}

func (m *transactionServiceMock) BttBalanceAt(ctx context.Context, address common.Address, block *big.Int) (*big.Int, error) {
	if m.bttBalanceAt != nil {
		return m.bttBalanceAt(ctx, address, block)
	}
	return big.NewInt(0), errors.New("Error")
}

func (m *transactionServiceMock) MyBttBalance(ctx context.Context) (*big.Int, error) {
	if m.myBttBalance != nil {
		return m.myBttBalance(ctx)
	}
	return big.NewInt(0), errors.New("Error")
}

// Option is the option passed to the mock Chequebook service
type Option interface {
	apply(*transactionServiceMock)
}

type optionFunc func(*transactionServiceMock)

func (f optionFunc) apply(r *transactionServiceMock) { f(r) }

func WithSenderAddress(f func(ctx context.Context) (addr common.Address)) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.senderAddress = f
	})
}

func WithEthereumAddress(f func(ctx context.Context) (addr common.Address)) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.ethereumAddress = f
	})
}

func WithOverlayEthAddress(f func(ctx context.Context) (addr common.Address)) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.overlayEthAddress = f
	})
}

func WithSendFunc(f func(ctx context.Context, request *transaction.TxRequest) (txHash common.Hash, err error)) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.send = f
	})
}

func WithWaitForReceiptFunc(f func(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error)) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.waitForReceipt = f
	})
}

func WithCallFunc(f func(ctx context.Context, request *transaction.TxRequest) (result []byte, err error)) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.call = f
	})
}

func WithStoredTransactionFunc(f func(txHash common.Hash) (*transaction.StoredTransaction, error)) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.storedTransaction = f
	})
}

func WithPendingTransactionsFunc(f func() ([]common.Hash, error)) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.pendingTransactions = f
	})
}

func WithResendTransactionFunc(f func(ctx context.Context, txHash common.Hash) error) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.resendTransaction = f
	})
}

func WithCancelTransactionFunc(f func(ctx context.Context, originalTxHash common.Hash) (common.Hash, error)) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.cancelTransaction = f
	})
}

func WithBttBalanceAt(f func(ctx context.Context, address common.Address, block *big.Int) (*big.Int, error)) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.bttBalanceAt = f
	})
}

func WithMyBttBalance(f func(ctx context.Context) (*big.Int, error)) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.myBttBalance = f
	})
}

func New(opts ...Option) transaction.Service {
	mock := new(transactionServiceMock)
	for _, o := range opts {
		o.apply(mock)
	}
	return mock
}

type Call struct {
	abi    *abi.ABI
	to     common.Address
	result []byte
	method string
	params []interface{}
}

func ABICall(abi *abi.ABI, to common.Address, result []byte, method string, params ...interface{}) Call {
	return Call{
		to:     to,
		abi:    abi,
		result: result,
		method: method,
		params: params,
	}
}

func WithABICallSequence(calls ...Call) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.call = func(ctx context.Context, request *transaction.TxRequest) ([]byte, error) {
			if len(calls) == 0 {
				return nil, errors.New("unexpected call")
			}

			call := calls[0]

			data, err := call.abi.Pack(call.method, call.params...)
			if err != nil {
				return nil, err
			}

			if !bytes.Equal(data, request.Data) {
				return nil, fmt.Errorf("wrong data. wanted %x, got %x", data, request.Data)
			}

			if request.To == nil {
				return nil, errors.New("call with no recipient")
			}
			if *request.To != call.to {
				return nil, fmt.Errorf("wrong recipient. wanted %x, got %x", call.to, *request.To)
			}

			calls = calls[1:]

			return call.result, nil
		}
	})
}

func WithABICall(abi *abi.ABI, to common.Address, result []byte, method string, params ...interface{}) Option {
	return WithABICallSequence(ABICall(abi, to, result, method, params...))
}

func WithABISend(abi *abi.ABI, txHash common.Hash, expectedAddress common.Address, expectedValue *big.Int, method string, params ...interface{}) Option {
	return optionFunc(func(s *transactionServiceMock) {
		s.send = func(ctx context.Context, request *transaction.TxRequest) (common.Hash, error) {
			data, err := abi.Pack(method, params...)
			if err != nil {
				return common.Hash{}, err
			}

			if !bytes.Equal(data, request.Data) {
				return common.Hash{}, fmt.Errorf("wrong data. wanted %x, got %x", data, request.Data)
			}

			if request.To != nil && *request.To != expectedAddress {
				return common.Hash{}, fmt.Errorf("sending to wrong contract. wanted %x, got %x", expectedAddress, request.To)
			}
			if request.Value.Cmp(expectedValue) != 0 {
				return common.Hash{}, fmt.Errorf("sending with wrong value. wanted %d, got %d", expectedValue, request.Value)
			}

			return txHash, nil
		}
	})
}
