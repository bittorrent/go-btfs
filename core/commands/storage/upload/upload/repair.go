package upload

import (
	"context"
	"errors"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/tron-us/go-btfs-common/crypto"
	guardpb "github.com/tron-us/go-btfs-common/protos/guard"
	"github.com/tron-us/go-btfs-common/utils/grpc"

	"github.com/libp2p/go-libp2p/core/peer"
)

var StorageUploadRepairCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Repair specific shards of a file.",
		ShortDescription: `
This command repairs the given shards of a file.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("file-hash", true, false, "Hash of file to upload."),
		cmds.StringArg("repair-shards", true, false, "Shard hashes to repair."),
		cmds.StringArg("renter-pid", true, false, "Original renter peer ID."),
		cmds.StringArg("blacklist", true, false, "Blacklist of hosts during upload."),
	},
	RunTimeout: 5 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}
		fileHash := req.Arguments[0]
		metaReq := &guardpb.CheckFileStoreMetaRequest{
			FileHash:     fileHash,
			RenterPid:    ctxParams.N.Identity.String(),
			RequesterPid: ctxParams.N.Identity.String(),
			RequestTime:  time.Now().UTC(),
		}
		sig, err := crypto.Sign(ctxParams.N.PrivateKey, metaReq)
		if err != nil {
			return err
		}
		metaReq.Signature = sig
		ctx, _ := helper.NewGoContext(req.Context)
		var meta *guardpb.FileStoreStatus
		err = grpc.GuardClient(ctxParams.Cfg.Services.GuardDomain).WithContext(ctx, func(ctx context.Context,
			client guardpb.GuardServiceClient) error {
			meta, err = client.CheckFileStoreMeta(ctx, metaReq)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		contracts := meta.Contracts
		if len(contracts) <= 0 {
			return errors.New("length of contracts is 0")
		}
		ssId, _ := uh.SplitContractId(contracts[0].ContractId)
		shardIndexes := make([]int, 0)
		i := 0
		shardHashes := strings.Split(req.Arguments[1], ",")
		for _, contract := range contracts {
			if contract.ShardHash == shardHashes[i] {
				shardIndexes = append(shardIndexes, int(contract.ShardIndex))
				i++
			}
		}
		rss, err := sessions.GetRenterSession(ctxParams, ssId, fileHash, shardHashes)
		if err != nil {
			return err
		}
		hp := uh.GetHostsProvider(ctxParams, strings.Split(req.Arguments[3], ","))
		m := contracts[0].ContractMeta
		renterPid, err := peer.Decode(req.Arguments[2])
		if err != nil {
			return err
		}

		// token: notice repair is dropped. This is just a compatible function of 'UploadShard'.
		UploadShard(rss, hp, m.Price, tokencfg.GetWbttToken(), m.ShardFileSize, -1, false, renterPid, -1,
			shardIndexes, &RepairParams{
				RenterStart: m.RentStart,
				RenterEnd:   m.RentEnd,
			})
		seRes := &Res{
			ID: ssId,
		}
		return res.Emit(seRes)
	},
	Type: Res{},
}
