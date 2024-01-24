package cache

import (
	"context"
	"errors"
	"fmt"
	"learn_geektime/orm"
	"strings"
	"testing"
	"time"
)

func TestReadThroughCache_Get(t *testing.T) {
	var db *orm.DB

	local := NewLocalCache(func(key string, val any) {

	})
	cache := &ReadThroughCache{
		Cache:      local,
		Expiration: time.Minute,
		LoadFunc: func(ctx context.Context, key string) (any, error) {
			if strings.HasPrefix(key, "/user/") {
				// 找用户的数据
				// key = /user/123 ， 其中123 是用户 id
				// 这是用户的
				return orm.NewSelector[User](db).Get(ctx)
			} else if strings.HasPrefix(key, "/order/") {
				// 找 order 数据
			}
			return nil, errors.New("不支持操作")
		},
	}
	cache.Get(context.Background(), "/user/123")

	//通用不行 就选择抽象出来 但此时我们自然就想到了 使用泛型
	userCache := &ReadThroughCache{
		Cache:      local,
		Expiration: time.Minute,
		LoadFunc: func(ctx context.Context, key string) (any, error) {
			if strings.HasPrefix(key, "/user/") {
				// 找用户的数据
				// key = /user/123 ， 其中123 是用户 id
				// 这是用户的
				return orm.NewSelector[User](db).Get(ctx)
			}
			return nil, errors.New("不支持操作")
		},
	}
	userCache.Get(context.Background(), "/user/123")

	userCacheV1 := &ReadThroughCacheV1[*User]{
		Cache:      local,
		Expiration: time.Minute,
		LoadFunc: func(ctx context.Context, key string) (*User, error) {
			if strings.HasPrefix(key, "/user/") {
				// 找用户的数据
				// key = /user/123 ， 其中123 是用户 id
				// 这是用户的
				return orm.NewSelector[User](db).Get(ctx)
			}
			return nil, errors.New("不支持操作")
		},
	}
	userCacheV1.Get(context.Background(), "/user/123")

	userCacheV2 := &ReadThroughCacheV2[*User]{
		//Cache: local,
		Expiration: time.Minute,
		LoadFunc: func(ctx context.Context, key string) (*User, error) {
			if strings.HasPrefix(key, "/user/") {
				// 找用户的数据
				// key = /user/123 ， 其中123 是用户 id
				// 这是用户的
				return orm.NewSelector[User](db).Get(ctx)
			}
			return nil, errors.New("不支持操作")
		},
	}
	user, _ := userCacheV2.Get(context.Background(), "/user/123")
	fmt.Print(user.Name)
}

type User struct {
	Name string
}
