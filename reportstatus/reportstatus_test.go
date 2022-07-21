package reportstatus

import (
	"encoding/hex"
	"fmt"
	"github.com/bittorrent/go-btfs/reportstatus/abi"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func TestReportStatusPackInput(t *testing.T) {
	peer := "1"
	createTime := uint32(1)
	version := "1"
	num := uint32(3)
	bttcAddress := common.HexToAddress("0x22df207EC3C8D18fEDeed87752C5a68E5b4f6FbD")
	signature, err := hex.DecodeString("3aab4d1631635d68bb8b9035c956b7e776dc972aa36e98177643a9dd47df7d3946459102b5678d73ad905b958cbf57ce6a001b3d27ecd204e6125a2543f897dc01") // can't contain 0x ...
	fmt.Println("...... ReportStatus, param = ", peer, createTime, version, num, bttcAddress, signature, err)

	statusABI := transaction.ParseABIUnchecked(abi.StatusHeartABI)
	callData, err := statusABI.Pack("reportStatus", peer, createTime, version, num, bttcAddress, signature)
	if err != nil {
		return
	}
	fmt.Println("...... ReportStatus, callData, err = ", common.Bytes2Hex(callData), err)
}
