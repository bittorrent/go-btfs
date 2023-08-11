package handlers

import (
	"errors"
	"net/http"
	"time"
)

type AccessKeyService interface {
	Generate() (record *AccessKeyRecord, err error)
	Enable(key string) (err error)
	Disable(key string) (err error)
	Reset(key string) (err error)
	Delete(key string) (err error)
	Get(key string) (record *AccessKeyRecord, err error)
	List() (list []*AccessKeyRecord, err error)
}

type AccessKeyRecord struct {
	Key       string    `json:"key"`
	Secret    string    `json:"secret"`
	Enable    bool      `json:"enable"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrAccessKeyIsNotFound = errors.New("access-key is not found")

type BucketService interface {
}

type ObjectService interface {
}

type MultipartService interface {
}

type SignService interface {
	Verify(r *http.Request) (err error)
}

var (
	ErrSignOutdated      = errors.New("sign is outdated")
	ErrSignKeyIsNotFound = errors.New("key is not found")
)
