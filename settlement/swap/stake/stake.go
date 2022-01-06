package stake

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	conabi "github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/settlement/swap/erc20"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("staking")

var (
	errDecodeABI = errors.New("could not decode stake abi data")
)

type stakeInfo struct {
	stakeAmount *big.Int
	lockTime    *big.Int
}

type service struct {
	stakeAddress       common.Address
	transactionService transaction.Service
	stake              *stakeInfo
	erc20Service       erc20.Service
	timeDivisor        int64
	issuer             common.Address
	quitC              chan struct{}
}

type Service interface {
	io.Closer
	CurrentStakeAmount() *big.Int
	CurrentStakeLockTime() *big.Int
	// GetStakeInfo retrieves latest available information from stake contract
	GetStakeInfo(ctx context.Context) (*stakeInfo, error)
	Start()
	//approve
	Approve(ctx context.Context, amount *big.Int) (common.Hash, error)
	//add stake
	AddStake(ctx context.Context, amount *big.Int) (common.Hash, error)
	//remove stake
	RmStake(ctx context.Context, amount *big.Int) (common.Hash, error)
}

var (
	stakeABI = transaction.ParseABIUnchecked(conabi.StakeAbi)
)

func New(
	stakeAddress common.Address,
	transactionService transaction.Service,
	erc20Service erc20.Service,
	overlayEthAddress common.Address,
	timeDivisor int64,
) Service {
	return &service{
		stakeAddress:       stakeAddress,
		transactionService: transactionService,
		stake:              nil,
		erc20Service:       erc20Service,
		quitC:              make(chan struct{}),
		timeDivisor:        timeDivisor,
		issuer:             overlayEthAddress,
	}
}

func (s *service) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		<-s.quitC
	}()

	go func() {
		defer cancel()
		for {
			stakeinfo, err := s.GetStakeInfo(ctx)
			if err != nil {
				log.Errorf("could not get stake info: %v", err)
			} else {
				log.Infof("updated stake info to %v", stakeinfo)
				s.stake = stakeinfo

				fmt.Println("stake info: amount = ")
			}

			ts := time.Now().Unix()

			// We poll the oracle in every timestamp divisible by constant 300 (timeDivisor)
			// in order to get latest version approximately at the same time on all nodes
			// and to minimize polling frequency
			// If the node gets newer information than what was applicable at last polling point at startup
			// this minimizes the negative scenario to less than 5 minutes
			// during which cheques can not be sent / accepted because of the asymmetric information
			timeUntilNextPoll := time.Duration(s.timeDivisor-ts%s.timeDivisor) * time.Second

			select {
			case <-s.quitC:
				return
			case <-time.After(timeUntilNextPoll):
			}
		}
	}()
}

func (s *service) GetStakeInfo(ctx context.Context) (*stakeInfo, error) {
	stake, err := s.getStakeAmount(ctx)
	if err != nil {
		return nil, err
	}

	lockTime, err := s.getStakeLockTime(ctx)
	if err != nil {
		return nil, err
	}

	//tm := time.Unix(lockTime.Int64(), 0)

	return &stakeInfo{stake, lockTime}, nil
}

func (s *service) CurrentStakeAmount() *big.Int {
	return s.stake.stakeAmount
}

func (s *service) CurrentStakeLockTime() *big.Int {
	return s.stake.lockTime
}

func (s *service) Approve(ctx context.Context, amount *big.Int) (common.Hash, error) {
	hash, err := s.erc20Service.Approve(ctx, s.stakeAddress, amount)
	if err != nil {
		return common.Hash{}, err
	}
	return hash, nil
}

//add stake
func (s *service) AddStake(ctx context.Context, amount *big.Int) (common.Hash, error) {
	approve_amount, err := s.erc20Service.Allowance(ctx, s.issuer, s.stakeAddress)
	if err != nil {
		return common.Hash{}, err
	}

	if approve_amount == big.NewInt(0) {
		return common.Hash{}, errors.New("pls approve token to stake contract")
	}
	callData, err := stakeABI.Pack("stake", amount)
	if err != nil {
		return common.Hash{}, err
	}

	request := &transaction.TxRequest{
		To:          &s.stakeAddress,
		Data:        callData,
		Value:       big.NewInt(0),
		Description: "stake",
	}

	txHash, err := s.transactionService.Send(ctx, request)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

//remove stake
func (s *service) RmStake(ctx context.Context, amount *big.Int) (common.Hash, error) {
	callData, err := stakeABI.Pack("unStake", amount)
	if err != nil {
		return common.Hash{}, err
	}

	request := &transaction.TxRequest{
		To:          &s.stakeAddress,
		Data:        callData,
		Value:       big.NewInt(0),
		Description: "unstake",
	}

	txHash, err := s.transactionService.Send(ctx, request)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (s *service) getStakeAmount(ctx context.Context) (*big.Int, error) {
	callData, err := stakeABI.Pack("selfStakeInfo")
	if err != nil {
		return nil, err
	}
	result, err := s.transactionService.Call(ctx, &transaction.TxRequest{
		To:   &s.stakeAddress,
		Data: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := stakeABI.Unpack("selfStakeInfo", result)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, errDecodeABI
	}

	stake, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || stake == nil {
		return nil, errDecodeABI
	}

	return stake, nil
}

func (s *service) getStakeLockTime(ctx context.Context) (*big.Int, error) {
	callData, err := stakeABI.Pack("selfTimeCanBeDecreased")
	if err != nil {
		return nil, err
	}
	result, err := s.transactionService.Call(ctx, &transaction.TxRequest{
		To:   &s.stakeAddress,
		Data: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := stakeABI.Unpack("selfTimeCanBeDecreased", result)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, errDecodeABI
	}

	lockTime, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || lockTime == nil {
		return nil, errDecodeABI
	}

	return lockTime, nil
}

func (s *service) Close() error {
	close(s.quitC)
	return nil
}
