package orm

import (
	"learn_geektime/orm/internal/errs"
)

var (
	DialectMySQL   Dialect = &mysqlDialect{}
	DialectSqlite  Dialect = &sqliteDialect{}
	Dialectpostgre Dialect = &postgreDialect{}
)

// Dialect 方言，构造不同数据库个性部分
type Dialect interface {
	// 引号
	quoter() byte
	buildUpsert(b *builder, upsert *Upsert) error
}

// standardSQL SQL标准的方言实现
type standardSQL struct {
}

func (s standardSQL) quoter() byte {
	//TODO implement me
	panic("implement me")
}

func (s standardSQL) buildUpsert(b *builder, upsert *Upsert) error {
	//TODO implement me
	panic("implement me")
}

// ---------------------------------------------------------------------------------------- 分隔符 ---------------------------------------------------------------------------------------------------

type UpsertBuilder[T any] struct {
	// 链式调用返回 INSERT 语句
	i       *Inserter[T]
	assigns []Assignable

	// 复杂语句
	//where []Predicate
	conflictColumns []string
}
type Upsert struct {
	assigns         []Assignable
	conflictColumns []string
	//doNothing bool
}

func (o *UpsertBuilder[T]) ConflictColumns(cols ...string) *UpsertBuilder[T] {
	o.conflictColumns = cols
	return o
}

//	func (o *onDuplicateKeyBuilder[T]) Where(ps ...Predicate) *onDuplicateKeyBuilder[T] {
//		o.where = ps
//		return o
//	}
func (o *UpsertBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	o.i.onDuplicate = &Upsert{
		conflictColumns: o.conflictColumns,
		assigns:         assigns,
	}
	return o.i
}

//func (o *onDuplicateKeyBuilder[T]) DoNothing(assigns ...Assignable) *Inserter[T] {
//	o.i.onDuplicate = &onDuplicateKey{
//		doNothing: true,
//	}
//	return o.i
//}

// ---------------------------------------------------------------------------------------- 分隔符 ---------------------------------------------------------------------------------------------------

// mysqlDialect 用于实现 mysql 的方言
type mysqlDialect struct {
	standardSQL
}

func (m *mysqlDialect) quoter() byte {
	return '`'
}

func (m *mysqlDialect) buildUpsert(b *builder, upsert *Upsert) error {
	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for idx, assign := range upsert.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}

		// buildAssignable
		switch expr := assign.(type) {
		case Assignment:
			err := b.buildAssignment(expr)
			if err != nil {
				return err
			}
		case Column:
			fd, ok := b.model.FieldMap[expr.name]
			if !ok {
				return errs.NewErrUnknownField(expr.name)
			}
			b.quote(fd.ColName)
			b.sb.WriteString(" = ")
			b.sb.WriteString("VALUES")
			b.sb.WriteByte('(')
			b.quote(fd.ColName)
			b.sb.WriteByte(')')
		}
	}
	return nil
}

// ---------------------------------------------------------------------------------------- 分隔符 ---------------------------------------------------------------------------------------------------

// sqliteDialect 用于实现 sqlite3 的方言
type sqliteDialect struct {
	standardSQL
}

func (dialect *sqliteDialect) quoter() byte {
	return '`'
}

func (dialect *sqliteDialect) buildUpsert(b *builder, upsert *Upsert) error {

	b.sb.WriteString(" ON CONFLICT")
	if len(upsert.conflictColumns) > 0 {
		b.sb.WriteByte('(')
		for i, column := range upsert.conflictColumns {
			if i > 0 {
				b.sb.WriteByte(',')
			}
			fd, ok := b.model.FieldMap[column]
			if !ok {
				return errs.NewErrUnknownField(column)
			}
			b.quote(fd.ColName)
		}
		b.sb.WriteByte(')')
	}
	b.sb.WriteString(" DO UPDATE SET ")
	for idx, assign := range upsert.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		// buildAssignable
		switch expr := assign.(type) {
		case Assignment:
			err := b.buildAssignment(expr)
			if err != nil {
				return err
			}
		case Column:
			fd, ok := b.model.FieldMap[expr.name]
			if !ok {
				return errs.NewErrUnknownField(expr.name)
			}
			b.quote(fd.ColName)
			b.sb.WriteString(" = ")
			b.sb.WriteString("excluded.")
			b.quote(fd.ColName)

		}
	}
	return nil
}

type postgreDialect struct {
	standardSQL
}
