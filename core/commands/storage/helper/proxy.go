package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bittorrent/go-btfs/core"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

const (
	ProxyStorageInfoPrefix           = "/proxy_storage/" // self or from network
	ProxyStoragePaymentPrefix        = "/proxy_payment/" // self or from network
	ProxyStoragePaymentBalancePrefix = "/proxy_payment_balance/"
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

type ProxyStoragePaymentInfo struct {
	From    string
	To      string
	Value   uint64
	Hash    string
	PayTime int64
	Balance uint64
}

func PutProxyStoragePayment(ctx context.Context, node *core.IpfsNode, ns *ProxyStoragePaymentInfo) error {
	rds := node.Repo.Datastore()
	b, err := json.Marshal(ns)
	if err != nil {
		return fmt.Errorf("cannot put current proxy storage settings: %s", err.Error())
	}
	// /proxy_payment/address/txHash
	return rds.Put(ctx, GetProxyStoragePaymentKey(node.Identity.String()+"/"+ns.From+"/"+ns.Hash), b)
}

func GetProxyStoragePayment(ctx context.Context, node *core.IpfsNode) (*ProxyStoragePaymentInfo, error) {
	rds := node.Repo.Datastore()
	b, err := rds.Get(ctx, GetProxyStoragePaymentKey(node.Identity.String()))
	if err != nil {
		return nil, err
	}
	ns := new(ProxyStoragePaymentInfo)
	err = json.Unmarshal(b, ns)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func GetProxyStoragePaymentByTxHash(ctx context.Context, node *core.IpfsNode, from, txHash string) (*ProxyStoragePaymentInfo, error) {
	rds := node.Repo.Datastore()
	b, err := rds.Get(ctx, GetProxyStoragePaymentKey(node.Identity.String()+"/"+from+"/"+txHash))
	if err != nil && !errors.Is(err, ds.ErrNotFound) {
		return nil, err
	}
	if errors.Is(err, ds.ErrNotFound) {
		return nil, nil
	}
	ns := new(ProxyStoragePaymentInfo)
	err = json.Unmarshal(b, ns)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func ChargeBalance(ctx context.Context, node *core.IpfsNode, from string, value uint64) (uint64, error) {
	rds := node.Repo.Datastore()
	b, err := rds.Get(ctx, GetProxyStoragePaymentBalanceKey(node.Identity.String()+"/"+strings.ToLower(from)))
	if err != nil && !errors.Is(err, ds.ErrNotFound) {
		return 0, err
	}
	if errors.Is(err, ds.ErrNotFound) {
		return value, rds.Put(ctx, GetProxyStoragePaymentBalanceKey(node.Identity.String()+"/"+strings.ToLower(from)), []byte(fmt.Sprintf("%d", value)))
	}

	balance, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return 0, err
	}
	balance += value
	return balance, rds.Put(ctx, GetProxyStoragePaymentBalanceKey(node.Identity.String()+"/"+strings.ToLower(from)), []byte(fmt.Sprintf("%d", balance)))
}

func SubBalance(ctx context.Context, node *core.IpfsNode, from string, value uint64) error {
	rds := node.Repo.Datastore()
	b, err := rds.Get(ctx, GetProxyStoragePaymentBalanceKey(node.Identity.String()+"/"+strings.ToLower(from)))
	if err != nil {
		return err
	}
	balance, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return err
	}
	balance -= value
	return rds.Put(ctx, GetProxyStoragePaymentBalanceKey(node.Identity.String()+"/"+strings.ToLower(from)), []byte(fmt.Sprintf("%d", balance)))
}

func GetBalance(ctx context.Context, node *core.IpfsNode, from string) (uint64, error) {
	rds := node.Repo.Datastore()
	b, err := rds.Get(ctx, GetProxyStoragePaymentBalanceKey(node.Identity.String()+"/"+strings.ToLower(from)))
	if err != nil {
		return 0, err
	}
	balance, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func GetBalanceList(ctx context.Context, node *core.IpfsNode) (map[string]uint64, error) {
	rds := node.Repo.Datastore()
	qr, err := rds.Query(ctx, query.Query{
		Prefix: GetProxyStoragePaymentBalanceKey(node.Identity.String()).String(),
	})
	if err != nil {
		return nil, err
	}
	ret := make(map[string]uint64)
	for r := range qr.Next() {
		if r.Error != nil {
			return nil, r.Error
		}
		balance, err := strconv.ParseUint(string(r.Entry.Value), 10, 64)
		if err != nil {
			return nil, err
		}
		// key is /proxy_payment_balance/peerId/address
		// set key to address only
		key := r.Entry.Key[len(GetProxyStoragePaymentBalanceKey(node.Identity.String()).String())+1:]
		ret[key] = balance
	}
	return ret, nil
}

func GetProxyStoragePaymentKey(peerId string) ds.Key {
	return NewKeyHelper(ProxyStoragePaymentPrefix, peerId)
}

func GetProxyStoragePaymentBalanceKey(peerId string) ds.Key {
	return NewKeyHelper(ProxyStoragePaymentBalancePrefix, peerId)
}
