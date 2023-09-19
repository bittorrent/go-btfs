package accesskey

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/api/providers"
	"github.com/bittorrent/go-btfs/s3/ctxmu"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/bittorrent/go-btfs/utils"
	"github.com/google/uuid"
	"time"
)

var _ Service = (*service)(nil)

type service struct {
	providers       providers.Providerser
	secretLength    int
	storeKeyPrefix  string
	lock            ctxmu.MultiCtxRWLocker
	waitLockTimeout time.Duration
}

func NewService(providers providers.Providerser, options ...Option) Service {
	svc := &service{
		providers:       providers,
		secretLength:    defaultSecretLength,
		storeKeyPrefix:  defaultStoreKeyPrefix,
		lock:            defaultLock,
		waitLockTimeout: defaultWaitLockTimout,
	}
	for _, option := range options {
		option(svc)
	}
	return svc
}

func (svc *service) Generate() (record *AccessKey, err error) {
	now := time.Now()
	record = &AccessKey{
		Key:       svc.newKey(),
		Secret:    svc.newSecret(),
		Enable:    true,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = svc.providers.StateStore().Put(svc.getStoreKey(record.Key), record)
	return
}

func (svc *service) Enable(key string) (err error) {
	enable := true
	err = svc.update(key, &updateArgs{
		Enable: &enable,
	})
	return
}

func (svc *service) Disable(key string) (err error) {
	enable := false
	err = svc.update(key, &updateArgs{
		Enable: &enable,
	})
	return
}

func (svc *service) Reset(key string) (err error) {
	secret := svc.newSecret()
	err = svc.update(key, &updateArgs{
		Secret: &secret,
	})
	return
}

func (svc *service) Delete(key string) (err error) {
	isDelete := true
	err = svc.update(key, &updateArgs{
		IsDelete: &isDelete,
	})
	return
}

func (svc *service) Get(key string) (ack *AccessKey, err error) {
	ack = &AccessKey{}
	err = svc.providers.StateStore().Get(svc.getStoreKey(key), ack)
	if err != nil && !errors.Is(err, providers.ErrStateStoreNotFound) {
		return
	}
	if errors.Is(err, providers.ErrStateStoreNotFound) || ack.IsDeleted {
		err = ErrNotFound
	}
	return
}

func (svc *service) List() (list []*AccessKey, err error) {
	err = svc.providers.StateStore().Iterate(svc.storeKeyPrefix, func(key, _ []byte) (stop bool, er error) {
		record := &AccessKey{}
		er = svc.providers.StateStore().Get(string(key), record)
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

func (svc *service) newKey() (key string) {
	key = uuid.NewString()
	return
}

func (svc *service) newSecret() (secret string) {
	secret = utils.RandomString(svc.secretLength)
	return
}

func (svc *service) getStoreKey(key string) (storeKey string) {
	storeKey = svc.storeKeyPrefix + key
	return
}

type updateArgs struct {
	Enable   *bool
	Secret   *string
	IsDelete *bool
}

func (svc *service) update(key string, args *updateArgs) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), svc.waitLockTimeout)
	defer cancel()

	err = svc.lock.Lock(ctx, key)
	if err != nil {
		return
	}
	defer svc.lock.Unlock(key)

	ack := &AccessKey{}
	stk := svc.getStoreKey(key)

	err = svc.providers.StateStore().Get(stk, ack)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return
	}
	if errors.Is(err, storage.ErrNotFound) || ack.IsDeleted {
		err = ErrNotFound
		return
	}

	if args.Enable != nil {
		ack.Enable = *args.Enable
	}
	if args.Secret != nil {
		ack.Secret = *args.Secret
	}
	if args.IsDelete != nil {
		ack.IsDeleted = *args.IsDelete
	}

	ack.UpdatedAt = time.Now()

	err = svc.providers.StateStore().Put(stk, ack)

	return
}
