package orm

// Expression 代表语句，或者是语句的一部分
// expr() 表明该接口为一个标记接口 并不具备实际的含义
type Expression interface {
	expr()
}

// RawExpr 代表原生表达式 用于提供用户自定义 SQL 语句
// 注意：orm 框架不对原生表达式进行任何校验
type RawExpr struct {
	raw  string
	args []any
}

func (r RawExpr) expr() {}

func (r RawExpr) selectable() {}

func Raw(expr string, args ...any) RawExpr {
	return RawExpr{
		raw:  expr,
		args: args,
	}
}

// s.Where(Raw("a = ? and b = ? and c = ? " , 1 ,2 ,3))
func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}

type SubQueryExpr struct {
	s SubQuery
	// 谓词 ALL ANY SOME
	pred string
}

func (s SubQueryExpr) expr() {}

func Any(query SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    query,
		pred: "ANY",
	}
}
func All(query SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    query,
		pred: "ALL",
	}
}
func Some(query SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    query,
		pred: "SOME",
	}
}
