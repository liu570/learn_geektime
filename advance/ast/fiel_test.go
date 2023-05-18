package ast

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestAst(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go",
		`
// annotation go through the source code and extra the annotation
// @author Deng Ming
/* @multiple first line
second line
*/
// @date 2022/04/02
package annotation

`, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	fv := &FileVisitor{
		ans: map[string]string{},
	}
	ast.Walk(fv, f)
	// 期望获取的数据
	res := map[string]string{
		"date":   "2022/04/02",
		"author": "Deng Ming",
	}
	assert.Equal(t, res, fv.ans)
}

/*




package ast

import(
	"fmt"
	"go/ast"
	"reflect"
)

type printVisitor struct{
}

func (t *printVisitor) Visit(node ast.Node) (w ast.Visitor){
	if node == nil{
		fmt.println(nil)
		return t
	}
	val := reflect.ValueOf(node)
	typ := reflect.TypeOf(node)
	for typ.Kind() == reflect.Ptr{
		typ = typ.Elem()
		val = val.Elem()
	}
	fmt.Printf("val: %+v,type: %s \n",val.Interface(),typ.Name())
	return t
}



*/
