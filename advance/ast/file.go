package ast

import (
	"go/ast"
	"strings"
)

type FileVisitor struct {
	ans   map[string]string
	types []*TypeSpecVisitor
}

func (f *FileVisitor) Visit(node ast.Node) (w ast.Visitor) {

	switch n := node.(type) {
	case *ast.File:
		if n.Doc == nil || len(n.Doc.List) == 0 {
			return f
		}

		for _, doc := range n.Doc.List {
			if !strings.HasPrefix(doc.Text, "// @") {
				continue
			}
			text := strings.TrimPrefix(doc.Text, "// @")
			segs := strings.SplitN(text, " ", 2)
			value := ""
			key := segs[0]
			if len(segs) > 1 {
				value = segs[1]
			}
			f.ans[key] = value
		}
	case *ast.TypeSpec:
		v := &TypeSpecVisitor{}
		f.types = append(f.types, v)
		return v
	default:
		return f
	}

	return f
}

type TypeSpecVisitor struct {
	ans map[string]string
}

func (f *TypeSpecVisitor) Visit(node ast.Node) (w ast.Visitor) {
	n, ok := node.(*ast.TypeSpec)
	if !ok {
		return f
	}
	for _, doc := range n.Doc.List {
		if !strings.HasPrefix(doc.Text, "// @") {
			continue
		}
		text := strings.TrimPrefix(doc.Text, "// @")
		segs := strings.SplitN(text, " ", 2)
		value := ""
		key := segs[0]
		if len(segs) > 1 {
			value = segs[1]
		}
		f.ans[key] = value
	}
	return f
}
