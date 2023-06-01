package cache

// mockgen -package=mocks -destination=mocks/redis_cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
// -package=mocks : 指定包名为 mocks
// -destination=mocks/redis_cmdable.mock.go ： 指定文件位置是当前目录下的 mocks 下的 redis_cmdable.mock.go
// github.com/redis/go-redis/v9 Cmdable : mock 的是 github.com/redis/go-redis/v9 库的 Cmdable 类
import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client redis.Cmdable
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
