package orm

import (
	"context"
	"database/sql"
	"learn_geektime/orm/internal/errs"
	"learn_geektime/orm/internal/valuer"
	"learn_geektime/orm/model"
)

// core 统一 Tx 与 DB 都需要的数据
// 同理几个 QueryBuilder 也都需要该数据
type core struct {
	r model.Registry
	// valCreator 来确定我们是使用 反射 还是 unsafe
	valCreator valuer.Creator
	// dialect 确定我们数据库使用哪个数据库的方言
	dialect Dialect
	// AOP 方案
	ms []Middleware
}

func getMulti[T any](ctx context.Context, c core, sess Session, qc *QueryContext) *QueryResult {
	var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
		q, err := qc.Builder.Build()
		if err != nil {
			return &QueryResult{
				Result: nil,
				Err:    err,
			}
		}
		rows, err := sess.queryContext(ctx, q.SQL, q.Args...)
		if err != nil {
			return &QueryResult{
				Result: nil,
				Err:    err,
			}
		}
		var ans []*T
		for rows != nil {
			t := new(T)
			res := c.valCreator(t, qc.Model)
			err = res.SetColumn(rows)
			if err != nil {
				if err == errs.ErrNoRows {
					break
				}
				return &QueryResult{
					Result: nil,
					Err:    err,
				}
			}
			ans = append(ans, t)
		}

		return &QueryResult{
			Result: ans,
			Err:    nil,
		}
	}

	for i := len(c.ms) - 1; i >= 0; i-- {
		root = c.ms[i](root)
	}
	res := root(ctx, qc)
	return res
}

// get 该方法是 QueryBuilder 中的 Get方法的具体实现, 各个对象通过传入不同的参数来获得不同的实现
func get[T any](ctx context.Context, c core, sess Session, qc *QueryContext) *QueryResult {
	var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
		q, err := qc.Builder.Build()
		if err != nil {
			return &QueryResult{
				Result: nil,
				Err:    err,
			}
		}
		rows, err := sess.queryContext(ctx, q.SQL, q.Args...)
		if err != nil {
			return &QueryResult{
				Result: nil,
				Err:    err,
			}
		}
		t := new(T)
		res := c.valCreator(t, qc.Model)
		err = res.SetColumn(rows)
		return &QueryResult{
			Result: t,
			Err:    err,
		}
	}
	for i := len(c.ms) - 1; i >= 0; i-- {
		root = c.ms[i](root)
	}

	return root(ctx, qc)

}

func exec[T any](ctx context.Context, c core, sess Session, qc *QueryContext) *QueryResult {
	var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
		q, err := qc.Builder.Build()
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}

		res, err := sess.execContext(ctx, q.SQL, q.Args...)
		return &QueryResult{
			Result: res,
			Err:    err,
		}
	}
	for i := len(c.ms) - 1; i >= 0; i-- {
		root = c.ms[i](root)
	}
	qr := root(ctx, qc)
	var res Result
	if qr.Result != nil {
		res.res = qr.Result.(sql.Result)
	}
	return &QueryResult{
		Err:    qr.Err,
		Result: res,
	}
}
