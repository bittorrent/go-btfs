package reportstatus

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	onlinePb "github.com/tron-us/go-btfs-common/protos/online"

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
	//ReportStatusTime = 60 * time.Second // 10 * time.Minute
)

func Init(transactionService transaction.Service, cfg *config.Config, statusAddress common.Address) Service {
	return New(statusAddress, transactionService, cfg)
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

	// CheckLastOnlineInfo check last online info.
	CheckLastOnlineInfo(peerId, bttcAddr string) error
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
	nonce := lastOnline.LastSignedInfo.Nonce
	bttcAddress := common.HexToAddress(lastOnline.LastSignedInfo.BttcAddress)
	signedTime := lastOnline.LastSignedInfo.SignedTime
	signature, err := hex.DecodeString(strings.Replace(lastOnline.LastSignature, "0x", "", -1))
	//fmt.Printf("... ReportStatus, lastOnline = %+v \n", lastOnline)

	callData, err := statusABI.Pack("reportStatus", peer, createTime, version, nonce, bttcAddress, signedTime, signature)
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
	//fmt.Println("... ReportStatus, txHash, err = ", txHash, err)

	now := time.Now()
	_, err = chain.SetReportStatusOK()
	if err != nil {
		return common.Hash{}, err
	}

	stx, err := s.transactionService.WaitForReceipt(context.Background(), txHash)
	if err != nil {
		return common.Hash{}, err
	}
	gasPrice := getGasPrice(request)
	gasTotal := big.NewInt(1).Mul(gasPrice, big.NewInt(int64(stx.GasUsed)))
	//fmt.Println("... ReportStatus, gasPrice, stx.GasUsed, gasTotal = ", gasPrice.String(), stx.GasUsed, gasTotal.String())

	r := &chain.LevelDbReportStatusInfo{
		Peer:           peer,
		BttcAddress:    bttcAddress.String(),
		StatusContract: s.statusAddress.String(),
		Nonce:          nonce,
		TxHash:         txHash.String(),
		GasSpend:       gasTotal.String(),
		ReportTime:     now,
	}
	_, err = chain.SetReportStatusListOK(r)
	if err != nil {
		return common.Hash{}, err
	}
	return txHash, nil
}

func getGasPrice(request *transaction.TxRequest) *big.Int {
	var gasPrice *big.Int
	if request.GasPrice == nil {
		gasPrice = big.NewInt(300000000000000)
	} else {
		gasPrice = request.GasPrice
	}
	return gasPrice
}

func CmdReportStatus() error {
	_, err := serv.ReportStatus()
	if err != nil {
		log.Errorf("ReportStatus err:%+v", err)
		return err
	}
	return nil
}

// CheckLastOnlineInfo report heart status
func (s *service) CheckLastOnlineInfo(peerId, bttcAddr string) error {
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
		//fmt.Println("")
		//fmt.Println("... ReportStatus, CheckReportStatus ...")

		report, err := chain.GetReportStatus()
		//fmt.Printf("... ReportStatus, CheckReportStatus report: %+v err:%+v \n", report, err)
		if err != nil {
			log.Errorf("GetReportStatus err:%+v", err)
			continue
		}

		now := time.Now()
		nowUnixMod := now.Unix() % 86400
		// report only 1 hour every, and must after 10 hour.
		if nowUnixMod > report.ReportStatusSeconds &&
			nowUnixMod < report.ReportStatusSeconds+3600 &&
			now.Sub(report.LastReportTime) > 10*time.Hour {

			err = serv.CheckReportStatus()
			if err != nil {
				log.Errorf("CheckReportStatus err:%+v", err)
				continue
			}
		}
	}
}
