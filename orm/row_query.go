package orm

import (
	"context"
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
		Type:      "RAW",
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
		Type:      "RAW",
		Builder:   r,
		Model:     model,
		TableName: model.TableName,
	})
	if res.Result != nil {
		return res.Result.([]*T), res.Err
	}
	return nil, res.Err
}

func (r *RawQuerier[T]) Exec(ctx context.Context) Result {
	model, err := r.r.Get(new(T))
	if err != nil {
		return Result{
			err: err,
		}
	}
	res := exec[T](ctx, r.core, r.sess, &QueryContext{
		Type:      "RAW",
		Builder:   r,
		Model:     model,
		TableName: model.TableName,
	})

	if res.Err != nil {
		return Result{err: res.Err}
	}
	return res.Result.(Result)
}
