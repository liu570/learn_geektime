package memory

import (
	"context"
	"errors"
	"sync"
)

type Session struct {
	id string
	// 用户打开多个标签页 同时操作你的 get 和 set 就需要加锁保护
	mutex  sync.RWMutex
	values map[string]string
}

func (s *Session) Get(ctx context.Context, key string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	val, ok := s.values[key]
	if ok {
		return val, nil
	}
	return "", errors.New("web: key not found")
}

func (s *Session) Set(ctx context.Context, key string, val string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.values[key] = val
	return nil
}

func (s *Session) ID() string {
	return s.id
}
