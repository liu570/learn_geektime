package orm

// TableReference 抽象的接口用于支持复杂查询 join 查询，table，子查询
type TableReference interface {
	tableAlias() string
}

//-----Table-------------------------------------------------------------------------------------------------------------------------------------------------------

// Table 普通表
type Table struct {
	entity any
	alias  string
}

// TableOf 传入一个 any 类型的实体，将其封装为 Table 类型
func TableOf(entity any) Table {
	return Table{
		entity: entity,
	}
}

func (t Table) C(name string) Column {
	return Column{
		name:  name,
		table: t,
	}
}
func (t Table) As(alias string) Table {
	return Table{
		entity: t.entity,
		alias:  alias,
	}
}

// Join
// 用法 A JOIN B
// 或 (A JOIN B) JOIN (C JOIN D) JOIN (SubQuery)
func (t Table) Join(right Table) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		typ:   "JOIN",
		right: right,
	}
}
func (t Table) LeftJoin(right Table) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		typ:   "LEFT JOIN",
		right: right,
	}
}
func (t Table) RightJoin(right Table) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		typ:   "RIGHT JOIN",
		right: right,
	}
}
func (t Table) tableAlias() string {
	return t.alias
}

//-----Join-------------------------------------------------------------------------------------------------------------------------------------------------------

// Join 查询
type Join struct {
	left TableReference
	// JOIN, LEFT JOIN, LEFT OUTER JOIN
	typ   string
	right TableReference

	on    []Predicate
	using []string
}

// tableAlias join 中没有 (A join B) as xxxx 这种用法
// 所以我们返回空字符串
func (j Join) tableAlias() string {
	return ""
}

func (j Join) Join(right Table) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		typ:   "JOIN",
		right: right,
	}
}
func (j Join) LeftJoin(right Table) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		typ:   "LEFT JOIN",
		right: right,
	}
}
func (j Join) RightJoin(right Table) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		typ:   "RIGHT JOIN",
		right: right,
	}
}

type JoinBuilder struct {
	left TableReference
	// JOIN, LEFT JOIN, LEFT OUTER JOIN
	typ   string
	right TableReference
}

func (jb *JoinBuilder) On(ps ...Predicate) Join {
	return Join{
		left:  jb.left,
		typ:   jb.typ,
		right: jb.right,
		on:    ps,
	}
}

func (jb *JoinBuilder) Using(cols ...string) Join {
	return Join{
		left:  jb.left,
		typ:   jb.typ,
		right: jb.right,
		using: cols,
	}
}

//-----SubQuery-------------------------------------------------------------------------------------------------------------------------------------------------------

// SubQuery 子查询
type SubQuery struct {
	q     QueryBuilder
	alias string
}

func (s SubQuery) expr() {}

func (s SubQuery) tableAlias() string {
	return s.alias
}
