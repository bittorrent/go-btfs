package mock

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/bittorrent/go-btfs/transaction/crypto"
	"github.com/bittorrent/go-btfs/transaction/crypto/eip712"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type signerMock struct {
	signTxFunc          func(transaction *types.Transaction, chainID *big.Int) (*types.Transaction, error)
	signTypedDataFunc   func(*eip712.TypedData) ([]byte, error)
	ethereumAddressFunc func() (common.Address, error)
	signFuncFunc        func([]byte) ([]byte, error)
	publicKeyFunc       func() (*ecdsa.PublicKey, error)
}

func (m *signerMock) EthereumAddress() (common.Address, error) {
	if m.ethereumAddressFunc != nil {
		return m.ethereumAddressFunc()
	}
	return common.Address{}, errors.New("signerMock.ethereumAddressFunc not implemented")
}

func (m *signerMock) Sign(data []byte) ([]byte, error) {
	if m.signFuncFunc != nil {
		return m.signFuncFunc(data)
	}
	return nil, errors.New("signerMock.signFuncFunc not implemented")
}

func (m *signerMock) SignTx(transaction *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	if m.signTxFunc != nil {
		return m.signTxFunc(transaction, chainID)
	}
	return nil, errors.New("signerMock.signTxFunc not implemented")
}

func (m *signerMock) PublicKey() (*ecdsa.PublicKey, error) {
	if m.publicKeyFunc != nil {
		return m.publicKeyFunc()
	}
	return nil, errors.New("signerMock.publicKeyFunc not implemented")
}

func (m *signerMock) SignTypedData(d *eip712.TypedData) ([]byte, error) {
	if m.signTypedDataFunc != nil {
		return m.signTypedDataFunc(d)
	}
	return nil, errors.New("signerMock.signTypedDataFunc not implemented")
}

func (m *signerMock) PrivKey() *ecdsa.PrivateKey {
	panic("implement me")
}

func New(opts ...Option) crypto.Signer {
	mock := new(signerMock)
	for _, o := range opts {
		o.apply(mock)
	}
	return mock
}

// Option is the option passed to the mock Chequebook service
type Option interface {
	apply(*signerMock)
}

type optionFunc func(*signerMock)

func (f optionFunc) apply(r *signerMock) { f(r) }

func WithSignFunc(f func(data []byte) ([]byte, error)) Option {
	return optionFunc(func(s *signerMock) {
		s.signFuncFunc = f
	})
}

func WithSignTxFunc(f func(transaction *types.Transaction, chainID *big.Int) (*types.Transaction, error)) Option {
	return optionFunc(func(s *signerMock) {
		s.signTxFunc = f
	})
}

func WithSignTypedDataFunc(f func(*eip712.TypedData) ([]byte, error)) Option {
	return optionFunc(func(s *signerMock) {
		s.signTypedDataFunc = f
	})
}

func WithEthereumAddressFunc(f func() (common.Address, error)) Option {
	return optionFunc(func(s *signerMock) {
		s.ethereumAddressFunc = f
	})
}

func WithPublicKeyFunc(f func() (*ecdsa.PublicKey, error)) Option {
	return optionFunc(func(s *signerMock) {
		s.publicKeyFunc = f
	})
}
