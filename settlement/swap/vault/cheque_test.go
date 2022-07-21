package vault_test

import (
	"bytes"

	"encoding/hex"
	"math/big"
	"testing"

	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	"github.com/bittorrent/go-btfs/transaction/crypto"
	"github.com/bittorrent/go-btfs/transaction/crypto/eip712"
	signermock "github.com/bittorrent/go-btfs/transaction/crypto/mock"
	"github.com/ethereum/go-ethereum/common"
)

func TestSignCheque(t *testing.T) {
	vaultAddress := common.HexToAddress("0x8d3766440f0d7b949a5e32995d09619a7f86e632")
	beneficiaryAddress := common.HexToAddress("0xb8d424e9662fe0837fb1d728f1ac97cebb1085fe")
	signature := common.Hex2Bytes("abcd")
	cumulativePayout := big.NewInt(10)
	chainId := int64(1)
	cheque := &vault.Cheque{
		Vault:            vaultAddress,
		Beneficiary:      beneficiaryAddress,
		CumulativePayout: cumulativePayout,
	}

	signer := signermock.New(
		signermock.WithSignTypedDataFunc(func(data *eip712.TypedData) ([]byte, error) {

			if data.Message["beneficiary"].(string) != beneficiaryAddress.Hex() {
				t.Fatal("signing cheque with wrong beneficiary")
			}

			if data.Message["vault"].(string) != vaultAddress.Hex() {
				t.Fatal("signing cheque for wrong vault")
			}

			if data.Message["cumulativePayout"].(string) != cumulativePayout.String() {
				t.Fatal("signing cheque with wrong cumulativePayout")
			}

			return signature, nil
		}),
	)

	chequeSigner := vault.NewChequeSigner(signer, chainId)

	result, err := chequeSigner.Sign(cheque)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, signature) {
		t.Fatalf("returned wrong signature. wanted %x, got %x", signature, result)
	}
}

func TestSignChequeIntegration(t *testing.T) {
	vaultAddress := common.HexToAddress("0xfa02D396842E6e1D319E8E3D4D870338F791AA25")
	beneficiaryAddress := common.HexToAddress("0x98E6C644aFeB94BBfB9FF60EB26fc9D83BBEcA79")
	cumulativePayout := big.NewInt(500)
	chainId := int64(1)

	data, err := hex.DecodeString("634fb5a872396d9693e5c9f9d7233cfa93f395c093371017ff44aa9ae6564cdd")
	if err != nil {
		t.Fatal(err)
	}

	privKey, err := crypto.DecodeSecp256k1PrivateKey(data)
	if err != nil {
		t.Fatal(err)
	}

	signer := crypto.NewDefaultSigner(privKey)

	cheque := &vault.Cheque{
		Vault:            vaultAddress,
		Beneficiary:      beneficiaryAddress,
		CumulativePayout: cumulativePayout,
	}

	chequeSigner := vault.NewChequeSigner(signer, chainId)

	result, err := chequeSigner.Sign(cheque)
	if err != nil {
		t.Fatal(err)
	}

	// computed using ganache
	expectedSignature, err := hex.DecodeString("3305964770f9b66463d58b61e7de9bf2b784098fc715338083946aabf69c7dec0cae8345705d4bf2e556482bd27f11ceb0e5ae1eb191a907cde7aee3989a3ad51c")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, expectedSignature) {
		t.Fatalf("returned wrong signature. wanted %x, got %x", expectedSignature, result)
	}
}
