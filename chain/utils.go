package chain

import (
	"encoding/base64"
	"errors"
	"fmt"
	config "github.com/bittorrent/go-btfs-config"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs-common/crypto"
	onlinePb "github.com/bittorrent/go-btfs-common/protos/online"
	oldcmds "github.com/bittorrent/go-btfs/commands"
	"github.com/bittorrent/go-btfs/core/commands/storage/path"
	"github.com/bittorrent/go-btfs/settlement/swap/vault"
	cpt "github.com/bittorrent/go-btfs/transaction/crypto"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/ethereum/go-ethereum/common"
)

// after btfs init
func GetBttcNonDaemon(env cmds.Environment) (defaultAddr string, _err error) {
	cctx := env.(*oldcmds.Context)
	_, b := os.LookupEnv(path.BtfsPathKey)
	if !b {
		c := cctx.ConfigRoot
		if bs, err := ioutil.ReadFile(path.PropertiesFileName); err == nil && len(bs) > 0 {
			c = string(bs)
		}
		cctx.ConfigRoot = c
	}

	cfg, err := cctx.GetConfig()
	if err != nil {
		return defaultAddr, err
	}

	// decode from string
	pkbytesOri, err := base64.StdEncoding.DecodeString(cfg.Identity.PrivKey)
	if err != nil {
		return defaultAddr, err
	}

	//new singer
	pk := cpt.Secp256k1PrivateKeyFromBytes(pkbytesOri[4:])
	singer := cpt.NewDefaultSigner(pk)

	address0x, err := singer.EthereumAddress()
	if err != nil {
		return defaultAddr, err
	}
	return address0x.Hex(), nil
}

func GetBttcByKey(privKey string) (defaultAddr string, _err error) {
	// decode from string
	pkbytesOri, err := base64.StdEncoding.DecodeString(privKey)
	if err != nil {
		return defaultAddr, err
	}

	//new singer
	pk := cpt.Secp256k1PrivateKeyFromBytes(pkbytesOri[4:])
	singer := cpt.NewDefaultSigner(pk)

	address0x, err := singer.EthereumAddress()
	if err != nil {
		return defaultAddr, err
	}
	return address0x.Hex(), nil
}

// after btfs init
func GetVaultNonDaemon(env cmds.Environment) (defaultAddr string, err error) {
	cctx := env.(*oldcmds.Context)
	_, b := os.LookupEnv(path.BtfsPathKey)
	if !b {
		c := cctx.ConfigRoot
		if bs, err := ioutil.ReadFile(path.PropertiesFileName); err == nil && len(bs) > 0 {
			c = string(bs)
		}
		cctx.ConfigRoot = c
	}

	// btfs id cmd, not node process
	statestore, err := InitStateStore(cctx.ConfigRoot)
	if err != nil {
		return defaultAddr, err
	}

	var vaultAddress common.Address
	err = statestore.Get(vault.VaultKey, &vaultAddress)
	if err != nil {
		if err == storage.ErrNotFound {
			return defaultAddr, nil
		}
		return defaultAddr, err
	}

	return vaultAddress.Hex(), nil
}

// after btfs init
func GetWalletImportPrvKey(env cmds.Environment) (string, error) {
	cctx := env.(*oldcmds.Context)
	cfg, err := cctx.GetConfig()
	if err != nil {
		return "", err
	}
	privKey, err := crypto.ToPrivKey(cfg.Identity.PrivKey)
	if err != nil {
		return "", err
	}
	keys, err := crypto.FromIcPrivateKey(privKey)
	if err != nil {
		return "", err
	}

	return keys.HexPrivateKey, nil
}

var chainIdKey = "ChainIdKey"
var DefaultStoreChainId = int64(-1)

// add chain id into leveldb
func StoreChainIdToDisk(ChainId int64, stateStorer storage.StateStorer) error {
	err := stateStorer.Put(chainIdKey, ChainId)
	if err != nil {
		return err
	}
	return nil
}

// get chain id from leveldb
func GetChainIdFromDisk(stateStorer storage.StateStorer) (int64, error) {
	var storeChainId int64
	err := stateStorer.Get(chainIdKey, &storeChainId)
	if err != nil {
		if err == storage.ErrNotFound {
			return DefaultStoreChainId, nil
		}
		return 0, err
	}
	return storeChainId, nil
}

func StoreChainIdIfNotExists(chainID int64, statestore storage.StateStorer) error {
	storeChainid, err := GetChainIdFromDisk(statestore)
	if err != nil {
		return err
	}

	if storeChainid <= 0 {
		err = StoreChainIdToDisk(chainID, statestore)
		if err != nil {
			fmt.Println("StoreChainIdIfNotExists: init StoreChainId err: ", err)
			return err
		}
	}

	return nil
}

// GetReportStatus from leveldb
var keyReportStatus = "keyReportStatus"

type ReportStatusInfo struct {
	ReportStatusSeconds int64
	LastReportTime      time.Time
}

func GetReportStatus() (*ReportStatusInfo, error) {
	if StateStore == nil {
		return nil, errors.New("please start btfs node, at first! ")
	}

	var reportStatusInfo ReportStatusInfo
	err := StateStore.Get(keyReportStatus, &reportStatusInfo)
	if err != nil {
		if err == storage.ErrNotFound {
			reportStatusInfo = ReportStatusInfo{ReportStatusSeconds: int64(rand.Intn(100000000) % 86400), LastReportTime: time.Time{}}
			err := StateStore.Put(keyReportStatus, reportStatusInfo)
			if err != nil {
				fmt.Println("StoreChainIdIfNotExists: init StoreChainId err: ", err)
				return nil, err
			}
		}
		return nil, err
	}
	return &reportStatusInfo, nil
}

func SetReportStatusOK() (*ReportStatusInfo, error) {
	var reportStatusInfo ReportStatusInfo
	err := StateStore.Get(keyReportStatus, &reportStatusInfo)
	if err != nil {
		return nil, err
	}
	reportStatusInfo.LastReportTime = time.Now()
	err = StateStore.Put(keyReportStatus, reportStatusInfo)
	if err != nil {
		return nil, err
	}
	//fmt.Println("... ReportStatus, SetReportStatus: ok! ")
	return &reportStatusInfo, nil
}

// GetReportStatus from leveldb
var keyReportStatusList = "keyReportStatusList"

type LevelDbReportStatusInfo struct {
	Peer           string    `json:"peer"`
	BttcAddress    string    `json:"bttc_addr"`
	StatusContract string    `json:"status_contract"`
	Nonce          uint32    `json:"nonce"`
	TxHash         string    `json:"tx_hash"`
	GasSpend       string    `json:"gas_spend"`
	ReportTime     time.Time `json:"report_time"`
	IncreaseNonce  uint32    `json:"increase_nonce"`
}

// SetReportStatusListOK store tx list
func SetReportStatusListOK(r *LevelDbReportStatusInfo) ([]*LevelDbReportStatusInfo, error) {
	if StateStore == nil {
		return nil, errors.New("please start btfs node, at first! ")
	}

	init := false

	rList := make([]*LevelDbReportStatusInfo, 0)
	err := StateStore.Get(keyReportStatusList, &rList)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			init = true
			// continue
		} else {
			return nil, err
		}
	}

	if init {
		r.IncreaseNonce = r.Nonce
	} else {
		r.IncreaseNonce = r.Nonce - rList[len(rList)-1].Nonce
	}

	rList = append(rList, r)
	err = StateStore.Put(keyReportStatusList, rList)
	if err != nil {
		return nil, err
	}
	//fmt.Println("... ReportStatus, SetReportStatusListOK: ok! rList = ", rList)
	return rList, nil
}

// GetReportStatusListOK store tx list
func GetReportStatusListOK() ([]*LevelDbReportStatusInfo, error) {
	if StateStore == nil {
		return nil, errors.New("please start btfs node, at first! ")
	}

	rList := make([]*LevelDbReportStatusInfo, 0)
	err := StateStore.Get(keyReportStatusList, &rList)
	if err != nil {
		if err.Error() == "storage: not found" {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return rList, nil
}

// GetLastOnline from leveldb
var keyLastOnline = "keyLastOnline"

type LastOnlineInfo struct {
	LastSignedInfo onlinePb.SignedInfo
	LastSignature  string
	LastTime       time.Time
}

type LastOnlineInfoRet struct {
	LastSignedInfo onlinePb.SignedInfo `json:"last_signed_info"`
	LastSignature  string              `json:"last_signature"`
	LastTime       time.Time           `json:"last_time""`
}

func GetLastOnline() (*LastOnlineInfo, error) {
	if StateStore == nil {
		return nil, errors.New("please start btfs node, at first! ")
	}

	var lastOnlineInfo LastOnlineInfo
	err := StateStore.Get(keyLastOnline, &lastOnlineInfo)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &lastOnlineInfo, nil
}
func StoreOnline(lastOnlineInfo *LastOnlineInfo) error {
	err := StateStore.Put(keyLastOnline, *lastOnlineInfo)
	if err != nil {
		fmt.Println("StoreOnline: init StoreChainId err: ", err)
		return err
	}

	return nil
}

func GetOnlineServer(chainId int64) string {
	if chainId == 199 {
		return config.DefaultServicesConfig().OnlineServerDomain
	} else {
		return config.DefaultServicesConfigTestnet().OnlineServerDomain
	}
}

// GetReportOnlineDailyLast from leveldb
var keyReportOnlineLastTimeDaily = "keyReportOnlineLastTimeDaily"

type ReportOnlineLastTimeDaily struct {
	EveryDaySeconds int64
	LastReportTime  time.Time
}

func GetReportOnlineLastTimeDaily() (*ReportOnlineLastTimeDaily, error) {
	if StateStore == nil {
		return nil, errors.New("please start btfs node, at first! ")
	}

	var info ReportOnlineLastTimeDaily
	err := StateStore.Get(keyReportOnlineLastTimeDaily, &info)
	if err != nil {
		if err == storage.ErrNotFound {
			v := ReportOnlineLastTimeDaily{EveryDaySeconds: int64(rand.Intn(100000000) % 86400), LastReportTime: time.Time{}}
			err := StateStore.Put(keyReportOnlineLastTimeDaily, v)
			if err != nil {
				fmt.Println("GetReportOnlineLastTimeDaily: init leveldb err: ", err)
				return nil, err
			}
			return &v, nil
		} else {
			return nil, err
		}
	}
	return &info, nil
}
func SetReportOnlineLastTimeDailyOK() (*ReportOnlineLastTimeDaily, error) {
	var info ReportOnlineLastTimeDaily
	err := StateStore.Get(keyReportOnlineLastTimeDaily, &info)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			// continue
		} else {
			return nil, err
		}
	}
	info.LastReportTime = time.Now()
	err = StateStore.Put(keyReportOnlineLastTimeDaily, info)
	if err != nil {
		return nil, err
	}
	//fmt.Println("... SetReportOnlineLastTimeDailyOK: ok! ")
	return &info, nil
}

// GetReportStatus from leveldb
var keyReportOnlineListDaily = "keyReportOnlineListDaily"

// SetReportOnlineListDailyOK store online daily list
func SetReportOnlineListDailyOK(r *LastOnlineInfo) ([]*LastOnlineInfo, error) {
	if StateStore == nil {
		return nil, errors.New("please start btfs node, at first! ")
	}

	rList := make([]*LastOnlineInfo, 0)
	err := StateStore.Get(keyReportOnlineListDaily, &rList)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			// continue
		} else {
			return nil, err
		}
	}

	rList = append(rList, r)
	err = StateStore.Put(keyReportOnlineListDaily, rList)
	if err != nil {
		return nil, err
	}
	fmt.Println("SetReportOnlineListDailyOK: ok! nonce =", r.LastSignedInfo.Nonce)
	return rList, nil
}

// GetReportOnlineListDailyOK store tx list
func GetReportOnlineListDailyOK() ([]*LastOnlineInfo, error) {
	if StateStore == nil {
		return nil, errors.New("please start btfs node, at first! ")
	}

	rList := make([]*LastOnlineInfo, 0)
	err := StateStore.Get(keyReportOnlineListDaily, &rList)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return rList, nil
}
