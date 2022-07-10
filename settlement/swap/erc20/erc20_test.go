package erc20_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/bittorrent/go-btfs/settlement/swap/erc20"
	"github.com/bittorrent/go-btfs/transaction"
	backendmock "github.com/bittorrent/go-btfs/transaction/backendmock"
	transactionmock "github.com/bittorrent/go-btfs/transaction/mock"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const Erc20ABI = `[
			{
				"inputs": [
					{
						"internalType": "string",
						"name": "name_",
						"type": "string"
					},
					{
						"internalType": "string",
						"name": "symbol_",
						"type": "string"
					}
				],
				"stateMutability": "nonpayable",
				"type": "constructor"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "owner",
						"type": "address"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "spender",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "value",
						"type": "uint256"
					}
				],
				"name": "Approval",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "from",
						"type": "address"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "to",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "value",
						"type": "uint256"
					}
				],
				"name": "Transfer",
				"type": "event"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "owner",
						"type": "address"
					},
					{
						"internalType": "address",
						"name": "spender",
						"type": "address"
					}
				],
				"name": "allowance",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "spender",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "approve",
				"outputs": [
					{
						"internalType": "bool",
						"name": "",
						"type": "bool"
					}
				],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "account",
						"type": "address"
					}
				],
				"name": "balanceOf",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "decimals",
				"outputs": [
					{
						"internalType": "uint8",
						"name": "",
						"type": "uint8"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "spender",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "subtractedValue",
						"type": "uint256"
					}
				],
				"name": "decreaseAllowance",
				"outputs": [
					{
						"internalType": "bool",
						"name": "",
						"type": "bool"
					}
				],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "spender",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "addedValue",
						"type": "uint256"
					}
				],
				"name": "increaseAllowance",
				"outputs": [
					{
						"internalType": "bool",
						"name": "",
						"type": "bool"
					}
				],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "account",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "mint",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "name",
				"outputs": [
					{
						"internalType": "string",
						"name": "",
						"type": "string"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "symbol",
				"outputs": [
					{
						"internalType": "string",
						"name": "",
						"type": "string"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "totalSupply",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "recipient",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "transfer",
				"outputs": [
					{
						"internalType": "bool",
						"name": "",
						"type": "bool"
					}
				],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "sender",
						"type": "address"
					},
					{
						"internalType": "address",
						"name": "recipient",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					}
				],
				"name": "transferFrom",
				"outputs": [
					{
						"internalType": "bool",
						"name": "",
						"type": "bool"
					}
				],
				"stateMutability": "nonpayable",
				"type": "function"
			}
		]`

const OracleAbi = `[
			{
				"inputs": [
					{
						"internalType": "uint256",
						"name": "_price",
						"type": "uint256"
					}
				],
				"stateMutability": "nonpayable",
				"type": "constructor"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "previousOwner",
						"type": "address"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "newOwner",
						"type": "address"
					}
				],
				"name": "OwnershipTransferred",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "price",
						"type": "uint256"
					}
				],
				"name": "PriceUpdate",
				"type": "event"
			},
			{
				"inputs": [],
				"name": "getPrice",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "owner",
				"outputs": [
					{
						"internalType": "address",
						"name": "",
						"type": "address"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "price",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "renounceOwnership",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "newOwner",
						"type": "address"
					}
				],
				"name": "transferOwnership",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "uint256",
						"name": "newPrice",
						"type": "uint256"
					}
				],
				"name": "updatePrice",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			}
		]`

var (
	endPoint = "http://18.144.29.246:8110"

	SenderAdd = common.HexToAddress("0xA4E7663A031ca1f67eEa828E4795653504d38c6e")
	erc20Add  = common.HexToAddress("0xD26c3d45a805a5f7809E27Bd18949d559e281900")

	erc20ABI          = transaction.ParseABIUnchecked(Erc20ABI)
	defaultPrivateKey = "e484c4373db5c55a9813e4abbb74a15edd794019b8db4365a876ed538622bcf9"
)

func DialBackend() (*ethclient.Client, error) {
	client, err := ethclient.Dial(endPoint)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getChainID(client *ethclient.Client) (*big.Int, error) {
	chainID, err := client.ChainID(context.Background())

	if err != nil {
		return nil, err
	}

	return chainID, nil
}

func TestBalanceOf(t *testing.T) {
	erc20Address := common.HexToAddress("00")
	account := common.HexToAddress("01")
	expectedBalance := big.NewInt(100)

	erc20 := erc20.New(
		backendmock.New(),
		transactionmock.New(
			transactionmock.WithABICall(
				&erc20ABI,
				erc20Address,
				expectedBalance.FillBytes(make([]byte, 32)),
				"balanceOf",
				account,
			),
		),
		erc20Address,
	)

	balance, err := erc20.BalanceOf(context.Background(), account)
	if err != nil {
		t.Fatal(err)
	}

	if expectedBalance.Cmp(balance) != 0 {
		t.Fatalf("got wrong balance. wanted %d, got %d", expectedBalance, balance)
	}
}

//TODO: FIX ME
func TestRealBalance(t *testing.T) {
	backend, err := DialBackend()

	if err != nil {
		t.Fatal(err)
	}

	callData, err := erc20ABI.Pack("balanceOf", SenderAdd)
	if err != nil {
		t.Fatal(err)
	}

	msg := ethereum.CallMsg{
		From:     SenderAdd,
		To:       &erc20Add,
		Data:     callData,
		GasPrice: big.NewInt(0),
		Gas:      0,
		Value:    big.NewInt(0),
	}
	data, err := backend.CallContract(context.Background(), msg, nil)
	if err != nil {
		t.Fatal(err)
	}

	results, err := erc20ABI.Unpack("balanceOf", data)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatal(err)
	}

	balance, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || balance == nil {
		t.Fatal(err)
	}

	t.Log("real balance is: ", balance)
}

func TestTransfer(t *testing.T) {
	address := common.HexToAddress("0xabcd")
	account := common.HexToAddress("01")
	value := big.NewInt(20)
	txHash := common.HexToHash("0xdddd")

	erc20 := erc20.New(
		backendmock.New(),
		transactionmock.New(
			transactionmock.WithABISend(&erc20ABI, txHash, address, big.NewInt(0), "transfer", account, value),
		),
		address,
	)

	returnedTxHash, err := erc20.Transfer(context.Background(), account, value)
	if err != nil {
		t.Fatal(err)
	}

	if txHash != returnedTxHash {
		t.Fatalf("returned wrong transaction hash. wanted %v, got %v", txHash, returnedTxHash)
	}
}
