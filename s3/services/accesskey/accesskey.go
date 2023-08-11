package accesskey

import (
	"errors"
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/services"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/bittorrent/go-btfs/utils"
	"github.com/google/uuid"
	"sync"
	"time"
)

const (
	defaultSecretLength   = 32
	defaultStoreKeyPrefix = "access-keys:"
)

var _ handlers.AccessKeyService = (*AccessKey)(nil)

type AccessKey struct {
	providers      services.Providerser
	secretLength   int
	storeKeyPrefix string
	locks          sync.Map
}

func NewAccessKey(providers services.Providerser, options ...Option) (ack *AccessKey) {
	ack = &AccessKey{
		providers:      providers,
		secretLength:   defaultSecretLength,
		storeKeyPrefix: defaultStoreKeyPrefix,
		locks:          sync.Map{},
	}
	for _, option := range options {
		option(ack)
	}
	return ack
}

func (ack *AccessKey) Generate() (record *handlers.AccessKeyRecord, err error) {
	now := time.Now()
	record = &handlers.AccessKeyRecord{
		Key:       ack.newKey(),
		Secret:    ack.newSecret(),
		Enable:    true,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = ack.providers.GetStateStore().Put(ack.getStoreKey(record.Key), record)
	return
}

func (ack *AccessKey) Enable(key string) (err error) {
	enable := true
	err = ack.update(key, &updateArgs{
		Enable: &enable,
	})
	return
}

func (ack *AccessKey) Disable(key string) (err error) {
	enable := false
	err = ack.update(key, &updateArgs{
		Enable: &enable,
	})
	return
}

func (ack *AccessKey) Reset(key string) (err error) {
	secret := ack.newSecret()
	err = ack.update(key, &updateArgs{
		Secret: &secret,
	})
	return
}

func (ack *AccessKey) Delete(key string) (err error) {
	isDelete := true
	err = ack.update(key, &updateArgs{
		IsDelete: &isDelete,
	})
	return
}

func (ack *AccessKey) Get(key string) (record *handlers.AccessKeyRecord, err error) {
	record = &handlers.AccessKeyRecord{}
	err = ack.providers.GetStateStore().Get(ack.getStoreKey(key), record)
	if err != nil && !errors.Is(err, services.ErrStateStoreNotFound) {
		return
	}
	if errors.Is(err, services.ErrStateStoreNotFound) || record.IsDeleted {
		err = handlers.ErrAccessKeyIsNotFound
	}
	return
}

func (ack *AccessKey) List() (list []*handlers.AccessKeyRecord, err error) {
	err = ack.providers.GetStateStore().Iterate(ack.storeKeyPrefix, func(key, _ []byte) (stop bool, er error) {
		record := &handlers.AccessKeyRecord{}
		er = ack.providers.GetStateStore().Get(string(key), record)
		if er != nil {
			return
		}
		if record.IsDeleted {
			return
		}
		list = append(list, record)
		return
	})
	return
}

func (ack *AccessKey) newKey() (key string) {
	key = uuid.NewString()
	return
}

func (ack *AccessKey) newSecret() (secret string) {
	secret = utils.RandomString(ack.secretLength)
	return
}

func (ack *AccessKey) getStoreKey(key string) (storeKey string) {
	storeKey = ack.storeKeyPrefix + key
	return
}

func (ack *AccessKey) lock(key string) (unlock func()) {
	loaded := true
	for loaded {
		_, loaded = ack.locks.LoadOrStore(key, nil)
		time.Sleep(10 * time.Millisecond)
	}
	unlock = func() {
		ack.locks.Delete(key)
	}
	return
}

type updateArgs struct {
	Enable   *bool
	Secret   *string
	IsDelete *bool
}

func (ack *AccessKey) update(key string, args *updateArgs) (err error) {
	unlock := ack.lock(key)
	defer unlock()

	record := &handlers.AccessKeyRecord{}
	stk := ack.getStoreKey(key)

	err = ack.providers.GetStateStore().Get(stk, record)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return
	}
	if errors.Is(err, storage.ErrNotFound) || record.IsDeleted {
		err = handlers.ErrAccessKeyIsNotFound
		return
	}

	if args.Enable != nil {
		record.Enable = *args.Enable
	}
	if args.Secret != nil {
		record.Secret = *args.Secret
	}
	if args.IsDelete != nil {
		record.IsDeleted = *args.IsDelete
	}

	record.UpdatedAt = time.Now()

	err = ack.providers.GetStateStore().Put(stk, record)

	return
}
