package vault_test

import (
	"context"
	"math/big"

	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	"github.com/ethereum/go-ethereum/common"
)

type chequeSignerMock struct {
	sign func(cheque *vault.Cheque) ([]byte, error)
}

func (m *chequeSignerMock) Sign(cheque *vault.Cheque) ([]byte, error) {
	return m.sign(cheque)
}

type factoryMock struct {
	erc20Address   func(ctx context.Context) (common.Address, error)
	deploy         func(ctx context.Context, issuer common.Address, defaultHardDepositTimeoutDuration *big.Int, nonce common.Hash) (common.Hash, error)
	waitDeployed   func(ctx context.Context, txHash common.Hash) (common.Address, error)
	verifyBytecode func(ctx context.Context) error
	verifyVault    func(ctx context.Context, vault common.Address) error
}

// ERC20Address returns the token for which this factory deploys vaults.
func (m *factoryMock) ERC20Address(ctx context.Context) (common.Address, error) {
	return m.erc20Address(ctx)
}

func (m *factoryMock) Deploy(ctx context.Context, issuer common.Address, defaultHardDepositTimeoutDuration *big.Int, nonce common.Hash) (common.Hash, error) {
	return m.deploy(ctx, issuer, defaultHardDepositTimeoutDuration, nonce)
}

func (m *factoryMock) WaitDeployed(ctx context.Context, txHash common.Hash) (common.Address, error) {
	return m.waitDeployed(ctx, txHash)
}

// VerifyBytecode checks that the factory is valid.
func (m *factoryMock) VerifyBytecode(ctx context.Context) error {
	return m.verifyBytecode(ctx)
}

// VerifyVault checks that the supplied vault has been deployed by this factory.
func (m *factoryMock) VerifyVault(ctx context.Context, vault common.Address) error {
	return m.verifyVault(ctx, vault)
}
