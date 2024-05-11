package orm

// Assignable 用来框定可以使用在 UPDATE 下的赋值数据类型
type Assignable interface {
	//assign 标记实现 assignable 接口
	assign()
}

// Assignment 构建用于实现 Assignable 接口
type Assignment struct {
	Column string
	val    Expression
}

func (a Assignment) assign() {}

// Assign 用于生成一个 Assignment 结构体
func Assign(column string, val any) Assignment {
	v, ok := val.(Expression)
	if !ok {
		v = Value{val: val}
	}
	return Assignment{
		Column: column,
		//val:    NewValue(val),
		val: v,
	}

}
