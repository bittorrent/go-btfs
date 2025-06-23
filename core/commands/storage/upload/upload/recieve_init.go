package upload

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/bittorrent/go-btfs/protos/metadata"
	"github.com/bittorrent/go-btfs/utils"

	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/ethereum/go-ethereum/common"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/core/commands/rm"
	"github.com/bittorrent/go-btfs/core/commands/storage/challenge"
	"github.com/bittorrent/go-btfs/core/commands/storage/helper"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	"github.com/bittorrent/go-btfs/core/corehttp/remote"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs-common/crypto"
	"github.com/bittorrent/protobuf/proto"

	"github.com/alecthomas/units"
	"github.com/cenkalti/backoff/v4"
	cidlib "github.com/ipfs/go-cid"
	ic "github.com/libp2p/go-libp2p/core/crypto"
)

var StorageUploadInitCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Initialize storage handshake with inquiring client.",
		ShortDescription: `
Storage host opens this endpoint to accept incoming upload/storage requests,
If current host is interested and all validation checks out, host downloads
the shard and replies back to client for the next challenge step.`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("session-id", true, false, "ID for the entire storage upload session."),
		cmds.StringArg("file-hash", true, false, "Root file storage node should fetch (the DAG)."),
		cmds.StringArg("shard-hash", true, false, "Shard the storage node should fetch."),
		cmds.StringArg("price", true, false, "Per GiB per day in µBTT (=0.000001BTT) for storing this shard offered by client."),
		cmds.StringArg("contract-meta", true, false, "Client's initial contract meta."),
		cmds.StringArg("storage-length", true, false, "Store file for certain length in days."),
		cmds.StringArg("shard-size", true, false, "Size of each shard received in bytes."),
		cmds.StringArg("shard-index", true, false, "Index of shard within the encoding scheme."),
		cmds.StringArg("upload-peer-id", false, false, "Peer id when upload sign is used."),
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

		if !ctxParams.Cfg.Experimental.StorageHostEnabled {
			return fmt.Errorf("storage host api not enabled")
		}
		requestPid, ok := remote.GetStreamRequestRemotePeerID(req, ctxParams.N)
		if !ok {
			return fmt.Errorf("fail to get peer ID from request")
		}

		// if my vault is not compatible with the peer's one, reject uploading
		myPeerId, err := peer.Decode(ctxParams.Cfg.Identity.PeerID)
		if err != nil {
			return fmt.Errorf("parse your peerId error: %v", err)
		}

		isVaultCompatible, err := chain.SettleObject.Factory.IsVaultCompatibleBetween(ctxParams.Ctx, myPeerId, requestPid)
		if err != nil {
			return err
		}
		if !isVaultCompatible {
			return fmt.Errorf("vault factory not compatible, please upgrade your node if possible")
		}

		// reject contract if holding contracts is above threshold
		hm := NewHostManager(ctxParams.Cfg)
		shardSize, err := strconv.ParseInt(req.Arguments[6], 10, 64)
		if err != nil {
			return err
		}

		accept, err := hm.AcceptContract(ctxParams.N.Repo.Datastore(), ctxParams.N.Identity.String(), shardSize)
		if err != nil {
			return err
		}
		if !accept {
			return errors.New("too many initialized contracts")
		}
		_, err = strconv.ParseInt(req.Arguments[3], 10, 64)
		if err != nil {
			return err
		}
		settings, err := helper.GetHostStorageConfig(ctxParams.Ctx, ctxParams.N)
		if err != nil {
			return err
		}

		storeLen, err := strconv.Atoi(req.Arguments[5])
		if err != nil {
			return err
		}
		if uint64(storeLen) < settings.StorageTimeMin {
			return fmt.Errorf("storage length invalid: want: >=%d, got: %d", settings.StorageTimeMin, storeLen)
		}
		ssId := req.Arguments[0]
		shardHash := req.Arguments[2]
		shardIndex, err := strconv.Atoi(req.Arguments[7])
		if err != nil {
			return err
		}

		fmt.Printf("upload init: start, shardSize:%v, requestPid:%v, shardIndex:%v . \n",
			shardSize, requestPid, shardIndex)

		halfSignedContractString := req.Arguments[4]
		if halfSignedContractString == "" {
			return fmt.Errorf("half signed contract is empty")
		}
		halfSignedContractBytes := []byte(halfSignedContractString)
		defer func() {
			err := recover()
			if err != nil {
				fmt.Printf("upload init: panic, err:%v, shardIndex:%v, requestPid:%v. \n", err, shardIndex, requestPid)
			}
		}()
		halfSignedContract := &metadata.Contract{}
		if err = proto.Unmarshal(halfSignedContractBytes, halfSignedContract); err != nil {
			return fmt.Errorf("unmarshal half signed contract error: %v", err)
		}
		if err != nil {
			return err
		}
		contractMeta := halfSignedContract.Meta
		// get renter's public key
		pid, ok := remote.GetStreamRequestRemotePeerID(req, ctxParams.N)
		if !ok {
			return fmt.Errorf("fail to get peer ID from request")
		}
		var peerId string
		if peerId = pid.String(); len(req.Arguments) >= 10 {
			peerId = req.Arguments[8]
		}
		payerPubKey, err := crypto.GetPubKeyFromPeerId(peerId)
		if err != nil {
			return err
		}
		s := halfSignedContract.GetUserSignature()
		ok, err = crypto.Verify(payerPubKey, contractMeta, s)
		if !ok || err != nil {
			return fmt.Errorf("can't verify guard contract: %v", err)
		}

		signedContract, err := signContract(contractMeta, halfSignedContract, ctxParams.N.PrivateKey)
		if err != nil {
			return err
		}
		signedContractBytes, err := proto.Marshal(signedContract)
		if err != nil {
			return err
		}

		var price int64
		var amount int64
		var rate *big.Int
		{
			// check renter-token
			token := common.HexToAddress(halfSignedContract.Meta.Token)
			_, bl := tokencfg.MpTokenStr[token]
			if !bl {
				err = errors.New("receive upload init, your input token is not supported. " + token.String())
				return err
			}

			// check renter-price
			price = int64(contractMeta.Price)
			priceOnline, err := chain.SettleObject.OracleService.CurrentPrice(token)
			if err != nil {
				return err
			}
			fmt.Printf("receive init, token[%s] renter-price[%v], online-price[%v],  \n", token.String(), price, priceOnline)

			if price < priceOnline.Int64() {
				return fmt.Errorf("receive init, your renter-price[%v] is less than online-price[%v]",
					price, priceOnline)
			}

			// check renter-amount
			rate, err = chain.SettleObject.OracleService.CurrentRate(token)
			if err != nil {
				return err
			}
			amount = int64(contractMeta.Amount)
			amountCal, err := uh.TotalPay(int64(contractMeta.ShardSize), price, storeLen, rate)
			if err != nil {
				return err
			}
			// fmt.Printf("receive init, renter-amount[%v], cal-amount[%v] \n", amount, amountCal)
			if amount < amountCal {
				return fmt.Errorf("receive init, your renter-amount[%v] is less than cal-amount[%v]. ",
					amount, amountCal)
			}
		}

		go func() {
			tmp := func() error {
				shard, err := sessions.GetSPShard(ctxParams, contractMeta.ContractId, price, amount, rate)
				if err != nil {
					return err
				}

				_, err = remote.P2PCall(ctxParams.Ctx, ctxParams.N, ctxParams.Api, requestPid, "/storage/upload/recvcontract",
					ssId,
					shardHash,
					shardIndex,
					signedContractBytes,
				)
				if err != nil {
					return err
				}

				if err := shard.UpdateToContractStatus(signedContract); err != nil {
					return err
				}

				fileHash := req.Arguments[1]
				err = downloadShardFromClient(ctxParams, halfSignedContract, fileHash, shardHash, false)
				if err != nil {
					return err
				}

				for i := 0; i < 20; i++ {
					// get first if exist update status
					status, err := chain.SettleObject.FileMetaService.GetContractStatus(contractMeta.ContractId)
					if err != nil {
						fmt.Printf("get contract status failed, err:%v \n", err)
						continue
					}

					if status == metadata.Contract_INVALID {
						time.Sleep(30 * time.Second)
						continue
					}

					err = chain.SettleObject.FileMetaService.UpdateContractStatus(contractMeta.ContractId)
					if err != nil {
						fmt.Printf("update contract status failed, err:%v \n", err)
						time.Sleep(3 * time.Second)
						continue
					} else {
						break
					}
				}

				fmt.Printf("upload init: send /storage/upload/recvcontract ok, wait for pay status, requestPid:%v, shardIndex:%v. \n",
					requestPid, shardIndex)

				blPay := false
				var wg sync.WaitGroup
				wg.Add(1)
				go func() {
					// every 30s check pay status
					ticker := time.NewTicker(30 * time.Second)
					defer ticker.Stop()

					// total timeout for checking pay status
					timeoutPay := time.NewTimer(10 * time.Minute)
					for {
						select {
						case <-ticker.C:
							if bl := shard.IsPayStatus(); bl {
								blPay = true
								wg.Done()
								return
							}
						case <-timeoutPay.C:
							return
						}
					}
				}()
				wg.Wait()

				if blPay {
					_ = shard.UpdateContractStatus()
					// pin shardHash
					err = pinShard(ctxParams, halfSignedContract, fileHash, shardHash)
					if err != nil {
						return err
					}
					fmt.Printf("upload init: pin shard ok, requestPid:%v, shardIndex:%v. \n", requestPid, shardIndex)
				} else {
					// rm shardHash
					err = rmShard(ctxParams, req, env, shardHash)
					if err != nil {
						return err
					}
					fmt.Printf("upload init: timeout, remove Shard, requestPid:%v, shardIndex:%v. \n", requestPid, shardIndex)
				}

				fmt.Printf("upload init: Complete! requestPid:%v, shardIndex:%v. \n", requestPid, shardIndex)
				if err := shard.Complete(); err != nil {
					return err
				}

				return nil
			}()
			if tmp != nil {
				log.Debug(tmp)
			}
		}()
		return nil
	},
}

func signContract(meta *metadata.ContractMeta, cont *metadata.Contract, privKey ic.PrivKey) (*metadata.Contract, error) {
	signedBytes, err := crypto.Sign(privKey, meta)
	if err != nil {
		return cont, err
	}
	if cont == nil {
		cont = &metadata.Contract{
			Meta:       meta,
			CreateTime: uint64(time.Now().Unix()),
		}
	} else {
		cont.CreateTime = uint64(time.Now().Unix())
	}
	cont.SpSignature = signedBytes
	return cont, err
}

func pinShard(ctxParams *uh.ContextParams, guardContract *metadata.Contract, fileHash string,
	shardHash string) error {

	err := downloadShardFromClient(ctxParams, guardContract, fileHash, shardHash, true)
	if err != nil {
		return errors.New("pinShard, stale contracts clean up error:" + err.Error())
	}

	return nil
}

func rmShard(ctxParams *uh.ContextParams, req *cmds.Request, env cmds.Environment, shardHash string) error {

	_, err := rm.RmDag(context.Background(), []string{shardHash}, ctxParams.N, req, env, true)
	if err != nil {
		// may have been cleaned up already, ignore
		return errors.New("rmShard, stale contracts clean up error:" + err.Error())
	}

	return nil
}

func downloadShardFromClient(ctxParams *uh.ContextParams, guardContract *metadata.Contract, fileHash string,
	shardHash string, blPin bool) error {

	// Get + pin to make sure it does not get accidentally deleted
	// Sharded scheme as special pin logic to add
	// file root dag + shard root dag + metadata full dag + only this shard dag
	fileCid, err := cidlib.Parse(fileHash)
	if err != nil {
		return err
	}
	shardCid, err := cidlib.Parse(shardHash)
	if err != nil {
		return err
	}
	// Need to compute a time to download shard that's fair for small vs large files
	low := 30 * time.Second
	high := 5 * time.Minute
	scaled := time.Duration(float64(guardContract.Meta.ShardSize) / float64(units.GiB) * float64(high))
	if scaled < low {
		scaled = low
	} else if scaled > high {
		scaled = high
	}
	// Also need to account for renter going up and down, to give an overall retry time limit
	lowRetry := 30 * time.Minute
	highRetry := 24 * time.Hour
	scaledRetry := time.Duration(float64(guardContract.Meta.ShardSize) / float64(units.GiB) * float64(highRetry))
	if scaledRetry < lowRetry {
		scaledRetry = lowRetry
	} else if scaledRetry > highRetry {
		scaledRetry = highRetry
	}
	expir := uint64(guardContract.Meta.StorageEnd)

	err = backoff.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), scaled)
		defer cancel()
		_, err = challenge.NewStorageChallengeResponse(ctx, ctxParams.N, ctxParams.Api, fileCid, shardCid, "", blPin, expir)
		return err
	}, uh.DownloadShardBo(scaledRetry))

	if err != nil {
		return fmt.Errorf("failed to download shard %s from file %s with contract id %s: [%v]",
			guardContract.Meta.ShardHash, fileCid, guardContract.Meta.ContractId, err)
	}
	return nil
}

func setPaidStatus(ctxParams *uh.ContextParams, contractId string) error {
	shard, err := sessions.GetSPShard(ctxParams, contractId, 0, 0, new(big.Int))
	if err != nil {
		return err
	}

	if bl := shard.IsContractStatus(); bl {
		if err := shard.ReceivePayCheque(); err != nil {
			return err
		}
	}

	return nil
}

func getInputPriceAmountRate(ctxParams *uh.ContextParams, contractId string) (int64, int64, *big.Int, error) {
	shard, err := sessions.GetSPShard(ctxParams, contractId, 0, 0, new(big.Int))
	if err != nil {
		return 0, 0, new(big.Int), err
	}

	return shard.GetInputPrice(), shard.GetInputAmount(), shard.GetInputRate(), nil
}
