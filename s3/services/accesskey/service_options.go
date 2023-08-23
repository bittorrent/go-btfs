package accesskey

type Option func(svc *service)

func WithSecretLength(length int) Option {
	return func(svc *service) {
		svc.secretLength = length
	}
}

func WithStoreKeyPrefix(prefix string) Option {
	return func(svc *service) {
		svc.storeKeyPrefix = prefix
	}
}
