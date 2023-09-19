package accesskey

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type Service interface {
	Generate() (record *AccessKey, err error)
	Enable(key string) (err error)
	Disable(key string) (err error)
	Reset(key string) (err error)
	Delete(key string) (err error)
	Get(key string) (ack *AccessKey, err error)
	List() (list []*AccessKey, err error)
}

type AccessKey struct {
	Key       string    `json:"key"`
	Secret    string    `json:"secret"`
	Enable    bool      `json:"enable"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
