package cache

import (
	"context"
	"errors"
	"learn_geektime/cache/internal/errs"
	"time"
)

// 布隆过滤器，在缓存的未命中的时候询问一下这些过滤器

type BloomCache struct {
	BloomFilter
	Cache
	LoadFunc func(ctx context.Context, key string) (any, error)
}

func (b *BloomCache) Get(ctx context.Context, key string) (any, error) {
	val, err := b.Cache.Get(ctx, key)
	if err != nil && err != errs.ErrKeyNotFound {
		return nil, err
	}
	if err == errs.ErrKeyNotFound {
		exist := b.BloomFilter.Exist(key)
		if exist {
			val, err = b.LoadFunc(ctx, key)
			b.Cache.Set(ctx, key, val, time.Minute)
		}
	}
	return val, err
}

type BloomCacheV1 struct {
	*ReadThroughCache
}

func NewBloomCacheV1(cache Cache, b BloomFilter, loadFunc func(ctx context.Context, key string) (any, error)) *BloomCacheV1 {
	return &BloomCacheV1{
		ReadThroughCache: &ReadThroughCache{
			Cache: cache,
			LoadFunc: func(ctx context.Context, key string) (any, error) {
				exist := b.Exist(key)
				if exist {
					return loadFunc(ctx, key)
				}
				return nil, errors.New("数据不存在")
			},
		},
	}
}

type BloomFilter interface {
	Exist(key string) bool
}
