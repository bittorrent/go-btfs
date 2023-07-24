package s3

import (
	cfg "github.com/bittorrent/go-btfs-config"
)

type Request struct {
}

type Service interface {
	Start(config *cfg.S3CompatibleAPI) (err error)
	Stop() (err error)
}
