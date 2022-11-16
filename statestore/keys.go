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

func addToken(s string, token string) string {
	if token == "WBTT" {
		return s
	}
	return fmt.Sprintf("%s_%s", s, token)
}

func PeerReceivedUncashRecordsCountKey(vault common.Address, token string) string {
	return fmt.Sprintf("%s%s", addToken(PeerReceivedUncashRecordsCountKeyPrefix, token), vault.String())
}

func GetTodayTotalDailyReceivedKey(token string) string {
	return fmt.Sprintf("%s%d", addToken(TotalDailyReceivedKey, token), utils.TodayUnix())
}

func GetTotalDailyReceivedKeyByTime(timestamp int64, token string) string {
	return fmt.Sprintf("%s%d", addToken(TotalDailyReceivedKey, token), timestamp)
}

func GetTodayTotalDailyReceivedCashedKey(token string) string {
	return fmt.Sprintf("%s%d", addToken(TotalDailyReceivedCashedKey, token), utils.TodayUnix())
}

func GetTotalDailyReceivedCashedKeyByTime(timestamp int64, token string) string {
	return fmt.Sprintf("%s%d", addToken(TotalDailyReceivedCashedKey, token), timestamp)
}

func GetTodayTotalDailySentKey(token string) string {
	return GetTotalDailySentKeyByTime(utils.TodayUnix(), token)
}

func GetTotalDailySentKeyByTime(timestamp int64, token string) string {
	return fmt.Sprintf("%s%d", addToken(TotalDailySentKey, token), timestamp)
}

func CashoutResultPrefixKey(token string) string {
	return addToken("swap_cashout_result_", token)
}

func CashoutResultKey(vault common.Address, token string) string {
	return fmt.Sprintf("%s%x_%d", CashoutResultPrefixKey(token), vault, time.Now().Unix())
}
