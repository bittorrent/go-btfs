package accesskey

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/s3/ctxmu"
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/services"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/bittorrent/go-btfs/utils"
	"github.com/google/uuid"
	"time"
)

const (
	defaultSecretLength    = 32
	defaultStoreKeyPrefix  = "access-keys:"
	defaultUpdateTimeoutMS = 200
)

var _ handlers.AccessKeyService = (*Service)(nil)

type Service struct {
	providers      services.Providerser
	secretLength   int
	storeKeyPrefix string
	locks          *ctxmu.MultiCtxRWMutex
	updateTimeout  time.Duration
}

func NewService(providers services.Providerser, options ...Option) (svc *Service) {
	svc = &Service{
		providers:      providers,
		secretLength:   defaultSecretLength,
		storeKeyPrefix: defaultStoreKeyPrefix,
		locks:          ctxmu.NewDefaultMultiCtxRWMutex(),
		updateTimeout:  time.Duration(defaultUpdateTimeoutMS) * time.Millisecond,
	}
	for _, option := range options {
		option(svc)
	}
	return svc
}

func (svc *Service) Generate() (record *handlers.AccessKeyRecord, err error) {
	now := time.Now()
	record = &handlers.AccessKeyRecord{
		Key:       svc.newKey(),
		Secret:    svc.newSecret(),
		Enable:    true,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = svc.providers.GetStateStore().Put(svc.getStoreKey(record.Key), record)
	return
}

func (svc *Service) Enable(key string) (err error) {
	enable := true
	err = svc.update(key, &updateArgs{
		Enable: &enable,
	})
	return
}

func (svc *Service) Disable(key string) (err error) {
	enable := false
	err = svc.update(key, &updateArgs{
		Enable: &enable,
	})
	return
}

func (svc *Service) Reset(key string) (err error) {
	secret := svc.newSecret()
	err = svc.update(key, &updateArgs{
		Secret: &secret,
	})
	return
}

func (svc *Service) Delete(key string) (err error) {
	isDelete := true
	err = svc.update(key, &updateArgs{
		IsDelete: &isDelete,
	})
	return
}

func (svc *Service) Get(key string) (record *handlers.AccessKeyRecord, err error) {
	record = &handlers.AccessKeyRecord{}
	err = svc.providers.GetStateStore().Get(svc.getStoreKey(key), record)
	if err != nil && !errors.Is(err, services.ErrStateStoreNotFound) {
		return
	}
	if errors.Is(err, services.ErrStateStoreNotFound) || record.IsDeleted {
		err = handlers.ErrAccessKeyIsNotFound
	}
	return
}

func (svc *Service) List() (list []*handlers.AccessKeyRecord, err error) {
	err = svc.providers.GetStateStore().Iterate(svc.storeKeyPrefix, func(key, _ []byte) (stop bool, er error) {
		record := &handlers.AccessKeyRecord{}
		er = svc.providers.GetStateStore().Get(string(key), record)
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

func (svc *Service) newKey() (key string) {
	key = uuid.NewString()
	return
}

func (svc *Service) newSecret() (secret string) {
	secret = utils.RandomString(svc.secretLength)
	return
}

func (svc *Service) getStoreKey(key string) (storeKey string) {
	storeKey = svc.storeKeyPrefix + key
	return
}

type updateArgs struct {
	Enable   *bool
	Secret   *string
	IsDelete *bool
}

func (svc *Service) update(key string, args *updateArgs) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), svc.updateTimeout)
	defer cancel()

	err = svc.locks.Lock(ctx, key)
	if err != nil {
		return
	}
	defer svc.locks.Unlock(key)

	record := &handlers.AccessKeyRecord{}
	stk := svc.getStoreKey(key)

	err = svc.providers.GetStateStore().Get(stk, record)
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

	err = svc.providers.GetStateStore().Put(stk, record)

	return
}
