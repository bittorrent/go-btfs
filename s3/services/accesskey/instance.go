package accesskey

import (
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/services"
	"sync"
)

var instance services.AccessKeyService

var once sync.Once

func InitInstance(providers providers.Providerser, options ...Option) {
	once.Do(func() {
		instance = NewAccessKey(providers, options...)
	})
}

func GetInstance() services.AccessKeyService {
	return instance
}

func Generate() (record *services.AccessKeyRecord, err error) {
	return instance.Generate()
}

func Enable(key string) (err error) {
	return instance.Enable(key)
}

func Disable(key string) (err error) {
	return instance.Disable(key)
}

func Reset(key string) (err error) {
	return instance.Reset(key)
}

func Delete(key string) (err error) {
	return instance.Delete(key)
}

func Get(key string) (record *services.AccessKeyRecord, err error) {
	return instance.Get(key)
}

func List() (list []*services.AccessKeyRecord, err error) {
	return instance.List()
}
