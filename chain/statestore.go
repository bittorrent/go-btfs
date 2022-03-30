package chain

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/bittorrent/go-btfs/statestore/leveldb"
	"github.com/bittorrent/go-btfs/statestore/mock"
	"github.com/bittorrent/go-btfs/transaction/storage"
)

func InitStateStore(dataDir string) (ret storage.StateStorer, err error) {
	if dataDir == "" {
		ret = mock.NewStateStore()
		log.Warn("using in-mem state store, no node state will be persisted")
		return ret, nil
	}
	return leveldb.NewStateStore(GetStateStorePath(dataDir))
}

func GetStateStorePath(dataDir string) string {
	return filepath.Join(dataDir, "statestore")
}

func BackUpStateStore(dataDir string, suffix string) error {
	if suffix == "" {
		suffix = fmt.Sprintf("backup%d", rand.Intn(100))
	}

	stateStorePath := GetStateStorePath(dataDir)
	finfo, err := os.Stat(stateStorePath)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	if !finfo.IsDir() {
		return errors.New("statestore is not a folder")
	}

	backUpPath := fmt.Sprintf("%s.%s", stateStorePath, suffix)
	err = os.Rename(stateStorePath, backUpPath)
	if err != nil {
		return err
	}
	fmt.Printf("backup statestore folder successfully to %s\n", backUpPath)
	return nil
}
