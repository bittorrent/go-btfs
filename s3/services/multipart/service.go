package multipart

import (
	"github.com/bittorrent/go-btfs/s3/services"
	"io"
)

var _ services.MultipartService = (*Service)(nil)

type Service struct {
}

func NewService(options ...Option) (svc *Service) {
	svc = &Service{}
	for _, option := range options {
		option(svc)
	}
	return
}

func (svc *Service) multiReader() io.Reader {
	var (
		r1 io.Reader
		r2 io.Reader
		r3 io.Reader
	)

	return io.MultiReader(r1, r2, r3)
}
