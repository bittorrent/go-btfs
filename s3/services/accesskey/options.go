package accesskey

type Option func(ack *AccessKey)

func WithSecretLength(length int) Option {
	return func(ack *AccessKey) {
		ack.secretLength = length
	}
}

func WithStoreKeyPrefix(prefix string) Option {
	return func(ack *AccessKey) {
		ack.storeKeyPrefix = prefix
	}
}
