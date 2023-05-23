package orm

import (
	"context"
	"learn_geektime/orm/internal/errs"
	"reflect"
)

type Updater[T any] struct {
	builder
	core
	sess   Session
	values []*T
	// 存储可以用在 SET 语句后面的结构体切片
	assigns []Assignable
	// 存储可以使用在 where 语句后面的结构体切片
	where []Predicate
}

func (u *Updater[T]) Exec(ctx context.Context) Result {
	// 这里调用 Get 方法获取元数据，是怕 Updater 中的元数据为空
	model, err := u.r.Get(new(T))
	if err != nil {
		return Result{
			err: err,
		}
	}
	qr := exec[T](ctx, u.sess.getCore(), u.sess, &QueryContext{
		Type:      "UPDATE",
		Builder:   u,
		Model:     model,
		TableName: model.TableName,
	})
	if qr.Err != nil {
		return Result{
			err: qr.Err,
		}
	}
	return qr.Result.(Result)
}

func NewUpdater[T any](sess Session) *Updater[T] {
	c := sess.getCore()
	return &Updater[T]{
		sess: sess,
		core: c,
		builder: builder{
			dialect: c.dialect,
		},
	}
}

func (u *Updater[T]) Update(vals ...*T) *Updater[T] {
	u.values = vals
	return u
}

func (u *Updater[T]) Set(cols ...Assignable) *Updater[T] {
	u.assigns = cols
	return u
}

func (u *Updater[T]) Where(ps ...Predicate) *Updater[T] {
	u.where = ps
	return u
}

func (u *Updater[T]) Build() (*Query, error) {
	// UPDATE
	u.sb.WriteString("UPDATE ")
	if len(u.values) == 0 {
		return nil, errs.ErrNoUpdatedColumns
	}
	var err error
	u.model, err = u.r.Get(u.values[0])
	if err != nil {
		return nil, err
	}
	u.builder.quote(u.model.TableName)

	// SET
	if len(u.assigns) > 0 {
		u.sb.WriteString(" SET ")
		for i, col := range u.assigns {
			if i > 0 {
				u.sb.WriteByte(',')
			}
			switch expr := col.(type) {
			case Column:
				fd, ok := u.model.FieldMap[expr.name]
				if !ok {
					return nil, errs.NewErrUnknownField(expr.name)
				}
				u.builder.quote(fd.ColName)
				u.sb.WriteString(" = ?")
				// TODO 有待商榷
				u.args = append(u.args, reflect.ValueOf(u.values[0]).Elem().Field(fd.Index).Interface())
			case Assignment:
				err = u.buildAssignment(expr)
				if err != nil {
					return nil, err
				}
			}

		}
	}

	// WHERE
	if len(u.where) > 0 {
		// 构造 WHERE
		u.sb.WriteString(" WHERE ")
		err = u.buildExpression(u.UnionPredicates(u.where...))
		if err != nil {
			return nil, err
		}

	}

	// ORDER BY

	// LIMIT

	u.sb.WriteByte(';')
	return &Query{
		SQL:  u.sb.String(),
		Args: u.args,
	}, nil
}
