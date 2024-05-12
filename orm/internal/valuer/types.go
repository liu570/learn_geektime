package valuer

import (
	"database/sql"
	"learn_geektime/orm/model"
)

// 我们提供两种实现方法 一种 反射 ， 一种 unsafe 所以我们需要提供 接口选择使用哪种方法

// Value 结果集映射接口
// 将 DB 查询的结果映射到 go 中对应的结构体中去
// -- SetColumn 具体的映射方法
// -- Field 用于获取该名字对应的字段具体数据
type Value interface {
	//Field(name string) (any, error)

	// SetColumn 将数据库查询到的结果映射到结构体中
	SetColumn(rows *sql.Rows) error

	// Field 用于该字段名对应的字段值
	// 诞生背景，插入语句的时候我们需要根据给出的数据 将数据 加入 builder 中的 args 数据中去，因此引入该方法抽象
	// 目前有反射与unsafe两种实现
	Field(name string) (any, error)
}

// Creator 用于创建 Value 接口对象 本框架 提供两个 Creator 的实现
// -- NewReflectValue
// -- NewUnsafeValue
// 本质上也是 Factory 模式，只不过其十分简单
// 需要传入两个参数 一个是对象，一个是该对象对应的元数据
type Creator func(val any, meta *model.Model) Value
