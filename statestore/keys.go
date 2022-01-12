package statestore

import (
	"fmt"
	"time"
)

var (
	TotalReceivedKey      = "swap_vault_total_received"
	TotalReceivedCountKey = "swap_vault_total_received_count"

	TotalReceivedCashedKey      = "swap_vault_total_received_uncashed"
	TotalReceivedCashedCountKey = "swap_vault_total_received_cashed_count"
	TotalDailyReceivedKey       = "swap_vault_total_daily_received_"
	TotalDailyReceivedCashedKey = "swap_vault_total_daily_received_cashed_"
)

func GetTodayTotalDailyReceivedKey() string {
	y, m, d := time.Now().Date()
	todayStart := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return fmt.Sprintf("%s%d", TotalDailyReceivedKey, todayStart.Unix())
}

func GetTotalDailyReceivedKeyByTime(timestamp int64) string {
	return fmt.Sprintf("%s%d", TotalDailyReceivedKey, timestamp)
}

func GetTodayTotalDailyReceivedCashedKey() string {
	y, m, d := time.Now().Date()
	todayStart := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return fmt.Sprintf("%s%d", TotalDailyReceivedCashedKey, todayStart.Unix())
}

func GetTotalDailyReceivedCashedKeyByTime(timestamp int64) string {
	return fmt.Sprintf("%s%d", TotalDailyReceivedCashedKey, timestamp)
}
