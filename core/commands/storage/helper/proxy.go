package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/core"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

const (
	ProxyStorageInfoPrefix           = "/btfs/proxy_storage/"
	ProxyStoragePaymentPrefix        = "/btfs/proxy_payment/"
	ProxyStoragePaymentBalancePrefix = "/btfs/proxy_payment_balance/"
	ProxyNeedPayCIDPrefix            = "/btfs/proxy_need_pay_cid/"
	ProxyUploadedFileInfoPrefix      = "/btfs/proxy_uploaded_file/"
)

const (
	DefaultPayTimeout = 30 * time.Minute
)

type ProxyStorageInfo struct {
	Price uint64 `json:"price"`
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
	return rds.Put(ctx, GetProxyStoragePaymentKey(node.Identity.String()+"/"+strings.ToLower(ns.From)+"/"+ns.Hash), b)
}

func GetProxyStoragePayment(ctx context.Context, node *core.IpfsNode) ([]*ProxyStoragePaymentInfo, error) {
	rds := node.Repo.Datastore()
	qr, err := rds.Query(ctx, query.Query{
		Prefix: GetProxyStoragePaymentKey(node.Identity.String()).String(),
	})
	if err != nil {
		return nil, err
	}
	ret := make([]*ProxyStoragePaymentInfo, 0)
	for r := range qr.Next() {
		if r.Error != nil {
			return nil, r.Error
		}
		var ns ProxyStoragePaymentInfo
		err = json.Unmarshal(r.Entry.Value, &ns)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &ns)
	}
	return ret, nil
}

func GetProxyStoragePaymentList(ctx context.Context, node *core.IpfsNode, from string) ([]*ProxyStoragePaymentInfo, error) {
	rds := node.Repo.Datastore()
	qr, err := rds.Query(ctx, query.Query{
		Prefix: GetProxyStoragePaymentKey(node.Identity.String() + "/" + strings.ToLower(from)).String(),
	})
	if err != nil {
		return nil, err
	}
	ret := make([]*ProxyStoragePaymentInfo, 0)
	for r := range qr.Next() {
		if r.Error != nil {
			return nil, r.Error
		}
		ns := new(ProxyStoragePaymentInfo)
		err = json.Unmarshal(r.Entry.Value, ns)
		if err != nil {
			return nil, err
		}
		ret = append(ret, ns)
	}
	return ret, nil
}

func GetProxyStoragePaymentByTxHash(ctx context.Context, node *core.IpfsNode, from, txHash string) (*ProxyStoragePaymentInfo, error) {
	rds := node.Repo.Datastore()
	b, err := rds.Get(ctx, GetProxyStoragePaymentKey(node.Identity.String()+"/"+strings.ToLower(from)+"/"+txHash))
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

type ProxyNeedPaymentInfo struct {
	CID      string `json:"cid"`
	FileSize int64  `json:"file_size"`
	Price    int64  `json:"price"`
	ExpireAt int64  `json:"expire_at"`
	NeedBTT  uint64 `json:"need_btt"`
}

func PutProxyNeedPaymentCID(ctx context.Context, node *core.IpfsNode, needPayInfo *ProxyNeedPaymentInfo) error {
	rds := node.Repo.Datastore()
	b, err := json.Marshal(needPayInfo)
	if err != nil {
		return fmt.Errorf("cannot put current proxy storage settings: %s", err.Error())
	}
	return rds.Put(ctx, GetProxyNeedPaymentKey(node.Identity.String()+"/"+needPayInfo.CID), b)
}

func GetProxyNeedPaymentCID(ctx context.Context, node *core.IpfsNode, cid string) (*ProxyNeedPaymentInfo, error) {
	rds := node.Repo.Datastore()
	p, err := rds.Get(ctx, GetProxyNeedPaymentKey(node.Identity.String()+"/"+cid))
	if err != nil {
		return nil, err
	}
	var ns ProxyNeedPaymentInfo
	err = json.Unmarshal(p, &ns)
	return &ns, err
}

func DeleteProxyNeedPaymentCID(ctx context.Context, node *core.IpfsNode, cid string) error {
	rds := node.Repo.Datastore()
	return rds.Delete(ctx, GetProxyNeedPaymentKey(node.Identity.String()+"/"+cid))
}

func GetProxyNeedPaymentKey(cid string) ds.Key {
	return NewKeyHelper(ProxyNeedPayCIDPrefix, cid)
}

type ProxyUploadFileInfo struct {
	From      string `json:"from"`
	CID       string `json:"cid"`
	FileSize  int64  `json:"file_size"`
	Price     int64  `json:"price"`
	ExpireAt  int64  `json:"expire_at"`
	TotalPay  uint64 `json:"total_pay"`
	CreatedAt int64  `json:"created_at"`
}

func PutProxyUploadedFileInfo(ctx context.Context, node *core.IpfsNode, uploadedCidInfo *ProxyUploadFileInfo) error {
	rds := node.Repo.Datastore()
	b, err := json.Marshal(uploadedCidInfo)
	if err != nil {
		return fmt.Errorf("cannot put current proxy storage settings: %s", err.Error())
	}
	return rds.Put(ctx, GetProxyUploadedFileInfoKey(node.Identity.String()+"/"+uploadedCidInfo.CID), b)
}

func ListProxyUploadedFileInfo(ctx context.Context, node *core.IpfsNode) ([]*ProxyUploadFileInfo, error) {
	rds := node.Repo.Datastore()
	qr, err := rds.Query(ctx, query.Query{
		Prefix: GetProxyUploadedFileInfoKey(node.Identity.String()).String(),
	})
	if err != nil {
		return nil, err
	}
	ret := make([]*ProxyUploadFileInfo, 0)
	for r := range qr.Next() {
		if r.Error != nil {
			return nil, r.Error
		}
		var ns ProxyUploadFileInfo
		err = json.Unmarshal(r.Entry.Value, &ns)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &ns)
	}
	return ret, nil
}

func GetProxyUploadedFileInfoKey(cid string) ds.Key {
	return NewKeyHelper(ProxyUploadedFileInfoPrefix, cid)
}
