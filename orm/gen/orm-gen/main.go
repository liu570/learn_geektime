package main

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	src := os.Args[1]
	dstDir := filepath.Dir(src)
	fileName := filepath.Base(src)
	idx := strings.LastIndexByte(fileName, '.')
	dst := filepath.Join(dstDir, fileName[:idx]) + ".gen.go"
	f, err := os.Create(dst)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gen(f, src)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("generated file:", dst)
}

//go:embed tpl.gohtml
var genOrm string

func gen(w io.Writer, srcFile string) error {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, srcFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	s := NewSingleFileEntryVisitor()
	ast.Walk(s, file)
	f := s.Get()

	// 模板编程
	tpl := template.New("gen-orm")
	tpl, err = tpl.Parse(genOrm)
	if err != nil {
		return err
	}
	return tpl.Execute(w, f)
}
