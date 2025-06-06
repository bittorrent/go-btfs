package upload

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/bittorrent/go-btfs/protos/metadata"
	"github.com/bittorrent/go-btfs/utils"

	"github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/core/corehttp/remote"

	cmds "github.com/bittorrent/go-btfs-cmds"

	"github.com/gogo/protobuf/proto"
)

var StorageUploadRecvContractCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "For renter client to receive half signed contracts.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("session-id", true, false, "Session ID which renter uses to storage all shards information."),
		cmds.StringArg("shard-hash", true, false, "Shard the storage node should fetch."),
		cmds.StringArg("shard-index", true, false, "Index of shard within the encoding scheme."),
		cmds.StringArg("agreement", true, false, "Signed agreement."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		contractId, err := doRecv(req, env)
		if contractId != "" {
			if ch, ok := ShardErrChanMap.Get(contractId); ok {
				go func() {
					ch.(chan error) <- err
				}()
			}
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func doRecv(req *cmds.Request, env cmds.Environment) (agreementId string, err error) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("receive agreement err: ", err)
		}
	}()
	ssID := req.Arguments[0]
	ctxParams, err := helper.ExtractContextParams(req, env)
	if err != nil {
		return
	}
	requestPid, ok := remote.GetStreamRequestRemotePeerID(req, ctxParams.N)
	if !ok {
		err = errors.New("failed to get remote peer id")
		return
	}
	rpk, err := requestPid.ExtractPublicKey()
	if err != nil {
		return
	}

	agreementBytes := []byte(req.Arguments[3])
	agreement := new(metadata.Contract)
	err = proto.Unmarshal(agreementBytes, agreement)
	if err != nil {
		return
	}
	bytes, err := proto.Marshal(agreement.Meta)
	if err != nil {
		return
	}

	valid, err := rpk.Verify(bytes, agreement.SpSignature)
	if err != nil {
		return
	}

	fmt.Println("receive contract valid: ", valid)
	fmt.Println("receive contract host pid: ", agreement.Meta.GetSpId())
	fmt.Println("receive contract remote peer id: ", requestPid.String())
	if !valid || agreement.Meta.GetSpId() != requestPid.String() {
		err = errors.New("invalid guard contract bytes")
		return
	}
	agreementId = agreement.Meta.ContractId

	shardHash := req.Arguments[1]
	index, err := strconv.Atoi(req.Arguments[2])
	if err != nil {
		return
	}
	shard, err := sessions.GetUserShard(ctxParams, ssID, shardHash, index)
	if err != nil {
		return
	}
	// ignore error
	_ = shard.UpdateShardToContractStatus(agreement)
	return
}
