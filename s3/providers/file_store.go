package providers

import (
	shell "github.com/bittorrent/go-btfs-api"
)

var _ FileStorer = (*FileStore)(nil)

type FileStore struct {
	*shell.Shell
}

func NewFileStore() *FileStore {
	return &FileStore{
		Shell: shell.NewLocalShell(),
	}
}
