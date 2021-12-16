package swap

import (
	"errors"
	"fmt"

	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/ethereum/go-ethereum/common"
)

var (
	peerPrefix            = "swap_vault_peer_"
	peerVaultPrefix       = "swap_peer_vault_"
	beneficiaryPeerPrefix = "swap_beneficiary_peer_"
	peerBeneficiaryPrefix = "swap_peer_beneficiary_"
	deductedForPeerPrefix = "swap_deducted_for_peer_"
	deductedByPeerPrefix  = "swap_deducted_by_peer_"
)

// Addressbook maps peers to beneficaries, vaults and in reverse.
type Addressbook interface {
	// Beneficiary returns the beneficiary for the given peer.
	Beneficiary(peer string) (beneficiary common.Address, known bool, err error)
	// Vault returns the vault for the given peer.
	Vault(peer string) (vaultAddress common.Address, known bool, err error)
	// BeneficiaryPeer returns the peer for a beneficiary.
	BeneficiaryPeer(beneficiary common.Address) (peer string, known bool, err error)
	// VaultPeer returns the peer for a vault.
	VaultPeer(vault common.Address) (peer string, known bool, err error)
	// PutBeneficiary stores the beneficiary for the given peer.
	PutBeneficiary(peer string, beneficiary common.Address) error
	// PutVault stores the vault for the given peer.
	PutVault(peer string, vault common.Address) error
	// MigratePeer returns whether a peer have already received a cheque that has been deducted
	MigratePeer(oldPeer, newPeer string) error
}

type addressbook struct {
	store storage.StateStorer
}

// NewAddressbook creates a new addressbook using the store.
func NewAddressbook(store storage.StateStorer) Addressbook {
	return &addressbook{
		store: store,
	}
}

func (a *addressbook) MigratePeer(oldPeer, newPeer string) error {
	ba, known, err := a.Beneficiary(oldPeer)
	if err != nil {
		return err
	}
	if !known {
		return errors.New("old beneficiary not known")
	}

	cb, known, err := a.Vault(oldPeer)
	if err != nil {
		return err
	}

	if err := a.PutBeneficiary(newPeer, ba); err != nil {
		return err
	}

	if err := a.store.Delete(peerBeneficiaryKey(oldPeer)); err != nil {
		return err
	}

	if known {
		if err := a.PutVault(newPeer, cb); err != nil {
			return err
		}
		if err := a.store.Delete(peerKey(oldPeer)); err != nil {
			return err
		}
	}

	return nil
}

// Beneficiary returns the beneficiary for the given peer.
func (a *addressbook) Beneficiary(peer string) (beneficiary common.Address, known bool, err error) {
	err = a.store.Get(peerBeneficiaryKey(peer), &beneficiary)
	if err != nil {
		if err != storage.ErrNotFound {
			return common.Address{}, false, err
		}
		return common.Address{}, false, nil
	}
	return beneficiary, true, nil
}

// BeneficiaryPeer returns the peer for a beneficiary.
func (a *addressbook) BeneficiaryPeer(beneficiary common.Address) (peer string, known bool, err error) {
	err = a.store.Get(beneficiaryPeerKey(beneficiary), &peer)
	if err != nil {
		if err != storage.ErrNotFound {
			return "", false, err
		}
		return "", false, nil
	}
	return peer, true, nil
}

// Vault returns the vault for the given peer.
func (a *addressbook) Vault(peer string) (vaultAddress common.Address, known bool, err error) {
	err = a.store.Get(peerKey(peer), &vaultAddress)
	if err != nil {
		if err != storage.ErrNotFound {
			return common.Address{}, false, err
		}
		return common.Address{}, false, nil
	}
	return vaultAddress, true, nil
}

// VaultPeer returns the peer for a vault.
func (a *addressbook) VaultPeer(vault common.Address) (peer string, known bool, err error) {
	err = a.store.Get(vaultPeerKey(vault), &peer)
	if err != nil {
		if err != storage.ErrNotFound {
			return "", false, err
		}
		return "", false, nil
	}
	return peer, true, nil
}

// PutBeneficiary stores the beneficiary for the given peer.
func (a *addressbook) PutBeneficiary(peer string, beneficiary common.Address) error {
	err := a.store.Put(peerBeneficiaryKey(peer), beneficiary)
	if err != nil {
		return err
	}
	return a.store.Put(beneficiaryPeerKey(beneficiary), peer)
}

// PutVault stores the vault for the given peer.
func (a *addressbook) PutVault(peer string, vault common.Address) error {
	err := a.store.Put(peerKey(peer), vault)
	if err != nil {
		return err
	}
	return a.store.Put(vaultPeerKey(vault), peer)
}

// peerKey computes the key where to store the vault from a peer.
func peerKey(peer string) string {
	return fmt.Sprintf("%s%s", peerPrefix, peer)
}

// vaultPeerKey computes the key where to store the peer for a vault.
func vaultPeerKey(vault common.Address) string {
	return fmt.Sprintf("%s%x", peerVaultPrefix, vault)
}

// peerBeneficiaryKey computes the key where to store the beneficiary for a peer.
func peerBeneficiaryKey(peer string) string {
	return fmt.Sprintf("%s%s", peerBeneficiaryPrefix, peer)
}

// beneficiaryPeerKey computes the key where to store the peer for a beneficiary.
func beneficiaryPeerKey(peer common.Address) string {
	return fmt.Sprintf("%s%x", beneficiaryPeerPrefix, peer)
}

func peerDeductedByKey(peer string) string {
	return fmt.Sprintf("%s%s", deductedByPeerPrefix, peer)
}

func peerDeductedForKey(peer string) string {
	return fmt.Sprintf("%s%s", deductedForPeerPrefix, peer)
}
