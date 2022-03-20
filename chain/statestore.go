package chain

import (
	"path/filepath"

	"github.com/bittorrent/go-btfs/statestore/leveldb"
	"github.com/bittorrent/go-btfs/statestore/mock"
	"github.com/bittorrent/go-btfs/transaction/storage"
)

var store storage.StateStorer

func InitStateStore(dataDir string) (ret storage.StateStorer, err error) {
	if dataDir == "" {
		ret = mock.NewStateStore()
		log.Warn("using in-mem state store, no node state will be persisted")
		return ret, nil
	}

	store, err = leveldb.NewStateStore(filepath.Join(dataDir, "statestore"))
	return store, err
}
