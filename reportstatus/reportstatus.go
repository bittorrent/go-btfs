package reportstatus

import (
	"context"
	"encoding/hex"
	"fmt"
	onlinePb "github.com/tron-us/go-btfs-common/protos/online"
	"math/big"
	"strings"
	"time"

	config "github.com/TRON-US/go-btfs-config"
	"github.com/bittorrent/go-btfs/chain"
	"github.com/bittorrent/go-btfs/reportstatus/abi"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/ethereum/go-ethereum/common"

	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("report-status-contract:")

var (
	statusABI = transaction.ParseABIUnchecked(abi.StatusHeartABI)
	serv      *service
)

const (
	ReportStatusTime = 10 * time.Minute
	//ReportStatusTime = 10 * time.Second // 10 * time.Minute
)

func Init(transactionService transaction.Service, cfg *config.Config, configRoot string, statusAddress common.Address, chainId int64) error {
	New(statusAddress, transactionService, cfg)

	err := CheckExistLastOnline(cfg, configRoot, chainId)
	if err != nil {
		return err
	}
	return nil
}

func isReportStatusEnabled(cfg *config.Config) bool {
	return cfg.Experimental.StorageHostEnabled || cfg.Experimental.ReportStatusContract
}

type service struct {
	statusAddress      common.Address
	transactionService transaction.Service
}

type Service interface {
	// ReportStatus report status heart info to statusContract
	ReportStatus() (common.Hash, error)

	// CheckReportStatus check report status heart info to statusContract
	CheckReportStatus() error
}

func New(statusAddress common.Address, transactionService transaction.Service, cfg *config.Config) Service {
	serv = &service{
		statusAddress:      statusAddress,
		transactionService: transactionService,
	}

	if isReportStatusEnabled(cfg) {
		go func() {
			cycleCheckReport()
		}()
	}
	return serv
}

// ReportStatus report heart status
func (s *service) ReportStatus() (common.Hash, error) {
	lastOnline, err := chain.GetLastOnline()
	if err != nil {
		return common.Hash{}, err
	}

	if lastOnline == nil {
		return common.Hash{}, nil
	}
	if len(lastOnline.LastSignedInfo.Peer) <= 0 {
		return common.Hash{}, nil
	}

	peer := lastOnline.LastSignedInfo.Peer
	createTime := lastOnline.LastSignedInfo.CreatedTime
	version := lastOnline.LastSignedInfo.Version
	num := lastOnline.LastSignedInfo.Nonce
	bttcAddress := common.HexToAddress(lastOnline.LastSignedInfo.BttcAddress)
	signedTime := lastOnline.LastSignedInfo.SignedTime
	signature, err := hex.DecodeString(strings.Replace(lastOnline.LastSignature, "0x", "", -1))
	//fmt.Println("... ReportStatus, param = ", peer, createTime, version, num, bttcAddress, signedTime, signature)
	fmt.Printf("... ReportStatus, lastOnline = %+v \n", lastOnline)

	callData, err := statusABI.Pack("reportStatus", peer, createTime, version, num, bttcAddress, signedTime, signature)
	if err != nil {
		return common.Hash{}, err
	}

	request := &transaction.TxRequest{
		To:          &s.statusAddress,
		Data:        callData,
		Value:       big.NewInt(0),
		Description: "Report Heart Status",
	}

	txHash, err := s.transactionService.Send(context.Background(), request)
	if err != nil {
		return common.Hash{}, err
	}
	fmt.Println("... ReportStatus, txHash, err = ", txHash, err)
	_, err = chain.SetReportStatusOK()
	if err != nil {
		return common.Hash{}, err
	}

	// WaitForReceipt takes long time
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("ReportHeartStatus recovered:%+v", err)
			}
		}()
	}()
	return txHash, nil
}

// report heart status
func (s *service) checkLastOnlineInfo(peerId, bttcAddr string) error {
	callData, err := statusABI.Pack("getStatus", peerId)
	if err != nil {
		return err
	}
	request := &transaction.TxRequest{
		To:   &s.statusAddress,
		Data: callData,
	}

	result, err := s.transactionService.Call(context.Background(), request)
	if err != nil {
		return err
	}
	v, err := statusABI.Unpack("getStatus", result)
	if err != nil {
		return err
	}
	//fmt.Printf("...... getStatus - result v = %+v, err = %v \n", v, err)

	nonce := v[3].(uint32)
	if nonce > 0 {
		lastOnlineInfo := chain.LastOnlineInfo{
			LastSignedInfo: onlinePb.SignedInfo{
				Peer:        v[0].(string),
				CreatedTime: v[1].(uint32),
				Version:     v[2].(string),
				Nonce:       v[3].(uint32),
				BttcAddress: bttcAddr,
				SignedTime:  v[4].(uint32),
			},
			LastSignature: "0x" + hex.EncodeToString(v[5].([]byte)),
			LastTime:      time.Now(),
		}
		fmt.Printf("... init reset lastOnlineInfo = %+v \n", lastOnlineInfo)

		err = chain.StoreOnline(&lastOnlineInfo)
		if err != nil {
			return err
		}
	}

	// WaitForReceipt takes long time
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("getStatus recovered:%+v", err)
			}
		}()
	}()
	return nil
}

// report heart status
func (s *service) genHashExt(ctx context.Context) (common.Hash, error) {
	peer := "1"
	createTime := uint32(1)
	version := "1"
	num := uint32(3)
	bttcAddress := "0x22df207EC3C8D18fEDeed87752C5a68E5b4f6FbD"
	fmt.Println("...... genHashExt, param = ", peer, createTime, version, num, bttcAddress)

	callData, err := statusABI.Pack("genHashExt", peer, createTime, version, num, common.HexToAddress(bttcAddress))
	if err != nil {
		return common.Hash{}, err
	}

	request := &transaction.TxRequest{
		To:   &s.statusAddress,
		Data: callData,
	}

	result, err := s.transactionService.Call(ctx, request)
	fmt.Println("...... genHashExt - totalStatus, result, err = ", common.Bytes2Hex(result), err)

	if err != nil {
		return common.Hash{}, err
	}

	// WaitForReceipt takes long time
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("genHashExt recovered:%+v", err)
			}
		}()
	}()
	return common.Hash{}, nil
}

func (s *service) CheckReportStatus() error {
	_, err := s.ReportStatus()
	if err != nil {
		log.Errorf("ReportStatus err:%+v", err)
		return err
	}
	return nil
}

func cycleCheckReport() {
	tick := time.NewTicker(ReportStatusTime)
	defer tick.Stop()

	// Force tick on immediate start
	// CheckReport in the for loop
	for ; true; <-tick.C {
		fmt.Println("")
		fmt.Println("... CheckReportStatus ...")

		report, err := chain.GetReportStatus()
		if err != nil {
			continue
		}
		fmt.Printf("... CheckReportStatus report: %+v \n", report)

		now := time.Now()
		nowUnixMod := now.Unix() % 86400
		// report only 1 hour every, and must after 10 hour.
		if nowUnixMod > report.ReportStatusSeconds &&
			nowUnixMod < report.ReportStatusSeconds+3600 &&
			now.Sub(report.LastReportTime) > 10*time.Hour {

			err = serv.CheckReportStatus()
			if err != nil {
				continue
			}
		}
	}
}
