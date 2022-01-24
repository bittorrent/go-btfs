package vault

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/bittorrent/go-btfs/transaction/crypto"
	"github.com/bittorrent/go-btfs/transaction/crypto/eip712"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

// Cheque represents a cheque for a SimpleSwap vault
type Cheque struct {
	Vault            common.Address
	Beneficiary      common.Address
	CumulativePayout *big.Int
}

type ChequeRecord struct {
	Vault       common.Address
	Beneficiary common.Address
	Amount      *big.Int
	ReceiveTime int64 //time.now().Unix()
}

type DailyReceivedStats struct {
	Amount *big.Int
	Count  int
	Date   int64 //time.now().Unix()
}

type DailySentStats struct {
	Amount *big.Int
	Count  int
	Date   int64 //time.now().Unix()
}

// the valid range is [MinIndex, MaxIndex)
type IndexRange struct {
	MinIndex uint64
	MaxIndex uint64
}

// SignedCheque represents a cheque together with its signature
type SignedCheque struct {
	Cheque
	Signature []byte
}

// vaultDomain computes chainId-dependant EIP712 domain
func vaultDomain(chainID int64) eip712.TypedDataDomain {
	return eip712.TypedDataDomain{
		Name:    "Vault",
		Version: "1.0",
		ChainId: math.NewHexOrDecimal256(chainID),
	}
}

// ChequeTypes are the needed type descriptions for cheque signing
var ChequeTypes = eip712.Types{
	"EIP712Domain": eip712.EIP712DomainType,
	"Cheque": []eip712.Type{
		{
			Name: "vault",
			Type: "address",
		},
		{
			Name: "beneficiary",
			Type: "address",
		},
		{
			Name: "cumulativePayout",
			Type: "uint256",
		},
	},
}

// ChequeSigner signs cheque
type ChequeSigner interface {
	// Sign signs a cheque
	Sign(cheque *Cheque) ([]byte, error)
}

type chequeSigner struct {
	signer  crypto.Signer // the underlying signer used
	chainID int64         // the chainID used for EIP712
}

// NewChequeSigner creates a new cheque signer for the given chainID.
func NewChequeSigner(signer crypto.Signer, chainID int64) ChequeSigner {
	return &chequeSigner{
		signer:  signer,
		chainID: chainID,
	}
}

// eip712DataForCheque converts a cheque into the correct TypedData structure.
func eip712DataForCheque(cheque *Cheque, chainID int64) *eip712.TypedData {
	return &eip712.TypedData{
		Domain: vaultDomain(chainID),
		Types:  ChequeTypes,
		Message: eip712.TypedDataMessage{
			"vault":            cheque.Vault.Hex(),
			"beneficiary":      cheque.Beneficiary.Hex(),
			"cumulativePayout": cheque.CumulativePayout.String(),
		},
		PrimaryType: "Cheque",
	}
}

// Sign signs a cheque.
func (s *chequeSigner) Sign(cheque *Cheque) ([]byte, error) {
	return s.signer.SignTypedData(eip712DataForCheque(cheque, s.chainID))
}

func (cheque *Cheque) String() string {
	return fmt.Sprintf("Contract: %x Beneficiary: %x CumulativePayout: %v", cheque.Vault, cheque.Beneficiary, cheque.CumulativePayout)
}

func (cheque *Cheque) Equal(other *Cheque) bool {
	if cheque.Beneficiary != other.Beneficiary {
		return false
	}
	if cheque.CumulativePayout.Cmp(other.CumulativePayout) != 0 {
		return false
	}
	return cheque.Vault == other.Vault
}

func (cheque *SignedCheque) Equal(other *SignedCheque) bool {
	if !bytes.Equal(cheque.Signature, other.Signature) {
		return false
	}
	return cheque.Cheque.Equal(&other.Cheque)
}
