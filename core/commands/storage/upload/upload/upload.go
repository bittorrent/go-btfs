package upload

import (
	"context"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"strconv"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/settlement/swap/swapprotocol"

	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/storage/hosts"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/offline"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	renterpb "github.com/bittorrent/go-btfs/protos/renter"

	cmds "github.com/bittorrent/go-btfs-cmds"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p/core/peer"
	cmap "github.com/orcaman/concurrent-map"
)

const (
	replicationFactorOptionName      = "replication-factor"
	hostSelectModeOptionName         = "host-select-mode"
	hostSelectionOptionName          = "host-selection"
	testOnlyOptionName               = "host-search-local"
	customizedPayoutOptionName       = "customize-payout"
	customizedPayoutPeriodOptionName = "customize-payout-period"
	copyName                         = "copy"

	defaultRepFactor     = 3
	defaultStorageLength = 30

	uploadPriceOptionName   = "price"
	storageLengthOptionName = "storage-length"
)

var (
	ShardErrChanMap = cmap.New()
	log             = logging.Logger("upload")
)

var StorageUploadCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Store files on BTFS network nodes through BTT payment.",
		ShortDescription: `
By default, BTFS selects hosts based on overall score according to the current client's environment.
To upload a file, <file-hash> must refer to a reed-solomon encoded file.

To create a reed-solomon encoded file from a normal file:

    $ btfs add --chunker=reed-solomon <file>
    added <file-hash> <file>

Run command to upload:

    $ btfs storage upload <file-hash>

To custom upload and storage a file on specific hosts:
    Use -m with 'custom' mode, and put host identifiers in -s, with multiple hosts separated by ','.

    # Upload a file to a set of hosts
    # Total # of hosts (N) must match # of shards in the first DAG level of root file hash
    $ btfs storage upload <file-hash> -m=custom -s=<host1-peer-id>,<host2-peer-id>,...,<hostN-peer-id>

    # Upload specific shards to a set of hosts
    # Total # of hosts (N) must match # of shards given
    $ btfs storage upload <shard-hash1> <shard-hash2> ... <shard-hashN> -l -m=custom -s=<host1-peer-id>,<host2-peer-id>,...,<hostN-peer-id>

Use status command to check for completion:
    $ btfs storage upload status <session-id> | jq`,
	},
	Subcommands: map[string]*cmds.Command{
		"init":              StorageUploadInitCmd,
		"supporttokens":     StorageUploadSupportTokensCmd,
		"cheque":            StorageUploadChequeCmd,
		"recvcontract":      StorageUploadRecvContractCmd,
		"status":            StorageUploadStatusCmd,
		"repair":            StorageUploadRepairCmd,
		"getcontractbatch":  offline.StorageUploadGetContractBatchCmd,
		"signcontractbatch": offline.StorageUploadSignContractBatchCmd,
		"getunsigned":       offline.StorageUploadGetUnsignedCmd,
		"sign":              offline.StorageUploadSignCmd,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("file-hash", true, false, "Hash of file to upload."),
		cmds.StringArg("upload-peer-id", false, false, "Peer id when upload upload."),
		cmds.StringArg("upload-nonce-ts", false, false, "Nounce timestamp when upload upload."),
		cmds.StringArg("upload-signature", false, false, "Session signature when upload upload."),
	},
	Options: []cmds.Option{
		cmds.Int64Option(uploadPriceOptionName, "p", "Max price per GiB per day of storage in µBTT (=0.000001BTT)."),
		cmds.IntOption(replicationFactorOptionName, "r", "Replication factor for the file with erasure coding built-in.").WithDefault(defaultRepFactor),
		cmds.StringOption(hostSelectModeOptionName, "m", "Based on this mode to select hosts and upload automatically. Default: mode set in config option Experimental.HostsSyncMode."),
		cmds.StringOption(hostSelectionOptionName, "s", "Use only these selected hosts in order on 'custom' mode. Use ',' as delimiter."),
		cmds.BoolOption(testOnlyOptionName, "t", "Enable host search under all domains 0.0.0.0 (useful for local test)."),
		cmds.IntOption(storageLengthOptionName, "len", "File storage period on hosts in days.").WithDefault(defaultStorageLength),
		cmds.BoolOption(customizedPayoutOptionName, "Enable file storage customized payout schedule.").WithDefault(false),
		cmds.IntOption(customizedPayoutPeriodOptionName, "Period of customized payout schedule.").WithDefault(1),
		cmds.IntOption(copyName, "copy num of file hash.").WithDefault(0),
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	RunTimeout: 15 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		swapprotocol.Req = req
		swapprotocol.Env = env

		ssId := uuid.New().String()
		ctxParams, err := helper.ExtractContextParams(req, env)
		if err != nil {
			return err
		}
		renterId := ctxParams.N.Identity
		offlineSigning := false
		if len(req.Arguments) > 1 {
			if len(req.Arguments) < 4 {
				return fmt.Errorf("not enough arguments, expect: %v, actual:%v", 4, len(req.Arguments))
			}
			renterId, err = peer.Decode(req.Arguments[1])
			if err != nil {
				return err
			}
			offlineSigning = true
		}
		err = backoff.Retry(func() error {
			peersLen := len(ctxParams.N.PeerHost.Network().Peers())
			if peersLen <= 0 {
				err = errors.New("failed to find any peer in table")
				log.Error(err)
				return err
			}
			return nil
		}, helper.WaitingForPeersBo)

		var shardHashes []string
		var fileSize int64
		var shardSize int64

		// token: parse token argument
		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}
		fmt.Println("token =", token, tokenStr)

		fileHash := req.Arguments[0]
		shardHashes, fileSize, shardSize, err = helper.GetShardHashes(ctxParams, fileHash)
		if len(shardHashes) == 0 && fileSize == -1 && shardSize == -1 &&
			strings.HasPrefix(err.Error(), "invalid hash: file must be reed-solomon encoded") {
			if copyNum, ok := req.Options[copyName].(int); ok {
				shardHashes, fileSize, shardSize, err = helper.GetShardHashesCopy(ctxParams, fileHash, copyNum)
				fmt.Printf("copy get, shardHashes:%v fileSize:%v, shardSize:%v, copy:%v err:%v \n",
					shardHashes, fileSize, shardSize, copyNum, err)
			}
		}
		if err != nil {
			return err
		}
		_, storageLength, err := helper.GetPriceAndMinStorageLength(ctxParams)
		if err != nil {
			return err
		}

		// token: get new price
		priceObj, err := chain.SettleObject.OracleService.CurrentPrice(token)
		if err != nil {
			return err
		}
		price := priceObj.Int64()
		// token: get new rate
		rate, err := chain.SettleObject.OracleService.CurrentRate(token)
		if err != nil {
			return err
		}
		_, err = helper.TotalPay(shardSize, price, storageLength, rate)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		// sync hosts from hub hosts.
		if !ctxParams.Cfg.Experimental.HostsSyncEnabled {
			_ = SyncHosts(ctxParams)
		}
		hp := helper.GetHostsProvider(ctxParams, make([]string, 0))
		if mode, ok := req.Options[hostSelectModeOptionName].(string); ok {
			var hostIDs []string
			if mode == "custom" {
				if hosts, ok := req.Options[hostSelectionOptionName].(string); ok {
					hostIDs = strings.Split(hosts, ",")
				}
				if len(hostIDs) != len(shardHashes) {
					return fmt.Errorf("custom mode hosts length must match shard hashes length")
				}
				hp = helper.GetCustomizedHostsProvider(ctxParams, hostIDs)
			}
		}
		rss, err := sessions.GetRenterSessionWithToken(ctxParams, ssId, fileHash, shardHashes, token)
		if err != nil {
			return err
		}
		if offlineSigning {
			offNonceTimestamp, err := strconv.ParseUint(req.Arguments[2], 10, 64)
			if err != nil {
				return err
			}
			err = rss.SaveOfflineMeta(&renterpb.OfflineMeta{
				OfflinePeerId:    req.Arguments[1],
				OfflineNonceTs:   offNonceTimestamp,
				OfflineSignature: req.Arguments[3],
			})
			if err != nil {
				return err
			}
		}
		shardIndexes := make([]int, 0)
		for i, _ := range rss.ShardHashes {
			shardIndexes = append(shardIndexes, i)
		}
		UploadShard(rss, hp, price, token, shardSize, storageLength, offlineSigning, renterId, fileSize, shardIndexes, nil)
		seRes := &Res{
			ID: ssId,
		}
		return res.Emit(seRes)
	},
	Type: Res{},
}

func SyncHosts(ctxParams *helper.ContextParams) error {
	cfg, err := ctxParams.N.Repo.Config()
	if err != nil {
		log.Errorf("Failed to get configuration %s", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	m := cfg.Experimental.HostsSyncMode
	_, err = hosts.SyncHosts(ctx, ctxParams.N, m)
	return err
}

type Res struct {
	ID string
}
