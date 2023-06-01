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
	//		}
	//	},
	//}
}

type User struct {
	name string
}
