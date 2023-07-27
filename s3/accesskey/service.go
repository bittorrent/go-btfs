package accesskey

import (
	"github.com/bittorrent/go-btfs/transaction/storage"
)

var svc Service

func InitService(config *Config, store storage.StateStorer) {
	svc = newService(config, store)
}

func Generate() (ack *AccessKey, err error) {
	return svc.Generate()
}

func Enable(key string) (err error) {
	return svc.Enable(key)
}

func Disable(key string) (err error) {
	return svc.Disable(key)
}

func Reset(key string) (err error) {
	return svc.Reset(key)
}

func Delete(key string) (err error) {
	return svc.Delete(key)
}

func Get(key string) (ack *AccessKey, err error) {
	return svc.Get(key)
}

func List() (list []*AccessKey, err error) {
	return svc.List()
}
