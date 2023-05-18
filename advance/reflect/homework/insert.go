package homework

import (
	"errors"
	"reflect"
	"strings"
)

var errInvalidEntity = errors.New("invalid entiity")

// InsertStmt 作业里面我们只是生成SQL，所以在处理sql.NullString 之类的接口
// 只需要判断有没有实现 driver.Valuer 就可以了

func InsertStmt(entity interface{}) (string, []interface{}, error) {
	//没有传入结构体返回错误
	if entity == nil {
		return "", nil, errInvalidEntity
	}
	val := reflect.ValueOf(entity)
	typ := reflect.TypeOf(entity)
	// 若传入的是指针则根据指针找到结构体
	for typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = val.Type()
	}
	// 无字段属性的结构体 返回非法输入
	numFields := val.NumField()
	if numFields == 0 {
		return "", nil, errInvalidEntity
	}

	var sql strings.Builder
	var args []interface{}

	fields, values := getFieldsAndValues(val)

	sql.WriteString("INSERT INTO `")
	sql.WriteString(typ.Name())
	sql.WriteString("`(")
	for i := 0; i < len(fields); i++ {
		if i != 0 {
			sql.WriteString(",")
		}
		sql.WriteString("`" + fields[i] + "`")
	}
	sql.WriteString(") VALUES(")
	for i := 0; i < len(fields); i++ {
		if i != 0 {
			sql.WriteString(",")
		}
		sql.WriteRune('?')
		args = append(args, values[fields[i]])
	}
	sql.WriteString(");")

	return sql.String(), args, nil
}

func getFieldsAndValues(val reflect.Value) ([]string, map[string]any) {
	typ := val.Type()
	fieldNum := typ.NumField()
	fields := make([]string, 0, fieldNum)
	// 一大神坑，不可以将空值放进一个空 map 里所以需要提前make好map
	// assignment to entry in nil map
	values := make(map[string]interface{}, fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		// 既是结构体也是匿名字段
		if fieldVal.Kind() == reflect.Struct && field.Anonymous {
			subFields, subValues := getFieldsAndValues(fieldVal)
			for _, k := range subFields {
				if _, ok := values[k]; ok {
					continue
				}
				fields = append(fields, k)
				values[k] = subValues[k]

			}
		} else {
			fields = append(fields, field.Name)
			values[field.Name] = fieldVal.Interface()
		}
	}
	return fields, values
}
