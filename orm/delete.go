package orm

import "context"

type Deleter[T any] struct {
	builder
	// DELETE *
	// or DELETE(*,*)
	sess    Session
	where   []Predicate
	orderBy []OrderAble

	limit int
}

func NewDeleter[T any](sess Session) *Deleter[T] {
	c := sess.getCore()
	return &Deleter[T]{
		sess: sess,
		builder: builder{
			dialect: c.dialect,
			core:    c,
		},
	}
}

func (d *Deleter[T]) Exec(ctx context.Context) Result {
	model, err := d.r.Get(new(T))
	if err != nil {
		return Result{
			err: err,
		}
	}
	qr := exec[T](ctx, d.sess.getCore(), d.sess, &QueryContext{
		Type:      "DELETE",
		Builder:   d,
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

func (d *Deleter[T]) From(tbl TableReference) *Deleter[T] {
	d.table = tbl
	return d
}

func (d *Deleter[T]) Where(ps ...Predicate) *Deleter[T] {
	d.where = ps
	return d
}
func (d *Deleter[T]) OrderBy(ods ...OrderAble) *Deleter[T] {
	d.orderBy = ods
	return d
}

func (d *Deleter[T]) Limit(num int) *Deleter[T] {
	d.limit = num
	return d
}

func (d *Deleter[T]) Build() (*Query, error) {
	model, err := d.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	d.model = model
	// DELETE
	d.sb.WriteString("DELETE")

	// FROM
	d.sb.WriteString(" FROM ")
	err = d.buildTableReference(d.table)
	if err != nil {
		return nil, err
	}

	// WHERE
	if len(d.where) > 0 {
		d.sb.WriteString(" WHERE ")
		err := d.buildExpression(d.UnionPredicates(d.where...))
		if err != nil {
			return nil, err
		}
	}

	// ORDER BY
	if len(d.orderBy) > 0 {
		d.sb.WriteString(" ORDER BY ")
		for i, order := range d.orderBy {
			if i > 0 {
				d.sb.WriteByte(',')
			}
			err = d.buildOrderAble(order)
			if err != nil {
				return nil, err
			}
		}
	}

	// LIMIT
	if d.limit > 0 {
		d.sb.WriteString(" LIMIT ?")
		d.addArgs(d.limit)
	}

	return &Query{
		SQL:  d.sb.String(),
		Args: d.args,
	}, nil
}
