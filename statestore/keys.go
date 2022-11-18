package statestore

import (
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
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

func PeerReceivedUncashRecordsCountKey(vault common.Address, token common.Address) string {
	return fmt.Sprintf("%s%s", tokencfg.AddToken(PeerReceivedUncashRecordsCountKeyPrefix, token), vault.String())
}

func GetTodayTotalDailyReceivedKey(token common.Address) string {
	return fmt.Sprintf("%s%d", tokencfg.AddToken(TotalDailyReceivedKey, token), utils.TodayUnix())
}

func GetTotalDailyReceivedKeyByTime(timestamp int64, token common.Address) string {
	return fmt.Sprintf("%s%d", tokencfg.AddToken(TotalDailyReceivedKey, token), timestamp)
}

func GetTodayTotalDailyReceivedCashedKey(token common.Address) string {
	return fmt.Sprintf("%s%d", tokencfg.AddToken(TotalDailyReceivedCashedKey, token), utils.TodayUnix())
}

func GetTotalDailyReceivedCashedKeyByTime(timestamp int64, token common.Address) string {
	return fmt.Sprintf("%s%d", tokencfg.AddToken(TotalDailyReceivedCashedKey, token), timestamp)
}

func GetTodayTotalDailySentKey(token common.Address) string {
	return GetTotalDailySentKeyByTime(utils.TodayUnix(), token)
}

func GetTotalDailySentKeyByTime(timestamp int64, token common.Address) string {
	return fmt.Sprintf("%s%d", tokencfg.AddToken(TotalDailySentKey, token), timestamp)
}

func CashoutResultPrefixKey(token common.Address) string {
	return tokencfg.AddToken("swap_cashout_result_", token)
}

func CashoutResultKey(vault common.Address, token common.Address) string {
	return fmt.Sprintf("%s%x_%d", CashoutResultPrefixKey(token), vault, time.Now().Unix())
}
