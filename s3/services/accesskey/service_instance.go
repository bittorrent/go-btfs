package accesskey

import (
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/services"
	"sync"
)

var service *Service

var once sync.Once

func InitService(providers services.Providerser, options ...Option) {
	once.Do(func() {
		service = NewService(providers, options...)
	})
}

func GetService() *Service {
	return service
}

func Generate() (record *handlers.AccessKeyRecord, err error) {
	return service.Generate()
}

func Enable(key string) (err error) {
	return service.Enable(key)
}

func Disable(key string) (err error) {
	return service.Disable(key)
}

func Reset(key string) (err error) {
	return service.Reset(key)
}

func Delete(key string) (err error) {
	return service.Delete(key)
}

func Get(key string) (record *handlers.AccessKeyRecord, err error) {
	return service.Get(key)
}

func List() (list []*handlers.AccessKeyRecord, err error) {
	return service.List()
}
