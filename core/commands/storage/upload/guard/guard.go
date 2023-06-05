package guard

import (
	"context"
	"fmt"
	"time"

	"github.com/bittorrent/go-btfs/core/commands/storage/challenge"
	uh "github.com/bittorrent/go-btfs/core/commands/storage/upload/helper"
	"github.com/bittorrent/go-btfs/core/commands/storage/upload/sessions"
	renterpb "github.com/bittorrent/go-btfs/protos/renter"

	cc "github.com/bittorrent/go-btfs-common/config"
	"github.com/bittorrent/go-btfs-common/crypto"
	guardpb "github.com/bittorrent/go-btfs-common/protos/guard"
	cgrpc "github.com/bittorrent/go-btfs-common/utils/grpc"
	config "github.com/bittorrent/go-btfs-config"
	"github.com/bittorrent/protobuf/proto"

	cidlib "github.com/ipfs/go-cid"
)

const (
	GuardTimeout = 360 * time.Second
)

// PrepFileChallengeQuestions checks and prepares all shard questions in one setting
func PrepFileChallengeQuestions(rss *sessions.RenterSession, fsStatus *guardpb.FileStoreStatus,
	fileHash string, offlineSigning bool, renterPid string) (
	[]*guardpb.ShardChallengeQuestions, error) {
	var shardHashes []cidlib.Cid
	var hostIDs []string
	for _, c := range fsStatus.Contracts {
		sh, err := cidlib.Parse(c.ShardHash)
		if err != nil {
			return nil, err
		}
		shardHashes = append(shardHashes, sh)
		hostIDs = append(hostIDs, c.HostPid)
	}

	questionsPerShard := cc.GetMinimumQuestionsCountPerShard(fsStatus)
	cid, err := cidlib.Parse(fileHash)
	if err != nil {
		return nil, err
	}
	return PrepCustomFileChallengeQuestions(rss, cid, shardHashes, hostIDs, questionsPerShard, offlineSigning, renterPid)
}

// PrepCustomFileChallengeQuestions is the inner version of PrepFileChallengeQuestions without
// using a real guard file contracts, but rather custom parameters (mostly for manual testing)
func PrepCustomFileChallengeQuestions(rss *sessions.RenterSession, fileHash cidlib.Cid, shardHashes []cidlib.Cid, hostIDs []string,
	questionsPerShard int, offlineSigning bool, renterPid string) ([]*guardpb.ShardChallengeQuestions, error) {
	// safety check
	if len(hostIDs) < len(shardHashes) {
		return nil, fmt.Errorf("hosts list must be at least %d", len(shardHashes))
	}
	// generate each shard's questions individually, then combine
	questions := make([]*guardpb.ShardChallengeQuestions, len(shardHashes))
	qc := make(chan questionRes)
	for i, sh := range shardHashes {
		go func(shardIndex int, shardHash cidlib.Cid, hostID string) {
			qs, err := PrepShardChallengeQuestions(rss, fileHash, shardHash, hostID, questionsPerShard, renterPid)
			qc <- questionRes{
				qs:  qs,
				err: err,
				i:   shardIndex,
			}
		}(i, sh, hostIDs[i])
	}
	for i := 0; i < len(questions); i++ {
		res := <-qc
		if res.err != nil {
			return nil, res.err
		}
		questions[res.i] = res.qs
	}
	fileQuestions := &guardpb.FileChallengeQuestions{
		FileHash:       rss.Hash,
		ShardQuestions: questions,
	}
	cb := make(chan []byte)
	uh.QuestionsChanMaps.Set(rss.SsId, cb)
	if offlineSigning {
		bytes, err := proto.Marshal(fileQuestions)
		if err != nil {
			return nil, err
		}
		err = rss.SaveOfflineSigning(&renterpb.OfflineSigning{Raw: bytes})
		if err != nil {
			return nil, err
		}
	} else {
		go func() {
			if bytes, err := func() ([]byte, error) {
				for _, sq := range fileQuestions.ShardQuestions {
					sig, err := crypto.Sign(rss.CtxParams.N.PrivateKey, sq)
					if err != nil {
						return nil, err
					}
					sq.PreparerSignature = sig
				}
				bytes, err := proto.Marshal(fileQuestions)
				if err != nil {
					return nil, err
				}
				return bytes, nil
			}(); err != nil {
				_ = rss.To(sessions.RssToErrorEvent, err)
				return
			} else {
				cb <- bytes
			}
		}()
	}
	signBytes := <-cb
	uh.FileMetaChanMaps.Remove(rss.SsId)
	if err := rss.To(sessions.RssToGuardQuestionsSignedEvent); err != nil {
		return nil, err
	}
	f := new(guardpb.FileChallengeQuestions)
	if err := proto.Unmarshal(signBytes, f); err != nil {
		return nil, err
	}
	return f.ShardQuestions, nil
}

type questionRes struct {
	qs  *guardpb.ShardChallengeQuestions
	err error
	i   int
}

// PrepShardChallengeQuestions checks and prepares an amount of random challenge questions
// and returns the necessary guard proto struct
func PrepShardChallengeQuestions(rss *sessions.RenterSession, fileHash cidlib.Cid, shardHash cidlib.Cid,
	hostID string, numQuestions int, renterPid string) (*guardpb.ShardChallengeQuestions, error) {
	ctx := rss.Ctx
	node := rss.CtxParams.N
	api := rss.CtxParams.Api
	var shardQuestions []*guardpb.ChallengeQuestion
	var sc *challenge.StorageChallenge
	var err error

	//FIXME: cache challange info
	//if shardInfo != nil {
	//	sc, err = shardInfo.GetChallengeOrNew(ctx, node, api, fileHash)
	//	if err != nil {
	//		return nil, err
	//	}
	//} else {
	// For testing, no session concept, create new challenge
	sc, err = challenge.NewStorageChallenge(ctx, node, api, fileHash, shardHash)
	if err != nil {
		return nil, err
	}
	//}
	// Generate questions
	sh := shardHash.String()
	for i := 0; i < numQuestions; i++ {
		err := sc.GenChallenge()
		if err != nil {
			return nil, err
		}
		q := &guardpb.ChallengeQuestion{
			ShardHash:    sh,
			HostPid:      hostID,
			ChunkIndex:   int32(sc.CIndex),
			Nonce:        sc.Nonce,
			ExpectAnswer: sc.Hash,
		}
		shardQuestions = append(shardQuestions, q)
	}
	return &guardpb.ShardChallengeQuestions{
		FileHash:      fileHash.String(),
		ShardHash:     sh,
		PreparerPid:   renterPid,
		QuestionCount: int32(numQuestions),
		Questions:     shardQuestions,
		PrepareTime:   time.Now(),
	}, nil
}

// SendChallengeQuestions combines all shard questions in a file and sends to guard service
func SendChallengeQuestions(ctx context.Context, cfg *config.Config, fileHash cidlib.Cid,
	questions []*guardpb.ShardChallengeQuestions) error {
	fileQuestions := &guardpb.FileChallengeQuestions{
		FileHash:       fileHash.String(),
		ShardQuestions: questions,
	}
	return sendChallengeQuestions(ctx, cfg, fileQuestions)
}

// sendChallengeQuestions opens a grpc connection, sends questions, and closes (short) connection
func sendChallengeQuestions(ctx context.Context, cfg *config.Config, req *guardpb.FileChallengeQuestions) error {
	cb := cgrpc.GuardClient(cfg.Services.GuardDomain)
	cb.Timeout(GuardTimeout)
	return cb.WithContext(ctx, func(ctx context.Context, client guardpb.GuardServiceClient) error {
		res, err := client.SendQuestions(ctx, req)
		if err != nil {
			return err
		}
		if res.Code != guardpb.ResponseCode_SUCCESS {
			return fmt.Errorf("failed to send questions: %v", res.Message)
		}
		return nil
	})
}
