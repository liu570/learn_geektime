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
	opIN    = "IN"
	opNOTIN = "NOT IN"
	opEXIST = "EXIST"
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

func Exist(val ...any) Predicate {
	return Predicate{
		op:    opEXIST,
		right: valueOf(val),
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

func valueOf(val any) Expression {
	switch v := val.(type) {
	case Expression:
		return v
	case []any:
		return IN{
			values: v,
		}
	default:
		return NewValue(val)
	}
}
