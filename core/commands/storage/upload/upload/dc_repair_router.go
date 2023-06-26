package upload

import (
	"context"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bittorrent/go-btfs/core/commands/cmdenv"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/core/corehttp/remote"
	"github.com/bittorrent/go-btfs/core/corerepo"

	cmds "github.com/bittorrent/go-btfs-cmds"

	"github.com/bittorrent/go-btfs-common/crypto"
	guardpb "github.com/bittorrent/go-btfs-common/protos/guard"
	"github.com/bittorrent/go-btfs-common/utils/grpc"
	"github.com/bittorrent/protobuf/proto"

	"github.com/alecthomas/units"
	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p/core/peer"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type RepairContractParams struct {
	FileHash             string
	FileSize             int64
	RepairPid            string
	LostShardHashes      []string
	RepairRewardAmount   int64
	DownloadRewardAmount int64
}

const (
	requestInterval = 5 * time.Minute
	repairerTimeout = 60 * time.Minute
)

var (
	requestFileHash string
	logger          = logging.Logger("repair")
)

var StorageDcRepairRouterCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Interact with Host repair requests and responses for negotiation of decentralized shards repair.",
		ShortDescription: `
Guard broadcasts repair request to potential nodes, after negotiation, nodes prepare the contracts for 
such repair job, sign and send to the guard for confirmation `,
	},
	Subcommands: map[string]*cmds.Command{
		"request":  hostRepairRequestCmd,
		"response": HostRepairResponseCmd,
	},
}

var hostRepairRequestCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Negotiate with hosts for repair jobs.",
		ShortDescription: `
This command sends request to mining host to negotiate the repair works.`,
	},
	Arguments: append([]cmds.Argument{
		cmds.StringArg("peer-ids", true, false, "Host Peer IDs to send repair requests."),
	}, HostRepairResponseCmd.Arguments...), // append pass-through arguments
	RunTimeout: 20 * time.Second,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		fileHash := req.Arguments[1]
		lostShardHashes := strings.Split(req.Arguments[2], ",")
		fileSize, err := strconv.ParseInt(req.Arguments[3], 10, 64)
		if err != nil {
			logger.Errorf("file size parse error: [%v]", err)
			return err
		}
		downloadRewardAmount, err := strconv.ParseInt(req.Arguments[4], 10, 64)
		if err != nil {
			logger.Errorf("download reward amount parse error: [%v]", err)
			return err
		}
		repairRewardAmount, err := strconv.ParseInt(req.Arguments[5], 10, 64)
		if err != nil {
			logger.Errorf("repair reward amount parse error: [%v]", err)
			return err
		}
		err = emptyCheck(fileHash, lostShardHashes, fileSize, downloadRewardAmount, repairRewardAmount)
		if err != nil {
			logger.Errorf("required field is empty: [%v]", err)
			return err
		}
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}
		api, err := cmdenv.GetApi(env, req)
		if err != nil {
			return err
		}
		peerIds := strings.Split(req.Arguments[0], ",")
		repairIds := make([]string, len(peerIds))
		var wg sync.WaitGroup
		wg.Add(len(peerIds))
		for index, peerId := range peerIds {
			go func() {
				err = func() error {
					pi, err := remote.FindPeer(req.Context, n, peerId)
					if err != nil {
						return err
					}
					_, err = remote.P2PCallStrings(req.Context, n, api, pi.ID, "/storage/dcrepair/response",
						req.Arguments[1:]...)
					if err != nil {
						return err
					}
					repairIds[index] = peerId
					return nil
				}()
				if err != nil {
					logger.Error("P2P call error", zap.Error(err), zap.String("peerId", peerId))
				}
				wg.Done()
			}()
		}
		wg.Wait()
		return cmds.EmitOnce(res, &peerIdList{repairIds})
	},
	Type: peerIdList{},
}

type peerIdList struct {
	Strings []string
}

var HostRepairResponseCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Host responds to repair jobs.",
		ShortDescription: `
This command enquires the repairer with the contract, if agrees with the contract after negotiation,
returns the repairer's signed contract to the invoker.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("file-hash", true, false, "File hash for the host prepares to repair."),
		cmds.StringArg("repair-shard-hash", true, false, "Shard hash for the host prepares to repair."),
		cmds.StringArg("file-size", true, false, "Size of the repair file."),
		cmds.StringArg("download-reward-amount", true, false, "Reward amount for download workload."),
		cmds.StringArg("repair-reward-amount", true, false, "Reward amount for repair workload."),
	},
	RunTimeout: 1 * time.Minute,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		ctxParams, err := uh.ExtractContextParams(req, env)
		repairId := ctxParams.N.Identity.Pretty()
		fileHash := req.Arguments[0]
		if requestFileHash == fileHash {
			return fmt.Errorf("file {%s} has been repairing on the host {%s}", fileHash, repairId)
		}
		requestFileHash = fileHash
		lostShardHashes := strings.Split(req.Arguments[1], ",")
		fileSize, err := strconv.ParseInt(req.Arguments[2], 10, 64)
		if err != nil {
			logger.Errorf("file size parse error: [%v]", err)
			return err
		}
		downloadRewardAmount, err := strconv.ParseInt(req.Arguments[3], 10, 64)
		if err != nil {
			logger.Errorf("download reward amount parse error: [%v]", err)
			return err
		}
		repairRewardAmount, err := strconv.ParseInt(req.Arguments[4], 10, 64)
		if err != nil {
			logger.Errorf("repair reward amount parse error: [%v]", err)
			return err
		}
		err = emptyCheck(fileHash, lostShardHashes, fileSize, downloadRewardAmount, repairRewardAmount)
		if err != nil {
			logger.Errorf("required field is empty: [%v]", err)
			return err
		}
		if !ctxParams.Cfg.Experimental.StorageHostEnabled {
			return fmt.Errorf("storage host api not enabled")
		}
		if !ctxParams.Cfg.Experimental.HostRepairEnabled {
			return fmt.Errorf("storage repair api not enabled")
		}
		sizeStat, err := corerepo.RepoSize(req.Context, ctxParams.N)
		if err != nil {
			return err
		}
		remainStorage := sizeStat.StorageMax - sizeStat.RepoSize
		if fileSize > int64(remainStorage) {
			errMsg := fmt.Sprintf("remaining storage space [%d] is less than required storage size [%d]", remainStorage, fileSize)
			logger.Errorf(errMsg)
			return fmt.Errorf(errMsg)
		}
		ns, err := helper.GetHostStorageConfig(ctxParams.Ctx, ctxParams.N)
		if err != nil {
			logger.Errorf("get host storage config error: [%v]", err)
			return err
		}
		var price uint64
		if ns.RepairCustomizedPricing {
			price = ns.RepairPriceCustomized
		} else {
			price = ns.RepairPriceDefault
		}
		totalPrice := uint64(float64(fileSize) / float64(units.GiB) * float64(price))
		if totalPrice <= 0 {
			totalPrice = 1
		}
		if totalPrice > uint64(repairRewardAmount) {
			errMsg := fmt.Sprintf("host desired price [%d] is greater than request price [%d]", totalPrice, repairRewardAmount)
			logger.Errorf(errMsg)
			return fmt.Errorf(errMsg)
		}
		err = doRepair(ctxParams, res, &RepairContractParams{
			FileHash:             fileHash,
			FileSize:             fileSize,
			RepairPid:            repairId,
			LostShardHashes:      lostShardHashes,
			RepairRewardAmount:   repairRewardAmount,
			DownloadRewardAmount: downloadRewardAmount,
		})

		if err != nil {
			logger.Errorf("repair job failed for peer id [%s]: [%v]", repairId, err)
			return err
		}

		logger.Info("repair lost shards done", zap.String("peerId", ctxParams.N.Identity.Pretty()))
		return nil
	},
}

func doRepair(ctxParams *uh.ContextParams, res cmds.ResponseEmitter, params *RepairContractParams) error {
	ctx, cf := context.WithTimeout(ctxParams.Ctx, repairerTimeout)
	defer cf()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		repairContractResp, err := submitSignedRepairContract(ctxParams, params)
		if err != nil {
			logger.Error("submit repair contract error", zap.Error(err), zap.String("peerId", ctxParams.N.Identity.Pretty()))
			return err
		}
		if repairContractResp.Status == guardpb.RepairContractResponse_BOTH_SIGNED {
			repairContract := repairContractResp.Contract
			arguments := []string{repairContract.FileHash}
			signedContractID, err := signContractID(repairContract.RepairContractId, ctxParams.N.PrivateKey)
			if err != nil {
				return err
			}
			paidIn := make(chan bool)
			go checkPaymentFromClient(ctxParams, paidIn, signedContractID)
			paid := <-paidIn
			if !paid {
				logger.Debugf("contract is not paid", zap.String("contractId", repairContract.RepairContractId))
				return fmt.Errorf("contract is not paid: %s", repairContract.RepairContractId)
			}
			lostShards := strings.Join(repairContract.LostShardHash, ",")
			err = cmdenv.DownloadAndRebuildFile(ctxParams.Req, res, ctxParams.Api, repairContract.FileHash, lostShards)
			if err != nil {
				logger.Debugf("download and rebuild file error", zap.String("lostShards", lostShards))
				return err
			}
			err = challengeLostShards(ctxParams, repairContract)
			if err != nil {
				logger.Errorf("challenge lost shards error: [%v]", err)
				return err
			}
			fileStatus, err := uploadShards(ctxParams, repairContract)
			if err != nil {
				logger.Errorf("repairer uploads shards error: [%v]", err)
				return err
			}
			if fileStatus == nil {
				cmdenv.RmPin(ctxParams.Req, res, ctxParams.N, ctxParams.Api, ctxParams.Cfg, arguments, true, false)
				return nil
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			err = submitFileStatus(ctx, ctxParams.Cfg, fileStatus)
			if err != nil {
				logger.Errorf("submit file store status error: [%v]", err)
				return err
			}
			cmdenv.RmPin(ctxParams.Req, res, ctxParams.N, ctxParams.Api, ctxParams.Cfg, arguments, true, false)
		}
		return nil
	})
	return eg.Wait()
}

func submitSignedRepairContract(ctxParams *uh.ContextParams, params *RepairContractParams) (*guardpb.RepairContractResponse, error) {
	repairReq := &guardpb.RepairContract{
		FileHash:             params.FileHash,
		FileSize:             params.FileSize,
		RepairPid:            params.RepairPid,
		LostShardHash:        params.LostShardHashes,
		RepairSignTime:       time.Now().UTC(),
		RepairRewardAmount:   params.RepairRewardAmount,
		DownloadRewardAmount: params.DownloadRewardAmount,
	}
	sig, err := crypto.Sign(ctxParams.N.PrivateKey, repairReq)
	if err != nil {
		return nil, err
	}
	repairReq.RepairSignature = sig
	var repairResp *guardpb.RepairContractResponse
	err = grpc.GuardClient(ctxParams.Cfg.Services.GuardDomain).WithContext(ctxParams.Ctx, func(ctx context.Context,
		client guardpb.GuardServiceClient) error {
		repairResp, err = client.SubmitRepairContract(ctx, repairReq)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return repairResp, nil
}

func uploadShards(ctxParams *uh.ContextParams, repairContract *guardpb.RepairContract) (*guardpb.FileStoreStatus, error) {
	repairContractReq := &guardpb.RequestRepairContracts{
		FileHash:       repairContract.FileHash,
		RepairNode:     repairContract.RepairPid,
		RepairSignTime: repairContract.RepairSignTime,
	}
	repairSignature, err := crypto.Sign(ctxParams.N.PrivateKey, repairContractReq)
	if err != nil {
		return nil, err
	}
	repairContractReq.RepairSignature = repairSignature
	var statusMeta *guardpb.FileStoreStatus
	var repairContractResp *guardpb.ResponseRepairContracts
	tick := time.Tick(requestInterval)
	ctxParams.Ctx = context.Background()

FOR:
	for true {
		select {
		case <-tick:
			err = grpc.GuardClient(ctxParams.Cfg.Services.GuardDomain).WithContext(ctxParams.Ctx, func(ctx context.Context,
				client guardpb.GuardServiceClient) error {
				repairContractResp, err = client.RequestForRepairContracts(ctx, repairContractReq)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			if repairContractResp.State == guardpb.ResponseRepairContracts_REQUEST_AGAIN {
				logger.Info(fmt.Sprintf("request repair contract again in %d minutes", requestInterval))
				continue
			}
			if repairContractResp.State == guardpb.ResponseRepairContracts_CONTRACT_READY {
				statusMeta = repairContractResp.Status
				break FOR
			}
			if repairContractResp.State == guardpb.ResponseRepairContracts_DOWNLOAD_NOT_DONE {
				logger.Info("download and challenge can not be completed for lost shards", zap.String("repairerId", repairContract.RepairPid))
				return nil, nil
			}
			if repairContractResp.State == guardpb.ResponseRepairContracts_CONTRACT_CLOSED {
				logger.Info("repair contract has been closed", zap.String("contractId", repairContract.RepairContractId))
				return nil, nil
			}
		}
	}

	if statusMeta == nil {
		return nil, fmt.Errorf("file store status is nil")
	}

	contracts := statusMeta.Contracts
	shardHashes := repairContract.LostShardHash
	if len(contracts) == 0 {
		return nil, errors.New("length of contracts is 0")
	}
	if len(contracts) != len(shardHashes) {
		return nil, errors.New("contracts size is not equal to lost shard hashes count")
	}

	ssId := uuid.New().String()
	rss, err := sessions.GetRenterSession(ctxParams, ssId, repairContract.FileHash, shardHashes)
	if err != nil {
		return nil, err
	}
	hp := uh.GetHostsProvider(ctxParams, make([]string, 0))
	shardMap := make(map[int]string)
	ctx, _ := context.WithTimeout(rss.Ctx, 10*time.Minute)
	for _, contract := range contracts {
		if isContained(shardHashes, contract.ShardHash) {
			shardMap[int(contract.ShardIndex)] = contract.ShardHash
			downloadAndSignContracts(contract, rss, ctx, hp)
		}
	}

	updatedContracts, err := retrieveUpdatedContracts(shardMap, rss)
	if err != nil {
		logger.Errorf("retrieve updated contracts error: [%v]", err)
		return nil, err
	}
	fileStatus, err := NewFileStatus(updatedContracts, rss.CtxParams.Cfg, updatedContracts[0].ContractMeta.RenterPid, rss.Hash, statusMeta.FileSize)
	if err != nil {
		logger.Errorf("fail to create file store status error: [%v]", err)
		return nil, err
	}
	sign, err := crypto.Sign(rss.CtxParams.N.PrivateKey, &fileStatus.FileStoreMeta)
	if err != nil {
		return nil, err
	}
	fileStatus.PreparerSignature = sign
	return fileStatus, nil
}

func downloadAndSignContracts(contract *guardpb.Contract, rss *sessions.RenterSession, ctx context.Context, hp uh.IHostsProvider) {
	go func() {
		err := backoff.Retry(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				break
			}
			price := contract.ContractMeta.Price
			shardHash := contract.ContractMeta.ShardHash
			shardIndex := contract.ContractMeta.ShardIndex
			shardFileSize := contract.ContractMeta.ShardFileSize
			host, err := hp.NextValidHost()
			if err != nil {
				terr := rss.To(sessions.RssToErrorEvent, err)
				if terr != nil {
					logger.Debugf("original err: %s, transition err: %s", err.Error(), terr.Error())
				}
				return nil
			}
			repairPid := rss.CtxParams.N.Identity
			contract.ContractMeta.HostPid = host
			contract.State = guardpb.Contract_RENEWED
			contract.PreparerPid = repairPid.Pretty()
			sign, err := crypto.Sign(rss.CtxParams.N.PrivateKey, &contract.ContractMeta)
			if err != nil {
				return err
			}
			contract.PreparerSignature = sign
			contract.RenterSignature = nil
			contract.GuardSignature = nil
			guardContractBytes, err := proto.Marshal(contract)
			if err != nil {
				return err
			}

			hostPid, err := peer.Decode(host)
			if err != nil {
				logger.Errorf("shard %s decodes host_pid error: %s", shardHash, err.Error())
				return err
			}
			cb := make(chan error)
			ShardErrChanMap.Set(contract.ContractMeta.ContractId, cb)
			go func() {
				newCtx, _ := context.WithTimeout(ctx, 10*time.Second)
				_, err := remote.P2PCall(newCtx, rss.CtxParams.N, rss.CtxParams.Api, hostPid, "/storage/upload/init",
					rss.SsId,
					rss.Hash,
					shardHash,
					price,
					"",
					guardContractBytes,
					-1,
					shardFileSize,
					shardIndex,
					repairPid,
				)
				if err != nil {
					cb <- err
					logger.Errorf("init P2P call to host id [%s] with error: [%v]", host, err)
				}
			}()
			tick := time.Tick(30 * time.Second)
			select {
			case err = <-cb:
				ShardErrChanMap.Remove(contract.ContractMeta.ContractId)
				return err
			case <-tick:
				return errors.New("host timeout")
			}
		}, uh.HandleShardBo)
		if err != nil {
			logger.Errorf("fail to setup contract for shard [%s]: [%v]", contract.ShardHash, err)
			_ = rss.To(sessions.RssToErrorEvent,
				errors.New("timeout: fail to setup contract in "+uh.HandleShardBo.MaxElapsedTime.String()))
		}
	}()
}

func retrieveUpdatedContracts(shardMap map[int]string, rss *sessions.RenterSession) ([]*guardpb.Contract, error) {
	tick := time.Tick(5 * time.Second)
	updatedContracts := make([]*guardpb.Contract, 0)
	ctx, cancel := context.WithTimeout(rss.Ctx, 10*time.Minute)
	defer cancel()
	for true {
		select {
		case <-tick:
			remainNum, queryContracts, err := queryUpdatedContracts(shardMap, rss)
			if err != nil {
				logger.Errorf("query updated contracts error: [%v]", err)
				return nil, err
			}
			updatedContracts = append(updatedContracts, queryContracts...)
			if remainNum == 0 {
				return updatedContracts, nil
			}
		case <-ctx.Done():
			logger.Infof("session %s done", rss.SsId)
			return updatedContracts, nil
		}
	}
	return updatedContracts, nil
}

func queryUpdatedContracts(shardMap map[int]string, rss *sessions.RenterSession) (int, []*guardpb.Contract, error) {
	contractStatus := "contract"
	queryContracts := make([]*guardpb.Contract, 0)
	for index, shardHash := range shardMap {
		shard, err := sessions.GetRenterShard(rss.CtxParams, rss.SsId, shardHash, index)
		if err != nil {
			continue
		}
		s, err := shard.Status()
		if err != nil {
			return -1, nil, err
		}
		if s.Status == contractStatus {
			contracts, err := shard.Contracts()
			if err != nil {
				return -1, nil, err
			}
			guardContract := contracts.SignedGuardContract
			if guardContract != nil {
				queryContracts = append(queryContracts, guardContract)
				delete(shardMap, index)
			}
		}
	}
	return len(shardMap), queryContracts, nil
}

/*func downloadAndRebuildFile(ctxParams *uh.ContextParams, fileHash string, repairShards []string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	url := fmt.Sprint(strings.Split(ctxParams.Cfg.Addresses.API[0], "/")[2], ":", strings.Split(ctxParams.Cfg.Addresses.API[0], "/")[4])
	resp, err := shell.NewShell(url).Request("get", fileHash).Option("rs", strings.Join(repairShards, ",")).Send(ctx)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	if resp.Error != nil {
		return nil, resp.Error
	}
	return ioutil.ReadAll(resp.Output)
}*/

/*func unpinLocalStorage(ctxParams *uh.ContextParams, fishHash string) error {
	url := fmt.Sprint(strings.Split(ctxParams.Cfg.Addresses.API[0], "/")[2], ":", strings.Split(ctxParams.Cfg.Addresses.API[0], "/")[4])
	err := shell.NewShell(url).Request("pin/rm", fishHash).Option("r", true).Exec(ctxParams.Ctx, nil)
	if err != nil {
		return err
	}
	return nil
}*/

func challengeLostShards(ctxParams *uh.ContextParams, repairReq *guardpb.RepairContract) error {
	eg, _ := errgroup.WithContext(ctxParams.Ctx)
	for _, lostShardHash := range repairReq.LostShardHash {
		shardHash := lostShardHash
		eg.Go(func() error {
			fileHash := repairReq.FileHash
			return challengeShard(ctxParams, fileHash, true, &guardpb.ContractMeta{
				//RenterPid:   repairReq.RepairPid,
				FileHash:   fileHash,
				ShardHash:  shardHash,
				ContractId: repairReq.RepairContractId,
				HostPid:    repairReq.RepairPid,
				RentStart:  repairReq.RepairSignTime,
			})
		})
	}
	return eg.Wait()
}

func emptyCheck(fileHash string, lostShardHashes []string, fileSize int64, DownloadRewardAmount int64, RepairRewardAmount int64) error {
	if fileHash == "" {
		return fmt.Errorf("file Hash is empty")
	}
	if fileSize == 0 {
		return fmt.Errorf("file size is 0")
	}
	if len(lostShardHashes) == 0 {
		return fmt.Errorf("lost shard hashes do not specify")
	}
	if DownloadRewardAmount <= 0 {
		return fmt.Errorf("download reward amount should be bigger than 0")
	}
	if RepairRewardAmount <= 0 {
		return fmt.Errorf("repair reward amount should be bigger than 0")
	}
	return nil
}

func isContained(lostShardHashes []string, shardHash string) bool {
	for _, lostShardHash := range lostShardHashes {
		if lostShardHash == shardHash {
			return true
		}
	}
	return false
}
