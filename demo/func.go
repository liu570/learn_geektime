package demo

import (
	"errors"
	"reflect"
)

// IterateFuncs 输出方法信息，并执行调用
// 考虑可能输入的类型：nil,基本类型，内置类型（切片，map，channel）之类的
// 结构体，结构体指针，多级指针

// 决策支持结构体，或者结构体指针
func IterateFuncs(val any) (map[string]*FuncInfo, error) {
	if val == nil {
		return nil, errors.New("输入 nil")
	}
	typ := reflect.TypeOf(val)
	refVal := reflect.ValueOf(val)

	// 由于有许多结构体方法是定义在指针之上，所以在这里不可以单纯的将指针转化成为结构体
	// 访问字段的时候一定要将指针转化为结构体因为指针没有字段
	// 但是在这里指针是有方法的
	//if typ.Kind() == reflect.Ptr {
	//	typ = typ.Elem()
	//	refVal = refVal.Elem()
	//}

	//这里判断条件表示为当类型不会结构体且不为一级指针的时候返回error
	if typ.Kind() != reflect.Struct &&
		!(typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct) {
		return nil, errors.New("不支持类型")
	}

	numMethod := typ.NumMethod()
	res := make(map[string]*FuncInfo, numMethod)
	for i := 0; i < numMethod; i++ {
		method := typ.Method(i)
		//
		mt := method.Type

		numIn := mt.NumIn()
		// 构造返回值里面的int
		in := make([]reflect.Type, 0, numIn)
		for j := 0; j < numIn; j++ {
			in = append(in, mt.In(j))
		}

		numOut := mt.NumOut()
		out := make([]reflect.Type, 0, numOut)
		for j := 0; j < numOut; j++ {
			out = append(out, mt.Out(j))
		}

		callRes := method.Func.Call([]reflect.Value{refVal})
		retVals := make([]any, 0, len(callRes))
		for _, cr := range callRes {
			retVals = append(retVals, cr.Interface())
		}
		res[method.Name] = &FuncInfo{
			Name:   method.Name,
			In:     in,
			Out:    out,
			Result: retVals,
		}
	}
	return res, nil
}

type FuncInfo struct {
	Name string
	In   []reflect.Type
	Out  []reflect.Type

	// 反射调用得到的结果
	Result []any
}
