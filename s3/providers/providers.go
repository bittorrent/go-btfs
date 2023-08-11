package providers

import (
	"github.com/bittorrent/go-btfs/s3/services"
)

var _ services.Providerser = (*Providers)(nil)

type Providers struct {
	statestore services.StateStorer
	filestore  services.FileStorer
}

func NewProviders(statestore services.StateStorer, filestore services.FileStorer, options ...Option) (providers *Providers) {
	providers = &Providers{
		statestore: statestore,
		filestore:  filestore,
	}
	for _, option := range options {
		option(providers)
	}
	return
}

func (p *Providers) GetStateStore() services.StateStorer {
	return p.statestore
}

func (p *Providers) GetFileStore() services.FileStorer {
	return p.filestore
}
