package orm

import (
	"errors"
	"learn_geektime/orm/internal/errs"
	"learn_geektime/orm/model"
	"strings"
)

// builder 用于封装一些轻量级的代码 用于简化代码
// 同时持有一些 QueryBuilder 的公共字段 即是 Selector、Insertor、Updater、Delector 中共有的字段
type builder struct {
	// core 是用来整合 DB, Tx 所共需的数据
	core
	sb      strings.Builder
	model   *model.Model
	args    []any
	dialect Dialect
	// 确定所对应的表名 如果为空则表明为对应的结构体名字
	table TableReference
}

// quote 给 name 标记上引号
func (b *builder) quote(name string) {
	b.sb.WriteByte(b.dialect.quoter())
	b.sb.WriteString(name)
	b.sb.WriteByte(b.dialect.quoter())
}

//func (b *builder) buildAssignable(assignable Assignable) error {
//	switch expr := assignable.(type) {
//	case Column:
//		fd, ok := b.model.FieldMap[expr.name]
//		if !ok {
//			return errs.NewErrUnknownField(expr.name)
//		}
//		b.quote(fd.ColName)
//		b.sb.WriteString(" = ?")
//		// TODO 有待商榷
//		b.args = append(b.args, reflect.ValueOf(b.values[0]).Elem().Field(fd.Index).Interface())
//	case Assignment:
//		return b.buildAssignment(expr)
//
//	}
//}

func (b *builder) buildSelectable(se Selectable) error {
	switch expr := se.(type) {
	case Column:
		return b.buildColumn(expr)
	case Aggregate:
		return b.buildAggregate(expr)
	case RawExpr:
		return b.buildRawExpr(expr)
	default:
		return errs.NewErrUnsupportedSelectable(expr)
	}

}

func (b *builder) buildOrderAble(order OrderAble) error {
	switch expr := order.(type) {
	case Column:
		return b.buildColumn(expr)
	case OrderBy:
		return b.buildOrderBy(expr)
	}
	return nil
}

func (b *builder) buildExpression(expression Expression) error {
	// 因为 有多个结构体实现了 Expression 接口所以我们需要 switch 判断传入的类型
	switch expr := expression.(type) {
	case Column:
		return b.buildColumn(expr)
	case Value:
		return b.buildValue(expr)
	case Predicate:
		return b.buildPredicate(expr)
	case Aggregate:
		return b.buildAggregate(expr)
	case RawExpr:
		return b.buildRawExpr(expr)
	case nil:
		return nil
	default:
		return errs.NewErrUnsupportedExpressionType(expr)
	}
}

func (b *builder) buildTableReference(t TableReference) error {
	switch tbl := t.(type) {
	case nil:
		b.quote(b.model.TableName)
	case Table:
		m, err := b.r.Get(tbl.entity)
		if err != nil {
			return err
		}
		b.quote(m.TableName)
		if len(tbl.alias) > 0 {
			b.sb.WriteString(" AS ")
			b.quote(tbl.alias)
		}
	case Join:
		b.sb.WriteByte('(')
		// 左边一张表，右边一张表
		err := b.buildTableReference(tbl.left)
		if err != nil {
			return err
		}
		b.sb.WriteByte(' ')
		b.sb.WriteString(tbl.typ)
		b.sb.WriteByte(' ')
		err = b.buildTableReference(tbl.right)
		if err != nil {
			return err
		}

		// ON
		if len(tbl.on) > 0 {
			b.sb.WriteString(" ON ")
			err = b.buildExpression(b.UnionPredicates(tbl.on...))
			if err != nil {
				return err
			}
		}

		// USING
		if len(tbl.using) > 0 {
			b.sb.WriteString(" USING ")
			b.sb.WriteByte('(')
			for i, us := range tbl.using {
				if i > 0 {
					b.sb.WriteByte(',')
				}
				colName, err := b.colName(tbl, us)
				if err != nil {
					return err
				}
				b.quote(colName)
			}
			b.sb.WriteByte(')')
		}

		b.sb.WriteByte(')')
	}
	return nil
}

// -------------------------------------------------------------------------------------分割线-------------------------------------------------------------------------------------
// 下面是为实现结构体的方法

func (b *builder) buildAssignment(assignment Assignment) error {
	fd, ok := b.model.FieldMap[assignment.Column]
	if !ok {
		return errs.NewErrUnknownField(assignment.Column)
	}
	b.quote(fd.ColName)
	b.sb.WriteString(" = ")
	switch expr := assignment.val.(type) {
	case Predicate:
		return b.buildExpression(expr)
	case RawExpr:
		b.sb.WriteString(expr.raw)
		b.args = append(b.args, expr.args...)
	case Value:
		b.sb.WriteByte('?')
		if b.args == nil {
			b.args = make([]any, 0, 8)
		}
		b.args = append(b.args, expr.val)

		//default:
		//b.sb.WriteByte('?')
		//b.args = append(b.args, expr.val)
	}
	return nil
}

// buildColumn 通过该方法构建 列名
func (b *builder) buildColumn(c Column) error {
	alias := ""
	if c.table != nil {
		alias = c.table.tableAlias()
	}
	colName, err := b.colName(c.table, c.name)
	if err != nil {
		return err
	}
	if len(alias) > 0 {
		b.quote(c.table.tableAlias())
		b.sb.WriteByte('.')
	}
	b.quote(colName)

	return err
}

// colName 返回元数据对应的列名 colName.
// 传入的参数为 TableReference, Column.name
func (b *builder) colName(table TableReference, fd string) (string, error) {
	switch tbl := table.(type) {
	case nil:
		// 用户没有调用 FROM 方法
		meta, ok := b.model.FieldMap[fd]
		if !ok {
			return "", errs.NewErrUnknownField(fd)
		}
		return meta.ColName, nil
	case Table:
		m, err := b.r.Get(tbl.entity)
		if err != nil {
			return "", err
		}
		meta, ok := m.FieldMap[fd]
		if !ok {
			return "", errs.NewErrUnknownField(fd)
		}
		return meta.ColName, nil
	case Join:
		colName, err := b.colName(tbl.left, fd)
		if err == nil {
			return colName, nil
		}
		return b.colName(tbl.right, fd)
	default:
		return "", errors.New("orm: 错误的表")
	}
}

func (b *builder) buildValue(value Value) error {
	b.sb.WriteByte('?')
	if b.args == nil {
		b.args = make([]any, 0, 8)
	}
	b.args = append(b.args, value.val)

	return nil
}

func (b *builder) buildPredicate(p Predicate) error {
	_, ok := p.left.(Predicate)
	if ok {
		b.sb.WriteByte('(')
	}
	if err := b.buildExpression(p.left); err != nil {
		return err
	}
	if ok {
		b.sb.WriteByte(')')
	}

	if p.op != "" {
		b.sb.WriteByte(' ')
		b.sb.WriteString(p.op.String())
		b.sb.WriteByte(' ')
	}
	_, ok = p.right.(Predicate)
	if ok {
		b.sb.WriteByte('(')
	}
	if err := b.buildExpression(p.right); err != nil {
		return err
	}
	if ok {
		b.sb.WriteByte(')')
	}
	return nil
}

// buildAggregate 用于构建 SQL 语句中的聚合函数部分
func (b *builder) buildAggregate(agg Aggregate) error {
	// 使用元数据校验 该列是否是 结构体中有的列
	fd, ok := b.model.FieldMap[agg.arg]
	if !ok {
		return errs.NewErrUnknownField(agg.arg)
	}
	b.sb.WriteString(agg.fn)
	b.sb.WriteByte('(')
	b.quote(fd.ColName)
	b.sb.WriteByte(')')
	return nil
}

func (b *builder) buildRawExpr(raw RawExpr) error {
	b.sb.WriteString(raw.raw)
	b.addArgs(raw.args...)
	return nil
}

func (b *builder) buildOrderBy(order OrderBy) error {
	fd, ok := b.model.FieldMap[order.col]
	if !ok {
		return errs.NewErrUnknownField(order.col)
	}
	b.quote(fd.ColName)
	b.sb.WriteByte(' ')
	b.sb.WriteString(order.order)
	return nil
}

func (b *builder) UnionPredicates(pds ...Predicate) Predicate {
	var pred Predicate
	for i, p := range pds {
		if i == 0 {
			pred = p
			continue
		}
		pred = pred.And(p)
	}
	return pred
}

func (b *builder) addArgs(args ...any) {
	if len(args) == 0 {
		return
	}
	if b.args == nil {
		b.args = make([]any, 0, 8)
	}
	b.args = append(b.args, args...)
}
