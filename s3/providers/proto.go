package providers

import (
	"errors"
	"io"
)

var (
	ErrStateStoreNotFound = errors.New("not found in state store")
	ErrFileStoreNotFound  = errors.New("not found in file store")
)

type Providerser interface {
	FileStore() FileStorer
	StateStore() StateStorer
}

type FileStorer interface {
	Store(r io.Reader) (id string, err error)
	Remove(id string) (err error)
	Cat(id string) (readCloser io.ReadCloser, err error)
}

type StateStorer interface {
	Get(key string, i interface{}) (err error)
	Put(key string, i interface{}) (err error)
	Delete(key string) (err error)
	Iterate(prefix string, iterFunc StateStoreIterFunc) (err error)
}

type StateStoreIterFunc func(key, value []byte) (stop bool, err error)
