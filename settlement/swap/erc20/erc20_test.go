package erc20_test

import (
	"context"
	"math/big"
	"testing"

	conabi "github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/settlement/swap/erc20"
	"github.com/bittorrent/go-btfs/transaction"
	backendmock "github.com/bittorrent/go-btfs/transaction/backendmock"
	transactionmock "github.com/bittorrent/go-btfs/transaction/mock"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	endPoint = "https://rpc.bittorrentchain.io"

	SenderAdd = common.HexToAddress("0x44721adf10BB3a76Ce9B456f53Ce9F652be9a2e6")
	erc20Add  = common.HexToAddress("0x23181F21DEa5936e24163FFABa4Ea3B316B57f3C")

	erc20ABI          = transaction.ParseABIUnchecked(conabi.Erc20ABI)
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

func TestDeposit(t *testing.T) {
	erc20Address := common.HexToAddress("0x00")
	expectedTxHash := common.BytesToHash(common.HexToAddress("0x01").Bytes())
	depositAmount := big.NewInt(200)

	erc20 := erc20.New(
		backendmock.New(
			backendmock.WithEstimateGasFunc(func(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
				return 0, nil
			}),
		),
		transactionmock.New(
			transactionmock.WithABISend(
				&erc20ABI,
				expectedTxHash,
				erc20Address,
				depositAmount,
				"deposit",
			),
		),
		erc20Address,
	)

	txHash, err := erc20.Deposit(context.Background(), depositAmount)
	if err != nil {
		t.Fatal(err)
	}

	if expectedTxHash != txHash {
		t.Fatalf("got wrong balance. wanted %d, got %d", expectedTxHash, txHash)
	}
}

func TestWithdraw(t *testing.T) {
	erc20Address := common.HexToAddress("0x00")
	expectedTxHash := common.BytesToHash(common.HexToAddress("0x01").Bytes())
	withdrawAmount := big.NewInt(200)

	erc20 := erc20.New(
		backendmock.New(),
		transactionmock.New(
			transactionmock.WithABISend(
				&erc20ABI,
				expectedTxHash,
				erc20Address,
				big.NewInt(0),
				"withdraw",
				withdrawAmount,
			),
		),
		erc20Address,
	)

	txHash, err := erc20.Withdraw(context.Background(), withdrawAmount)
	if err != nil {
		t.Fatal(err)
	}

	if expectedTxHash != txHash {
		t.Fatalf("got wrong balance. wanted %d, got %d", expectedTxHash, txHash)
	}
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

func TestAllowance(t *testing.T) {
	vaultAddress := common.HexToAddress("0xabcd")
	issue := common.HexToAddress("01")
	// value := big.NewInt(20)
	txHash := common.HexToHash("0xdddd")
	result := big.NewInt(1)

	erc20 := erc20.New(
		backendmock.New(),
		transactionmock.New(
			transactionmock.WithABICall(&erc20ABI, vaultAddress, result.FillBytes(make([]byte, 32)), "allowance", issue, vaultAddress),
		),
		vaultAddress,
	)

	num, err := erc20.Allowance(context.Background(), issue, vaultAddress)
	if err != nil {
		t.Fatal(err)
	}

	if num.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("returned wrong transaction hash. wanted %v, got %v", txHash, num)
	}
}

func TestApprove(t *testing.T) {
	toAddress := common.HexToAddress("0xabcd")
	value := big.NewInt(20)
	txHash := common.HexToHash("0xdddd")

	erc20 := erc20.New(
		backendmock.New(),
		transactionmock.New(
			transactionmock.WithABISend(&erc20ABI, txHash, toAddress, big.NewInt(0), "approve", toAddress, value),
		),
		toAddress,
	)

	resultHashTx, err := erc20.Approve(context.Background(), toAddress, value)
	if err != nil {
		t.Fatal(err)
	}

	if resultHashTx != txHash {
		t.Fatalf("returned wrong transaction hash. wanted %v, got %v", txHash, resultHashTx)
	}
}

func TestTransferFrom(t *testing.T) {
	issue := common.HexToAddress("0xabcdfg")
	vaultAddress := common.HexToAddress("0xabcd")
	value := big.NewInt(20)
	txHash := common.HexToHash("0xdddd")

	erc20 := erc20.New(
		backendmock.New(),
		transactionmock.New(
			transactionmock.WithABISend(&erc20ABI, txHash, vaultAddress, big.NewInt(0), "transferFrom", issue, vaultAddress, value),
		),
		vaultAddress,
	)

	resultHashTx, err := erc20.TransferFrom(context.Background(), issue, vaultAddress, value)
	if err != nil {
		t.Fatal(err)
	}

	if resultHashTx != txHash {
		t.Fatalf("returned wrong transaction hash. wanted %v, got %v", txHash, resultHashTx)
	}
}
