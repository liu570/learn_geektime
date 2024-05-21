package ast

import (
	"github.com/stretchr/testify/require"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestPrintVisitor_Visit(t *testing.T) {
	type args struct {
		node ast.Node
	}
	tests := []struct {
		name     string
		filename string
		src      any
	}{
		{
			name:     "PrintVisitor",
			filename: "print_visitor.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileSet := token.NewFileSet()
			file, err := parser.ParseFile(fileSet, tt.filename, tt.src, parser.ParseComments)
			require.NoError(t, err)
			v := &PrintVisitor{}
			ast.Walk(v, file)
		})
	}
}
