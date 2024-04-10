package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ClientE2ESuite struct {
	suite.Suite
	rdb redis.Cmdable
}

func (s *ClientE2ESuite) SetupSuite() {
	s.rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	for s.rdb.Ping(context.Background()).Err() != nil {

	}
}

func TestClientE2E(t *testing.T) {
	suite.Run(t, &ClientE2ESuite{})
}

func (s *ClientE2ESuite) TestTryLock() {
	t := s.T()
	rdb := s.rdb
	client := &Client{
		client: rdb,
	}
	testCases := []struct {
		name string

		key        string
		expiration time.Duration

		wantLock *Lock
		wantErr  error

		before func()
		after  func()
	}{
		{
			// 加锁成功
			name:       "locked",
			key:        "locked-key",
			expiration: time.Minute,
			before:     func() {},
			after: func() {
				res, err := rdb.Del(context.Background(), "locked-key").Result()
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
			},
			wantLock: &Lock{
				key: "locked-key",
			},
		}, {
			//	模拟并发竞争失败
			name:       "failed",
			key:        "failed-key",
			expiration: time.Minute,
			before: func() {
				//	假设已经有人设置了分布式锁
				val, err := rdb.Set(context.Background(), "failed-key", "123", time.Minute).Result()
				require.NoError(t, err)
				require.Equal(t, "OK", val)
			},
			after: func() {
				res, err := rdb.Del(context.Background(), "failed-key").Result()
				require.NoError(t, err)
				require.Equal(t, int64(1), res)
			},
			wantErr: ErrFailedToPreemptLock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			l, err := client.TryLock(context.Background(), tc.key, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.key, l.key)
			assert.NotEmpty(t, l.value)
			tc.after()
		})
	}
}

func (s *ClientE2ESuite) TestUnLock() {
	t := s.T()
	rdb := s.rdb
	client := &Client{
		client: rdb,
	}
	testCases := []struct {
		name string

		lock *Lock

		before func()
		after  func()

		wantLock *Lock
		wantErr  error
	}{
		{
			name: "unlocked",
			lock: func() *Lock {
				l, err := client.TryLock(context.Background(), "unlocked-key", time.Minute)
				require.NoError(t, err)
				return l
			}(),
			before: func() {},
			after: func() {
				res, err := rdb.Exists(context.Background(), "unlocked-key").Result()
				require.NoError(t, err)
				require.Equal(t, 1, res)
			},
		},
		{
			name: "lock not hold",
			lock: &Lock{
				key:   "not-hold-key",
				value: "123",
			},
			wantErr: ErrLockNotHold,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.lock.unLock(context.Background())
			require.Equal(t, tc.wantErr, err)
		})
	}
}
