package accesskey

import (
	"github.com/bittorrent/go-btfs/s3/handlers"
)

var instance handlers.AccessKeyer

func InitInstance(storer handlers.StateStorer, options ...Option) {
	instance = NewAccessKey(storer, options...)
}

func GetInstance() handlers.AccessKeyer {
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
