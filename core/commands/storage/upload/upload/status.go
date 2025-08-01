package upload

import (
	"fmt"

	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/utils"

	"github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"

	cmds "github.com/bittorrent/go-btfs-cmds"
	guardpb "github.com/bittorrent/go-btfs-common/protos/guard"

	"github.com/ipfs/go-datastore"
)

var StorageUploadStatusCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Check storage upload and payment status (From client's perspective).",
		ShortDescription: `
This command print upload and payment status by the time queried.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("session-id", true, false, "ID for the entire storage upload session.").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		status := &StatusRes{}
		// check and get session info from sessionMap
		ssId := req.Arguments[0]

		ctxParams, err := helper.ExtractContextParams(req, env)
		if err != nil {
			return err
		}

		// check if checking request from host or client
		if !ctxParams.Cfg.Experimental.StorageClientEnabled && !ctxParams.Cfg.Experimental.StorageHostEnabled {
			return fmt.Errorf("storage client/host api not enabled")
		}

		session, err := sessions.GetRenterSession(ctxParams, ssId, "", make([]string, 0))
		if err != nil {
			return err
		}
		sessionStatus, err := session.GetRenterSessionStatus()
		if err != nil {
			return err
		}
		status.Status = sessionStatus.Status
		status.Message = sessionStatus.Message
		info, err := session.GetAdditionalInfo()
		if err == nil {
			status.AdditionalInfo = info.Info
		} else {
			// NOP
		}

		// get shards info from session
		shards := make(map[string]*ShardStatus)
		status.FileHash = session.Hash
		fullyCompleted := true
		for i, h := range session.ShardHashes {
			shard, err := sessions.GetUserShard(ctxParams, ssId, h, i)
			if err != nil {
				return err
			}
			st, err := shard.GetShardStatus()
			if err != nil {
				return err
			}
			additionalInfo, err := shard.GetAdditionalInfo()
			if err != nil && err != datastore.ErrNotFound {
				return err
			}
			switch additionalInfo.Info {
			case guardpb.Contract_UPLOADED.String(), guardpb.Contract_CANCELED.String(), guardpb.Contract_CLOSED.String():
				// NOP
			default:
				fullyCompleted = false
			}
			c := &ShardStatus{
				ContractID:     "",
				Price:          0,
				Host:           "",
				Status:         st.Status,
				Message:        st.Message,
				AdditionalInfo: additionalInfo.Info,
			}
			// if contracts.SignedGuardContract != nil {
			// c.ContractID = contracts.SignedGuardContract.ContractId
			// c.Price = contracts.SignedGuardContract.Price
			// c.Host = contracts.SignedGuardContract.HostPid
			// }
			shards[sessions.GetShardId(ssId, h, i)] = c
		}
		if (status.Status == sessions.RssWaitUploadReqSignedStatus || status.Status == sessions.RssCompleteStatus) && !fullyCompleted {
			meta, err := chain.SettleObject.FileMetaService.GetFileMetaByCID(session.Hash)
			if err != nil {
				log.Debug(err)
				return err
			}
			for _, c := range meta.Contracts {
				shards[sessions.GetShardId(ssId, c.Meta.ShardHash, int(c.Meta.ShardIndex))].AdditionalInfo = c.Status.String()
			}
		}
		status.Shards = shards
		if len(status.Shards) == 0 && status.Status == sessions.RssInitStatus {
			status.Message = "session not found"
		}
		return res.Emit(status)
	},
	Type: StatusRes{},
}

type StatusRes struct {
	Status         string
	Message        string
	AdditionalInfo string
	FileHash       string
	Shards         map[string]*ShardStatus
}

type ShardStatus struct {
	ContractID     string
	Price          int64
	Host           string
	Status         string
	Message        string
	AdditionalInfo string
}
