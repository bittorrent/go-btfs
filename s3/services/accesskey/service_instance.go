package accesskey

import (
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/services"
	"sync"
)

var service *Service

var once sync.Once

func InitService(providers providers.Providerser, options ...Option) {
	once.Do(func() {
		service = NewService(providers, options...)
	})
}

func GetService() *Service {
	return service
}

func Generate() (ack *services.AccessKey, err error) {
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

func Get(key string) (record *services.AccessKey, err error) {
	return service.Get(key)
}

func List() (list []*services.AccessKey, err error) {
	return service.List()
}
