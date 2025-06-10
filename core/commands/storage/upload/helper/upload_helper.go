package helper

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"

	cmds "github.com/bittorrent/go-btfs-cmds"
	config "github.com/bittorrent/go-btfs-config"
	iface "github.com/bittorrent/interface-go-btfs-core"
	"github.com/bittorrent/interface-go-btfs-core/path"

	"github.com/alecthomas/units"
	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	cidlib "github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
)

const (
	uploadPriceOptionName   = "price"
	storageLengthOptionName = "storage-length"
)

var (
	log = logging.Logger("upload")

	WaitUploadBo = func(maxTime time.Duration) *backoff.ExponentialBackOff {
		bo := backoff.NewExponentialBackOff()
		bo.InitialInterval = 10 * time.Second
		bo.MaxElapsedTime = maxTime
		bo.Multiplier = 1.5
		bo.MaxInterval = 10 * time.Minute
		return bo
	}
	HandleShardBo = func() *backoff.ExponentialBackOff {
		bo := backoff.NewExponentialBackOff()
		bo.InitialInterval = 1 * time.Second
		bo.MaxElapsedTime = 300 * time.Second
		bo.Multiplier = 1
		bo.MaxInterval = 1 * time.Second
		return bo
	}()
	CheckPaymentBo = func() *backoff.ExponentialBackOff {
		bo := backoff.NewExponentialBackOff()
		bo.InitialInterval = 10 * time.Second
		bo.MaxElapsedTime = 5 * time.Minute
		bo.Multiplier = 1.5
		bo.MaxInterval = 60 * time.Second
		return bo
	}()
	DownloadShardBo = func(maxTime time.Duration) *backoff.ExponentialBackOff {
		bo := backoff.NewExponentialBackOff()
		bo.InitialInterval = 10 * time.Second
		bo.MaxElapsedTime = maxTime
		bo.Multiplier = 1.5
		bo.MaxInterval = 30 * time.Minute
		return bo
	}
	WaitingForPeersBo = func() *backoff.ExponentialBackOff {
		bo := backoff.NewExponentialBackOff()
		bo.InitialInterval = 1 * time.Second
		bo.MaxElapsedTime = 300 * time.Second
		bo.Multiplier = 1.2
		bo.MaxInterval = 5 * time.Second
		return bo
	}()
)

type ContextParams struct {
	Req *cmds.Request
	Env cmds.Environment
	Ctx context.Context
	N   *core.IpfsNode
	Cfg *config.Config
	Api iface.CoreAPI
}

func ExtractContextParams(req *cmds.Request, env cmds.Environment) (*ContextParams, error) {
	// get node
	n, err := cmdenv.GetNode(env)
	if err != nil {
		return nil, err
	}
	// get config settings
	cfg, err := n.Repo.Config()
	if err != nil {
		return nil, err
	}
	// get core api
	api, err := cmdenv.GetApi(env, req)
	if err != nil {
		return nil, err
	}
	ctx, _ := helper.NewGoContext(req.Context)
	return &ContextParams{
		Req: req,
		Env: env,
		Ctx: ctx,
		N:   n,
		Cfg: cfg,
		Api: api,
	}, nil
}

func GetShardHashes(params *ContextParams, fileHash string) (shardHashes []string, fileSize int64,
	shardSize int64, err error) {
	fileCid, err := cidlib.Parse(fileHash)
	if err != nil {
		return nil, -1, -1, err
	}
	cids, fileSize, err := helper.CheckAndGetReedSolomonShardHashes(params.Ctx, params.N, params.Api, fileCid)
	if err != nil || len(cids) == 0 {
		return nil, -1, -1, fmt.Errorf("invalid hash: %s", err)
	}

	shardHashes = make([]string, 0)
	for _, c := range cids {
		shardHashes = append(shardHashes, c.String())
	}
	shardCid, err := cidlib.Parse(shardHashes[0])
	if err != nil {
		return nil, -1, -1, err
	}
	sz, err := getNodeSizeFromCid(params.Ctx, shardCid, params.Api)
	if err != nil {
		return nil, -1, -1, err
	}
	shardSize = int64(sz)
	return
}

func GetShardHashesCopy(params *ContextParams, fileHash string, copyNum int) (shardHashes []string, fileSize int64,
	shardSize int64, err error) {

	fileCid, err := cidlib.Parse(fileHash)
	if err != nil {
		return nil, -1, -1, err
	}
	sz, err := getNodeSizeFromCid(params.Ctx, fileCid, params.Api)
	if err != nil {
		return nil, -1, -1, err
	}
	fileSize = int64(sz)
	shardSize = fileSize

	shardHashes = make([]string, 0)
	for i := 0; i < copyNum+1; i++ {
		shardHashes = append(shardHashes, fileHash)
	}

	return
}

func GetPriceAndMinStorageLength(params *ContextParams) (price int64, storageLength int, err error) {
	ns, err := helper.GetHostStorageConfig(params.Ctx, params.N)
	if err != nil {
		return -1, -1, err
	}
	price, found := params.Req.Options[uploadPriceOptionName].(int64)
	if !found {
		price = int64(ns.StoragePriceAsk)
	}
	storageLength = params.Req.Options[storageLengthOptionName].(int)
	if uint64(storageLength) < ns.StorageTimeMin {
		return -1, -1, fmt.Errorf("invalid storage len. want: >= %d, got: %d",
			ns.StorageTimeMin, storageLength)
	}
	return
}

func TotalPay(shardSize int64, price int64, storageLength int, rate *big.Int) (int64, error) {
	if rate.Cmp(big.NewInt(1)) > 0 {
		return TotalPayReal(shardSize, price, storageLength, rate)
	}
	return TotalPayRound(shardSize, price, storageLength, rate)
}

func TotalPayReal(shardSize int64, price int64, storageLength int, rate *big.Int) (int64, error) {
	totalPay := int64(float64(shardSize) / float64(units.GiB) * float64(price) * float64(storageLength))
	if totalPay <= 0 {
		totalPay = 1
	}
	//fmt.Printf("size:%v GB, price:%v000000000000 , storageLength:%v,  TotalPay:%v000000000000 \n", float64(shardSize)/float64(units.GiB), price, storageLength, totalPay)
	fmt.Printf("size:%v GB, price:%v*%v , storageLength:%v,  TotalPay:%v*%v \n", float64(shardSize)/float64(units.GiB), price, rate.String(), storageLength, totalPay, rate.String())

	return totalPay, nil
}

func TotalPayRound(shardSize int64, price int64, storageLength int, rate *big.Int) (int64, error) {
	totalPayFloat := float64(shardSize) / float64(units.GiB) * float64(price) * float64(storageLength)

	totalPay := int64(math.Floor(totalPayFloat + 0.5))
	//fmt.Printf("size:%v GB, price:%v*%v , storageLength:%v,  TotalPay:%v*%v \n", float64(shardSize)/float64(units.GiB), price, rate.String(), storageLength, totalPay, rate.String())
	fmt.Printf("size:%v GB, price:%v*%v , storageLength:%v,  TotalPay:%v*%v, TotalPayFloat:%v*%v \n", float64(shardSize)/float64(units.GiB), price, rate.String(), storageLength, totalPay, rate.String(), totalPayFloat, rate.String())
	if totalPay < 1 {
		return 0, errors.New("Your file upload fee is less than 1 precision minimum unit, please use WBTT token to upload. ")
	}

	return totalPay, nil
}

func NewContractID(sessionId string) string {
	id := uuid.New().String()
	return sessionId + "," + id
}

func SplitContractId(contractId string) (ssId string, shardHash string) {
	splits := strings.Split(contractId, ",")
	ssId = splits[0]
	shardHash = splits[1]
	return
}

func getNodeSizeFromCid(ctx context.Context, hash cidlib.Cid, api iface.CoreAPI) (uint64, error) {
	leafPath := path.IpfsPath(hash)
	ipldNode, err := api.ResolveNode(ctx, leafPath)
	if err != nil {
		return 0, err
	}
	return ipldNode.Size()
}
