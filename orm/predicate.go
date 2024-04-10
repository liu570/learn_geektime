// 该结构体 用来 对 where 之后的语句进行结构化

package orm

type op string

const (
	opEQ    = "="
	opLT    = "<"
	opGT    = ">"
	opNOT   = "NOT"
	opAND   = "AND"
	opOR    = "OR"
	opADD   = "+"
	opSUB   = "-"
	opMulti = "*"
)

func (o op) String() string {
	return string(o)
}

// Ent 的做法
// eg: Eq("id",12)
//func Eq(column string, arg any) Predicate {
//	return Predicate{
//		Column: column,
//		Op:     "=",
//		Data:    arg,
//	}
//}

// Column 用来描述 数据库查询语句中的列名
type Column struct {
	// 用来确定该列是哪个表的列名
	table TableReference
	name  string
	alias string
}

func (c Column) OrderAble() {}

// assign 标记实现 assignable 接口
func (c Column) assign() {}

// selectable 标记实现 Selectable 接口
func (c Column) selectable() {}

// expr 标记实现 Expression 接口
func (c Column) expr() {}

func C(name string) Column {
	return Column{name: name}
}

// EQ  用法 C("id").EQ(12)
func (c Column) EQ(val any) Predicate {
	right, ok := val.(Expression)
	if !ok {
		right = NewValue(val)
	}
	return Predicate{
		left:  c,
		op:    opEQ,
		right: right,
	}
}

func (c Column) LT(val any) Predicate {
	right, ok := val.(Expression)
	if !ok {
		right = NewValue(val)
	}
	return Predicate{
		left:  c,
		op:    opLT,
		right: right,
	}
}

func (c Column) GT(val any) Predicate {
	right, ok := val.(Expression)
	if !ok {
		right = NewValue(val)
	}
	return Predicate{
		left:  c,
		op:    opGT,
		right: right,
	}
}

func (c Column) Add(val any) Predicate {
	right, ok := val.(Expression)
	if !ok {
		right = NewValue(val)
	}
	return Predicate{
		left:  c,
		op:    opADD,
		right: right,
	}
}

func (c Column) Sub(val any) Predicate {
	right, ok := val.(Expression)
	if !ok {
		right = NewValue(val)
	}
	return Predicate{
		left:  c,
		op:    opSUB,
		right: right,
	}
}

// Predicate 代表一个查询条件
// Predicate 可以通过 Predicate 组合构成一个复杂的查询条件
type Predicate struct {
	left  Expression
	op    op
	right Expression
}

func (p Predicate) expr() {}

func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNOT,
		right: p,
	}
}

// And( C("id").EQ(12),C("name").Eq("Tom") )
//func And(p1 Predicate, p2 Predicate) Predicate {
//	return Predicate{
//		left:  p1,
//		op:    opAND,
//		right: p2,
//	}
//}

// C("id").EQ(12).And(C("name").Eq("Tom"))
func (p1 Predicate) And(p2 Predicate) Predicate {
	return Predicate{
		left:  p1,
		op:    opAND,
		right: p2,
	}
}

func (p1 Predicate) Or(p2 Predicate) Predicate {
	return Predicate{
		left:  p1,
		op:    opOR,
		right: p2,
	}
}

// Value 用于代表 SQL 语句参数的值
type Value struct {
	val any
}

func (v Value) expr() {}

func NewValue(val any) Value {
	return Value{val: val}
}
