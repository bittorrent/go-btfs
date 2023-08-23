package s3

import (
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/routers"
	"github.com/bittorrent/go-btfs/s3/server"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"github.com/bittorrent/go-btfs/s3/services/auth"
	"github.com/bittorrent/go-btfs/s3/services/bucket"
	"github.com/bittorrent/go-btfs/s3/services/cors"
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
	corsSvc := cors.NewService()
	accessKeySvc := accesskey.NewService(ps)
	authSvc := auth.NewService(ps, accessKeySvc)
	bucketSvc := bucket.NewService(ps)

	// handlers
	hs := handlers.NewHandlers(corsSvc, authSvc, bucketSvc, nil, nil)

	// routers
	rs := routers.NewRouters(hs)

	// server
	svr := server.NewServer(rs)

	return svr
}
