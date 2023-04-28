package offline

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/utils"
	"strconv"

	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"

	cmds "github.com/bittorrent/go-btfs-cmds"

	cmap "github.com/orcaman/concurrent-map"
)

var StorageUploadSignContractBatchCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get the unsigned contracts from the upload session.",
		ShortDescription: `
This command reads all the unsigned contracts from the upload session 
(From BTFS SDK application's perspective) and returns the contracts to the caller.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("session-id", true, false, "ID for the entire storage upload session."),
		cmds.StringArg("peer-id", true, false, "Offline signs needed for this particular client."),
		cmds.StringArg("nonce-timestamp", true, false, "Nonce timestamp string for this upload signing."),
		cmds.StringArg("upload-session-signature", true, false, "Private key-signed string of peer-id:nonce-timestamp"),
		cmds.StringArg("contracts-type", true, false, "get guard or escrow contracts"),
		cmds.StringArg("signed-data-items", true, false, "signed data items."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		ssID := req.Arguments[0]
		ctxParams, err := uh.ExtractContextParams(req, env)
		if err != nil {
			return err
		}
		rss, err := sessions.GetRenterSession(ctxParams, ssID, "", make([]string, 0))
		if err != nil {
			return err
		}
		err = verifyReceivedMessage(req, rss)
		if err != nil {
			return err
		}

		var signedContracts []contract
		signedContractsString := req.Arguments[5]
		err = json.Unmarshal([]byte(signedContractsString), &signedContracts)
		if err != nil {
			return err
		}
		if len(signedContracts) != len(rss.ShardHashes) {
			return fmt.Errorf("number of received signed data items %d does not match number of shards %d",
				len(signedContracts), len(rss.ShardHashes))
		}
		var cm cmap.ConcurrentMap
		if cm = uh.EscrowChanMaps; req.Arguments[4] == contractsTypeGuard {
			cm = uh.GuardChanMaps
		}
		for i := 0; i < len(rss.ShardHashes); i++ {
			shardId := signedContracts[i].Key
			ch, found := cm.Get(shardId)
			if !found {
				return fmt.Errorf("can not find an entry for key %s", shardId)
			}
			by, err := helper.StringToBytes(signedContracts[i].ContractData, helper.Base64)
			if err != nil {
				return err
			}
			ch.(chan []byte) <- by
		}
		return nil
	},
}

func verifyReceivedMessage(req *cmds.Request, rss *sessions.RenterSession) error {
	meta, err := rss.OfflineMeta()
	if err != nil {
		return err
	}
	if meta.OfflinePeerId != req.Arguments[1] {
		return errors.New("peerIDs do not match")
	}
	offlineNonceTimestamp, err := strconv.ParseUint(req.Arguments[2], 10, 64)
	if err != nil {
		return err
	}
	if meta.OfflineNonceTs != offlineNonceTimestamp {
		return errors.New("Nonce timestamps do not match")
	}
	if meta.OfflineSignature != req.Arguments[3] {
		return errors.New("Session signature do not match")
	}
	return nil
}
