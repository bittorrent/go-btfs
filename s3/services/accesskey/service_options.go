package accesskey

type Option func(svc *Service)

func WithSecretLength(length int) Option {
	return func(svc *Service) {
		svc.secretLength = length
	}
}

func WithStoreKeyPrefix(prefix string) Option {
	return func(svc *Service) {
		svc.storeKeyPrefix = prefix
	}
}
