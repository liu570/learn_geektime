package cache

import (
	"time"
)

type RetryStrategy interface {
	// Next 返回下一次重试的间隔 如果下一次不需要的第二参数返回 false
	Next() (time.Duration, bool)
}

type FixIntervalRetry struct {
	//	最大间隔
	Interval time.Duration
	Max      int
	Cnt      int
}

func (f *FixIntervalRetry) Next() (time.Duration, bool) {
	f.Cnt++
	return f.Interval, f.Cnt <= f.Max
}
