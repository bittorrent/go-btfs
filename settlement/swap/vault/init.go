package vault

import (
	"context"
	"crypto/rand"
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
	vaultKey           = "swap_vault"
	VaultDeploymentKey = "swap_vault_transaction_deployment"

	balanceCheckBackoffDuration = 20 * time.Second
	balanceCheckMaxRetries      = 10
	maxApprove                  = 10000000000000000
	decimals                    = 100000000000000000
)

func checkBalance(
	ctx context.Context,
	swapBackend transaction.Backend,
	chainId int64,
	overlayEthAddress common.Address,
	erc20Token erc20.Service,
) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, balanceCheckBackoffDuration*time.Duration(balanceCheckMaxRetries))
	defer cancel()
	for {
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
			case <-timeoutCtx.Done():
				return fmt.Errorf("insufficient BTT for initial deposit")
			}
			continue
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
	overlayEthAddress common.Address,
	chequeSigner ChequeSigner,
	chequeStore ChequeStore,
) (vaultService Service, err error) {
	// verify that the supplied factory is valid
	err = vaultFactory.VerifyBytecode(ctx)
	if err != nil {
		return nil, err
	}

	erc20Address, err := vaultFactory.ERC20Address(ctx)
	if err != nil {
		return nil, err
	}

	erc20Service := erc20.New(swapBackend, transactionService, erc20Address)

	var vaultAddress common.Address
	err = stateStore.Get(vaultKey, &vaultAddress)
	if err != nil {
		if err != storage.ErrNotFound {
			return nil, err
		}

		var txHash common.Hash
		err = stateStore.Get(VaultDeploymentKey, &txHash)
		if err != nil && err != storage.ErrNotFound {
			return nil, err
		}

		if err == storage.ErrNotFound {
			log.Infof("no vault found, deploying new one.")
			err = checkBalance(ctx, swapBackend, chainId, overlayEthAddress, erc20Service)
			if err != nil {
				return nil, err
			}

			nonce := make([]byte, 32)
			_, err = rand.Read(nonce)
			if err != nil {
				return nil, err
			}

			// if we don't yet have a vault, deploy a new one
			txHash, err = vaultFactory.Deploy(ctx, overlayEthAddress, common.BytesToHash(nonce), peerId)
			if err != nil {
				return nil, err
			}

			log.Infof("deploying new vault in transaction %x", txHash)

			err = stateStore.Put(VaultDeploymentKey, txHash)
			if err != nil {
				return nil, err
			}
		} else {
			log.Infof("waiting for vault deployment in transaction %x", txHash)
		}

		vaultAddress, err = vaultFactory.WaitDeployed(ctx, txHash)
		if err != nil {
			return nil, err
		}

		log.Infof("deployed vault at address 0x%x", vaultAddress)

		// save the address for later use
		err = stateStore.Put(vaultKey, vaultAddress)
		if err != nil {
			return nil, err
		}

		vaultService, err = New(transactionService, vaultAddress, overlayEthAddress, stateStore, chequeSigner, erc20Service, chequeStore)
		if err != nil {
			return nil, err
		}
	} else {
		vaultService, err = New(transactionService, vaultAddress, overlayEthAddress, stateStore, chequeSigner, erc20Service, chequeStore)
		if err != nil {
			return nil, err
		}

		log.Infof("using existing vault 0x%x", vaultAddress)
	}

	fmt.Printf("self vault: 0x%x \n", vaultAddress)

	// regardless of how the vault service was initialised make sure that the vault is valid
	err = vaultFactory.VerifyVault(ctx, vaultService.Address())
	if err != nil {
		return nil, err
	}

	// approve to vaultAddress
	allowance, err := erc20Service.Allowance(ctx, overlayEthAddress, vaultAddress)
	if err != nil {
		return nil, err
	}

	if allowance.Cmp(big.NewInt(0)) == 0 {
		var value big.Int
		value.Mul(big.NewInt(maxApprove), big.NewInt(decimals))
		hash, err := erc20Service.Approve(ctx, vaultAddress, &value)
		if err != nil {
			return nil, err
		}

		fmt.Printf("approve WBTT to vault [0x%x] at tx [0x%x] \n", vaultAddress, hash)
	}
	return vaultService, nil
}
