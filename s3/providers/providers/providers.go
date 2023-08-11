package providers

import "github.com/bittorrent/go-btfs/s3/providers"

var _ providers.Providerser = (*Providers)(nil)

type Providers struct {
	statestore providers.StateStorer
	filestore  providers.FileStorer
}

func NewProviders(statestore providers.StateStorer, filestore providers.FileStorer, options ...Option) *Providers {
	p := &Providers{
		statestore: statestore,
		filestore:  filestore,
	}
	for _, option := range options {
		option(p)
	}
	return p
}

func (p *Providers) GetStateStore() providers.StateStorer {
	return p.statestore
}

func (p *Providers) GetFileStore() providers.FileStorer {
	return p.filestore
}
