package chain

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/chain/tokencfg"

	"github.com/bittorrent/go-btfs/accounting"
	"github.com/bittorrent/go-btfs/chain/config"
	"github.com/bittorrent/go-btfs/settlement"
	"github.com/bittorrent/go-btfs/settlement/swap"
	"github.com/bittorrent/go-btfs/settlement/swap/bttc"
	"github.com/bittorrent/go-btfs/settlement/swap/erc20"
	"github.com/bittorrent/go-btfs/settlement/swap/priceoracle"
	"github.com/bittorrent/go-btfs/settlement/swap/swapprotocol"
	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/bittorrent/go-btfs/transaction/crypto"
	"github.com/bittorrent/go-btfs/transaction/sctx"
	"github.com/bittorrent/go-btfs/transaction/storage"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	logging "github.com/ipfs/go-log"
)

var (
	log          = logging.Logger("chain")
	ChainObject  ChainInfo
	SettleObject SettleInfo

	StateStore storage.StateStorer
)

const (
	MaxDelay          = 1 * time.Minute
	CancellationDepth = 6
)

type ChainInfo struct {
	Chainconfig        config.ChainConfig
	Backend            transaction.Backend
	OverlayAddress     common.Address
	Signer             crypto.Signer
	ChainID            int64
	PeerID             string
	TransactionMonitor transaction.Monitor
	TransactionService transaction.Service
}

type SettleInfo struct {
	Factory        vault.Factory
	VaultService   vault.Service
	ChequeStore    vault.ChequeStore
	CashoutService vault.CashoutService
	SwapService    *swap.Service
	OracleService  priceoracle.Service
	BttcService    bttc.Service
}

// InitChain will initialize the Ethereum backend at the given endpoint and
// set up the Transaction Service to interact with it using the provided signer.
func InitChain(
	ctx context.Context,
	stateStore storage.StateStorer,
	signer crypto.Signer,
	pollingInterval time.Duration,
	chainID int64,
	peerid string,
	chainconfig *config.ChainConfig,
) (*ChainInfo, error) {

	StateStore = stateStore

	backend, err := ethclient.Dial(chainconfig.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("dial eth client: %w", err)
	}

	_, err = backend.BlockNumber(context.Background())
	if err != nil {
		errMsg := "Could not connect to blockchain rpc, please check your network connection"
		if err == io.EOF {
			return nil, errors.New(errMsg)

		}
		return nil, fmt.Errorf("%s.%w", errMsg, err)
	}

	overlayEthAddress, err := signer.EthereumAddress()
	if err != nil {
		return nil, fmt.Errorf("eth address: %w", err)
	}

	transactionMonitor := transaction.NewMonitor(backend, overlayEthAddress, pollingInterval, CancellationDepth)

	transactionService, err := transaction.NewService(backend, signer, stateStore, big.NewInt(chainID), transactionMonitor)
	if err != nil {
		return nil, fmt.Errorf("new transaction service: %w", err)
	}

	ChainObject = ChainInfo{
		Chainconfig:        *chainconfig,
		Backend:            backend,
		OverlayAddress:     overlayEthAddress,
		ChainID:            chainID,
		PeerID:             peerid,
		Signer:             signer,
		TransactionMonitor: transactionMonitor,
		TransactionService: transactionService,
	}

	return &ChainObject, nil
}

func InitSettlement(
	ctx context.Context,
	stateStore storage.StateStorer,
	chaininfo *ChainInfo,
	deployGasPrice string,
	chainID int64,
) (*SettleInfo, error) {
	//InitVaultFactory
	factory, err := initVaultFactory(chaininfo.Backend, chaininfo.ChainID, chaininfo.TransactionService,
		chaininfo.Chainconfig.CurrentFactory.String())

	if err != nil {
		return nil, errors.New("init vault factory error")
	}

	// init wbtt service
	erc20Address, err := factory.ERC20Address(ctx)
	if err != nil {
		return nil, err
	}
	erc20Service := erc20.New(chaininfo.Backend, chaininfo.TransactionService, erc20Address)

	// muti tokens
	mpErc20Service := make(map[string]erc20.Service)
	for k, tokenAddr := range tokencfg.MpTokenAddr {
		mpErc20Service[k] = erc20.New(chaininfo.Backend, chaininfo.TransactionService, tokenAddr)
	}

	// init bttc service
	bttcService := bttc.New(chaininfo.TransactionService, erc20Service, mpErc20Service)

	//initChequeStoreCashout
	chequeStore, cashoutService := initChequeStoreCashout(
		stateStore,
		chaininfo.Backend,
		factory,
		chainID,
		chaininfo.OverlayAddress,
		chaininfo.TransactionService,
	)

	//new accounting
	accounting, err := accounting.NewAccounting(stateStore)

	if err != nil {
		return nil, errors.New("new accounting service error")
	}

	//InitVaultService
	vaultService, err := initVaultService(
		ctx,
		stateStore,
		chaininfo.Signer,
		chaininfo.ChainID,
		chaininfo.PeerID,
		chaininfo.Chainconfig.VaultLogicAddress,
		chaininfo.Backend,
		chaininfo.OverlayAddress,
		chaininfo.TransactionService,
		factory,
		deployGasPrice,
		chequeStore,
		erc20Service,
		mpErc20Service,
	)

	if err != nil {
		return nil, fmt.Errorf("init vault service: %w", err)
	}

	//InitSwap
	swapService, priceOracleService, err := initSwap(
		stateStore,
		chaininfo.OverlayAddress,
		vaultService,
		chequeStore,
		cashoutService,
		accounting,
		chaininfo.Chainconfig.PriceOracleAddress.String(),
		chaininfo.ChainID,
		chaininfo.TransactionService,
	)

	if err != nil {
		return nil, errors.New("init swap service error , " + err.Error())
	}

	accounting.SetPayFunc(swapService.Pay)

	SettleObject = SettleInfo{
		Factory:        factory,
		VaultService:   vaultService,
		ChequeStore:    chequeStore,
		CashoutService: cashoutService,
		SwapService:    swapService,
		OracleService:  priceOracleService,
		BttcService:    bttcService,
	}

	return &SettleObject, nil
}

// InitVaultFactory will initialize the vault factory with the given
// chain backend.
func initVaultFactory(
	backend transaction.Backend,
	chainID int64,
	transactionService transaction.Service,
	factoryAddress string,
) (vault.Factory, error) {
	var currentFactory common.Address

	if factoryAddress == "" {
		log.Infof("none factory address for chain id %d: %x", chainID, currentFactory)
		return nil, errors.New("none factory address for chain id")
	} else if !common.IsHexAddress(factoryAddress) {
		return nil, errors.New("malformed factory address")
	} else {
		currentFactory = common.HexToAddress(factoryAddress)
		log.Infof("using custom factory address: %x", currentFactory)
	}

	return vault.NewFactory(
		backend,
		transactionService,
		currentFactory,
	), nil
}

// InitVaultService will initialize the vault service with the given
// vault factory and chain backend.
func initVaultService(
	ctx context.Context,
	stateStore storage.StateStorer,
	signer crypto.Signer,
	chainID int64,
	peerId string,
	vaultLogic common.Address,
	backend transaction.Backend,
	overlayEthAddress common.Address,
	transactionService transaction.Service,
	vaultFactory vault.Factory,
	deployGasPrice string,
	chequeStore vault.ChequeStore,
	erc20Service erc20.Service,
	mpErc20Service map[string]erc20.Service,
) (vault.Service, error) {
	chequeSigner := vault.NewChequeSigner(signer, chainID)

	if deployGasPrice != "" {
		gasPrice, ok := new(big.Int).SetString(deployGasPrice, 10)
		if !ok {
			return nil, fmt.Errorf("deploy gas price \"%s\" cannot be parsed", deployGasPrice)
		}
		ctx = sctx.SetGasPrice(ctx, gasPrice)
	}

	vaultService, err := vault.Init(
		ctx,
		vaultFactory,
		stateStore,
		transactionService,
		backend,
		chainID,
		peerId,
		vaultLogic,
		overlayEthAddress,
		chequeSigner,
		chequeStore,
		erc20Service,
		mpErc20Service,
	)
	if err != nil {
		return nil, fmt.Errorf("vault init: %w", err)
	}

	return vaultService, nil
}

func initChequeStoreCashout(
	stateStore storage.StateStorer,
	swapBackend transaction.Backend,
	vaultFactory vault.Factory,
	chainID int64,
	overlayEthAddress common.Address,
	transactionService transaction.Service,
) (vault.ChequeStore, vault.CashoutService) {
	chequeStore := vault.NewChequeStore(
		stateStore,
		vaultFactory,
		chainID,
		overlayEthAddress,
		transactionService,
		vault.RecoverCheque,
	)

	cashout := vault.NewCashoutService(
		stateStore,
		swapBackend,
		transactionService,
		chequeStore,
	)

	return chequeStore, cashout
}

// InitSwap will initialize and register the swap service.
func initSwap(
	stateStore storage.StateStorer,
	overlayEthAddress common.Address,
	vaultService vault.Service,
	chequeStore vault.ChequeStore,
	cashoutService vault.CashoutService,
	accounting settlement.Accounting,
	priceOracleAddress string,
	chainID int64,
	transactionService transaction.Service,
) (*swap.Service, priceoracle.Service, error) {

	var currentPriceOracleAddress common.Address
	if priceOracleAddress == "" {
		return nil, nil, errors.New("no known price oracle address for this network")
	} else {
		currentPriceOracleAddress = common.HexToAddress(priceOracleAddress)
	}

	priceOracle := priceoracle.New(currentPriceOracleAddress, transactionService)
	_, err := priceOracle.CheckNewPrice(tokencfg.GetWbttToken()) // CheckNewPrice when node starts
	if err != nil {
		return nil, nil, errors.New("CheckNewPrice error, it may happens when contract call failed if bttc chain rpc is down, please try again")
	}

	swapProtocol := swapprotocol.New(overlayEthAddress, priceOracle)
	swapAddressBook := swap.NewAddressbook(stateStore)

	swapService := swap.New(
		swapProtocol,
		stateStore,
		vaultService,
		chequeStore,
		swapAddressBook,
		chainID,
		cashoutService,
		accounting,
	)

	swapProtocol.SetSwap(swapService)
	swapprotocol.SwapProtocol = swapProtocol

	return swapService, priceOracle, nil
}

func GetTxHash(stateStore storage.StateStorer) ([]byte, error) {
	var txHash common.Hash
	key := vault.VaultDeploymentKey
	if err := stateStore.Get(key, &txHash); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, errors.New("vault deployment transaction hash not found, please specify the transaction hash manually")
		}
		return nil, err
	}

	log.Infof("using the vault transaction hash %x", txHash)
	return txHash.Bytes(), nil
}

func GetTxNextBlock(ctx context.Context, backend transaction.Backend, monitor transaction.Monitor, duration time.Duration, trx []byte, blockHash string) ([]byte, error) {

	if blockHash != "" {
		blockHashTrimmed := strings.TrimPrefix(blockHash, "0x")
		if len(blockHashTrimmed) != 64 {
			return nil, errors.New("invalid length")
		}
		blockHash, err := hex.DecodeString(blockHashTrimmed)
		if err != nil {
			return nil, err
		}
		log.Infof("using the provided block hash %x", blockHash)
		return blockHash, nil
	}

	// if not found in statestore, fetch from chain
	tx, err := backend.TransactionReceipt(ctx, common.BytesToHash(trx))
	if err != nil {
		return nil, err
	}

	block, err := transaction.WaitBlock(ctx, backend, duration, big.NewInt(0).Add(tx.BlockNumber, big.NewInt(1)))
	if err != nil {
		return nil, err
	}

	hash := block.Hash()
	hashBytes := hash.Bytes()

	log.Infof("using the next block hash from the blockchain %x", hashBytes)

	return hashBytes, nil
}
