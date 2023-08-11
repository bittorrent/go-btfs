package filestore

import (
	shell "github.com/bittorrent/go-btfs-api"
	"github.com/bittorrent/go-btfs/s3/providers"
)

var _ providers.FileStorer = (*LocalShell)(nil)

type LocalShell struct {
	*shell.Shell
}

func NewLocalShell() *LocalShell {
	return &LocalShell{
		Shell: shell.NewLocalShell(),
	}
}
