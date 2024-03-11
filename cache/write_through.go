package cache

import (
	"context"
	"time"
)

type WriteThrough struct {
	Cache
	StoreFunc func(ctx context.Context, key string, val any) error
}

func (w *WriteThrough) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	// 在这里开 goroutine就是全异步
	err := w.StoreFunc(ctx, key, val)
	if err != nil {
		return err
	}
	//在这里开 goroutine 就是半异步
	return w.Cache.Set(ctx, key, val, expiration)
}
