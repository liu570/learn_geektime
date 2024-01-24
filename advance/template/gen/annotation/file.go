package annotation

import (
	"go/ast"
)

// SingleFileEntryVisitor 这部分和课堂演示差不多，但是我建议你们自己试着写一些
type SingleFileEntryVisitor struct {
	//panic("implement me")
	file *fileVisitor
}

func (s *SingleFileEntryVisitor) Get() File {
	//panic("implement me")
	return s.file.Get()
}

func (s *SingleFileEntryVisitor) Visit(node ast.Node) ast.Visitor {

	switch n := node.(type) {
	case *ast.File:
		fv := &fileVisitor{
			ans: newAnnotations[*ast.File](n, n.Doc),
		}
		s.file = fv
	case *ast.TypeSpec:

	}
	return s
}

// -----------------------------------------------------------------------------------------------------

type fileVisitor struct {
	ans     Annotations[*ast.File]
	types   []*typeVisitor
	visited bool
}

func (f *fileVisitor) Get() File {
	//panic("implement me")
	return File{
		Annotations: f.ans,
	}
}

func (f *fileVisitor) Visit(node ast.Node) ast.Visitor {
	//panic("implement me")
	n := node.(*ast.File)
	f.ans = newAnnotations[*ast.File](n, n.Doc)
	return f
}

// -----------------------------------------------------------------------------------------------------

type File struct {
	Annotations[*ast.File]
	Types []Type
}

// -----------------------------------------------------------------------------------------------------

type typeVisitor struct {
	ans    Annotations[*ast.TypeSpec]
	fields []Field
}

func (t *typeVisitor) Get() Type {
	panic("implement me")
}

func (t *typeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	panic("implement me")
}

// -----------------------------------------------------------------------------------------------------

// Type ast中表明结构体类型
type Type struct {
	Annotations[*ast.TypeSpec]
	Fields []Field
}

// -----------------------------------------------------------------------------------------------------

// Field 字段结构体
type Field struct {
	Annotations[*ast.Field]
}
