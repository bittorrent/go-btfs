package vault_test

import (
	"context"

	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/peer"
)

type chequeSignerMock struct {
	sign func(cheque *vault.Cheque) ([]byte, error)
}

func (m *chequeSignerMock) Sign(cheque *vault.Cheque) ([]byte, error) {
	return m.sign(cheque)
}

type factoryMock struct {
	erc20Address             func(ctx context.Context) (common.Address, error)
	deploy                   func(ctx context.Context, issuer common.Address, vaultLogic common.Address, peerId string, tokenAddress common.Address) (vault common.Address, trx common.Hash, err error)
	waitDeployed             func(ctx context.Context, txHash common.Hash) (common.Address, error)
	verifyBytecode           func(ctx context.Context) error
	verifyVault              func(ctx context.Context, vault common.Address) error
	getPeerVault             func(ctx context.Context, peerID peer.ID) (vault common.Address, err error)
	getPeerVaultWithCache    func(ctx context.Context, peerID peer.ID) (vault common.Address, err error)
	isVaultCompatibleBetween func(ctx context.Context, peerID1, peerID2 peer.ID) (isCompatible bool, err error)
}

// ERC20Address returns the token for which this factory deploys vaults.
func (m *factoryMock) ERC20Address(ctx context.Context) (common.Address, error) {
	return m.erc20Address(ctx)
}

func (m *factoryMock) Deploy(ctx context.Context, issuer common.Address, vaultLogic common.Address, peerId string, tokenAddress common.Address) (vault common.Address, trx common.Hash, err error) {
	return m.deploy(ctx, issuer, vaultLogic, peerId, tokenAddress)
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

func (m *factoryMock) GetPeerVault(ctx context.Context, peerID peer.ID) (vault common.Address, err error) {
	return m.getPeerVault(ctx, peerID)
}

func (m *factoryMock) GetPeerVaultWithCache(ctx context.Context, peerID peer.ID) (vault common.Address, err error) {
	return m.getPeerVaultWithCache(ctx, peerID)
}

func (m *factoryMock) IsVaultCompatibleBetween(ctx context.Context, peerID1, peerID2 peer.ID) (isCompatible bool, err error) {
	return m.isVaultCompatibleBetween(ctx, peerID1, peerID2)
}
