package s3

import (
	config "github.com/bittorrent/go-btfs-config"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/s3/handlers"
	"github.com/bittorrent/go-btfs/s3/providers"
	"github.com/bittorrent/go-btfs/s3/routers"
	"github.com/bittorrent/go-btfs/s3/server"
	"github.com/bittorrent/go-btfs/s3/services/accesskey"
	"github.com/bittorrent/go-btfs/s3/services/object"
	"github.com/bittorrent/go-btfs/s3/services/sign"
	"sync"
)

var (
	ps   *providers.Providers
	once sync.Once
)

func initProviders() {
	once.Do(func() {
		sstore := providers.NewStorageStateStoreProxy(chain.StateStore)
		fstore := providers.NewBtfsAPI("")
		ps = providers.NewProviders(sstore, fstore)
	})
}

func GetProviders() *providers.Providers {
	initProviders()
	return ps
}

func NewServer(cfg config.S3CompatibleAPI) *server.Server {
	// providers
	initProviders()

	// services
	acksvc := accesskey.NewService(ps)
	sigsvc := sign.NewService()
	objsvc := object.NewService(ps)

	// handlers
	hs := handlers.NewHandlers(
		acksvc, sigsvc, objsvc, handlers.WithHeaders(cfg.HTTPHeaders),
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
