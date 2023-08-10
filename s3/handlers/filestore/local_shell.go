package filestore

import (
	shell "github.com/bittorrent/go-btfs-api"
	"github.com/bittorrent/go-btfs/s3/handlers"
)

var _ handlers.FileStorer = (*LocalShell)(nil)

type LocalShell struct {
	*shell.Shell
}

func NewFileStore() *LocalShell {
	return &LocalShell{
		Shell: shell.NewLocalShell(),
	}
}
