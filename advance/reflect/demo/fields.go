package demo

import (
	"errors"
	"fmt"
	"reflect"
)

func IterateFields(val any) {
	// 复杂逻辑
	res, err := iterateFields(val)

	//简单逻辑
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range res {
		fmt.Println(k, v)
	}

}

func iterateFields(val any) (map[string]any, error) {
	if val == nil {
		return nil, errors.New("不能为 nil")
	}

	typ := reflect.TypeOf(val)
	refVal := reflect.ValueOf(val)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		refVal = refVal.Elem()
	}
	numFields := typ.NumField()
	res := make(map[string]any, numFields)
	for i := 0; i < numFields; i++ {
		res[typ.Field(i).Name] = refVal.Field(i).Interface()
	}
	return res, nil

}
