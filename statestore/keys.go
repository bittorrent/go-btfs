package statestore

import (
	"fmt"
	"time"

	"github.com/bittorrent/go-btfs/utils"
	"github.com/ethereum/go-ethereum/common"
)

var (
	TotalReceivedKey      = "swap_vault_total_received"       // 收到支票总额度
	TotalReceivedCountKey = "swap_vault_total_received_count" // 收到支票总数量

	TotalReceivedCashedKey      = "swap_vault_total_received_cashed"       // 收到支票兑现总额度
	TotalReceivedCashedCountKey = "swap_vault_total_received_cashed_count" // 收到支票兑现总数量

	TotalDailyReceivedKey       = "swap_vault_total_daily_received_"        // 单日收到支票总额度+总数量
	TotalDailyReceivedCashedKey = "swap_vault_total_daily_received_cashed_" // 单日收到支票兑现总额度
	TotalDailySentKey           = "swap_vault_total_daily_sent_"            // 单日发出支票总额度/总数量

	PeerReceivedUncashRecordsCountKeyPrefix = "swap_vault_peer_received_uncashed_records_count_" // 每个peer收到支票未兑现数量
)

func PeerReceivedUncashRecordsCountKey(vault common.Address) string {
	return fmt.Sprintf("%s%s", PeerReceivedUncashRecordsCountKeyPrefix, vault.String())
}

func GetTodayTotalDailyReceivedKey() string {
	return fmt.Sprintf("%s%d", TotalDailyReceivedKey, utils.TodayUnix())
}

func GetTotalDailyReceivedKeyByTime(timestamp int64) string {
	return fmt.Sprintf("%s%d", TotalDailyReceivedKey, timestamp)
}

func GetTodayTotalDailyReceivedCashedKey() string {
	return fmt.Sprintf("%s%d", TotalDailyReceivedCashedKey, utils.TodayUnix())
}

func GetTotalDailyReceivedCashedKeyByTime(timestamp int64) string {
	return fmt.Sprintf("%s%d", TotalDailyReceivedCashedKey, timestamp)
}

func GetTodayTotalDailySentKey() string {
	return GetTotalDailySentKeyByTime(utils.TodayUnix())
}

func GetTotalDailySentKeyByTime(timestamp int64) string {
	return fmt.Sprintf("%s%d", TotalDailySentKey, timestamp)
}

func CashoutResultPrefixKey() string {
	return "swap_cashout_result_"
}

func CashoutResultKey(vault common.Address) string {
	return fmt.Sprintf("%s%x_%d", CashoutResultPrefixKey(), vault, time.Now().Unix())
}
