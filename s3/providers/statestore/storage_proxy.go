package statestore

import (
	"errors"
	"github.com/bittorrent/go-btfs/s3/services"
	"github.com/bittorrent/go-btfs/transaction/storage"
)

var _ services.StateStorer = (*StorageProxy)(nil)

type StorageProxy struct {
	proxy storage.StateStorer
}

func NewStorageStateStoreProxy(proxy storage.StateStorer) services.StateStorer {
	return &StorageProxy{
		proxy: proxy,
	}
}

func (s *StorageProxy) Put(key string, val interface{}) (err error) {
	return s.proxy.Put(key, val)
}

func (s *StorageProxy) Get(key string, i interface{}) (err error) {
	err = s.proxy.Get(key, i)
	if errors.Is(err, storage.ErrNotFound) {
		err = services.ErrStateStoreNotFound
	}
	return
}

func (s *StorageProxy) Delete(key string) (err error) {
	err = s.proxy.Delete(key)
	if errors.Is(err, storage.ErrNotFound) {
		err = services.ErrStateStoreNotFound
	}
	return
}

func (s *StorageProxy) Iterate(prefix string, iterFunc services.StateStoreIterFunc) (err error) {
	return s.proxy.Iterate(prefix, storage.StateIterFunc(iterFunc))
}
