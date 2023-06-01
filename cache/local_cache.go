package cache

import (
	"context"
	"learn_geektime/cache/internal/errs"
	"sync"
	"time"
)

type LocalCache struct {
	// 该方式是使用 map 来对缓存进行存储，若想提高内存利用率最好不要使用 map 可以 使用序列化协议来存储数据
	data      map[string]any
	mutex     sync.RWMutex
	close     chan struct{}
	closeOnce sync.Once

	// 回调方法
	//onEvicted func(key string, val any) error
	onEvicted func(key string, val any)
	//onEvicted func(ctx context.Context, key string, val any) error
	//onEvicted func(ctx context.Context, key string, val any)
	//onEvicteds []func(key string, val any)
}
type KVOption func(key string, val any)

func NewLocalCache(option ...KVOption) *LocalCache {
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
				// 引入一个计数器，防止程序在删除缓存的时候占用系统太多资源而导致系统反应速度降低
				cnt := 0
				res.mutex.Lock()
				for k, v := range res.data {
					if v.(*item).deadline.Before(time.Now()) {
						//delete(res.data, k)
						res.delete(k, v.(*item).val)
					}
					cnt++
					if cnt >= 1000 {
						break
					}
				}
				res.mutex.Unlock()
			}
		}
	}()
	return res
}

func (l *LocalCache) Get(ctx context.Context, key string) (any, error) {

	// double check 写法 先查看数据在不在，在的话再查看数据有没有过期
	l.mutex.RLock()
	val, ok := l.data[key]
	l.mutex.RUnlock()
	if !ok {
		return nil, errs.NewErrKeyNotFound(key)
	}

	itm := val.(*item)
	if itm.deadline.Before(time.Now()) {
		l.mutex.Lock()
		defer l.mutex.Unlock()
		val, ok = l.data[key]
		if !ok {
			return nil, errs.NewErrKeyNotFound(key)
		}
		itm = val.(*item)
		if itm.deadline.Before(time.Now()) {
			l.delete(key, itm.val)
			return nil, errs.NewErrKeyNotFound(key)
		}
	}
	return itm.val, nil
}

func (l *LocalCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.data[key] = &item{
		val:      val,
		deadline: time.Now().Add(expiration),
	}
	return nil
}

func (l *LocalCache) Delete(ctx context.Context, key string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	val, ok := l.data[key]
	if !ok {
		return nil
	}
	l.delete(key, val.(*item).val)
	return nil
}
func (l *LocalCache) delete(key string, val any) {
	delete(l.data, key)
	if l.onEvicted != nil {
		l.onEvicted(key, val)
	}
}
func (l *LocalCache) Close() error {

	// 关闭只能关闭一次
	l.closeOnce.Do(func() {
		l.close <- struct{}{}
		close(l.close)
	})

	// 也是一个方法 但是没有上面方法保险
	// 采用的是 select + default 防止多次 close 阻塞调用者
	//select {
	//case l.close <- struct{}{}:
	// // 关闭channel的 时候需要小心，如果发送了数据到已经关闭的 channel 会引起 panic
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
