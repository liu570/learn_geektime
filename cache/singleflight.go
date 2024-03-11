package cache

// singleflight 模式 当 大量的对同一健值的访问请求落到数据库中时只允许其中一个进行访问的方法

type SingleflightCache struct {
	ReadThroughCache
}

/*
func NewSingleflightCache(cache Cache,
	loadFunc func(ctx context.Context, key string) (any, error)) Cache {
	g := &Singleflight.Group{}
	return &SingleflightCache{
		ReadThroughCache: ReadThroughCache{
			Cache: cache,
			LoadFunc: func(ctx context.Context, key string) (any, error) {
				defer func() {
					g.Forget(key)
				}()
				val, err, _ = g.Do(key, func() (interface{}, error) {
					return loadFunc(ctx, key)
				})
				return val, err
			},
		},
	}
}
*/
