package accesskey

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type Config struct {
	SecretLength int
	StorePrefix  string
}

type AccessKey struct {
	Key       string    `json:"key"`
	Secret    string    `json:"secret"`
	Enable    bool      `json:"enable"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Bucket struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	ACL       string `json:"acl"`
	CID       string `json:"cid"`
	IsDeleted bool   `json:"is_deleted"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Service interface {
	Generate() (ack *AccessKey, err error)
	Enable(key string) (err error)
	Disable(key string) (err error)
	Reset(key string) (err error)
	Delete(key string) (err error)
	Get(key string) (ack *AccessKey, err error)
	List() (list []*AccessKey, err error)
}
