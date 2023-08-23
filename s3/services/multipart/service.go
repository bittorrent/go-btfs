package multipart

import (
	"github.com/bittorrent/go-btfs/s3/services"
	"io"
)

var _ services.MultipartService = (*service)(nil)

type service struct {
}

func NewService(options ...Option) Service {
	svc := &service{}
	for _, option := range options {
		option(svc)
	}
	return svc
}

func (svc *service) multiReader() io.Reader {
	var (
		r1 io.Reader
		r2 io.Reader
		r3 io.Reader
	)

	return io.MultiReader(r1, r2, r3)
}
