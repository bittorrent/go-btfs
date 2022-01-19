package utils

import "time"

// TodayUnix truncate today to date,discards hour,minute and second
// and return unix timestamp
func TodayUnix() int64 {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix()
}
