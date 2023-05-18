package valuer

import (
	"database/sql"
	"learn_geektime/orm/model"
)

// 我们提供两种实现方法 一种 反射 ， 一种 unsafe 所以我们需要提供 接口选择使用哪种方法

// 先来一种反射和 unsafe 的抽象

type Value interface {
	//Field(name string) (any, error)

	// SetColumn 将数据库查询到的结果映射到结构体中
	SetColumn(rows *sql.Rows) error

	// Field 用于获取该名字对应的字段
	Field(name string) (any, error)
}

// Creator 本质上也是Factory 模式 ， 只不过其十分简单
type Creator func(val any, meta *model.Model) Value
