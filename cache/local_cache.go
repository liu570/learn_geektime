package cache

import (
	"context"
	"learn_geektime/cache/internal/errs"
	"sync"
	"time"
)

type LocalCache struct {
	m sync.Map
}

func (l *LocalCache) Get(ctx context.Context, key string) (any, error) {
	val, ok := l.m.Load(key)
	if !ok {
		return nil, errs.NewErrKeyNotFound(key)
	}
	return val, nil
}

func (l *LocalCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	l.m.Store(key, val)
	return nil
}

func (l *LocalCache) Delete(ctx context.Context, key string, val any) error {
	l.m.Delete(key)
	return nil
}
