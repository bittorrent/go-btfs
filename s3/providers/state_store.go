package providers

import (
	"errors"
	"github.com/bittorrent/go-btfs/transaction/storage"
)

var _ StateStorer = (*StateStore)(nil)

type StateStore struct {
	proxy storage.StateStorer
}

func NewStorageStateStoreProxy(proxy storage.StateStorer) *StateStore {
	return &StateStore{
		proxy: proxy,
	}
}

func (s *StateStore) Put(key string, val interface{}) (err error) {
	return s.proxy.Put(key, val)
}

func (s *StateStore) Get(key string, i interface{}) (err error) {
	err = s.proxy.Get(key, i)
	if errors.Is(err, storage.ErrNotFound) {
		err = ErrStateStoreNotFound
	}
	return
}

func (s *StateStore) Delete(key string) (err error) {
	err = s.proxy.Delete(key)
	if errors.Is(err, storage.ErrNotFound) {
		err = ErrStateStoreNotFound
	}
	return
}

func (s *StateStore) Iterate(prefix string, iterFunc StateStoreIterFunc) (err error) {
	return s.proxy.Iterate(prefix, storage.StateIterFunc(iterFunc))
}
