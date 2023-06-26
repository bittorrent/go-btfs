package upload

import (
	"errors"
	"github.com/bittorrent/go-btfs/utils"
	"strconv"

	"github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/core/corehttp/remote"

	cmds "github.com/bittorrent/go-btfs-cmds"
	guardpb "github.com/bittorrent/go-btfs-common/protos/guard"

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
		cmds.StringArg("escrow-contract", true, false, "Signed Escrow contract."),
		cmds.StringArg("guard-contract", true, false, "Signed Guard contract."),
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

func doRecv(req *cmds.Request, env cmds.Environment) (contractId string, err error) {
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

	//escrowContractBytes := []byte(req.Arguments[3])
	escrowContractBytes := []byte{}
	guardContractBytes := []byte(req.Arguments[4])
	guardContract := new(guardpb.Contract)
	err = proto.Unmarshal(guardContractBytes, guardContract)
	if err != nil {
		return
	}
	bytes, err := proto.Marshal(&guardContract.ContractMeta)
	if err != nil {
		return
	}
	valid, err := rpk.Verify(bytes, guardContract.HostSignature)
	if err != nil {
		return
	}
	if !valid || guardContract.ContractMeta.GetHostPid() != requestPid.Pretty() {
		err = errors.New("invalid guard contract bytes")
		return
	}
	contractId = guardContract.ContractMeta.ContractId

	shardHash := req.Arguments[1]
	index, err := strconv.Atoi(req.Arguments[2])
	if err != nil {
		return
	}
	shard, err := sessions.GetRenterShard(ctxParams, ssID, shardHash, index)
	if err != nil {
		return
	}
	// ignore error
	_ = shard.Contract(escrowContractBytes, guardContract)
	return
}
