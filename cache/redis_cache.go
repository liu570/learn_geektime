package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client redis.Cmdable
}

func NewRedisCacheWithAddr(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // 没有设置密码
		DB:       0,  // 使用默认的DB
	})
	return NewRedisCache(client)
}

func NewRedisCache(client redis.Cmdable) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) (any, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	res, err := r.client.Set(ctx, key, val, expiration).Result()

	if err != nil {
		return err
	}
	if res != "OK" {
		return errors.New("cache: 设置键值对失败")
	}
	return err
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}
