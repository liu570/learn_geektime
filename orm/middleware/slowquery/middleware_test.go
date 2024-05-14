package slowquery

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"learn_geektime/orm"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := &MiddlewareBuilder{}
	builder.LogFunc(func(sql string, args ...any) {
		// 在此地你可以操作一些东西
		// 注意 此地 args 是敏感信息、不应该打在日志上 需要进行处理
		fmt.Println(sql)
	}).SlowQueryThreshold(100)
	db, err := orm.Open("sqlite3", "file:test.db?cache=sharad&mode=memory",
		orm.DBWithMiddlewares(builder.Build(), func(next orm.Handler) orm.Handler {
			return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
				time.Sleep(time.Millisecond)
				return next(ctx, qc)
			}
		}))
	if err != nil {
		t.Fatal(err)
	}
	// 如何指定 middleware
	_, err = orm.NewSelector[TestModel](db).Get(context.Background())
	assert.NotNil(t, err)
	//require.NoError(t, err)
}

type TestModel struct {
}
