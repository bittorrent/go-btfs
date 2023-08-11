package providers

import (
	"errors"
	"io"
)

type Providerser interface {
	GetFileStore() FileStorer
	GetStateStore() StateStorer
}

type FileStorer interface {
	AddWithOpts(r io.Reader, pin bool, rawLeaves bool) (hash string, err error)
	Remove(hash string) (removed bool)
	Cat(path string) (readCloser io.ReadCloser, err error)
	Unpin(path string) (err error)
}

type StateStorer interface {
	Get(key string, i interface{}) (err error)
	Put(key string, i interface{}) (err error)
	Delete(key string) (err error)
	Iterate(prefix string, iterFunc StateStoreIterFunc) (err error)
}

type StateStoreIterFunc func(key, value []byte) (stop bool, err error)

var ErrStateStoreNotFound = errors.New("not found")
