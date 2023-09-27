package accesskey

import (
	"github.com/bittorrent/go-btfs/s3/api/providers"
	"sync"
)

var svcInstance Service

var once sync.Once

func InitService(providers providers.Providerser, options ...Option) {
	once.Do(func() {
		svcInstance = NewService(providers, options...)
	})
}

func GetServiceInstance() Service {
	return svcInstance
}

func Generate() (ack *AccessKey, err error) {
	return svcInstance.Generate()
}

func Enable(key string) (err error) {
	return svcInstance.Enable(key)
}

func Disable(key string) (err error) {
	return svcInstance.Disable(key)
}

func Reset(key string) (err error) {
	return svcInstance.Reset(key)
}

func Delete(key string) (err error) {
	return svcInstance.Delete(key)
}

func Get(key string) (record *AccessKey, err error) {
	return svcInstance.Get(key)
}

func List() (list []*AccessKey, err error) {
	return svcInstance.List()
}
