package model

import (
	"database/sql"
	"learn_geektime/orm/internal/errs"
	"reflect"
)

// 我们支持的全部标签上的 key 都放在这里
// 方便用户查找，和我们后期维护
const (
	tagKeyColumn = "column"
)

type ModelOpt func(m *Model) error

type Model struct {
	TableName string
	// 使用切片保持结构体字段的顺序
	Fields    []*Field
	FieldMap  map[string]*Field
	ColumnMap map[string]*Field
}

func WithTableName(name string) ModelOpt {
	return func(m *Model) error {
		if name == "" {
			return errs.ErrEmptyTableName
		}
		m.TableName = name
		return nil
	}
}

func WithColumnName(field string, colName string) ModelOpt {
	return func(m *Model) error {
		fd, ok := m.FieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.ColName = colName
		return nil
	}
}

func WithColumn(field string, col *Field) ModelOpt {
	return func(m *Model) error {
		m.FieldMap[field] = col
		return nil
	}
}

type Field struct {
	//对应的 go 里面的字段名
	GoName string
	//字段对应的列名
	ColName string
	//对应的结构体类型
	Type reflect.Type

	// uintptr 可以被gc管理地址的
	// uintptr 表达字段的相对量(偏移量)
	Offset uintptr

	// 标记该字段在结构体中的顺序
	Index int
	//AutoIncrement bool
}

// TableName 接口定义表名
type TableName interface {
	TableName() string
}

// TestModel 用于测试
type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
