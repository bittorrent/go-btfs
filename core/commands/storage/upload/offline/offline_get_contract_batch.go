package offline

import (
	"fmt"
	"github.com/bittorrent/go-btfs/utils"

	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"

	cmds "github.com/bittorrent/go-btfs-cmds"

	cmap "github.com/orcaman/concurrent-map"
)

const (
	contractsTypeGuard = "guard"
)

var StorageUploadGetContractBatchCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get all the contracts from the upload session	(From BTFS SDK application's perspective).",
		ShortDescription: `
This command (on client) reads the unsigned contracts and returns 
the contracts to the caller.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("session-id", true, false, "ID for the entire storage upload session."),
		cmds.StringArg("peer-id", true, false, "Offline signs needed for this particular client."),
		cmds.StringArg("nonce-timestamp", true, false, "Nonce timestamp string for this upload signing."),
		cmds.StringArg("upload-session-signature", true, false, "Private key-signed string of <peer-id>:<nonce-timstamp>"),
		cmds.StringArg("contracts-type", true, false, "get guard or escrow contracts"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		ssId := req.Arguments[0]
		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}
		if !ctxParams.Cfg.Experimental.StorageClientEnabled {
			return fmt.Errorf("storage client api not enabled")
		}
		rss, err := sessions.GetRenterSession(ctxParams, ssId, "", make([]string, 0))
		if err != nil {
			return err
		}
		status, err := rss.Status()
		if err != nil {
			return err
		}
		err = verifyReceivedMessage(req, rss)
		if err != nil {
			return err
		}
		contracts := make([]*contract, 0)
		var cm cmap.ConcurrentMap
		if cm = uh.EscrowContractMaps; req.Arguments[4] == contractsTypeGuard {
			cm = uh.GuardContractMaps
		}
		for i, h := range status.ShardHashes {
			shardId := sessions.GetShardId(ssId, h, i)
			c := &contract{
				Key: shardId,
			}
			if bytes, ok := cm.Get(shardId); ok {
				c.ContractData, err = helper.BytesToString(bytes.([]byte), helper.Base64)
			} else {
				return res.Emit(&getContractBatchRes{
					Contracts: make([]*contract, 0),
				})
			}
			contracts = append(contracts, c)
		}
		return res.Emit(&getContractBatchRes{
			Contracts: contracts,
		})
	},
	Type: getContractBatchRes{},
}

type getContractBatchRes struct {
	Contracts []*contract
}

type contract struct {
	Key          string `json:"key"`
	ContractData string `json:"contract"`
}
