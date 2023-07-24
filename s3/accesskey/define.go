package accesskey

import "time"

type AccessKey struct {
	Key       string    `json:"key"`
	Secret    string    `json:"secret"`
	Root      string    `json:"root"`
	Enable    bool      `json:"enable"`
	CreatedAt time.Time `json:"created_at"`
}

type Service interface {
	Generate() (*AccessKey, error)
	Get(key string) (*AccessKey, error)
	Disable(key string) error
	Reset(key string) error
	Delete(key string) error
	List() ([]*AccessKey, error)
}
