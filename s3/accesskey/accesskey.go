package accesskey

import (
	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/bittorrent/go-mfs"
	"path"
	"time"
)

var _ Service = &service{}

const (
	KeyLen      = 8
	SecretLen   = 32
	RootPrefix  = "s3_buckets"
	StorePrefix = "s3-access-key-"
)

type service struct {
	store storage.StateStorer
}

func (svc *service) Generate() (ak *AccessKey, err error) {
	ak := &AccessKey{
		Key:       GetRandStr(KeyLen),
		Secret:    GetRandStr(SecretLen),
		Root:      path.Join("/", RootPrefix, GetRandStr(8)),
		Enable:    true,
		CreatedAt: time.Now(),
	}

	// create root dir
	mfs.Mkdir()

	// store accessKey
	err = svc.store.Put(StorePrefix+ak.Key, ak)
	return
}

func (svc *service) Get(key string) (ak *AccessKey, err error) {
	return
}

func (svc *service) Disable(key string) (err error) {
	return
}

func (svc *service) Reset(key string) (err error) {
	return
}

func (svc *service) Delete(key string) (err error) {
	return
}

func (svc *service) List() (aks []*AccessKey, err error) {
	return
}
