package cache

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"learn_geektime/cache/mocks"
	"testing"
	"time"
)

// mockgen -package=mocks -destination=mocks/redis_cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
// -package=mocks : 指定包名为 mocks
// -destination=mocks/redis_cmdable.mock.go ： 指定文件位置是当前目录下的 mocks 下的 redis_cmdable.mock.go
// github.com/redis/go-redis/v9 Cmdable : mock 的是 github.com/redis/go-redis/v9 库的 Cmdable 类

func TestRedisCache_Set(t *testing.T) {
	// ctrl 是 gomock 机制要求我们要有这个
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name       string
		mock       func() redis.Cmdable
		key        string
		val        string
		expiration time.Duration

		wantErr error
	}{
		{
			name: "return ok",
			mock: func() redis.Cmdable {
				res := mocks.NewMockCmdable(ctrl)
				// 因为 Set 方法需要返回一个 *StatusCmd 所以 调用该方法
				cmd := redis.NewStatusCmd(nil)
				cmd.SetVal("OK")
				// 下列代表我们期待 当 调用set 的时候会返回一个 cmd (cmd 是 redis 中的 set 的期待返回 所以我们直接用 redis 里面的)
				res.EXPECT().Set(gomock.Any(), "key1", "value1", time.Minute).
					Return(cmd)
				return res
			},
			key:        "key1",
			val:        "value1",
			expiration: time.Minute,
		},
		{
			name: "time out",
			mock: func() redis.Cmdable {
				res := mocks.NewMockCmdable(ctrl)
				// 因为 Set 方法需要返回一个 *StatusCmd 所以 调用该方法
				cmd := redis.NewStatusCmd(nil)
				cmd.SetErr(context.DeadlineExceeded)
				res.EXPECT().Set(gomock.Any(), "key1", "value1", time.Minute).
					Return(cmd)
				return res
			},
			key:        "key1",
			val:        "value1",
			expiration: time.Minute,
			wantErr:    context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmdable := tc.mock()
			client := NewRedisCache(cmdable)
			err := client.Set(context.Background(), tc.key, tc.val, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
		})
	}
}
