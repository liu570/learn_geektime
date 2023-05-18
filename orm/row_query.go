package orm

import (
	"context"
	"learn_geektime/orm/internal/errs"
)

// RawQuerier 原生查询器
// 该结构体 与 Inserter, Updater, Selector 是同一类型的对象，都实现了 QueryBuilder 接口
// 但是该对象的功能与前几个对象不同, RawQuerier 同一实现了 Executor, Queier 接口
// 且上述几个类型的对象的两个接口的实现都是交由  RawQuerier 对象实现
type RawQuerier[T any] struct {
	core
	sess Session

	sql  string
	args []any

	typ string
}

func (r *RawQuerier[T]) Build() (*Query, error) {
	return &Query{
		SQL:  r.sql,
		Args: r.args,
	}, nil
}

// RawQuery 原生查询，将获取的 SQL 语句传入，交由该方法执行 SQL 语句
func RawQuery[T any](sess Session, sql string, args ...any) *RawQuerier[T] {
	return &RawQuerier[T]{
		sql:  sql,
		args: args,
		sess: sess,
		core: sess.getCore(),
	}
}

func (r *RawQuerier[T]) Get(ctx context.Context) (*T, error) {

	model, err := r.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	res := get[T](ctx, r.core, r.sess, &QueryContext{
		Type:      r.typ,
		Builder:   r,
		Model:     model,
		TableName: model.TableName,
	})
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

func (r *RawQuerier[T]) GetMulti(ctx context.Context) ([]*T, error) {
	model, err := r.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	res := getMulti[T](ctx, r.core, r.sess, &QueryContext{
		Type:      r.typ,
		Builder:   r,
		Model:     model,
		TableName: model.TableName,
	})
	if res.Result != nil {
		return res.Result.([]*T), res.Err
	}
	return nil, res.Err
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
