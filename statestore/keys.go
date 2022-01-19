package statestore

import (
	"fmt"
	"time"

	"github.com/bittorrent/go-btfs/utils"
	"github.com/ethereum/go-ethereum/common"
)

var (
	TotalReceivedKey      = "swap_vault_total_received"
	TotalReceivedCountKey = "swap_vault_total_received_count"

	TotalReceivedCashedKey      = "swap_vault_total_received_uncashed"
	TotalReceivedCashedCountKey = "swap_vault_total_received_cashed_count"
	TotalDailyReceivedKey       = "swap_vault_total_daily_received_"
	TotalDailySentKey           = "swap_vault_total_daily_sent_"
	TotalDailyReceivedCashedKey = "swap_vault_total_daily_received_cashed_"

	PeerReceivedUncashRecordsCountKeyPrefix = "swap_vault_peer_received_uncashed_records_count_"
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
