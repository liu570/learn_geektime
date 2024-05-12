package orm

type OrderAble interface {
	OrderAble()
}

// OrderBy 用于封装在 ORDER BY 下后的语句结构体
type OrderBy struct {
	col string
	// 正常来说有 ASC 与 DESC
	order string
}

func (o OrderBy) OrderAble() {}

func (o OrderBy) selectable() {}

func Asc(name string) OrderBy {
	return OrderBy{
		col:   name,
		order: "ASC",
	}
}

func Desc(name string) OrderBy {
	return OrderBy{
		col:   name,
		order: "DESC",
	}
}
