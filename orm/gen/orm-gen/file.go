package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"reflect"
	"strings"
)

type SingleFileEntryVisitor struct {
	fileVisitor *FileVisitor
}

func NewSingleFileEntryVisitor() *SingleFileEntryVisitor {
	return &SingleFileEntryVisitor{}
}

func (s *SingleFileEntryVisitor) Get(operators ...string) *File {
	types := make([]Type, len(s.fileVisitor.Types))
	for i, visitor := range s.fileVisitor.Types {
		types[i] = Type{
			Name:   visitor.Name,
			Fields: visitor.Fields,
		}
	}
	var opts []string
	if len(operators) == 0 {
		opts = []string{"LT", "GT", "EQ"}
	}
	opts = append(opts, operators...)
	return &File{
		Package:   s.fileVisitor.Package,
		Imports:   s.fileVisitor.Imports,
		Types:     types,
		Operators: opts,
	}
}

func (s *SingleFileEntryVisitor) Visit(node ast.Node) (w ast.Visitor) {
	fn, ok := node.(*ast.File)
	if !ok {
		// 不是我们需要的代表文件的节点
		return s
	}
	s.fileVisitor = &FileVisitor{
		Package: fn.Name.String(),
	}
	return s.fileVisitor
}

type File struct {
	Package   string
	Imports   []string
	Types     []Type
	Operators []string
}
type FileVisitor struct {
	Package string
	Imports []string
	Types   []*TypeVisitor
}

func (f *FileVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch node := node.(type) {
	case *ast.ImportSpec:
		path := node.Path.Value
		if node.Name != nil && node.Name.String() != "" {
			path = node.Name.String() + " " + path
		}
		f.Imports = append(f.Imports, path)
	case *ast.TypeSpec:
		if node.Name.String() == "" {
			return f
		}
		v := &TypeVisitor{
			Name: node.Name.String(),
		}
		f.Types = append(f.Types, v)
		return v
	}
	return f
}

type Type struct {
	Name   string
	Fields []Field
}
type Field struct {
	Name         string
	Type         string
	IsComparable bool
}
type TypeVisitor struct {
	Name   string
	Fields []Field
}

func (t *TypeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch node := node.(type) {
	case *ast.StructType:
		for _, field := range node.Fields.List {
			typ := exprToString(field.Type)
			isComparable := isComparable(field.Type)
			for _, name := range field.Names {
				t.Fields = append(t.Fields, Field{
					Name:         name.String(),
					Type:         typ,
					IsComparable: isComparable,
				})
			}
		}
	}
	return t
}

// Function to convert ast.Expr to string
func exprToString(expr ast.Expr) string {
	var sb strings.Builder
	err := printer.Fprint(&sb, token.NewFileSet(), expr)
	if err != nil {
		log.Fatalf("Failed to convert expression to string: %v", err)
	}
	return sb.String()
}

// Helper function to determine if a type is comparable
func isComparable(expr ast.Expr) bool {
	// Convert ast.Expr to string representation of the type
	typeStr := exprToString(expr)
	// Parse the type string as Go code and get its reflect.Type
	src := "package main\ntype _ " + typeStr
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		log.Fatalf("Failed to parse type string: %v", err)
	}
	var typ reflect.Type
	ast.Inspect(file, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			typ = reflect.TypeOf(ts.Type)
			return false
		}
		return true
	})
	// Check if the type is comparable using reflect
	return typ.Comparable()
}

func exprToStringV1(expr ast.Expr) string {
	var typ string
	switch expr := expr.(type) {
	case *ast.Ident:
		typ = "string"
	case *ast.StarExpr:
		switch xt := expr.X.(type) {
		case *ast.Ident:
			typ = "*" + xt.String()
		case *ast.SelectorExpr:
			typ = "*" + xt.X.(*ast.Ident).Name + "." + xt.Sel.String()
		case *ast.ArrayType:
			typ = "*[]byte"
		}
	case *ast.ArrayType:
		typ = "[]byte"
	default:
		panic(fmt.Sprintf("expr type %T not supported", expr))
	}
	return typ
}

// Helper function to determine if a type is comparable
func isComparableV1(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
			"float32", "float64", "complex64", "complex128",
			"bool", "string":
			return true
		default:
			// Assume user-defined types are comparable
			return true
		}
	case *ast.ArrayType:
		// Arrays are comparable if their element type is comparable
		return isComparable(t.Elt)
	case *ast.StructType:
		// Structs are comparable if all their fields are comparable
		for _, field := range t.Fields.List {
			if !isComparable(field.Type) {
				return false
			}
		}
		return true
	case *ast.StarExpr, *ast.InterfaceType, *ast.MapType, *ast.FuncType, *ast.ChanType, *ast.SliceExpr:
		// Pointers, interfaces, maps, functions, channels, and slices are not comparable
		return false
	default:
		return false
	}
}
