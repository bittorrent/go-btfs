package providers

var _ Providerser = (*Providers)(nil)

type Providers struct {
	stateStore StateStorer
	fileStore  FileStorer
}

func NewProviders(stateStore StateStorer, fileStore FileStorer, options ...Option) (providers *Providers) {
	providers = &Providers{
		stateStore: stateStore,
		fileStore:  fileStore,
	}
	for _, option := range options {
		option(providers)
	}
	return
}

func (p *Providers) StateStore() StateStorer {
	return p.stateStore
}

func (p *Providers) FileStore() FileStorer {
	return p.fileStore
}
