package vault

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	cp "github.com/bittorrent/go-btfs-common/crypto"
	"github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/protos/metadata"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/bittorrent/go-btfs/transaction/crypto"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p/core/peer"
	"golang.org/x/crypto/sha3"
)

type FileMeta interface {
	AddFileMeta(cid string, meta *metadata.FileMetaInfo) error
	UpdateContractStatus(contractId string) error
	GetFileMeta(cid string, contractIds []string) (*metadata.FileMetaInfo, error)
	GetFileMetaByCID(cid string) (*metadata.FileMetaInfo, error)
	GetContractStatus(contractId string) (metadata.Contract_ContractStatus, error)
	UpdateAutoRenewal(cid string, autoRenewal bool) error
}

type fileMeta struct {
	FileMetaAbi     *abi.FileMetaContract
	Singer          crypto.Signer
	backend         transaction.Backend
	chainId         *big.Int
	contractAddress common.Address
	lock            sync.Mutex
}

func NewFileMetaService(address common.Address, backend transaction.Backend, singer crypto.Signer, chainId *big.Int) FileMeta {
	fileMetaContract, err := abi.NewFileMetaContract(address, backend)
	if err != nil {
		return nil
	}
	return &fileMeta{
		FileMetaAbi:     fileMetaContract,
		Singer:          singer,
		backend:         backend,
		chainId:         chainId,
		contractAddress: address,
	}
}

func (fm *fileMeta) AddFileMeta(cid string, meta *metadata.FileMetaInfo) error {
	if cid == "" {
		return fmt.Errorf("cid cannot be empty")
	}
	if meta == nil {
		return fmt.Errorf("meta cannot be nil")
	}

	opts, err := bind.NewKeyedTransactorWithChainID(fm.Singer.PrivKey(), fm.chainId)
	if err != nil {
		fmt.Printf("Failed to create transactor: %v\n", err)
		return err
	}

	mb, err := proto.Marshal(meta)
	if err != nil {
		fmt.Printf("Failed to marshal metadata: %v\n", err)
		return err
	}

	pairs := make([]abi.FileMetaContractSPPair, 0)
	for _, c := range meta.Contracts {
		if c.Meta.ContractId == "" {
			fmt.Printf("Warning: empty contract ID found\n")
			continue
		}

		hostAddress, err := getPublicAddressFromHostID(c.Meta.SpId)
		if err != nil {
			fmt.Printf("Warning: failed to get host address for contract ID %s: %v\n", c.Meta.ContractId, err)
			continue
		}
		var hash [32]byte
		copy(hash[:], keccak256([]byte(c.Meta.ContractId)))
		pairs = append(pairs, abi.FileMetaContractSPPair{
			ContractId: hash,
			Sp:         common.HexToAddress(hostAddress), // Convert hostID to Ethereum address
		})
	}

	fmt.Printf("Adding file meta - CID: %s, Metadata size: %d bytes, Contracts count: %d\n",
		cid, len(mb), len(pairs))

	tx, err := fm.FileMetaAbi.AddFileMeta(opts, cid, mb, new(big.Int).SetUint64(meta.FileSize), pairs)
	if err != nil {
		fmt.Printf("Failed to add file meta: %v\n", err)
		fmt.Printf("Contracts address: %s\n", fm.contractAddress.Hex())
		fmt.Printf("Gas limit: %d\n", opts.GasLimit)
		return fmt.Errorf("smart contract execution failed: %w", err)
	}
	fmt.Printf("Successfully added file meta, transaction hash: %s\n", tx.Hash())
	return nil
}

func getPublicAddressFromHostID(hostID string) (string, error) {
	peerID, err := peer.Decode(hostID)
	if err != nil {
		return "", fmt.Errorf("failed to decode hostID: %v", err)
	}

	pubKey, err := peerID.ExtractPublicKey()
	if err != nil {
		return "", fmt.Errorf("failed to extract public key: %v", err)
	}

	pkBytes, err := cp.Secp256k1PublicKeyRaw(pubKey)
	if err != nil {
		panic(err)
	}

	ethPk, err := ethCrypto.UnmarshalPubkey(pkBytes)
	if err != nil {
		return "", err
	}

	return ethCrypto.PubkeyToAddress(*ethPk).Hex(), nil

}

func (fm *fileMeta) UpdateContractStatus(contractId string) error {
	fm.lock.Lock()
	defer fm.lock.Unlock()
	opts, err := bind.NewKeyedTransactorWithChainID(fm.Singer.PrivKey(), fm.chainId)
	if err != nil {
		fmt.Printf("Failed to create transactor: %v\n", err)
		return err
	}
	tx, err := fm.FileMetaAbi.UpdateStatus(opts, contractId, uint8(metadata.Contract_COMPLETED))
	if err != nil {
		fmt.Printf("update status error:%v, %s\n", err, contractId)
		return err
	}

	for i := 0; i < 10; i++ {
		receipt, err := fm.backend.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				time.Sleep(10 * time.Second)
				continue
			}
			fmt.Printf("update status error:%v, %s\n", err, contractId)
			return err
		}
		if receipt.Status != 1 {
			fmt.Printf("transaction failed on-chain (likely out of gas), txHash: %s\n", tx.Hash())
			return err
		}
	}
	fmt.Println("update contract:", contractId, "status ok:", tx.Hash())
	return nil
}

func (fm *fileMeta) GetFileMeta(cid string, contractIds []string) (*metadata.FileMetaInfo, error) {
	meta, err := fm.FileMetaAbi.GetFileMeta(nil, cid, contractIds)

	if err != nil {
		return nil, err
	}
	fss := &metadata.FileMetaInfo{}
	err = proto.Unmarshal(meta.MetaData, fss)
	if err != nil {
		return nil, err
	}
	for _, c := range fss.Contracts {
		for i, contractId := range contractIds {
			if c.Meta.ContractId == contractId {
				status := meta.Statuses[i]
				c.Status = metadata.Contract_ContractStatus(int32(status))
			}
		}
	}
	return fss, nil
}

func (fm *fileMeta) GetFileMetaByCID(cid string) (*metadata.FileMetaInfo, error) {
	var hash [32]byte
	copy(hash[:], keccak256([]byte(cid)))
	meta, err := fm.FileMetaAbi.FileMeta(nil, hash)
	if err != nil {
		return nil, err
	}
	fss := &metadata.FileMetaInfo{}
	err = proto.Unmarshal(meta, fss)
	if err != nil {
		return nil, err
	}
	return fss, nil
}

func (fm *fileMeta) GetContractStatus(contractId string) (metadata.Contract_ContractStatus, error) {
	var hash [32]byte
	copy(hash[:], keccak256([]byte(contractId)))
	meta, err := fm.FileMetaAbi.ContractStatus(nil, hash)
	if err != nil {
		return metadata.Contract_ContractStatus(0), err
	}
	return metadata.Contract_ContractStatus(int32(meta)), nil
}

func (fm *fileMeta) UpdateAutoRenewal(cid string, autoRenewal bool) error {
	// opts, err := bind.NewKeyedTransactorWithChainID(fm.Singer.PrivKey(), fm.chainId)
	// if err != nil {
	// 	fmt.Printf("Failed to create transactor: %v\n", err)
	// 	return err
	// }
	// TODO
	// tx, err := fm.FileMetaAbi.UpdateAutoRenewal(opts, cid, autoRenewal)
	return nil
}

func keccak256(data ...[]byte) []byte {
	h := sha3.NewLegacyKeccak256()
	for _, b := range data {
		h.Write(b)
	}
	return h.Sum(nil)
}
