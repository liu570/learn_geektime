package orm

import (
	"context"
	"learn_geektime/orm/internal/errs"
)

type Selector[T any] struct {
	// 组合了多种语法需要使用的相同的数据
	builder
	sess Session

	// 存储可以使用在 where 语句后面的结构体切片
	where []Predicate
	// 存储可以使用在 SELECT 语句后面的结构体切片
	columns []Selectable
	// 存储可以使用在 HAVING 语句后面的结构体切片 该数据与 WHERE 语法需要的数据几乎一样
	having []Predicate
	// 存储可以使用在 GROUP BY 语句后面的结构体切片
	groupBy []Column
	// 存储可以使用在 ORDER BY 语句后面的结构体数据
	orderBy []OrderAble

	limit  int
	offset int

	// AOP 方案设计 移动至 core 里面去
	//ms []Middleware
}

// Selectable 标记适用于 SELECT 语句下的合法结构
// 用于严格 检测适合用于 Selector 下的内容
// -- Column 列可以
// -- Aggregate
// -- OrderBy
//
//	-- RawExpr
type Selectable interface {
	// 标记方法
	selectable()
}

func NewSelector[T any](sess Session) *Selector[T] {
	c := sess.getCore()
	return &Selector[T]{
		sess: sess,
		builder: builder{
			dialect: c.dialect,
			core:    c,
		},
	}
}

func (s *Selector[T]) Use(ms ...Middleware) *Selector[T] {
	s.ms = ms
	return s
}

// Select s.Select("id", "age")
func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

func (s *Selector[T]) From(tbl TableReference) *Selector[T] {
	s.table = tbl
	return s
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) GroupBy(cs ...Column) *Selector[T] {
	s.groupBy = cs
	return s
}

func (s *Selector[T]) Having(ps ...Predicate) *Selector[T] {
	s.having = ps
	return s
}

func (s *Selector[T]) OrderBy(ods ...OrderAble) *Selector[T] {
	s.orderBy = ods
	return s
}

func (s *Selector[T]) Limit(num int) *Selector[T] {
	s.limit = num
	return s
}

func (s *Selector[T]) Offset(num int) *Selector[T] {
	s.offset = num
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	model, err := s.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	res := get[T](ctx, s.core, s.sess, &QueryContext{
		Type:    "SELECT",
		Builder: s,
		Model:   model,
		//TableName: model.TableName,
	})
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err

}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	model, err := s.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	res := getMulti[T](ctx, s.core, s.sess, &QueryContext{
		Type:      "SELECT",
		Builder:   s,
		Model:     model,
		TableName: model.TableName,
	})
	if res.Result != nil {
		return res.Result.([]*T), res.Err
	}
	return nil, res.Err
}

func (s *Selector[T]) Build() (*Query, error) {
	t := new(T)
	var err error
	s.model, err = s.core.r.Get(t)
	if err != nil {
		return nil, err
	}

	// SELECT
	s.sb.WriteString("SELECT ")
	if len(s.columns) > 0 {
		for i, column := range s.columns {
			if i > 0 {
				s.sb.WriteByte(',')
			}
			err := s.buildSelectable(column)
			if err != nil {
				return nil, err
			}
		}
	} else {
		s.sb.WriteByte('*')
	}

	// FROM
	s.sb.WriteString(" FROM ")
	err = s.buildTableReference(s.table)
	if err != nil {
		return nil, err
	}

	// WHERE
	if len(s.where) > 0 {
		// 构造 WHERE
		s.sb.WriteString(" WHERE ")
		err = s.buildExpression(s.UnionPredicates(s.where...))
		if err != nil {
			return nil, err
		}
	}

	// GROUP BY
	if len(s.groupBy) > 0 {
		s.sb.WriteString(" GROUP BY ")
		for i, column := range s.groupBy {
			if i > 0 {
				s.sb.WriteByte(',')
			}
			fd, ok := s.model.FieldMap[column.name]
			if !ok {
				return nil, errs.NewErrUnknownField(column.name)
			}
			s.quote(fd.ColName)
		}
	}

	// HAVING
	if len(s.having) > 0 {
		s.sb.WriteString(" HAVING ")
		err = s.buildExpression(s.UnionPredicates(s.having...))
		if err != nil {
			return nil, err
		}
	}

	// ORDER BY
	if len(s.orderBy) > 0 {
		s.sb.WriteString(" ORDER BY ")
		for i, order := range s.orderBy {
			if i > 0 {
				s.sb.WriteByte(',')
			}
			err = s.buildOrderAble(order)
			if err != nil {
				return nil, err
			}
		}
	}

	// LIMIT
	if s.limit > 0 {
		s.sb.WriteString(" LIMIT ?")
		s.addArgs(s.limit)
	}

	// OFFSET
	if s.offset > 0 {
		s.sb.WriteString(" OFFSET ?")
		s.addArgs(s.offset)
	}

	s.sb.WriteByte(';')
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}
