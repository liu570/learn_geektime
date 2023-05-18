package orm

import "context"

// RawQuerier 原生查询器
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
