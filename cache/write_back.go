package cache

import (
	"context"
	"log"
	"time"
)

type writeBackCache struct {
	*LocalCache
}

func newWriteBackCache(store func(ctx context.Context, key string, val any) error) *writeBackCache {
	return &writeBackCache{
		LocalCache: NewLocalCache(func(key string, val any) {
			// 这个地方 context 不好设置 不可用下面的 background
			// err 不好处理
			err := store(context.Background(), key, val)
			if err != nil {
				log.Fatalln(err)
			}
		}),
	}
}

func (w *writeBackCache) Close() error {
	// 遍历所有的 key，将值刷新到数据库
	// 回写模式缓存，需要在关闭缓存的时候保持数据一致性将数据写入数据库

	return nil
}

// 与回写对应的是缓存的预加载
type PreloadCache struct {
	Cache
	sentinelCache *LocalCache
}

func NewPreloadCache(c Cache, loadFunc func(ctx context.Context, key string) (any, error)) *PreloadCache {

	// sentinelCache 上的 key value 过期
	// 就把主 Cache 上的数据刷新
	return &PreloadCache{
		Cache: c,
		sentinelCache: NewLocalCache(func(key string, val any) {
			// context.Background() 可以做成参数
			val, err := loadFunc(context.Background(), key)
			if err == nil {
				// context.Background()、time.Minute 可以做成参数
				err = c.Set(context.Background(), key, val, time.Minute)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}),
	}
}

func (c *PreloadCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	// 哨兵缓存在这里只存一个键值不存具体数据
	err := c.sentinelCache.Set(ctx, key, "", expiration-time.Second*3)
	if err != nil {
		log.Fatalln(err)
	}
	return c.Cache.Set(ctx, key, val, expiration)
	//panic("implement me!")
}

/*
var group = &singleflight.Group{}

func Biz(key string) {
	val, err := cache.Get(context.Background(), key)
	if err == KeyNotFound {
		val, err, _ = group.Do(key, func() (interface{}, error) {
			newVal, err := QueryFromDB(key)
			if err != nil {
				return nil, err
			}
			err = cache.Set(context.Background(), key, newVal)
			return newVal, er
		})
	}
	println(val)
}
*/
