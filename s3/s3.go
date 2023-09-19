package s3

import (
	config "github.com/bittorrent/go-btfs-config"
	"github.com/bittorrent/go-btfs/s3/api/handlers"
	"github.com/bittorrent/go-btfs/s3/api/providers"
	"github.com/bittorrent/go-btfs/s3/api/routers"
	"github.com/bittorrent/go-btfs/s3/api/server"
	"github.com/bittorrent/go-btfs/s3/api/services/accesskey"
	"github.com/bittorrent/go-btfs/s3/api/services/object"
	"github.com/bittorrent/go-btfs/s3/api/services/sign"
	"github.com/bittorrent/go-btfs/s3/ctxmu"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"sync"
)

var (
	ps   *providers.Providers
	once sync.Once
)

func InitProviders(stateStore storage.StateStorer) (err error) {
	once.Do(func() {
		var (
			sstore providers.StateStorer
			fstore providers.FileStorer
		)
		sstore = providers.NewStorageStateStoreProxy(stateStore)
		fstore, err = providers.NewBtfsAPI()
		if err != nil {
			return
		}
		ps = providers.NewProviders(sstore, fstore)
	})
	return
}

func GetProviders() *providers.Providers {
	return ps
}

func NewServer(cfg config.S3CompatibleAPI) *server.Server {
	// global multiple keys read write lock
	lock := ctxmu.NewDefaultMultiCtxRWMutex()

	// services
	sigsvc := sign.NewService()
	acksvc := accesskey.NewService(ps, accesskey.WithLock(lock))
	objsvc := object.NewService(ps, object.WithLock(lock))

	// handlers
	hs := handlers.NewHandlers(
		acksvc, sigsvc, objsvc,
		handlers.WithHeaders(cfg.HTTPHeaders),
	)

	// routers
	rs := routers.NewRouters(hs)

	// server
	svr := server.NewServer(
		rs,
		server.WithAddress(cfg.Address),
	)

	return svr
}
