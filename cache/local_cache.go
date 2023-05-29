package cache

import (
	"context"
	"learn_geektime/cache/internal/errs"
	"sync"
	"time"
)

type LocalCache struct {
	m         sync.Map
	close     chan struct{}
	closeOnce sync.Once
}

func NewLocalCache() *LocalCache {
	res := &LocalCache{
		close: make(chan struct{}, 1),
	}
	// 间隔时间, 过长则过期的缓存迟迟得不到删除
	// 过短，则频繁执行，效果不好（过期的 key 很少）
	ticker := time.NewTicker(time.Second)
	go func() {
		// 没有时间间隔，不断遍历
		for {
			select {
			case <-ticker.C:
				res.m.Range(func(key, value any) bool {
					// 如果过期了
					itm := value.(*item)
					// time.Now() 是一个很慢的调用
					if itm.deadline.Before(time.Now()) {
						res.m.Delete(key)
					}
					return true
				})
			}
		}
	}()
	return res
}

func (l *LocalCache) Get(ctx context.Context, key string) (any, error) {
	val, ok := l.m.Load(key)
	if !ok {
		return nil, errs.NewErrKeyNotFound(key)
	}
	return val, nil
}

func (l *LocalCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	l.m.Store(key, &item{
		val:      val,
		deadline: time.Now().Add(expiration),
	})
	return nil
}

func (l *LocalCache) Delete(ctx context.Context, key string, val any) error {
	l.m.Delete(key)
	return nil
}

func (l *LocalCache) Close() error {

	// 关闭只能关闭一次
	l.closeOnce.Do(func() {
		l.close <- struct{}{}
		close(l.close)
	})

	// 也是一个方法 但是没有上面方法保险
	//select {
	//case l.close <- struct{}{}:
	//	close(l.close)
	//default:
	//	return errors.New("cache: 已经被关闭了")
	//}
	return nil
}

// 可以考虑 用 sync.pool 来复用
type item struct {
	val      any
	deadline time.Time
}
