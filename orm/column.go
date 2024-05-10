package orm

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
	return Predicate{
		left:  c,
		op:    opEQ,
		right: valueOf(val),
	}
}

func (c Column) LT(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opLT,
		right: valueOf(val),
	}
}

func (c Column) GT(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opGT,
		right: valueOf(val),
	}
}

func (c Column) Add(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opADD,
		right: valueOf(val),
	}
}

func (c Column) Sub(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opSUB,
		right: valueOf(val),
	}
}
