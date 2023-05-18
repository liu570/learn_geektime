package memory

import (
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"learn_geektime/web/session"
	"sync"
	"time"
)

type Store struct {
	// 如果我们可以确保下列操作中生成的都是不同的 id 则我们不加锁也没什么问题
	// 在这里我认为是不用加锁的，因为我发现 cache 包已经自己加了锁是线程安全的（大明老师加了锁）
	// 我认为我们不需要重复加锁
	mutex      sync.RWMutex
	expiration time.Duration
	cache      *cache.Cache
}

func NewStore(expiration time.Duration) *Store {
	return &Store{
		// cache.New 返回一个过期时间为 expiration 和 定期清理时间间隔为 time.second 的缓存
		cache: cache.New(expiration, time.Second),
	}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	res := &Session{
		id:     id,
		values: make(map[string]string),
	}
	s.cache.Set(id, res, s.expiration)
	return res, nil
}

func (s *Store) Remove(ctx context.Context, id string) error {

	s.cache.Delete(id)
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	// cache.Get 帮我们在获取 session 的时候帮我们检查 session 是否过期
	res, ok := s.cache.Get(id)
	if !ok {
		return nil, errors.New("web/store: session not found")
	}
	return res.(session.Session), nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	res, ok := s.cache.Get(id)
	if !ok {
		return errors.New("web/store: session not found")
	}
	s.cache.Set(id, res, s.expiration)
	return nil
}
