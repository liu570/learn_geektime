package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	//Incre
}

type CacheV2[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	Set(ctx context.Context, key string, val T, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	//Incre
}
