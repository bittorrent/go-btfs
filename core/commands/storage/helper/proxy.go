package helper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bittorrent/go-btfs/core"
	ds "github.com/ipfs/go-datastore"
)

const (
	ProxyStorageInfoPrefix = "/proxy_storage/" // self or from network
)

type ProxyStorageInfo struct {
	Price uint64
}

// PutProxyStorageConfig server as a proxy node save storage config.
func PutProxyStorageConfig(ctx context.Context, node *core.IpfsNode, ns *ProxyStorageInfo) error {
	rds := node.Repo.Datastore()
	b, err := json.Marshal(ns)
	if err != nil {
		return fmt.Errorf("cannot put current proxy storage settings: %s", err.Error())
	}
	return rds.Put(ctx, GetProxyStorageKey(node.Identity.String()), b)
}

func GetProxyStorageConfig(ctx context.Context, node *core.IpfsNode) (*ProxyStorageInfo, error) {
	rds := node.Repo.Datastore()
	b, err := rds.Get(ctx, GetProxyStorageKey(node.Identity.String()))
	if err != nil {
		return nil, err
	}
	ns := new(ProxyStorageInfo)
	err = json.Unmarshal(b, ns)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func GetProxyStorageKey(pid string) ds.Key {
	return NewKeyHelper(ProxyStorageInfoPrefix, pid)
}
