package transaction

import "time"

var (
	StoredTransactionKey = storedTransactionKey
)

func (s *Matcher) SetTimeNow(f func() time.Time) {
	s.timeNow = f
}

func (s *Matcher) SetTime(k int64) {
	s.SetTimeNow(func() time.Time {
		return time.Unix(k, 0)
	})
}
