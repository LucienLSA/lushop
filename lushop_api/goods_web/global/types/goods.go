package types

import "time"

func Uint64ToTime(value uint64) time.Time {
	t := time.Unix(int64(value), 0)
	return t
}
