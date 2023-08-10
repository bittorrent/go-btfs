package handlers

import (
	"errors"
	"io"
	"time"
)

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

type AccessKeyer interface {
	Generate() (record *AccessKeyRecord, err error)
	Enable(key string) (err error)
	Disable(key string) (err error)
	Reset(key string) (err error)
	Delete(key string) (err error)
	Get(key string) (record *AccessKeyRecord, err error)
	List() (list []*AccessKeyRecord, err error)
}

type AccessKeyRecord struct {
	Key       string    `json:"key"`
	Secret    string    `json:"secret"`
	Enable    bool      `json:"enable"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrAccessKeyIsNotFound = errors.New("access-key is not found")
