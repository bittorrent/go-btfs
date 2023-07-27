package accesskey

import (
	"errors"
	"github.com/bittorrent/go-btfs/s3/utils/random"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"github.com/google/uuid"
	"sync"
	"time"
)

var _ Service = &service{}

type service struct {
	config *Config
	store  storage.StateStorer
	locks  sync.Map
}

func newService(config *Config, store storage.StateStorer) *service {
	return &service{
		config: config,
		store:  store,
		locks:  sync.Map{},
	}
}

func (s *service) Generate() (ack *AccessKey, err error) {
	now := time.Now()
	ack = &AccessKey{
		Key:       s.newKey(),
		Secret:    s.newSecret(),
		Enable:    true,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = s.store.Put(s.getStoreKey(ack.Key), ack)
	return
}

func (s *service) Enable(key string) (err error) {
	enable := true
	err = s.update(key, &updateArgs{
		Enable: &enable,
	})
	return
}

func (s *service) Disable(key string) (err error) {
	enable := false
	err = s.update(key, &updateArgs{
		Enable: &enable,
	})
	return
}

func (s *service) Reset(key string) (err error) {
	secret := s.newSecret()
	err = s.update(key, &updateArgs{
		Secret: &secret,
	})
	return
}

func (s *service) Delete(key string) (err error) {
	isDelete := true
	err = s.update(key, &updateArgs{
		IsDelete: &isDelete,
	})
	return
}

func (s *service) Get(key string) (ack *AccessKey, err error) {
	ack = &AccessKey{}
	err = s.store.Get(s.getStoreKey(key), ack)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return
	}
	if errors.Is(err, storage.ErrNotFound) || ack.IsDeleted {
		err = ErrNotFound
	}
	return
}

func (s *service) List() (list []*AccessKey, err error) {
	err = s.store.Iterate(s.config.StorePrefix, func(key, _ []byte) (stop bool, er error) {
		ack := &AccessKey{}
		er = s.store.Get(string(key), ack)
		if er != nil {
			return
		}
		if ack.IsDeleted {
			return
		}
		list = append(list, ack)
		return
	})
	return
}

func (s *service) newKey() (key string) {
	key = uuid.NewString()
	return
}

func (s *service) newSecret() (secret string) {
	secret = random.NewString(s.config.SecretLength)
	return
}

func (s *service) getStoreKey(key string) (storeKey string) {
	storeKey = s.config.StorePrefix + key
	return
}

func (s *service) lock(key string) (unlock func()) {
	loaded := true
	for loaded {
		_, loaded = s.locks.LoadOrStore(key, nil)
		time.Sleep(10 * time.Millisecond)
	}
	unlock = func() {
		s.locks.Delete(key)
	}
	return
}

type updateArgs struct {
	Enable   *bool
	Secret   *string
	IsDelete *bool
}

func (s *service) update(key string, args *updateArgs) (err error) {
	unlock := s.lock(key)
	defer unlock()

	ack := &AccessKey{}
	stk := s.getStoreKey(key)

	err = s.store.Get(stk, ack)
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

	err = s.store.Put(stk, ack)

	return
}
