package accesskey

import (
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/services"
	"sync"
)

var instance handlers.AccessKeyService

var once sync.Once

func InitInstance(providers services.Providerser, options ...Option) {
	once.Do(func() {
		instance = NewAccessKey(providers, options...)
	})
}

func GetInstance() handlers.AccessKeyService {
	return instance
}

func Generate() (record *handlers.AccessKeyRecord, err error) {
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

func Get(key string) (record *handlers.AccessKeyRecord, err error) {
	return instance.Get(key)
}

func List() (list []*handlers.AccessKeyRecord, err error) {
	return instance.List()
}
