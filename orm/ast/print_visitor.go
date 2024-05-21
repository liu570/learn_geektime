package ast

import (
	"fmt"
	"go/ast"
	"reflect"
)

type PrintVisitor struct {
}

func (p PrintVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		fmt.Println("nil ast.Node")
		return p
	}
	val := reflect.ValueOf(node)
	for val.Kind() != reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	fmt.Printf("val: %v type: %v\n", val.Interface(), typ.Name())
	return p

}
