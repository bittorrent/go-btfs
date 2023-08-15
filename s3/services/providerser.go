package services

import (
	"context"
	"github.com/bittorrent/go-btfs/transaction/storage"
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
	ReadAllChan(ctx context.Context, prefix string, seekKey string) (<-chan *storage.Entry, error)
}
