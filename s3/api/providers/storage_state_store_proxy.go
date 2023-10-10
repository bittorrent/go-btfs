package providers

import (
	"errors"
	"github.com/bittorrent/go-btfs/transaction/storage"
)

var _ StateStorer = (*StorageStateStoreProxy)(nil)

type StorageStateStoreProxy struct {
	to storage.StateStorer
}

func NewStorageStateStoreProxy(to storage.StateStorer) *StorageStateStoreProxy {
	return &StorageStateStoreProxy{
		to: to,
	}
}

func (s *StorageStateStoreProxy) Put(key string, val interface{}) (err error) {
	return s.to.Put(key, val)
}

func (s *StorageStateStoreProxy) Get(key string, i interface{}) (err error) {
	err = s.to.Get(key, i)
	if errors.Is(err, storage.ErrNotFound) {
		err = ErrStateStoreNotFound
	}
	return
}

func (s *StorageStateStoreProxy) Delete(key string) (err error) {
	err = s.to.Delete(key)
	if errors.Is(err, storage.ErrNotFound) {
		err = ErrStateStoreNotFound
	}
	return
}

func (s *StorageStateStoreProxy) Iterate(prefix string, iterFunc StateStoreIterFunc) (err error) {
	return s.to.Iterate(prefix, storage.StateIterFunc(iterFunc))
}
