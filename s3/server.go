package s3

import (
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/routers"
	"github.com/bittorrent/go-btfs/s3/server"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"github.com/bittorrent/go-btfs/s3/services/bucket"
	"github.com/bittorrent/go-btfs/s3/services/cors"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"github.com/bittorrent/go-btfs/s3/services/sign"
	"github.com/bittorrent/go-btfs/transaction/storage"
	"sync"
)

var (
	ps   *providers.Providers
	once sync.Once
)

func GetProviders(storageStore storage.StateStorer) *providers.Providers {
	once.Do(func() {
		sstore := providers.NewStorageStateStoreProxy(storageStore)
		fstore := providers.NewFileStore()
		ps = providers.NewProviders(sstore, fstore)
	})
	return ps
}

func NewServer(storageStore storage.StateStorer) *server.Server {
	_ = GetProviders(storageStore)

	// services
	corsvc := cors.NewService()
	acksvc := accesskey.NewService(ps)
	sigsvc := sign.NewService()
	bucsvc := bucket.NewService(ps)
	bucsvc.SetEmptyBucket(bucsvc.EmptyBucket) //todo EmptyBucket参数后续更新为object对象
	objsvc := object.NewService(ps)

	// handlers
	hs := handlers.NewHandlers(corsvc, acksvc, sigsvc, bucsvc, objsvc)

	// routers
	rs := routers.NewRouters(hs)

	// server
	svr := server.NewServer(rs)

	return svr
}
