package upload

import (
	"errors"
	"strings"
	"time"

	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/protos/metadata"
	"github.com/bittorrent/go-btfs/utils"
	"github.com/libp2p/go-libp2p/core/peer"

	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"

	cmds "github.com/bittorrent/go-btfs-cmds"
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
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}
		fileHash := req.Arguments[0]
		if err != nil {
			return err
		}

		var meta *metadata.FileMetaInfo
		meta, err = chain.SettleObject.FileMetaService.GetFileMeta(fileHash, []string{})
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		contracts := meta.Contracts
		if len(contracts) <= 0 {
			return errors.New("length of contracts is 0")
		}
		ssId, _ := uh.SplitContractId(contracts[0].Meta.ContractId)
		shardIndexes := make([]int, 0)
		i := 0
		shardHashes := strings.Split(req.Arguments[1], ",")
		for _, contract := range contracts {
			if contract.Meta.ShardHash == shardHashes[i] {
				shardIndexes = append(shardIndexes, int(contract.Meta.ShardIndex))
				i++
			}
		}
		rss, err := sessions.GetRenterSession(ctxParams, ssId, fileHash, shardHashes)
		if err != nil {
			return err
		}
		hp := uh.GetSPsProvider(ctxParams, strings.Split(req.Arguments[3], ","))
		m := contracts[0].Meta
		renterPid, err := peer.Decode(req.Arguments[2])
		if err != nil {
			return err
		}

		// token: notice repair is dropped. This is just a compatible function of 'UploadShard'.
		UploadShard(&ShardUploadContext{
			Rss:           rss,
			HostsProvider: hp,
			Price:         int64(m.Price),
			Token:         tokencfg.GetWbttToken(),
			ShardSize:     int64(m.ShardSize),
			StorageLength: -1,
			ShardIndexes:  shardIndexes,
			RepairParams: &RepairParams{
				RenterStart: time.Unix(int64(m.StorageStart), 0),
				RenterEnd:   time.Unix(int64(m.StorageEnd), 0),
			},
			RenterId: renterPid,
		})
		seRes := &Res{
			ID: ssId,
		}
		return res.Emit(seRes)
	},
	Type: Res{},
}
