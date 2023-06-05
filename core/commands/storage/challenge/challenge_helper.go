package challenge

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"math/big"
	"sync"

	core "github.com/bittorrent/go-btfs/core"
	coreiface "github.com/bittorrent/interface-go-btfs-core"
	options "github.com/bittorrent/interface-go-btfs-core/options"
	path "github.com/bittorrent/interface-go-btfs-core/path"
	uuid "github.com/google/uuid"
	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
)

type StorageChallenge struct {
	Ctx  context.Context
	Node *core.IpfsNode
	API  coreiface.CoreAPI

	// Internal fields for parsing chunks
	seenCIDs map[string]bool
	allCIDs  []cid.Cid // Not storing node for faster fetch
	sync.Mutex

	// Selections and generations for challenge
	ID         string  // Challenge ID (randomly generated on creation)
	RID        cid.Cid // Root hash in multihash format (initial file root)
	SID        cid.Cid // Shard hash in multihash format (selected shard on each host)
	CIndex     int     // Chunk index within each shard (selected index on each challenge request)
	CID        cid.Cid // Chunk hash in multihash format (selected chunk on each challenge request)
	Nonce      string  // Random nonce for each challenge request (uuidv4)
	Hash       string  // Generated SHA-256 hash (chunk bytes + nonce bytes) for proof-of-file-existence
	Expiration uint64  // End date of the pinned chunks
}

// newStorageChallengeHelper creates a challenge object with new ID, resolves the cid path
// and initializes underlying CIDs to be ready for challenge generation.
// When used by storage client: challengeID is "", will be randomly genenerated
// When used by storage host: challengeID is a valid uuid v4
// When host needs to get blocks and initialize the challenge for the very first time, use init=true
func newStorageChallengeHelper(ctx context.Context, node *core.IpfsNode, api coreiface.CoreAPI,
	rootHash, shardHash cid.Cid, challengeID string, init bool, expir uint64) (*StorageChallenge, error) {
	if challengeID == "" {
		chid, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		challengeID = chid.String()
	}
	sc := &StorageChallenge{
		Ctx:        ctx,
		Node:       node,
		API:        api,
		seenCIDs:   map[string]bool{},
		ID:         challengeID,
		RID:        rootHash,
		SID:        shardHash,
		Expiration: expir,
	}
	if err := sc.getAllCIDsRecursive(rootHash, init); err != nil {
		return nil, err
	}
	return sc, nil
}

func NewStorageChallenge(ctx context.Context, node *core.IpfsNode, api coreiface.CoreAPI,
	rootHash, shardHash cid.Cid) (*StorageChallenge, error) {
	return newStorageChallengeHelper(ctx, node, api, rootHash, shardHash, "", false, 0)
}

func NewStorageChallengeResponse(ctx context.Context, node *core.IpfsNode, api coreiface.CoreAPI,
	rootHash, shardHash cid.Cid, challengeID string, init bool, expir uint64) (*StorageChallenge, error) {
	return newStorageChallengeHelper(ctx, node, api, rootHash, shardHash, challengeID, init, expir)
}

// getAllCIDsRecursive traverses the full DAG to find all cids and
// stores them in allCIDs for quick challenge regeneration/retry.
// This function skips the other shard hashes at the same level as
// the selected shard hash, since all other shard hash sub-DAGs belong
// to other hosts, however, everything else is shared for (meta) info
// retrieval and DAG traversal.
func (sc *StorageChallenge) getAllCIDsRecursive(blockHash cid.Cid, init bool) error {
	ncs := string(blockHash.Bytes()) // shorter/faster key
	// Already seen
	if _, ok := sc.seenCIDs[ncs]; ok {
		return nil
	}
	// Check if can be resolved
	rp, err := sc.API.ResolvePath(sc.Ctx, path.IpfsPath(blockHash))
	if err != nil {
		return err
	}
	// Mark as seen
	sc.seenCIDs[ncs] = true
	sc.allCIDs = append(sc.allCIDs, blockHash)
	// Recurse
	links, err := sc.API.Object().Links(sc.Ctx, rp)
	if err != nil {
		return err
	}
	// Check if we are at the shard hash level - only check selected shard
	for _, l := range links {
		if l.Cid == sc.SID {
			links = []*ipld.Link{l}
			break
		}
	}
	isSelectedShard := blockHash == sc.SID
	for _, l := range links {
		// All children of selected shard won't need any pin as they will
		// be recursively pinned below
		if err := sc.getAllCIDsRecursive(l.Cid, init && !isSelectedShard); err != nil {
			return err
		}
	}
	// If shard dag, recursively pin the full dag to store relation
	// Otherwise, just singlely pin a dag
	if init {
		return sc.API.Pin().Add(sc.Ctx, rp,
			options.Pin.Recursive(isSelectedShard), options.Pin.Expiration(sc.Expiration))
	}
	return nil
}

// GenChallenge picks a random cid from the internal cid list,
// then generates a random nonce, combined for the final hash.
// It is safe to call this function multiple times to re-gen a new challenge
// while re-using the same data structure and ID.
func (sc *StorageChallenge) GenChallenge() error {
	sc.Lock()
	defer sc.Unlock()

	// Randomly select a (new) CID
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(sc.allCIDs))))
	if err != nil {
		return err
	}
	sc.CID = sc.allCIDs[n.Int64()]
	sc.CIndex = int(n.Int64())

	// Fetch the raw data
	r, _, err := sc.API.Object().Data(sc.Ctx, path.IpfsPath(sc.CID), true, false)
	if err != nil {
		return err
	}

	nonce, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	// Give user readable string
	sc.Nonce = nonce.String()

	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return err
	}
	// But use bytes for hashing
	nb := [16]byte(nonce)
	h.Write(nb[:])
	sc.Hash = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}

// SolveChallenge solves the given challenge chunk index + nonce's request
// and records the solved/responding value in Hash field.
func (sc *StorageChallenge) SolveChallenge(chIndex int, chNonce string) error {
	sc.Lock()
	defer sc.Unlock()

	// Get the cid for the chunk
	if chIndex < 0 || chIndex >= len(sc.allCIDs) {
		return fmt.Errorf("chunk index is out of range")
	}

	// Fetch the raw data
	chHash := sc.allCIDs[chIndex]
	r, _, err := sc.API.Object().Data(sc.Ctx, path.IpfsPath(chHash), true, false)
	if err != nil {
		return err
	}
	sc.CID = chHash
	sc.CIndex = chIndex

	// Decode nonce
	nonce, err := uuid.Parse(chNonce)
	if err != nil {
		return err
	}
	sc.Nonce = chNonce

	// Re-hash to solve challenge
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return err
	}
	nb := [16]byte(nonce)
	h.Write(nb[:])
	sc.Hash = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}
