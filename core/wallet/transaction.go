package wallet

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	walletpb "github.com/bittorrent/go-btfs/protos/wallet"

	config "github.com/TRON-US/go-btfs-config"
	escrowPb "github.com/tron-us/go-btfs-common/protos/escrow"
	exPb "github.com/tron-us/go-btfs-common/protos/exchange"
	ledgerPb "github.com/tron-us/go-btfs-common/protos/ledger"
	tronPb "github.com/tron-us/go-btfs-common/protos/protocol/api"
	corePb "github.com/tron-us/go-btfs-common/protos/protocol/core"
	"github.com/tron-us/go-btfs-common/utils/grpc"

	"github.com/gogo/protobuf/proto"
	ds "github.com/ipfs/go-datastore"
)

var (
	ErrInsufficientExchangeBalanceOnTron   = errors.New("exchange balance on Tron network is not sufficient")
	ErrInsufficientUserBalanceOnTron       = errors.New(fmt.Sprint("User balance on tron network is not sufficient."))
	ErrInsufficientUserBalanceOnLedger     = errors.New("rpc error: code = ResourceExhausted desc = NSF")
	ErrInsufficientExchangeBalanceOnLedger = errors.New("exchange balance on Private Ledger is not sufficient")
)

// Do the deposit action, integrate exchange's PrepareDeposit and Deposit API.
func Deposit(ctx context.Context, n *core.IpfsNode, ledgerAddr []byte, amount int64,
	privateKey *ecdsa.PrivateKey, runDaemon bool, async bool) (*exPb.PrepareDepositResponse, error) {
	log.Debug("Deposit begin!")
	//PrepareDeposit
	prepareResponse, err := PrepareDeposit(ctx, ledgerAddr, amount)
	if err != nil {
		return nil, err
	}
	if prepareResponse.Response.Code != exPb.Response_SUCCESS {
		return prepareResponse, errors.New(string(prepareResponse.Response.ReturnMessage))
	}
	log.Debug(fmt.Sprintf("PrepareDeposit success, id: [%d]", prepareResponse.GetId()))

	//Do the DepositRequest.
	depositResponse, err := DepositRequest(ctx, prepareResponse, privateKey)
	if err != nil {
		return prepareResponse, err
	}
	if depositResponse.Response.Code != exPb.Response_SUCCESS {
		return prepareResponse, errors.New(string(depositResponse.Response.ReturnMessage))
	}
	log.Debug(fmt.Sprintf("Call Deposit API success, id: [%d]", prepareResponse.GetId()))

	err = PersistTx(n.Repo.Datastore(), n.Identity.Pretty(), strconv.FormatInt(prepareResponse.GetId(), 10),
		amount, BttWallet, InAppWallet, StatusPending, walletpb.TransactionV1_EXCHANGE)
	if err != nil {
		return nil, err
	}

	if runDaemon && async {
		go ConfirmDepositProcess(context.Background(), n, prepareResponse, privateKey)
	} else {
		err := ConfirmDepositProcess(ctx, n, prepareResponse, privateKey)
		if err != nil {
			return nil, err
		}
	}
	// Doing confirm deposit.

	log.Debug("Deposit end!")
	return prepareResponse, nil
}

// Call exchange's PrepareDeposit API.
func PrepareDeposit(ctx context.Context, ledgerAddr []byte, amount int64) (*exPb.PrepareDepositResponse, error) {
	// Prepare to Deposit.
	var err error
	var prepareResponse *exPb.PrepareDepositResponse
	err = grpc.ExchangeClient(exchangeService).WithContext(ctx,
		func(ctx context.Context, client exPb.ExchangeClient) error {
			prepareDepositRequest := &exPb.PrepareDepositRequest{Amount: amount, OutTxId: time.Now().UnixNano(),
				UserAddress: ledgerAddr}
			prepareResponse, err = client.PrepareDeposit(ctx, prepareDepositRequest)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return prepareResponse, nil
}

// Call exchange's Deposit API
func DepositRequest(ctx context.Context, prepareResponse *exPb.PrepareDepositResponse, privateKey *ecdsa.PrivateKey) (*exPb.DepositResponse, error) {
	// Sign Tron Transaction.
	tronTransaction := prepareResponse.GetTronTransaction()
	for range tronTransaction.GetRawData().GetContract() {
		signature, err := Sign(tronTransaction, privateKey)
		if err != nil {
			return nil, err
		}
		tronTransaction.Signature = append(tronTransaction.GetSignature(), signature)
	}

	var err error
	var depositResponse *exPb.DepositResponse
	err = grpc.ExchangeClient(exchangeService).WithContext(ctx,
		func(ctx context.Context, client exPb.ExchangeClient) error {
			depositRequest := &exPb.DepositRequest{Id: prepareResponse.GetId(), SignedTronTransaction: tronTransaction}

			depositResponse, err = client.Deposit(ctx, depositRequest)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return depositResponse, nil
}

// Continuous call ConfirmDeposit until it responses a FAILED or SUCCESS.
func ConfirmDepositProcess(ctx context.Context, n *core.IpfsNode, prepareResponse *exPb.PrepareDepositResponse,
	privateKey *ecdsa.PrivateKey) error {
	log.Debug(fmt.Sprintf("[Id:%d] ConfirmDepositProcess begin.", prepareResponse.GetId()))

	// Continuous call ConfirmDeposit until it responses a FAILED or SUCCESS.
	now := time.Now().UnixNano() / 1e6
	for time.Now().UnixNano()/1e6 < (now + 7*86400*1000) {
		time.Sleep(5 * time.Second)

		// ConfirmDeposit after 1min, because watcher need wait for 1min to confirm tron transaction.
		if time.Now().UnixNano()/1e6 < (prepareResponse.GetTronTransaction().GetRawData().GetTimestamp() + 60*1000) {
			continue
		}
		log.Debug(fmt.Sprintf("[Id:%d] ConfirmDeposit begin.", prepareResponse.GetId()))

		confirmDepositResponse, err := ConfirmDeposit(ctx, prepareResponse.GetId())
		if err != nil {
			continue
		}
		if confirmDepositResponse == nil || confirmDepositResponse.SuccessChannelState == nil {
			continue
		}
		err = save(n.Repo.Datastore(), confirmDepositResponse.SuccessChannelState.Channel)
		if err != nil {
			log.Debug(err)
		}
		txId := strconv.FormatInt(prepareResponse.GetId(), 10)
		if confirmDepositResponse.GetResponse().GetCode() == exPb.Response_TRANSACTION_PENDING {
			log.Debug(fmt.Sprintf("[Id:%d]TronTransaction is PENDING.", prepareResponse.GetId()))

			err := UpdateStatus(n.Repo.Datastore(), n.Identity.Pretty(), txId, StatusPending)
			if err != nil {
				return err
			}
			continue
		}
		if confirmDepositResponse.GetResponse().GetCode() == exPb.Response_TRANSACTION_FAILED {
			return UpdateStatus(n.Repo.Datastore(), n.Identity.Pretty(), txId, StatusFailed)
		}
		if confirmDepositResponse.GetResponse().GetCode() == exPb.Response_SUCCESS {
			signSuccessChannelState := confirmDepositResponse.GetSuccessChannelState()
			if signSuccessChannelState != nil {
				toSignature, err := Sign(signSuccessChannelState.GetChannel(), privateKey)
				if err != nil {
					return err
				}
				signSuccessChannelState.ToSignature = toSignature
			} else {
				return fmt.Errorf("[Id:%d] SignSuccessChannelState is nil", prepareResponse.GetId())
			}

			for i := 0; i < 10; i++ {
				err = grpc.EscrowClient(escrowService).WithContext(ctx,
					func(ctx context.Context, client escrowPb.EscrowServiceClient) error {
						_, err = client.CloseChannel(ctx, signSuccessChannelState)
						if err != nil {
							return err
						}
						return nil
					})
				if err == nil {
					break
				}
				time.Sleep(5 * time.Second)
			}
			if err != nil {
				return err
			}
			log.Info(fmt.Sprintf("[Id:%d] Close SuccessChannelState succeed.", prepareResponse.GetId()))
			err = UpdateStatus(n.Repo.Datastore(), n.Identity.Pretty(), txId, StatusSuccess)
			if err != nil {
				return err
			}
			return nil
		}
	}
	if time.Now().UnixNano()/1e6 >= (prepareResponse.GetTronTransaction().GetRawData().GetExpiration() + 240*1000) {
		err := fmt.Errorf("[Id:%d] Didn't get the tron transaction results until the expiration time.",
			prepareResponse.GetId())
		return err
	}
	return nil
}

// Call exchange's ConfirmDeposit API.
func ConfirmDeposit(ctx context.Context, logId int64) (*exPb.ConfirmDepositResponse, error) {
	var err error
	var confirmDepositResponse *exPb.ConfirmDepositResponse
	err = grpc.ExchangeClient(exchangeService).WithContext(ctx,
		func(ctx context.Context, client exPb.ExchangeClient) error {
			confirmDepositRequest := &exPb.ConfirmDepositRequest{Id: logId}
			confirmDepositResponse, err = client.ConfirmDeposit(ctx, confirmDepositRequest)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return confirmDepositResponse, nil
}

// Do the withdraw action, integrate exchange's PrepareWithdraw and Withdraw API, return channel id and error.
// If Withdraw succeed, with return channel id and logInfo id, error is nil; otherwise will return error, channel id
// and logInfo id is 0.
func Withdraw(ctx context.Context, n *core.IpfsNode, ledgerAddr, externalAddr []byte, amount int64,
	privateKey *ecdsa.PrivateKey) (int64, int64, error) {
	log.Debug("Withdraw begin!")
	outTxId := time.Now().UnixNano()
	//PrepareWithdraw
	prepareResponse, err := PrepareWithdraw(ctx, ledgerAddr, externalAddr, amount, outTxId)
	if err != nil {
		return 0, 0, err
	}
	if prepareResponse.Response.Code != exPb.Response_SUCCESS {
		return 0, 0, errors.New(string(prepareResponse.Response.ReturnMessage))
	}
	log.Debug(fmt.Sprintf("Prepare withdraw success, id: [%d]", prepareResponse.GetId()))

	channelCommit := &ledgerPb.ChannelCommit{
		Payer:     &ledgerPb.PublicKey{Key: ledgerAddr},
		Recipient: &ledgerPb.PublicKey{Key: prepareResponse.GetLedgerExchangeAddress()},
		Amount:    amount,
		PayerId:   time.Now().UnixNano() + prepareResponse.GetId(),
	}
	//Sign channel commit.
	signature, err := Sign(channelCommit, privateKey)
	if err != nil {
		return 0, 0, err
	}

	var channelId *ledgerPb.ChannelID
	err = grpc.EscrowClient(escrowService).WithContext(ctx,
		func(ctx context.Context, client escrowPb.EscrowServiceClient) error {
			channelId, err = client.CreateChannel(ctx,
				&ledgerPb.SignedChannelCommit{Channel: channelCommit, Signature: signature})
			if err != nil {
				if err.Error() == ErrInsufficientUserBalanceOnLedger.Error() {
					return ErrInsufficientUserBalanceOnLedger
				}
				return err
			}
			return nil
		})
	if err != nil {
		return 0, 0, err
	}
	log.Debug(fmt.Sprintf("CreateChannel success, channelId: [%d]", channelId.GetId()))

	//Do the WithdrawRequest.
	withdrawResponse, err := WithdrawRequest(ctx, channelId, ledgerAddr, amount, prepareResponse, privateKey)
	if err != nil {
		return 0, 0, err
	}

	txId := strconv.FormatInt(prepareResponse.GetId(), 10)
	err = PersistTx(n.Repo.Datastore(), n.Identity.Pretty(), txId, amount,
		InAppWallet, BttWallet, StatusPending, walletpb.TransactionV1_EXCHANGE)
	if err != nil {
		return 0, 0, err
	}

	if withdrawResponse.Response.Code != exPb.Response_SUCCESS {
		err := UpdateStatus(n.Repo.Datastore(), n.Identity.Pretty(), txId, StatusFailed)
		if err != nil {
			return 0, 0, err
		}
		return 0, 0, errors.New(string(withdrawResponse.Response.ReturnMessage))
	}
	log.Debug("Withdraw end!")
	err = UpdateStatus(n.Repo.Datastore(), n.Identity.Pretty(), txId, StatusSuccess)
	if err != nil {
		return 0, 0, err
	}
	return channelId.Id, prepareResponse.GetId(), nil
}

// Call exchange's Withdraw API
func PrepareWithdraw(ctx context.Context, ledgerAddr, externalAddr []byte, amount, outTxId int64) (
	*exPb.PrepareWithdrawResponse, error) {
	var err error
	var prepareResponse *exPb.PrepareWithdrawResponse
	err = grpc.ExchangeClient(exchangeService).WithContext(ctx,
		func(ctx context.Context, client exPb.ExchangeClient) error {
			prepareWithdrawRequest := &exPb.PrepareWithdrawRequest{
				Amount: amount, OutTxId: outTxId, UserAddress: ledgerAddr, UserExternalAddress: externalAddr}
			prepareResponse, err = client.PrepareWithdraw(ctx, prepareWithdrawRequest)
			if err != nil {
				return err
			}
			log.Debug(prepareResponse)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return prepareResponse, nil
}

// Call exchange's PrepareWithdraw API
func WithdrawRequest(ctx context.Context, channelId *ledgerPb.ChannelID, ledgerAddr []byte, amount int64,
	prepareResponse *exPb.PrepareWithdrawResponse, privateKey *ecdsa.PrivateKey) (*exPb.WithdrawResponse, error) {
	//make signed success channel state.
	successChannelState := &ledgerPb.ChannelState{
		Id:       channelId,
		Sequence: 1,
		From: &ledgerPb.Account{
			Address: &ledgerPb.PublicKey{
				Key: ledgerAddr,
			},
			Balance: 0,
		},
		To: &ledgerPb.Account{
			Address: &ledgerPb.PublicKey{
				Key: prepareResponse.GetLedgerExchangeAddress(),
			},
			Balance: amount,
		},
	}
	successSignature, err := Sign(successChannelState, privateKey)
	if err != nil {
		return nil, err
	}
	successChannelStateSigned := &ledgerPb.SignedChannelState{Channel: successChannelState, FromSignature: successSignature}

	//make signed fail channel state.
	failChannelState := &ledgerPb.ChannelState{
		Id:       channelId,
		Sequence: 1,
		From: &ledgerPb.Account{
			Address: &ledgerPb.PublicKey{
				Key: ledgerAddr,
			},
			Balance: amount,
		},
		To: &ledgerPb.Account{
			Address: &ledgerPb.PublicKey{
				Key: prepareResponse.GetLedgerExchangeAddress(),
			},
			Balance: 0,
		},
	}
	failSignature, err := Sign(failChannelState, privateKey)
	if err != nil {
		return nil, err
	}

	var withdrawResponse *exPb.WithdrawResponse
	err = grpc.ExchangeClient(exchangeService).WithContext(ctx,
		func(ctx context.Context, client exPb.ExchangeClient) error {
			failChannelStateSigned := &ledgerPb.SignedChannelState{Channel: failChannelState, FromSignature: failSignature}
			//Post the withdraw request.
			withdrawRequest := &exPb.WithdrawRequest{
				Id:                  prepareResponse.GetId(),
				SuccessChannelState: successChannelStateSigned,
				FailureChannelState: failChannelStateSigned,
			}
			withdrawResponse, err = client.Withdraw(ctx, withdrawRequest)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return withdrawResponse, nil
}

// Get the token balance on tron blockchain
func GetTokenBalance(ctx context.Context, addr []byte, tokenId string) (int64, error) {
	var tokenBalance int64 = 0
	err := grpc.SolidityClient(solidityService).WithContext(ctx,
		func(ctx context.Context, client tronPb.WalletSolidityClient) error {
			account := &corePb.Account{Address: addr}
			myAccount, err := client.GetAccount(ctx, account)
			if err != nil {
				return err
			}
			tokenMap := myAccount.GetAssetV2()
			if tokenMap == nil || len(tokenMap) == 0 {
				return nil
			}
			tokenBalance = tokenMap[tokenId]
			return nil
		})
	if err != nil {
		return 0, err
	}
	return tokenBalance, nil
}

var (
	walletTransactionKeyPrefix = "/btfs/%v/wallet/transactions/"
	walletTransactionKey       = walletTransactionKeyPrefix + "%v/"

	walletTransactionV1KeyPrefix = "/btfs/%v/wallet/v1/transactions/"
	walletTransactionV1Key       = walletTransactionV1KeyPrefix + "%v/"
)

func PersistTx(d ds.Datastore, peerId string, txId string, amount int64,
	from string, to string, status string, txType walletpb.TransactionV1_Type) error {
	return PersistTxWithTime(d, peerId, txId, amount, from, to, status, txType, time.Now())
}

func PersistTxWithTime(d ds.Datastore, peerId string, txId string, amount int64,
	from string, to string, status string, txType walletpb.TransactionV1_Type, timeCreate time.Time) error {
	return sessions.Save(d, fmt.Sprintf(walletTransactionV1Key, peerId, txId),
		&walletpb.TransactionV1{
			Id:         txId,
			TimeCreate: timeCreate,
			Amount:     amount,
			From:       from,
			To:         to,
			Status:     status,
			Type:       txType,
		})
}

func UpdateStatus(d ds.Datastore, peerId string, txId string, status string) error {
	key := fmt.Sprintf(walletTransactionV1Key, peerId, txId)
	s := &walletpb.TransactionV1{}
	err := sessions.Get(d, key, s)
	if err != nil {
		return err
	}
	if s.Status != status {
		s.Status = status
		return sessions.Save(d, key, s)
	}
	return nil
}

func GetTransactions(d ds.Datastore, peerId string) ([]*walletpb.TransactionV1, error) {
	txs := make(TxSlice, 0)
	v0Txs, err := loadV0Txs(d, peerId)
	if err != nil || v0Txs == nil {
		//ignore, NOP
		err = nil
	}
	txs = append(txs, v0Txs...)
	list, err := sessions.List(d, fmt.Sprintf(walletTransactionV1KeyPrefix, peerId))
	if err != nil {
		return nil, err
	}
	for _, bytes := range list {
		tx := new(walletpb.TransactionV1)
		err := proto.Unmarshal(bytes, tx)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	sort.Sort(txs)
	return txs, nil
}

type TxSlice []*walletpb.TransactionV1

func (p TxSlice) Len() int {
	return len(p)
}

func (p TxSlice) Less(i, j int) bool {
	return p[i].TimeCreate.After(p[j].TimeCreate)
}

func (p TxSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func loadV0Txs(d ds.Datastore, peerId string) ([]*walletpb.TransactionV1, error) {
	list, err := sessions.List(d, fmt.Sprintf(walletTransactionKeyPrefix, peerId))
	if err != nil {
		return nil, err
	}
	txs := make([]*walletpb.TransactionV1, 0)
	for _, bytes := range list {
		tx := new(walletpb.Transaction)
		err := proto.Unmarshal(bytes, tx)
		if err != nil {
			return nil, err
		}
		txs = append(txs, &walletpb.TransactionV1{
			Id:         strconv.FormatInt(tx.Id, 10),
			TimeCreate: tx.TimeCreate,
			Amount:     tx.Amount,
			From:       tx.From,
			To:         tx.To,
			Status:     tx.Status,
			Type:       walletpb.TransactionV1_EXCHANGE,
		})
	}
	return txs, nil
}

func UpdatePendingTransactions(ctx context.Context, d ds.Datastore, cfg *config.Config, peerId string) (int, int, error) {
	scv1, ecv1, err := updateV1Txs(ctx, d, cfg, peerId)
	if err != nil {
		return 0, 0, err
	}
	scv0, ecv0, err := updateV0Txs(ctx, d, cfg, peerId)
	if err != nil {
		return 0, 0, err
	}
	return scv1 + scv0, ecv1 + ecv0, nil
}

func CloseLedgerChannel(ctx context.Context, d ds.Datastore, conf *config.Config) error {
	if err := Init(ctx, conf); err != nil {
		return err
	}
	states, err := list(d)
	if err != nil {
		return err
	}
	if len(states) > 0 {
		err = grpc.EscrowClient(escrowService).WithContext(ctx, func(ctx context.Context, c escrowPb.EscrowServiceClient) error {
			now := time.Now()
			errors := make([]error, 0)
			for _, e := range states {
				if now.Sub(e.TimeCreate) < 5*time.Minute {
					continue
				}
				sig, err := Sign(e.State, hostWallet.privateKey)
				if err != nil {
					errors = append(errors, err)
					continue
				}
				_, err = c.CloseChannel(ctx, &ledgerPb.SignedChannelState{
					Channel:     e.State,
					ToSignature: sig,
				})
				if err != nil {
					errors = append(errors, err)
					continue
				}
				if err := rm(ctx, d, e.State.Id.Id); err != nil {
					errors = append(errors, err)
				}
			}
			if len(errors) > 0 {
				return errors[len(errors)-1]
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func updateV0Txs(ctx context.Context, d ds.Datastore, cfg *config.Config, peerId string) (int, int, error) {
	list, err := sessions.List(d, fmt.Sprintf(walletTransactionKeyPrefix, peerId))
	if err != nil {
		return 0, 0, err
	}
	successCount := 0
	errorCount := 0
	for _, bytes := range list {
		tx := new(walletpb.Transaction)
		err := proto.Unmarshal(bytes, tx)
		if err != nil {
			errorCount++
			continue
		}
		if tx.Status != StatusPending {
			continue
		}

		txId := strconv.FormatInt(tx.Id, 10)
		status, err := getExchangeTxStatus(ctx, cfg, txId)
		if err != nil {
			errorCount++
			continue
		}
		if status == StatusPending {
			continue
		}
		err = UpdateStatus(d, peerId, txId, status)
		if err != nil {
			errorCount++
			continue
		}
		successCount++
	}
	return successCount, errorCount, nil
}

func updateV1Txs(ctx context.Context, d ds.Datastore, cfg *config.Config, peerId string) (int, int, error) {
	list, err := sessions.List(d, fmt.Sprintf(walletTransactionV1KeyPrefix, peerId))
	if err != nil {
		return 0, 0, err
	}
	successCount := 0
	errorCount := 0
	for _, bytes := range list {
		tx := new(walletpb.TransactionV1)
		err := proto.Unmarshal(bytes, tx)
		if err != nil {
			errorCount++
			continue
		}
		if tx.Status != StatusPending {
			continue
		}
		switch tx.Type {
		case walletpb.TransactionV1_EXCHANGE:
			status, err := getExchangeTxStatus(ctx, cfg, tx.Id)
			if err != nil {
				errorCount++
				continue
			}
			if status != StatusPending && status != exPb.Response_TRANSACTION_PENDING.String() {
				err := UpdateStatus(d, peerId, tx.Id, status)
				if err != nil {
					errorCount++
					continue
				}
			}
		case walletpb.TransactionV1_ON_CHAIN:
			status, err := getOnChainTxStatus(ctx, d, cfg, peerId, tx.Id)
			if err != nil {
				errorCount++
				continue
			}
			if status != StatusPending {
				err := UpdateStatus(d, peerId, tx.Id, status)
				if err != nil {
					errorCount++
					continue
				}
			}
		case walletpb.TransactionV1_OFF_CHAIN:
		}
	}
	successCount++
	return successCount, errorCount, nil
}

func getExchangeTxStatus(ctx context.Context, cfg *config.Config, txIdStr string) (string, error) {
	txId, err := strconv.ParseInt(txIdStr, 10, 64)
	if err != nil {
		return "", err
	}
	in := &exPb.QueryTransactionRequest{
		Id: txId,
	}
	var resp *exPb.QueryTransactionResponse
	err = grpc.ExchangeClient(cfg.Services.ExchangeDomain).WithContext(ctx, func(ctx context.Context, client exPb.ExchangeClient) error {
		resp, err = client.QueryTransaction(ctx, in)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	status := resp.Response.Code.String()
	switch status {
	case exPb.Response_TRANSACTION_PENDING.String():
		status = StatusPending
	case exPb.Response_TRANSACTION_FAILED.String():
		status = StatusFailed
	case exPb.Response_SUCCESS.String():
		status = StatusSuccess
	}
	return status, err
}

func getOnChainTxStatus(ctx context.Context, d ds.Datastore, cfg *config.Config, peerId string, txId string) (string, error) {
	status := StatusPending
	err := grpc.WalletClient(cfg.Services.FullnodeDomain).WithContext(ctx, func(ctx context.Context, client tronPb.WalletClient) error {
		bytes, err := hex.DecodeString(txId)
		if err != nil {
			return err
		}
		in := &tronPb.BytesMessage{
			Value: bytes,
		}
		resp, err := client.GetTransactionInfoById(ctx, in)
		if err != nil {
			return err
		}
		status = resp.Result.String()
		return nil
	})
	if err != nil {
		return "", err
	}
	if status == corePb.TransactionInfo_SUCESS.String() {
		status = StatusSuccess
	} else if status == corePb.TransactionInfo_FAILED.String() {
		status = StatusFailed
	}
	return status, nil
}
