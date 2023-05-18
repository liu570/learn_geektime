package demo

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestAstPrinter(t *testing.T) {
	// 此地token代表代码里面最小的有意义的组件
	// eg import、 ( 、 func 、 TestAstPrinter 等都是一个token
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go",
		`
package ast

import(
	"fmt"
	"go/ast"
	"reflect"
)

type printVisitor struct{
	name string
}

func (t *printVisitor) Visit(node ast.Node) (w ast.Visitor){
	return t
}

`, parser.ParseComments)
	if err != nil {
		t.Fatal()
	}
	ast.Walk(&printVisitor{}, f)

}
