package cache

import (
	"context"
	"learn_geektime/cache/internal/errs"
	"log"
	"time"
)

// 该类用来实现 read through 缓存模式
// read through 即是 读 DB 交由 cache 实现

type ReadThroughCache struct {
	Cache
	Expiration time.Duration

	// 我们把最常见的“捞DB”这种说法抽象为“加载数据”
	LoadFunc func(ctx context.Context, key string) (any, error)
}

func (c *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	// 逻辑: 先捞缓存，再捞 DB
	val, err := c.Cache.Get(ctx, key)

	// 未知错误
	if err != nil && err != errs.NewErrKeyNotFound(key) {
		return nil, errs.NewErrKeyNotFound(key)
	}

	// 缓存没有数据
	if err == errs.NewErrKeyNotFound(key) {
		// 捞 db
		val, err = c.LoadFunc(ctx, key)
		if err != nil {
			return nil, err
		}
		// 这里 err 可以考虑忽略掉，或者输出 warn 日志
		err = c.Set(ctx, key, val, c.Expiration)
		if err != nil {
			// 这里要考虑刷新缓存失败，究竟要不要返回 val
			//return nil, err

			// 这里只输出错误日志
			log.Fatalln(err)
		}
		return val, nil
	}
	return val, nil
}
