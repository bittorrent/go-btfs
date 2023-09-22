package vault

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

var RestartFixCashOutStatusLock bool = true
var RestartWaitCashOutOnlineTime time.Duration = 30 //seconds

// CashOutStatus from leveldb
var prefixKeyCashOutStatusStore = "keyCashOutStatusStore" // + txHash.
type CashOutStatusStoreInfo struct {
	Token            common.Address
	Vault            common.Address
	Beneficiary      common.Address
	CumulativePayout *big.Int
	TxHash           string
}

func getkeyCashOutStatusStore(txHash string) string {
	return fmt.Sprintf("%s-%s", prefixKeyCashOutStatusStore, txHash)
}

// AddCashOutStatusStore .
func (s *cashoutService) AddCashOutStatusStore(info CashOutStatusStoreInfo) (err error) {
	if s.store == nil {
		return errors.New("please start btfs node, at first! ")
	}

	err = s.store.Put(getkeyCashOutStatusStore(info.TxHash), info)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCashOutStatusStore .
func (s *cashoutService) DeleteCashOutStatusStore(info CashOutStatusStoreInfo) (err error) {
	if s.store == nil {
		return errors.New("please start btfs node, at first! ")
	}

	err = s.store.Delete(getkeyCashOutStatusStore(info.TxHash))
	if err != nil {
		if err.Error() == "storage: not found" {
			return nil
		} else {
			return err
		}
	}
	return
}

// GetCashOutStatusStore .
func (s *cashoutService) GetCashOutStatusStore(txHash string) (bl bool, err error) {
	if s.store == nil {
		return bl, errors.New("please start btfs node, at first! ")
	}

	var info CashOutStatusStoreInfo
	err = s.store.Get(getkeyCashOutStatusStore(txHash), &info)
	if err != nil {
		if err.Error() == "storage: not found" {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

// GetAllCashOutStatusStore .
func (s *cashoutService) GetAllCashOutStatusStore() (infoList []CashOutStatusStoreInfo, err error) {
	if s.store == nil {
		return nil, errors.New("please start btfs node, at first! ")
	}

	infoList = make([]CashOutStatusStoreInfo, 0)
	err = s.store.Iterate(prefixKeyCashOutStatusStore, func(key, val []byte) (stop bool, err error) {
		var info CashOutStatusStoreInfo
		err = s.store.Get(string(key), &info)
		if err != nil {
			return false, err
		}
		infoList = append(infoList, info)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return infoList, nil
}
