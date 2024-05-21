package cache

import (
	"context"
	_ "embed"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var (
	//go:embed lua/unlock.lua
	luaUnlock string

	//go:embed lua/lock.lua
	luaLock string

	//go:embed lua/refresh.lua
	luaRefresh             string
	ErrFailedToPreemptLock = errors.New("redis-lock:抢锁失败")
	//ErrLockNotHold 一般是出现在你预期你本来持有锁，结果却没有持有锁的地方//比如说当你尝试释放锁的时候，可能得到这个错误
	//这一般意味着有人绕开了 rlock 的控制，直接操作了 Redis
	ErrLockNotHold = errors.New("redis-lock:未持有锁")
)

// 分布式锁的概念，与寻常锁不同的是分布式锁需要进行网络通信
// 在 redis 中分布式锁的本质就是一个键值对

// 在分布式锁中有以下问题：
// 锁的过期时间：如果没有过期时间 1 拿到锁后 2 一直等待 若 1 崩溃则锁一直在 2 无法拿到锁执行任务
// 存在过期时间： 1 拿到锁， 但 1 未完成任务时 锁就过期了 这时如果 2 拿到了锁 那么 1 会释放 2 的锁
// 所以如何判断锁是不是自己的锁：

type Client struct {
	client redis.Cmdable
}

func NewClient(cmd redis.Cmdable) *Client {
	return &Client{
		client: cmd,
	}
}

func (c *Client) Lock(ctx context.Context, key string,
	expiration time.Duration, retry RetryStrategy, timeout time.Duration) (*Lock, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	val := uuid.New().String()
	for {
		lctx, cancel := context.WithTimeout(ctx, timeout)
		res, err := c.client.Eval(lctx, luaLock, []string{key}, val, expiration.Seconds()).Result()
		cancel()

		// 成功
		if res == "OK" {
			return NewLock(c.client, key, val, expiration), nil
		}

		// 未知错误
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}

		// 超时/锁被占用
		interval, ok := retry.Next()
		if !ok {
			return nil, ErrFailedToPreemptLock
		}

		// 同时监听睡眠或 ctx 超时
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(interval):

		}
	}
}

func (c *Client) TryLock(ctx context.Context, key string,
	expiration time.Duration) (*Lock, error) {
	val := uuid.New().String()
	ok, err := c.client.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrFailedToPreemptLock
	}
	return NewLock(c.client, key, val, expiration), nil
}

type Lock struct {
	client     redis.Cmdable
	value      string
	key        string
	expiration time.Duration
	unlock     chan struct{}
	unlockOnce sync.Once
}

func NewLock(client redis.Cmdable, key string, val string, expiration time.Duration) *Lock {
	return &Lock{
		client:     client,
		value:      val,
		key:        key,
		expiration: expiration,
		unlock:     make(chan struct{}, 1),
	}
}

func (l *Lock) AutoRefresh(internal time.Duration, timeout time.Duration) error {

	ticker := time.NewTicker(internal)
	defer ticker.Stop()
	// 不断续约，直到收到退出信号
	retrySignal := make(chan struct{}, 1)
	defer close(retrySignal)
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			if err == context.DeadlineExceeded {
				// 一直重试失败怎么办
				retrySignal <- struct{}{}
				continue
			}
			if err != nil {
				return err
			}
		case <-retrySignal:
			// 重试信号
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			if err == context.DeadlineExceeded {
				// 一直重试失败怎么办
				retrySignal <- struct{}{}
				continue
			}
			if err != nil {
				return err
			}
		case <-l.unlock:
			return nil
		}
	}

}

// Refresh redis 在缓存过期时间上采取了续约的设置，
// 过期时间设置的不长，但可以在业务未执行完毕的时候进行续约操作
func (l *Lock) Refresh(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaRefresh, []string{l.key},
		l.value, l.expiration.Seconds()).Int64()
	if err == redis.Nil {
		return ErrLockNotHold
	}
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotHold
	}
	return nil
}

func (l *Lock) UnLock(ctx context.Context) error {
	// close 关闭时会向所有的 <-l.unlock 发送一次零值信号 ，用来处理并发时有多个 AutoRefresh 调用 只传一个信号不够
	l.unlockOnce.Do(func() {
		close(l.unlock)
	})
	// 由于下列释放锁 有对键值对进行 get del 两个操作有时间差别的问题
	// 这里考虑 redis 是 单线程 所以使用 lua 脚本 实现类似原子操作的功能
	res, err := l.client.Eval(ctx, luaUnlock, []string{l.key}, l.value).Int64()
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotHold
	}
	return nil
	/*	value, err := l.client.Get(ctx, key).Result()
		if err != nil {
			return err
		}
		if value != l.value {
			// 这里代表锁不是你的锁
			return ErrLockNotHold
		}
		res, err := l.client.Del(ctx, key).Result()
		if err != nil {
			return err
		}
		if res != 1 {
			return ErrLockNotHold
		}
		return nil
	*/
}
