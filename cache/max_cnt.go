package cache

import (
	"context"
	"errors"
	"learn_geektime/cache/internal/errs"
	"sync"
	"sync/atomic"
	"time"
)

// MaxCntCacheDecorator 对本地缓存的封装 使用装饰器模式增加功能
type MaxCntCacheDecorator struct {
	mutex  sync.Mutex
	MaxCnt int32
	Cnt    int32
	Cache  *LocalCache
}

func NewMaxCntCache(maxCnt int32) *MaxCntCacheDecorator {
	res := &MaxCntCacheDecorator{MaxCnt: maxCnt}
	c := NewLocalCache(func(key string, val any) {
		atomic.AddInt32(&res.Cnt, -1)
	})
	res.Cache = c
	return res
}

func (c *MaxCntCacheDecorator) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, err := c.Cache.Get(ctx, key)
	if err != nil && err != errs.NewErrKeyNotFound(key) {
		// 这个错误如何处理
		return err
	}
	if err == errs.NewErrKeyNotFound(key) {
		// 判断有没有超过最大值
		cnt := atomic.AddInt32(&c.Cnt, 1)
		if cnt > c.MaxCnt {
			atomic.AddInt32(&c.Cnt, -1)
			return errors.New("cache 已经满了")
		}
	}
	return c.Cache.Set(ctx, key, val, expiration)
}
