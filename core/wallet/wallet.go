package wallet

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/escrow"

	config "github.com/TRON-US/go-btfs-config"
	"github.com/tron-us/go-btfs-common/crypto"
	"github.com/tron-us/go-btfs-common/ledger"
	escrowpb "github.com/tron-us/go-btfs-common/protos/escrow"
	ledgerpb "github.com/tron-us/go-btfs-common/protos/ledger"
	"github.com/tron-us/go-btfs-common/utils/grpc"
	"github.com/tron-us/protobuf/proto"

	logging "github.com/ipfs/go-log"
	ic "github.com/libp2p/go-libp2p/core/crypto"
)

var log = logging.Logger("core/wallet")

var (
	WithdrawMinAmount int64 = 1
	WithdrawMaxAmount int64 = 1000000000000
	DepositMinAmount  int64 = 1
	DepositMaxAmount  int64 = 1000000000000
	TokenId                 = "1002000"
	TokenIdDev              = "1000252"
	hostWallet        Wallet

	escrowService   string
	exchangeService string
	solidityService string
)

type Wallet struct {
	privKeyIC     ic.PrivKey
	privateKey    *ecdsa.PrivateKey
	tronAddress   []byte // 41***
	ledgerAddress []byte // address in ledger
}

// withdraw from ledger to tron
func WalletWithdraw(ctx context.Context, configuration *config.Config, n *core.IpfsNode, amount int64) error {
	err := Init(ctx, configuration)
	if err != nil {
		return err
	}

	if hostWallet.privateKey == nil {
		log.Error("wallet is not initialized")
		return errors.New("wallet is not initialized")
	}

	if amount < WithdrawMinAmount || amount > WithdrawMaxAmount {
		return errors.New(fmt.Sprintf("withdraw amount should between %d ~ %d", WithdrawMinAmount, WithdrawMaxAmount))
	}

	// get ledger balance before withdraw
	ledgerBalance, err := Balance(ctx, configuration)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get ledger balance, reason: %v", err))
	}
	log.Info(fmt.Sprintf("Get ledger account success, balance: [%d]", ledgerBalance))

	if amount > ledgerBalance {
		return errors.New(fmt.Sprintf("not enough ledger balance, current balance is %d", ledgerBalance))
	}

	// Doing withdraw request.
	channelId, id, err := Withdraw(ctx, n, hostWallet.ledgerAddress, hostWallet.tronAddress, amount, hostWallet.privateKey)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Withdraw submitted! ChannelId: [%d], id [%d]\n", channelId, id))
	return nil
}

const (
	InAppWallet = "BTFS Wallet"
	BttWallet   = "BTT Wallet"

	StatusPending = "Pending"
	StatusSuccess = "Success"
	StatusFailed  = "Failed"
)

func WalletDeposit(ctx context.Context, configuration *config.Config, n *core.IpfsNode,
	amount int64, runDaemon bool, async bool) error {
	err := Init(ctx, configuration)
	if err != nil {
		return err
	}

	if hostWallet.privateKey == nil {
		log.Error("wallet is not initialized")
		return errors.New("wallet is not initialized")
	}

	if amount < DepositMinAmount || amount > DepositMaxAmount {
		return errors.New(fmt.Sprintf("deposit amount should between %d ~ %d", DepositMinAmount, DepositMaxAmount))
	}

	_, err = Balance(ctx, configuration)
	if err != nil {
		return err
	}

	prepareResponse, err := Deposit(ctx, n, hostWallet.ledgerAddress, amount, hostWallet.privateKey, runDaemon, async)
	if err != nil {
		log.Error("Failed to Deposit, ERR[%v]\n", err)
		return err
	}

	fmt.Println(fmt.Sprintf("Deposit Submitted: Id [%d]\n", prepareResponse.GetId()))
	return nil
}

// GetBalance both on ledger and Tron.
func GetBalance(ctx context.Context, configuration *config.Config) (int64, int64, error) {
	err := Init(ctx, configuration)
	if err != nil {
		return 0, 0, err
	}

	if hostWallet.privateKey == nil {
		log.Error("wallet is not initialized")
		return 0, 0, errors.New("wallet is not initialized")
	}

	// get tron balance
	tokenId := TokenId
	if strings.Contains(configuration.Services.EscrowDomain, "dev") ||
		strings.Contains(configuration.Services.EscrowDomain, "staging") {
		tokenId = TokenIdDev
	}

	tronBalance, err := GetTokenBalance(ctx, hostWallet.tronAddress, tokenId)
	if err != nil {
		return 0, 0,
			errors.New(fmt.Sprintf("Failed to get exchange tron balance, reason: %v", err))
	}

	log.Info(fmt.Sprintf("Get exchange tron account success, balance: [%d]", tronBalance))

	// get ledger balance from escrow
	ledgerBalance, err := Balance(ctx, configuration)
	if err != nil {
		return 0, 0,
			errors.New(fmt.Sprintf("Failed to get ledger balance, reason: %v", err))
	}

	log.Info(fmt.Sprintf("Get ledger account success, balance: [%d]", ledgerBalance))

	return tronBalance, ledgerBalance, nil
}

func Init(ctx context.Context, configuration *config.Config) error {
	if configuration == nil {
		fmt.Println("Init wallet, configuration is nil")
		log.Error("init wallet failed, input nil configuration")
		return errors.New("init wallet failed")
	}

	// get service name
	escrowService = configuration.Services.EscrowDomain
	exchangeService = configuration.Services.ExchangeDomain
	solidityService = configuration.Services.SolidityDomain

	// get key
	privKeyIC, err := configuration.Identity.DecodePrivateKey("")
	if err != nil {
		log.Error("wallet get private key failed")
		return err
	}
	// base64 key
	privKeyRaw, err := privKeyIC.Raw()
	if err != nil {
		log.Error("wallet get private key raw failed")
		return err
	}
	// hex key
	hexPrivKey := hex.EncodeToString(privKeyRaw)
	// hex key to ecdsa
	privateKey, err := crypto.HexToECDSA(hexPrivKey)
	if err != nil {
		log.Error("error when convent private key to edca, ERR[%v]\n", err)
		return err
	}
	if privateKey == nil {
		log.Error("wallet get private key ecdsa failed")
		return err
	}
	hostWallet.privateKey = privateKey

	// tron key 41****
	addr, err := crypto.PublicKeyToAddress(privateKey.PublicKey)
	if err != nil {
		log.Error("wallet get tron address failed, ERR[%v]\n ", err)
		return err
	}
	addBytes := addr.Bytes()
	hostWallet.tronAddress = addBytes

	ledgerAddress, err := ic.MarshalPublicKey(privKeyIC.GetPublic())
	if err != nil {
		fmt.Println("get ledger address failed, ERR: \n", err)
		return err
	}

	hostWallet.ledgerAddress = ledgerAddress
	return nil
}

func Balance(ctx context.Context, configuration *config.Config) (int64, error) {
	privKey, err := configuration.Identity.DecodePrivateKey("")
	if err != nil {
		return 0, err
	}
	lgSignedPubKey, err := ledger.NewSignedPublicKey(privKey, privKey.GetPublic())
	if err != nil {
		return 0, err
	}

	return BalanceHelper(ctx, configuration, false, nil, lgSignedPubKey)
}

func BalanceHelper(ctx context.Context, configuration *config.Config, offsign bool, signedBytes []byte, lgSignedPubKey *ledgerpb.SignedPublicKey) (int64, error) {
	if offsign {
		var ledgerSignedPubKey ledgerpb.SignedPublicKey
		err := proto.Unmarshal(signedBytes, &ledgerSignedPubKey)
		if err != nil {
			return 0, err
		}
		lgSignedPubKey = &ledgerSignedPubKey
	}

	var balance int64 = 0
	err := grpc.EscrowClient(configuration.Services.EscrowDomain).WithContext(ctx,
		func(ctx context.Context, client escrowpb.EscrowServiceClient) error {
			res, err := client.BalanceOf(ctx, ledger.NewSignedCreateAccountRequest(lgSignedPubKey.Key, lgSignedPubKey.Signature))
			if err != nil {
				return err
			}
			err = escrow.VerifyEscrowRes(configuration, res.Result, res.EscrowSignature)
			if err != nil {
				return err
			}
			balance = res.Result.Balance
			log.Debug("balance of account is ", balance)
			return nil
		})
	if err != nil {
		return 0, err
	}
	return balance, nil
}
