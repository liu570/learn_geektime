package demo

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRefelct(t *testing.T) {
	u := &BaseEntity{}
	val := reflect.ValueOf(u)
	typ := val.Type()
	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = val.Type()
	}
	fmt.Println("typ.Name() : " + typ.Name())
	fmt.Println("typ.String() : " + typ.String())
	fmt.Printf("typ.NumField() : %d\n", typ.NumField())
	fmt.Printf("val.Type() : %v\n", val.Type())
	fmt.Printf("val.Kind() : %v\n", val.Kind())
	fmt.Printf("val.Interface() : %v\n", val.Interface())
	var args []interface{}
	for i := 0; i < typ.NumField(); i++ {
		fmt.Printf("typ.Field(%d).Name : %v\n", i, typ.Field(i).Name)
		fmt.Printf("typ.Field(%d).Type : %v\n", i, typ.Field(i).Type)
		fmt.Printf("typ.Field(%d).Tag : %v\n", i, typ.Field(i).Tag)
		if val.Field(i).Interface() == nil {
			println("yyyyyyyyyyyyyyyyyyyy")
		}
		fmt.Printf("----------------------\n")
		fmt.Printf("val.Field(%d).Interface() : %v\n", i, val.Field(i).Interface())
		args = append(args, val.Field(i).Interface())
		fmt.Printf("val.Field(%d).String() : %s\n", i, val.Field(i).String())
		fmt.Printf("val.Field(%d).Type() : %v\n", i, val.Field(i).Type())
		fmt.Printf("val.Field(%d).Kind() : %v\n", i, val.Field(i).Kind())
		fmt.Printf("----------------------\n")
	}

}

func TestSlice(t *testing.T) {
	s := []string{"123", "456"}
	s1 := []string{"789", "876"}
	s = append(s, "s1")
	fmt.Println(s)
	fmt.Println(s1)
}

func TestRefelct2(t *testing.T) {
	use := &User2{
		Name:     "liu",
		password: "123456",
		Age:      22,
	}
	val := reflect.ValueOf(use)
	typ := val.Type()
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = val.Type()
	}
	fields := make([]string, 0, val.NumField())
	values := make(map[string]any, val.NumField())
	for i := 0; i < typ.NumField(); i++ {

		//fmt.Println("---------------", i, "  ", val.Field(i).Interface())
		fields = append(fields, typ.Field(i).Name)
		values[typ.Field(i).Name] = val.FieldByName(typ.Field(i).Name).Interface()
		fmt.Println(fields)
		fmt.Println(values)
	}

}

type User2 struct {
	Name     string
	Age      int
	password string
}

type BaseEntity struct {
	CreateTime int64
	UpdateTime *int64
}
