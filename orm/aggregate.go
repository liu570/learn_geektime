package orm

// Aggregate 实现聚合函数的结构体
type Aggregate struct {
	arg string
	fn  string
}

func (a Aggregate) selectable() {}

func (a Aggregate) expr() {}

func (a Aggregate) EQ(val any) Predicate {
	return Predicate{
		left:  a,
		op:    opEQ,
		right: NewValue(val),
	}
}

func (a Aggregate) LT(val any) Predicate {
	return Predicate{
		left:  a,
		op:    opLT,
		right: NewValue(val),
	}
}

func (a Aggregate) GT(val any) Predicate {
	return Predicate{
		left:  a,
		op:    opGT,
		right: NewValue(val),
	}
}

func Avg(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  "AVG",
	}
}

func Min(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  "MIN",
	}
}

func Max(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  "MAX",
	}
}
