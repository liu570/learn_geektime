package cache

import (
	"context"
	"fmt"
	"learn_geektime/cache/internal/errs"
	"log"
	"sync"
	"time"
)

//	ReadThroughCache 该类用来实现 read through 缓存模式
//
// read through 即是 读 DB 交由 cache 实现
type ReadThroughCache struct {
	mutex sync.RWMutex
	Cache
	Expiration time.Duration

	// 我们把最常见的“捞DB”这种说法抽象为“加载数据”
	LoadFunc func(ctx context.Context, key string) (any, error)
}

func (c *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	// 逻辑: 先捞缓存，再捞 DB

	c.mutex.RLock()
	val, err := c.Cache.Get(ctx, key)
	c.mutex.RUnlock()

	// 未知错误
	if err != nil && err != errs.NewErrKeyNotFound(key) {
		return nil, err
	}

	// 缓存没有数据
	if err != nil && err == errs.NewErrKeyNotFound(key) {
		// 加锁问题
		c.mutex.Lock()
		defer c.mutex.Unlock()
		// 捞 db
		val, err = c.LoadFunc(ctx, key)
		if err != nil {
			// 这里会暴露 LoadFunc 底层
			// 例如如果 LoadFunc 是数据库查询，这里就回暴露数据库的错误信息（或者 orm 框架的）
			return nil, fmt.Errorf("cache: 无法加载数据, %w", err)
			//return nil, err
		}

		// 同步刷新缓存
		// 这里 err 可以考虑忽略掉，或者输出 warn 日志
		err = c.Set(ctx, key, val, c.Expiration)
		if err != nil {
			// 这里要考虑刷新缓存失败，究竟要不要返回 value
			//return nil, err

			// 这里只输出错误日志
			log.Fatalln(err)
		}
		return val, nil
	}

	// 缓存中有数据 直接返回
	return val, nil
}

func (c *ReadThroughCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	// TODO: 注意这里加锁需要斟酌 加上锁还是有数据不一致的问题
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Cache.Set(ctx, key, val, expiration)
	//panic("implement me!")
}

type ReadThroughCacheV1[T any] struct {
	Cache
	Expiration time.Duration

	// 我们把最常见的“捞DB”这种说法抽象为“加载数据”
	LoadFunc func(ctx context.Context, key string) (T, error)
}

func (c *ReadThroughCacheV1[T]) Get(ctx context.Context, key string) (any, error) {
	// 逻辑: 先捞缓存，再捞 DB
	val, err := c.Cache.Get(ctx, key)

	// 未知错误
	if err != nil && err != errs.NewErrKeyNotFound(key) {
		return nil, err
	}

	// 缓存没有数据
	if err != nil && err == errs.NewErrKeyNotFound(key) {
		// 捞 db
		val, err = c.LoadFunc(ctx, key)
		if err != nil {
			return nil, err
		}

		// 这里 err 可以考虑忽略掉，或者输出 warn 日志
		err = c.Set(ctx, key, val, c.Expiration)
		if err != nil {
			// 这里要考虑刷新缓存失败，究竟要不要返回 value
			//return nil, err

			// 这里只输出错误日志
			log.Fatalln(err)
		}
		return val, nil
	}

	// 缓存中有数据 直接返回
	return val, nil
}

type ReadThroughCacheV2[T any] struct {
	CacheV2[T]
	Expiration time.Duration

	// 我们把最常见的“捞DB”这种说法抽象为“加载数据”
	LoadFunc func(ctx context.Context, key string) (T, error)
}

type ReadThroughCacheV3 struct {
	mutex sync.RWMutex
	Cache
	Expiration time.Duration

	// 我们把最常见的“捞DB”这种说法抽象为“加载数据”
	//LoadFunc func(ctx context.Context, key string) (any, error)
	Loader
}

type Loader interface {
	Load(ctx context.Context, key string) (any, error)
}

type LoadFunc func(ctx context.Context, key string) (any, error)

func (l LoadFunc) Load(ctx context.Context, key string) (any, error) {
	return l(ctx, key)
}
