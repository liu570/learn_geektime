package orm

const (
	fnAVG   = "AVG"
	fnSUM   = "SUM"
	fnMAX   = "MAX"
	fnMIN   = "MIN"
	fnCOUNT = "COUNT"
)

// Aggregate 实现聚合函数的结构体
// eg. 常见的聚合函数有 AVG("age")、SUM("age")、MAX("age")、MIN("age")、COUNT("age")、
type Aggregate struct {
	arg   string
	fn    string
	alias string
}

func (a Aggregate) selectable() {}

func (a Aggregate) expr() {}

// AS （列名的别名功能）这里新建一个对象 相当于将 Aggregate 设置为一个不可变对象 可以避免并发问题
// 同时这里不使用指针也可以减少内存逃逸现象
func (a Aggregate) AS(alias string) Aggregate {
	return Aggregate{
		arg:   a.arg,
		fn:    a.fn,
		alias: alias,
	}
}

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
		fn:  fnAVG,
	}
}

func Min(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  fnMIN,
	}
}

func Max(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  fnMAX,
	}
}

func Count(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  fnCOUNT,
	}
}

func Sum(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  fnSUM,
	}
}
