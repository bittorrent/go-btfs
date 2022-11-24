package vault

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/bittorrent/go-btfs/settlement/swap/erc20"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/ethereum/go-ethereum/common"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("vault:init")

const (
	VaultKey           = "swap_vault"
	VaultDeploymentKey = "swap_vault_transaction_deployment"

	balanceCheckBackoffDuration = 20 * time.Second
	balanceCheckMaxRetries      = 10
	maxApprove                  = 10000000000000000
	decimals                    = 100000000000000000
)

func checkBalance(ctx context.Context, swapBackend transaction.Backend, overlayEthAddress common.Address) error {
	for {
		timeoutCtx, _ := context.WithTimeout(ctx, balanceCheckBackoffDuration*time.Duration(balanceCheckMaxRetries))
		ethBalance, err := swapBackend.BalanceAt(timeoutCtx, overlayEthAddress, nil)
		if err != nil {
			return err
		}

		gasPrice, err := swapBackend.SuggestGasPrice(timeoutCtx)
		if err != nil {
			return err
		}

		minimumEth := gasPrice.Mul(gasPrice, big.NewInt(300000))
		insufficientETH := ethBalance.Cmp(minimumEth) < 0
		if insufficientETH {
			fmt.Printf("cannot continue until there is sufficient (100 Suggested) BTT (for Gas) available on 0x%x \n", overlayEthAddress)
			select {
			case <-time.After(balanceCheckBackoffDuration):
				continue
			}
		}
		return nil
	}
}

// Init initialises the vault service.
func Init(
	ctx context.Context,
	vaultFactory Factory,
	stateStore storage.StateStorer,
	transactionService transaction.Service,
	swapBackend transaction.Backend,
	chainId int64,
	peerId string,
	vaultLogicAddress common.Address,
	overlayEthAddress common.Address,
	chequeSigner ChequeSigner,
	chequeStore ChequeStore,
	erc20Service erc20.Service,
	mpErc20Service map[string]erc20.Service,
) (vaultService Service, err error) {

	// verify that the supplied factory is valid
	err = vaultFactory.VerifyBytecode(ctx)
	if err != nil {
		return nil, err
	}

	// deploy vault if it not exist
	var vaultAddress common.Address
	tokenAddress := erc20Service.Address(ctx)
	vaultAddress, err = deployVaultIfNotExist(ctx, peerId, vaultFactory, tokenAddress, vaultLogicAddress, overlayEthAddress, stateStore)
	if err != nil {
		return nil, err
	}
	fmt.Printf("self vault: 0x%x \n", vaultAddress)

	// regardless of how the vault service was initialised make sure that the vault is valid
	err = vaultFactory.VerifyVault(ctx, vaultAddress)
	if err != nil {
		return nil, err
	}

	// approve to vaultAddress
	err = erc20tokenApprove(ctx, "WBTT", erc20Service, overlayEthAddress, vaultAddress)
	if err != nil {
		return nil, err
	}

	// muti tokens
	for tokenStr, erc20Svr := range mpErc20Service {
		err = erc20tokenApprove(ctx, tokenStr, erc20Svr, overlayEthAddress, vaultAddress)
		if err != nil {
			return nil, err
		}
	}

	vaultService, err = New(transactionService, vaultAddress, overlayEthAddress, stateStore, chequeSigner, erc20Service, mpErc20Service, chequeStore)
	return vaultService, err
}

// GetStoredVaultAddr returns vault address stored in stateStore. If you want exactly result, pls query the bttc chain.
func GetStoredVaultAddr(stateStore storage.StateStorer) (vault common.Address, err error) {
	err = stateStore.Get(VaultKey, &vault)
	if err == nil {
		return
	} else {
		if err == storage.ErrNotFound {
			return common.Address{}, nil
		}
		return
	}
}

func deployVaultIfNotExist(
	ctx context.Context,
	peerId string,
	vaultFactory Factory,
	tokenAddress common.Address,
	vaultLogicAddress common.Address,
	overlayEthAddress common.Address,
	stateStore storage.StateStorer,
) (vaultAddress common.Address, err error) {

	zeroAddr := common.Address{}

	err = stateStore.Get(VaultKey, &vaultAddress)
	if err == nil {
		log.Infof("using existing vault 0x%x", vaultAddress)
		return vaultAddress, nil
	}
	if err != storage.ErrNotFound {
		return zeroAddr, err
	}

	var txHash common.Hash
	err = stateStore.Get(VaultDeploymentKey, &txHash)
	if err != nil && err != storage.ErrNotFound {
		return zeroAddr, err
	}

	if err == storage.ErrNotFound {
		log.Infof("no vault found, deploying new one.")
		vaultAddress, txHash, err = vaultFactory.Deploy(ctx, overlayEthAddress, vaultLogicAddress, peerId, tokenAddress)
		if err != nil {
			return zeroAddr, err
		}
		// existing vault got from factory
		if vaultAddress != zeroAddr {
			log.Infof("using existing vault 0x%x !", vaultAddress)
			err = stateStore.Put(VaultKey, vaultAddress)
			return vaultAddress, nil
		}

		log.Infof("deploying new vault in transaction %x", txHash)
		err = stateStore.Put(VaultDeploymentKey, txHash)
		if err != nil {
			return zeroAddr, err
		}
	} else {
		log.Infof("waiting for vault deployment in transaction %x", txHash)
	}

	vaultAddress, err = vaultFactory.WaitDeployed(ctx, txHash)
	if err != nil {
		return zeroAddr, err
	}

	log.Infof("deployed vault at address 0x%x", vaultAddress)
	err = stateStore.Put(VaultKey, vaultAddress)
	return vaultAddress, err
}

func erc20tokenApprove(ctx context.Context, tokenStr string, erc20Service erc20.Service, issuer, vault common.Address) error {
	allowance, err := erc20Service.Allowance(ctx, issuer, vault)
	if err != nil {
		return err
	}

	if allowance.Cmp(big.NewInt(0)) == 0 {
		var value big.Int
		value.Mul(big.NewInt(maxApprove), big.NewInt(decimals))
		hash, err := erc20Service.Approve(ctx, vault, &value)
		if err != nil {
			return err
		}
		log.Infof("approve %s to vault [0x%x] at tx [0x%x] \n", tokenStr, vault, hash)
	}
	return nil
}
