package wallet

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/thedevsaddam/gojsonq/v2"

	"github.com/bittorrent/go-btfs/core"
	walletpb "github.com/bittorrent/go-btfs/protos/wallet"

	config "github.com/TRON-US/go-btfs-config"
	"github.com/tron-us/go-btfs-common/crypto"
	tronPb "github.com/tron-us/go-btfs-common/protos/protocol/api"

	protocol_core "github.com/tron-us/go-btfs-common/protos/protocol/core"
	"github.com/tron-us/go-btfs-common/utils/grpc"

	"github.com/gogo/protobuf/proto"
	"github.com/ipfs/go-datastore"
	ic "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/mr-tron/base58/base58"
	"github.com/status-im/keycard-go/hexutils"
)

var (
	txUrlTemplate        = "%s/v1/accounts/%s/transactions?limit=200&only_to=true&order_by=block_timestamp,asc&min_block_timestamp=%d"
	curBlockTimestampKey = "/accounts/%s/transactions/current/block_timestamp"
	client               = http.DefaultClient
)

func SyncTxFromTronGrid(ctx context.Context, cfg *config.Config, ds datastore.Datastore) ([]*TxData, error) {
	keys, err := crypto.FromPrivateKey(cfg.Identity.PrivKey)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(txUrlTemplate, cfg.Services.TrongridDomain, keys.Base58Address, 0)
	if v, err := ds.Get(ctx, datastore.NewKey(fmt.Sprintf(curBlockTimestampKey, keys.Base58Address))); err == nil {
		blockTimestamp, err := strconv.ParseInt(string(v), 10, 64)
		url = fmt.Sprintf(txUrlTemplate, cfg.Services.TrongridDomain, keys.Base58Address, blockTimestamp)
		if err != nil {
			return nil, err
		}
	}
	log.Debug("sync tx called", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	jq := gojsonq.New().FromString(string(body))
	if s, ok := jq.Find("success").(bool); !ok || !s {
		return nil, errors.New("fail to get latest transactions")
	}
	i := -1
	tds := make([]*TxData, 0)
	var lastBlockTimestamp int64 = 0
	for true {
		i++
		pfx := fmt.Sprintf("data.[%d]", i)
		if v := jq.Reset().Find(pfx + ".raw_data.contract.[0].parameter.value"); v == nil {
			break
		} else {
			m, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if m["asset_name"] != getTokenId(cfg) {
				continue
			}
			from, err := hexToBase58(m["owner_address"].(string))
			if err != nil {
				continue
			}
			to, err := hexToBase58(m["to_address"].(string))
			if err != nil {
				continue
			}
			td := &TxData{
				amount:    int64(m["amount"].(float64)),
				assetName: m["asset_name"].(string),
				from:      from,
				to:        to,
			}
			if t := jq.Reset().Find(pfx + ".raw_data.timestamp"); t == nil {
				continue
			} else {
				td.timestamp = int64(t.(float64))
				lastBlockTimestamp = td.timestamp
			}
			if txId := jq.Reset().Find(pfx + ".txID"); txId == nil {
				continue
			} else {
				if err := PersistTxWithTime(ds, cfg.Identity.PeerID, txId.(string), td.amount, td.from, td.to,
					StatusSuccess, walletpb.TransactionV1_ON_CHAIN, time.Unix(td.timestamp/1000, td.timestamp%1000)); err != nil {
					log.Error(err)
				}
			}
			tds = append(tds, td)
		}
	}
	if len(tds) > 0 {
		err := ds.Put(ctx, datastore.NewKey(fmt.Sprintf(curBlockTimestampKey, keys.Base58Address)),
			[]byte(strconv.FormatInt(lastBlockTimestamp, 10)))
		if err != nil {
			log.Debug(err)
		}
	}
	return tds, nil
}

func hexToBase58(h string) (string, error) {
	bs, err := hex.DecodeString(h)
	if err != nil {
		return "", err
	}
	rs, err := crypto.Encode58Check(bs)
	if err != nil {
		return "", err
	}
	return rs, nil
}

type TxData struct {
	amount    int64
	assetName string
	from      string
	to        string
	timestamp int64
}

func TransferBTT(ctx context.Context, n *core.IpfsNode, cfg *config.Config, privKey ic.PrivKey,
	from string, to string, amount int64) (*TronRet, error) {
	return TransferBTTWithMemo(ctx, n, cfg, privKey, from, to, amount, "")
}

func TransferBTTWithMemo(ctx context.Context, n *core.IpfsNode, cfg *config.Config, privKey ic.PrivKey,
	from string, to string, amount int64, memo string) (*TronRet, error) {
	var err error
	if privKey == nil {
		privKey, err = crypto.ToPrivKey(cfg.Identity.PrivKey)
		if err != nil {
			return nil, err
		}
	}
	if from == "" {
		keys, err := crypto.FromIcPrivateKey(privKey)
		if err != nil {
			return nil, err
		}
		from = keys.HexAddress
	}
	tx, err := PrepareTx(ctx, cfg, from, to, amount, memo)
	if err != nil {
		return nil, err
	}
	txId := ""
	if tx.Txid != nil {
		txId = hex.EncodeToString(tx.Txid)
	}
	raw, err := privKey.Raw()
	if err != nil {
		return nil, err
	}
	ecdsa, err := crypto.HexToECDSA(hex.EncodeToString(raw))
	if err != nil {
		return nil, err
	}
	bs, err := proto.Marshal(tx.Transaction.RawData)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.EcdsaSign(ecdsa, bs)
	if err != nil {
		return nil, err
	}
	rawBytes, err := proto.Marshal(tx.Transaction.RawData)
	if err != nil {
		return nil, err
	}
	err = SendRawTransaction(ctx, cfg.Services.FullnodeDomain, rawBytes, sig)
	if err != nil {
		return nil, err
	}
	err = PersistTx(n.Repo.Datastore(), n.Identity.String(), txId, amount,
		BttWallet, to, StatusPending, walletpb.TransactionV1_ON_CHAIN)
	if err != nil {
		return nil, err
	}
	go func() {
		// confirmed after 19 * 3 second/block
		time.Sleep(1 * time.Minute)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		status, err := GetStatus(ctx, cfg.Services.SolidityDomain, txId)
		if err != nil {
			log.Error(err)
			return
		}
		err = PersistTx(n.Repo.Datastore(), n.Identity.String(), txId, amount,
			BttWallet, to, status, walletpb.TransactionV1_ON_CHAIN)
		if err != nil {
			log.Error(err)
			return
		}
	}()
	return &TronRet{
		Message: string(tx.Result.Message),
		Result:  tx.Result.Result,
		Code:    tx.Result.Code.String(),
		TxId:    hex.EncodeToString(tx.Txid),
	}, nil
}

// base58/hex -> hex
func toHex(address string) (string, error) {
	if strings.HasPrefix(address, "T") {
		bytes, err := base58.Decode(address)
		if err != nil {
			return "", err
		}
		if len(bytes) <= 4 {
			return "", errors.New("invalid address")
		}
		address = hexutils.BytesToHex(bytes[:len(bytes)-4])
	}
	return address, nil
}

type TronRet struct {
	Message string
	Result  bool
	Code    string
	TxId    string
}

func PrepareTx(ctx context.Context, cfg *config.Config, from string, to string, amount int64, memo string) (*tronPb.TransactionExtention, error) {
	var (
		tx  *tronPb.TransactionExtention
		err error
	)
	from, err = toHex(from)
	if err != nil {
		return nil, err
	}
	to, err = toHex(to)
	if err != nil {
		return nil, err
	}
	oa, err := hex.DecodeString(from)
	if err != nil {
		return nil, err
	}
	ta, err := hex.DecodeString(to)
	if err != nil {
		return nil, err
	}
	tokenId := TokenId
	if strings.Contains(cfg.Services.EscrowDomain, "dev") ||
		strings.Contains(cfg.Services.EscrowDomain, "staging") {
		tokenId = TokenIdDev
	}
	err = grpc.WalletClient(cfg.Services.FullnodeDomain).WithContext(ctx, func(ctx context.Context, client tronPb.WalletClient) error {
		tx, err = client.TransferAsset2(ctx, &protocol_core.TransferAssetContract{
			AssetName:    []byte(tokenId),
			OwnerAddress: oa,
			ToAddress:    ta,
			Amount:       amount,
		})
		if err != nil {
			return err
		}
		if !tx.Result.Result {
			return errors.New(string(tx.Result.Message))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	tx.Transaction.RawData.Data = []byte(memo)
	bs, err := proto.Marshal(tx.Transaction.RawData)
	if err != nil {
		return nil, err
	}
	hashed := sha256.Sum256(bs)
	tx.Txid = hashed[:]
	return tx, nil
}

func SendRawTransaction(ctx context.Context, url string, raw []byte, sig []byte) error {
	rawMsg := &protocol_core.TransactionRaw{}
	err := proto.Unmarshal(raw, rawMsg)
	if err != nil {
		return err
	}
	tx := &protocol_core.Transaction{
		RawData:   rawMsg,
		Signature: [][]byte{sig},
	}
	return grpc.WalletClient(url).WithContext(ctx, func(ctx context.Context, client tronPb.WalletClient) error {
		_, err = client.BroadcastTransaction(ctx, tx)
		return err
	})
}

func GetStatus(ctx context.Context, url string, txId string) (string, error) {
	txIdBytes, err := hex.DecodeString(txId)
	if err != nil {
		return "", err
	}
	var info *protocol_core.Transaction
	err = grpc.SolidityClient(url).WithContext(ctx, func(ctx context.Context, client tronPb.WalletSolidityClient) error {
		info, err = client.GetTransactionById(ctx, &tronPb.BytesMessage{Value: txIdBytes})
		return err
	})
	if err != nil {
		return StatusFailed, err
	}
	status := StatusPending
	if info != nil && info.Ret != nil && len(info.Ret) > 0 &&
		info.Ret[0].ContractRet == protocol_core.Transaction_Result_SUCCESS {
		status = StatusSuccess
	} else {
		status = StatusFailed
	}
	return status, nil
}

func GetBalanceByWalletAddress(ctx context.Context, solidityUrl string, walletAddress string) (tokenMap map[string]int64, err error) {
	address, err := toHex(walletAddress)
	if err != nil {
		return nil, err
	}

	oa, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	err = grpc.SolidityClient(solidityUrl).WithContext(ctx, func(ctx context.Context, client tronPb.WalletSolidityClient) error {
		account := &protocol_core.Account{Address: oa}
		myAccount, err := client.GetAccount(ctx, account)
		if err != nil {
			return err
		}
		tokenMap = myAccount.GetAssetV2()
		if tokenMap == nil {
			tokenMap = make(map[string]int64)
		}
		tokenMap["TRX"] = myAccount.Balance
		return nil
	})

	return tokenMap, err
}
