package services

type StateStoreIterFunc func(key, value []byte) (stop bool, err error)
