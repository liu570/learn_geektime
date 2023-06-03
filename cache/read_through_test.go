package cache

import (
	"testing"
)

func TestReadThroughCache_Get(t *testing.T) {
	//var db *orm.DB
	//cache := &ReadThroughCache{
	//	Expiration: time.Minute,
	//	LoadFunc: func(ctx context.Context, key string) (any, error) {
	//		if strings.HasPrefix(key, "/user/") {
	//			// 找用户的数据
	//			// key = /user/123 ， 其中123 是用户 id
	//			// 这是用户的
	//			return orm.NewSelector[User](db).Get(ctx)
	//		} else if strings.HasPrefix(key, "/order/") {
	//			// 找 order 数据
	//		}
	//	},
	//}

	// 通用不行 就选择抽象出来 但此时我们自然就想到了 使用泛型
	//userCache := &ReadThroughCache{
	//	Expiration: time.Minute,
	//	LoadFunc: func(ctx context.Context, key string) (any, error) {
	//		if strings.HasPrefix(key, "/user/") {
	//			// 找用户的数据
	//			// key = /user/123 ， 其中123 是用户 id
	//			// 这是用户的
	//			return orm.NewSelector[User](db).Get(ctx)
	//		}
	//	}
	//}
}

type User struct {
	name string
}
