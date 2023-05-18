package demo

import "go/ast"

type printVisitor struct {
}

func (p *printVisitor) Visit(node ast.Node) (w ast.Visitor) {
	return p
}
