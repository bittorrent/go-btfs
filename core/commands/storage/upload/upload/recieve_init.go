package upload

import (
	"context"
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/utils"
	"math/big"
	"strconv"
	"sync"
	"time"

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
	"github.com/bittorrent/go-btfs-common/ledger"
	escrowpb "github.com/bittorrent/go-btfs-common/protos/escrow"
	guardpb "github.com/bittorrent/go-btfs-common/protos/guard"
	"github.com/bittorrent/go-btfs-common/utils/grpc"
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
		cmds.StringArg("escrow-contract", true, false, "Client's initial escrow contract data."),
		cmds.StringArg("guard-contract-meta", true, false, "Client's initial guard contract meta."),
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
		isVaultCompatible, err := chain.SettleObject.Factory.IsVaultCompatibleBetween(ctxParams.Ctx, myPeerId, requestPid)
		if err != nil {
			return err
		}
		if !isVaultCompatible {
			return fmt.Errorf("vault factory not compatible, please upgrade your node if possible")
		}

		// reject contract if holding contracts is above threshold
		hm := NewHostManager(ctxParams.Cfg)
		shardSize, err := strconv.ParseInt(req.Arguments[7], 10, 64)
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
		//if uint64(price) < settings.StoragePriceAsk {
		//	return fmt.Errorf("price invalid: want: >=%d, got: %d", settings.StoragePriceAsk, price)
		//}

		storeLen, err := strconv.Atoi(req.Arguments[6])
		if err != nil {
			return err
		}
		if uint64(storeLen) < settings.StorageTimeMin {
			return fmt.Errorf("storage length invalid: want: >=%d, got: %d", settings.StorageTimeMin, storeLen)
		}
		ssId := req.Arguments[0]
		shardHash := req.Arguments[2]
		shardIndex, err := strconv.Atoi(req.Arguments[8])
		if err != nil {
			return err
		}

		fmt.Printf("--- upload init: start, shardSize:%v, requestPid:%v, shardIndex:%v . \n",
			shardSize, requestPid, shardIndex)

		//halfSignedEscrowContString := req.Arguments[4]
		halfSignedGuardContString := req.Arguments[5]

		var halfSignedGuardContBytes []byte
		halfSignedGuardContBytes = []byte(halfSignedGuardContString)
		halfSignedGuardContract := &guardpb.Contract{}
		err = proto.Unmarshal(halfSignedGuardContBytes, halfSignedGuardContract)
		if err != nil {
			return err
		}
		guardContractMeta := halfSignedGuardContract.ContractMeta
		// get renter's public key
		pid, ok := remote.GetStreamRequestRemotePeerID(req, ctxParams.N)
		if !ok {
			return fmt.Errorf("fail to get peer ID from request")
		}
		var peerId string
		if peerId = pid.String(); len(req.Arguments) >= 10 {
			peerId = req.Arguments[9]
		}
		payerPubKey, err := crypto.GetPubKeyFromPeerId(peerId)
		if err != nil {
			return err
		}
		s := halfSignedGuardContract.GetRenterSignature()
		if s == nil {
			s = halfSignedGuardContract.GetPreparerSignature()
		}
		ok, err = crypto.Verify(payerPubKey, &guardContractMeta, s)
		if !ok || err != nil {
			return fmt.Errorf("can't verify guard contract: %v", err)
		}

		signedGuardContract, err := signGuardContract(&guardContractMeta, halfSignedGuardContract, ctxParams.N.PrivateKey)
		if err != nil {
			return err
		}
		signedGuardContractBytes, err := proto.Marshal(signedGuardContract)
		if err != nil {
			return err
		}

		var price int64
		var amount int64
		var rate *big.Int
		{
			// check renter-token
			token := common.HexToAddress(halfSignedGuardContract.Token)
			_, bl := tokencfg.MpTokenStr[token]
			if !bl {
				err = errors.New("receive upload init, your input token is not supported. " + token.String())
				return err
			}

			// check renter-price
			price = guardContractMeta.Price
			priceOnline, err := chain.SettleObject.OracleService.CurrentPrice(token)
			if err != nil {
				return err
			}
			fmt.Printf("receive init, token[%s] renter-price[%v], online-price[%v],  \n", token.String(), price, priceOnline)

			if price < priceOnline.Int64() {
				return errors.New(
					fmt.Sprintf("receive init, your renter-price[%v] is less than online-price[%v]. ",
						price, priceOnline),
				)
			}

			// check renter-amount
			rate, err = chain.SettleObject.OracleService.CurrentRate(token)
			if err != nil {
				return err
			}
			amount = guardContractMeta.Amount
			amountCal, err := uh.TotalPay(guardContractMeta.ShardFileSize, price, storeLen, rate)
			if err != nil {
				return err
			}
			//fmt.Printf("receive init, renter-amount[%v], cal-amount[%v] \n", amount, amountCal)
			if amount < amountCal {
				return errors.New(
					fmt.Sprintf("receive init, your renter-amount[%v] is less than cal-amount[%v]. ",
						amount, amountCal),
				)
			}
		}

		go func() {
			tmp := func() error {
				shard, err := sessions.GetHostShard(ctxParams, signedGuardContract.ContractId, price, amount, rate)
				if err != nil {
					return err
				}

				_, err = remote.P2PCall(ctxParams.Ctx, ctxParams.N, ctxParams.Api, requestPid, "/storage/upload/recvcontract",
					ssId,
					shardHash,
					shardIndex,
					nil,
					signedGuardContractBytes,
				)
				if err != nil {
					return err
				}

				if err := shard.Contract(nil, signedGuardContract); err != nil {
					return err
				}

				fileHash := req.Arguments[1]
				err = downloadShardFromClient(ctxParams, halfSignedGuardContract, fileHash, shardHash, false)
				if err != nil {
					return err
				}

				err = challengeShard(ctxParams, fileHash, false, &guardContractMeta)
				if err != nil {
					return err
				}

				fmt.Printf("upload init: send /storage/upload/recvcontract ok, wait for pay status, requestPid:%v, shardIndex:%v. \n",
					requestPid, shardIndex)

				blPay := false
				var wg sync.WaitGroup
				wg.Add(1)
				go func() {
					// every 30s check pay status
					tick := time.Tick(30 * time.Second)

					// total timeout for checking pay status
					timeoutPay := time.NewTimer(10 * time.Minute)
					for true {
						select {
						case <-tick:
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

				if blPay == true {
					// pin shardHash
					err = pinShard(ctxParams, halfSignedGuardContract, fileHash, shardHash)
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

func challengeShard(ctxParams *uh.ContextParams, fileHash string, isRepair bool, guardContractMeta *guardpb.ContractMeta) error {
	in := &guardpb.ReadyForChallengeRequest{
		RenterPid:   guardContractMeta.RenterPid,
		FileHash:    guardContractMeta.FileHash,
		ShardHash:   guardContractMeta.ShardHash,
		ContractId:  guardContractMeta.ContractId,
		HostPid:     guardContractMeta.HostPid,
		PrepareTime: guardContractMeta.RentStart,
		IsRepair:    isRepair,
	}

	sign, err := crypto.Sign(ctxParams.N.PrivateKey, in)
	if err != nil {
		return err
	}
	in.Signature = sign
	// Need to renew another 6 mins due to downloading shard could have already made
	// req.Context obsolete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	var question *guardpb.RequestChallengeQuestion
	err = grpc.GuardClient(ctxParams.Cfg.Services.GuardDomain).WithContext(ctx,
		func(ctx context.Context, client guardpb.GuardServiceClient) error {
			for i := 0; i < 20; i++ {
				question, err = client.RequestChallenge(ctx, in)
				if err == nil {
					break
				}
				time.Sleep(30 * time.Second)
			}
			return err
		})
	if err != nil {
		return fmt.Errorf("request challenge questions error: [%v]", err)
	}
	if question == nil {
		return errors.New("question is nil")
	}

	fileHashCid, err := cidlib.Parse(fileHash)
	if err != nil {
		return err
	}
	shardHashCid, err := cidlib.Parse(question.Question.ShardHash)
	if err != nil {
		return err
	}
	sc, err := challenge.NewStorageChallengeResponse(ctx, ctxParams.N, ctxParams.Api,
		fileHashCid, shardHashCid, "", false, 0)
	if err != nil {
		return err
	}
	err = sc.SolveChallenge(int(question.Question.ChunkIndex), question.Question.Nonce)
	if err != nil {
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	resp := &guardpb.ResponseChallengeQuestion{
		Answer: &guardpb.ChallengeQuestion{
			ShardHash:    question.Question.ShardHash,
			HostPid:      question.Question.HostPid,
			ChunkIndex:   int32(sc.CIndex),
			Nonce:        sc.Nonce,
			ExpectAnswer: sc.Hash,
		},
		FileHash:    fileHash,
		HostPid:     question.Question.HostPid,
		ResolveTime: time.Now(),
		IsRepair:    isRepair,
	}

	privKey, err := ctxParams.Cfg.Identity.DecodePrivateKey("")
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(privKey, resp)
	if err != nil {
		return err
	}
	resp.Signature = sig
	err = grpc.GuardClient(ctxParams.Cfg.Services.GuardDomain).WithContext(ctx,
		func(ctx context.Context, client guardpb.GuardServiceClient) error {
			_, err := client.ResponseChallenge(ctx, resp)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		log.Debug(err)
		return fmt.Errorf("response challenge error: [%v]", err)
	}
	return nil
}

//func signEscrowContractAndMarshal(contract *escrowpb.EscrowContract, signedContract *escrowpb.SignedEscrowContract,
//	privKey ic.PrivKey) ([]byte, error) {
//	sig, err := crypto.Sign(privKey, contract)
//	if err != nil {
//		return nil, err
//	}
//	if signedContract == nil {
//		signedContract = escrow.NewSignedContract(contract)
//	}
//	signedContract.SellerSignature = sig
//	signedBytes, err := proto.Marshal(signedContract)
//	if err != nil {
//		return nil, err
//	}
//	return signedBytes, nil
//}

func signGuardContract(meta *guardpb.ContractMeta, cont *guardpb.Contract, privKey ic.PrivKey) (*guardpb.Contract, error) {
	signedBytes, err := crypto.Sign(privKey, meta)
	if err != nil {
		return nil, err
	}
	if cont == nil {
		cont = &guardpb.Contract{
			ContractMeta:   *meta,
			LastModifyTime: time.Now(),
		}
	} else {
		cont.LastModifyTime = time.Now()
	}
	cont.HostSignature = signedBytes
	return cont, err
}

//func signGuardContractAndMarshal(meta *guardpb.ContractMeta, cont *guardpb.Contract, privKey ic.PrivKey) ([]byte, error) {
//	signedBytes, err := crypto.Sign(privKey, meta)
//	if err != nil {
//		return nil, err
//	}
//
//	if cont == nil {
//		cont = &guardpb.Contract{
//			ContractMeta:   *meta,
//			LastModifyTime: time.Now(),
//		}
//	} else {
//		cont.LastModifyTime = time.Now()
//	}
//	cont.HostSignature = signedBytes
//	return proto.Marshal(cont)
//}

// call escrow service to check if payment is received or not
func checkPaymentFromClient(ctxParams *uh.ContextParams, paidIn chan bool, contractID *escrowpb.SignedContractID) {
	var err error
	paid := false
	err = backoff.Retry(func() error {
		paid, err = isPaidin(ctxParams, contractID)
		if err != nil {
			return err
		}
		if paid {
			paidIn <- true
			return nil
		}
		return errors.New("reach max retry times")
	}, uh.CheckPaymentBo)
	if err != nil {
		paidIn <- paid
	}
}

func pinShard(ctxParams *uh.ContextParams, guardContract *guardpb.Contract, fileHash string,
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

func downloadShardFromClient(ctxParams *uh.ContextParams, guardContract *guardpb.Contract, fileHash string,
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
	scaled := time.Duration(float64(guardContract.ShardFileSize) / float64(units.GiB) * float64(high))
	if scaled < low {
		scaled = low
	} else if scaled > high {
		scaled = high
	}
	// Also need to account for renter going up and down, to give an overall retry time limit
	lowRetry := 30 * time.Minute
	highRetry := 24 * time.Hour
	scaledRetry := time.Duration(float64(guardContract.ShardFileSize) / float64(units.GiB) * float64(highRetry))
	if scaledRetry < lowRetry {
		scaledRetry = lowRetry
	} else if scaledRetry > highRetry {
		scaledRetry = highRetry
	}
	expir := uint64(guardContract.RentEnd.Unix())

	err = backoff.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), scaled)
		defer cancel()
		_, err = challenge.NewStorageChallengeResponse(ctx, ctxParams.N, ctxParams.Api, fileCid, shardCid, "", blPin, expir)
		return err
	}, uh.DownloadShardBo(scaledRetry))

	if err != nil {
		return fmt.Errorf("failed to download shard %s from file %s with contract id %s: [%v]",
			guardContract.ShardHash, guardContract.FileHash, guardContract.ContractId, err)
	}
	return nil
}

func setPaidStatus(ctxParams *uh.ContextParams, contractId string) error {
	shard, err := sessions.GetHostShard(ctxParams, contractId, 0, 0, new(big.Int))
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
	shard, err := sessions.GetHostShard(ctxParams, contractId, 0, 0, new(big.Int))
	if err != nil {
		return 0, 0, new(big.Int), err
	}

	return shard.GetInputPrice(), shard.GetInputAmount(), shard.GetInputRate(), nil
}

func isPaidin(ctxParams *uh.ContextParams, contractID *escrowpb.SignedContractID) (bool, error) {
	//var signedPayinRes *escrowpb.SignedPayinStatus
	//ctx, _ := helper.NewGoContext(ctxParams.Ctx)
	//err := grpc.EscrowClient(ctxParams.Cfg.Services.EscrowDomain).WithContext(ctx,
	//	func(ctx context.Context, client escrowpb.EscrowServiceClient) error {
	//		res, err := client.IsPaid(ctx, contractID)
	//		if err != nil {
	//			return err
	//		}
	//		err = escrow.VerifyEscrowRes(ctxParams.Cfg, res.Status, res.EscrowSignature)
	//		if err != nil {
	//			return err
	//		}
	//		signedPayinRes = res
	//		return nil
	//	})
	//if err != nil {
	//	return false, err
	//}
	//return signedPayinRes.Status.Paid, nil

	return true, nil
}

func signContractID(id string, privKey ic.PrivKey) (*escrowpb.SignedContractID, error) {
	contractID, err := ledger.NewContractID(id, privKey.GetPublic())
	if err != nil {
		return nil, err
	}
	// sign contractID
	sig, err := crypto.Sign(privKey, contractID)
	if err != nil {
		return nil, err
	}
	return ledger.NewSingedContractID(contractID, sig), nil
}
